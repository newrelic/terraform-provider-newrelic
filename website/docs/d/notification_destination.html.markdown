---
layout: "newrelic"
page_title: "New Relic: newrelic_notification_destination"
sidebar_current: "docs-newrelic-datasource-notification-destination"
description: |-
  Looks up the information about a notifications' destination data source in New Relic.
---

# Data Source: newrelic\_notification\_destination

Use this data source to get information about a specific notification destination in New Relic that already exists. More information on Terraform's data sources can be found [here](https://www.terraform.io/language/data-sources).

## ID Example Usage

```hcl
# Data source
data "newrelic_notification_destination" "foo" {
  id = "1e543419-0c25-456a-9057-fb0eb310e60b"
}

# Resource
resource "newrelic_notification_channel" "foo-channel" {
  name           = "webhook-example"
  type           = "WEBHOOK"
  destination_id = data.newrelic_notification_destination.foo.id
  product        = "IINT"

  property {
    key   = "payload"
    value = "{\n\t\"name\": \"foo\"\n}"
    label = "Payload Template"
  }
}
```

## Name Example Usage

```hcl
# Data source
data "newrelic_notification_destination" "foo" {
  name = "webhook-destination"
}

# Resource
resource "newrelic_notification_channel" "foo-channel" {
  name           = "webhook-example"
  type           = "WEBHOOK"
  destination_id = data.newrelic_notification_destination.foo.id
  product        = "IINT"

  property {
    key   = "payload"
    value = "{\n\t\"name\": \"foo\"\n}"
    label = "Payload Template"
  }
}
```

## Argument Reference

The following arguments are supported:

Either of the following two attributes are required, and not both:
* `id` - (Optional) The id of the notification destination in New Relic.
* `name` - (Optional) The name of the notification destination.

Optional:
* `account_id` - (Optional) The New Relic account ID to operate on.  This allows you to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The name of the notification destination.
* `type` - The notification destination type, either: `EMAIL`, `SERVICE_NOW`, `SERVICE_NOW_APP`, `WEBHOOK`, `JIRA`, `MOBILE_PUSH`, `EVENT_BRIDGE`, `PAGERDUTY_ACCOUNT_INTEGRATION` or `PAGERDUTY_SERVICE_INTEGRATION`, `SLACK` and `SLACK_COLLABORATION`.
* `property` - A nested block that describes a notification destination property.
* `active` - An indication whether the notification destination is active or not.
* `status` - The status of the notification destination.
* `guid` - The unique entity identifier of the destination in New Relic.
* `secure_url` - The URL in secure format, showing only the `prefix`, as the `secure_suffix` is a secret. 


```
Warning: This data source will use the account ID linked to your API key. At the moment it is not possible to dynamically set the account ID.
```
