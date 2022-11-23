---
layout: "newrelic"
page_title: "New Relic: newrelic_obfuscation_expression"
sidebar_current: "docs-newrelic-datasource-obfuscation-expression"
description: |-
Grabs a Obfuscation Expression by name.
---

# Data Source: newrelic\_obfuscation\_expression

Use this data source to get information about a specific Obfuscation Expression in New Relic that already exists.

## Example Usage

```hcl
data "newrelic_obfuscation_expression" "expression" {
  account_id = 123456
  name       = "The expression"
}

resource "newrelic_obfuscation_rule" "rule" {
  name        = "ruleName"
  description = "description of the rule"
  filter      = "hostStatus=running"
  enabled     = true
  // Reference the obfuscation expression data source in the obfuscation rule resource
  action {
    attribute    = ["message"]
    expression_id = data.newrelic_obfuscation_expression.expression.id
    method       = "MASK"
  }
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The account id associated with the obfuscation expression. If left empty will default to account ID specified in provider level configuration.
* `name` - (Required) Name of expression.
