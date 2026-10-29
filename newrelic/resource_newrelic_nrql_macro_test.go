//go:build integration || NGEP

package newrelic

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/newrelic/newrelic-client-go/v2/pkg/nrqlmacros"
)

// TestAccNewRelicNrqlMacro_Basic tests the complete lifecycle (Create, Update, Import, Delete) of the resource.
func TestAccNewRelicNrqlMacro_Basic(t *testing.T) {
	t.Parallel()

	var (
		macroName       = fmt.Sprintf("testMacro%s", acctest.RandString(8))
		macroExpression = "environment = 'production' AND appName NOT LIKE '%test%'"
		resourceName    = "newrelic_nrql_macro.foo"

		macroExpressionUpdated = "environment = 'staging' AND appName NOT LIKE '%test%'"
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if v := os.Getenv("NEW_RELIC_ORGANIZATION_ID"); v == "" {
				t.Log("NEW_RELIC_ORGANIZATION_ID not set; organization_id will be auto-resolved from provider credentials")
			}
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlMacroDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNrqlMacroConfig(macroName, macroExpression),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlMacroExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", macroName),
					resource.TestCheckResourceAttr(resourceName, "expression", macroExpression),
					resource.TestCheckResourceAttrSet(resourceName, "organization_id"),
				),
			},
			// Test: Update expression
			{
				Config: testAccNrqlMacroConfig(macroName, macroExpressionUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlMacroExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", macroName),
					resource.TestCheckResourceAttr(resourceName, "expression", macroExpressionUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "organization_id"),
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

// testAccCheckNewRelicNrqlMacroDestroy verifies that the NRQL macro resource
// has been successfully destroyed in New Relic.
func testAccCheckNewRelicNrqlMacroDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "newrelic_nrql_macro" {
			continue
		}

		_, err := client.Nrqlmacros.GetEntityWithContext(context.Background(), rs.Primary.ID)
		if err != nil {
			return nil // Success: resource has been destroyed.
		}

		return fmt.Errorf("NRQL macro '%s' still exists", rs.Primary.ID)
	}

	return nil
}

// testAccCheckNewRelicNrqlMacroExists verifies that the NRQL macro resource
// has been created in New Relic by querying the API.
func testAccCheckNewRelicNrqlMacroExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found in state: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no NRQL macro ID is set in state")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		resp, err := client.Nrqlmacros.GetEntityWithContext(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching NRQL macro with ID %s: %w", rs.Primary.ID, err)
		}

		if resp == nil {
			return fmt.Errorf("entity with macroID %s returned nil response", rs.Primary.ID)
		}

		if _, ok := (*resp).(*nrqlmacros.EntityManagementNRQLMacroEntity); !ok {
			return fmt.Errorf("entity %s is not of type EntityManagementNRQLMacroEntity", rs.Primary.ID)
		}

		return nil
	}
}

// testAccNrqlMacroConfig is a helper function to generate the HCL configuration
// for a newrelic_nrql_macro resource.
func testAccNrqlMacroConfig(name, expression string) string {
	return fmt.Sprintf(`
resource "newrelic_nrql_macro" "foo" {
  name       = "%s"
  expression = "%s"
}
`, name, expression)
}
