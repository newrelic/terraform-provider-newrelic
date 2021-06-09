// +build integration

package newrelic

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicPluginsAlertCondition_NameGreaterThan64Char(t *testing.T) {
	avoidEmptyAccountID()
	expectedErrorMsg, _ := regexp.Compile(`expected length of name to be in the range \(1 \- 64\)`)
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicPluginsAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicPluginsAlertConditionNameGreaterThan64Char(rName),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicPluginsAlertCondition_NameLessThan1Char(t *testing.T) {
	avoidEmptyAccountID()
	expectedErrorMsg, _ := regexp.Compile(`expected length of name to be in the range \(1 \- 64\)`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicPluginsAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicPluginsAlertConditionNameLessThan1Char(),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicPluginsAlertCondition_TermDurationGreaterThan120(t *testing.T) {
	avoidEmptyAccountID()
	expectedErrorMsg, _ := regexp.Compile(`expected term.0.duration to be in the range \(5 - 120\)`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicPluginsAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicPluginsAlertConditionTermDurationGreaterThan120(),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicPluginsAlertCondition_TermDurationLessThan5(t *testing.T) {
	avoidEmptyAccountID()
	expectedErrorMsg, _ := regexp.Compile(`expected term.0.duration to be in the range \(5 - 120\)`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicPluginsAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicPluginsAlertConditionTermDurationLessThan5(),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}
