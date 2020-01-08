package client

import (
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	insightsInsertURL = "https://insights-collector.newrelic.com/v1/accounts"
	insightsQueryURL  = "https://insights-api.newrelic.com/v1/accounts"

	// Minimum length check for a valid NRQL statement
	minValidNRQLLength = 8 // "SELECT 1"

	// DefaultBatchTimeout is the amount of time to submit batches even if the event count hasn't been hit
	DefaultBatchTimeout = 1 * time.Minute
	// DefaultBatchEventCount is the maximum number of events before sending a batch (fuzzy)
	DefaultBatchEventCount = 950
	// DefaultWorkerCount is the number of background workers consuming and sending events
	DefaultWorkerCount = 1

	// DefaultInsertRequestTimeout is the amount of seconds to wait for a insert response
	DefaultInsertRequestTimeout = 10 * time.Second
	// DefaultQueryRequestTimeout is the amount of seconds to wait for a query response
	DefaultQueryRequestTimeout = 20 * time.Second

	// DefaultRetries is how many times to attempt the query
	DefaultRetries = 3
	// DefaultRetryWaitTime is the amount of seconds between query attempts
	DefaultRetryWaitTime = 5 * time.Second
)

// Compression to use during transport.
type Compression int32

// Supported / recognized types of compression
const (
	None    Compression = iota
	Deflate Compression = iota
	Gzip    Compression = iota
	Zlib    Compression = iota
)

// Client is the building block of the insert and query clients
type Client struct {
	URL            *url.URL
	Logger         *log.Logger
	RequestTimeout time.Duration
	RetryCount     int
	RetryWait      time.Duration
}

// InsertClient contains all of the configuration required for inserts
type InsertClient struct {
	InsertKey   string
	eventQueue  chan []byte
	eventTimer  *time.Timer
	flushQueue  chan bool
	WorkerCount int
	BatchSize   int
	BatchTime   time.Duration
	Compression Compression
	Client
	Statistics
}

// Statistics about the inserted data
type Statistics struct {
	// the number of events added using EnqueueEvent
	EventCount int64
	// the number of events that finished processing (both successfully and not) in batch mode
	ProcessedEventCount int64
	// the number of times a Flush has been requested
	FlushCount int64
	// the number of bytes that have been attempted to be written to the insights API
	ByteCount int64
	// the number of times a flush has occurred with a full batch of events
	FullFlushCount int64
	// the number of times a flush has occurred with a partial batch of events
	PartialFlushCount int64
	// the number of times flush has been called because the BatchTime timer expired
	TimerExpiredCount int64
	// the number of times failed batches have been retried
	InsightsRetryCount int64
	// not used
	HTTPErrorCount int64
}

// Assumption here that responses from insights are either success or error.
type insertResponse struct {
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success,omitempty"`
}

// QueryClient contains all of the configuration required for queries
type QueryClient struct {
	QueryKey string
	Client
}

// QueryResponse used to decode the JSON response from Insights
type QueryResponse struct {
	Results  []map[string]interface{} `json:"results"`
	Facets   []map[string]interface{} `json:"facets"`
	Metadata QueryMetadata            `json:"metadata"`
}

// QueryMetadata used to decode the JSON response metadata from Insights
type QueryMetadata struct {
	Contents        interface{} `json:"contents"`
	EventType       string      `json:"eventType"`
	OpenEnded       bool        `json:"openEnded"`
	BeginTime       time.Time   `json:"beginTime"`
	EndTime         time.Time   `json:"endTime"`
	BeginTimeMillis int64       `json:"beginTimeMillis"`
	EndTimeMillis   int64       `json:"endTimeMillis"`
	RawSince        string      `json:"rawSince"`
	RawUntil        string      `json:"rawUntil"`
	RawCompareWith  string      `json:"rawCompareWith"`
}
