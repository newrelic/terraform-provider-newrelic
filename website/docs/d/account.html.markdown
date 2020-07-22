---
layout: "newrelic"
page_title: "New Relic: newrelic_account"
sidebar_current: "docs-newrelic-datasource-account"
description: |-
  Grabs a New Relic account.
---

# Data Source: newrelic\_account

Use this data source to get information about a specific account in New Relic.
Accounts can be located by ID or name.  Exactly one of the two attributes is
required.

## Example Usage

```hcl
data "newrelic_account" "acc" {
  scope = "global"
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The account ID in New Relic.
* `name` - (Optional) The account name in New Relic.
* `scope` - (Optional) The scope of the account in New Relic.  Valid values are "global" and "in_region".  Defaults to "in_region".