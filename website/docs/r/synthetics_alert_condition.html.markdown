---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_alert_condition"
sidebar_current: "docs-newrelic-resource-synthetics-alert-condition"
description: |-
  Create and manage a Synthetics alert condition for a policy in New Relic.
---

# Resource: newrelic\_synthetics\_alert\_condition

## Example Usage

```hcl
data "newrelic_synthetics_monitor" "foo" {
  name = "foo"
}

resource "newrelic_synthetics_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name        = "foo"
  monitor_id  = "${data.newrelic_synthetics_monitor.foo.id}"
  runbook_url = "https://www.example.com"
}
```

## Argument Reference

The following arguments are supported:

  * `policy_id` - (Required) The ID of the policy where this condition should be used.
  * `name` - (Required) The title of this condition.
  * `monitor_id` - (Required) The ID of the Synthetics monitor to be referenced in the alert condition. 
  * `runbook_url` - (Optional) Runbook URL to display in notifications.
  * `enabled` - (Optional) Set whether to enable the alert condition. Defaults to `true`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The ID of the Synthetics alert condition.