//go:build integration

package newrelic

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
)

func TestAccNewRelicSyntheticsCertCheckMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_cert_check_monitor.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsCertCheckMonitorResourceDestroy,
		Steps: []resource.TestStep{
			//Create
			{
				Config: testAccNewRelicSyntheticsCertCheckMonitorConfig(rName, "EVERY_5_MINUTES", "ENABLED", 30),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicSyntheticsCertCheckMonitorExists(resourceName),
				),
			},
			// Update
			{
				Config: testAccNewRelicSyntheticsCertCheckMonitorConfig(fmt.Sprintf("%s-updated", rName), "EVERY_10_MINUTES", "DISABLED", 20),
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
					"locations_private",
					"certificate_expiration",
					"domain",
					"tag",
				},
			},
		},
	})
}

func testAccNewRelicSyntheticsCertCheckMonitorConfig(name string, period string, status string, certExp int) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_cert_check_monitor" "foo" {
	name                   = "%[1]s"
	domain                 = "newrelic.com"
	period                 = "%[2]s"
	status                 = "%[3]s"
	certificate_expiration = %[4]d
	locations_public       = ["AP_SOUTH_1"]
	tag {
		key    = "cars"
		values = ["audi"]
	}
}
`, name, period, status, certExp)
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
		time.Sleep(60 * time.Second)

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

		// Unfortunately we still have to wait due to async delay with entity indexing :(
		time.Sleep(60 * time.Second)

		found, _ := client.Entities.GetEntity(common.EntityGUID(r.Primary.ID))
		if (*found) != nil {
			return fmt.Errorf("synthetics monitor still exists")
		}
	}
	return nil
}
