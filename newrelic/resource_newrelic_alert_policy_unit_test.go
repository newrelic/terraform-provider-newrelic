//go:build unit
// +build unit

package newrelic

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicAlertPolicy_ErrorThrownWhenNameEmpty(t *testing.T) {
	avoidEmptyAccountID()
	expectedErrorMsg, _ := regexp.Compile(`name must not be empty`)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAlertPolicyConfigNameEmpty(),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}
