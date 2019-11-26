# client
--
    import "github.com/newrelic/go-insights/client"


## Usage

```go
const (

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
```

#### type Client

```go
type Client struct {
	URL            *url.URL
	Logger         *log.Logger
	RequestTimeout time.Duration
	RetryCount     int
	RetryWait      time.Duration
}
```

Client is the building block of the insert and query clients

#### func (*Client) UseCustomURL

```go
func (c *Client) UseCustomURL(customURL string)
```
UseCustomURL allows overriding the default Insights Host / Scheme.

#### type Compression

```go
type Compression int32
```

Compression to use during transport.

```go
const (
	None    Compression = iota
	Deflate Compression = iota
	Gzip    Compression = iota
	Zlib    Compression = iota
)
```
Supported / recognized types of compression

#### type InsertClient

```go
type InsertClient struct {
	InsertKey string

	WorkerCount int
	BatchSize   int
	BatchTime   time.Duration
	Compression Compression
	Client
	Statistics
}
```

InsertClient contains all of the configuration required for inserts

#### func  NewInsertClient

```go
func NewInsertClient(insertKey string, accountID string) *InsertClient
```
NewInsertClient makes a new client for the user to send data with

#### func (*InsertClient) EnqueueEvent

```go
func (c *InsertClient) EnqueueEvent(data interface{}) (err error)
```
EnqueueEvent handles the queueing. Only works in batch mode.

#### func (*InsertClient) Flush

```go
func (c *InsertClient) Flush() error
```
Flush gives the user a way to manually flush the queue in the foreground. This
is also used by watchdog when the timer expires.

#### func (*InsertClient) PostEvent

```go
func (c *InsertClient) PostEvent(data interface{}) error
```
PostEvent allows sending a single event directly.

#### func (*InsertClient) SetCompression

```go
func (c *InsertClient) SetCompression(compression Compression)
```
SetCompression allows modification of the compression type used in communication

#### func (*InsertClient) Start

```go
func (c *InsertClient) Start() error
```
Start runs the insert client in batch mode.

#### func (*InsertClient) StartListener

```go
func (c *InsertClient) StartListener(inputChannel chan interface{}) (err error)
```
StartListener creates a goroutine that consumes from a channel and Enqueues
events as to not block the writing of events to the channel

#### func (*InsertClient) Validate

```go
func (c *InsertClient) Validate() error
```
Validate makes sure the InsertClient is configured correctly for use

#### type QueryClient

```go
type QueryClient struct {
	QueryKey string
	Client
}
```

QueryClient contains all of the configuration required for queries

#### func  NewQueryClient

```go
func NewQueryClient(queryKey, accountID string) *QueryClient
```
NewQueryClient makes a new client for the user to query with.

#### func (*QueryClient) Query

```go
func (c *QueryClient) Query(nrqlQuery string, response interface{}) (err error)
```
Query initiates an Insights query, with the JSON parsed into 'response' struct

#### func (*QueryClient) QueryEvents

```go
func (c *QueryClient) QueryEvents(nrqlQuery string) (response *QueryResponse, err error)
```
QueryEvents initiates an Insights query, returns a response for parsing

#### func (*QueryClient) Validate

```go
func (c *QueryClient) Validate() error
```
Validate makes sure the QueryClient is configured correctly for use

#### type QueryMetadata

```go
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
```

QueryMetadata used to decode the JSON response metadata from Insights

#### type QueryResponse

```go
type QueryResponse struct {
	Results  []map[string]interface{} `json:"results"`
	Facets   []map[string]interface{} `json:"facets"`
	Metadata QueryMetadata            `json:"metadata"`
}
```

QueryResponse used to decode the JSON response from Insights

#### type Statistics

```go
type Statistics struct {
	EventCount         int64
	FlushCount         int64
	ByteCount          int64
	FullFlushCount     int64
	PartialFlushCount  int64
	TimerExpiredCount  int64
	InsightsRetryCount int64
	HTTPErrorCount     int64
}
```

Statistics about the inserted data
