# http
--
    import "."


## Usage

```go
var (
	// ErrNotFound is returned when the resource was not found in New Relic.
	ErrNotFound = errors.New("newrelic: Resource not found")
)
```

#### type ErrorDetail

```go
type ErrorDetail struct {
	Title string `json:"title,omitempty"`
}
```

ErrorDetail represents the details of an ErrorResponse from New Relic.

#### type ErrorResponse

```go
type ErrorResponse struct {
	Detail *ErrorDetail `json:"error,omitempty"`
}
```

ErrorResponse represents an error response from New Relic.

#### func (*ErrorResponse) Error

```go
func (e *ErrorResponse) Error() string
```

#### type LinkHeaderPager

```go
type LinkHeaderPager struct{}
```

LinkHeaderPager represents a pagination implementation that adheres to RFC 5988.

#### func (*LinkHeaderPager) Parse

```go
func (l *LinkHeaderPager) Parse(res *resty.Response) Paging
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
	Parse(res *resty.Response) Paging
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
