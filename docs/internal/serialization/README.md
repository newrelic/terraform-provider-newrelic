# serialization
--
    import "."


## Usage

#### type Epoch

```go
type Epoch time.Time
```

Epoch is a type used for unmarshaling timestamps represented in epoch time. Its
underlying type is time.Time.

#### func (Epoch) Equal

```go
func (e Epoch) Equal(u Epoch) bool
```
Equal provides a comparator for the Epoch type.

#### func (Epoch) MarshalJSON

```go
func (e Epoch) MarshalJSON() ([]byte, error)
```
MarshalJSON is responsible for marshaling the Epoch type.

#### func (*Epoch) UnmarshalJSON

```go
func (e *Epoch) UnmarshalJSON(s []byte) (err error)
```
UnmarshalJSON is responsible for unmarshaling the Epoch type.
