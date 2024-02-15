---
layout: "newrelic"
page_title: "New Relic: newrelic_alert_policy"
sidebar_current: "docs-newrelic-datasource-alert-policy"
description: |-
  Looks up the information about an alert policy in New Relic.
---

# Data Source: newrelic\_alert\_policy

Use this data source to get information about a specific alert policy in New Relic that already exists.
More information on Terraform's data sources can be found [here](https://www.terraform.io/language/data-sources).

## Example Usage

```hcl
data "newrelic_alert_channel" "foo" {
  name = "foo@example.com"
}

data "newrelic_alert_policy" "foo" {
  name = "foo policy"
}

resource "newrelic_alert_policy_channel" "foo" {
  policy_id  = data.newrelic_alert_policy.foo.id
  channel_id = data.newrelic_alert_channel.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the alert policy in New Relic.
* `account_id` - (Optional) The New Relic account ID to operate on.  This allows you to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the alert policy.
* `incident_preference` - The rollup strategy for the policy, which can have one of the following values:
  * `PER_POLICY` - Represents the incident grouping preference **One issue per policy**. Refer to [this page](https://docs.newrelic.com/docs/alerts-applied-intelligence/new-relic-alerts/alert-policies/specify-when-alerts-create-incidents/#preference-policy) for more details on this incident grouping preference.
  * `PER_CONDITION` - Represents the incident grouping preference **One issue per condition**. Refer to [this page](https://docs.newrelic.com/docs/alerts-applied-intelligence/new-relic-alerts/alert-policies/specify-when-alerts-create-incidents/#preference-condition) for more details on this incident grouping preference.
  * `PER_CONDITION_AND_TARGET` - Represents the incident grouping preference **One issue per condition and signal**. Refer to [this page](https://docs.newrelic.com/docs/alerts-applied-intelligence/new-relic-alerts/alert-policies/specify-when-alerts-create-incidents/#preference-signal) for more details on this incident grouping preference.
* `created_at` - The time the policy was created.
* `updated_at` -  The time the policy was last updated.
