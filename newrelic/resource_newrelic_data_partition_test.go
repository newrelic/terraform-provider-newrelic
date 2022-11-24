//go:build integration
// +build integration

package newrelic

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

//Checking the creation, update, import and deletion of obfuscation expression
func TestAccNewRelicDataPartitionRule_Basic(t *testing.T) {
	resourceName := "newrelic_data_partition_rule.foo"
	rName := acctest.RandString(7)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicDataPartitionRuleDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config: testAccNewRelicDataPartitionRuleConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDataPartitionRuleExists(resourceName)),
			},
			////update
			//{
			//	Config: testAccNewRelicDataPartitionRuleUpdate(rName),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheckNewRelicDataPartitionRuleExists(resourceName)),
			//},
			////import
			//{
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	ResourceName:      resourceName,
			//},
		},
	})
}

////Must fail if given the same name
//func TestAccNewRelicDataPartitionRule_Validation(t *testing.T) {
//	rName := acctest.RandString(7)
//	expectedMsg, _ := regexp.Compile("Invalid input: There is another obfuscation expression with the same name in this account")
//	resource.ParallelTest(t, resource.TestCase{
//		PreCheck:     func() { testAccPreCheck(t) },
//		Providers:    testAccProviders,
//		CheckDestroy: testAccCheckNewRelicDataPartitionRuleDestroy,
//		Steps: []resource.TestStep{
//			//create
//			{
//				Config:      testAccNewRelicDataPartitionRule_ValidateName(rName),
//				ExpectError: expectedMsg,
//			},
//		},
//	})
//}

func testAccCheckNewRelicDataPartitionRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "newrelic_data_partition_rule" {
			continue
		}
		_, err := getDataPartitionByID(context.Background(), client, testAccountID, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("data partition rule still exists: %s", err)
		}
	}

	return nil
}

func testAccNewRelicDataPartitionRuleConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_data_partition_rule" "foo"{
	account_id = %[1]d
	description = "%[3]s"
	enabled = true
    attribute_name = "hostname"
    matching_expression = "localhost"
    matching_method = "EQUALS"
    retention_policy = "SECONDARY"
    target_data_partition = "Log_Test_%[2]s"
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccCheckNewRelicDataPartitionRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		_, err := getDataPartitionByID(context.Background(), client, testAccountID, rs.Primary.ID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccNewRelicDataPartitionRuleUpdate(name string) string {
	return fmt.Sprintf(`
resource "newrelic_obfuscation_expression" "foo"{
	account_id = %[1]d
	name = "%[2]s-updated"
	description = "%[3]s-updated"
	regex = "(^http)"
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicDataPartitionRule_ValidateName(name string) string {
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
