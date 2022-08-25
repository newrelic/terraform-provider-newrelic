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

	resourceName := "newrelic_workflow.foo"
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
				Config: testAccNewRelicWorkflowConfiguration(testAccountID, rName),
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

func testAccNewRelicWorkflowConfiguration(accountID int, name string) string {
	return fmt.Sprintf(`
resource "newrelic_notification_destination" "foo" {
  account_id = %[1]d
  name = "tf-test-destination"
  type = "WEBHOOK"

  property {
    key   = "url"
    value = "https://webhook.site/"
  }

  auth_basic {
    user     = "username"
    password = "password"
  }
}

resource "newrelic_notification_channel" "foo" {
  account_id     = newrelic_notification_destination.foo.account_id
  name           = "webhook-example"
  type           = "WEBHOOK"
  product        = "IINT"
  destination_id = newrelic_notification_destination.foo.id

  property {
    key   = "payload"
    value = "{\n\t\"name\": \"foo\"\n}"
    label = "Payload Template"
  }
}

resource "newrelic_workflow" "foo" {
  account_id            = newrelic_notification_destination.foo.account_id
  name                  = "%[2]s"
  enrichments_enabled   = true
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

  enrichments {
    nrql {
      name = "Log"
      configurations {
        query = "SELECT count(*) FROM Log"
      }
    }
  }

  destination_configuration {
    channel_id = newrelic_notification_channel.foo.id
  }
}
`, accountID, name)
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

		resp, err := client.Workflows.GetWorkflows(testAccountID, "", filters)
		if len(resp.Entities) > 0 {
			return fmt.Errorf("workflow still exists")
		}

		if err != nil {
			return err
		}

	}
	return nil
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
