---
layout: "newrelic"
page_title: "New Relic: newrelic_composite_service_level"
sidebar_current: "docs-newrelic-resource-composite-service-level"
description: |-
  Create and manage a New Relic Composite Service Level.
---

# Resource: newrelic\_composite\_service\_level

Use this resource to create, update, and delete New Relic Composite Service Levels — a Service Level that rolls up multiple existing Service Levels (or other composites) into a single, weighted reliability metric.

A New Relic User API key is required to provision this resource. Set the `api_key`
attribute in the `provider` block or the `NEW_RELIC_API_KEY` environment
variable with your User API key.

Important:
- Only roles that provide [permissions](https://docs.newrelic.com/docs/accounts/accounts-billing/new-relic-one-user-management/new-relic-one-user-model-understand-user-structure/) to create events to metric rules can create composite Service Levels.
- Only [Full users](https://docs.newrelic.com/docs/accounts/accounts-billing/new-relic-one-user-management/new-relic-one-user-model-understand-user-structure/#user-type) can view composite Service Levels.
- **Reference other resources' GUID output for `component.guid` (e.g. `newrelic_service_level.foo.sli_guid`, `newrelic_composite_service_level.bar.guid`), never a hardcoded GUID string.** A composite cannot be deleted while it's referenced as a component in another composite — referencing a resource's own output lets Terraform's dependency graph order creates/deletes correctly, where a hardcoded string would let Terraform attempt an out-of-order delete and hit that server-side guard.

## Example Usage

```hcl
resource "newrelic_service_level" "api_latency" {
  guid        = "MXxBUE18QVBQTElDQVRJT058MQ"
  name        = "API latency"
  description = "Proportion of API requests served faster than a threshold."

  events {
    account_id = 12345678
    valid_events {
      from = "Transaction"
    }
    good_events {
      from  = "Transaction"
      where = "duration < 0.5"
    }
  }

  objective {
    target = 99.9
    time_window {
      rolling {
        count = 7
        unit  = "DAY"
      }
    }
  }
}

resource "newrelic_service_level" "api_error_rate" {
  guid        = "MXxBUE18QVBQTElDQVRJT058Mg"
  name        = "API error rate"
  description = "Proportion of API requests that succeed."

  events {
    account_id = 12345678
    valid_events {
      from = "Transaction"
    }
    good_events {
      from  = "Transaction"
      where = "error IS FALSE"
    }
  }

  objective {
    target = 99.9
    time_window {
      rolling {
        count = 7
        unit  = "DAY"
      }
    }
  }
}

resource "newrelic_composite_service_level" "api_reliability" {
  account_id  = 12345678
  name        = "API reliability"
  description = "Composite reliability metric across API latency and error rate."

  component {
    guid = newrelic_service_level.api_latency.sli_guid
  }

  component {
    guid = newrelic_service_level.api_error_rate.sli_guid
  }

  objective {
    target = 99.9
    time_window {
      rolling {
        count = 7
        unit  = "DAY"
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

  * `account_id` - (Required) The New Relic account ID where the composite will reside. Note that changing the account ID will force a new resource.
  * `name` - (Required) A short name for the composite that will help anyone understand what it is about.
  * `component` - (Required) The component Service Levels that make up this composite. Must contain 2-20 components, with no duplicate `guid`. See [Component](#component) below for details.
  * `objective` - (Required) The objective of the composite, only one can be defined.
  See [Objective](#objective) below for details.
  * `description` - (Optional) The description of the composite.

### Component

All `component` blocks support the following arguments:

  * `guid` - (Required) The unique entity identifier of the component Service Level. **Reference another resource's GUID output rather than a hardcoded string** — `newrelic_service_level.foo.sli_guid` for a plain Service Level (note: `sli_guid`, not `guid` — on `newrelic_service_level`, `guid` is the *input* entity the SLI attaches to, not the SLI's own identity), or `newrelic_composite_service_level.bar.guid` for a nested composite. See the delete-ordering note above. Components may belong to different accounts than the composite itself, and a composite may reference another composite as a component (nesting is allowed up to a depth of 5, with cycle detection).
  * `weight` - (Optional) The weight of this component in the composite calculation, expressed as a percentage (0-100). Weights are all-or-nothing: if any component has a weight, every component must have one, and they must sum to exactly 100. If no weights are given, the server distributes weight equally across all components.

### Objective

  * `target` - (Required) The target of the objective, valid values between `0` and `100`.
  * `time_window` - (Required) Time window is the period of the objective. Must match the time window of every component (including nested composites').
    * `rolling` - (Required) Rolling window.
      * `count` - (Required) Valid values are `1`, `7` and `28`.
      * `unit` - (Required) The only supported value is `DAY`.

## Attributes Reference

The following attributes are exported:

  * `guid` - The unique entity identifier of the composite. This is also the resource's import ID.

## Import

New Relic Composite Service Levels can be imported using the composite's GUID. Unlike `newrelic_service_level`, which encodes the account ID and the related entity's GUID into a composite ID string, a composite reads, updates, and deletes directly by its own GUID — so the import ID is just that GUID, nothing else.

Example:

```bash
$ terraform import newrelic_composite_service_level.api_reliability MXxBUE18Q09NUE9TSVRFfDE
```
