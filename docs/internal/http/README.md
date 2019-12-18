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

ErrorDetail represents a New Relic response error detail.

#### type ErrorNotFound

```go
type ErrorNotFound struct{}
```

ErrorNotFound is returned when a 404 response is returned from New Relic's APIs.

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

ErrorResponse provides an interface for obtaining a single error message from an
error response object.

#### type ErrorUnexpectedStatusCode

```go
type ErrorUnexpectedStatusCode struct {
}
```

ErrorUnexpectedStatusCode is returned when an unexpected status code is returned
from New Relic's APIs.

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

NewRelicClient represents a client for communicating with the New Relic APIs.

#### func  NewClient

```go
func NewClient(config config.Config) NewRelicClient
```
NewClient is used to create a new instance of NewRelicClient.

#### func  NewTestAPIClient

```go
func NewTestAPIClient(handler http.Handler) NewRelicClient
```
NewTestAPIClient returns a test NewRelicClient instance that is configured to
communicate with a mock server.

#### func (*NewRelicClient) Delete

```go
func (c *NewRelicClient) Delete(url string, params *map[string]string, reqBody interface{}, value interface{}) (*http.Response, error)
```
Delete represents an HTTP DELETE request to a New Relic API.

#### func (*NewRelicClient) Get

```go
func (c *NewRelicClient) Get(url string, params *map[string]string, reqBody interface{}, value interface{}) (*http.Response, error)
```
Get represents an HTTP GET request to a New Relic API.

#### func (*NewRelicClient) Post

```go
func (c *NewRelicClient) Post(url string, params *map[string]string, reqBody interface{}, value interface{}) (*http.Response, error)
```
Post represents an HTTP POST request to a New Relic API.

#### func (*NewRelicClient) Put

```go
func (c *NewRelicClient) Put(url string, params *map[string]string, reqBody interface{}, value interface{}) (*http.Response, error)
```
Put represents an HTTP PUT request to a New Relic API.

#### func (*NewRelicClient) SetErrorValue

```go
func (c *NewRelicClient) SetErrorValue(v ErrorResponse) *NewRelicClient
```
SetErrorValue is used to unmarshal error body responses in JSON format.

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
