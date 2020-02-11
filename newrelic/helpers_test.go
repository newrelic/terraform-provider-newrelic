package newrelic

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

var (
	nrInternalAccount = os.Getenv("NR_ACC_TESTING") != ""
)

func TestParseIDs_Basic(t *testing.T) {
	ids, err := parseIDs("1:2", 2)

	require.NoError(t, err)
	require.Equal(t, 2, len(ids))
	require.Equal(t, 1, ids[0])
	require.Equal(t, 2, ids[1])
}

func TestParseIDs_BadIDs(t *testing.T) {
	_, err := parseIDs("12", 2)
	require.Error(t, err)

	_, err = parseIDs("a:b", 2)
	require.Error(t, err)
}

func TestSerializeIDs_Basic(t *testing.T) {
	id := serializeIDs([]int{1, 2})

	require.Equal(t, "1:2", id)
}

func TestStripWhitespace(t *testing.T) {
	json := " { \"key\": \"value\" } "
	e := "{\"key\":\"value\"}"
	a := stripWhitespace(json)

	require.Equal(t, e, a)
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
