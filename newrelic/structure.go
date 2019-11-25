package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Takes the result of flatmap.Expand for an array of ints
// and returns a []*int
func expandIntList(configured []interface{}) []int {
	vs := make([]int, 0, len(configured))
	for _, v := range configured {
		if val, ok := v.(int); ok {
			vs = append(vs, val)
		}
	}
	return vs
}

// Takes the result of schema.Set of strings and returns a []int
func expandIntSet(configured *schema.Set) []int {
	return expandIntList(configured.List())
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []*string
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
func expandStringSet(configured *schema.Set) []string {
	return expandStringList(configured.List())
}
