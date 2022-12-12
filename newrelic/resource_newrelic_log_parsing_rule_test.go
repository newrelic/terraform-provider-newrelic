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

//Checking the creation, update, import and deletion of logging parsing rule
func TestAccNewRelicLogParsingRule_Basic(t *testing.T) {
	resourceName := "newrelic_log_parsing_rule.foo"
	rName := acctest.RandString(7)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicLogParsingRuleDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config: testAccNewRelicLogParsingRuleConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicLogParsingRuleExists(resourceName)),
			},
			//update
			/*{
				Config: testAccNewRelicLogParsingRuleUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicLogParsingRuleExists(resourceName)),
			},
			//import
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},*/
		},
	})
}

func testAccCheckNewRelicLogParsingRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "newrelic_log_parsing_rule" {
			continue
		}
		_, err := getLogParsingRuleByID(context.Background(), client, testAccountID, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Log parsing rule still exists: %s", err)
		}
	}

	return nil
}

func testAccNewRelicLogParsingRuleConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_log_parsing_rule" "foo"{
	account_id = %[1]d
	name = "%[2]s"
	attribute = "%[3]s"
	enabled     = true
    grok        = "sampleattribute='%%{NUMBER:test:int}'"
    lucene      = "logtype:linux_messages"
    nrql        = "SELECT * FROM Log WHERE logtype = 'linux_messages'"
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccCheckNewRelicLogParsingRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		_, err := getLogParsingRuleByID(context.Background(), client, testAccountID, rs.Primary.ID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccNewRelicLogParsingRuleUpdate(name string) string {
	return fmt.Sprintf(`
resource "newrelic_log_parsing_rule" "foo"{
	account_id = %[1]d
	name = "%[2]s"+ "_update"
	attribute = "%[3]s"
	enabled     = false
    grok        = "sampleattribute=%%{NUMBER:test:int}"
    lucene      = "logtype:linux_messages"
    nrql        = "SELECT * FROM Log WHERE logtype = 'linux_messages'"                                                
}
`, testAccountID, name, testAccExpectedApplicationName)
}
