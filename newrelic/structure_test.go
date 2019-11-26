package newrelic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestExpandIntList_Basic(t *testing.T) {
	testVals := []interface{}{1, 2}
	ids := expandIntList(testVals)

	for i, v := range ids {
		if v != testVals[i].(int) {
			t.Fatal("int list expansion failed")
		}
	}
}

func TestExpandIntSet_Basic(t *testing.T) {
	testSchema := &schema.Set{F: schema.HashInt}
	testSchema.Add(1)

	ids := expandIntSet(testSchema)

	if len(ids) != 1 {
		t.Fatal("int set expansion failed")
	}
}

func TestExpandStringList_Basic(t *testing.T) {
	testVals := []interface{}{"one", "two"}
	ids := expandStringList(testVals)

	for i, v := range ids {
		if v != testVals[i].(string) {
			t.Fatal("string list expansion failed")
		}
	}

}

func TestExpandStringSet_Basic(t *testing.T) {
	testSchema := &schema.Set{F: schema.HashString}
	testSchema.Add("one")

	ids := expandStringSet(testSchema)

	if len(ids) != 1 {
		t.Fatal("string set expansion failed")
	}
}
