# infrastructure
--
    import "."


## Usage

#### type ErrorDetail

```go
type ErrorDetail struct {
	Status string `json:"status,omitempty"`
	Detail string `json:"detail,omitempty"`
}
```

ErrorDetail represents the details of an error response from New Relic
Infrastructure.

#### type ErrorResponse

```go
type ErrorResponse struct {
	Errors []*ErrorDetail `json:"errors,omitempty"`
}
```

ErrorResponse represents the an error response from New Relic Infrastructure.

#### func (*ErrorResponse) Error

```go
func (e *ErrorResponse) Error() string
```
Error surfaces an error message from the Infrastructure error response

#### type Infrastructure

```go
type Infrastructure struct {
}
```


#### func  New

```go
func New(config config.Config) Infrastructure
```
New is used to create a new Synthetics client instance.
