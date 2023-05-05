//go:build integration
// +build integration

package newrelic

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Checking the creation, update, import and deletion of obfuscation expression
func TestAccNewRelicObfuscationExpression_Basic(t *testing.T) {
	resourceName := "newrelic_obfuscation_expression.foo"
	rName := acctest.RandString(7)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicObfuscationExpressionDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config: testAccNewRelicObfuscationExpressionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicObfuscationExpressionExists(resourceName)),
			},
			//update
			{
				Config: testAccNewRelicObfuscationExpressionUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicObfuscationExpressionExists(resourceName)),
			},
			//import
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
		},
	})
}

// Must fail if given the same name
func TestAccNewRelicObfuscationExpression_Validation(t *testing.T) {
	rName := acctest.RandString(7)
	expectedMsg, _ := regexp.Compile("Invalid input: There is another obfuscation expression with the same name in this account")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicObfuscationExpressionDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config:      testAccNewRelicObfuscationExpression_ValidateName(rName),
				ExpectError: expectedMsg,
			},
		},
	})
}

func testAccCheckNewRelicObfuscationExpressionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "newrelic_obfuscation_expression" {
			continue
		}
		_, err := getObfuscationExpressionByID(context.Background(), client, testAccountID, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Obfuscation expression still exists: %s", err)
		}
	}

	return nil
}

func testAccNewRelicObfuscationExpressionConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_obfuscation_expression" "foo"{
	account_id = %[1]d
	name = "%[2]s"
	description = "%[3]s"
	regex = "(^http)"
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccCheckNewRelicObfuscationExpressionExists(n string) resource.TestCheckFunc {
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

func testAccNewRelicObfuscationExpressionUpdate(name string) string {
	return fmt.Sprintf(`
resource "newrelic_obfuscation_expression" "foo"{
	account_id = %[1]d
	name = "%[2]s-updated"
	description = "%[3]s-updated"
	regex = "(^http)"
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicObfuscationExpression_ValidateName(name string) string {
	return fmt.Sprintf(`
resource "newrelic_obfuscation_expression" "foo"{
	account_id = %[1]d
	name = "%[2]s"
	description = "%[3]s"
	regex = "(^http)"
}

resource "newrelic_obfuscation_expression" "bar"{
	account_id = %[1]d
	name = "%[2]s"
	description = "%[3]s"
	regex = "(^http)"
}
`, testAccountID, name, testAccExpectedApplicationName)
}
