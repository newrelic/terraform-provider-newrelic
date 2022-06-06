//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

func TestAccNewRelicSyntheticsPrivateLocation_Basic(t *testing.T) {
	resourceName := "newrelic_synthetics_private_location.bar"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsPrivateLocationDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsPrivateLocationConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsPrivateLocationExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsPrivateLocationConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsPrivateLocationExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_id", "description", "domain_id", "guid", "key", "location_id", "name", "verified_script_execution"},
			},
		},
	})
}

func testAccCheckNewRelicSyntheticsPrivateLocationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no synthetics private location is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		queryString := fmt.Sprintf("domain = 'SYNTH' AND type = 'PRIVATE_LOCATION' AND name = '%s'", rs.Primary.Attributes["name"])
		time.Sleep(10 * time.Second)

		found, err := client.Entities.GetEntitySearchByQuery(entities.EntitySearchOptions{}, queryString, []entities.EntitySearchSortCriteria{})
		if err != nil {
			return err
		}

		for _, e := range found.Results.Entities {
			if !strings.EqualFold(string(e.GetGUID()), rs.Primary.ID) {
				return fmt.Errorf("synthetics private location not found: %v - %v", rs.Primary.ID, found)
			}
		}

		return nil
	}
}

func testAccCheckNewRelicSyntheticsPrivateLocationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_synthetics_private_location" {
			continue
		}

		queryString := fmt.Sprintf("domain = 'SYNTH' AND type = 'PRIVATE_LOCATION' AND name = '%s'", r.Primary.Attributes["name"])

		time.Sleep(10 * time.Second)

		found, err := client.Entities.GetEntitySearchByQuery(entities.EntitySearchOptions{}, queryString, []entities.EntitySearchSortCriteria{})
		if err != nil {
			return err
		}

		if found.Count != 0 {
			return fmt.Errorf("synthetics private location still exists")
		}

	}
	return nil
}

func testAccNewRelicSyntheticsPrivateLocationConfig(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_synthetics_private_location" "bar" {
	account_id	 = 2520528
	description  = "Test Description"
	name		 =	"%[1]s"
	verified_script_execution = true
}
`, name)
}

func testAccNewRelicSyntheticsPrivateLocationConfigUpdated(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_synthetics_private_location" "bar" {
	account_id	 = 2520528
	description  = "Test Description-Updated"
	name		 =	"%[1]s-updated"
	verified_script_execution = false
}
`, name)
}
