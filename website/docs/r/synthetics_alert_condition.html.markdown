---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_alert_condition"
sidebar_current: "docs-newrelic-resource-synthetics-alert-condition"
description: |-
  Create and manage a Synthetics alert condition for a policy in New Relic.
---

# newrelic\_synthetics\_alert\_condition

## Example Usage

```hcl

# Discover an existing monitor by name and create a new alert condition
data "newrelic_synthetics_monitor" "existing" {
  name = "Existing Monitor"
}

resource "newrelic_synthetics_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"
  name        = "foo"
  monitor_id  = "${data.newrelic_synthetics_monitor.existing.id}"
  runbook_url = "https://www.example.com"
}


# Or create a new monitor resource from scratch and create a new alert condition
resource "newrelic_synthetics_monitor" "new" {
    name                = "New Monitor"
    frequency           = 10
    uri                 = "https://www.example.com"
    locations           = ["AWS_US_EAST_1","AWS_US_WEST_2"]
    status              = "ENABLED"
    sla_threshold       = 2
}

resource "newrelic_synthetics_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"
  name        = "foo"
  monitor_id  = "${newrelic_synthetics_monitor.new.id}"
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

The following attributes are exported:

  * `id` - The ID of the Synthetics alert condition.