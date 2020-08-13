// +build integration

package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicAlertMutingRule_Basic(t *testing.T) {
	resourceName := "newrelic_alert_muting_rule.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicAlertMutingRuleBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertMutingRuleExists(resourceName),
				),
			},
			// Test: Update
			// {
			// 	Config: testAccNewRelicNrqlAlertConditionConfigBasic(rName, "5", "180", ""),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
			// 	),
			// },
			// // Test: Import
			// {
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	// Ignore items with deprecated fields because
			// 	// we don't set deprecated fields on import
			// 	ImportStateVerifyIgnore: []string{"term", "nrql", "violation_time_limit"},
			// 	ImportStateIdFunc:       testAccImportStateIDFunc(resourceName, "static"),
			// },
		},
	})
}

func testAccNewRelicAlertMutingRuleBasic(
	name string,
) string {
	return fmt.Sprintf(`

resource "newrelic_alert_muting_rule" "foo" {
	name = "tf-test-%[1]s"
	enabled = true
	description = "muting rule test."
	condition {
		conditions {
			attribute 	= "product"
			operator 	= "EQUALS"
			values 		= ["APM"]
		}
		conditions {
			attribute 	= "event"
			operator 	= "EQUALS"
			values 		= ["Muted"]
		}
		operator = "AND"
	}
}
`, name)
}

func testAccCheckNewRelicAlertMutingRuleExists(n string) resource.TestCheckFunc {
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

		var accountID int
		var err error

		ids, err := parseHashedIDs(rs.Primary.ID)
		if err != nil {
			return err
		}

		conditionID := ids[1]
		accountID = providerConfig.AccountID

		if rs.Primary.Attributes["account_id"] != "" {
			accountID, err = strconv.Atoi(rs.Primary.Attributes["account_id"])
			if err != nil {
				return err
			}
		}

		found, err := client.Alerts.GetNrqlConditionQuery(accountID, strconv.Itoa(conditionID))
		if err != nil {
			return err
		}

		if found.ID != strconv.Itoa(conditionID) {
			return fmt.Errorf("alert condition not found: %v - %v", conditionID, found)
		}

		return nil
	}
}

//to do
// func testAccCheckNewRelicAlertMutingRuleDestroy(s *terraform.State) error {
// 	providerConfig := testAccProvider.Meta().(*ProviderConfig)
// 	client := providerConfig.NewClient

// 	for _, r := range s.RootModule().Resources {
// 		if r.Type != "newrelic_nrql_alert_condition" {
// 			continue
// 		}

// 		var accountID int
// 		var err error

// 		ids, err := parseHashedIDs(r.Primary.ID)
// 		if err != nil {
// 			return err
// 		}

// 		conditionID := ids[1]
// 		accountID = providerConfig.AccountID

// 		if r.Primary.Attributes["account_id"] != "" {
// 			accountID, err = strconv.Atoi(r.Primary.Attributes["account_id"])
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		if _, err = client.Alerts.GetNrqlConditionQuery(accountID, strconv.Itoa(conditionID)); err == nil {
// 			return fmt.Errorf("NRQL Alert condition still exists") //nolint:golint
// 		}
// 	}

// 	return nil
// }
