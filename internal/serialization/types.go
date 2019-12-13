package serialization

import (
	"strconv"
	"time"
)

// Epoch is a type used for unmarshaling timestamps represented in epoch time.
// Its underlying type is time.Time.
type Epoch time.Time

// MarshalJSON is responsible for marshaling the Epoch type.
func (e Epoch) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(e).Unix(), 10)), nil
}

// UnmarshalJSON is responsible for unmarshaling the Epoch type.
func (e *Epoch) UnmarshalJSON(s []byte) (err error) {
	q, err := strconv.ParseInt(string(s), 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(e) = time.Unix(q, 0)
	return
}

// Equal provides a comparator for the Epoch type.
func (e Epoch) Equal(u Epoch) bool {
	return time.Time(e).Equal(time.Time(u))
}
