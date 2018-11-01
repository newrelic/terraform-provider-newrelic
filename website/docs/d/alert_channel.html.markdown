---
layout: "newrelic"
page_title: "New Relic: newrelic_alert_channel"
sidebar_current: "docs-newrelic-datasource-alert-channel"
description: |-
  Looks up the information about an alert channel in New Relic.
---

# newrelic\_alert\_channel

Use this data source to get information about an specific alert channel in New Relic which already exists (e.g newrelic user).

## Example Usage

```hcl
data "newrelic_alert_channel" "foo" {
  name = "foo@example.com"
}

resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_alert_policy_channel" "foo" {
  policy_id  = "${newrelic_alert_policy.foo.id}"
  channel_id = "${newrelic_alert_channel.foo.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the alert channel in New Relic.

## Attributes Reference
* `id` - The ID of the alert channel.
* `type` - Alert channel type, either: `campfire`, `email`, `hipchat`, `opsgenie`, `pagerduty`, `slack`, `victorops`, or `webhook`..
* `policy_ids` - A list of policy IDs associated with the alert channel.
