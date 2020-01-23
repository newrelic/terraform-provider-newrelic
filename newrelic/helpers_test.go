package newrelic

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

var (
	nrInternalAccount = os.Getenv("NR_ACC_TESTING") != ""
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

func TestParseIDs_BadIDs(t *testing.T) {
	_, err := parseIDs("12", 2)
	if err == nil {
		t.Fatal(err)
	}

	_, err = parseIDs("a:b", 2)
	if err == nil {
		t.Fatal(err)
	}
}

func TestSerializeIDs_Basic(t *testing.T) {
	id := serializeIDs([]int{1, 2})

	if id != "1:2" {
		t.Fatal(id)
	}
}

func testAccDeleteNewRelicAlertPolicy(name string) func() {
	return func() {
		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		params := alerts.ListPoliciesParams{
			Name: name,
		}
		policies, _ := client.Alerts.ListPolicies(&params)

		for _, p := range policies {
			if p.Name == name {
				_, _ = client.Alerts.DeletePolicy(p.ID)
				break
			}
		}
	}
}

// A custom check function to log the internal state during a test run.
// nolint:deadcode,unused
func logState(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		t.Logf("State: %s\n", s)

		return nil
	}
}
