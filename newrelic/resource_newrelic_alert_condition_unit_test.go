// +build unit

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicAlertCondition_ShortTermDuration(t *testing.T) {
	avoidEmptyAccountID()
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	expectedErrorMsg, _ := regexp.Compile(`expected term.0.duration to be in the range \(5 - 120\)`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAlertConditionConfigDuration(rName, 4),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicAlertCondition_LongTermDuration(t *testing.T) {
	avoidEmptyAccountID()
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	expectedErrorMsg, _ := regexp.Compile(`expected term.0.duration to be in the range \(5 - 120\)`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAlertConditionConfigDuration(rName, 121),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicAlertCondition_LongName(t *testing.T) {
	avoidEmptyAccountID()
	expectedErrorMsg, _ := regexp.Compile(`expected length of name to be in the range \(1 \- 128\)`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAlertConditionConfig("really-long-name-longer-than-one-hundred-and-twenty-eight-characters-so-it-causes-an-error-because-really-long-name-causes-an-error"),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicAlertCondition_EmptyName(t *testing.T) {
	avoidEmptyAccountID()
	expectedErrorMsg, _ := regexp.Compile(`name must not be empty`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAlertConditionConfig(""),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}
