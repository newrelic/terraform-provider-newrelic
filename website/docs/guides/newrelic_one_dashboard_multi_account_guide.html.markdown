---
layout: "newrelic"
page_title: "New Relic: `newrelic_one_dashboard` - Multi-Account Query Upgrade Guide"
sidebar_current: "docs-newrelic-provider-one-dashboard-multi-account-guide"
description: |-
  A guide for using the enhanced `account_id` attribute to support both single and multi-account NRQL queries in `newrelic_one_dashboard` widgets.
---

## Guide: Upgrading to Multi-Account Queries with the **account_id** Attribute in newrelic_one_dashboard resource

The `newrelic_one_dashboard` resource allows you to create and manage dashboards in New Relic. Within this resource, the [`account_id`](providers/newrelic/newrelic/latest/docs/resources/one_dashboard#account_id-2) attribute in the [`nrql_query`](providers/newrelic/newrelic/latest/docs/resources/one_dashboard#nested-nrql_query-blocks) block (used with various widgets) is used to specify the account from which data is queried for widgets. To enhance its functionality, with v3.65.0 of the New Relic Terraform Provider, the `account_id` attribute has been upgraded to support querying across multiple accounts in a single widget.

This guide explains what has changed, why it changed, and how to use the new, more powerful functionality.

### Why the Change?

Previously, each `nrql_query` could only target a single New Relic account. The primary motivation for this update was to **enable multi-account NRQL queries**, a highly requested addition that allows you to build much more powerful and consolidated dashboards, comprising widgets powered by cross-account queries.

To achieve this without introducing breaking changes or confusing new attributes, we've made the existing `account_id` attribute more flexible.

--- 

### How to Use the New `account_id`

The `account_id` attribute now accepts both a single ID (as before) and a list of IDs for multi-account queries.

### For a Single Account ID

There is **no change** to how you configure a single account. You can continue to provide the account ID as a plain number.

```hcl
resource "newrelic_one_dashboard" "example" {
  # ...
  page {
    widget_bar {
      title = "Single Account Widget"
      nrql_query {
        # This syntax remains the same.
        account_id = 1234567
        query      = "SELECT count(*) FROM Transaction"
      }
    }
  }
}
```

### For Multiple Account IDs (New)

To query multiple accounts, provide a list of numbers to the `account_id` attribute using the Terraform built-in [`jsonencode()`](https://developer.hashicorp.com/terraform/language/functions/jsonencode) function.

💡 Tip: Using `jsonencode()` is the standard and safest way to provide a list or complex type as a string to a Terraform attribute - and is hence, the recommended/supported approach for the `account_id` attribute to contain multiple account IDs.

```hcl
resource "newrelic_one_dashboard" "example" {
  page {
    widget_line {
      title = "Multi-Account Widget"
      nrql_query {
        # Use jsonencode() to provide a list of account IDs.
        account_id = jsonencode([1234567, 9876543, 5554443])
        query      = "SELECT count(*) FROM Transaction"
      }
    }
  }
}
```

## 🚀 Upgrading from a Previous Version
We have designed this change to be 100% backward-compatible. Your existing dashboard configurations will continue to work as usual, without any manual changes required.

However, here's a full example showing how you might evolve your HCL file to use `nrql_query` with multiple account IDs.

### Before (Old Provider Version, prior to v3.65.0)

```hcl
resource "newrelic_one_dashboard" "my_dashboard" {
  name = "Production Overview"
  page {
    name = "Services"
    widget_table {
      title  = "Transaction Errors (App A)"
      row    = 1
      column = 1
      nrql_query {
        account_id = 1111111
        query      = "SELECT count(*) FROM TransactionError FACET appName"
      }
    }
  }
}
```
### After (New Provider Version, v.3.65.0+)
Notice how the existing widget for App A is untouched. We've simply added a new multi-account widget - you can modify/update the existing single account query to contain multiple account IDs too, with the `jsonencode()` syntax as shown.

```hcl
resource "newrelic_one_dashboard" "my_dashboard" {
  name = "Production Overview"
  page {
    name = "Services"

    # This existing widget requires NO changes.
    widget_table {
      title  = "Transaction Errors (App A)"
      row    = 1
      column = 1
      nrql_query {
        account_id = 1111111 # Still works perfectly.
        query      = "SELECT count(*) FROM TransactionError FACET appName"
      }
    }

    # NEW: A second widget querying multiple accounts.
    widget_line {
      title  = "Combined API Gateway Throughput"
      row    = 1
      column = 2
      nrql_query {
        account_id = jsonencode([1111111, 2222222, 3333333])
        query      = "SELECT rate(count(*), 1 minute) FROM ApiGatewaySample"
      }
    }
  }
}
```

