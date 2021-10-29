//go:build integration
// +build integration

package newrelic

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestValidateNaiveDateTime_Validates(t *testing.T) {
	validDate := "2021-02-21T15:30:00"
	resourceName := "schedule.0.end_repeat"

	warns, errs := validateNaiveDateTime(validDate, resourceName)

	require.Equal(t, []string([]string(nil)), warns)
	require.Equal(t, []error([]error(nil)), errs)

}

func TestValidateNaiveDateTime_RejectsNumericOffset(t *testing.T) {
	// It should reject any 8601 time with an offset
	invalidDate := "2021-02-21T15:30:00-08:00"
	resourceName := "schedule.0.end_repeat"

	warns, errs := validateNaiveDateTime(invalidDate, resourceName)
	expectedErrs := []error{errors.New("\"schedule.0.end_repeat\" of \"2021-02-21T15:30:00-08:00\" must be in the format 2006-01-02T15:04:05")}

	require.Equal(t, []string([]string(nil)), warns)
	require.Equal(t, expectedErrs, errs)

}

func TestValidateNaiveDateTime_RejectsGMTOffset(t *testing.T) {
	// It should reject an 8601 time with GMT designation
	invalidDate := "2021-02-21T15:30:00Z"
	resourceName := "schedule.0.end_repeat"

	warns, errs := validateNaiveDateTime(invalidDate, resourceName)
	expectedErrs := []error{errors.New("\"schedule.0.end_repeat\" of \"2021-02-21T15:30:00Z\" must be in the format 2006-01-02T15:04:05")}

	require.Equal(t, []string([]string(nil)), warns)
	require.Equal(t, expectedErrs, errs)

}

func TestValidateNaiveDateTime_RejectsUnixDateTime(t *testing.T) {
	// It should reject an 8601 time with GMT designation
	invalidDate := "123456789123456"
	resourceName := "schedule.0.end_repeat"

	warns, errs := validateNaiveDateTime(invalidDate, resourceName)
	expectedErrs := []error{errors.New("\"schedule.0.end_repeat\" of \"123456789123456\" must be in the format 2006-01-02T15:04:05")}

	require.Equal(t, []string([]string(nil)), warns)
	require.Equal(t, expectedErrs, errs)

}

func TestAccNewRelicAlertMutingRule_Basic(t *testing.T) {
	resourceName := "newrelic_alert_muting_rule.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertMutingRuleDestroy,
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

func TestAccNewRelicAlertMutingRule_BadInput(t *testing.T) {
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertMutingRuleDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicAlertMutingRuleBadInput(rName, "new muting rule", "product", "EQUALS", "APM", "BROWSER"),
				ExpectError: regexp.MustCompile("Validation Error: BAD_USER_INPUT"),
			},
		},
	})
}

func TestAccNewRelicAlertMutingRule_WithSchedule(t *testing.T) {
	resourceName := "newrelic_alert_muting_rule.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertMutingRuleDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicAlertMutingRuleWithSchedule(
					rName,
					"new muting rule",
					"product",
					"EQUALS",
					"APM",
					`
						start_time         = "2021-01-21T15:30:00"
						end_time           = "2021-01-21T16:30:00"
						time_zone          = "America/Los_Angeles"
						repeat             = "WEEKLY"
						end_repeat         = "2022-06-11T12:00:00"
						weekly_repeat_days = ["FRIDAY", "TUESDAY"]
					`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertMutingRuleExists(resourceName),
				),
			},
			//Test: Update to null out schedule completely
			{
				Config: testAccNewRelicAlertMutingRuleBasic(
					rName,
					"updated without schedule",
					"conditionName",
					"EQUALS",
					"My cool condition",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertMutingRuleExists(resourceName),
				),
			},
			{
				//Test: Update to add a schedule with default WEEKLY repeat (empty slice of days)
				Config: testAccNewRelicAlertMutingRuleWithSchedule(rName,
					"updated muting rule with schedule",
					"conditionType",
					"NOT_EQUALS",
					"baseline",
					`
						start_time         = "2021-02-21T15:30:00"
						end_time           = "2021-02-21T16:30:00"
						end_repeat         = "2022-06-11T12:00:00"
						repeat             = "WEEKLY"
						time_zone          = "America/Los_Angeles"
						weekly_repeat_days = []
					`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertMutingRuleExists(resourceName),
				),
			},
			{
				//Test: Update to add a schedule with DAILY repeat, new timezone & repeat_count
				Config: testAccNewRelicAlertMutingRuleWithSchedule(rName,
					"updated muting rule with schedule daily repeat, 42 times",
					"conditionType",
					"NOT_EQUALS",
					"baseline",
					`
						start_time         = "2021-02-21T15:30:00"
						end_time           = "2021-02-21T16:30:00"
						repeat_count       = 42
						repeat             = "DAILY"
						time_zone          = "Asia/Bangkok"
					`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertMutingRuleExists(resourceName),
				),
			},

			//Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
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

func testAccNewRelicAlertMutingRuleBadInput(
	name string,
	description string,
	attribute string,
	operator string,
	value1 string,
	value2 string,
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
			values 		= ["%[5]s", "%[6]s"]
		}
		conditions {
			attribute 	= "conditionType"
			operator 	= "%[4]s"
			values 		= ["static"]
		}
		operator = "AND"
	}
}
`, name, description, attribute, operator, value1, value2)
}

func testAccNewRelicAlertMutingRuleWithSchedule(
	name string,
	description string,
	attribute string,
	operator string,
	values string,
	schedule string,
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
	schedule { 
		%[6]s
	}
}`, name, description, attribute, operator, values, schedule)
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

func TestValidateMutingRuleConditionAttribute_ValidTags(t *testing.T) {
	// It should accept tags.value as well
	validTag := "tags.EC2Timezone"
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
	expectedErrs := []error{errors.New("\"condition.0.conditions.0.attribute\" of \"tag\" must be in the format tag.tag_name")}
	expectedWarns := []string([]string(nil))

	require.Equal(t, expectedWarns, warns)
	require.Equal(t, expectedErrs, errs)
}

func TestValidateMutingRuleConditionAttribute_entityGuid(t *testing.T) {
// It should accept entity.guid as an attribute
validAttr := "entity.guid"
resourceName := "condition.0.conditions.0.attribute"

warns, errs := validateMutingRuleConditionAttribute(validAttr, resourceName)
expectedErrs := []error(nil)
expectedWarns := []string([]string(nil))

require.Equal(t, expectedWarns, warns)
require.Equal(t, expectedErrs, errs)

}
