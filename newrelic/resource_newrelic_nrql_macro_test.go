//go:build integration || NGEP

package newrelic

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/newrelic/newrelic-client-go/v2/pkg/nrqlmacros"
)

// TestAccNewRelicNrqlMacro_Basic covers the full lifecycle of a non-parameterized macro:
// create, verify, update expression, verify again, then import.
func TestAccNewRelicNrqlMacro_Basic(t *testing.T) {
	t.Parallel()

	macroName := fmt.Sprintf("testMacro%s", acctest.RandString(8))
	resourceName := "newrelic_nrql_macro.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlMacroDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNrqlMacroConfig(macroName, "environment = 'production'"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlMacroExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", macroName),
					resource.TestCheckResourceAttr(resourceName, "expression", "environment = 'production'"),
					resource.TestCheckResourceAttrSet(resourceName, "organization_id"),
					resource.TestCheckResourceAttr(resourceName, "parameter_count", "0"),
				),
			},
			{
				Config: testAccNrqlMacroConfig(macroName, "environment = 'staging'"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlMacroExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "expression", "environment = 'staging'"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccNewRelicNrqlMacro_Parameterized verifies that parameter_count is accepted
// and enforced correctly for expressions that use positional parameters.
func TestAccNewRelicNrqlMacro_Parameterized(t *testing.T) {
	t.Parallel()

	macroName := fmt.Sprintf("testMacroP%s", acctest.RandString(8))
	resourceName := "newrelic_nrql_macro.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlMacroDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNrqlMacroConfigParameterized(macroName, "status = $1 AND region = $2", 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlMacroExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "parameter_count", "2"),
				),
			},
			// Update to remove parameters (go from 2 to 0)
			{
				Config: testAccNrqlMacroConfig(macroName, "environment = 'production'"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlMacroExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "parameter_count", "0"),
				),
			},
		},
	})
}

func testAccCheckNewRelicNrqlMacroDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "newrelic_nrql_macro" {
			continue
		}

		_, err := client.Nrqlmacros.GetEntityWithContext(context.Background(), rs.Primary.ID)
		if err != nil {
			return nil
		}

		return fmt.Errorf("NRQL macro '%s' still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckNewRelicNrqlMacroExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no NRQL macro ID is set in state")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		resp, err := client.Nrqlmacros.GetEntityWithContext(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching NRQL macro %s: %w", rs.Primary.ID, err)
		}

		if resp == nil {
			return fmt.Errorf("NRQL macro %s returned nil from the API", rs.Primary.ID)
		}

		if _, ok := (*resp).(*nrqlmacros.EntityManagementNRQLMacroEntity); !ok {
			return fmt.Errorf("entity %s is not a NRQL macro", rs.Primary.ID)
		}

		return nil
	}
}

func testAccNrqlMacroConfig(name, expression string) string {
	return fmt.Sprintf(`
resource "newrelic_nrql_macro" "test" {
  name       = %q
  expression = %q
}
`, name, expression)
}

func testAccNrqlMacroConfigParameterized(name, expression string, parameterCount int) string {
	return fmt.Sprintf(`
resource "newrelic_nrql_macro" "test" {
  name            = %q
  expression      = %q
  parameter_count = %d
}
`, name, expression, parameterCount)
}
