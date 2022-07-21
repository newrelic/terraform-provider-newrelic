---
layout: "newrelic"
page_title: "New Relic: newrelic_notification_channel"
sidebar_current: "docs-newrelic-resource-notification-channel"
description: |-
Create and manage a notification channel for notifications in New Relic.
---

# Resource: newrelic\_notification\_channel

Use this resource to create and manage New Relic notification channels.

## Example Usage

##### Webhook
```hcl
resource "newrelic_notification_channel" "foo" {
  name = "webhook-example"
  type = "WEBHOOK"
  destination_id = "1234"
  product = "IINT"

  property {
    key = "payload"
    value = "{\n\t\"name\": \"foo\"\n}"
    label = "Payload Template"
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) Determines the New Relic account where the notification channel will be created. Defaults to the account associated with the API key used.
* `name` - (Required) The name of the channel.
* `type` - (Required) The type of channel.  One of: `EMAIL`, `SERVICENOW_INCIDENTS`, `WEBHOOK`, `PAGERDUTY_ACCOUNT_INTEGRATION` or `PAGERDUTY_SERVICE_INTEGRATION`.
* `destination_id` - (Required) The id of the destination.
* `product` - (Required) The type of product.  One of: `DISCUSSIONS`, `ERROR_TRACKING` or `IINT` (workflows).
* `property` - (Required) A nested block that describes a notification channel property. See [Nested property blocks](#nested-property-blocks) below for details.

### Nested `property` blocks

Each notification channel type supports a specific set of arguments for the `property` block:

* `WEBHOOK`
  * `headers` - (Optional) A map of key/value pairs that represents the webhook headers.
  * `payload` - (Required) A map of key/value pairs that represents the webhook payload.
* `SERVICENOW_INCIDENTS`
  * `description` - (Required) A map of key/value pairs that represents a description.
  * `short_description` - (Required) A map of key/value pairs that represents a short description.
* `PAGERDUTY_SERVICE_INTEGRATION`
  * `summary` - (Required) A map of key/value pairs that represents the summery.
* `PAGERDUTY_ACCOUNT_INTEGRATION`
  * `summary` - (Required) A map of key/value pairs that represents the summery.
  * `pagerduty_service` - (Required) Specifies the service key for integrating with Pagerduty.
  * `user_for_comment` - (Required) Specifies the service key for integrating with Pagerduty.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the channel.

## Additional Examples

##### ServiceNow
```hcl
resource "newrelic_notification_channel" "foo" {
  name = "servicenow-incident-example"
  type = "SERVICENOW_INCIDENTS"
  destination_id = "1234"
  product = "PD"

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

##### Email
```hcl
resource "newrelic_notification_channel" "foo" {
  name = "email-example"
  type = "EMAIL"
  destination_id = "1234"
  product = "ERROR_TRACKING"
}
```

##### PagerDuty with account integration
```hcl
resource "newrelic_notification_channel" "foo" {
  name = "pagerduty-account-example"
  type = "PAGERDUTY_ACCOUNT_INTEGRATION"
  destination_id = "1234"
  product = "IINT"

  property {
    key = "summary"
    value = "General summary"
  }

  property {
    key = "service"
    value = "1234"
  }

  property {
    key = "email"
    value = "test@test.com"
  }
}
```

##### PagerDuty with service integration
```hcl
resource "newrelic_notification_channel" "foo" {
  name = "pagerduty-account-example"
  type = "PAGERDUTY_SERVICE_INTEGRATION"
  destination_id = "1234"
  product = "IINT"

  property {
    key = "summary"
    value = "General summary"
  }
}
```

~> **NOTE:** Sensitive data such as channel API keys, service keys, etc are not returned from the underlying API for security reasons and may not be set in state when importing.

## Additional Information
More information can be found in NewRelic [documentation](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/notification-integrations/).
