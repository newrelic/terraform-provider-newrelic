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

//Checking the creation, update, import and deletion of obfuscation expression
func TestAccNewRelicObfuscationRule_Basic(t *testing.T) {
	resourceName := "newrelic_obfuscation_rule.foo"
	rName := acctest.RandString(7)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicObfuscationRuleDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config: testAccNewRelicObfuscationRuleConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicObfuscationRuleExists(resourceName)),
			},
			//update
			{
				Config: testAccNewRelicObfuscationRuleUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicObfuscationRuleExists(resourceName)),
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

// Checking for same name error
func TestAccNewRelicObfuscationRule_NameValid(t *testing.T) {
	rName := acctest.RandString(7)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicObfuscationRuleDestroy,
		Steps: []resource.TestStep{
			//Create with the same name
			{
				Config:      testAccNewRelicObfuscationRule_validName(rName),
				ExpectError: regexp.MustCompile(` There is another obfuscation rule with the same name in this account.`),
			},
		},
	})
}

//if the obfuscation expression is updated
func TestAccNewRelicObfuscationRule_ExpressionUpdate(t *testing.T) {
	resourceName := "newrelic_obfuscation_rule.foo"
	rName := acctest.RandString(7)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicObfuscationRuleDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config: testAccNewRelicObfuscationRuleConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicObfuscationRuleExists(resourceName)),
			},
			//update
			{
				Config: testAccNewRelicObfuscationRule_ExpressionUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicObfuscationRuleExists(resourceName)),
			},
		},
	})
}

//actions updated
func TestAccNewRelicObfuscationRule_ActionsUpdate(t *testing.T) {
	resourceName := "newrelic_obfuscation_rule.foo"
	rName := acctest.RandString(7)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicObfuscationRuleDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config: testAccNewRelicObfuscationRuleConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicObfuscationRuleExists(resourceName)),
			},
			//update
			{
				Config: testAccNewRelicObfuscationRule_ActionsUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicObfuscationRuleExists(resourceName)),
			},
		},
	})
}

func testAccCheckNewRelicObfuscationRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "newrelic_obfuscation_rule" {
			continue
		}
		_, err := getObfuscationRuleByID(context.Background(), client, testAccountID, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Obfuscation rule still exists: %s", err)
		}
	}

	return nil
}
func testAccCheckNewRelicObfuscationRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		_, err := getObfuscationRuleByID(context.Background(), client, testAccountID, rs.Primary.ID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccNewRelicObfuscationRuleConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_obfuscation_expression" "bar"{
account_id = %[1]d
name = "%[2]s"
description = "%[3]s"
regex = "(^http)"
}

resource "newrelic_obfuscation_rule" "foo"{
account_id = %[1]d
name = "%[2]s"
description = "%[3]s"
filter = "hostStatus=running"
enabled = true
action{
	attribute = ["message"]
	expression_id = newrelic_obfuscation_expression.bar.id
	method = "MASK"
}
action{
	attribute = ["hostName,ip"]
	expression_id = newrelic_obfuscation_expression.bar.id
	method = "HASH_SHA256"
}
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicObfuscationRuleUpdate(name string) string {
	return fmt.Sprintf(`
resource "newrelic_obfuscation_expression" "bar"{
	account_id = %[1]d
	name = "%[2]s"
	description = "%[3]s"
	regex = "(^http)"
}

resource "newrelic_obfuscation_rule" "foo"{
account_id = %[1]d
name = "%[2]s-update"
description = "%[3]s-update"
filter = "hostStatus=running"
enabled = true
action{
	attribute = ["message"]
	expression_id = newrelic_obfuscation_expression.bar.id
	method = "MASK"
}
action{
	attribute = ["hostName,ip"]
	expression_id = newrelic_obfuscation_expression.bar.id
	method = "HASH_SHA256"
}
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicObfuscationRule_ExpressionUpdate(name string) string {
	return fmt.Sprintf(`
resource "newrelic_obfuscation_expression" "bar"{
	account_id = %[1]d
	name = "%[2]s-update"
	description = "%[3]s-update"
	regex = "(^http)"
}

resource "newrelic_obfuscation_rule" "foo"{
account_id = %[1]d
name = "%[2]s"
description = "%[3]s"
filter = "hostStatus=running"
enabled = true
action{
	attribute = ["hostName,ip"]
	expression_id = newrelic_obfuscation_expression.bar.id
	method = "HASH_SHA256"
}
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicObfuscationRule_validName(name string) string {
	return fmt.Sprintf(`
resource "newrelic_obfuscation_expression" "bar"{
account_id = %[1]d
name = "%[2]s"
description = "%[3]s"
regex = "(^http)"
}

resource "newrelic_obfuscation_rule" "foo"{
account_id = %[1]d
name = "%[2]s"
description = "%[3]s"
filter = "hostStatus=running"
enabled = true
action{
	attribute = ["message"]
	expression_id = newrelic_obfuscation_expression.bar.id
	method = "MASK"
}
}

resource "newrelic_obfuscation_rule" "foo1"{
account_id = %[1]d
name = "%[2]s"
description = "%[3]s"
filter = "hostStatus=running"
enabled = true
action{
	attribute = ["message"]
	expression_id = newrelic_obfuscation_expression.bar.id
	method = "MASK"
}
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicObfuscationRule_ActionsUpdate(name string) string {
	return fmt.Sprintf(`
resource "newrelic_obfuscation_expression" "bar"{
account_id = %[1]d
name = "%[2]s"
description = "%[3]s"
regex = "(^http)"
}

resource "newrelic_obfuscation_rule" "foo"{
account_id = %[1]d
name = "%[2]s"
description = "%[3]s"
filter = "hostStatus=running"
enabled = true
action{
	attribute = ["message-update"]
	expression_id = newrelic_obfuscation_expression.bar.id
	method = "MASK"
}
}
`, testAccountID, name, testAccExpectedApplicationName)
}
