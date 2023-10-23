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
  violation_time_limit_seconds = 3600

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
  * `violation_time_limit_seconds` - (Optional) The maximum number of seconds a violation can remain open before being closed by the system. The value must be between 300 seconds (5 minutes) to 2592000 seconds (30 days), both inclusive. Defaults to 259200 seconds (3 days) if this argument is not specified in the configuration, in accordance with the characteristics of this field in NerdGraph, as specified in the [docs](https://docs.newrelic.com/docs/alerts-applied-intelligence/new-relic-alerts/advanced-alerts/rest-api-alerts/alerts-conditions-api-field-names/#violation_time_limit_seconds).
  * `entities` - (Required) The Monitor GUID's of the Synthetics monitors to alert on.
  * `critical` - (Required) A condition term with the priority set to critical.
  * `warning` - (Optional) A condition term with the priority set to warning.


-> **WARNING:** This resource will use the account ID linked to your API key. At the moment it is not possible to dynamically set the account ID.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `entity_guid` - The unique entity identifier of the condition in New Relic.

## Import

New Relic Synthetics MultiLocation Conditions can be imported using a concatenated string of the format
 `<policy_id>:<condition_id>`, e.g.

```bash
$ terraform import newrelic_synthetics_multilocation_alert_condition.example 12345678:1456
```

## Tags

Manage synthetics multilocation alert condition tags with `newrelic_entity_tags`. For up-to-date documentation about the tagging resource, please check [newrelic_entity_tags](entity_tags.html#example-usage)

```hcl
resource "newrelic_alert_policy" "foo" {
  name = "foo policy"
}

resource "newrelic_synthetics_monitor" "foo" {
  status           = "ENABLED"
  name             = "foo monitor"
  period           = "EVERY_MINUTE"
  uri              = "https://www.one.newrelic.com"
  type             = "SIMPLE"
  locations_public = ["AP_EAST_1"]

  custom_header {
    name  = "some_name"
    value = "some_value"
  }

  treat_redirect_as_failure = true
  validation_string         = "success"
  bypass_head_request       = true
  verify_ssl                = true

  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}

resource "newrelic_synthetics_multilocation_alert_condition" "foo" {
  policy_id = newrelic_alert_policy.foo.id

  name                         = "foo condition"
  runbook_url                  = "https://example.com"
  enabled                      = true
  violation_time_limit_seconds = 3600

  entities = [
    newrelic_synthetics_monitor.foo.id
  ]

  critical {
    threshold = 2
  }

  warning {
    threshold = 1
  }
}


resource "newrelic_entity_tags" "my_condition_entity_tags" {
  guid = newrelic_synthetics_multilocation_alert_condition.foo.entity_guid

  tag {
    key = "my-key"
    values = ["my-value", "my-other-value"]
  }

  tag {
    key = "my-key-2"
    values = ["my-value-2"]
  }
}
```

