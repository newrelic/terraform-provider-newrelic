//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

func TestAccNewRelicServiceLevel_Basic(t *testing.T) {
	resourceName := "newrelic_service_level.sli"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicServiceLevelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicServiceLevelConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicServiceLevelExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicServiceLevelConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicServiceLevelExists(resourceName),
				),
			},
		},
	})
}

func testAccNewRelicServiceLevelConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_workload" "workload" {
	name = "%[2]s"
	account_id = %[1]d
	entity_search_query {
		query = "tags.namespace like '%%App%%' "
	}
	scope_account_ids =  [%[1]d]
}

resource "newrelic_service_level" "sli" {
	guid = newrelic_workload.workload.guid
	name = "%[2]s"

	events {
		account_id = %[1]d
		valid_events {
			from = "Transaction"
		}
		good_events {
			from = "Transaction"
			select {
				attribute = "duration"
				function = "COUNT"
			}
		}
	}

	objective {
		target = 99.00
		time_window {
			rolling {
				count = 7
				unit = "DAY"
			}
		}
	}
}

resource "newrelic_service_level" "cdf_sli" {
  guid = newrelic_workload.workload.guid
  name = "%[2]s using cdf"

  events {
    account_id = %[1]d
    valid_events {
      from = "Metric"
      select {
        attribute = "query.wallClockTime.negative.distribution"
        function = "GET_FIELD"
      }
      where = "metricName = 'query.wallClockTime.negative.distribution'"
    }

    good_events {
      from = "Metric"
      select {
        attribute = "query.wallClockTime.negative.distribution"
        function =  "GET_CDF_COUNT"
        threshold = 7
      }
      where = "metricName = 'query.wallClockTime.negative.distribution'"
    }
  }

  objective {
    target = 49
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

func testAccNewRelicServiceLevelConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_workload" "workload" {
	name = "%[2]s"
	account_id = %[1]d
	entity_search_query {
		query = "tags.namespace like '%%App%%' "
	}
	scope_account_ids =  [%[1]d]
}

resource "newrelic_service_level" "sli" {
	guid = newrelic_workload.workload.guid
	name = "%[2]s-updated"

	events {
		account_id = %[1]d
		valid_events {
			from = "Transaction"
			select {
				attribute = "duration"
				function = "SUM"
			}
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
				unit = "DAY"
			}
		}
	}
}
`, testAccountID, name)
}

func testAccCheckNewRelicServiceLevelExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No SLI ID is set")
		}

		identifier, err := parseIdentifier(rs.Primary.ID)
		if err != nil {
			return err
		}

		time.Sleep(30 * time.Second)

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		indicators, err := client.ServiceLevel.GetIndicators(common.EntityGUID(getSliGUID(identifier)))
		if err != nil {
			return err
		}

		if len(*indicators) == 1 && (*indicators)[0].ID == identifier.ID {
			return nil
		}

		return fmt.Errorf("SLI not found: %v", rs.Primary.ID)
	}
}

func testAccCheckNewRelicServiceLevelDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_service_level" {
			continue
		}

		identifier, err := parseIdentifier(r.Primary.ID)
		if err != nil {
			return err
		}

		_, err = client.ServiceLevel.GetIndicators(common.EntityGUID(getSliGUID(identifier)))
		if err == nil {
			return fmt.Errorf("SLI still exists")
		}
	}

	return nil
}

func cleanupDanglingServiceLevelResources() error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	query := "domain = 'EXT' AND type = 'SERVICE_LEVEL' AND (name LIKE '%tf-test-%' OR name LIKE '%tf_test_%')"

	fmt.Printf("\n[INFO] cleaning up any dangling integration test resources... \n")
	time.Sleep(1 * time.Second)

	matches, err := client.Entities.GetEntitySearchByQuery(
		entities.EntitySearchOptions{},
		query,
		[]entities.EntitySearchSortCriteria{},
	)

	if err != nil {
		return fmt.Errorf("error cleaning up dangling synthetics resources: %s", err)
	}

	if matches != nil {
		resources := matches.Results.Entities
		for _, r := range resources {
			_, err := client.ServiceLevel.ServiceLevelDelete(common.EntityGUID(string(r.GetGUID())))
			if err != nil {
				log.Printf("[ERROR] error deleting dangling resource: %s", err)
			}
		}

		fmt.Printf("\n[INFO] deleted %d dangling resources \n", len(resources))
	}

	return nil
}
