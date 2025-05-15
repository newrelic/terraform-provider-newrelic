package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Takes the result of flatmap.Expand for an array of ints
// and returns a []*int
// nolint:unused,deadcode
func expandIntList(configured []interface{}) []int {
	vs := make([]int, 0, len(configured))
	for _, v := range configured {
		if val, ok := v.(int); ok {
			vs = append(vs, val)
		}
	}
	return vs
}

// nolint:unused,deadcode
// Takes the result of schema.Set of strings and returns a []int
func expandIntSet(configured *schema.Set) []int {
	return expandIntList(configured.List())
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []*string
// nolint:unused,deadcode
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}

// Takes the result of schema.Set of strings and returns a []string
// nolint:unused,deadcode
func expandStringSet(configured *schema.Set) []string {
	return expandStringList(configured.List())
}
