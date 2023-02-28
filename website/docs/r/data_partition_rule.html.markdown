---
layout: "newrelic"
page_title: "New Relic: newrelic_data_partition_rule"
sidebar_current: "docs-newrelic-resource-data-partition-rule"
description: |-
Create and manage Data partition rule.
---

# Resource: newrelic\_data\_partition\_rule

Use this resource to create, update and delete New Relic Data partition rule.


## Example Usage

```hcl

resource "newrelic_data_partition_rule" "foo"{
  description = "description"
  enabled = true
  attribute_name = "Name"
  matching_expression = "expression"
  matching_method = "EQUALS"
  nrql = "logtype='node'"
  retention_policy = "STANDARD"
  target_data_partition = "Log_name"
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The account id associated with the data partition rule.
* `description` - (Optional) The description of the data partition rule.
* `enabled` - (Required) Whether or not this data partition rule is enabled.
* `attribute_name` - (Required) The attribute name against which this matching condition will be evaluated.
* `matching_expression` - (Required) The matching expression of the data partition rule definition.
* `matching_method` - (Required) The matching method of the data partition rule definition. Valid values are `EQUALS` and `LIKE`.
* `nrql` - (Required) The NRQL to match events for this data partition rule. Logs matching this criteria will be routed to the specified data partition.
* `retention_policy` - (Required) The retention policy of the data partition data. Valid values are `SECONDARY` and `STANDARD`.
* `target_data_partition` - (Required) The name of the data partition where logs will be allocated once the rule is enabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the data partition rule.
* `deleted` - Whether or not this data partition rule is deleted. Deleting a data partition rule does not delete the already persisted data. This data will be retained for a given period of time specified in the retention policy field.

## Import

New Relic data partition rule can be imported using the rule ID, e.g.

```bash
$ terraform import newrelic_data_partition_rule.foo <id>
```

## Additional Information

More details about the data partition can be found [here](https://docs.newrelic.com/docs/logs/ui-data/data-partitions/)