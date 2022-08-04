//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/ai"
)

func TestNewRelicWorkflow_Basic(t *testing.T) {
	t.Skip("Skipping TestNewRelicWorkflow_Basic.  AWAITING FINAL IMPLEMENTATION!")

	resourceName := "newrelic_workflow.test-workflow"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-workflow-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testNewRelicWorkflowConfigByType(rName, "WEBHOOK", "IINT", "b1e90a32-23b7-4028-b2c7-ffbdfe103852", `{
					key = "payload"
					value = "{\n\t\"id\": \"test\"\n}"
					label = "Payload Template"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testNewRelicWorkflowConfigByType(rName, "WEBHOOK", "IINT", "b1e90a32-23b7-4028-b2c7-ffbdfe103852", `{
					key = "payload"
					value = "{\n\t\"id\": \"test-update\"\n}"
					label = "Payload Template Update"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// Test: Import
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
		},
	})
}

func testAccNewRelicWorkflowDestroy(s *terraform.State) error {
	providerConfig := testAccProvider.Meta().(*ProviderConfig)
	client := providerConfig.NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_workflow" {
			continue
		}

		var accountID int
		id := r.Primary.ID
		accountID = providerConfig.AccountID
		filters := ai.AiWorkflowsFilters{
			ID: id,
		}

		_, err := client.Workflows.GetWorkflows(accountID, "", filters)
		if err == nil {
			return fmt.Errorf("workflow still exists")
		}

	}
	return nil
}

func testNewRelicWorkflowConfigByType(name string, channelType string, product string, destinationId string, properties string) string {
	return fmt.Sprintf(`
		resource "newrelic_workflow" "test-workflow" {
			name = "%s"
			enrichments_enabled = "%s"
			destinations_enabled = "%s"
			workflow_enabled = "%s"
			muting_rules_handling %s
			enrichments
			issues_filter
			destination_configurations
		}
	`, name, channelType, product, destinationId, properties)
}

func testAccCheckNewRelicWorkflowExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providerConfig := testAccProvider.Meta().(*ProviderConfig)
		client := providerConfig.NewClient

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no workflow ID is set")
		}

		var accountID int
		id := rs.Primary.ID
		accountID = providerConfig.AccountID
		filters := ai.AiWorkflowsFilters{
			ID: id,
		}

		found, err := client.Workflows.GetWorkflows(accountID, "", filters)
		if err != nil {
			return err
		}

		if string(found.Entities[0].ID) != rs.Primary.ID {
			return fmt.Errorf("workflow not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}
