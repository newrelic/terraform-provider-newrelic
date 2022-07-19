//go:build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/common"
)

func TestAccNewRelicSyntheticsBrokenLinksMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_broken_links_monitor.foo"
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicSyntheticsBrokenLinksMonitorConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorEntityExists(resourceName),
				),
			},
			// Update
			{
				Config: testAccNewRelicSyntheticsBrokenLinksMonitorConfig(fmt.Sprintf("%s-updated", rName)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorEntityExists(resourceName),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true, //name,type
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"period",
					"locations_public",
					"locations_private",
					"status",
					"tag",
					"uri",
				},
			},
		},
	})
}

func testAccNewRelicSyntheticsBrokenLinksMonitorConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_broken_links_monitor" "foo" {
  name	=	"%[1]s"
  period	=	"EVERY_HOUR"
  status	=	"ENABLED"
  locations_public	=	["Mumbai, IN"]
  uri = "https://www.google.com"

  tag {
    key	= "tf-test"
    values	= ["tf-acc-test"]
  }
}`, name)
}

func testAccCheckNewRelicSyntheticsMonitorEntityExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no synthetics monitor ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		result, err := client.Entities.GetEntity(common.EntityGUID(rs.Primary.ID))
		if err != nil {
			return err
		}
		if string((*result).GetGUID()) != rs.Primary.ID {
			fmt.Errorf("the monitor is not found %v - %v", (*result).GetGUID(), rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckNewRelicSyntheticsMonitorResourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_synthetics_broken_links_monitor" {
			continue
		}

		found, _ := client.Entities.GetEntity(common.EntityGUID(r.Primary.ID))
		if (*found) != nil {
			return fmt.Errorf("synthetics monitor still exists")
		}
	}
	return nil
}
