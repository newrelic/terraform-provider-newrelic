//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
)

func TestNewRelicWorkflow_Webhook(t *testing.T) {
	resourceName := "newrelic_workflow.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-workflow-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{

			// Test: Create workflow
			{
				Config: testAccNewRelicWorkflowConfigurationWebhook(testAccountID, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicWorkflowConfigurationWebhook(testAccountID, fmt.Sprintf("%s-updated", rName)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// Test: Import
			// {
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ResourceName:      resourceName,
			// },
		},
	})
}

func TestNewRelicWorkflow_Email(t *testing.T) {
	resourceName := "newrelic_workflow.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-workflow-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{

			// Test: Create workflow
			{
				Config: testAccNewRelicWorkflowConfigurationEmail(testAccountID, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicWorkflowConfigurationEmail(testAccountID, fmt.Sprintf("%s-updated", rName)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// Test: Import
			// {
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ResourceName:      resourceName,
			// },
		},
	})
}

func TestNewRelicWorkflow_MinimalConfig(t *testing.T) {
	resourceName := "newrelic_workflow.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-workflow-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{

			// Test: Create workflow
			{
				Config: testAccNewRelicWorkflowConfigurationMinimal(testAccountID, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicWorkflowConfigurationMinimal(testAccountID, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
		},
	})
}

func testAccNewRelicWorkflowConfigurationMinimal(accountID int, name string) string {
	return fmt.Sprintf(`
resource "newrelic_notification_destination" "foo" {
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
  name           = "webhook-example"
  type           = "WEBHOOK"
  product        = "IINT"
  destination_id = newrelic_notification_destination.foo.id

  property {
    key   = "payload"
    value = "{}"
  }
}

resource "newrelic_workflow" "foo" {
  name                  = "%[2]s"
  muting_rules_handling = "NOTIFY_ALL_ISSUES"

  issues_filter {
    name = "filter-name"
    type = "FILTER"

    predicate {
      attribute = "source"
      operator  = "EQUAL"
      values    = ["newrelic"]
    }
  }

  destination {
    channel_id = newrelic_notification_channel.foo.id
  }
}
`, accountID, name)
}

func testAccNewRelicWorkflowConfigurationWebhook(accountID int, name string) string {
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
  enabled      = true
  muting_rules_handling = "NOTIFY_ALL_ISSUES"

  issues_filter {
    name = "filter-name"
    type = "FILTER"

    predicate {
      attribute = "source"
      operator  = "EQUAL"
      values    = ["newrelic", "pagerduty"]
    }
  }

  enrichments {
    nrql {
      name = "Log"
      configuration {
        query = "SELECT count(*) FROM Log"
      }
    }
  }

  destination {
    channel_id = newrelic_notification_channel.foo.id
  }
}
`, accountID, name)
}

func testAccNewRelicWorkflowConfigurationEmail(accountID int, name string) string {
	return fmt.Sprintf(`
resource "newrelic_notification_destination" "foo" {
  account_id = %[1]d
  name = "tf-test-destination-email"
  type = "EMAIL"

  property {
    key   = "email"
    value = "noreply+terraform-test@newrelic.com"
  }

  auth_basic {
    user     = "username"
    password = "password"
  }
}

resource "newrelic_notification_channel" "foo" {
  account_id     = newrelic_notification_destination.foo.account_id
  name           = "tf-test-notification-channel-email"
  type           = "EMAIL"
  product        = "IINT"
  destination_id = newrelic_notification_destination.foo.id

  property {
    key   = "subject"
    value = "{{ issueTitle }}"
  }

	property {
    key   = "customDetailsEmail"
    value = "This text is a part of the email body"
  }
}

resource "newrelic_workflow" "foo" {
  account_id            = newrelic_notification_destination.foo.account_id
  name                  = "%[2]s"
  enrichments_enabled   = true
  destinations_enabled  = true
  enabled      = true
  muting_rules_handling = "NOTIFY_ALL_ISSUES"

  issues_filter {
    name = "example-filter-by-team-name"
    type = "FILTER"

    predicate {
      attribute = "accumulations.tag.Team"
      operator  = "EXACTLY_MATCHES"
      values    = ["developer-toolkit"]
    }
  }

  enrichments {
    nrql {
      name = "Log"
      configuration {
        query = "SELECT count(*) FROM Log"
      }
    }
  }

  destination {
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
