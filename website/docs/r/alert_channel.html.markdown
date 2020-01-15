---
layout: "newrelic"
page_title: "New Relic: newrelic_alert_channel"
sidebar_current: "docs-newrelic-resource-alert-channel"
description: |-
  Create and manage a notification channel for alerts in New Relic.
---

# Resource: newrelic\_alert\_channel

Use this resource to create and manage New Relic alert policies.

## Example Usage

##### Email
```hcl
resource "newrelic_alert_channel" "foo" {
  name = "foo"
  type = "email"

  config {
    recipients              = "foo@example.com"
    include_json_attachment = "1"
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the channel.
  * `type` - (Required) The type of channel.  One of: `email`, `slack`, `opsgenie`, `pagerduty`, `victorops`, or `webhook`.
  * `config` - (Optional) A nested block that describes an alert channel configuration.  Only one config block is permitted per alert channel definition.  See [Nested config blocks](#nested-`config`-blocks) below for details.
  * `configuration` - **Deprecated** (Optional) A map of key/value pairs with channel type specific values. This argument is deprecated.  Use the `config` argument instead.

### Nested `config` blocks

Each alert channel type supports a specific set of arguments for the `config` block:

  * `email`
    * `recipients` - (Required) Specifies the targeted email recipients.
    * `include_json_attachment` - (Optional) Allows a JSON attachment to be included with an alert email.
  * `webhook`
    * `base_url` - (Required) The base URL of the webhook destination.
    * `auth_password` - (Optional) Specifies an authentication password for use with a channel.  Supported by the `webhook` channel type.
    * `auth_type` - (Optional) Specifies an authentication method for use with a channel.  Supported by the `webhook` channel type.  Only HTTP basic authentication is currently supported via the value `BASIC`.
    * `auth_username` - (Optional) Specifies an authentication username for use with a channel.  Supported by the `webhook` channel type.
    * `headers` - (Optional) A map of key/value pairs that represents extra HTTP headers to be sent along with the webhook payload
    * `payload` - (Optional) A map of key/value pairs that represents the webhook payload.
  * `pagerduty`
    * `service_key` - (Required) Specifies the service key for integrating with Pagerduty.
  * `victorops`
    * `key` - (Required) The key for integrating with VictorOps.
    * `route_key` - (Required) The route key for integrating with VictorOps.
  * `slack`
    * `url` - (Required) The URL for integrating with Slack.
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

## Channel Configurations

Each supported channel supports a particular set of configuration arguments.

  * `email`
    * `recipients` - (Required) Comma delimited list of email addresses.
    * `include_json_attachment` - (Optional) `0` or `1`. Flag for whether or not to attach a JSON document containing information about the associated alert to the email that is sent to recipients. Default: `0`
  * `slack`
    * `url` - (Required) Your organization's Slack URL.
    * `channel` - (Required) The Slack channel for which to send notifications.
  * `opsgenie`
    * `api_key` - (Required) Your OpsGenie API key.
    * `teams` - (Optional) Comma delimited list of teams.
    * `tags` - (Optional) Comma delimited list of tags.
    * `recipients` - (Optional) Comma delimited list of email addresses.
  * `pagerduty`
    * `service_key` - (Required) Your PagerDuty service key.
  * `victorops`
    * `key` - (Required) Your VictorOps key.
    * `route_key` - (Required) The route for which to send notifications.

## Additional Examples

##### Slack
```hcl
resource "newrelic_alert_channel" "foo" {
  name = "slack-example"
  type = "slack"

  configuration = {
    url     = "https://<YourOrganization>.slack.com"
    channel = "example-alerts-channel"
  }
}
```

##### OpsGenie
```hcl
resource "newrelic_alert_channel" "foo" {
  name = "opsgenie-example"
  type = "opsgenie"

  configuration = {
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

  configuration = {
    service_key = "abc123"
  }
}
```

##### VictorOps
```hcl
resource "newrelic_alert_channel" "foo" {
  name = "victorops-example"
  type = "victorops"

  configuration = {
    key       = "abc123"
    route_key = "/example"
  }
}
```

## Import

Alert channels can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_alert_channel.main <id>
```
