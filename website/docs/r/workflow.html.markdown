---
layout: "newrelic"
page_title: "New Relic: newrelic_workflow"
sidebar_current: "docs-newrelic-resource-workflow"
description: |-
Create and manage a workflow in New Relic.
---

# Resource: newrelic\_workflow

Use this resource to create and manage New Relic workflow.

## Example Usage

##### Workflow
```hcl
# Workflows
resource "newrelic_workflow" "foo" {
  name = "workflow-example"
  enrichments_enabled = true
  destinations_enabled = true
  workflow_enabled = true
  muting_rules_handling = "NOTIFY_ALL_ISSUES"

  enrichments = {
    nrql = {
      name = "Log"
      configurations = {
        query = "SELECT * FROM Log"
      }
    }

    nrql = {
      name = "Metric"
      configurations = {
        query = "SELECT * FROM Metric"
      }
    }
  }

  issues_filter = {
    name = "Filter1"
    type = "FILTER"
    predicates = {
      attribute = "source"
      operator = "EQUAL"
      values = "newrelic"
    }
  }

  destination_configurations {
    channel_id = "d8ad79ce-c8e9-4451-8f7e-1f04a613997e"
  }
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) Determines the New Relic account where the workflow will be created. Defaults to the account associated with the API key used.
* `name` - (Required) The name of the workflow.
* `enrichments_enabled` - (Optional) Whether enrichments are enabled..
* `destinations_enabled` - (Optional) Whether destinations are enabled..
* `workflow_enabled` - (Optional) Whether workflow is enabled.
* `muting_rules_handling` - (Optional) Which muting rule handling this workflow has.
* `destination_configurations` - A nested block that describes a notification channel properties. See [Nested properties blocks](#nested-properties-blocks) below for details.
* `issues_filter` - (Required) The type of product.  One of: `DISCUSSIONS`, `ERROR_TRACKING` or `IINT` (workflows).
* `enrichments` - (Optional) A nested block that describes a notification channel properties. See [Nested properties blocks](#nested-properties-blocks) below for details.

### Nested `destination_configurations` blocks

Each workflow type supports a specific set of arguments for the `destination_configurations` block:

* `WEBHOOK`
  * `headers` - (Optional) A map of key/value pairs that represents the webhook headers.
  * `payload` - (Required) A map of key/value pairs that represents the webhook payload.
* `SERVICENOW_INCIDENTS`
  * `description` - (Optional) A map of key/value pairs that represents a description.
  * `short_description` - (Optional) A map of key/value pairs that represents a short description.
* `JIRA_CLASSIC`, `JIRA_NEXTGEN`
  * `project` - (Required) A map of key/value pairs that represents the jira project id.
  * `issuetype` - (Required) A map of key/value pairs that represents the issue type id.
* `EMAIL`
  * `subject` - (Optional) A map of key/value pairs that represents the email subject title.
  * `customDetailsEmail` - (Optional) A map of key/value pairs that represents the email custom details.
* `PAGERDUTY_SERVICE_INTEGRATION`
  * `summary` - (Required) A map of key/value pairs that represents the summery.
* `PAGERDUTY_ACCOUNT_INTEGRATION`
  * `summary` - (Required) A map of key/value pairs that represents the summery.
  * `service` - (Required) Specifies the service id for integrating with Pagerduty.
  * `email` - (Required) Specifies the user email for integrating with Pagerduty.

### Nested `issues_filter` blocks

Each workflow type supports a specific set of arguments for the `issues_filter` block:

* `WEBHOOK`
  * `headers` - (Optional) A map of key/value pairs that represents the webhook headers.
  * `payload` - (Required) A map of key/value pairs that represents the webhook payload.
* `SERVICENOW_INCIDENTS`
  * `description` - (Optional) A map of key/value pairs that represents a description.
  * `short_description` - (Optional) A map of key/value pairs that represents a short description.
* `JIRA_CLASSIC`, `JIRA_NEXTGEN`
  * `project` - (Required) A map of key/value pairs that represents the jira project id.
  * `issuetype` - (Required) A map of key/value pairs that represents the issue type id.
* `EMAIL`
  * `subject` - (Optional) A map of key/value pairs that represents the email subject title.
  * `customDetailsEmail` - (Optional) A map of key/value pairs that represents the email custom details.
* `PAGERDUTY_SERVICE_INTEGRATION`
  * `summary` - (Required) A map of key/value pairs that represents the summery.
* `PAGERDUTY_ACCOUNT_INTEGRATION`
  * `summary` - (Required) A map of key/value pairs that represents the summery.
  * `service` - (Required) Specifies the service id for integrating with Pagerduty.
  * `email` - (Required) Specifies the user email for integrating with Pagerduty.

### Nested `enrichments` blocks

Each workflow type supports a specific set of arguments for the `enrichments` block:

* `WEBHOOK`
  * `headers` - (Optional) A map of key/value pairs that represents the webhook headers.
  * `payload` - (Required) A map of key/value pairs that represents the webhook payload.
* `SERVICENOW_INCIDENTS`
  * `description` - (Optional) A map of key/value pairs that represents a description.
  * `short_description` - (Optional) A map of key/value pairs that represents a short description.
* `JIRA_CLASSIC`, `JIRA_NEXTGEN`
  * `project` - (Required) A map of key/value pairs that represents the jira project id.
  * `issuetype` - (Required) A map of key/value pairs that represents the issue type id.
* `EMAIL`
  * `subject` - (Optional) A map of key/value pairs that represents the email subject title.
  * `customDetailsEmail` - (Optional) A map of key/value pairs that represents the email custom details.
* `PAGERDUTY_SERVICE_INTEGRATION`
  * `summary` - (Required) A map of key/value pairs that represents the summery.
* `PAGERDUTY_ACCOUNT_INTEGRATION`
  * `summary` - (Required) A map of key/value pairs that represents the summery.
  * `service` - (Required) Specifies the service id for integrating with Pagerduty.
  * `email` - (Required) Specifies the user email for integrating with Pagerduty.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the workflow.

## Additional Information
More details about the workflows can be found [here](https://docs.newrelic.com/docs/alerts-applied-intelligence/applied-intelligence/incident-workflows/incident-workflows/).