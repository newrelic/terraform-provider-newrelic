//go:build integration
// +build integration

package newrelic

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Checking the creation, update, import and deletion of logging parsing rule
func TestAccNewRelicLogParsingRule_Basic(t *testing.T) {
	resourceName := "newrelic_log_parsing_rule.foo"
	rName := generateNameForIntegrationTestResource()
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
			{
				Config: testAccNewRelicLogParsingRuleUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicLogParsingRuleExists(resourceName)),
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

func TestAccNewRelicLogParsingRule_Unique_Name_Update(t *testing.T) {
	resourceName := "newrelic_log_parsing_rule.foo"
	expectedErrorMsg := regexp.MustCompile("name is already in use by another rule")
	rName1 := generateNameForIntegrationTestResource()
	rName2 := generateNameForIntegrationTestResource()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicLogParsingRuleDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config: testAccNewRelicLogParsingRuleConfigCreateUniqueName(rName1, rName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicLogParsingRuleExists(resourceName)),
			},
			{
				Config:      testAccNewRelicLogParsingRuleConfigUpdateUniqueName(rName1),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicLogParsingRule_Unique_Name_Create(t *testing.T) {
	expectedErrorMsg := regexp.MustCompile("name is already in use by another rule")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicLogParsingRuleDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config:      testAccNewRelicLogParsingRuleUniqueNameConfig("test_tf_DescUnique"),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicLogParsingRule_Grok_Test_Matched(t *testing.T) {
	resourceName := "newrelic_log_parsing_rule.foo"
	rName := generateNameForIntegrationTestResource()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicLogParsingRuleDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config: testAccNewRelicLogParsingRuleGrokConfigMatched(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicLogParsingRuleExists(resourceName)),
			},
		},
	})
}

func TestAccNewRelicLogParsingRule_Grok_Test_Unmatched(t *testing.T) {
	resourceName := "newrelic_log_parsing_rule.foo"
	rName := generateNameForIntegrationTestResource()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicLogParsingRuleDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config: testAccNewRelicLogParsingRuleGrokConfigUnmatched(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicLogParsingRuleExists(resourceName)),
			},
		},
	})
}

func TestAccNewRelicLogParsingRule_Invalid_Grok_Test(t *testing.T) {
	rName := generateNameForIntegrationTestResource()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicLogParsingRuleDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config:      testAccNewRelicLogParsingRuleInvalidGrokConfig(rName),
				ExpectError: regexp.MustCompile("Invalid Grok pattern"),
			},
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
	account_id  = %[1]d
	name        = "%[2]s"
	attribute   = "%[3]s"
	enabled     = true
	grok        = "sampleattribute='%%%%{NUMBER:test:int}'"
	lucene      = "logtype:linux_messages"
	nrql        = "SELECT * FROM Log WHERE logtype = 'linux_messages'"
}
`, testAccountID, name, testAccExpectedApplicationName)
}
func testAccNewRelicLogParsingRuleConfigCreateUniqueName(name1 string, name2 string) string {
	return fmt.Sprintf(`
resource "newrelic_log_parsing_rule" "foo"{
	account_id  = %[1]d
	name        = "%[2]s"
	attribute   = "%[3]s"
	enabled     = true
	grok        = "sampleattribute='%%%%{NUMBER:test:int}'"
	lucene      = "logtype:linux_messages"
	nrql        = "SELECT * FROM Log WHERE logtype = 'linux_messages'"
}
resource "newrelic_log_parsing_rule" "bar"{
	account_id  = %[1]d
	name        = "%[4]s"
	attribute   = "%[3]s"
	enabled     = true
	grok        = "sampleattribute='%%%%{NUMBER:test:int}'"
	lucene      = "logtype:linux_messages"
	nrql        = "SELECT * FROM Log WHERE logtype = 'linux_messages'"
}
`, testAccountID, name1, testAccExpectedApplicationName, name2)
}

func testAccNewRelicLogParsingRuleConfigUpdateUniqueName(name1 string) string {
	return fmt.Sprintf(`
resource "newrelic_log_parsing_rule" "foo"{
	account_id  = %[1]d
	name        = "%[2]s"
	attribute   = "%[3]s"
	enabled     = true
	grok        = "sampleattribute='%%%%{NUMBER:test:int}'"
	lucene      = "logtype:linux_messages"
	nrql        = "SELECT * FROM Log WHERE logtype = 'linux_messages'"
}
resource "newrelic_log_parsing_rule" "bar"{
	account_id  = %[1]d
	name        = "%[2]s"
	attribute   = "%[3]s"
	enabled     = true
	grok        = "sampleattribute='%%%%{NUMBER:test:int}'"
	lucene      = "logtype:linux_messages"
	nrql        = "SELECT * FROM Log WHERE logtype = 'linux_messages'"
}
`, testAccountID, name1, testAccExpectedApplicationName)
}

func testAccNewRelicLogParsingRuleUniqueNameConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_log_parsing_rule" "foo"{
	account_id  = %[1]d
	name        = "%[2]s"
	attribute   = "%[3]s"
	enabled     = true
	grok        = "sampleattribute='%%%%{NUMBER:test:int}'"
	lucene      = "logtype:linux_messages"
	nrql        = "SELECT * FROM Log WHERE logtype = 'linux_messages'"
}
resource "newrelic_log_parsing_rule" "bar"{
	account_id  = %[1]d
	name        = "%[2]s"
	attribute   = "%[3]s"
	enabled     = true
	grok        = "sampleattribute='%%%%{NUMBER:test:int}'"
	lucene      = "logtype:linux_messages"
	nrql        = "SELECT * FROM Log WHERE logtype = 'linux_messages'"
	depends_on  = ["newrelic_log_parsing_rule.foo"
  ]

}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicLogParsingRuleGrokConfigMatched(name string) string {
	return fmt.Sprintf(`
data "newrelic_test_grok_pattern" "grok"{
	account_id  = %[1]d
	grok        = "%%%%{IP:host_ip}"
	log_lines   = ["host_ip: 43.3.120.2"]
}
resource "newrelic_log_parsing_rule" "foo"{
	account_id  = %[1]d
	name        = "%[2]s"
	attribute   = "%[3]s"
	enabled     = true
	grok        = data.newrelic_test_grok_pattern.grok.grok
	lucene      = "logtype:linux_messages"
	nrql        = "SELECT * FROM Log WHERE logtype = 'linux_messages'"
	matched     = data.newrelic_test_grok_pattern.grok.test_grok[0].matched
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicLogParsingRuleGrokConfigUnmatched(name string) string {
	return fmt.Sprintf(`
data "newrelic_test_grok_pattern" "grok"{
	account_id  = %[1]d
	grok        = "%%%%{IP:host_ip}"
	log_lines   = ["bytes_received: 2048"]
}
resource "newrelic_log_parsing_rule" "foo"{
	account_id  = %[1]d
	name        = "%[2]s"
	attribute   = "%[3]s"
	enabled     = true
	grok        = data.newrelic_test_grok_pattern.grok.grok
	lucene      = "logtype:linux_messages"
	nrql        = "SELECT * FROM Log WHERE logtype = 'linux_messages'"
	matched     = data.newrelic_test_grok_pattern.grok.test_grok[0].matched
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicLogParsingRuleInvalidGrokConfig(name string) string {
	return fmt.Sprintf(`
data "newrelic_test_grok_pattern" "grok"{
	account_id  = %[1]d
	grok        = "{IP:host_ip}"
	log_lines   = ["host_ip: 43.3.120.2","bytes_received: 2048"]
}
resource "newrelic_log_parsing_rule" "foo"{
	account_id  = %[1]d
	name        = "%[2]s"
	attribute   = "%[3]s"
	enabled     = true
	grok        = data.newrelic_test_grok_pattern.grok.grok
	lucene      = "logtype:linux_messages"
	nrql        = "SELECT * FROM Log WHERE logtype = 'linux_messages'"
	matched     = data.newrelic_test_grok_pattern.grok.test_grok[0].matched
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
	account_id  = %[1]d
	name        = "%[2]s_update"
	attribute   = "%[3]s"
	enabled     = false
	grok        = "sampleattribute='%%%%{NUMBER:test:int}'"
	lucene      = "logtype:linux_messages"
	nrql        = "SELECT * FROM Log WHERE logtype = 'linux_messages'"                                                
}
`, testAccountID, name, testAccExpectedApplicationName)
}
