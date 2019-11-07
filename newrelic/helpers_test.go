package newrelic

import (
	"testing"
)

func TestParseIDs_Basic(t *testing.T) {
	ids, err := parseIDs("1:2", 2)
	if err != nil {
		t.Fatal(err)
	}

	if len(ids) != 2 {
		t.Fatal(len(ids))
	}

	if ids[0] != 1 || ids[1] != 2 {
		t.Fatal(ids)
	}
}

func TestSerializeIDs_Basic(t *testing.T) {
	id := serializeIDs([]int{1, 2})

	if id != "1:2" {
		t.Fatal(id)
	}
}

func deletePolicy(name string) func() {
	return func() {
		client := testAccProvider.Meta().(*ProviderConfig).Client
		policies, _ := client.ListAlertPolicies()

		for _, p := range policies {
			if p.Name == name {
				_ = client.DeleteAlertPolicy(p.ID)
				break
			}
		}
	}
}
