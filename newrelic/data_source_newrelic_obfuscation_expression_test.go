//go:build integration || LOGGING_INTEGRATIONS

package newrelic

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicObfuscationExpressionDataSource_Basic(t *testing.T) {
	name := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//read
			{
				Config: testAccNewRelicObfuscationExpressionDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicObfuscationExpressionDataSourceExists("data.newrelic_obfuscation_expression.exp")),
			},
		},
	})
}

func testAccNewRelicObfuscationExpressionDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_obfuscation_expression" "foo"{
	account_id = %[1]d
	name = "%[2]s"
	description = "%[3]s"
	regex = "(^http)"
}

data "newrelic_obfuscation_expression" "exp"{
	account_id = %[1]d
	name = newrelic_obfuscation_expression.foo.name
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccCheckNewRelicObfuscationExpressionDataSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		_, err := getObfuscationExpressionByID(context.Background(), client, testAccountID, rs.Primary.ID)
		if err != nil {
			return err
		}

		return nil
	}
}
