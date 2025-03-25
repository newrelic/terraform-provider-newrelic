---
layout: "newrelic"
page_title: "New Relic: newrelic_notification_channel"
sidebar_current: "docs-newrelic-resource-notification-channel"
description: |-
Create and manage a notification channel for notifications in New Relic.
---

# Resource: newrelic\_notification\_channel

Use this resource to create and manage New Relic notification channels. Details regarding supported products and permissions can be found [here](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/destinations).

A channel is an entity that is used to configure notifications. It is also called a message template. It is a separate entity from workflows, but a channel is required in order to create a workflow.

## Example Usage

##### [Webhook](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-api-notifications-channels/#webhook)
```hcl
resource "newrelic_notification_channel" "foo" {
  account_id = 12345678
  name = "webhook-example"
  type = "WEBHOOK"
  destination_id = "00b6bd1d-ac06-4d3d-bd72-49551e70f7a8"
  product = "IINT" // (Workflows)

  // must be valid json
  property {
    key = "payload"
    value = "name: {{ foo }}"
    label = "Payload Template"
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) Determines the New Relic account where the notification channel will be created. Defaults to the account associated with the API key used.
* `name` - (Required) The name of the channel.
* `type` - (Required) The type of channel.  One of: `EMAIL`, `SERVICENOW_INCIDENTS`, `SERVICE_NOW_APP`, `WEBHOOK`, `JIRA_CLASSIC`, `MOBILE_PUSH`, `EVENT_BRIDGE`, `SLACK` and `SLACK_COLLABORATION`, `PAGERDUTY_ACCOUNT_INTEGRATION` or `PAGERDUTY_SERVICE_INTEGRATION`, `MICROSOFT_TEAMS`.
* `destination_id` - (Required) The id of the destination.
* `product` - (Required) The type of product.  One of: `DISCUSSIONS`, `ERROR_TRACKING` or `IINT` (workflows).
* `property` - A nested block that describes a notification channel property. See [Nested property blocks](#nested-property-blocks) below for details.

### Nested `property` blocks
Most properties can use variables, which will be filled at the time of sending the notification with data from the issue. The properties where this is not available generally correlate to identifiers in the third party, such as Slack channel id or Jira project id. 

* `key` - (Required) The notification property key.
* `value` - (Required) The notification property value.
* `label` - (Optional) The notification property label.
* `display_value` - (Optional) The notification property display value.

Each notification channel type supports a specific set of arguments for the `property` block:

* `WEBHOOK`
  * `headers` - (Optional) A map of key/value pairs that represents the webhook headers.
  * `payload` - (Required) A map of key/value pairs that represents the webhook payload.
* `SERVICENOW_INCIDENTS`
  * `description` - (Optional) Free text that represents a description.
  * `short_description` - (Optional) Free text that represents a short description.
* `SERVICE_NOW_APP` - No Properties.
* `JIRA_CLASSIC`, `JIRA_NEXTGEN`
  * `project` - (Required) Identifier that specifies jira project id.
  * `issuetype` - (Required) Identifier that specifies the issue type id.
  * `description` - (Required) Free text that represents a description.
  * `summary` - (Required) Free text that represents the summary.
* `EMAIL`
  * `subject` - (Optional) Free text that represents the email subject title.
  * `customDetailsEmail` - (Optional) Free text that represents the email custom details.
* `PAGERDUTY_SERVICE_INTEGRATION`
  * `summary` - (Required) Free text that represents the summary.
  * `customDetails` - (Optional) Free text that *replaces* the content of the alert.
* `PAGERDUTY_ACCOUNT_INTEGRATION`
  * `summary` - (Required) Free text that represents the summary.
  * `service` - (Required) Identifier that specifies the service id to alert to.
  * `email` - (Required) Specifies the user email for integrating with Pagerduty.
  * `customDetails` - (Optional) Free text that *replaces* the content of the alert.
* `SLACK`
  * `channelId` - (Required) Specifies the Slack channel id. This can be found in slack browser via the url. Example - https://app.slack.com/client/\<UserId>/\<ChannelId>.
  * `customDetailsSlack` - (Optional) A map of key/value pairs that represents the slack custom details. Must be compatible with Slack's blocks api. 
* `MICROSOFT_TEAMS`
  * `teamId` - (Required) Free text that represents the team id.
  * `channelId` - (Required) Identifier that specifies the channel id to alert to.
  * `customDetails` - (Optional) Free text that *replaces* the content of the alert.
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the channel.

## Additional Examples

~> **NOTE:** We support all properties. The mentioned properties are just an example.

##### [ServiceNow](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-api-notifications-channels/#servicenow)
To see the properties’ keys for your account, check ServiceNow incidents table.

```hcl
resource "newrelic_notification_channel" "foo" {
  account_id = 12345678
  name = "servicenow-incident-example"
  type = "SERVICENOW_INCIDENTS"
  destination_id = "00b6bd1d-ac06-4d3d-bd72-49551e70f7a8"
  product = "IINT"

  property {
    key = "description"
    value = "General description"
  }

  property {
    key = "short_description"
    value = "Short description"
  }
}
```
##### [ServiceNowApp](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-api-notifications-channels/#servicenowapp)
```hcl
resource "newrelic_notification_channel" "foo" {
  account_id = 12345678
  name = "servicenow-app-example"
  type = "SERVICE_NOW_APP"
  destination_id = "00b6bd1d-ac06-4d3d-bd72-49551e70f7a8"
  product = "IINT"
}
```

##### [Email](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-api-notifications-channels/#email)
```hcl
resource "newrelic_notification_channel" "foo" {
  account_id = 12345678
  name = "email-example"
  type = "EMAIL"
  destination_id = "00b6bd1d-ac06-4d3d-bd72-49551e70f7a8"
  product = "IINT"

  property {
    key = "subject"
    value = "New Subject Title"
  }

  property {
    key = "customDetailsEmail"
    value = "issue id - {{issueId}}"
  }
}
```

##### [Jira Classic](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-api-notifications-channels/#jira)
```hcl
resource "newrelic_notification_channel" "foo" {
  account_id = 12345678
  name = "jira-example"
  type = "JIRA_CLASSIC"
  destination_id = "00b6bd1d-ac06-4d3d-bd72-49551e70f7a8"
  product = "ERROR_TRACKING"

  property {
    key = "project"
    value = "10000" // id of the project, can be found using https://confluence.atlassian.com/jirakb/how-to-get-project-id-from-the-jira-user-interface-827341414.html
    label = "Project-Name"
  }

  property {
    key = "issuetype"
    value = "10004" // id of the issuetype, can be found using https://confluence.atlassian.com/jirakb/finding-the-id-for-issue-types-646186508.html
    label = "Bug"
  }

  property {
    key = "description"
    value = "Issue ID: {{ issueId }}"
  }

  property {
    key = "summary"
    value = "{{ annotations.title.[0] }}"
  }
}
```

##### [PagerDuty with account integration](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-api-notifications-channels/#pagerduty)
```hcl
resource "newrelic_notification_channel" "foo" {
  account_id = 12345678
  name = "pagerduty-account-example"
  type = "PAGERDUTY_ACCOUNT_INTEGRATION"
  destination_id = "00b6bd1d-ac06-4d3d-bd72-49551e70f7a8"
  product = "IINT"

  property {
    key = "summary"
    value = "General summary"
  }

  property {
    key = "service"
    label = "Service Name" // Label of service drop-down value in the channel builder. This will show as blank if left empty.
    value = "PTQK3FM"
  }

  property {
    key = "email"
    value = "example@email.com"
  }
  
  // must be valid json
  property {
    key   = "customDetails"
    value = <<-EOT
            {
            "id":{{json issueId}},
            "IssueURL":{{json issuePageUrl}},
            "NewRelic priority":{{json priority}},
            "Total Incidents":{{json totalIncidents}},
            "Impacted Entities":"{{#each entitiesData.names}}{{this}}{{#unless @last}}, {{/unless}}{{/each}}",
            "Runbook":"{{#each accumulations.runbookUrl}}{{this}}{{#unless @last}}, {{/unless}}{{/each}}",
            "Description":"{{#each annotations.description}}{{this}}{{#unless @last}}, {{/unless}}{{/each}}",
            "isCorrelated":{{json isCorrelated}},
            "Alert Policy Names":"{{#each accumulations.policyName}}{{this}}{{#unless @last}}, {{/unless}}{{/each}}",
            "Alert Condition Names":"{{#each accumulations.conditionName}}{{this}}{{#unless @last}}, {{/unless}}{{/each}}",
            "Workflow Name":{{json workflowName}}
            }
        EOT
  }
}
```

##### [PagerDuty with service integration](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-api-notifications-channels/#pagerduty)
```hcl
resource "newrelic_notification_channel" "foo" {
  account_id = 12345678
  name = "pagerduty-account-example"
  type = "PAGERDUTY_SERVICE_INTEGRATION"
  destination_id = "00b6bd1d-ac06-4d3d-bd72-49551e70f7a8"
  product = "IINT"

  property {
    key = "summary"
    value = "General summary"
  }
  // must be valid json
  property {
    key   = "customDetails"
    value = <<-EOT
            {
            "id":{{json issueId}},
            "IssueURL":{{json issuePageUrl}},
            "NewRelic priority":{{json priority}},
            "Total Incidents":{{json totalIncidents}},
            "Impacted Entities":"{{#each entitiesData.names}}{{this}}{{#unless @last}}, {{/unless}}{{/each}}",
            "Runbook":"{{#each accumulations.runbookUrl}}{{this}}{{#unless @last}}, {{/unless}}{{/each}}",
            "Description":"{{#each annotations.description}}{{this}}{{#unless @last}}, {{/unless}}{{/each}}",
            "isCorrelated":{{json isCorrelated}},
            "Alert Policy Names":"{{#each accumulations.policyName}}{{this}}{{#unless @last}}, {{/unless}}{{/each}}",
            "Alert Condition Names":"{{#each accumulations.conditionName}}{{this}}{{#unless @last}}, {{/unless}}{{/each}}",
            "Workflow Name":{{json workflowName}}
            }
        EOT
  }
}
```

#### Mobile Push
```hcl
resource "newrelic_notification_channel" "foo" {
  account_id = 12345678
  name = "mobile-push-example"
  type = "MOBILE_PUSH"
  destination_id = "00b6bd1d-ac06-4d3d-bd72-49551e70f7a8"
  product = "IINT"
}
```

#### [AWS Event Bridge](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-api-notifications-channels/#eventBridge)
```hcl
resource "newrelic_notification_channel" "foo" {
  account_id = 12345678
  name = "event-bridge-example"
  type = "EVENT_BRIDGE"
  destination_id = "00b6bd1d-ac06-4d3d-bd72-49551e70f7a8"
  product = "IINT"

  property {
    key = "eventSource"
    value = "aws.partner/mydomain/myaccountid/name"
  }
  // must be valid json
  property {
    key = "eventContent"
    value = "{ id: {{ json issueId }} }"
  }
}
```

#### [SLACK](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-api-notifications-channels/#slack)
```hcl
resource "newrelic_notification_channel" "foo" {
  account_id = 12345678
  name = "slack-example"
  type = "SLACK"
  destination_id = "00b6bd1d-ac06-4d3d-bd72-49551e70f7a8"
  product = "IINT"

  property {
    key = "channelId"
    value = "123456"
  }

  property {
    key = "customDetailsSlack"
    value = "issue id - {{issueId}}"
  }
}
```

#### [MICROSOFT_TEAMS]
```hcl
resource "newrelic_notification_channel" "foo" {
  account_id = 12345678
  name = "ms-teams-example"
  type = "MICROSOFT_TEAMS"
  destination_id = "52d36b54-3d68-4ac4-9d0a-581febb91f2c"
  product = "IINT"

  property {
    key = "teamId"
    value = "906379b4-f5ac-40fd-b242-d4faaa4d3963"
  }

  property {
    key = "channelId"
    value = "19:wk9tU4tSr335Y1cNiXOynredbi3lFoeabu0kybfmbBA1@thread.tacv2"
  }

  property {
    key = "customDetails"
    value = "{{ nrAccountId }}"
  }
}
```

~> **NOTE:** Sensitive data such as channel API keys, service keys, etc are not returned from the underlying API for security reasons and may not be set in state when importing.

## Full Scenario Example
Create a destination resource and reference that destination to the channel resource:

### Create a destination
```hcl
resource "newrelic_notification_destination" "webhook-destination" {
  account_id = 12345678
  name = "destination-webhook"
  type = "WEBHOOK"

  property {
    key = "url"
    value = "https://webhook.mywebhook.com"
  }

  auth_basic {
    user = "username"
    password = "password"
  }
}
```

### Create a channel
```hcl
resource "newrelic_notification_channel" "webhook-channel" {
  account_id = 12345678
  name = "channel-webhook"
  type = "WEBHOOK"
  destination_id = newrelic_notification_destination.webhook-destination.id
  product = "IINT"

  property {
    key = "payload"
    value = "{name: foo}"
    label = "Payload Template"
  }
}
```

## Import
Channels can only be used by a single workflow, therefore importing them is not particularly useful, because in the UI channels are created upon workflow creation. 
Additionally, the channel id isn't available via the UI, and you'd need to look it up with the `channels` query in the NerdGraph API.
That being said, importing is possible using -
```terraform import newrelic_notification_destination.foo <destination_id>```

## Additional Information
More details about the channels API can be found [here](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-api-notifications-channels).
### Moving from Legacy Alert Channels to Notification Channels
As described in the documentation of this resource, channels can be created and managed using [`newrelic_notification_destination`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/notification_destination) and [`newrelic_notification_channel`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/notification_channel). A combination of these resources is an alternative to the legacy resource [`newrelic_alert_channel`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/alert_channel), which is **deprecated** and will be **removed in a future major release**, as stated in the documentation of the resource.

If you're currently using `newrelic_alert_channel` to manage channels, we **strongly recommend** migrating to these notifications-based resources at the earliest.

Please refer to the examples in this page, or [this example](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/getting_started#add-a-notification-channel) for illustrations on setting up channels with these resources.
