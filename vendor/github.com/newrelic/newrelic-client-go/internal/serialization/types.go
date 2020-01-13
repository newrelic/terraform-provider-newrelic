package serialization

import (
	"strconv"
	"time"
)

// EpochTime is a type used for unmarshaling timestamps represented in epoch time.
// Its underlying type is time.Time.
type EpochTime time.Time

// MarshalJSON is responsible for marshaling the EpochTime type.
func (e EpochTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(e).Unix(), 10)), nil
}

// UnmarshalJSON is responsible for unmarshaling the EpochTime type.
func (e *EpochTime) UnmarshalJSON(s []byte) (err error) {
	q, err := strconv.ParseInt(string(s), 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(e) = time.Unix(q, 0)
	return
}

// Equal provides a comparator for the EpochTime type.
func (e EpochTime) Equal(u EpochTime) bool {
	return time.Time(e).Equal(time.Time(u))
}
