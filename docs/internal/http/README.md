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


#### func (*LinkHeaderPager) Parse

```go
func (l *LinkHeaderPager) Parse(res *resty.Response) Paging
```

#### type NewRelicClient

```go
type NewRelicClient struct {
	Client resty.Client
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

#### func (*NewRelicClient) Delete

```go
func (nr *NewRelicClient) Delete(path string) error
```
nolint

#### func (*NewRelicClient) Get

```go
func (nr *NewRelicClient) Get(path string, params *map[string]string, result interface{}) error
```
Get executes an HTTP GET request.

#### func (*NewRelicClient) Post

```go
func (nr *NewRelicClient) Post(path string, body interface{}, result interface{}) error
```
nolint

#### func (*NewRelicClient) Put

```go
func (nr *NewRelicClient) Put(path string, body interface{}, result interface{}) error
```
nolint

#### type Pager

```go
type Pager interface {
	Parse(res *resty.Response) Paging
}
```


#### type Paging

```go
type Paging struct {
	Next string
}
```
