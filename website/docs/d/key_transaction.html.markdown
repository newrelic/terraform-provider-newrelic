---
layout: "newrelic"
page_title: "New Relic: newrelic_key_transaction"
sidebar_current: "docs-newrelic-datasource-key-transaction"
description: |-
  Looks up the information about a key transaction in New Relic.
---

# Data Source: newrelic\_key\_transaction

Use this data source to get information about a specific key transaction in New Relic.

## Example Usage

```hcl
data "newrelic_key_transaction" "txn" {
  name = "txn"
}

resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name        = "foo"
  type        = "apm_kt_metric"
  entities    = ["${data.newrelic_key_transaction.txn.id}"]
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

* `name` - (Required) The name of the application in New Relic.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the application.
