package newrelic

import (
	"net/http"
	"net/url"
	"strings"
)

// Transaction instruments one logical unit of work: either an inbound web
// request or background task.  Start a new Transaction with the
// Application.StartTransaction method.
//
// All methods on Transaction are nil safe. Therefore, a nil Transaction
// pointer can be safely used as a mock.
type Transaction struct {
	Private interface{}
	thread  *thread
}

// End finishes the Transaction.  After that, subsequent calls to End or
// other Transaction methods have no effect.  All segments and
// instrumentation must be completed before End is called.
func (txn *Transaction) End() {
	if nil == txn {
		return
	}
	if nil == txn.thread {
		return
	}

	var r interface{}
	if txn.thread.Config.ErrorCollector.RecordPanics {
		// recover must be called in the function directly being deferred,
		// not any nested call!
		r = recover()
	}
	txn.thread.logAPIError(txn.thread.End(r), "end transaction", nil)
}

// Ignore prevents this transaction's data from being recorded.
func (txn *Transaction) Ignore() {
	if nil == txn {
		return
	}
	if nil == txn.thread {
		return
	}
	txn.thread.logAPIError(txn.thread.Ignore(), "ignore transaction", nil)
}

// SetName names the transaction.  Use a limited set of unique names to
// ensure that Transactions are grouped usefully.
func (txn *Transaction) SetName(name string) {
	if nil == txn {
		return
	}
	if nil == txn.thread {
		return
	}
	txn.thread.logAPIError(txn.thread.SetName(name), "set transaction name", nil)
}

// NoticeError records an error.  The Transaction saves the first five
// errors.  For more control over the recorded error fields, see the
// newrelic.Error type.
//
// In certain situations, using this method may result in an error being
// recorded twice.  Errors are automatically recorded when
// Transaction.WriteHeader receives a status code at or above 400 or strictly
// below 100 that is not in the IgnoreStatusCodes configuration list.  This
// method is unaffected by the IgnoreStatusCodes configuration list.
//
// NoticeError examines whether the error implements the following optional
// methods:
//
//   // StackTrace records a stack trace
//   StackTrace() []uintptr
//
//   // ErrorClass sets the error's class
//   ErrorClass() string
//
//   // ErrorAttributes sets the errors attributes
//   ErrorAttributes() map[string]interface{}
//
// The newrelic.Error type, which implements these methods, is the recommended
// way to directly control the recorded error's message, class, stacktrace,
// and attributes.
func (txn *Transaction) NoticeError(err error) {
	if nil == txn {
		return
	}
	if nil == txn.thread {
		return
	}
	txn.thread.logAPIError(txn.thread.NoticeError(err), "notice error", nil)
}

// AddAttribute adds a key value pair to the transaction event, errors,
// and traces.
//
// The key must contain fewer than than 255 bytes.  The value must be a
// number, string, or boolean.
//
// For more information, see:
// https://docs.newrelic.com/docs/agents/manage-apm-agents/agent-metrics/collect-custom-attributes
func (txn *Transaction) AddAttribute(key string, value interface{}) {
	if nil == txn {
		return
	}
	if nil == txn.thread {
		return
	}
	txn.thread.logAPIError(txn.thread.AddAttribute(key, value), "add attribute", nil)
}

// SetWebRequestHTTP marks the transaction as a web transaction.  If
// the request is non-nil, SetWebRequestHTTP will additionally collect
// details on request attributes, url, and method.  If headers are
// present, the agent will look for a distributed tracing header.
func (txn *Transaction) SetWebRequestHTTP(r *http.Request) {
	if nil == r {
		txn.SetWebRequest(WebRequest{})
		return
	}
	wr := WebRequest{
		Header:    r.Header,
		URL:       r.URL,
		Method:    r.Method,
		Transport: transport(r),
	}
	txn.SetWebRequest(wr)
}

func transport(r *http.Request) TransportType {
	if strings.HasPrefix(r.Proto, "HTTP") {
		if r.TLS != nil {
			return TransportHTTPS
		}
		return TransportHTTP
	}
	return TransportUnknown
}

// SetWebRequest marks the transaction as a web transaction.  SetWebRequest
// additionally collects details on request attributes, url, and method if these
// fields are set.  If headers are present, the agent will look for a
// distributed tracing header.  Use SetWebRequestHTTP if you have a
// *http.Request.
func (txn *Transaction) SetWebRequest(r WebRequest) {
	if nil == txn {
		return
	}
	if nil == txn.thread {
		return
	}
	txn.thread.logAPIError(txn.thread.SetWebRequest(r), "set web request", nil)
}

// SetWebResponse allows the Transaction to instrument response code and
// response headers.  Use the return value of this method in place of the input
// parameter http.ResponseWriter in your instrumentation.
//
// The returned http.ResponseWriter is safe to use even if the Transaction
// receiver is nil or has already been ended.
//
// The returned http.ResponseWriter implements the combination of
// http.CloseNotifier, http.Flusher, http.Hijacker, and io.ReaderFrom
// implemented by the input http.ResponseWriter.
//
// This method is used by WrapHandle, WrapHandleFunc, and most integration
// package middlewares.  Therefore, you probably want to use this only if you
// are writing your own instrumentation middleware.
func (txn *Transaction) SetWebResponse(w http.ResponseWriter) http.ResponseWriter {
	if nil == txn {
		return w
	}
	if nil == txn.thread {
		return w
	}
	return txn.thread.SetWebResponse(w)
}

// StartSegmentNow starts timing a segment.  The SegmentStartTime returned can
// be used as the StartTime field in Segment, DatastoreSegment, or
// ExternalSegment.  The returned SegmentStartTime is safe to use even  when the
// Transaction receiver is nil.  In this case, the segment will have no effect.
func (txn *Transaction) StartSegmentNow() SegmentStartTime {
	if nil == txn {
		return SegmentStartTime{}
	}
	if nil == txn.thread {
		return SegmentStartTime{}
	}
	return txn.thread.StartSegmentNow()
}

// StartSegment makes it easy to instrument segments.  To time a function, do
// the following:
//
//	func timeMe(txn newrelic.Transaction) {
//		defer txn.StartSegment("timeMe").End()
//		// ... function code here ...
//	}
//
// To time a block of code, do the following:
//
//	segment := txn.StartSegment("myBlock")
//	// ... code you want to time here ...
//	segment.End()
func (txn *Transaction) StartSegment(name string) *Segment {
	return &Segment{
		StartTime: txn.StartSegmentNow(),
		Name:      name,
	}
}

// InsertDistributedTraceHeaders adds a Distributed Trace header used to link
// transactions.  InsertDistributedTraceHeaders should be called every
// time an outbound call is made since the payload contains a timestamp.
//
// StartExternalSegment calls InsertDistributedTraceHeaders, so you
// don't need to use it for outbound HTTP calls: Just use
// StartExternalSegment!
func (txn *Transaction) InsertDistributedTraceHeaders(hdrs http.Header) {
	if nil == txn {
		return
	}
	if nil == txn.thread {
		return
	}
	insertDistributedTraceHeaders(txn.thread, hdrs)
}

// AcceptDistributedTraceHeaders links transactions by accepting a
// distributed trace payload from another transaction.
//
// Application.StartTransaction calls this method automatically if a
// payload is present in the request headers.  Therefore, this method
// does not need to be used for typical HTTP transactions.
//
// AcceptDistributedTraceHeaders should be used as early in the
// transaction as possible.  It may not be called after a call to
// CreateDistributedTracePayload.
func (txn *Transaction) AcceptDistributedTraceHeaders(t TransportType, hdrs http.Header) {
	if nil == txn {
		return
	}
	if nil == txn.thread {
		return
	}
	txn.thread.logAPIError(txn.thread.AcceptDistributedTraceHeaders(t, hdrs), "accept trace payload", nil)
}

// Application returns the Application which started the transaction.
func (txn *Transaction) Application() *Application {
	if nil == txn {
		return nil
	}
	if nil == txn.thread {
		return nil
	}
	return txn.thread.Application()
}

// BrowserTimingHeader generates the JavaScript required to enable New
// Relic's Browser product.  This code should be placed into your pages
// as close to the top of the <head> element as possible, but after any
// position-sensitive <meta> tags (for example, X-UA-Compatible or
// charset information).
//
// This function freezes the transaction name: any calls to SetName()
// after BrowserTimingHeader() will be ignored.
//
// The *BrowserTimingHeader return value will be nil if browser
// monitoring is disabled, the application is not connected, or an error
// occurred.  It is safe to call the pointer's methods if it is nil.
func (txn *Transaction) BrowserTimingHeader() *BrowserTimingHeader {
	if nil == txn {
		return nil
	}
	if nil == txn.thread {
		return nil
	}
	b, err := txn.thread.BrowserTimingHeader()
	txn.thread.logAPIError(err, "create browser timing header", nil)
	return b
}

// NewGoroutine allows you to use the Transaction in multiple
// goroutines.
//
// Each goroutine must have its own Transaction reference returned by
// NewGoroutine.  You must call NewGoroutine to get a new Transaction
// reference every time you wish to pass the Transaction to another
// goroutine. It does not matter if you call this before or after the
// other goroutine has started.
//
// All Transaction methods can be used in any Transaction reference.
// The Transaction will end when End() is called in any goroutine.
//
// Example passing a new Transaction reference directly to another
// goroutine:
//
//	go func(txn newrelic.Transaction) {
//		defer txn.StartSegment("async").End()
//		time.Sleep(100 * time.Millisecond)
//	}(txn.NewGoroutine())
//
// Example passing a new Transaction reference on a channel to another
// goroutine:
//
//	ch := make(chan newrelic.Transaction)
//	go func() {
//		txn := <-ch
//		defer txn.StartSegment("async").End()
//		time.Sleep(100 * time.Millisecond)
//	}()
//	ch <- txn.NewGoroutine()
//
func (txn *Transaction) NewGoroutine() *Transaction {
	if nil == txn {
		return nil
	}
	if nil == txn.thread {
		return nil
	}
	return txn.thread.NewGoroutine()
}

// GetTraceMetadata returns distributed tracing identifiers.  Empty
// string identifiers are returned if the transaction has finished.
func (txn *Transaction) GetTraceMetadata() TraceMetadata {
	if nil == txn {
		return TraceMetadata{}
	}
	if nil == txn.thread {
		return TraceMetadata{}
	}
	return txn.thread.GetTraceMetadata()
}

// GetLinkingMetadata returns the fields needed to link data to a trace or
// entity.
func (txn *Transaction) GetLinkingMetadata() LinkingMetadata {
	if nil == txn {
		return LinkingMetadata{}
	}
	if nil == txn.thread {
		return LinkingMetadata{}
	}
	return txn.thread.GetLinkingMetadata()
}

// IsSampled indicates if the Transaction is sampled.  A sampled
// Transaction records a span event for each segment.  Distributed tracing
// must be enabled for transactions to be sampled.  False is returned if
// the Transaction has finished.
func (txn *Transaction) IsSampled() bool {
	if nil == txn {
		return false
	}
	if nil == txn.thread {
		return false
	}
	return txn.thread.IsSampled()
}

const (
	// DistributedTraceNewRelicHeader is the header used by New Relic agents
	// for automatic trace payload instrumentation.
	DistributedTraceNewRelicHeader = "Newrelic"
)

// TransportType is used in Transaction.AcceptDistributedTraceHeaders() to
// represent the type of connection that the trace payload was transported
// over.
type TransportType string

// TransportType names used across New Relic agents:
const (
	TransportUnknown TransportType = "Unknown"
	TransportHTTP    TransportType = "HTTP"
	TransportHTTPS   TransportType = "HTTPS"
	TransportKafka   TransportType = "Kafka"
	TransportJMS     TransportType = "JMS"
	TransportIronMQ  TransportType = "IronMQ"
	TransportAMQP    TransportType = "AMQP"
	TransportQueue   TransportType = "Queue"
	TransportOther   TransportType = "Other"
)

func (tt TransportType) toString() string {
	switch tt {
	case TransportHTTP, TransportHTTPS, TransportKafka, TransportJMS, TransportIronMQ, TransportAMQP,
		TransportQueue, TransportOther:
		return string(tt)
	default:
		return string(TransportUnknown)
	}
}

// WebRequest is used to provide request information to Transaction.SetWebRequest.
type WebRequest struct {
	// Header may be nil if you don't have any headers or don't want to
	// transform them to http.Header format.
	Header http.Header
	// URL may be nil if you don't have a URL or don't want to transform
	// it to *url.URL.
	URL *url.URL
	// Method is the request's method.
	Method string
	// If a distributed tracing header is found in the WebRequest.Header,
	// this TransportType will be used in the distributed tracing metrics.
	Transport TransportType
}

// LinkingMetadata is returned by Transaction.GetLinkingMetadata.  It contains
// identifiers needed to link data to a trace or entity.
type LinkingMetadata struct {
	// TraceID identifies the entire distributed trace.  This field is empty
	// if distributed tracing is disabled.
	TraceID string
	// SpanID identifies the currently active segment.  This field is empty
	// if distributed tracing is disabled or the transaction is not sampled.
	SpanID string
	// EntityName is the Application name as set on the Config.  If multiple
	// application names are specified in the Config, only the first is
	// returned.
	EntityName string
	// EntityType is the type of this entity and is always the string
	// "SERVICE".
	EntityType string
	// EntityGUID is the unique identifier for this entity.
	EntityGUID string
	// Hostname is the hostname this entity is running on.
	Hostname string
}

// TraceMetadata is returned by Transaction.GetTraceMetadata.  It contains
// distributed tracing identifiers.
type TraceMetadata struct {
	// TraceID identifies the entire distributed trace.  This field is empty
	// if distributed tracing is disabled.
	TraceID string
	// SpanID identifies the currently active segment.  This field is empty
	// if distributed tracing is disabled or the transaction is not sampled.
	SpanID string
}
