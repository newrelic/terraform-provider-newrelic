---
layout: "newrelic"
page_title: "New Relic: newrelic_account"
sidebar_current: "docs-newrelic-datasource-account"
description: |-
  Grabs a New Relic account.
---

# Data Source: newrelic\_account

This data source allows you to retrieve information about a specific account in New Relic.

## Overview

You can locate accounts using either their `account_id` or `name`. However, only one of these attributes can be specified at a time. If neither attribute is provided, the provider's default `account_id` will be used.

## Example Usage

```hcl
data "newrelic_account" "example" {
  name = "Test Account"
}
```

## Argument Reference
The following arguments are supported:
- `account_id` - (Optional) The unique identifier of the account in New Relic. This must be an integer.
- `name` - (Optional) The name of the account in New Relic. This must be a string.