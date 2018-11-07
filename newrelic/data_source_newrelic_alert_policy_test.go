package newrelic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNewRelicAlertPolicyDataSource_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNewRelicAlertPolicyDataSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicAlertPolicy("data.newrelic_alert_policy.policy"),
				),
			},
		},
	})
}

func testAccNewRelicAlertPolicy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get an alert policy from New Relic")
		}

		if strings.Contains(strings.ToLower(testAccExpectedAlertPolicyName), strings.ToLower(a["name"])) {
			return fmt.Errorf("Expected the alert policy name to be: %s, but got: %s", testAccExpectedAlertPolicyName, a["name"])
		}

		return nil
	}
}

func testAccNewRelicAlertPolicyDataSourceConfig(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%s"
}

data "newrelic_alert_policy" "policy" {
	name = "${newrelic_alert_policy.foo.name}"
}
`, rName)
}
