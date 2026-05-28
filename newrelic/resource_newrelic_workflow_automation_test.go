//go:build integration || WORKFLOW_AUTOMATION

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Test basic CRUD operations for workflow automation with ACCOUNT scope
func TestAccNewRelicWorkflowAutomation_Basic_Account(t *testing.T) {
	resourceName := "newrelic_workflow_automation.foo"
	rName := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkflowAutomationDestroy,
		Steps: []resource.TestStep{
			// Test: Create workflow automation
			{
				Config: testAccNewRelicWorkflowAutomationConfig_Account(testAccountID, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowAutomationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "scope_type", "ACCOUNT"),
					resource.TestCheckResourceAttr(resourceName, "scope_id", fmt.Sprintf("%d", testAccountID)),
					resource.TestCheckResourceAttrSet(resourceName, "definition"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
		},
	})
}

// Test update operations for workflow automation
func TestAccNewRelicWorkflowAutomation_Update_Account(t *testing.T) {
	resourceName := "newrelic_workflow_automation.foo"
	rName := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkflowAutomationDestroy,
		Steps: []resource.TestStep{
			// Test: Create workflow automation
			{
				Config: testAccNewRelicWorkflowAutomationConfig_Account(testAccountID, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowAutomationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			// Test: Update workflow automation definition
			{
				Config: testAccNewRelicWorkflowAutomationConfig_AccountUpdated(testAccountID, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowAutomationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

// Test import functionality
func TestAccNewRelicWorkflowAutomation_Import_Account(t *testing.T) {
	resourceName := "newrelic_workflow_automation.foo"
	rName := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkflowAutomationDestroy,
		Steps: []resource.TestStep{
			// Test: Create workflow automation
			{
				Config: testAccNewRelicWorkflowAutomationConfig_Account(testAccountID, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowAutomationExists(resourceName),
				),
			},
			// Test: Import workflow automation
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test workflow automation with name specified in Terraform config matching YAML
func TestAccNewRelicWorkflowAutomation_WithNameInConfig_Account(t *testing.T) {
	resourceName := "newrelic_workflow_automation.foo"
	rName := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkflowAutomationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicWorkflowAutomationConfig_AccountWithName(testAccountID, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowAutomationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

// Test validation: name mismatch between Terraform config and YAML
func TestAccNewRelicWorkflowAutomation_NameMismatch(t *testing.T) {
	rName := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkflowAutomationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicWorkflowAutomationConfig_NameMismatch(testAccountID, rName),
				ExpectError: regexp.MustCompile("name in resource configuration .* does not match name in YAML definition"),
			},
		},
	})
}

// Test validation: invalid scope_type
func TestAccNewRelicWorkflowAutomation_InvalidScopeType(t *testing.T) {
	rName := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkflowAutomationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicWorkflowAutomationConfig_InvalidScopeType(testAccountID, rName),
				ExpectError: regexp.MustCompile("scope_type .* is not supported"),
			},
		},
	})
}

// Test validation: empty scope_type
func TestAccNewRelicWorkflowAutomation_EmptyScopeType(t *testing.T) {
	rName := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkflowAutomationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicWorkflowAutomationConfig_EmptyScopeType(testAccountID, rName),
				ExpectError: regexp.MustCompile("scope_type is required"),
			},
		},
	})
}

// Test validation: empty scope_id
func TestAccNewRelicWorkflowAutomation_EmptyScopeId(t *testing.T) {
	rName := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkflowAutomationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicWorkflowAutomationConfig_EmptyScopeId(rName),
				ExpectError: regexp.MustCompile("scope_id is required"),
			},
		},
	})
}

// Test validation: invalid YAML
func TestAccNewRelicWorkflowAutomation_InvalidYAML(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkflowAutomationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicWorkflowAutomationConfig_InvalidYAML(testAccountID),
				ExpectError: regexp.MustCompile("failed to parse YAML definition"),
			},
		},
	})
}

// Test validation: YAML without name field
func TestAccNewRelicWorkflowAutomation_YAMLWithoutName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkflowAutomationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicWorkflowAutomationConfig_YAMLWithoutName(testAccountID),
				ExpectError: regexp.MustCompile("name field not found in YAML definition"),
			},
		},
	})
}

// Helper function to check if workflow automation exists
func testAccCheckNewRelicWorkflowAutomationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no workflow automation ID is set")
		}

		providerConfig := testAccProvider.Meta().(*ProviderConfig)
		client := providerConfig.NewClient

		// Parse the resource ID
		scopeType := rs.Primary.Attributes["scope_type"]
		scopeID := rs.Primary.Attributes["scope_id"]
		name := rs.Primary.Attributes["name"]

		if scopeType == "ACCOUNT" {
			// For ACCOUNT scope, convert scopeID to int
			var accountID int
			_, err := fmt.Sscanf(scopeID, "%d", &accountID)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", scopeID)
			}

			workflow, err := client.WorkflowAutomation.GetWorkflow(accountID, name, 0)
			if err != nil {
				return fmt.Errorf("error fetching workflow automation: %s", err)
			}

			if workflow == nil {
				return fmt.Errorf("workflow automation not found")
			}
		} else if scopeType == "ORGANIZATION" {
			workflow, err := client.Organization.GetWorkflow(name, 0)
			if err != nil {
				return fmt.Errorf("error fetching workflow automation: %s", err)
			}

			if workflow == nil {
				return fmt.Errorf("workflow automation not found")
			}
		}

		return nil
	}
}

// Helper function to check if workflow automation is destroyed
func testAccCheckNewRelicWorkflowAutomationDestroy(s *terraform.State) error {
	providerConfig := testAccProvider.Meta().(*ProviderConfig)
	client := providerConfig.NewClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "newrelic_workflow_automation" {
			continue
		}

		scopeType := rs.Primary.Attributes["scope_type"]
		scopeID := rs.Primary.Attributes["scope_id"]
		name := rs.Primary.Attributes["name"]

		if scopeType == "ACCOUNT" {
			var accountID int
			_, err := fmt.Sscanf(scopeID, "%d", &accountID)
			if err != nil {
				// If we can't parse the ID, consider it destroyed
				continue
			}

			workflow, err := client.WorkflowAutomation.GetWorkflow(accountID, name, 0)
			if err == nil && workflow != nil {
				return fmt.Errorf("workflow automation still exists")
			}
		} else if scopeType == "ORGANIZATION" {
			workflow, err := client.Organization.GetWorkflow(name, 0)
			if err == nil && workflow != nil {
				return fmt.Errorf("workflow automation still exists")
			}
		}
	}

	return nil
}

// Test configuration helper functions

func testAccNewRelicWorkflowAutomationConfig_Account(accountID int, name string) string {
	return fmt.Sprintf(`
resource "newrelic_workflow_automation" "foo" {
  name       = "%[2]s"
  scope_id   = "%[1]d"
  scope_type = "ACCOUNT"

  definition = <<-EOT
name: %[2]s
description: This is a test workflow created by terraform
steps:
  - name: waitStep1
    type: wait
    seconds: 10
  - name: waitStep2
    type: wait
    seconds: 10
EOT
}
`, accountID, name)
}

func testAccNewRelicWorkflowAutomationConfig_AccountUpdated(accountID int, name string) string {
	return fmt.Sprintf(`
resource "newrelic_workflow_automation" "foo" {
  name       = "%[2]s"
  scope_id   = "%[1]d"
  scope_type = "ACCOUNT"

  definition = <<-EOT
name: %[2]s
description: This is an updated test workflow created by terraform
steps:
  - name: waitStep1
    type: wait
    seconds: 15
  - name: waitStep2
    type: wait
    seconds: 15
  - name: waitStep3
    type: wait
    seconds: 10
EOT
}
`, accountID, name)
}

func testAccNewRelicWorkflowAutomationConfig_AccountWithName(accountID int, name string) string {
	return fmt.Sprintf(`
resource "newrelic_workflow_automation" "foo" {
  name       = "%[2]s"
  scope_id   = "%[1]d"
  scope_type = "ACCOUNT"

  definition = <<-EOT
name: %[2]s
description: This is a test workflow created by terraform
steps:
  - name: waitStep1
    type: wait
    seconds: 10
  - name: waitStep2
    type: wait
    seconds: 10
EOT
}
`, accountID, name)
}

func testAccNewRelicWorkflowAutomationConfig_NameMismatch(accountID int, name string) string {
	return fmt.Sprintf(`
resource "newrelic_workflow_automation" "foo" {
  name       = "wrong_name"
  scope_id   = "%[1]d"
  scope_type = "ACCOUNT"

  definition = <<-EOT
name: %[2]s
description: This is a test workflow created by terraform
steps:
  - name: waitStep1
    type: wait
    seconds: 10
  - name: waitStep2
    type: wait
    seconds: 10
EOT
}
`, accountID, name)
}

func testAccNewRelicWorkflowAutomationConfig_InvalidScopeType(accountID int, name string) string {
	return fmt.Sprintf(`
resource "newrelic_workflow_automation" "foo" {
  name       = "%[2]s"
  scope_id   = "%[1]d"
  scope_type = "INVALID_SCOPE"

  definition = <<-EOT
name: %[2]s
description: This is a test workflow created by terraform
steps:
  - name: waitStep1
    type: wait
    seconds: 10
  - name: waitStep2
    type: wait
    seconds: 10
EOT
}
`, accountID, name)
}

func testAccNewRelicWorkflowAutomationConfig_EmptyScopeType(accountID int, name string) string {
	return fmt.Sprintf(`
resource "newrelic_workflow_automation" "foo" {
  name       = "%[2]s"
  scope_id   = "%[1]d"
  scope_type = ""

  definition = <<-EOT
name: %[2]s
description: This is a test workflow created by terraform
steps:
  - name: waitStep1
    type: wait
    seconds: 10
  - name: waitStep2
    type: wait
    seconds: 10
EOT
}
`, accountID, name)
}

func testAccNewRelicWorkflowAutomationConfig_EmptyScopeId(name string) string {
	return fmt.Sprintf(`
resource "newrelic_workflow_automation" "foo" {
  name       = "%[1]s"
  scope_id   = ""
  scope_type = "ACCOUNT"

  definition = <<-EOT
name: %[1]s
description: This is a test workflow created by terraform
steps:
  - name: waitStep1
    type: wait
    seconds: 10
  - name: waitStep2
    type: wait
    seconds: 10
EOT
}
`, name)
}

func testAccNewRelicWorkflowAutomationConfig_InvalidYAML(accountID int) string {
	return fmt.Sprintf(`
resource "newrelic_workflow_automation" "foo" {
  name       = "invalid_yaml_test"
  scope_id   = "%[1]d"
  scope_type = "ACCOUNT"

  definition = "this is not valid yaml: [[[{"
}
`, accountID)
}

func testAccNewRelicWorkflowAutomationConfig_YAMLWithoutName(accountID int) string {
	return fmt.Sprintf(`
resource "newrelic_workflow_automation" "foo" {
  name       = "no_name_in_yaml"
  scope_id   = "%[1]d"
  scope_type = "ACCOUNT"

  definition = <<-EOT
description: This is a test workflow created by terraform
steps:
  - name: waitStep1
    type: wait
    seconds: 10
  - name: waitStep2
    type: wait
    seconds: 10
EOT
}
`, accountID)
}
