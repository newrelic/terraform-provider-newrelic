//go:build integration

package newrelic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"testing"
	"time"
)

func TestAccNewRelicSyntheticsCertCheckMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_cert_check_monitor.foo"
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsCertCheckMonitorResourceDestroy,
		Steps: []resource.TestStep{
			//Create
			{
				Config: testAccNewRelicSyntheticsCertCheckMonitorConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicSyntheticsCertCheckMonitorExists(resourceName),
				),
			},
			// Update
			{
				PreConfig: func() {
					// Unfortunately we still have to wait due to async delay with entity indexing :(
					time.Sleep(10 * time.Second)
				},
				Config: testAccNewRelicSyntheticsCertCheckMonitorConfigUpdated(fmt.Sprintf("%s-updated", rName)),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicSyntheticsCertCheckMonitorExists(resourceName),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"locations_public",
					"locations_private",
					"certificate_expiration",
					"domain",
					"tag",
				},
			},
		},
	})
}

func testAccNewRelicSyntheticsCertCheckMonitorConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_cert_check_monitor" "foo" {
  name="%[1]s"
  domain="newrelic.com"
  period="EVERY_5_MINUTES"
  status="ENABLED"
  certificate_expiration=30
  location_public=["AP_SOUTH_1"]
  tag{
    key="cars"
    values=["audi"]
  }
}
`, name)
}

func testAccNewRelicSyntheticsCertCheckMonitorConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_cert_check_monitor" "foo" {
  name="%[1]s-updated"
  domain="newrelic.com"
  period="EVERY_MINUTE"
  status="DISABLED"
  certificate_expiration=20
  location_public=["AP_SOUTH_1","AP_EAST_1"]
  tag{
    key="cars"
    values=["audi","BMW"]
  }
}
`, name)
}

func testAccNewRelicSyntheticsCertCheckMonitorExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no synthetics monitor ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		// We also have to wait for the monitor's deletion to be indexed as well :(
		time.Sleep(5 * time.Second)

		result, err := client.Entities.GetEntity(common.EntityGUID(rs.Primary.ID))
		if err != nil {
			return err
		}
		if string((*result).GetGUID()) != rs.Primary.ID {
			return fmt.Errorf("the monitor is not found %v - %v", (*result).GetGUID(), rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckNewRelicSyntheticsCertCheckMonitorResourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_synthetics_cert_check_monitor" {
			continue
		}

		found, _ := client.Entities.GetEntity(common.EntityGUID(r.Primary.ID))
		if (*found) != nil {
			return fmt.Errorf("synthetics monitor still exists")
		}
	}
	return nil
}
