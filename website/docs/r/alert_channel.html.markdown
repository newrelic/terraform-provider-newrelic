---
layout: "newrelic"
page_title: "New Relic: newrelic_alert_channel"
sidebar_current: "docs-newrelic-resource-alert-channel"
description: |-
  Create and manage a notification channel for alerts in New Relic.
---

# newrelic\_alert\_channel

## Example Usage

Example email channel:
```hcl
resource "newrelic_alert_channel" "foo" {
  name = "foo"
  type = "email"

  configuration = {
    recipients              = "foo@example.com"
    include_json_attachment = "1"
  }
}
```

Example webhook channel:

```hcl
resource "newrelic_alert_channel" "bar" {
  name = "bar"
  type = "webhook"

  configuration = {
    base_url = "http://test.com",
    auth_username = "username",
    auth_password = "password",
    payload_type  = "application/json",
  }

  headers {
    api-key = "my-internal-api-key"
  }

  payload {
    # using all the default fields
    account_id                       = "$ACCOUNT_ID"
    account_name                     = "$ACCOUNT_NAME"
    closed_violations_count_critical = "$CLOSED_VIOLATIONS_COUNT_CRITICAL"
    closed_violations_count_warning  = "$CLOSED_VIOLATIONS_COUNT_WARNING"
    condition_family_id              = "$CONDITION_FAMILY_ID"
    #... some defaults not shown ...
    timestamp = "$TIMESTAMP"
    violation_callback_url           = "$VIOLATION_CALLBACK_URL"
    violation_chart_url              = "$VIOLATION_CHART_URL"

    # add our own custom field
    my_custom_field                  = "my custom value"
  }
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the channel.
  * `type` - (Required) The type of channel.  One of: `campfire`, `email`, `hipchat`, `opsgenie`, `pagerduty`, `slack`, `victorops`, or `webhook`.
  * `configuration` - (Required) A map of key / value pairs with channel type specific values.

Additionally, Webhook channels can have the following optional fields:

  * `headers` - (Optional) A map of key / value pairs of HTTP headers to add to the webhook request.
  * `payload` - (Optional) A map of key / value pairs to be sent as the webhook's payload. Fully replaces the [default payload](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/rest-api-calls-new-relic-alerts#webhook_json_channel). See [Webhook Documentation](https://docs.newrelic.com/docs/alerts/new-relic-alerts/managing-notification-channels/customize-your-webhook-payload#variables) for fields and values that New Relic provides.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the channel.

## Import

Alert channels can be imported using the `id`, e.g.

```
$ terraform import newrelic_alert_channel.main 12345
```
