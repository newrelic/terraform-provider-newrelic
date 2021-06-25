// +build unit

package newrelic

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type testCase struct {
	val          interface{}
	f            schema.SchemaValidateFunc
	expectedErr  *regexp.Regexp
	expectedWarn *regexp.Regexp
}

func TestValidationValidateViolationCloseTimer(t *testing.T) {
	runTestCases(t, []testCase{
		{
			val:         0,
			f:           validateViolationCloseTimer(),
			expectedWarn: regexp.MustCompile(`0 is no longer a valid value. Using the default value of 24`),
		},
		{
			val: 2,
			f:   validateViolationCloseTimer(),
		},
		{
			val:         13,
			f:           validateViolationCloseTimer(),
			expectedErr: regexp.MustCompile(`expected [\w]+ to be one of 1, 2, 4, 8, 12, 24, 48, 72, got 13`),
		},
		{
			val:         "foo",
			f:           intInSlice([]int{1, 2, 3}),
			expectedErr: regexp.MustCompile(`expected type of [\w]+ to be int`),
		},
	})
}

func TestValidationIntInInSlice(t *testing.T) {
	runTestCases(t, []testCase{
		{
			val: 2,
			f:   intInSlice([]int{1, 2, 3}),
		},
		{
			val:         4,
			f:           intInSlice([]int{1, 2, 3}),
			expectedErr: regexp.MustCompile(`expected [\w]+ to be one of \[1 2 3\], got 4`),
		},
		{
			val:         "foo",
			f:           intInSlice([]int{1, 2, 3}),
			expectedErr: regexp.MustCompile(`expected type of [\w]+ to be int`),
		},
	})
}

func TestValidationFloat64Gte(t *testing.T) {
	runTestCases(t, []testCase{
		{
			val: 1.1,
			f:   float64Gte(1.1),
		},
		{
			val: 1.2,
			f:   float64Gte(1.1),
		},
		{
			val:         "foo",
			f:           float64Gte(1.1),
			expectedErr: regexp.MustCompile(`expected type of [\w]+ to be float64`),
		},
		{
			val:         0.1,
			f:           float64Gte(1.1),
			expectedErr: regexp.MustCompile(`expected [\w]+ to be greater than or equal to 1.1, got 0.1`),
		},
	})
}

func TestValidationFloat64AtLeast(t *testing.T) {
	runTestCases(t, []testCase{
		{
			val: 1.1,
			f:   float64AtLeast(1.1),
		},
		{
			val: 1.2,
			f:   float64AtLeast(1.1),
		},
		{
			val:         "foo",
			f:           float64AtLeast(1.1),
			expectedErr: regexp.MustCompile(`expected type of [\w]+ to be float64`),
		},
		{
			val:         0.1,
			f:           float64AtLeast(1.1),
			expectedErr: regexp.MustCompile(`expected [\w]+ to be at least 1.100000, got 0.100000`),
		},
	})
}

func runTestCases(t *testing.T, cases []testCase) {
	matchErr := func(errs []error, r *regexp.Regexp) bool {
		// err must match one provided
		for _, err := range errs {
			if r.MatchString(err.Error()) {
				return true
			}
		}

		return false
	}

	matchWarn := func(warnings []string, r *regexp.Regexp) bool {
		// warning must match one provided
		for _, warn := range warnings {
			if r.MatchString(warn) {
				return true
			}
		}

		return false
	}

	for i, tc := range cases {
		warnings, errs := tc.f(tc.val, "test_property")
		if len(warnings) == 0 && tc.expectedWarn == nil {
			continue
		}
		if !matchWarn(warnings, tc.expectedWarn) {
			t.Fatalf("expected test case %d to produce warning matching \"%s\", got %v", i, tc.expectedWarn, warnings)
		}

		if len(errs) == 0 && tc.expectedErr == nil {
			continue
		}

		if !matchErr(errs, tc.expectedErr) {
			t.Fatalf("expected test case %d to produce error matching \"%s\", got %v", i, tc.expectedErr, errs)
		}
	}
}
