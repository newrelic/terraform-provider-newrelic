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
	Client resty.Client
}
```

NewRelicClient is the internal client for communicating with the New Relic APIs.

#### func  NewClient

```go
func NewClient(config config.Config) NewRelicClient
```
NewClient is used to create a new instance of the NewRelicClient type.

#### func  NewTestAPIClient

```go
func NewTestAPIClient(handler http.Handler) NewRelicClient
```
NewTestAPIClient returns a test NewRelicClient instance that is configured to
communicate with a mock server.

#### func (*NewRelicClient) Delete

```go
func (n *NewRelicClient) Delete(path string) error
```
Delete executes an HTTP DELETE request. nolint

#### func (*NewRelicClient) Get

```go
func (n *NewRelicClient) Get(path string, params *map[string]string, result interface{}) error
```
Get executes an HTTP GET request.

#### func (*NewRelicClient) Post

```go
func (n *NewRelicClient) Post(path string, body interface{}, result interface{}) error
```
Post executes an HTTP POST request. nolint

#### func (*NewRelicClient) Put

```go
func (n *NewRelicClient) Put(path string, body interface{}, result interface{}) error
```
Put executes an HTTP PUT request. nolint

#### func (NewRelicClient) SetError

```go
func (n NewRelicClient) SetError(err interface{}) NewRelicClient
```
SetError allows for registering different well-known error response structures.

#### func (NewRelicClient) SetPager

```go
func (n NewRelicClient) SetPager(pager Pager) NewRelicClient
```
SetPager allows for use of different pagination implementations.

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

#### type ReplacementClient

```go
type ReplacementClient struct {
	Client *retryablehttp.Client
	Config config.ReplacementConfig
}
```


#### func  NewReplacementClient

```go
func NewReplacementClient(config config.ReplacementConfig) ReplacementClient
```

#### func (*ReplacementClient) Get

```go
func (c *ReplacementClient) Get(url string, params *map[string]string, reqBody interface{}, value interface{}) (*http.Response, error)
```

#### func (*ReplacementClient) SetErrorValue

```go
func (c *ReplacementClient) SetErrorValue(v ErrorResponse) *ReplacementClient
```

#### type RestyErrorDetail

```go
type RestyErrorDetail struct {
	Title string `json:"title,omitempty"`
}
```

ErrorDetail represents the details of an ErrorResponse from New Relic.

#### type RestyErrorResponse

```go
type RestyErrorResponse struct {
	Detail *ErrorDetail `json:"error,omitempty"`
}
```

RestyErrorResponse represents an error response from New Relic.

#### func (*RestyErrorResponse) Error

```go
func (e *RestyErrorResponse) Error() string
```
