---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_account"
sidebar_current: "docs-newrelic-datasource-cloud-account"
description: |-
    Grabs a cloud account linked to New Relic.
---

# Data Source: newrelic\_cloud\_account

Use this data source to get information about a specific cloud account linked to New Relic.
Accounts can be located by a combination of New Relic Account ID, name and cloud provider (aws, gcp, azure, etc). Name and cloud provider are required attributes. If no account_id is specified on the resource the provider level account_id will be used. 

## Example Usage

```hcl
data "newrelic_cloud_account" "account" {
  account_id = 12345
  cloud_provider = "aws"
  name = "my aws account"
}
```

### GCP Dimensional Metrics account lookup

To look up a GCP account linked for **GCP Dimensional Metrics** (keyless / Workload Identity Federation), set `cloud_provider = "gcp"` and `is_dimensional_metrics = true`. This flag tells the data source to match only Dimensional Metrics linked accounts, distinguishing them from legacy GCP integration accounts of the same name.

```hcl
data "newrelic_cloud_account" "gcp_dm_account" {
  account_id             = 12345
  cloud_provider         = "gcp"
  name                   = "my gcp dimensional metrics account"
  is_dimensional_metrics = true
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The account ID in New Relic.
* `cloud_provider` - (Required) The cloud provider of the account (aws, gcp, azure, etc)
* `name` - (Required) The cloud account name in New Relic.
* `is_dimensional_metrics` - (Optional) Set to `true` to look up a GCP **Dimensional Metrics** (keyless / Workload Identity Federation) linked account, as opposed to a legacy GCP integration account of the same name. Can only be used when `cloud_provider` is `"gcp"`. Defaults to `false`.