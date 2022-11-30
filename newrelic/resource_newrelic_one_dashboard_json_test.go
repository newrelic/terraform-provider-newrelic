//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
)

// TestAccNewRelicOneDashboardRaw_Create Ensure that we can create a NR1 Dashboard
func TestAccNewRelicOneDashboardJson_Create(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicOneDashboardJsonConfig_OnePageFull(rName, strconv.Itoa(testAccountID)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard_json.bar", 0),
				),
			},
			// Import
			{
				ResourceName: "newrelic_one_dashboard_json.bar",
				ImportState:  true,
			},
		},
	})
}

// testAccCheckNewRelicOneDashboardRawConfig contains all the config options for a single page dashboard
func testAccCheckNewRelicOneDashboardJsonConfig_OnePageFull(dashboardName string, accountID string) string {
	return `
resource "newrelic_one_dashboard_json" "bar" {
  json = <<EOT
  	` + testAccCheckNewRelicOneDashboardJsonConfig_Full(dashboardName, accountID) + `
  EOT
}`
}

// testAccCheckNewRelicOneDashboardRawConfig_PageFull generates a TF config snippet that is
// an entire dashboard page, with all widget types
func testAccCheckNewRelicOneDashboardJsonConfig_Full(pageName string, accountID string) string {
	return `
	{
		"name": "APM and Infrastructure",
		"description": "",
		"permissions": "PUBLIC_READ_WRITE",
		"pages": [
		  {
			"name": "APM and Infrastructure",
			"description": "",
			"widgets": [
			  {
				"title": "Application Names 1234",
				"layout": {
				  "column": 1,
				  "row": 1,
				  "width": 4,
				  "height": 3
				},
				"linkedEntityGuids": [
				  "MTYwNjg2MnxWSVp8REFTSEJPQVJEfDUwNjQyNQ"
				],
				"visualization": {
				  "id": "viz.bar"
				},
				"rawConfiguration": {
				  "nrqlQueries": [
					{
					  "accountId": ` + accountID + `,
					  "query": "SELECT average(duration) FROM Transaction,ProcessSample   facet cases( where appName = 'WebPortal' OR 'nr.apmApplicationNames' LIKE '%WebPortal%' as 'Web Portal' ,WHERE appName = 'Billing Service' OR 'nr.apmApplicationNames' LIKE '%Billing Service%' as 'Billing Service', WHERE appName ='Fulfillment Service' OR 'nr.apmApplicationNames' LIKE '%Fulfillment%' as 'Fulfillment', WHERE appName = 'Plan Service' OR 'nr.apmApplicationNames' like '%Plan Service%' as 'Plan Service' )"
					}
				  ]
				}
			  }
			]
		  }
		],
		"variables": []
	}
	`
}

// testAccCheckNewRelicOneDashboardDestroy expects the dashboard read to fail,
// and errors if we DO get the dashboard back.
func testAccCheckNewRelicOneDashboardJsonDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_one_dashboard_json" {
			continue
		}

		_, err := client.Dashboards.GetDashboardEntity(common.EntityGUID(r.Primary.ID))
		if err == nil {
			return fmt.Errorf("newrelic_one_dashboard_json still exists")
		}

	}
	return nil
}
