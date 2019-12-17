# http
--
    import "github.com/newrelic/newrelic-client-go/internal/http"


## Usage

```go
var (
	// ErrNotFound is returned when the resource was not found in New Relic.
	ErrNotFound = errors.New("newrelic: Resource not found")
)
```

#### func  RetryPolicy

```go
func RetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error)
```
RetryPolicy provides a callback for retryablehttp's CheckRetry, which will retry
on connection errors and server errors.

#### type DefaultErrorResponse

```go
type DefaultErrorResponse struct {
	ErrorDetail ErrorDetail `json:"error"`
}
```

DefaultErrorResponse represents the default error response from New Relic.

#### func (*DefaultErrorResponse) Error

```go
func (e *DefaultErrorResponse) Error() string
```

#### type ErrorDetail

```go
type ErrorDetail struct {
	Title string `json:"title"`
}
```


#### type ErrorNotFound

```go
type ErrorNotFound struct{}
```


#### func (*ErrorNotFound) Error

```go
func (e *ErrorNotFound) Error() string
```

#### type ErrorResponse

```go
type ErrorResponse interface {
	Error() string
}
```


#### type ErrorUnexpectedStatusCode

```go
type ErrorUnexpectedStatusCode struct {
}
```


#### func (*ErrorUnexpectedStatusCode) Error

```go
func (e *ErrorUnexpectedStatusCode) Error() string
```

#### type LinkHeaderPager

```go
type LinkHeaderPager struct{}
```

LinkHeaderPager represents a pagination implementation that adheres to RFC 5988.

#### func (*LinkHeaderPager) Parse

```go
func (l *LinkHeaderPager) Parse(resp *http.Response) Paging
```
Parse is used to parse a pagination context from an HTTP response.

#### type NewRelicClient

```go
type NewRelicClient struct {
	Client *retryablehttp.Client
	Config config.Config
}
```


#### func  NewClient

```go
func NewClient(config config.Config) NewRelicClient
```

#### func  NewTestAPIClient

```go
func NewTestAPIClient(handler http.Handler) NewRelicClient
```
NewTestAPIClient returns a test NewRelicClient instance that is configured to
communicate with a mock server.

#### func (*NewRelicClient) Get

```go
func (c *NewRelicClient) Get(url string, params *map[string]string, reqBody interface{}, value interface{}) (*http.Response, error)
```

#### func (*NewRelicClient) SetErrorValue

```go
func (c *NewRelicClient) SetErrorValue(v ErrorResponse) *NewRelicClient
```

#### type Pager

```go
type Pager interface {
	Parse(res *http.Response) Paging
}
```

Pager represents a pagination implementation.

#### type Paging

```go
type Paging struct {
	Next string
}
```

Paging represents the pagination context returned from the Pager implementation.
