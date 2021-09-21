//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicAlertPolicyDataSource_Basic(t *testing.T) {
	resourceName := "newrelic_alert_policy.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertPolicyDataSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertPolicyDataSource("data.newrelic_alert_policy.policy"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-test-%s", rName)),
				),
			},
		},
	})
}

func TestAccNewRelicAlertPolicyDataSource_NameExactMatchOnly(t *testing.T) {
	rName := acctest.RandString(5)
	expectedErrorMsg := regexp.MustCompile(`the name '.*' does not match any New Relic alert policy`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAlertPolicyDataSourceConfigNameExactMatchOnly(rName),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func testAccNewRelicAlertPolicyDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%s"
}

data "newrelic_alert_policy" "policy" {
	name = newrelic_alert_policy.foo.name
}
`, name)
}

func testAccNewRelicAlertPolicyDataSourceConfigNameExactMatchOnly(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%s"
}

data "newrelic_alert_policy" "policy" {
	name = "tf-test-%s"
	depends_on = [newrelic_alert_policy.foo]
}
`, name, name[:len(name)-1])
}

func testAccCheckNewRelicAlertPolicyDataSource(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("expected to get an alert policy from New Relic")
		}

		return nil
	}
}
