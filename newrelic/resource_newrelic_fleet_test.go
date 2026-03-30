//go:build integration || FLEET

package newrelic

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicFleet_Basic(t *testing.T) {
	resourceName := "newrelic_fleet.foo"
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicFleetConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "managed_entity_type", "HOST"),
					resource.TestCheckResourceAttr(resourceName, "operating_system", "LINUX"),
					resource.TestCheckResourceAttr(resourceName, "description", "Test fleet"),
					resource.TestCheckResourceAttrSet(resourceName, "organization_id"),
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

func TestAccNewRelicFleet_Windows(t *testing.T) {
	resourceName := "newrelic_fleet.windows"
	rName := fmt.Sprintf("tf-test-win-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicFleetConfigWindows(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "managed_entity_type", "HOST"),
					resource.TestCheckResourceAttr(resourceName, "operating_system", "WINDOWS"),
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

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicFleetConfigKubernetes(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "managed_entity_type", "KUBERNETESCLUSTER"),
					resource.TestCheckNoResourceAttr(resourceName, "operating_system"),
				),
			},
			// Update
			{
				Config: testAccNewRelicFleetConfigKubernetesUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated Kubernetes fleet"),
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

func TestAccNewRelicFleet_WithTags(t *testing.T) {
	resourceName := "newrelic_fleet.tags"
	rName := fmt.Sprintf("tf-test-tags-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicFleetConfigWithTags(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
				),
			},
		},
	})
}

func TestAccNewRelicFleet_WithProduct(t *testing.T) {
	resourceName := "newrelic_fleet.product"
	rName := fmt.Sprintf("tf-test-prod-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicFleetConfigWithProduct(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "product", "INFRA"),
				),
			},
		},
	})
}

// Error case tests

func TestAccNewRelicFleet_MissingOperatingSystemForHost(t *testing.T) {
	rName := fmt.Sprintf("tf-test-error-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicFleetConfigMissingOS(rName),
				ExpectError: regexp.MustCompile("operating_system is required when managed_entity_type is HOST"),
			},
		},
	})
}

func TestAccNewRelicFleet_OperatingSystemForKubernetes(t *testing.T) {
	rName := fmt.Sprintf("tf-test-error-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicFleetConfigKubernetesWithOS(rName),
				ExpectError: regexp.MustCompile("operating_system should not be specified for KUBERNETESCLUSTER fleets"),
			},
		},
	})
}

// Helper functions

func setupFleetTestCredentials(t *testing.T) {
	t.Helper()

	// Set fleet credentials for this test
	originalAPIKey := os.Getenv("NEW_RELIC_API_KEY")
	originalAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	t.Cleanup(func() {
		os.Setenv("NEW_RELIC_API_KEY", originalAPIKey)
		os.Setenv("NEW_RELIC_ACCOUNT_ID", originalAccountID)
	})

	fleetAPIKey := os.Getenv("NEW_RELIC_FLEET_TEST_API_KEY")
	fleetAccountID := os.Getenv("NEW_RELIC_FLEET_TEST_ACCOUNT_ID")
	if fleetAPIKey != "" {
		os.Setenv("NEW_RELIC_API_KEY", fleetAPIKey)
	}
	if fleetAccountID != "" {
		os.Setenv("NEW_RELIC_ACCOUNT_ID", fleetAccountID)
	}
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

func testAccCheckNewRelicFleetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_fleet" {
			continue
		}

		// Try to get the fleet - should return error if not found
		entity, err := client.FleetControl.GetEntity(r.Primary.ID)
		if err == nil && entity != nil {
			return fmt.Errorf("fleet still exists: %s", r.Primary.ID)
		}
	}

	return nil
}

// Config functions

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

func testAccNewRelicFleetConfigWindows(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet" "windows" {
  name                = "%s"
  managed_entity_type = "HOST"
  operating_system    = "WINDOWS"
  description         = "Test Windows fleet"
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

func testAccNewRelicFleetConfigKubernetesUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet" "k8s" {
  name                = "%s"
  managed_entity_type = "KUBERNETESCLUSTER"
  description         = "Updated Kubernetes fleet"
}
`, name)
}

func testAccNewRelicFleetConfigWithTags(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet" "tags" {
  name                = "%s"
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
  description         = "Test fleet with tags"
  tags                = ["environment:test", "team:platform"]
}
`, name)
}

func testAccNewRelicFleetConfigWithProduct(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet" "product" {
  name                = "%s"
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
  description         = "Test fleet with product"
  product             = "INFRA"
}
`, name)
}

func testAccNewRelicFleetConfigMissingOS(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet" "error" {
  name                = "%s"
  managed_entity_type = "HOST"
  description         = "This should fail - missing operating_system"
}
`, name)
}

func testAccNewRelicFleetConfigKubernetesWithOS(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet" "error" {
  name                = "%s"
  managed_entity_type = "KUBERNETESCLUSTER"
  operating_system    = "LINUX"
  description         = "This should fail - operating_system not allowed for K8s"
}
`, name)
}
