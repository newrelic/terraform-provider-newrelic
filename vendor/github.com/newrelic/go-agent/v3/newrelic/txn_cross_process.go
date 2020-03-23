package newrelic

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/newrelic/go-agent/v3/internal"
	"github.com/newrelic/go-agent/v3/internal/cat"
)

// Bitfield values for the txnCrossProcess.Type field.
const (
	txnCrossProcessSynthetics = (1 << 0)
	txnCrossProcessInbound    = (1 << 1)
	txnCrossProcessOutbound   = (1 << 2)
)

var (
	// errAccountNotTrusted indicates that, while the inbound headers were valid,
	// the account ID within them is not trusted by the user's application.
	errAccountNotTrusted = errors.New("account not trusted")
)

// txnCrossProcess contains the metadata required for CAT and Synthetics
// headers, transaction events, and traces.
type txnCrossProcess struct {
	// The user side switch controlling whether CAT is enabled or not.
	Enabled bool

	// The user side switch controlling whether Distributed Tracing is enabled or not
	// This is required by synthetics support.  If Distributed Tracing is enabled,
	// any synthetics functionality that is triggered should not set nr.guid.
	DistributedTracingEnabled bool

	// Rather than copying in the entire ConnectReply, here are the fields that
	// we need to support CAT.
	CrossProcessID  []byte
	EncodingKey     []byte
	TrustedAccounts internal.TrustedAccountSet

	// CAT state for a given transaction.
	Type                uint8
	ClientID            string
	GUID                string
	TripID              string
	PathHash            string
	AlternatePathHashes map[string]bool
	ReferringPathHash   string
	ReferringTxnGUID    string
	Synthetics          *cat.SyntheticsHeader

	// The encoded synthetics header received as part of the request headers, if
	// any. By storing this here, we avoid needing to marshal the invariant
	// Synthetics struct above each time an external segment is created.
	SyntheticsHeader string
}

// crossProcessMetadata represents the metadata that must be transmitted with
// an external request for CAT to work.
type crossProcessMetadata struct {
	ID         string
	TxnData    string
	Synthetics string
}

// Init initialises a txnCrossProcess based on the given application connect
// reply.
func (txp *txnCrossProcess) Init(enabled bool, dt bool, reply *internal.ConnectReply) {
	txp.CrossProcessID = []byte(reply.CrossProcessID)
	txp.EncodingKey = []byte(reply.EncodingKey)
	txp.DistributedTracingEnabled = dt
	txp.Enabled = enabled
	txp.TrustedAccounts = reply.TrustedAccounts
}

// CreateCrossProcessMetadata generates request metadata that enable CAT and
// Synthetics support for an external segment.
func (txp *txnCrossProcess) CreateCrossProcessMetadata(txnName, appName string) (crossProcessMetadata, error) {
	metadata := crossProcessMetadata{}

	// Regardless of the user's CAT settings, if there was a synthetics header in
	// the inbound request, a synthetics header should always be included in the
	// outbound request headers.
	if txp.IsSynthetics() {
		metadata.Synthetics = txp.SyntheticsHeader
	}

	if txp.Enabled {
		txp.SetOutbound(true)
		txp.requireTripID()

		id, err := txp.outboundID()
		if err != nil {
			return metadata, err
		}

		txnData, err := txp.outboundTxnData(txnName, appName)
		if err != nil {
			return metadata, err
		}

		metadata.ID = id
		metadata.TxnData = txnData
	}

	return metadata, nil
}

// Finalise handles any end-of-transaction tasks. In practice, this simply
// means ensuring the path hash is set if it hasn't already been.
func (txp *txnCrossProcess) Finalise(txnName, appName string) error {
	if txp.Enabled && txp.Used() {
		_, err := txp.setPathHash(txnName, appName)
		return err
	}

	// If there was no CAT activity, then do nothing, successfully.
	return nil
}

// IsInbound returns true if the transaction had inbound CAT headers.
func (txp *txnCrossProcess) IsInbound() bool {
	return 0 != (txp.Type & txnCrossProcessInbound)
}

// IsOutbound returns true if the transaction has generated outbound CAT
// headers.
func (txp *txnCrossProcess) IsOutbound() bool {
	// We don't actually use this anywhere today, but it feels weird not having
	// it.
	return 0 != (txp.Type & txnCrossProcessOutbound)
}

// IsSynthetics returns true if the transaction had inbound Synthetics headers.
func (txp *txnCrossProcess) IsSynthetics() bool {
	// Technically, this is redundant: the presence of a non-nil Synthetics
	// pointer should be sufficient to determine if this is a synthetics
	// transaction. Nevertheless, it's convenient to have the Type field be
	// non-zero if any CAT behaviour has occurred.
	return 0 != (txp.Type&txnCrossProcessSynthetics) && nil != txp.Synthetics
}

// ParseAppData decodes the given appData value.
func (txp *txnCrossProcess) ParseAppData(encodedAppData string) (*cat.AppDataHeader, error) {
	if !txp.Enabled {
		return nil, nil
	}
	if encodedAppData != "" {
		rawAppData, err := deobfuscate(encodedAppData, txp.EncodingKey)
		if err != nil {
			return nil, err
		}

		appData := &cat.AppDataHeader{}
		if err := json.Unmarshal(rawAppData, appData); err != nil {
			return nil, err
		}

		return appData, nil
	}

	return nil, nil
}

// CreateAppData creates the appData value that should be sent with a response
// to ensure CAT operates as expected.
func (txp *txnCrossProcess) CreateAppData(name string, queueTime, responseTime time.Duration, contentLength int64) (string, error) {
	// If CAT is disabled, do nothing, successfully.
	if !txp.Enabled {
		return "", nil
	}

	data, err := json.Marshal(&cat.AppDataHeader{
		CrossProcessID:        string(txp.CrossProcessID),
		TransactionName:       name,
		QueueTimeInSeconds:    queueTime.Seconds(),
		ResponseTimeInSeconds: responseTime.Seconds(),
		ContentLength:         contentLength,
		TransactionGUID:       txp.GUID,
	})
	if err != nil {
		return "", err
	}

	obfuscated, err := obfuscate(data, txp.EncodingKey)
	if err != nil {
		return "", err
	}

	return obfuscated, nil
}

// Used returns true if any CAT or Synthetics related functionality has been
// triggered on the transaction.
func (txp *txnCrossProcess) Used() bool {
	return 0 != txp.Type
}

// SetInbound sets the inbound CAT flag. This function is provided only for
// internal and unit testing purposes, and should not be used outside of this
// package normally.
func (txp *txnCrossProcess) SetInbound(inbound bool) {
	if inbound {
		txp.Type |= txnCrossProcessInbound
	} else {
		txp.Type &^= txnCrossProcessInbound
	}
}

// SetOutbound sets the outbound CAT flag. This function is provided only for
// internal and unit testing purposes, and should not be used outside of this
// package normally.
func (txp *txnCrossProcess) SetOutbound(outbound bool) {
	if outbound {
		txp.Type |= txnCrossProcessOutbound
	} else {
		txp.Type &^= txnCrossProcessOutbound
	}
}

// SetSynthetics sets the Synthetics CAT flag. This function is provided only
// for internal and unit testing purposes, and should not be used outside of
// this package normally.
func (txp *txnCrossProcess) SetSynthetics(synthetics bool) {
	if synthetics {
		txp.Type |= txnCrossProcessSynthetics
	} else {
		txp.Type &^= txnCrossProcessSynthetics
	}
}

// handleInboundRequestHeaders parses the CAT headers from the given metadata
// and updates the relevant fields on the provided TxnData.
func (txp *txnCrossProcess) handleInboundRequestHeaders(metadata crossProcessMetadata) error {
	if txp.Enabled && metadata.ID != "" && metadata.TxnData != "" {
		if err := txp.handleInboundRequestEncodedCAT(metadata.ID, metadata.TxnData); err != nil {
			return err
		}
	}

	if metadata.Synthetics != "" {
		if err := txp.handleInboundRequestEncodedSynthetics(metadata.Synthetics); err != nil {
			return err
		}
	}

	return nil
}

func (txp *txnCrossProcess) handleInboundRequestEncodedCAT(encodedID, encodedTxnData string) error {
	rawID, err := deobfuscate(encodedID, txp.EncodingKey)
	if err != nil {
		return err
	}

	rawTxnData, err := deobfuscate(encodedTxnData, txp.EncodingKey)
	if err != nil {
		return err
	}

	if err := txp.handleInboundRequestID(rawID); err != nil {
		return err
	}

	return txp.handleInboundRequestTxnData(rawTxnData)
}

func (txp *txnCrossProcess) handleInboundRequestID(raw []byte) error {
	id, err := cat.NewIDHeader(raw)
	if err != nil {
		return err
	}

	if !txp.TrustedAccounts.IsTrusted(id.AccountID) {
		return errAccountNotTrusted
	}

	txp.SetInbound(true)
	txp.ClientID = string(raw)
	txp.setRequireGUID()

	return nil
}

func (txp *txnCrossProcess) handleInboundRequestTxnData(raw []byte) error {
	txnData := &cat.TxnDataHeader{}
	if err := json.Unmarshal(raw, txnData); err != nil {
		return err
	}

	txp.SetInbound(true)
	if txnData.TripID != "" {
		txp.TripID = txnData.TripID
	} else {
		txp.setRequireGUID()
		txp.TripID = txp.GUID
	}
	txp.ReferringTxnGUID = txnData.GUID
	txp.ReferringPathHash = txnData.PathHash

	return nil
}

func (txp *txnCrossProcess) handleInboundRequestEncodedSynthetics(encoded string) error {
	raw, err := deobfuscate(encoded, txp.EncodingKey)
	if err != nil {
		return err
	}

	if err := txp.handleInboundRequestSynthetics(raw); err != nil {
		return err
	}

	txp.SyntheticsHeader = encoded
	return nil
}

func (txp *txnCrossProcess) handleInboundRequestSynthetics(raw []byte) error {
	synthetics := &cat.SyntheticsHeader{}
	if err := json.Unmarshal(raw, synthetics); err != nil {
		return err
	}

	// The specced behaviour here if the account isn't trusted is to disable the
	// synthetics handling, but not CAT in general, so we won't return an error
	// here.
	if txp.TrustedAccounts.IsTrusted(synthetics.AccountID) {
		txp.SetSynthetics(true)
		txp.setRequireGUID()
		txp.Synthetics = synthetics
	}

	return nil
}

func (txp *txnCrossProcess) outboundID() (string, error) {
	return obfuscate(txp.CrossProcessID, txp.EncodingKey)
}

func (txp *txnCrossProcess) outboundTxnData(txnName, appName string) (string, error) {
	pathHash, err := txp.setPathHash(txnName, appName)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(&cat.TxnDataHeader{
		GUID:     txp.GUID,
		TripID:   txp.TripID,
		PathHash: pathHash,
	})
	if err != nil {
		return "", err
	}

	return obfuscate(data, txp.EncodingKey)
}

// setRequireGUID ensures that the transaction has a valid GUID, and sets the
// nr.guid and trip ID if they are not already set.  If the customer has enabled
// DistributedTracing, then the new style of guid will be set elsewhere.
func (txp *txnCrossProcess) setRequireGUID() {
	if txp.DistributedTracingEnabled {
		return
	}

	if txp.GUID != "" {
		return
	}

	txp.GUID = fmt.Sprintf("%x", randUint64())

	if txp.TripID == "" {
		txp.requireTripID()
	}
}

// requireTripID ensures that the transaction has a valid trip ID.
func (txp *txnCrossProcess) requireTripID() {
	if !txp.Enabled {
		return
	}
	if txp.TripID != "" {
		return
	}

	txp.setRequireGUID()
	txp.TripID = txp.GUID
}

// setPathHash generates a path hash, sets the transaction's path hash to
// match, and returns it. This function will also ensure that the alternate
// path hashes are correctly updated.
func (txp *txnCrossProcess) setPathHash(txnName, appName string) (string, error) {
	pathHash, err := cat.GeneratePathHash(txp.ReferringPathHash, txnName, appName)
	if err != nil {
		return "", err
	}

	if pathHash != txp.PathHash {
		if txp.PathHash != "" {
			// Lazily initialise the alternate path hashes if they haven't been
			// already.
			if txp.AlternatePathHashes == nil {
				txp.AlternatePathHashes = make(map[string]bool)
			}

			// The spec limits us to a maximum of 10 alternate path hashes.
			if len(txp.AlternatePathHashes) < 10 {
				txp.AlternatePathHashes[txp.PathHash] = true
			}
		}
		txp.PathHash = pathHash
	}

	return pathHash, nil
}
