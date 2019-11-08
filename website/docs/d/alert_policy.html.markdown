---
layout: "newrelic"
page_title: "New Relic: newrelic_alert_policy"
sidebar_current: "docs-newrelic-datasource-alert-policy"
description: |-
  Looks up the information about an alert policy in New Relic.
---

# newrelic\_alert\_policy

Use this data source to get information about an specific alert policy in New Relic which already exists.

## Example Usage

```hcl
data "newrelic_alert_channel" "foo" {
  name = "foo@example.com"
}

data "newrelic_alert_policy" "foo" {
  name = "foo policy"
}

resource "newrelic_alert_policy_channel" "foo" {
  policy_id  = "${data.newrelic_alert_policy.foo.id}"
  channel_id = "${data.newrelic_alert_channel.foo.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the alert policy in New Relic.

## Attributes Reference

* `id` - The ID of the alert policy.
* `incident_preference` - The rollup strategy for the policy. Options include: PER_POLICY, PER_CONDITION, or PER_CONDITION_AND_TARGET. The default is PER_POLICY.
* `created_at` - The time the policy was created.
* `updated_at` -  The time the policy was last updated.
