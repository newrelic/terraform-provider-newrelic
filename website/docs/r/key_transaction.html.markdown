---
layout: "newrelic"
page_title: "New Relic: newrelic_key_transaction"
sidebar_current: "docs-newrelic-resource-key-transaction"
description: |-
    Create a new New Relic Key Transaction.
---

# Resource: newrelic\_key\_transaction

Use this resource to create a new Key Transaction in New Relic.

-> **NOTE:** For more information on Key Transactions, head over to [this page](https://docs.newrelic.com/docs/apm/transactions/key-transactions/introduction-key-transactions/) in New Relic's docs.

## Example Usage

```hcl
resource "newrelic_key_transaction" "foo" {
  application_guid     = "MzgfNjUyNnxBUE19QVBQTElDQVHJT068NTUfNDT4MjUy"
  apdex_index          = 0.5
  browser_apdex_target = 0.5
  metric_name          = "WebTransaction/Function/__main__:foo_bar"
  name                 = "Sample Key Transaction"
}
```
## Argument Reference

The following arguments are supported by this resource.

* `application_guid` - (Required) The GUID of the APM Application comprising transactions, of which one would be made a key transaction.
* `metric_name` - (Required) - The name of the underlying metric monitored by the key transaction to be created.
* `name` - (Required) - The name of the key transaction.
* `apdex_index` - (Required) A decimal value, measuring user satisfaction with response times, ranging from 0 (frustrated) to 1 (satisfied).
* `browser_apdex_target` - (Required) A decimal value representing the response time threshold for satisfactory experience (e.g., 0.5 seconds).

-> **NOTE:** It may be noted that the `metric_name` and `application_guid` of a Key Transaction _cannot_ be updated in a key transaction that has already been created; since this is not supported. As a consequence, altering the values of `application_guid` and/or `metric_name` of a `newrelic_key_transaction` resource created (to try updating these values) would result in `terraform plan` prompting a forced destruction and re-creation of the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported by this resource.

* `id` - The GUID of the created key transaction.
* `domain` - The domain of the entity monitored by the key transaction.
*  `type` - The type of the entity monitored by the key transaction.

## Import
A Key Transaction in New Relic may be imported into Terraform using its GUID specified in the `<id>` field, in the following command.

```bash
$ terraform import newrelic_key_transaction.foo <id>
```
