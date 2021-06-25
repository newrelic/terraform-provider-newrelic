package newrelic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func float64Gte(gte float64) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(float64)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be float64", k))
			return
		}

		if v >= gte {
			return
		}

		es = append(es, fmt.Errorf("expected %s to be greater than or equal to %v, got %v", k, gte, v))
		return
	}
}

func validateViolationCloseTimer() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		val, ok := i.(int)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
			return
		}
		switch val {
		case 1, 2, 4, 8, 12, 24, 48, 72:
		case 0:
			warnings = append(warnings, "0 is no longer a valid value. Using the default value of 24")
		default:
			errors = append(errors, fmt.Errorf("expected %s to be one of %s, got %v", k, "1, 2, 4, 8, 12, 24, 48, 72", val))
		}
		return warnings, errors
	}
}



func intInSlice(valid []int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(int)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be int", k))
			return
		}

		for _, p := range valid {
			if v == p {
				return
			}
		}

		es = append(es, fmt.Errorf("expected %s to be one of %v, got %v", k, valid, v))
		return
	}
}

// float64AtLeast returns a SchemaValidateFunc which tests if the provided value
// is of type float64 and is at least min (inclusive)
func float64AtLeast(min float64) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(float64)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be float64", k))
			return
		}

		if v < min {
			es = append(es, fmt.Errorf("expected %s to be at least %f, got %f", k, min, v))
			return
		}

		return
	}
}
