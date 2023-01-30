---
layout: "newrelic"
page_title: "New Relic: newrelic_workflow"
sidebar_current: "docs-newrelic-resource-workflow"
description: |-
Create and manage a workflow in New Relic.
---

# Resource: newrelic\_workflow

Use this resource to create and manage New Relic workflows.

## Example Usage

##### Workflow
```hcl
resource "newrelic_workflow" "foo" {
  name = "workflow-example"
  muting_rules_handling = "NOTIFY_ALL_ISSUES"

  issues_filter {
    name = "filter-name"
    type = "FILTER"

    predicate {
      attribute = "accumulations.tag.team"
      operator = "EXACTLY_MATCHES"
      values = [ "growth" ]
    }
  }

  destination {
    channel_id = newrelic_notification_channel.some_channel.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the workflow.
* `issues_filter` - (Required) A filter used to identify issues handled by this workflow. See [Nested issues_filter blocks](#nested-issues_filter-blocks) below for details.
* `muting_rules_handling` - (Required) How to handle muted issues. See [Muting Rules](#muting-rules) below for details.
* `destination` - (Required) Notification configuration. See [Nested destination blocks](#nested-destination-blocks) below for details.
* `account_id` - (Optional) Determines the New Relic account in which the workflow is created. Defaults to the account defined in the provider section.
* `enrichments_enabled` - (Optional) Whether enrichments are enabled. Defaults to true.
* `destinations_enabled` - (Optional) **DEPRECATED** Whether destinations are enabled. Please use `enabled` instead:
these two are different flags, but they are functionally identical. Defaults to true.
* `enabled` - (Optional) Whether workflow is enabled. Defaults to true.
* `enrichments` - (Optional) Workflow's enrichments. See [Nested enrichments blocks](#nested-enrichments-blocks) below for details.

### Nested `issues_filter` blocks

Issue filter defines how to identify issues that should be handled by the workflow.
It consists of one or more predicates.

An issue must match **all** predicates to be handled by the workflow.

For example, we have a workflow with two predicates in the issue filter: 
- `tag.team EXACTLY_MATCHES my_team` meaning "the issue has a tag called `team` and the tag value is set to 'my_team'" 
- `labels.policyIds EXACTLY_MATCHES 123` meaning "the issue includes an incident triggered by an alert policy with id = 123"

In this case, an issue would have to **both** have the tag and include an incident triggered by the said policy in order to be processed by the workflow.  

Each `issues_filter` block supports the following arguments:

* `name` - (Required) The name of the filter. The name only serves a cosmetic purpose and can only be seen through Terraform and GraphQL API. It can't be empty.
* `type` - (Required) Type of the filter. Please just set this field to `FILTER`. The field is likely to be deprecated/removed in the near future. 
* `predicate` (Required) A condition an issue event should satisfy to be processed by the workflow 
  * `attribute` - (Required) Issue event attribute to check
  * `operator` - (Required) An operator to use to compare the attribute with the provided `values`, see supported operators below
  * `values` - (Required) The `attribute` must match **any** of the values in this list  

#### Issue Attribute Types

Each attribute can have one of the following types:
- Plain String
  - String attributes normally represent a text property of an NR issue (i.e. not a property of a child incident)
  - Example attributes: `state`, `priority`
- Numbers
  - Similarly to strings, number attributes represent a numerical property of an NR issue  
  - Most of the number attributes are timestamps, e.g.: `timestamp`, `acknowledgedAt`
- Lists
  - Attributes that belong to issue's child incidents are represented as lists of strings 
  - Each issue might consist of multiple incidents, which is why incident properties are lists of values
  - Examples:
    - `labels.policyIds` - a list of IDs of alert policies that triggered the incidents included in the issue
    - `accumulations.tags.X` - a list of values of a tag `X` collected from all incidents in the issue

While the descriptions above might allow you to guess the field type, you can also check the type by going to
the issue page in the UI and inspecting the issue payload using a button at the top.
- If the attribute value is `null`, you can try checking a different issue
- If the value is a number, then it is a number field
- If the value is a string that contains square brackets, then it is most likely a list field
  - Example: `"[\"some_value\"]"`
- Otherwise, it is a string field 

#### Operators

Depending on their type, different issue attributes support different operators.

**All operators are case-insensitive**

##### Plan String

Plain strings support the following operators:
- `CONTAINS`, `DOES_NOT_CONTAIN` - check if the value contains one of the given strings 
- `EQUAL`, `DOES_NOT_EQUAL` - check if the value is equal to one of the given strings
- `STARTS_WITH` - check if the value starts from one of the given strings
- `ENDS_WITH` - check if the value starts from one of the given strings

##### Numbers

Numbers support the following operators:
- `EQUAL`
- `DOES_NOT_EQUAL`
- `LESS_OR_EQUAL`
- `LESS_THAN`
- `GREATER_OR_EQUAL`
- `GREATER_THAN`

These operators directly correspond to the mathematical operations.

##### Lists

Lists support the following operators:
- `DOES_NOT_EXACTLY_MATCH`, `EXACTLY_MATCHES` - check if the array contains one of the given items
  - `["Aa","Bb"]` would "exactly match" `Aa`
  - `["Aa","Bb"]` would not "exactly match" `A`
- `CONTAINS`, `DOES_NOT_CONTAIN` - **[please avoid]** check if the array has an item with a value that contains any of the given values
  `["Aa","Bb"]` would "contain" **both** `A` **and** `Aa`

#### Other

All fields also support 
- `IS`, `IS_NOT` - **[please avoid]** check if something is `NULL` or not. Should only be used with `NULL` value

### Muting Rules

An issue can be either "Fully muted", "Partially muted" or "Not muted" (for more on these states please visit [our official documentation page about them](https://docs.newrelic.com/docs/alerts-applied-intelligence/new-relic-alerts/alert-notifications/muting-rules-suppress-notifications/#workflow-behavior)).

Muting rule defines how to handle issues that are partially or fully muted. 

Possible values:
* `DONT_NOTIFY_FULLY_MUTED_ISSUES` - Do not send notifications for fully muted issues, do send notifications for partially muted issues
* `DONT_NOTIFY_FULLY_OR_PARTIALLY_MUTED_ISSUES` - Do not send notifications for fully or partially muted issues.
* `NOTIFY_ALL_ISSUES` - Always send notifications, no matter whether the issue is muted or not

### Nested `destination` blocks

In order to get notified via a workflow you need two things:
- A [notification_destination](notification_destination.html) that defines reusable credentials/settings for a 
notification provider (e.g. webhook's `basic` credentials, Slack's OAuth credentials, PagerDuty API key, etc)
- A [notification_channel](notification_channel.html) that describes additional notification parameters to be used
for this specific workflow. Different destination types allow for different channel configuration options. 
For example, a webhook channel can define the payload template, while an email channel can define a subject as well as additional
details to include into the email body
  - **Please note that channels must be created with `product = "IINT"`** ("IINT" stands for "incident intelligence") in order to be used with workflows (see examples below)

Destinations can be reused across multiple channels. But each channel can only be used in a single workflow (at least for now).

When a workflow is deleted, all channels associated with it are also deleted. We might change this behaviour in the future, 
but please keep this behaviour in mind when managing workflows right now. 
If you only delete a workflow resource and not its channel resource, the next time you run TF it will report state drift.

Each workflow resource requires one or more `destination` blocks. These blocks define notification channels to use for the
workflow.

Block's arguments:
* `channel_id` - (Required) Id of a [notification_channel](notification_channel.html) to use for notifications. Please note that you have to use a 
**notification** channel, not an `alert_channel`.
* `notification_triggers` - (Optional) Issue events to notify on. The value is a list of possible issue events. See [Notification Triggers](#notification-triggers) below for details. 

### Notification Triggers

Each issue produces multiple events during its lifetime.
For example, issue activation, acknowledgement, and resolution are all separate events in issue's lifecycle.

It allows you to choose which events trigger notifications and which are ignored. This configuration is separate for each of the channels added to the workflow.
One could, for example, configure a workflow to open a Jira ticket and send Slack notifications once an issue is opened, but only send a Slack notification once it is closed.

Possible values:
* `ACTIVATED` - Send a notification when an issue is activated
* `ACKNOWLEDGED` - Send a notification when an issue is acknowledged
* `PRIORITY_CHANGED` - Send a notification when an issue's priority has been changed
* `CLOSED` - Send a notification when an issue is closed
* `OTHER_UPDATES` - Send a notification on other updates on the issue. These updates include:
1. An incident has been added to the issue
2. An incident in the issue has been closed
3. A different issue has been merged to this issue

    
### Nested `enrichments` blocks

Enrichments can give additional context on alert notifications by adding an [NRQL](https://docs.newrelic.com/docs/query-your-data/nrql-new-relic-query-language/get-started/introduction-nrql-new-relics-query-language/) query results to them.
Read more about enrichments in [workflows documentation](https://docs.newrelic.com/docs/alerts-applied-intelligence/applied-intelligence/incident-workflows/incident-workflows/#enrichments) 

`Enrichments` blocks have the following structure:
* `nrql` - a wrapper block 
  * `name` - A nrql enrichment name. This name can be used in your notification templates (see [notification_channel documentation](notification_channel.html))
  * `configuration` - Another wrapper block
    * `query` - An NRQL query to run


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the workflow.

## Import

Workflows can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_workflow.foo <id>
```
You can find the workflow ID from the workflow table by clicking on ... at the end of the row and choosing `Copy workflow id to clipboard`.

## Policy-Based Workflow Example
This scenario describes one of most common ways of using workflows by defining a set of policies the workflow handles

```hcl
// Create a policy to track
resource "newrelic_alert_policy" "my-policy" {
  name = "my_policy"
}

// Create a reusable notification destination
resource "newrelic_notification_destination" "webhook-destination" {
  name = "destination-webhook"
  type = "WEBHOOK"

  property {
    key = "url"
    value = "https://example.com"
  }

  auth_basic {
    user = "username"
    password = "password"
  }
}

// Create a notification channel to use in the workflow
resource "newrelic_notification_channel" "webhook-channel" {
  name = "channel-webhook"
  type = "WEBHOOK"
  destination_id = newrelic_notification_destination.webhook-destination.id
  product = "IINT" // Please note the product used!

  property {
    key = "payload"
    value = "{}"
    label = "Payload Template"
  }
}

// A workflow that matches issues that include incidents triggered by the policy
resource "newrelic_workflow" "workflow-example" {
  name = "workflow-example"
  muting_rules_handling = "NOTIFY_ALL_ISSUES"

  issues_filter {
    name = "Filter-name"
    type = "FILTER"

    predicate {
      attribute = "labels.policyIds"
      operator = "EXACTLY_MATCHES"
      values = [ newrelic_alert_policy.my-policy.id ]
    }
  }

  destination {
    channel_id = newrelic_notification_channel.webhook-channel.id
  }
}
```

### An example of a workflow with enrichments

```hcl
resource "newrelic_workflow" "workflow-example" {
  name = "workflow-enrichment-example"
  muting_rules_handling = "NOTIFY_ALL_ISSUES"
  
  issues_filter {
  name = "Filter-name"
  type = "FILTER"
      predicate {
        attribute = "accumulations.tag.team"
        operator = "EXACTLY_MATCHES"
        values = [ "my_team" ]
      }
  }
  
  enrichments {
    nrql {
      name = "Log Count"
      configuration {
        query = "SELECT count(*) FROM Log WHERE message like '%error%' since 10 minutes ago"
      }
    }
  }
  
  destination {
    channel_id = newrelic_notification_channel.webhook-channel.id
  }
}
```

### An example of a workflow with notification triggers

```hcl
resource "newrelic_workflow" "workflow-example" {
  name = "workflow-enrichment-example"
  muting_rules_handling = "NOTIFY_ALL_ISSUES"
  
  issues_filter {
  name = "Filter-name"
  type = "FILTER"
      predicate {
        attribute = "accumulations.tag.team"
        operator = "EXACTLY_MATCHES"
        values = [ "my_team" ]
      }
  }
  
  destination {
    channel_id = newrelic_notification_channel.webhook-channel.id
    notification_triggers = [ "ACTIVATED" ]
  }
}
```

## Additional Information
More details about the workflows can be found [here](https://docs.newrelic.com/docs/alerts-applied-intelligence/applied-intelligence/incident-workflows/incident-workflows/).

## v3.3 changes
In version v3.3 we renamed the following arguments:

- `workflow_enabled` changed to `enabled`.
- `destination_configuration` changed to `destination`.
- `predicates` changed to `predicate`.
- Enrichment's `configurations` changed to `configuration`.
