---
layout: "newrelic"
page_title: "New Relic: newrelic_notifications_destination"
sidebar_current: "docs-newrelic-datasource-notifications-destination"
description: |-
  Looks up the information about a notifications' destination data source in New Relic.
---

# Data Source: newrelic\_notifications\_destination

Use this data source to get information about a specific notifications destination in New Relic that already exists. More information on Terraform's data sources can be found [here](https://www.terraform.io/language/data-sources).

## Example Usage

```hcl
# Data source
data "newrelic_notifications_destination" "foo" {
  name = "foo@example.com"
}

# Resource
resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

# Using the data source and resource together
resource "newrelic_alert_policy_channel" "foo" {
  policy_id  = newrelic_alert_policy.foo.id
  channel_id = data.newrelic_alert_channel.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The id of the notification destination in New Relic.
* `account_id` - (Optional) The New Relic account ID to operate on.  This allows you to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The name of the notification destination.
* `type` - The notification destination type, either: `EMAIL`, `SERVICE_NOW`, `WEBHOOK`, `JIRA`, `MOBILE_PUSH`, `EVENT_BRIDGE`, `PAGERDUTY_ACCOUNT_INTEGRATION` or `PAGERDUTY_SERVICE_INTEGRATION`, `SLACK` and `SLACK_COLLABORATION`.
* `auth_basic` - A nested block that describes a basic authentication credentials.
* `auth_token` - A nested block that describes a token authentication credentials.
* `property` - A nested block that describes a notification destination property.


```
Warning: This data source will use the account ID linked to your API key. At the moment it is not possible to dynamically set the account ID.
```
