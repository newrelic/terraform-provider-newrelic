//go:build integration || WORKLOADS

package newrelic

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicMaintenanceWindow_Basic(t *testing.T) {
	resourceName := "newrelic_maintenance_window.foo"
	rName := generateNameForIntegrationTestResource()
	// Always schedule the window in the future so the test doesn't expire.
	startTime := time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02T15:04:05")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicMaintenanceWindowDestroy,
		Steps: []resource.TestStep{
			// Step 0: create only the Service Level dependency, then wait for it to be
			// indexed before any maintenance window references it as an affected entity.
			// Without this wait, maintenanceWindowCreate intermittently fails with a
			// generic "Internal server error" for an SLI that was created moments earlier.
			{
				Config: testAccNewRelicMaintenanceWindowSLIOnlyConfig(rName),
				Check: func(s *terraform.State) error {
					time.Sleep(20 * time.Second)
					return nil
				},
			},
			// Test: Create
			{
				Config: testAccNewRelicMaintenanceWindowConfig(rName, startTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicMaintenanceWindowExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "duration", "PT1H"),
					resource.TestCheckResourceAttr(resourceName, "affected_entity_type", "SERVICE_LEVEL"),
					resource.TestCheckResourceAttr(resourceName, "affected_entities.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "start_time", startTime),
					resource.TestCheckResourceAttr(resourceName, "scope.0.type", "ACCOUNT"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.id", strconv.Itoa(testAccountID)),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicMaintenanceWindowConfigUpdated(rName, startTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicMaintenanceWindowExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "duration", "PT2H"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "affected_entities.#", "1"),
					// The refresh that follows this step reads the window back via
					// GetListByIds, which can briefly return the pre-update values;
					// give the update time to become visible so the plan stays empty.
					func(s *terraform.State) error {
						time.Sleep(30 * time.Second)
						return nil
					},
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccNewRelicMaintenanceWindowSLIOnlyConfig creates just the Service Level
// dependency used as the maintenance window's affected entity, without the
// maintenance window itself. Its "newrelic_service_level" "mw_sli" block must
// stay identical to the one in testAccNewRelicMaintenanceWindowConfig so that
// applying that config next doesn't cause a diff/replace of this resource.
func testAccNewRelicMaintenanceWindowSLIOnlyConfig(name string) string {
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
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicMaintenanceWindowConfig(name, startTime string) string {
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
	description = "created by acceptance test"
	scope {
		id   = "%[1]d"
		type = "ACCOUNT"
	}

	start_time = "%[4]s"
	duration   = "PT1H"
	timezone   = "UTC"

	affected_entity_type = "SERVICE_LEVEL"
	affected_entities    = [newrelic_service_level.mw_sli.sli_guid]
}
`, testAccountID, name, testAccExpectedApplicationName, startTime)
}

func testAccNewRelicMaintenanceWindowConfigUpdated(name, startTime string) string {
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
	scope {
		id   = "%[1]d"
		type = "ACCOUNT"
	}

	start_time = "%[4]s"
	duration   = "PT2H"
	timezone   = "UTC"

	affected_entity_type = "SERVICE_LEVEL"
	affected_entities    = [newrelic_service_level.mw_sli.sli_guid]
}
`, testAccountID, name, testAccExpectedApplicationName, startTime)
}

// testAccCheckNewRelicMaintenanceWindowExists polls GetListByIds for up to a
// minute rather than checking once: a maintenance window that was just
// created is not always immediately resolvable by this query, the same
// indexing delay already accounted for elsewhere in this test.
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

		return resource.RetryContext(context.Background(), 60*time.Second, func() *resource.RetryError {
			resp, err := client.MaintenanceWindow.GetListByIds([]string{rs.Primary.ID})
			if err != nil {
				return resource.RetryableError(err)
			}
			if resp == nil || len(resp.MaintenanceWindows) != 1 {
				return resource.RetryableError(fmt.Errorf("maintenance window not found yet: %v", rs.Primary.ID))
			}

			return nil
		})
	}
}

// testAccCheckNewRelicMaintenanceWindowDestroy mirrors the same retry
// tolerance as the exists check, in case a just-deleted window briefly still
// shows up in GetListByIds before the index catches up.
func testAccCheckNewRelicMaintenanceWindowDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_maintenance_window" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), 60*time.Second, func() *resource.RetryError {
			resp, err := client.MaintenanceWindow.GetListByIds([]string{r.Primary.ID})
			if err == nil && resp != nil && len(resp.MaintenanceWindows) > 0 {
				return resource.RetryableError(fmt.Errorf("maintenance window still exists: %v", r.Primary.ID))
			}

			return nil
		})
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
