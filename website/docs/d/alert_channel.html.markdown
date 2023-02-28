---
layout: "newrelic"
page_title: "New Relic: newrelic_alert_channel"
sidebar_current: "docs-newrelic-datasource-alert-channel"
description: |-
  Looks up the information about an alert channel in New Relic.
---

# Data Source: newrelic\_alert\_channel

Use this data source to get information about a specific alert channel in New Relic that already exists. More information on Terraform's data sources can be found [here](https://www.terraform.io/language/data-sources).

-> **WARNING:** The `newrelic_alert_channel` data source is deprecated and will be removed in the next major release.


## Example Usage

```hcl
# Data source
data "newrelic_alert_channel" "foo" {
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

* `name` - (Required) The name of the alert channel in New Relic.
* `account_id` - (Optional) The New Relic account ID to operate on.  This allows you to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the alert channel.
* `type` - Alert channel type, either: `email`, `opsgenie`, `pagerduty`, `slack`, `victorops`, or `webhook`.
* `config` - Alert channel configuration.
* `policy_ids` - A list of policy IDs associated with the alert channel.


```
Warning: This data source will use the account ID linked to your API key. At the moment it is not possible to dynamically set the account ID.
```
