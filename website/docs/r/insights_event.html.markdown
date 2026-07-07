---
layout: "newrelic"
page_title: "New Relic: newrelic_insights_event"
sidebar_current: "docs-newrelic-resource-insights-event"
description: |-
  Create one or more Insights events.
---

# Resource: newrelic\_insights\_event

Use this resource to create one or more Insights events during a terraform run.

-> **Region-aware — no manual endpoint override needed.** As of the release that ships this note, `newrelic_insights_event` sends events through the same region-aware client the rest of the provider uses. Setting `provider "newrelic" { region = "US" | "EU" | "JP" }` is enough — the correct Insights collector endpoint is selected automatically (`insights-collector.newrelic.com` for US, `insights-collector.eu01.nr-data.net` for EU, `insights-collector.jp.nr-data.net` for JP). The `insights_insert_url` provider argument is now **ignored** (with a warning in the logs) and will be removed in a future major release. If you previously set `insights_insert_url` as a workaround for JP or EU, you can safely remove it.

## Example Usage

```hcl
resource "newrelic_insights_event" "foo" {
  event {
    type = "MyEvent"

    timestamp = 1232471100

    attribute {
      key   = "a_string_attribute"
      value = "a string"
    }
    attribute {
      key   = "an_integer_attribute"
      value = 42
      type  = "int"
    }
    attribute {
      key   = "a_float_attribute"
      value = 101.1
      type  = "float"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

  * `event` - (Required) An event to insert into Insights. Multiple event blocks can be defined. See [Events](#events) below for details.

## Events

The `event` mapping supports the following arguments:

  * `type` - (Required) The event's name. Can be a combination of alphanumeric characters, underscores, and colons.
  * `timestamp` - (Optional) Must be a Unix epoch timestamp. You can define timestamps either in seconds or in milliseconds.
  * `attribute` - (Required) An attribute to include in your event payload. Multiple attribute blocks can be defined for an event. See [Attributes](#attributes) below for details.

### Attributes

The `attribute` mapping supports the following arguments:

  * `key` - (Required) The name of the attribute.
  * `value` - (Required) The value of the attribute.
  * `type` - (Optional) Specify the type for the attribute value. This is useful when passing integer or float values to Insights. Allowed values are `string`, `int`, or `float`. Defaults to `string`.
