---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_multilocation_alert_condition"
sidebar_current: "docs-newrelic-synthetics-multilocation-alert-condition"
description: |-
  Create and manage a New Relic Synthetics Location Alerts.
---

# Resource: newrelic\_synthetics\_multilocation\_alert\_condition

Use this resource to create, update, and delete a New Relic Synthetics Location Alerts.

-> **NOTE:** This is a legacy resource. The [newrelic_nrql_alert_condition](nrql_alert_condition.html) resource is preferred for configuring alerts conditions. In most cases feature parity can be achieved with a NRQL query. This condition type may be deprecated in the future.

## Example Usage

```hcl
resource "newrelic_alert_policy" "policy" {
  name = "my-policy"
}

resource "newrelic_synthetics_monitor" "monitor" {
  locations_public = ["US_WEST_1"]
  name             = "my-monitor"
  period           = "EVERY_10_MINUTES"
  status           = "DISABLED"
  type             = "SIMPLE"
  uri              = "https://www.one.newrelic.com"
}

resource "newrelic_synthetics_multilocation_alert_condition" "example" {
  policy_id = newrelic_alert_policy.policy.id

  name                         = "Example condition"
  runbook_url                  = "https://example.com"
  enabled                      = true
  violation_time_limit_seconds = "3600"

  entities = [
    newrelic_synthetics_monitor.monitor.id
  ]

  critical {
    threshold = 2
  }

  warning {
    threshold = 1
  }
}
```
## Argument Reference

The following arguments are supported:

  * `name` - (Required) The title of the condition.
  * `policy_id` - (Required) The ID of the policy where this condition will be used.
  * `runbook_url` - (Optional) Runbook URL to display in notifications.
  * `enabled` - (Optional) Set whether to enable the alert condition.  Defaults to true.
  * `violation_time_limit_seconds` - (Required) The maximum number of seconds a violation can remain open before being closed by the system. Must be one of: 0, 3600, 7200, 14400, 28800, 43200, 86400.
  * `entities` - (Required) The Monitor GUID's of the Synthetics monitors to alert on.
  * `critical` - (Required) A condition term with the priority set to critical.
  * `warning` - (Optional) A condition term with the priority set to warning.

```
Warning: This resource will use the account ID linked to your API key. At the moment it is not possible to dynamically set the account ID.
```

## Import

New Relic Synthetics MultiLocation Conditions can be imported using a concatenated string of the format
 `<policy_id>:<condition_id>`, e.g.

```bash
$ terraform import newrelic_synthetics_multilocation_alert_condition.example 12345678:1456
```
