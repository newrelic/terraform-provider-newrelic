---
layout: "newrelic"
page_title: "New Relic: newrelic_insights_event"
sidebar_current: "docs-newrelic-resource-insights-event"
description: |-
  Create one or more Insights events.
---

# Resource: newrelic\_insights\_event

Use this resource to create one or more Insights events during a terraform run.

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
