//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"testing"
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
				ImportStateVerifyIgnore: []string{"description", "domain_id", "key", "location_id", "verified_script_execution"},
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

		found, err := client.Entities.GetEntity(common.EntityGUID(rs.Primary.ID))
		if err != nil {
			return err
		}

		if string((*found).GetGUID()) != rs.Primary.ID {
			fmt.Errorf("the private location was not found %v - %v", (*found).GetGUID(), rs.Primary.ID)
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

		found, err := client.Entities.GetEntity(common.EntityGUID(r.Primary.ID))
		if err != nil {
			return err
		}

		if (*found) != nil {
			fmt.Errorf("private location still exists")
		}
	}
	return nil
}

func testAccNewRelicSyntheticsPrivateLocationConfig(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_synthetics_private_location" "bar" {
		description               = "Test Description-Updated"
		name                      = "%[1]s"
		verified_script_execution = false
}
`, name)
}

func testAccNewRelicSyntheticsPrivateLocationConfigUpdated(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_synthetics_private_location" "bar" {
		description               = "Test Description-Updated"
		name                      = "%[1]s"
		verified_script_execution = false
}
`, name)
}
