//go:build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicFleet_Basic(t *testing.T) {
	resourceName := "newrelic_fleet.foo"
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicFleetConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "managed_entity_type", "HOST"),
					resource.TestCheckResourceAttr(resourceName, "operating_system", "LINUX"),
				),
			},
			// Update
			{
				Config: testAccNewRelicFleetConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-updated", rName)),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicFleet_Kubernetes(t *testing.T) {
	resourceName := "newrelic_fleet.k8s"
	rName := fmt.Sprintf("tf-test-k8s-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicFleetConfigKubernetes(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "managed_entity_type", "KUBERNETESCLUSTER"),
					resource.TestCheckNoResourceAttr(resourceName, "operating_system"),
				),
			},
		},
	})
}

func testAccCheckNewRelicFleetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no fleet ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		entity, err := client.FleetControl.GetEntity(rs.Primary.ID)
		if err != nil {
			return err
		}

		if entity == nil {
			return fmt.Errorf("fleet not found: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccNewRelicFleetConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet" "foo" {
  name                = "%s"
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
  description         = "Test fleet"
}
`, name)
}

func testAccNewRelicFleetConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet" "foo" {
  name                = "%s-updated"
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
  description         = "Updated description"
}
`, name)
}

func testAccNewRelicFleetConfigKubernetes(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet" "k8s" {
  name                = "%s"
  managed_entity_type = "KUBERNETESCLUSTER"
  description         = "Test Kubernetes fleet"
}
`, name)
}
