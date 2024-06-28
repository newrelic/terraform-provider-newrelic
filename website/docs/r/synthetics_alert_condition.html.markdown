---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_alert_condition"
sidebar_current: "docs-newrelic-resource-synthetics-alert-condition"
description: |-
  Create and manage a Synthetics alert condition for a policy in New Relic.
---

# Resource: newrelic\_synthetics\_alert\_condition

Use this resource to create and manage synthetics alert conditions in New Relic.

-> **WARNING:** The `newrelic_synthetics_alert_condition` resource is deprecated and will be removed in the next major release. The resource [newrelic_nrql_alert_condition](nrql_alert_condition.html) would be a preferred alternative to configure alert conditions - in most cases, feature parity can be achieved with a NRQL query.For more details and examples on moving away from Synthetics alert conditions to the NRQL based alternative, please check out [this](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/migration_guide_alert_conditions#migrating-from-synthetics-alert-conditions-to-nrql-alert-conditions) example.

## Example Usage

```hcl
resource "newrelic_synthetics_alert_condition" "foo" {
  policy_id = newrelic_alert_policy.foo.id

  name        = "foo"
  monitor_id  = newrelic_synthetics_monitor.foo.id
  runbook_url = "https://www.example.com"
}
```

## Argument Reference

The following arguments are supported:

  * `policy_id` - (Required) The ID of the policy where this condition should be used.
  * `name` - (Required) The title of this condition.
  * `monitor_id` - (Required) The GUID of the Synthetics monitor to be referenced in the alert condition.
  * `runbook_url` - (Optional) Runbook URL to display in notifications.
  * `enabled` - (Optional) Set whether to enable the alert condition. Defaults to `true`.

```
Warning: This resource will use the account ID linked to your API key. At the moment it is not possible to dynamically set the account ID.
```

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The ID of the Synthetics alert condition.
  * `entity_guid` - The unique entity identifier of the condition in New Relic.


## Import

Synthetics alert conditions can be imported using a composite ID of `<policy_id>:<condition_id>`, e.g.

```
$ terraform import newrelic_synthetics_alert_condition.main 12345:67890
```

## Tags

Manage synthetics alert condition tags with `newrelic_entity_tags`. For up-to-date documentation about the tagging resource, please check [newrelic_entity_tags](entity_tags.html#example-usage)

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

resource "newrelic_synthetics_alert_condition" "foo" {
  policy_id = newrelic_alert_policy.foo.id

  name        = "foo synthetics condition"
  monitor_id  = newrelic_synthetics_monitor.foo.id
  runbook_url = "https://www.example.com"
}

resource "newrelic_entity_tags" "my_condition_entity_tags" {
  guid = newrelic_synthetics_alert_condition.foo.entity_guid


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
