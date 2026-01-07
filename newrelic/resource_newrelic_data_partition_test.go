//go:build integration || LOGGING_INTEGRATIONS

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

// Checking the creation, update, import and deletion of data partition rule
func TestAccNewRelicDataPartitionRule_Basic(t *testing.T) {
	resourceName := "newrelic_data_partition_rule.foo"
	rName := acctest.RandString(7)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccLogDataPartitionsCleanup(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicDataPartitionRuleDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config: testAccNewRelicDataPartitionRuleConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDataPartitionRuleExists(resourceName)),
			},
			//update
			{
				Config: testAccNewRelicDataPartitionRuleUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDataPartitionRuleExists(resourceName)),
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

// Checking the creation of new resource on update name
func TestAccNewRelicDataPartitionRule_NameUpdate(t *testing.T) {
	resourceName := "newrelic_data_partition_rule.foo"
	rName := acctest.RandString(7)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccLogDataPartitionsCleanup(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicDataPartitionRuleDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config: testAccNewRelicDataPartitionRuleConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDataPartitionRuleExists(resourceName)),
			},
			//update
			{
				Config: testAccNewRelicDataPartitionRuleUpdate_Name(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDataPartitionRuleExists(resourceName)),
			},
		},
	})
}

// Checking the creation of new resource on update NRQL
func TestAccNewRelicDataPartitionRule_NRQLUpdate(t *testing.T) {
	resourceName := "newrelic_data_partition_rule.foo"
	rName := acctest.RandString(7)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccLogDataPartitionsCleanup(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicDataPartitionRuleDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config: testAccNewRelicDataPartitionRuleConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDataPartitionRuleExists(resourceName)),
			},
			//update
			{
				Config: testAccNewRelicDataPartitionRuleUpdate_NRQL(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDataPartitionRuleExists(resourceName)),
			},
		},
	})
}

// Must fail if given the same name
func TestAccNewRelicDataPartitionRule_DuplicateName(t *testing.T) {
	rName := acctest.RandString(7)
	expectedMsg, _ := regexp.Compile("DUPLICATE_DATA_PARTITION_RULE_NAME")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccLogDataPartitionsCleanup(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicDataPartitionRuleDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config:      testAccNewRelicDataPartitionRule_DuplicateName(rName),
				ExpectError: expectedMsg,
			},
		},
	})
}

// Must fail if given the invalid name
func TestAccNewRelicDataPartitionRule_Validation(t *testing.T) {
	rName := acctest.RandString(7)
	expectedMsg, _ := regexp.Compile("INVALID_DATA_PARTITION_INPUT")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccLogDataPartitionsCleanup(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicDataPartitionRuleDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config:      testAccNewRelicDataPartitionRule_ValidateName(rName),
				ExpectError: expectedMsg,
			},
		},
	})
}

func testAccCheckNewRelicDataPartitionRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "newrelic_data_partition_rule" {
			continue
		}
		_, err := getDataPartitionByID(context.Background(), client, testAccountID, rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("data partition rule still exists: %s", err)
		}
	}

	return nil
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

func testAccNewRelicDataPartitionRuleConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_data_partition_rule" "foo"{
	account_id = %[1]d
	description = "%[3]s"
	enabled = true
	nrql = "logtype='node'"
    retention_policy = "SECONDARY"
    target_data_partition = "Log_Test_%[2]s"
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicDataPartitionRuleUpdate(name string) string {
	return fmt.Sprintf(`
resource "newrelic_data_partition_rule" "foo"{
	account_id = %[1]d
	description = "%[3]s_update"
	enabled = true
	nrql = "logtype='linux_messages'"
    retention_policy = "SECONDARY"
    target_data_partition = "Log_Test_%[2]s"
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicDataPartitionRule_DuplicateName(name string) string {
	return fmt.Sprintf(`
resource "newrelic_data_partition_rule" "foo"{
	account_id = %[1]d
	description = "%[3]s"
	enabled = true
	nrql = "logtype='node'"
    retention_policy = "SECONDARY"
    target_data_partition = "Log_Test_%[2]s"
}
resource "newrelic_data_partition_rule" "bar"{
	account_id = %[1]d
	description = "%[3]s"
	enabled = true
	nrql = "logtype='node'"
    retention_policy = "SECONDARY"
    target_data_partition = "Log_Test_%[2]s"
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicDataPartitionRule_ValidateName(name string) string {
	return fmt.Sprintf(`
resource "newrelic_data_partition_rule" "foo"{
	account_id = %[1]d
	description = "%[3]s"
	enabled = true
	nrql = "logtype='node'"
    retention_policy = "SECONDARY"
    target_data_partition = "Test_%[2]s"
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicDataPartitionRuleUpdate_Enabled(name string) string {
	return fmt.Sprintf(`
resource "newrelic_data_partition_rule" "foo"{
	account_id = %[1]d
	description = "%[3]s_update"
	enabled = false
	nrql = "logtype='node'"
    retention_policy = "SECONDARY"
    target_data_partition = "Log_Test_%[2]s"
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicDataPartitionRuleUpdate_Name(name string) string {
	return fmt.Sprintf(`
resource "newrelic_data_partition_rule" "foo"{
	account_id = %[1]d
	description = "%[3]s_update"
	enabled = true
	nrql = "logtype='node'"
    retention_policy = "SECONDARY"
    target_data_partition = "Log_Test_%[2]s_update"
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicDataPartitionRuleUpdate_NRQL(name string) string {
	return fmt.Sprintf(`
resource "newrelic_data_partition_rule" "foo"{
	account_id = %[1]d
	description = "%[3]s_update"
	enabled = true
	nrql = "logtype='node'"
    retention_policy = "SECONDARY"
    target_data_partition = "Log_Test_%[2]s_update"
}
`, testAccountID, name, testAccExpectedApplicationName)
}
