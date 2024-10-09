---
layout: "newrelic"
page_title: "New Relic: newrelic_key_transaction"
sidebar_current: "docs-newrelic-datasource-key-transaction"
description: |-
  Looks up the information about a key transaction in New Relic.
---

# Data Source: newrelic\_key\_transaction

Use this data source to get information about a specific key transaction in New Relic that already exists.  More information on Terraform's data sources can be found [here](https://www.terraform.io/language/data-sources).

## Example Usage

```hcl
data "newrelic_key_transaction" "txn" {
  name = "txn"
}

resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_alert_condition" "foo" {
  policy_id = newrelic_alert_policy.foo.id

  name        = "foo"
  type        = "apm_kt_metric"
  entities    = [data.newrelic_key_transaction.txn.id]
  metric      = "error_percentage"
  runbook_url = "https://www.example.com"

  term {
    duration      = 5
    operator      = "below"
    priority      = "critical"
    threshold     = "0.75"
    time_function = "all"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the key transaction in New Relic.
* `guid` - (Optional) GUID of the key transaction in New Relic.

## Note
If two or more key transactions have the same name, the datasource fetch only one key transaction. To differentiate between them, use the GUID of the key transaction.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the application.
* `guid` - GUID of the key transaction in New Relic.
* `domain` - Domain of the key transaction in New Relic.
* `type` - Type of the key transaction in New Relic.
* `name` - Name of the key Transation in New Relic.

```
Warning: This data source will use the account ID linked to your API key. At the moment it is not possible to dynamically set the account ID.
```