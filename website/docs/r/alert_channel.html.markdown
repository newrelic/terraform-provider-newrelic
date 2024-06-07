---
layout: "newrelic"
page_title: "New Relic: newrelic_alert_channel"
sidebar_current: "docs-newrelic-resource-alert-channel"
description: |-
  Create and manage a notification channel for alerts in New Relic.
---

# Resource: newrelic\_alert\_channel

Use this resource to create and manage New Relic alert channels.

-> **WARNING:** The `newrelic_alert_channel` resource is **deprecated** and will be **removed in a future major release**. As an alternative, you can set up channels using a combination of the newer resources [`newrelic_notification_destination`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/notification_destination) and [`newrelic_notification_channel`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/notification_channel). We **strongly recommend** migrating to these new resources at the earliest. Please refer to [this example](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/getting_started#add-a-notification-channel) for a detailed illustration on setting up channels with these resources.

## Example Usage

##### Email
```hcl
resource "newrelic_alert_channel" "foo" {
  name = "foo"
  type = "email"

  config {
    recipients              = "foo@example.com"
    include_json_attachment = "true"
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

  * `account_id` - (Optional) Determines the New Relic account where the alert channel will be created. Defaults to the account associated with the API key used.
  * `name` - (Required) The name of the channel.
  * `type` - (Required) The type of channel.  One of: `email`, `slack`, `opsgenie`, `pagerduty`, `victorops`, or `webhook`.
  * `config` - (Optional) A nested block that describes an alert channel configuration.  Only one config block is permitted per alert channel definition.  See [Nested config blocks](#nested-config-blocks) below for details.

### Nested `config` blocks

Each alert channel type supports a specific set of arguments for the `config` block:

  * `email`
    * `recipients` - (Required) Comma delimited list of email addresses.
    * `include_json_attachment` - (Optional) `true` or `false`. Flag for whether or not to attach a JSON document containing information about the associated alert to the email that is sent to recipients.
  * `webhook`
    * `base_url` - (Required) The base URL of the webhook destination.
    * `auth_password` - (Optional) Specifies an authentication password for use with a channel.  Supported by the `webhook` channel type.
    * `auth_type` - (Optional) Specifies an authentication method for use with a channel.  Supported by the `webhook` channel type.  Only HTTP basic authentication is currently supported via the value `BASIC`.
    * `auth_username` - (Optional) Specifies an authentication username for use with a channel.  Supported by the `webhook` channel type.
    * `headers` - (Optional) A map of key/value pairs that represents extra HTTP headers to be sent along with the webhook payload.
    * `headers_string` - (Optional) Use instead of `headers` if the desired payload is more complex than a list of key/value pairs (e.g. a set of headers that makes use of nested objects).  The value provided should be a valid JSON string with escaped double quotes. Conflicts with `headers`.
    * `payload` - (Optional) A map of key/value pairs that represents the webhook payload.  Must provide `payload_type` if setting this argument.
    * `payload_string` - (Optional) Use instead of `payload` if the desired payload is more complex than a list of key/value pairs (e.g. a payload that makes use of nested objects).  The value provided should be a valid JSON string with escaped double quotes. Conflicts with `payload`.
    * `payload_type` - (Optional) Can either be `application/json` or `application/x-www-form-urlencoded`. The `payload_type` argument is _required_ if `payload` is set.
  * `pagerduty`
    * `service_key` - (Required) Specifies the service key for integrating with Pagerduty.
  * `victorops`
    * `key` - (Required) The key for integrating with VictorOps.
    * `route_key` - (Required) The route key for integrating with VictorOps.
  * `slack`
    * `url` - (Required) [Slack Webhook URL](https://api.slack.com/messaging/webhooks#create_a_webhook).
    * `channel` - (Optional) The Slack channel to send notifications to.
  * `opsgenie`
    * `api_key` - (Required) The API key for integrating with OpsGenie.
    * `region` - (Required) The data center region to store your data.  Valid values are `US` and `EU`.  Default is `US`.
    * `teams` - (Optional) A set of teams for targeting notifications. Multiple values are comma separated.
    * `tags` - (Optional) A set of tags for targeting notifications. Multiple values are comma separated.
    * `recipients` - (Optional) A set of recipients for targeting notifications.  Multiple values are comma separated.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The ID of the channel.

## Additional Examples

##### Slack
```hcl
resource "newrelic_alert_channel" "foo" {
  name = "slack-example"
  type = "slack"

  config {
    url     = "https://hooks.slack.com/services/XXXXXXX/XXXXXXX/XXXXXXXXXX"
    channel = "example-alerts-channel"
  }
}
```

-> **NOTE:** For instructions on setting up Webhooks with Slack, please visit the article linked under the argument `slack` in the aforementioned configuration, or [this article](https://docs.newrelic.com/docs/alerts-applied-intelligence/new-relic-alerts/alert-notifications/notification-channels-control-where-send-alerts/#slack) in New Relic's docs for additional details on setting up the `New Relic Alerts` Slack application, and subsequently using the generated Webhook URL.

##### OpsGenie
```hcl
resource "newrelic_alert_channel" "foo" {
  name = "opsgenie-example"
  type = "opsgenie"

  config {
    api_key    = "abc123"
    teams      = "team1, team2"
    tags       = "tag1, tag2"
    recipients = "user1@domain.com, user2@domain.com"
  }
}
```

##### PagerDuty
```hcl
resource "newrelic_alert_channel" "foo" {
  name = "pagerduty-example"
  type = "pagerduty"

  config {
    service_key = "abc123"
  }
}
```

##### VictorOps
```hcl
resource "newrelic_alert_channel" "foo" {
  name = "victorops-example"
  type = "victorops"

  config {
    key       = "abc123"
    route_key = "/example"
  }
}
```

##### Webhook
```hcl
resource "newrelic_alert_channel" "foo" {
  name = "webhook-example"
  type = "webhook"

  config {
    base_url = "http://www.test.com"
    payload_type = "application/json"
    payload = {
      condition_name = "$CONDITION_NAME"
      policy_name = "$POLICY_NAME"
    }

    headers = {
      header1 = value1
      header2 = value2
    }
  }
}
```

##### Webhook with complex payload
```hcl
resource "newrelic_alert_channel" "foo" {
  name = "webhook-example"
  type = "webhook"

  config {
    base_url = "http://www.test.com"
    payload_type = "application/json"
    payload_string = <<EOF
{
  "my_custom_values": {
    "condition_name": "$CONDITION_NAME",
    "policy_name": "$POLICY_NAME"
  }
}
EOF
  }
}
```

## Import

Alert channels can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_alert_channel.main <id>
```

~> **NOTE:** Sensitive data such as channel API keys, service keys, etc are not returned from the underlying API for security reasons and may not be set in state when importing.
