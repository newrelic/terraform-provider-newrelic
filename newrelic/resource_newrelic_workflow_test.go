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
	// t.Skip("Skipping TestNewRelicWorkflow_Basic.  AWAITING FINAL IMPLEMENTATION!")

	resourceName := "newrelic_workflow"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-workflow-test-%s", rand)
	// enrichments := `{
	// 					nrql {
	// 					  name = "Log"
	// 					  configurations {
	// 					   query = "SELECT count(*) FROM Log"
	// 					  }
	// 					}
	// 				}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{

			// Test: Create workflow
			{
				Config: testAccNewRelicWorkflowConfiguration(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// // Test: Update
			// {
			// 	Config: testAccNewRelicWorkflowConfiguration(rName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckNewRelicWorkflowExists(resourceName),
			// 	),
			// },
			// Test: Import
			// {
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ResourceName:      resourceName,
			// },
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

		filters := ai.AiWorkflowsFilters{
			ID: r.Primary.ID,
		}

		_, err := client.Workflows.GetWorkflows(testAccountID, "", filters)
		if err == nil {
			return fmt.Errorf("workflow still exists")
		}

	}
	return nil
}

func testAccNewRelicWorkflowConfiguration(name string) string {
	return fmt.Sprintf(`
resource "newrelic_notification_destination" "foo" {
	name = "tf-test-destination"
	type = "WEBHOOK"

	properties {
		key   = "url"
		value = "https://webhook.site/"
	}

	auth = {
		type     = "BASIC"
		user     = "username"
		password = "password"
	}
}

resource "newrelic_notification_channel" "foo" {
	name           = "webhook-example"
	type           = "WEBHOOK"
	product        = "IINT"
	destination_id = newrelic_notification_destination.foo.id

	properties {
		key   = "payload"
		value = "{\n\t\"name\": \"foo\"\n}"
		label = "Payload Template"
	}
}

resource "newrelic_workflow" "foo" {
	name                  = "%s"
	enrichments_enabled   = false
	destinations_enabled  = true
	workflow_enabled      = true
	muting_rules_handling = "NOTIFY_ALL_ISSUES"

	issues_filter {
		name = "filter-name"
		type = "FILTER"

		predicates {
			attribute = "source"
			operator  = "EQUAL"
			values    = ["newrelic", "pagerduty"]
		}
	}

	destination_configuration {
		channel_id = newrelic_notification_channel.foo.id
	}
}
`, name)
}

func testNewRelicWorkflowConfigByType(name string, enrichments_enabled string, destinations_enabled string, workflow_enabled string, muting_rules_handling string, enrichments string, issuesFilter string, destinationConfigurations string) string {
	if enrichments == "" {
		return fmt.Sprintf(`
		resource "newrelic_workflow" "test-workflow" {
			name = "%s"
			enrichments_enabled = "%s"
			destinations_enabled = "%s"
			workflow_enabled = "%s"
			muting_rules_handling %s
			issues_filter = {
				%s
			}
			destination_configurations {
				%s
			}
		}
	`, name, enrichments_enabled, destinations_enabled, workflow_enabled, muting_rules_handling, issuesFilter, destinationConfigurations)
	}

	return fmt.Sprintf(`
		resource "newrelic_workflow" "test-workflow" {
			name = "%s"
			enrichments_enabled = "%s"
			destinations_enabled = "%s"
			workflow_enabled = "%s"
			muting_rules_handling %s
			enrichments {
				%s
			}
			issues_filter = {
				%s
			}
			destination_configurations {
				%s
			}
		}
	`, name, enrichments_enabled, destinations_enabled, workflow_enabled, muting_rules_handling, enrichments, issuesFilter, destinationConfigurations)
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
