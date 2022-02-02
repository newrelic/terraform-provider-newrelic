---
layout: "newrelic"
page_title: "New Relic: newrelic_service_level"
sidebar_current: "docs-newrelic-resource-service-level"
description: |-
  Create and manage a New Relic Service Level.
---

-> **New Relic Service Level Management is in preview. [Read more](https://docs.newrelic.com/docs/service-level-management/intro-slm)**

# Resource: newrelic\_service\_level

Use this resource to create, update, and delete New Relic Service Level Indicators and Objectives.

A New Relic User API key is required to provision this resource.  Set the `api_key`
attribute in the `provider` block or the `NEW_RELIC_API_KEY` environment
variable with your User API key.

Important:
- Only roles that provide [permissions](https://docs.newrelic.com/docs/accounts/accounts-billing/new-relic-one-user-management/new-relic-one-user-model-understand-user-structure/) to create events to metric rules can create SLI/SLOs.
- Only [Full users](https://docs.newrelic.com/docs/accounts/accounts-billing/new-relic-one-user-management/new-relic-one-user-model-understand-user-structure/#user-type) can view SLI/SLOs.

## Example Usage

```hcl
resource "newrelic_service_level" "foo" {
    guid = "MXxBUE18QVBQTElDQVRJT058MQ"
    name = "Availability"
    description = "SLI that measures the availability of the service."

    events {
        account_id = 12345678
        valid_events {
            from = "Transaction"
            where = "appName = 'Example application'"
        }
        bad_events {
            from = "TransactionError"
            where = "appName = 'Example application' AND error.expected is false"
        }
    }

    objective {
        target = 99.00
        time_window {
            rolling {
                count = 7
                unit = "DAY"
            }
        }
    }
}
```

## Argument Reference

The following arguments are supported:

  * `guid` - (Required) The GUID of the entity (e.g, APM Service, Browser application, Workload, etc.) that you want to relate this SLI to. Note that changing the GUID will force a new resource.
  * `name` - (Required) A short name for the SLI that will help anyone understand what it is about.
  * `events` - (Required) The events that define the NRDB data for the SLI/SLO calculations.
  See [Events](#events) below for details.
  * `description` - (Optional) The description of the SLI.
  * `objective` - (Optional) An objective for the SLI. Multiple objective blocks can be defined for an SLI.
  See [Nested objective blocks](#nested-objective-blocks) below for details.

### Events

All nested `events` blocks support the following common arguments:

  * `account_id` - (Required) The ID of the account where the entity (e.g, APM Service, Browser application, Workload, etc.) belongs to,
  and that contains the NRDB data for the SLI/SLO calculations. 
  * `valid_events` - (Required) The definition of valid requests.
    * `from` - (Required) The event type where NRDB data will be fetched from.
    * `where` - (Optional) A filter that specifies all the NRDB events that are considered in this SLI (e.g, those that refer to a particular entity).
  * `good_events` - (Optional) The definition of good responses. If you define an SLI from valid and good events, you must leave the bad events argument empty.
    * `from` - (Required) The event type where NRDB data will be fetched from.
    * `where` - (Optional) A filter that narrows down the NRDB events just to those that are considered good responses (e.g, those that refer to
    a particular entity and were successful).
  * `bad_events` - (Optional) The definition of the bad responses. If you define an SLI from valid and bad events, you must leave the good events argument empty.
    * `from` - (Required) The event type where NRDB data will be fetched from.
    * `where` - (Optional) A filter that narrows down the NRDB events just to those that are considered bad responses (e.g, those that refer to
    a particular entity and returned an error).

### Nested `objective` blocks

  * `target` - (Required) The target for your SLO, valid values between `0` and `100`. Up to 5 decimals accepted.
  * `time_window` - (Required) Time window is the period for the SLO.
    * `rolling` - (Required) Rolling window.
      * `count` - (Required) Valid values are `1`, `7` and `28`.
      * `unit` - (Required) The only supported value is `DAY`.

## Attributes Reference

The following attributes are exported:

  * `sli_id` - The unique entity identifier of the Service Level Indicator.

## Import

New Relic Service Levels can be imported using a concatenated string of the format
 `<account_id>:<sli_id>:<guid>`, where the `guid` is the entity the SLI relates to.

Example:

```bash
$ terraform import newrelic_service_level.foo 12345678:4321:MXxBUE18QVBQTElDQVRJT058MQ
```
