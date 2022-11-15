---
layout: "newrelic"
page_title: "New Relic: newrelic_obfuscation_expression"
sidebar_current: "docs-newrelic-resource-obfuscation-expression"
description: |-
Create and manage Obfuscation Expression.
---

# Resource: newrelic\_obfuscation\_expression

Use this resource to create, update and delete New Relic Obfuscation Expressions.


## Example Usage

```hcl
resource "newrelic_obfuscation_expression" "foo"{ 
  account_id = 12345
  name = "OExp"
  description = "The description"
  regex = "(regex.*)"
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The account id associated with the obfuscation expression.
* `description` - (Optional) Description of expression.
* `name` - (Required) Name of expression.
* `regex` - (Required) Regex of expression.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the obfuscation expression.

## Import

New Relic obfuscation expression can be imported using the expression ID, e.g.

```bash
$ terraform import newrelic_obfuscation_expression.foo 34567
```