// +build integration

package newrelic

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
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
		// CheckDestroy: testAccCheckNewRelicAlertMutingRuleDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicAlertMutingRuleBasic(rName, "new muting rule", "product", "EQUALS", "APM"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertMutingRuleExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicAlertMutingRuleBasic(rName, "second muting rule", "conditionType", "NOT_EQUALS", "baseline"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertMutingRuleExists(resourceName),
				),
			},
			// // Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true},
		},
	})
}

func testAccNewRelicAlertMutingRuleBasic(
	name string,
	description string,
	attribute string,
	operator string,
	values string,
) string {
	return fmt.Sprintf(`

resource "newrelic_alert_muting_rule" "foo" {
	name = "tf-test-%[1]s"
	enabled = true
	description = "%[2]s"
	condition {
		conditions {
			attribute 	= "%[3]s"
			operator 	= "EQUALS"
			values 		= ["%[5]s"]
		}
		conditions {
			attribute 	= "conditionType"
			operator 	= "%[4]s"
			values 		= ["static"]
		}
		operator = "AND"
	}
}
`, name, description, attribute, operator, values)
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

		accountID = ids[0]
		mutingRuleID := ids[1]

		if rs.Primary.Attributes["account_id"] != "" {
			accountID, err = strconv.Atoi(rs.Primary.Attributes["account_id"])
			if err != nil {
				return err
			}
		}

		found, err := client.Alerts.GetMutingRule(accountID, mutingRuleID)
		if err != nil {
			return err
		}

		if found.ID != mutingRuleID {
			return fmt.Errorf("alert muting rule not found: %v - %v", mutingRuleID, found)
		}

		return nil
	}
}

func testAccCheckNewRelicAlertMutingRuleDestroy(s *terraform.State) error {
	providerConfig := testAccProvider.Meta().(*ProviderConfig)
	client := providerConfig.NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_alert_muting_rule" {
			continue
		}

		var accountID int
		var err error

		ids, err := parseHashedIDs(r.Primary.ID)
		if err != nil {
			return err
		}

		mutingRuleID := ids[1]
		accountID = providerConfig.AccountID

		if r.Primary.Attributes["account_id"] != "" {
			accountID, err = strconv.Atoi(r.Primary.Attributes["account_id"])
			if err != nil {
				return err
			}
		}

		if _, err = client.Alerts.GetMutingRule(accountID, mutingRuleID); err == nil {
			return fmt.Errorf("Alert muting rule still exists") //nolint:golint
		}
	}

	return nil
}

func TestValidateMutingRuleConditionAttribute_ValidTag(t *testing.T) {
	// It should accept a dot-separated tag and value
	validTag := "tag.EC2Timezone with sp aces"
	resourceName := "condition.0.conditions.0.attribute"

	warns, errs := validateMutingRuleConditionAttribute(validTag, resourceName)
	expectedErrs := []error(nil)
	expectedWarns := []string([]string(nil))

	require.Equal(t, expectedWarns, warns)
	require.Equal(t, expectedErrs, errs)

}

func TestValidateMutingRuleConditionAttribute_InvalidTag(t *testing.T) {
	// It should reject "tag" without value
	invalidTag := "tag"
	resourceName := "condition.0.conditions.0.attribute"

	warns, errs := validateMutingRuleConditionAttribute(invalidTag, resourceName)
	expectedErrs := []error{errors.New("\"condition.0.attribute\" of \"tag\" must be in the format tag.tag_name")}
	expectedWarns := []string([]string(nil))

	require.Equal(t, expectedWarns, warns)
	require.Equal(t, expectedErrs, errs)
}
