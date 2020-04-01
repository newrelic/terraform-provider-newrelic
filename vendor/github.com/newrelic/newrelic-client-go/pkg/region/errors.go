package region

import (
	"fmt"
)

// InvalidError returns when the Region is not valid
type InvalidError struct {
	Message string
}

// Error string reported when an InvalidError happens
func (e InvalidError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("invalid region: %s", e.Message)
	}

	return "invalid region"
}

// ErrorNil returns an InvalidError message saying the value was nil
func ErrorNil() InvalidError {
	return InvalidError{
		Message: "value is nil",
	}
}
