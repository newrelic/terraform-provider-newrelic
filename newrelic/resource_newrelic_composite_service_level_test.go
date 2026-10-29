//go:build integration || WORKLOADS

package newrelic

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
)

func TestAccNewRelicCompositeServiceLevel_Basic(t *testing.T) {
	resourceName := "newrelic_composite_service_level.composite"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCompositeServiceLevelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicCompositeServiceLevelConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCompositeServiceLevelExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicCompositeServiceLevelConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCompositeServiceLevelExists(resourceName),
				),
			},
			// Test: Import — unlike newrelic_service_level, composite's ID is
			// just its own bare GUID, so this works cleanly with no custom
			// identifier encoding to account for.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNewRelicCompositeServiceLevelConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_workload" "workload" {
	name       = "%[2]s"
	account_id = %[1]d
	entity_search_query {
		query = "tags.namespace like '%%App%%' "
	}
	scope_account_ids = [%[1]d]
}

resource "newrelic_service_level" "component_1" {
	guid = newrelic_workload.workload.guid
	name = "%[2]s-component-1"

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
		target = 99
		time_window {
			rolling {
				count = 7
				unit  = "DAY"
			}
		}
	}
}

resource "newrelic_service_level" "component_2" {
	guid = newrelic_workload.workload.guid
	name = "%[2]s-component-2"

	events {
		account_id = %[1]d
		valid_events {
			from = "Transaction"
		}
		good_events {
			from  = "Transaction"
			where = "error IS FALSE"
		}
	}

	objective {
		target = 99
		time_window {
			rolling {
				count = 7
				unit  = "DAY"
			}
		}
	}
}

resource "newrelic_composite_service_level" "composite" {
	account_id  = %[1]d
	name        = "%[2]s"
	description = "Composite SLI acceptance test"

	component {
		guid = newrelic_service_level.component_1.sli_guid
	}

	component {
		guid = newrelic_service_level.component_2.sli_guid
	}

	objective {
		target = 99
		time_window {
			rolling {
				count = 7
				unit  = "DAY"
			}
		}
	}
}
`, testAccountID, name)
}

func testAccNewRelicCompositeServiceLevelConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_workload" "workload" {
	name       = "%[2]s"
	account_id = %[1]d
	entity_search_query {
		query = "tags.namespace like '%%App%%' "
	}
	scope_account_ids = [%[1]d]
}

resource "newrelic_service_level" "component_1" {
	guid = newrelic_workload.workload.guid
	name = "%[2]s-component-1"

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
		target = 99
		time_window {
			rolling {
				count = 7
				unit  = "DAY"
			}
		}
	}
}

resource "newrelic_service_level" "component_2" {
	guid = newrelic_workload.workload.guid
	name = "%[2]s-component-2"

	events {
		account_id = %[1]d
		valid_events {
			from = "Transaction"
		}
		good_events {
			from  = "Transaction"
			where = "error IS FALSE"
		}
	}

	objective {
		target = 99
		time_window {
			rolling {
				count = 7
				unit  = "DAY"
			}
		}
	}
}

resource "newrelic_composite_service_level" "composite" {
	account_id  = %[1]d
	name        = "%[2]s-updated"
	description = "Composite SLI acceptance test (updated)"

	component {
		guid   = newrelic_service_level.component_1.sli_guid
		weight = 60
	}

	component {
		guid   = newrelic_service_level.component_2.sli_guid
		weight = 40
	}

	objective {
		target = 99
		time_window {
			rolling {
				count = 7
				unit  = "DAY"
			}
		}
	}
}
`, testAccountID, name)
}

func testAccCheckNewRelicCompositeServiceLevelExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no composite service level ID is set")
		}

		time.Sleep(30 * time.Second)

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		composites, err := client.ServiceLevel.GetComposites([]common.EntityGUID{common.EntityGUID(rs.Primary.ID)})
		if err != nil {
			return err
		}

		if len(*composites) == 1 && string((*composites)[0].GUID) == rs.Primary.ID {
			return nil
		}

		return fmt.Errorf("composite service level not found: %v", rs.Primary.ID)
	}
}

func testAccCheckNewRelicCompositeServiceLevelDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_composite_service_level" {
			continue
		}

		_, err := client.ServiceLevel.GetComposites([]common.EntityGUID{common.EntityGUID(r.Primary.ID)})
		if err == nil {
			return fmt.Errorf("composite service level still exists")
		}
	}

	return nil
}
