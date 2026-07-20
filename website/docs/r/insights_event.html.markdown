---
layout: "newrelic"
page_title: "New Relic: newrelic_insights_event"
sidebar_current: "docs-newrelic-resource-insights-event"
description: |-
  Create one or more Insights events.
---

# Resource: newrelic\_insights\_event

Use this resource to create one or more Insights events during a terraform run.

-> **NOTE** <span style="color:red;">Starting <b>v3.95.0</b> of the New Relic Terraform Provider, the `insights_insert_url` provider argument is deprecated and will be removed in a future major release.</span><br><br>The correct Insights collector endpoint is now picked automatically from the provider's `region` argument, so no manual override is needed. If you were setting `insights_insert_url` as a workaround for the EU or JP region, you can safely remove it - values passed there are ignored with a warning in the logs.

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
