---
layout: "newrelic"
page_title: "New Relic: newrelic_obfuscation_rule"
sidebar_current: "docs-newrelic-resource-obfuscation-rule"
description: |-
Create and manage Obfuscation Rule.
---

# Resource: newrelic\_obfuscation\_rule

Use this resource to create, update and delete New Relic Obfuscation Rule.


## Example Usage

```hcl

resource "newrelic_obfuscation_expression" "bar" {
  name        = "expressionName"
  description = "description of the expression"
  regex       = "(^http)"
}

resource "newrelic_obfuscation_rule" "foo" {
  name        = "ruleName"
  description = "description of the rule"
  filter      = "hostStatus=running"
  enabled     = true
  action {
    attribute    = ["message"]
    expression_id = newrelic_obfuscation_expression.bar.id
    method       = "MASK"
  }
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The account id associated with the obfuscation rule.
* `description` - (Optional) Description of rule.
* `name` - (Required) Name of rule.
* `filter` - (Required) NRQL for determining whether a given log record should have obfuscation actions applied.
* `enabled` - (Required) Whether the rule should be applied or not to incoming data.
* `action` - (Required) Actions for the rule. The actions will be applied in the order specified by this list.

### Nested `action` blocks

All nested `action` blocks support the following common arguments:

* `attribute` - (Required) Attribute names for action. An empty list applies the action to all the attributes.
* `expression_id` - (Required) Expression Id for action.
* `method` - (Required) Obfuscation method to use. Methods for replacing obfuscated values are `HASH_SHA256` and `MASK`.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the obfuscation rule.

## Import

New Relic obfuscation rule can be imported using the rule ID, e.g.

```bash
$ terraform import newrelic_obfuscation_rule.foo 34567
```