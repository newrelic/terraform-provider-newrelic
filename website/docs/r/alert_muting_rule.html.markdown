---
layout: 'newrelic'
page_title: 'New Relic: newrelic_alert_muting_rule'
sidebar_current: 'docs-newrelic-resource-alert-muting-rule'
description: |-
  Create a muting rule for New Relic Alerts violations.
---

# Resource: newrelic_alert_muting_rule

Use this resource to create a muting rule for New Relic Alerts violations.

-> **IMPORTANT!** Version 2.0.0 of the New Relic Terraform Provider introduces some [additional requirements](/docs/providers/newrelic/index.html) for configuring the provider.
<br><br>
Before upgrading to version 2.0.0 or later, it is recommended to upgrade to the most recent 1.x version of the provider and ensure that your environment successfully runs `terraform plan` without unexpected changes.

## Example Usage

##### Type: `static` (default)
```hcl
resource "newrelic_alert_muting_rule" "foo" {
	name = "Example Muting Rule"
	enabled = true
	description = "muting rule test."
	condition {
		conditions {
			attribute 	= "product"
			operator 	= "EQUALS"
			values 		= ["APM"]
		}
		conditions {
			attribute 	= "targetId"
			operator 	= "EQUALS"
			values 		= ["Muted"]
		}
		operator = "AND"
	}
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:
  * `account_id` - (Optional) The account id of the MutingRule.
  * `condition`  - (Required) The condition that defines which violations to target. See [Nested condition blocks](#nested-condition-blocks) below for details.
  * `enabled` - (Required) Whether the MutingRule is enabled.
  * `name` - The name of the MutingRule.
  * `description` - The description of the MutingRule.


### Nested `condition` blocks

All nested `condition` blocks support the following arguments:
  * `conditions` - (Optional) The individual MutingRuleConditions within the group. See [Nested conditions blocks](#nested-conditions-blocks) below for details.
  * `operator` - (Required) The operator used to combine all the MutingRuleConditions within the group.


### Nested `conditions` blocks
* `attribute` - (Required) The attribute on a violation.
* `operator` - (Required) The operator used to compare the attribute's value with the supplied value(s)
* `values` - (Required) The value(s) to compare against the attribute's value.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:
<!-- 
- `id` - The ID of the NRQL alert condition. This is a composite ID with the format `<policy_id>:<condition_id>` - e.g. `538291:6789035`. -->


## Import

<!-- Alert conditions can be imported using a composite ID of `<policy_id>:<condition_id>:<conditionType>`, e.g. -->

