package newrelic

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/newrelic/go-agent/v3/internal"
	"github.com/newrelic/go-agent/v3/internal/cat"
	"github.com/newrelic/go-agent/v3/internal/jsonx"
	"github.com/newrelic/go-agent/v3/internal/logger"
)

// txnEvent represents a transaction.
// https://source.datanerd.us/agents/agent-specs/blob/master/Transaction-Events-PORTED.md
// https://newrelic.atlassian.net/wiki/display/eng/Agent+Support+for+Synthetics%3A+Forced+Transaction+Traces+and+Analytic+Events
type txnEvent struct {
	FinalName          string
	Start              time.Time
	Duration           time.Duration
	TotalTime          time.Duration
	Queuing            time.Duration
	Zone               apdexZone
	Attrs              *attributes
	externalCallCount  uint64
	externalDuration   time.Duration
	datastoreCallCount uint64
	datastoreDuration  time.Duration
	CrossProcess       txnCrossProcess
	BetterCAT          betterCAT
	HasError           bool
}

// betterCAT stores the transaction's priority and all fields related
// to a DistributedTracer's Cross-Application Trace.
type betterCAT struct {
	Enabled       bool
	Priority      priority
	Sampled       bool
	Inbound       *payload
	TxnID         string
	TraceID       string
	TransportType string
}

// SetTraceAndTxnIDs takes a single 32 character ID and uses it to
// set both the trace (32 char) and transaction (16 char) ID.
func (bc *betterCAT) SetTraceAndTxnIDs(traceID string) {
	txnLength := 16
	bc.TraceID = traceID
	if len(traceID) <= txnLength {
		bc.TxnID = traceID
	} else {
		bc.TxnID = traceID[:txnLength]
	}
}

// txnData contains the recorded data of a transaction.
type txnData struct {
	txnEvent
	IsWeb          bool
	Name           string    // Work in progress name.
	Errors         txnErrors // Lazily initialized.
	Stop           time.Time
	ApdexThreshold time.Duration

	stamp           segmentStamp
	threadIDCounter uint64

	TraceIDGenerator        *internal.TraceIDGenerator
	ShouldCollectSpanEvents func() bool
	rootSpanID              string
	rootSpanErrData         *errorData
	SpanEvents              []*spanEvent

	customSegments    map[string]*metricData
	datastoreSegments map[datastoreMetricKey]*metricData
	externalSegments  map[externalMetricKey]*metricData
	messageSegments   map[internal.MessageMetricKey]*metricData

	TxnTrace txnTrace

	SlowQueriesEnabled bool
	SlowQueryThreshold time.Duration
	SlowQueries        *slowQueries

	// These better CAT supportability fields are left outside of
	// TxnEvent.BetterCAT to minimize the size of transaction event memory.
	DistributedTracingSupport distributedTracingSupport
}

func (t *txnData) saveTraceSegment(end segmentEnd, name string, attrs spanAttributeMap, externalGUID string) {
	attrs = t.Attrs.filterSpanAttributes(attrs, destSegment)
	t.TxnTrace.witnessNode(end, name, attrs, externalGUID)
}

// tracingThread contains a segment stack that is used to track segment parenting time
// within a single goroutine.
type tracingThread struct {
	threadID uint64
	stack    []segmentFrame
	// start and end are used to track the TotalTime this tracingThread was active.
	start time.Time
	end   time.Time
}

// RecordActivity indicates that activity happened at this time on this
// goroutine which helps track total time.
func (thread *tracingThread) RecordActivity(now time.Time) {
	if thread.start.IsZero() || now.Before(thread.start) {
		thread.start = now
	}
	if now.After(thread.end) {
		thread.end = now
	}
}

// TotalTime returns the amount to time that this thread contributes to the
// total time.
func (thread *tracingThread) TotalTime() time.Duration {
	if thread.start.Before(thread.end) {
		return thread.end.Sub(thread.start)
	}
	return 0
}

// newTracingThread returns a new tracingThread to track segments in a new goroutine.
func newTracingThread(txndata *txnData) *tracingThread {
	// Each thread needs a unique ID.
	txndata.threadIDCounter++
	return &tracingThread{
		threadID: txndata.threadIDCounter,
	}
}

type segmentStamp uint64

type segmentTime struct {
	Stamp segmentStamp
	Time  time.Time
}

// segmentStartTime is embedded into the top level segments (rather than
// segmentTime) to minimize the structure sizes to minimize allocations.
type segmentStartTime struct {
	Stamp segmentStamp
	Depth int
}

type stringJSONWriter string

func (s stringJSONWriter) WriteJSON(buf *bytes.Buffer) {
	jsonx.AppendString(buf, string(s))
}

type intJSONWriter int

func (i intJSONWriter) WriteJSON(buf *bytes.Buffer) {
	jsonx.AppendInt(buf, int64(i))
}

// spanAttributeMap is used for span attributes and segment attributes. The
// value is a jsonWriter to allow for segment query parameters.
type spanAttributeMap map[spanAttribute]jsonWriter

func (m *spanAttributeMap) addString(key spanAttribute, val string) {
	if "" != val {
		m.add(key, stringJSONWriter(val))
	}
}

func (m *spanAttributeMap) addInt(key spanAttribute, val int) {
	m.add(key, intJSONWriter(val))
}

func (m *spanAttributeMap) add(key spanAttribute, val jsonWriter) {
	if *m == nil {
		*m = make(spanAttributeMap)
	}
	(*m)[key] = val
}

func (m spanAttributeMap) copy() spanAttributeMap {
	if len(m) == 0 {
		return nil
	}
	cpy := make(spanAttributeMap, len(m))
	for k, v := range m {
		cpy[k] = v
	}
	return cpy
}

type segmentFrame struct {
	segmentTime
	children   time.Duration
	spanID     string
	attributes spanAttributeMap
}

type segmentEnd struct {
	start      segmentTime
	stop       segmentTime
	duration   time.Duration
	exclusive  time.Duration
	SpanID     string
	ParentID   string
	threadID   uint64
	attributes spanAttributeMap
}

func (end segmentEnd) spanEvent() *spanEvent {
	if "" == end.SpanID {
		return nil
	}
	return &spanEvent{
		GUID:         end.SpanID,
		ParentID:     end.ParentID,
		Timestamp:    end.start.Time,
		Duration:     end.duration,
		Attributes:   end.attributes,
		IsEntrypoint: false,
	}
}

const (
	datastoreProductUnknown   = "Unknown"
	datastoreOperationUnknown = "other"
)

// HasErrors indicates whether the transaction had errors.
func (t *txnData) HasErrors() bool {
	return len(t.Errors) > 0
}

func (t *txnData) time(now time.Time) segmentTime {
	// Update the stamp before using it so that a 0 stamp can be special.
	t.stamp++
	return segmentTime{
		Time:  now,
		Stamp: t.stamp,
	}
}

// AddAgentSpanAttribute allows attributes to be added to spans.
func (thread *tracingThread) AddAgentSpanAttribute(key spanAttribute, val string) {
	if len(thread.stack) > 0 {
		thread.stack[len(thread.stack)-1].attributes.addString(key, val)
	}
}

// RemoveErrorSpanAttribute allows attributes to be removed from spans.
func (thread *tracingThread) RemoveErrorSpanAttribute(key spanAttribute) {
	stackLen := len(thread.stack)
	if stackLen <= 0 {
		return
	}
	delete(thread.stack[stackLen-1].attributes, key)
}

// startSegment begins a segment.
func startSegment(t *txnData, thread *tracingThread, now time.Time) segmentStartTime {
	tm := t.time(now)
	thread.stack = append(thread.stack, segmentFrame{
		segmentTime: tm,
		children:    0,
	})

	return segmentStartTime{
		Stamp: tm.Stamp,
		Depth: len(thread.stack) - 1,
	}
}

// GetRootSpanID returns the root span ID.
func (t *txnData) GetRootSpanID() string {
	if "" == t.rootSpanID {
		t.rootSpanID = t.TraceIDGenerator.GenerateSpanID()
	}
	return t.rootSpanID
}

// CurrentSpanIdentifier returns the identifier of the span at the top of the
// segment stack.
func (t *txnData) CurrentSpanIdentifier(thread *tracingThread) string {
	if 0 == len(thread.stack) {
		return t.GetRootSpanID()
	}
	if "" == thread.stack[len(thread.stack)-1].spanID {
		thread.stack[len(thread.stack)-1].spanID = t.TraceIDGenerator.GenerateSpanID()
	}
	return thread.stack[len(thread.stack)-1].spanID
}

func (t *txnData) saveSpanEvent(e *spanEvent) {
	e.Attributes = t.Attrs.filterSpanAttributes(e.Attributes, destSpan)
	if len(t.SpanEvents) < maxSpanEvents {
		t.SpanEvents = append(t.SpanEvents, e)
	}
}

var (
	errMalformedSegment = errors.New("segment identifier malformed: perhaps unsafe code has modified it?")
	// errSegmentOrder indicates that segments have been ended in the
	// incorrect order.
	errSegmentOrder = errors.New(`improper segment use: segments must be ended in "last started first ended" order: ` +
		`use https://godoc.org/github.com/newrelic/go-agent/newrelic#Transaction.NewGoroutine to use the transaction in multiple goroutines`)
)

func endSegment(t *txnData, thread *tracingThread, start segmentStartTime, now time.Time) (segmentEnd, error) {
	if 0 == start.Stamp {
		return segmentEnd{}, errMalformedSegment
	}
	if start.Depth >= len(thread.stack) {
		return segmentEnd{}, errSegmentOrder
	}
	if start.Depth < 0 {
		return segmentEnd{}, errMalformedSegment
	}
	frame := thread.stack[start.Depth]
	if start.Stamp != frame.Stamp {
		return segmentEnd{}, errSegmentOrder
	}

	var children time.Duration
	for i := start.Depth; i < len(thread.stack); i++ {
		children += thread.stack[i].children
	}
	s := segmentEnd{
		stop:       t.time(now),
		start:      frame.segmentTime,
		attributes: frame.attributes,
	}
	if s.stop.Time.After(s.start.Time) {
		s.duration = s.stop.Time.Sub(s.start.Time)
	}
	if s.duration > children {
		s.exclusive = s.duration - children
	}

	// Note that we expect (depth == (len(t.stack) - 1)).  However, if
	// (depth < (len(t.stack) - 1)), that's ok: could be a panic popped
	// some stack frames (and the consumer was not using defer).

	if start.Depth > 0 {
		thread.stack[start.Depth-1].children += s.duration
	}

	thread.stack = thread.stack[0:start.Depth]

	if fn := t.ShouldCollectSpanEvents; nil != fn && fn() {
		s.SpanID = frame.spanID
		if "" == s.SpanID {
			s.SpanID = t.TraceIDGenerator.GenerateSpanID()
		}
		// Note that the current span identifier is the parent's
		// identifier because we've already popped the segment that's
		// ending off of the stack.
		s.ParentID = t.CurrentSpanIdentifier(thread)
	}

	s.threadID = thread.threadID

	thread.RecordActivity(s.start.Time)
	thread.RecordActivity(s.stop.Time)

	return s, nil
}

// endBasicSegment ends a basic segment.
func endBasicSegment(t *txnData, thread *tracingThread, start segmentStartTime, now time.Time, name string) error {
	end, err := endSegment(t, thread, start, now)
	if nil != err {
		return err
	}
	if nil == t.customSegments {
		t.customSegments = make(map[string]*metricData)
	}
	m := metricDataFromDuration(end.duration, end.exclusive)
	if data, ok := t.customSegments[name]; ok {
		data.aggregate(m)
	} else {
		// Use `new` in place of &m so that m is not
		// automatically moved to the heap.
		cpy := new(metricData)
		*cpy = m
		t.customSegments[name] = cpy
	}

	if t.TxnTrace.considerNode(end) {
		attributes := end.attributes.copy()
		t.saveTraceSegment(end, customSegmentMetric(name), attributes, "")
	}

	if evt := end.spanEvent(); evt != nil {
		evt.Name = customSegmentMetric(name)
		evt.Category = spanCategoryGeneric
		t.saveSpanEvent(evt)
	}

	return nil
}

// endExternalParams contains the parameters for endExternalSegment.
type endExternalParams struct {
	TxnData    *txnData
	Thread     *tracingThread
	Start      segmentStartTime
	Now        time.Time
	Logger     logger.Logger
	Response   *http.Response
	URL        *url.URL
	Host       string
	Library    string
	Method     string
	StatusCode *int
}

// endExternalSegment ends an external segment.
func endExternalSegment(p endExternalParams) error {
	t := p.TxnData
	end, err := endSegment(t, p.Thread, p.Start, p.Now)
	if nil != err {
		return err
	}

	// Use the Host field if present, otherwise use host in the URL.
	if p.Host == "" && p.URL != nil {
		p.Host = p.URL.Host
	}
	if p.Host == "" {
		p.Host = "unknown"
	}
	if p.Library == "" {
		p.Library = "http"
	}

	var appData *cat.AppDataHeader
	if p.Response != nil {
		hdr := httpHeaderToAppData(p.Response.Header)
		appData, err = t.CrossProcess.ParseAppData(hdr)
		if err != nil {
			if p.Logger.DebugEnabled() {
				p.Logger.Debug("failure to parse cross application response header", map[string]interface{}{
					"err":    err.Error(),
					"header": hdr,
				})
			}
		}
	}

	var crossProcessID string
	var transactionName string
	var transactionGUID string
	if appData != nil {
		crossProcessID = appData.CrossProcessID
		transactionName = appData.TransactionName
		transactionGUID = appData.TransactionGUID
	}

	key := externalMetricKey{
		Host:                    p.Host,
		Library:                 p.Library,
		Method:                  p.Method,
		ExternalCrossProcessID:  crossProcessID,
		ExternalTransactionName: transactionName,
	}
	if nil == t.externalSegments {
		t.externalSegments = make(map[externalMetricKey]*metricData)
	}
	t.externalCallCount++
	t.externalDuration += end.duration
	m := metricDataFromDuration(end.duration, end.exclusive)
	if data, ok := t.externalSegments[key]; ok {
		data.aggregate(m)
	} else {
		// Use `new` in place of &m so that m is not
		// automatically moved to the heap.
		cpy := new(metricData)
		*cpy = m
		t.externalSegments[key] = cpy
	}

	if t.TxnTrace.considerNode(end) {
		attributes := end.attributes.copy()
		if p.Library == "http" {
			attributes.addString(spanAttributeHTTPURL, safeURL(p.URL))
		}
		t.saveTraceSegment(end, key.scopedMetric(), attributes, transactionGUID)
	}

	if evt := end.spanEvent(); evt != nil {
		evt.Name = key.scopedMetric()
		evt.Category = spanCategoryHTTP
		evt.Kind = "client"
		evt.Component = p.Library
		if p.Library == "http" {
			evt.Attributes.addString(spanAttributeHTTPURL, safeURL(p.URL))
			evt.Attributes.addString(spanAttributeHTTPMethod, p.Method)
		}
		if p.StatusCode != nil {
			evt.Attributes.addInt(spanAttributeHTTPStatusCode, *p.StatusCode)
		} else if p.Response != nil {
			evt.Attributes.addInt(spanAttributeHTTPStatusCode, p.Response.StatusCode)
		}
		t.saveSpanEvent(evt)
	}

	return nil
}

// endMessageParams contains the parameters for endMessageSegment.
type endMessageParams struct {
	TxnData         *txnData
	Thread          *tracingThread
	Start           segmentStartTime
	Now             time.Time
	Logger          logger.Logger
	DestinationName string
	Library         string
	DestinationType string
	DestinationTemp bool
}

// endMessageSegment ends an external segment.
func endMessageSegment(p endMessageParams) error {
	t := p.TxnData
	end, err := endSegment(t, p.Thread, p.Start, p.Now)
	if nil != err {
		return err
	}

	key := internal.MessageMetricKey{
		Library:         p.Library,
		DestinationType: p.DestinationType,
		DestinationName: p.DestinationName,
		DestinationTemp: p.DestinationTemp,
	}

	if nil == t.messageSegments {
		t.messageSegments = make(map[internal.MessageMetricKey]*metricData)
	}
	m := metricDataFromDuration(end.duration, end.exclusive)
	if data, ok := t.messageSegments[key]; ok {
		data.aggregate(m)
	} else {
		// Use `new` in place of &m so that m is not
		// automatically moved to the heap.
		cpy := new(metricData)
		*cpy = m
		t.messageSegments[key] = cpy
	}

	if t.TxnTrace.considerNode(end) {
		attributes := end.attributes.copy()
		t.saveTraceSegment(end, key.Name(), attributes, "")
	}

	if evt := end.spanEvent(); evt != nil {
		evt.Name = key.Name()
		evt.Category = spanCategoryGeneric
		t.saveSpanEvent(evt)
	}

	return nil
}

// endDatastoreParams contains the parameters for endDatastoreSegment.
type endDatastoreParams struct {
	TxnData            *txnData
	Thread             *tracingThread
	Start              segmentStartTime
	Now                time.Time
	Product            string
	Collection         string
	Operation          string
	ParameterizedQuery string
	QueryParameters    map[string]interface{}
	Host               string
	PortPathOrID       string
	Database           string
	ThisHost           string
}

const (
	unknownDatastoreHost         = "unknown"
	unknownDatastorePortPathOrID = "unknown"
)

var (
	hostsToReplace = map[string]struct{}{
		"localhost":       {},
		"127.0.0.1":       {},
		"0.0.0.0":         {},
		"0:0:0:0:0:0:0:1": {},
		"::1":             {},
		"0:0:0:0:0:0:0:0": {},
		"::":              {},
	}
)

func (t txnData) slowQueryWorthy(d time.Duration) bool {
	return t.SlowQueriesEnabled && (d >= t.SlowQueryThreshold)
}

func datastoreSpanAddress(host, portPathOrID string) string {
	if "" != host && "" != portPathOrID {
		return host + ":" + portPathOrID
	}
	if "" != host {
		return host
	}
	return portPathOrID
}

// endDatastoreSegment ends a datastore segment.
func endDatastoreSegment(p endDatastoreParams) error {
	end, err := endSegment(p.TxnData, p.Thread, p.Start, p.Now)
	if nil != err {
		return err
	}
	if p.Operation == "" {
		p.Operation = datastoreOperationUnknown
	}
	if p.Product == "" {
		p.Product = datastoreProductUnknown
	}
	if p.Host == "" && p.PortPathOrID != "" {
		p.Host = unknownDatastoreHost
	}
	if p.PortPathOrID == "" && p.Host != "" {
		p.PortPathOrID = unknownDatastorePortPathOrID
	}
	if _, ok := hostsToReplace[p.Host]; ok {
		p.Host = p.ThisHost
	}

	// We still want to create a slowQuery if the consumer has not provided
	// a Query string (or it has been removed by LASP) since the stack trace
	// has value.
	if p.ParameterizedQuery == "" {
		collection := p.Collection
		if "" == collection {
			collection = "unknown"
		}
		p.ParameterizedQuery = fmt.Sprintf(`'%s' on '%s' using '%s'`,
			p.Operation, collection, p.Product)
	}

	key := datastoreMetricKey{
		Product:      p.Product,
		Collection:   p.Collection,
		Operation:    p.Operation,
		Host:         p.Host,
		PortPathOrID: p.PortPathOrID,
	}
	if nil == p.TxnData.datastoreSegments {
		p.TxnData.datastoreSegments = make(map[datastoreMetricKey]*metricData)
	}
	p.TxnData.datastoreCallCount++
	p.TxnData.datastoreDuration += end.duration
	m := metricDataFromDuration(end.duration, end.exclusive)
	if data, ok := p.TxnData.datastoreSegments[key]; ok {
		data.aggregate(m)
	} else {
		// Use `new` in place of &m so that m is not
		// automatically moved to the heap.
		cpy := new(metricData)
		*cpy = m
		p.TxnData.datastoreSegments[key] = cpy
	}

	scopedMetric := datastoreScopedMetric(key)
	// errors in QueryParameters must not stop the recording of the segment
	queryParams, err := vetQueryParameters(p.QueryParameters)

	if p.TxnData.TxnTrace.considerNode(end) {
		attributes := end.attributes.copy()
		attributes.addString(spanAttributeDBStatement, p.ParameterizedQuery)
		attributes.addString(spanAttributeDBInstance, p.Database)
		attributes.addString(spanAttributePeerAddress, datastoreSpanAddress(p.Host, p.PortPathOrID))
		attributes.addString(spanAttributePeerHostname, p.Host)
		if len(queryParams) > 0 {
			attributes.add(spanAttributeQueryParameters, queryParams)
		}
		p.TxnData.saveTraceSegment(end, scopedMetric, attributes, "")
	}

	if p.TxnData.slowQueryWorthy(end.duration) {
		if nil == p.TxnData.SlowQueries {
			p.TxnData.SlowQueries = newSlowQueries(maxTxnSlowQueries)
		}
		p.TxnData.SlowQueries.observeInstance(slowQueryInstance{
			Duration:           end.duration,
			DatastoreMetric:    scopedMetric,
			ParameterizedQuery: p.ParameterizedQuery,
			QueryParameters:    queryParams,
			Host:               p.Host,
			PortPathOrID:       p.PortPathOrID,
			DatabaseName:       p.Database,
			StackTrace:         getStackTrace(),
		})
	}

	if evt := end.spanEvent(); evt != nil {
		evt.Name = scopedMetric
		evt.Category = spanCategoryDatastore
		evt.Kind = "client"
		evt.Component = p.Product
		evt.Attributes.addString(spanAttributeDBStatement, p.ParameterizedQuery)
		evt.Attributes.addString(spanAttributeDBInstance, p.Database)
		evt.Attributes.addString(spanAttributePeerAddress, datastoreSpanAddress(p.Host, p.PortPathOrID))
		evt.Attributes.addString(spanAttributePeerHostname, p.Host)
		evt.Attributes.addString(spanAttributeDBCollection, p.Collection)
		p.TxnData.saveSpanEvent(evt)
	}

	return err
}

// MergeBreakdownMetrics creates segment metrics.
func mergeBreakdownMetrics(t *txnData, metrics *metricTable) {
	scope := t.FinalName
	isWeb := t.IsWeb
	// Custom Segment Metrics
	for key, data := range t.customSegments {
		name := customSegmentMetric(key)
		// Unscoped
		metrics.add(name, "", *data, unforced)
		// Scoped
		metrics.add(name, scope, *data, unforced)
	}

	// External Segment Metrics
	for key, data := range t.externalSegments {
		metrics.add(externalRollupMetric.all, "", *data, forced)
		metrics.add(externalRollupMetric.webOrOther(isWeb), "", *data, forced)

		hostMetric := externalHostMetric(key)
		metrics.add(hostMetric, "", *data, unforced)
		if "" != key.ExternalCrossProcessID && "" != key.ExternalTransactionName {
			txnMetric := externalTransactionMetric(key)

			// Unscoped CAT metrics
			metrics.add(externalAppMetric(key), "", *data, unforced)
			metrics.add(txnMetric, "", *data, unforced)
		}

		// Scoped External Metric
		metrics.add(key.scopedMetric(), scope, *data, unforced)
	}

	// Datastore Segment Metrics
	for key, data := range t.datastoreSegments {
		metrics.add(datastoreRollupMetric.all, "", *data, forced)
		metrics.add(datastoreRollupMetric.webOrOther(isWeb), "", *data, forced)

		product := datastoreProductMetric(key)
		metrics.add(product.all, "", *data, forced)
		metrics.add(product.webOrOther(isWeb), "", *data, forced)

		if key.Host != "" && key.PortPathOrID != "" {
			instance := datastoreInstanceMetric(key)
			metrics.add(instance, "", *data, unforced)
		}

		operation := datastoreOperationMetric(key)
		metrics.add(operation, "", *data, unforced)

		if "" != key.Collection {
			statement := datastoreStatementMetric(key)

			metrics.add(statement, "", *data, unforced)
			metrics.add(statement, scope, *data, unforced)
		} else {
			metrics.add(operation, scope, *data, unforced)
		}
	}
	// Message Segment Metrics
	for key, data := range t.messageSegments {
		metric := key.Name()
		metrics.add(metric, scope, *data, unforced)
		metrics.add(metric, "", *data, unforced)
	}
}
