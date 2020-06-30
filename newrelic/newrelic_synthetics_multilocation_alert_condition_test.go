package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccNewRelicSyntheticsMultiLocationAlertCondition_Basic(t *testing.T) {
	resourceName := "newrelic_synthetics_multilocation_alert_condition.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsMultiLocationConditionConfigBasic(rName, "1", "2", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMultiLocationAlertConditionExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsMultiLocationConditionConfigBasic(rName, "11", "12", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMultiLocationAlertConditionExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore items with deprecated fields because
				// we don't set deprecated fields on import
				// ImportStateVerifyIgnore: []string{"term", "nrql", "violation_time_limit"},
				// ImportStateIdFunc: testAccImportStateIDFunc(resourceName, "static"),
			},
		},
	})
}

func testAccCheckNewRelicSyntheticsMultiLocationAlertConditionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providerConfig := testAccProvider.Meta().(*ProviderConfig)
		client := providerConfig.NewClient

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no alert condition ID is set")
		}

		var err error

		ids, err := parseHashedIDs(rs.Primary.ID)
		if err != nil {
			return err
		}

		conditionID := ids[1]
		policyID := ids[0]

		found, err := client.Alerts.GetMultiLocationSyntheticsCondition(policyID, conditionID)
		if err != nil {
			return err
		}

		if found.ID != conditionID {
			return fmt.Errorf("synthetics multi-location alert condition not found: %v - %v", conditionID, found)
		}

		return nil
	}
}

func testAccNewRelicSyntheticsMultiLocationConditionConfigBasic(
	name string,
	criticalThreshold string,
	warningThreshold string,
	conditionalAttrs string,
) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_synthetics_multilocation_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id

  name                         = "tf-test-%[1]s"
  runbook_url                  = "https://foo.example.com"
  enabled                      = true
  violation_time_limit_seconds = "3600"

	entities = [
		"b62bcdde-6c73-4b7c-afb8-e18bae3cf4db"
	]

	critical {
    threshold = %[2]s
	}

	warning {
    threshold = %[3]s
	}

	%[4]s
}
`, name, criticalThreshold, warningThreshold, conditionalAttrs)
}

func TestExpandMultiLocationSyntheticsCondition(t *testing.T) {
	var criticalTerms []map[string]interface{}
	criticalTerms = append(criticalTerms, map[string]interface{}{
		"threshold": 1,
	})

	var warningTerms []map[string]interface{}
	warningTerms = append(warningTerms, map[string]interface{}{
		"threshold": 9,
	})

	cases := map[string]struct {
		Data         map[string]interface{}
		ExpectErr    bool
		ExpectReason string
		Expanded     *alerts.MultiLocationSyntheticsCondition
	}{
		"valid minimal": {
			Data: map[string]interface{}{},
		},
		"with critical term": {
			Data: map[string]interface{}{
				// "nrql": []interface{}{nrql},
				"name":                         "yoo hoo",
				"enabled":                      true,
				"violation_time_limit_seconds": 100,
				"critical":                     criticalTerms,
			},
			Expanded: &alerts.MultiLocationSyntheticsCondition{
				Name:                      "yoo hoo",
				Enabled:                   true,
				ViolationTimeLimitSeconds: 100,
				Terms: []alerts.MultiLocationSyntheticsConditionTerm{
					{
						Threshold: 1,
						Priority:  "critical",
					},
				},
			},
		},
		"with critical and warning term": {
			Data: map[string]interface{}{
				"name":                         "yoo hoo",
				"enabled":                      true,
				"violation_time_limit_seconds": 100,
				"critical":                     criticalTerms,
				"warning":                      warningTerms,
			},
			Expanded: &alerts.MultiLocationSyntheticsCondition{
				Name:                      "yoo hoo",
				Enabled:                   true,
				ViolationTimeLimitSeconds: 100,
				Terms: []alerts.MultiLocationSyntheticsConditionTerm{
					{
						Threshold: 1,
						Priority:  "critical",
					},
					{
						Threshold: 9,
						Priority:  "warning",
					},
				},
			},
		},
	}

	r := resourceNewRelicSyntheticsMultiLocationAlertCondition()

	for _, tc := range cases {
		d := r.TestResourceData()

		for k, v := range tc.Data {
			if k == "critical" || k == "warning" {
				var terms []map[string]interface{}

				terms = append(terms, v.([]map[string]interface{})...)

				if err := d.Set(k, terms); err != nil {
					t.Fatalf("err: %s", err)
				}
			} else {
				if err := d.Set(k, v); err != nil {
					t.Fatalf("err: %s", err)
				}
			}
		}

		expanded, err := expandMultiLocationSyntheticsCondition(d)

		if tc.ExpectErr {
			require.NotNil(t, err)
			require.Equal(t, err.Error(), tc.ExpectReason)
		} else {
			require.Nil(t, err)
		}

		if tc.Expanded != nil {
			if tc.Expanded.Name != "" {
				require.Equal(t, tc.Expanded.Name, expanded.Name)
			}

			require.Equal(t, tc.Expanded.Enabled, expanded.Enabled)
			require.Equal(t, tc.Expanded.ViolationTimeLimitSeconds, expanded.ViolationTimeLimitSeconds)

			if len(tc.Expanded.Terms) > 0 {
				assert.Equal(t, tc.Expanded.Terms, expanded.Terms)
			}
		}
	}

}

func TestFlattenMultiLocationSyntheticsCondition(t *testing.T) {
	var criticalTerms []map[string]interface{}
	criticalTerms = append(criticalTerms, map[string]interface{}{
		"threshold": 1,
	})

	var warningTerms []map[string]interface{}
	warningTerms = append(warningTerms, map[string]interface{}{
		"threshold": 9,
	})

	r := resourceNewRelicSyntheticsMultiLocationAlertCondition()

	cases := map[string]struct {
		Data         map[string]interface{}
		ExpectErr    bool
		ExpectReason string
		Flattened    *alerts.MultiLocationSyntheticsCondition
	}{
		"minimal": {
			Data: map[string]interface{}{
				"policy_id": 123,
				"name":      "testing123",
			},
			Flattened: &alerts.MultiLocationSyntheticsCondition{
				Name: "testing123",
			},
		},
		"less minimal": {
			Data: map[string]interface{}{
				"policy_id": 123,
				"name":      "testing123",
				"enabled":   true,
			},
			Flattened: &alerts.MultiLocationSyntheticsCondition{
				Name:    "testing123",
				Enabled: true,
			},
		},
		"with critical and warning terms": {
			Data: map[string]interface{}{
				"policy_id":                    123,
				"name":                         "testing123",
				"enabled":                      true,
				"violation_time_limit_seconds": 100,
				"critical":                     criticalTerms,
				"warning":                      warningTerms,
				"entities":                     []string{"one", "two"},
			},
			Flattened: &alerts.MultiLocationSyntheticsCondition{
				ID:                        503,
				Name:                      "testing123",
				Enabled:                   true,
				ViolationTimeLimitSeconds: 100,
				Terms: []alerts.MultiLocationSyntheticsConditionTerm{
					{
						Threshold: 1,
						Priority:  "critical",
					},
					{
						Threshold: 9,
						Priority:  "warning",
					},
				},
				Entities: []string{"one", "two"},
			},
		},
	}

	for _, tc := range cases {
		if tc.Flattened != nil {
			d := r.TestResourceData()
			id := fmt.Sprintf("%d:%d", tc.Data["policy_id"], tc.Flattened.ID)
			d.SetId(id)

			err := flattenMultiLocationSyntheticsCondition(tc.Flattened, d)
			assert.NoError(t, err)

			for k, v := range tc.Data {
				var x interface{}
				var ok bool
				if x, ok = d.GetOk(k); !ok {
					t.Fatalf("failed to get k: %s: %s", k, err)
				}

				if k == "critical" || k == "warning" {
					for _, term := range tc.Flattened.Terms {
						if k == term.Priority {
							assert.Equal(t, v.([]map[string]interface{})[0]["threshold"], term.Threshold)
						}
					}
				} else if k == "entities" {
					assert.Equal(t, v.([]string), v)
				} else {
					assert.Equal(t, x, v)
				}
			}
		}
	}
}
