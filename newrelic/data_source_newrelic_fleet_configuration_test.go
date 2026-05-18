//go:build integration || FLEET

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNewRelicFleetConfigurationDataSource_ByConfigurationID fetches configuration
// content using the configuration GUID. Expects the latest version's content.
func TestAccNewRelicFleetConfigurationDataSource_ByConfigurationID(t *testing.T) {
	rName := fmt.Sprintf("tf-test-ds-id-%s", acctest.RandString(5))
	resourceName := "newrelic_fleet_configuration.src"
	dataSourceName := "data.newrelic_fleet_configuration.by_id"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetConfigDataSourceByID(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "configuration_content"),
					resource.TestCheckResourceAttrSet(dataSourceName, "organization_id"),
					// Data source content must match the resource's latest version content
					resource.TestCheckResourceAttrPair(
						dataSourceName, "configuration_content",
						resourceName, "version.0.configuration_content",
					),
				),
			},
		},
	})
}

// TestAccNewRelicFleetConfigurationDataSource_ByVersionEntityID fetches configuration
// content for a specific version using the version entity GUID.
func TestAccNewRelicFleetConfigurationDataSource_ByVersionEntityID(t *testing.T) {
	rName := fmt.Sprintf("tf-test-ds-ver-%s", acctest.RandString(5))
	dataSourceName := "data.newrelic_fleet_configuration.by_version"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetConfigDataSourceByVersionID(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "configuration_content"),
					resource.TestCheckResourceAttrSet(dataSourceName, "organization_id"),
				),
			},
		},
	})
}

// TestAccNewRelicFleetConfigurationDataSource_ByName searches for a configuration by name
// and verifies the data source resolves to the correct entity.
func TestAccNewRelicFleetConfigurationDataSource_ByName(t *testing.T) {

	rName := fmt.Sprintf("tf-test-ds-name-%s", acctest.RandString(5))
	dataSourceName := "data.newrelic_fleet_configuration.by_name"
	resourceName := "newrelic_fleet_configuration.src"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetConfigDataSourceByName(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "configuration_content"),
					resource.TestCheckResourceAttrSet(dataSourceName, "organization_id"),
					// Name-based lookup must resolve the same content as the resource
					resource.TestCheckResourceAttrPair(
						dataSourceName, "configuration_content",
						resourceName, "version.0.configuration_content",
					),
				),
			},
		},
	})
}

// TestAccNewRelicFleetConfigurationDataSource_NotFound verifies that querying a
// non-existent configuration by name returns an appropriate error.
func TestAccNewRelicFleetConfigurationDataSource_NotFound(t *testing.T) {
	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckFleetEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccFleetConfigDataSourceNotFound(),
				ExpectError: regexp.MustCompile(`no fleet configuration found with name`),
			},
		},
	})
}

// TestAccNewRelicFleetConfigurationDataSource_MutuallyExclusive verifies that specifying
// two lookup inputs simultaneously is rejected at plan time.
func TestAccNewRelicFleetConfigurationDataSource_MutuallyExclusive(t *testing.T) {
	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckFleetEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccFleetConfigDataSourceConflict(),
				ExpectError: regexp.MustCompile(`conflicts with`),
			},
		},
	})
}

// ── config templates ──────────────────────────────────────────────────────────

func testAccFleetConfigDataSourceByID(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "src" {
  name                = %q
  agent_type          = "NRInfra"
  managed_entity_type = "HOST"

  version {
    configuration_content = "log:\n  level: info\n# ds-by-id\n"
  }
}

data "newrelic_fleet_configuration" "by_id" {
  configuration_id = newrelic_fleet_configuration.src.configuration_id
}
`, name)
}

func testAccFleetConfigDataSourceByVersionID(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "src" {
  name                = %q
  agent_type          = "NRInfra"
  managed_entity_type = "HOST"

  version {
    configuration_content = "log:\n  level: info\n# ds-by-version\n"
  }
}

data "newrelic_fleet_configuration" "by_version" {
  version_entity_id = newrelic_fleet_configuration.src.latest_version_entity_id
}
`, name)
}

func testAccFleetConfigDataSourceByName(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "src" {
  name                = %q
  agent_type          = "NRInfra"
  managed_entity_type = "HOST"

  version {
    configuration_content = "log:\n  level: info\n# ds-by-name\n"
  }
}

data "newrelic_fleet_configuration" "by_name" {
  name = newrelic_fleet_configuration.src.name
  depends_on = [newrelic_fleet_configuration.src]
}
`, name)
}

func testAccFleetConfigDataSourceNotFound() string {
	return `
data "newrelic_fleet_configuration" "missing" {
  name = "tf-test-does-not-exist-zzzzzzzz"
}
`
}

func testAccFleetConfigDataSourceConflict() string {
	return `
data "newrelic_fleet_configuration" "conflict" {
  configuration_id  = "some-guid"
  version_entity_id = "other-guid"
}
`
}
