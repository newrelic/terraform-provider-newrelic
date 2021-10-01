//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/common"
)

func TestAccNewRelicServiceLevel_Basic(t *testing.T) {
	resourceName := "newrelic_service_level.sli"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
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
			// TODO
			{
				Config: testAccNewRelicServiceLevelConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicServiceLevelExists(resourceName),
				),
			},
			// Test: Import
			// TODO
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNewRelicServiceLevelConfig(name string) string {
	return fmt.Sprintf(`

resource "newrelic_workload" "workload" {
	name = "%[2]s"
	account_id = %[1]d
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
	scope_account_ids =  [%[1]d]
}

resource "newrelic_service_level" "sli" {
	guid = newrelic_workload.workload.guid
	name = "%[2]s-updated"
	
	events {
		account_id = %[1]d
		valid_events {
			from = "Transaction"
		}
		good_events {
			from = "Transaction"
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

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		indicators, err := client.ServiceLevel.GetIndicators(common.EntityGUID(identifier.GUID))
		if err != nil {
			return err
		}

		for _, indicator := range *indicators {
			if identifier.ID == indicator.ID {
				return nil
			}
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

		_, err = client.ServiceLevel.GetIndicators(common.EntityGUID(identifier.GUID))
		if err == nil {
			return fmt.Errorf("SLI still exists")
		}
	}

	return nil
}
