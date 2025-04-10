//go:build integration || WORKFLOW_INTEGRATIONS
// +build integration WORKFLOW_INTEGRATIONS

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
)

func TestNewRelicWorkflow_MicrosoftTeams(t *testing.T) {
	resourceName := "newrelic_workflow.foo"
	rName := generateNameForIntegrationTestResource()

	if testAccMSTeamsDestinationSecurityCode == "" {
		t.Skipf("Skipping this test, as NEW_RELIC_MS_TEAMS_DESTINATION_SECURITY_CODE must be set for this test to run.")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{
			// Test: Create workflow
			{
				Config: testAccNewRelicWorkflowConfigurationMicrosoftTeams(testAccountID, testAccMSTeamsDestinationSecurityCode, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestNewRelicWorkflow_WithInvestigatingNotificationTriggers(t *testing.T) {
	resourceName := "newrelic_workflow.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{
			// Test: Create workflow
			{
				Config: testAccNewRelicWorkflowConfigurationWithNotificationTriggers(testAccountID, rName, `["ACTIVATED", "INVESTIGATING"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// Test: Update workflow
			{
				Config: testAccNewRelicWorkflowConfigurationWithNotificationTriggers(testAccountID, fmt.Sprintf("%s-updated", rName), `["ACTIVATED", "INVESTIGATING", "CLOSED"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
		},
	})
}

func TestNewRelicWorkflow_Webhook(t *testing.T) {
	resourceName := "newrelic_workflow.foo"
	rName := generateNameForIntegrationTestResource()

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
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicWorkflowConfigurationWebhook(testAccountID, fmt.Sprintf("%s-updated", rName)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
		},
	})
}

func TestNewRelicWorkflow_Email(t *testing.T) {
	resourceName := "newrelic_workflow.foo"
	rName := generateNameForIntegrationTestResource()

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
		},
	})
}

func TestNewRelicWorkflow_DeleteWorkflowWithoutRemovingChannels(t *testing.T) {
	resourceName := "newrelic_workflow.foo"
	channelResourceName := "foo"
	channelResources := testAccNewRelicChannelConfigurationEmail(channelResourceName)
	workflowResource := testAccNewRelicWorkflowConfiguration(channelResourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{
			// Test: Create workflow
			{
				Config: channelResources + workflowResource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// Test: Delete workflow
			{
				Config: channelResources,
			},
		},
	})
}

func TestNewRelicWorkflow_MinimalConfig(t *testing.T) {
	resourceName := "newrelic_workflow.foo"
	rName := generateNameForIntegrationTestResource()

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

func TestNewRelicWorkflow_InvalidIssuesFilterAttr(t *testing.T) {
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicWorkflowConfigurationInvalidIssuesFilterAttr(testAccountID, rName),
				ExpectError: regexp.MustCompile("VALIDATION_ERROR"),
			},
		},
	})
}

func TestNewRelicWorkflow_UpdateChannels(t *testing.T) {
	resourceName := "newrelic_workflow.foo"
	channelResourceName := "oldChannel"
	workflowName := acctest.RandString(10)
	channelResources := testAccNewRelicChannelConfigurationEmail(channelResourceName)

	// On order for the test to work, it is really important to make sure that TF does not decide to delete this
	// workflow instead of updating it. Using a minimal workflow allows us to minimise the chance that we touch
	// something that triggers re-creation (accountId is the most important one to avoid)
	workflowResource := testAccNewRelicOnlyWorkflowConfigurationMinimal(workflowName, channelResourceName)

	channelResourceNameNew := "newChannel"
	channelResourcesNew := testAccNewRelicChannelConfigurationEmail(channelResourceNameNew)
	workflowResourceUpdated := testAccNewRelicOnlyWorkflowConfigurationMinimal(workflowName, channelResourceNameNew)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{

			// Test: Create workflow
			{
				Config: channelResources + workflowResource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// Test: Replace channel in the workflow
			{
				Config: channelResources + channelResourcesNew + workflowResourceUpdated,
			},
		},
	})
}

func TestNewRelicWorkflow_RemoveEnrichments(t *testing.T) {
	resourceName := "newrelic_workflow.foo"
	channelResourceName := "foo"
	issuesFilterName := "bar"
	workflowName := generateNameForIntegrationTestResource()
	enrichmentsSection := `
enrichments {
	nrql {
		name = "Log Count"
		configuration {
			query = "SELECT count(*) FROM Log"
		}
	}
}
`
	channelResources := testAccNewRelicChannelConfigurationEmail(channelResourceName)
	workflowWithEnrichment := testAccNewRelicWorkflowConfigurationCustom(workflowName, issuesFilterName, channelResourceName, enrichmentsSection)
	workflowWithoutEnrichment := testAccNewRelicWorkflowConfigurationCustom(workflowName, issuesFilterName, channelResourceName, "")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{
			// Test: Create workflow
			{
				Config: channelResources + workflowWithEnrichment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// Test: Remove enrichments, verify that the plan is empty afterwards
			{
				Config: channelResources + workflowWithoutEnrichment,
			},
		},
	})
}

// This test doesnt seem valid anymore because it looks like the API automatically generates a
// filter name is one is not provided. Skipping this test for now, but we probably can remove
// this test at some point.
func TestNewRelicWorkflow_EmptyIssuesFilterName(t *testing.T) {
	t.Skip("Skipping do due to new API behavior which automatically generates a filter name if one is not provided.")

	channelResourceName := "foo"
	workflowName := acctest.RandString(5)
	issuesFilterName := ""
	channelResources := testAccNewRelicChannelConfigurationEmail(channelResourceName)
	workflowWithEmptyIssuesFilterName := testAccNewRelicWorkflowConfigurationCustom(workflowName, issuesFilterName, channelResourceName, "")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{

			// Test: Create workflow with empty issuesFilter name
			{
				Config:      channelResources + workflowWithEmptyIssuesFilterName,
				ExpectError: regexp.MustCompile(`expected \"issues_filter.0.name\" to not be an empty string or whitespace`),
			},
		},
	})
}

func TestNewRelicWorkflow_WithCreatedNotificationTriggers(t *testing.T) {
	resourceName := "newrelic_workflow.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{
			// Test: Create workflow
			{
				Config: testAccNewRelicWorkflowConfigurationWithNotificationTriggers(testAccountID, rName, `["ACTIVATED"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// Test: Update workflow
			{
				Config: testAccNewRelicWorkflowConfigurationWithNotificationTriggers(testAccountID, fmt.Sprintf("%s-updated", rName), `["ACTIVATED", "CLOSED"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
		},
	})
}

func TestNewRelicWorkflow_WithNonCreatedNotificationTriggers(t *testing.T) {
	resourceName := "newrelic_workflow.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{
			// Test: Create workflow
			{
				Config: testAccNewRelicWorkflowConfigurationWithNotificationTriggers(testAccountID, rName, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicWorkflowConfigurationWithNotificationTriggers(testAccountID, fmt.Sprintf("%s-updated", rName), `["ACTIVATED"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
		},
	})
}

func TestNewRelicWorkflow_WithUpdatedNotificationTriggers(t *testing.T) {
	resourceName := "newrelic_workflow.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{
			// Test: Create workflow
			{
				Config: testAccNewRelicWorkflowConfigurationWithNotificationTriggers(testAccountID, rName, `["ACTIVATED"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicWorkflowConfigurationWithNotificationTriggers(testAccountID, fmt.Sprintf("%s-updated", rName), `["ACTIVATED", "CLOSED"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
		},
	})
}

func TestNewRelicWorkflow_WithCreatedUpdateOriginalMessage(t *testing.T) {
	resourceName := "newrelic_workflow.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{
			// Test: Create workflow
			{
				Config: testAccNewRelicWorkflowConfigurationWithCustomDestination(testAccountID, rName, `update_original_message = true`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// Test: Update workflow
			{
				Config: testAccNewRelicWorkflowConfigurationWithCustomDestination(testAccountID, fmt.Sprintf("%s-updated", rName), `update_original_message = true`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
		},
	})
}

func TestNewRelicWorkflow_WithNonCreatedUpdateOriginalMessage(t *testing.T) {
	resourceName := "newrelic_workflow.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{
			// Test: Create workflow
			{
				Config: testAccNewRelicWorkflowConfigurationWithCustomDestination(testAccountID, rName, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// Test: Update workflow
			{
				Config: testAccNewRelicWorkflowConfigurationWithCustomDestination(testAccountID, fmt.Sprintf("%s-updated", rName), `update_original_message = true`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
		},
	})
}

func TestNewRelicWorkflow_BooleanFlags_DisableOnUpdate(t *testing.T) {
	channelResourceName := "foo"
	workflowName := acctest.RandString(10)
	channelResources := testAccNewRelicChannelConfigurationEmail(channelResourceName)
	workflowResource := testAccNewRelicOnlyWorkflowConfigurationMinimal(workflowName, channelResourceName)
	disabledWorkflowResource := testAccNewRelicWorkflowMinimalDisabledConfiguration(workflowName, channelResourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{

			// Test: Create workflow
			{
				Config: channelResources + workflowResource,
			},
			// Test: Disable workflow
			{
				Config: channelResources + disabledWorkflowResource,
			},
		},
	})
}

func TestNewRelicWorkflow_BooleanFlags_DisableOnCreation(t *testing.T) {
	channelResourceName := "foo"
	workflowName := acctest.RandString(10)
	channelResources := testAccNewRelicChannelConfigurationEmail(channelResourceName)
	disabledWorkflowResource := testAccNewRelicWorkflowMinimalDisabledConfiguration(workflowName, channelResourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{
			// Test: Create workflow
			{
				Config: channelResources + disabledWorkflowResource,
			},
		},
	})
}

func TestNewRelicWorkflow_NotificationTriggerShouldIgnoreOrder(t *testing.T) {
	resourceName := "newrelic_workflow.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicWorkflowDestroy,
		Steps: []resource.TestStep{
			// Test: Create workflow with non-standard trigger order
			{
				Config: testAccNewRelicWorkflowConfigurationWithNotificationTriggers(testAccountID, rName, `["CLOSED", "ACTIVATED"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkflowExists(resourceName),
				),
			},
			// Test: Update trigger order
			{
				Config: testAccNewRelicWorkflowConfigurationWithNotificationTriggers(testAccountID, rName, `["ACKNOWLEDGED", "ACTIVATED", "CLOSED"]`),
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
      attribute = "accumulations.sources"
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

func testAccNewRelicWorkflowConfigurationInvalidIssuesFilterAttr(accountID int, name string) string {
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
      attribute = "invalid.attribute"
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
      attribute = "priority"
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

func testAccNewRelicWorkflowConfigurationMicrosoftTeams(accountID int, msTeamsSecurityCode string, name string) string {
	return fmt.Sprintf(`
resource "newrelic_notification_destination" "foo" {
  account_id = %[1]d
  name = "tf-test-destination"
  type = "MICROSOFT_TEAMS"

  property {
    key   = "securityCode"
    value = "%[2]s"
  }
}

resource "newrelic_notification_channel" "foo" {
  account_id     = newrelic_notification_destination.foo.account_id
  name           = "ms-teams-example"
  type           = "MICROSOFT_TEAMS"
  product        = "IINT"
  destination_id = newrelic_notification_destination.foo.id

  property {
    key = "teamId"
    value = "045a1764-e6a2-47a1-becb-e2ee56605901"
  }
  property {
    key = "channelId"
    value = "19:56d407dbb6bf403ebc88122ad39826cc@thread.tacv2"
  }
}

resource "newrelic_workflow" "foo" {
  account_id            = newrelic_notification_destination.foo.account_id
  name                  = "%[3]s"
  enrichments_enabled   = true
  destinations_enabled  = true
  enabled      = true
  muting_rules_handling = "NOTIFY_ALL_ISSUES"

  issues_filter {
    name = "filter-name"
    type = "FILTER"

    predicate {
      attribute = "priority"
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
`, accountID, msTeamsSecurityCode, name)
}

func testAccNewRelicWorkflowConfigurationWithNotificationTriggers(accountID int, name string, notificationTriggers string) string {
	parsedNotificationTriggers := ""
	if notificationTriggers != "" {
		parsedNotificationTriggers = fmt.Sprintf(`notification_triggers = %[1]s`, notificationTriggers)
	}
	return testAccNewRelicWorkflowConfigurationWithCustomDestination(accountID, name, parsedNotificationTriggers)
}

func testAccNewRelicWorkflowConfigurationWithCustomDestination(accountID int, name string, customDestination string) string {
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
      attribute = "priority"
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
    %[3]s
  }
}
`, accountID, name, customDestination)
}

func testAccNewRelicChannelConfigurationEmail(channelResourceName string) string {
	destinationResourceName := "destination_" + acctest.RandString(5)
	return fmt.Sprintf(`
resource "newrelic_notification_destination" "%[2]s" {
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

resource "newrelic_notification_channel" "%[1]s" {
  name           = "tf-test-notification-channel-email"
  type           = "EMAIL"
  product        = "IINT"
  destination_id = newrelic_notification_destination.%[2]s.id

  property {
    key   = "subject"
    value = "{{ issueTitle }}"
  }
}`, channelResourceName, destinationResourceName)
}

func testAccNewRelicWorkflowConfiguration(channelResourceName string) string {
	return testAccNewRelicWorkflowConfigurationCustom(
		acctest.RandString(5),
		acctest.RandString(5),
		channelResourceName,
		"")
}

func testAccNewRelicWorkflowConfigurationCustom(workflowName string, issuesFilterName string, channelResourceName string, customSections string) string {
	return fmt.Sprintf(`
resource "newrelic_workflow" "foo" {
  name                  = "%[3]s"
  enrichments_enabled   = true
  destinations_enabled  = true
  enabled      = true
  muting_rules_handling = "NOTIFY_ALL_ISSUES"

  issues_filter {
    name = "%[2]s"
    type = "FILTER"

    predicate {
      attribute = "priority"
      operator  = "EQUAL"
      values    = ["test"]
    }
  }

  %[4]s

  destination {
    channel_id = newrelic_notification_channel.%[1]s.id
  }
}`, channelResourceName, issuesFilterName, workflowName, customSections)
}

func testAccNewRelicOnlyWorkflowConfigurationMinimal(workflowName string, channelResourceName string) string {
	return fmt.Sprintf(`
resource "newrelic_workflow" "foo" {
  name                  = "%[2]s"
  muting_rules_handling = "NOTIFY_ALL_ISSUES"

  issues_filter {
    name = "filter-name"
    type = "FILTER"

    predicate {
      attribute = "priority"
      operator  = "EQUAL"
      values    = ["test"]
    }
  }

  destination {
    channel_id = newrelic_notification_channel.%[1]s.id
  }
}`, channelResourceName, workflowName)
}

func testAccNewRelicWorkflowMinimalDisabledConfiguration(workflowName, channelResourceName string) string {
	return fmt.Sprintf(`
resource "newrelic_workflow" "foo" {
	name                  = "%[2]s"
	muting_rules_handling = "NOTIFY_ALL_ISSUES"
	enrichments_enabled   = false
	destinations_enabled  = false
	enabled      = false

	issues_filter {
		name = "filter-name"
		type = "FILTER"

		predicate {
			attribute = "priority"
			operator  = "EQUAL"
			values    = ["test"]
		}
	}

	destination {
		channel_id = newrelic_notification_channel.%[1]s.id
	}
}`, channelResourceName, workflowName)

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
