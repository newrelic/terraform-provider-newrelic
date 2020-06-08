---
layout: 'newrelic'
page_title: 'New Relic: newrelic_alert_policy_channel'
sidebar_current: 'docs-newrelic-resource-alert-policy-channel'
description: |-
  Map alert policies to alert channels in New Relic.
---

# Resource: newrelic_alert_policy_channel

Use this resource to map alert policies to alert channels in New Relic.

## Example Usage

The example below will apply multiple alert channels to an existing New Relic alert policy.

```hcl
# Fetches the data for this policy from your New Relic account
# and is referenced in the newrelic_alert_policy_channel block below.
data "newrelic_alert_policy" "example_policy" {
  name = "my-alert-policy"
}

# Creates an email alert channel.
resource "newrelic_alert_channel" "email_channel" {
  name = "bar"
  type = "email"

  config {
    recipients              = "foo@example.com"
    include_json_attachment = "1"
  }
}

# Creates a Slack alert channel.
resource "newrelic_alert_channel" "slack_channel" {
  name = "slack-channel-example"
  type = "slack"

  config {
    channel = "#example-channel"
    url     = "http://example-org.slack.com"
  }
}

# Applies the created channels above to the alert policy
# referenced at the top of the config.
resource "newrelic_alert_policy_channel" "foo" {
  policy_id  = newrelic_alert_policy.example_policy.id
  channel_ids = [
    data.newrelic_alert_channel.email_channel.id,
    data.newrelic_alert_channel.slack_channel.id
  ]
}
```

## Argument Reference

The following arguments are supported:

- `policy_id` - (Required) The ID of the policy.
- `channel_ids` - (Required) Array of channel IDs to apply to the specified policy. We recommended sorting channel IDs in ascending order to avoid drift your Terraform state.

## Import

Alert policy channels can be imported using the following notation: `<policyID>:<channelID>:<channelID>`, e.g.

```
$ terraform import newrelic_alert_policy_channel.foo 123456:3462754:2938324
```

When importing `newrelic_alert_policy_channel` resource, the attribute `channel_ids`\* will be set in your Terraform state. You can import multiple channels as long as those channel IDs are included as part of the import ID hash.

