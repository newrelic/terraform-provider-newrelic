//go:build integration || ALERTS
// +build integration ALERTS

package newrelic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicAlertCompoundConditionDataSource_Basic(t *testing.T) {
	resourceName := "newrelic_alert_compound_condition.foo"
	dataSourceName := "data.newrelic_alert_compound_condition.bar"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertCompoundConditionDataSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "policy_id", resourceName, "policy_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "trigger_expression", resourceName, "trigger_expression"),
					resource.TestCheckResourceAttr(dataSourceName, "component_conditions.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "facet_matching_behavior", "FACETS_IGNORED"),
				),
			},
		},
	})
}

func testAccNewRelicAlertCompoundConditionDataSourceConfig(name string) string {
	return testAccNewRelicAlertCompoundConditionConfigBasic(name, "A AND B") + `
data "newrelic_alert_compound_condition" "bar" {
	id = newrelic_alert_compound_condition.foo.id
}
`
}
