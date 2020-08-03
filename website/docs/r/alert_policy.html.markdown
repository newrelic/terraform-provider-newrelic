---
layout: "newrelic"
page_title: "New Relic: newrelic_alert_policy"
sidebar_current: "docs-newrelic-resource-alert-policy"
description: |-
  Create and manage alert policies in New Relic.
---

# Resource: newrelic\_alert\_policy

Use this resource to create and manage New Relic alert policies.

## Example Usage

```hcl
resource "newrelic_alert_policy" "foo" {
  name = "example"
  incident_preference = "PER_POLICY" # PER_POLICY is default
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the policy.
  * `incident_preference` - (Optional) The rollup strategy for the policy.  Options include: `PER_POLICY`, `PER_CONDITION`, or `PER_CONDITION_AND_TARGET`.  The default is `PER_POLICY`.
  * `channel_ids` - (Optional) An array of channel IDs (integers) to assign to the policy. Adding or removing channel IDs from this array will result in a new alert policy resource being created and the old one being destroyed. Also note that channel IDs _cannot_ be imported via `terraform import` (see [Import](#import) for info).
  * `account_id` - (Optional) The New Relic account ID to operate on.  This allows the user to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The ID of the policy.

## Additional Examples

##### Provision multiple notification channels and add those channels to a policy
```hcl
# Provision a Slack notification channel.
resource "newrelic_alert_channel" "slack_channel" {
  name = "slack-example"
  type = "slack"

  config {
    url     = "https://hooks.slack.com/services/<*****>/<*****>"
    channel = "example-alerts-channel"
  }
}

# Provision an email notification channel.
resource "newrelic_alert_channel" "email_channel" {
  name = "email-example"
  type = "email"

  config {
    recipients              = "example@testing.com"
    include_json_attachment = "1"
  }
}

# Provision the alert policy.
resource "newrelic_alert_policy" "policy_with_channels" {
  name                = "example-with-channels"
  incident_preference = "PER_CONDITION"

  # Add the provisioned channels to the policy.
  channel_ids = [
    newrelic_alert_channel.slack_channel.id,
    newrelic_alert_channel.email_channel.id,
  ]
}
```
<br>

##### Reference existing notification channels and add those channel to a policy
```hcl
# Reference an existing Slack notification channel.
data "newrelic_alert_channel" "slack_channel" {
  name = "slack-channel-notification"
}

# Reference an existing email notification channel.
data "newrelic_alert_channel" "email_channel" {
  name = "test@example.com"
}

# Provision the alert policy.
resource "newrelic_alert_policy" "policy_with_channels" {
  name                = "example-with-channels"
  incident_preference = "PER_CONDITION"

  # Add the referenced channels to the policy.
  channel_ids = [
    data.newrelic_alert_channel.slack_channel.id,
    data.newrelic_alert_channel.email_channel.id,
  ]
}
```

## Import

Alert policies can be imported using a composite ID of `<id>:<account_id>`, where `account_id` is the account number scoped to the alert policy resource.

Example import:

```
$ terraform import newrelic_alert_policy.foo 23423556:4593020
```

Please note that channel IDs (`channel_ids`) _cannot_ be imported due channels being a separate resource. However, to add channels to an imported alert policy, you can import the policy, add the `channel_ids` attribute with the associated channel IDs, then run `terraform apply`. This will result in the original alert policy being destroyed and a new alert policy being created along with the channels being added to the policy.
