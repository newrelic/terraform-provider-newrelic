---
layout: "newrelic"
page_title: "New Relic: newrelic_notification_channel"
sidebar_current: "docs-newrelic-resource-notification-channel"
description: |-
Create and manage a notification channel for notifications in New Relic.
---

# Resource: newrelic\_notification\_channel

Use this resource to create and manage New Relic notification channels. Details regarding supported products and permissions can be found [here](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/destinations).

## Example Usage

##### [Webhook](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-api-notifications-channels/#webhook)
```hcl
resource "newrelic_notification_channel" "foo" {
  account_id = 12345678
  name = "webhook-example"
  type = "WEBHOOK"
  destination_id = "00b6bd1d-ac06-4d3d-bd72-49551e70f7a8"
  product = "IINT"

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
* `type` - (Required) The type of channel.  One of: `EMAIL`, `SERVICENOW_INCIDENTS`, `WEBHOOK`, `JIRA_CLASSIC`, `JIRA_NEXTGEN`, `PAGERDUTY_ACCOUNT_INTEGRATION` or `PAGERDUTY_SERVICE_INTEGRATION`.
* `destination_id` - (Required) The id of the destination.
* `product` - (Required) The type of product.  One of: `DISCUSSIONS`, `ERROR_TRACKING` or `IINT` (workflows).
* `property` - A nested block that describes a notification channel property. See [Nested property blocks](#nested-property-blocks) below for details.

### Nested `property` blocks

* `key` - (Required) The notification property key.
* `value` - (Required) The notification property value.
* `label` - (Optional) The notification property label.
* `display_value` - (Optional) The notification property display value.
* 
Each notification channel type supports a specific set of arguments for the `property` block:

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
* `SLACK`
  * `channelId` - (Required) The slack channel id you want to sent to. Must be of the referenced destination.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the channel.

## Additional Examples

##### [ServiceNow](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-api-notifications-channels/#servicenow)
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

##### [Email](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-api-notifications-channels/#email)
```hcl
resource "newrelic_notification_channel" "foo" {
  account_id = 12345678
  name = "email-example"
  type = "EMAIL"
  destination_id = "00b6bd1d-ac06-4d3d-bd72-49551e70f7a8"
  product = "ERROR_TRACKING"

  property {
    key = "subject"
    value = "New Subject Title"
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
    value = "10000"
  }

  property {
    key = "issuetype"
    value = "10004"
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
    value = "PTQK3FM"
  }

  property {
    key = "email"
    value = "example@email.com"
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
    value = "https://webhook.site/94193c01-4a81-4782-8f1b-554d5230395b"
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

## Additional Information
More details about the channels API can be found [here](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-api-notifications-channels).