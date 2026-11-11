//go:build integration || WORKLOADS

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicMaintenanceWindow_Basic(t *testing.T) {
	resourceName := "newrelic_maintenance_window.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicMaintenanceWindowDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicMaintenanceWindowConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicMaintenanceWindowExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "duration", "PT1H"),
					resource.TestCheckResourceAttr(resourceName, "affected_entity_type", "SERVICE_LEVEL"),
					resource.TestCheckResourceAttr(resourceName, "affected_entities.#", "1"),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicMaintenanceWindowConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicMaintenanceWindowExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "duration", "PT2H"),
				),
			},
		},
	})
}

func testAccNewRelicMaintenanceWindowConfig(name string) string {
	return fmt.Sprintf(`
data "newrelic_entity" "app" {
	name   = "%[3]s"
	domain = "APM"
	type   = "APPLICATION"
}

resource "newrelic_service_level" "mw_sli" {
	guid = data.newrelic_entity.app.guid
	name = "%[2]s-sli"

	events {
		account_id = %[1]d
		valid_events {
			from = "Transaction"
		}
		good_events {
			from = "Transaction"
		}
	}

	objective {
		target = 99.00
		time_window {
			rolling {
				count = 7
				unit  = "DAY"
			}
		}
	}
}

resource "newrelic_maintenance_window" "foo" {
	name       = "%[2]s"
	scope_id   = "%[1]d"
	scope_type = "ACCOUNT"

	start_time = "2026-12-01T22:00:00"
	duration   = "PT1H"
	timezone   = "UTC"

	affected_entity_type = "SERVICE_LEVEL"
	affected_entities    = [newrelic_service_level.mw_sli.sli_guid]
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicMaintenanceWindowConfigUpdated(name string) string {
	return fmt.Sprintf(`
data "newrelic_entity" "app" {
	name   = "%[3]s"
	domain = "APM"
	type   = "APPLICATION"
}

resource "newrelic_service_level" "mw_sli" {
	guid = data.newrelic_entity.app.guid
	name = "%[2]s-sli"

	events {
		account_id = %[1]d
		valid_events {
			from = "Transaction"
		}
		good_events {
			from = "Transaction"
		}
	}

	objective {
		target = 99.00
		time_window {
			rolling {
				count = 7
				unit  = "DAY"
			}
		}
	}
}

resource "newrelic_maintenance_window" "foo" {
	name        = "%[2]s"
	description = "updated by acceptance test"
	scope_id    = "%[1]d"
	scope_type  = "ACCOUNT"

	start_time = "2026-12-01T22:00:00"
	duration   = "PT2H"
	timezone   = "UTC"

	affected_entity_type = "SERVICE_LEVEL"
	affected_entities    = [newrelic_service_level.mw_sli.sli_guid]
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccCheckNewRelicMaintenanceWindowExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no maintenance window ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		resp, err := client.MaintenanceWindow.GetListByIds([]string{rs.Primary.ID})
		if err != nil {
			return err
		}
		if resp == nil || len(resp.MaintenanceWindows) != 1 {
			return fmt.Errorf("maintenance window not found: %v", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckNewRelicMaintenanceWindowDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_maintenance_window" {
			continue
		}

		resp, err := client.MaintenanceWindow.GetListByIds([]string{r.Primary.ID})
		if err == nil && resp != nil && len(resp.MaintenanceWindows) > 0 {
			return fmt.Errorf("maintenance window still exists: %v", r.Primary.ID)
		}
	}

	return nil
}
