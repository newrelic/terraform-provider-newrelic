---
layout: 'newrelic'
page_title: 'New Relic: newrelic_api_access_key'
sidebar_current: 'docs-newrelic-datasource-api-access-key'
description: |-
  Look up the value of an existing New Relic API access key.
---

# Data Source: newrelic\_api\_access\_key

Use this data source to retrieve details of an existing New Relic API access key, which can be either an [Ingest (license) key](https://docs.newrelic.com/docs/accounts/install-new-relic/account-setup/license-key) or a [User API key](https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys#user-api-key). A common use case is to fetch the default license key that New Relic creates for every account (for instance, via the [`newrelic_account_management`](../resources/account_management) resource) and inject its value into other resources, without having to manage the key itself with the [`newrelic_api_access_key`](../resources/api_access_key) resource.

A key may be looked up in one of two ways: directly by its `key_id`, or by searching with a combination of `account_id`, `key_type`, `ingest_type`, `user_id` and `name`.

Refer to the New Relic article ['Use NerdGraph to manage license keys and User API keys'](https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-user-keys) for detailed information.

-> **WARNING:** If the key being fetched was created for a user other than the one whose API key is being used to run Terraform, the New Relic API returns a <span style="color:tomato;">truncated key value</span> for security reasons. In such a case, the `key` attribute is not populated in state and a warning is emitted. See the [Attributes Reference](#attributes-reference) section below, and the New Relic article ['Use NerdGraph to manage license keys and User API keys'](https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-user-keys/#query-keys) for further details.

## Example Usage

### Example: Look up the default license key of an account

```hcl
data "newrelic_api_access_key" "default" {
  account_id  = 1234567
  key_type    = "INGEST"
  ingest_type = "LICENSE"
  name        = "License Key for <Account name>"
}
```

### Example: Look up a key directly by its ID

```hcl
data "newrelic_api_access_key" "by_id" {
  key_id   = "131313133A331313130B5F13DF01313FDB13B13133EE5E133D13EAAB3A3C13D3"
  key_type = "INGEST"
}
```

### Example: Inject a license key into another resource

```hcl
data "newrelic_api_access_key" "ingest" {
  account_id  = 1234567
  key_type    = "INGEST"
  ingest_type = "LICENSE"
}

output "license_key" {
  value     = data.newrelic_api_access_key.ingest.key
  sensitive = true
}
```

## Argument Reference

The following arguments are supported:

- `key_type` - (Required) The type of the key to look up. Valid options are `INGEST` or `USER` (case-sensitive).
- `key_id` - (Optional) The ID of the key to look up. When specified, the key is fetched directly by its ID and all other search arguments are ignored.
- `account_id` - (Optional) The New Relic account ID to search for keys in. Defaults to the `account_id` set in the `provider{}` block (or the `NEW_RELIC_ACCOUNT_ID` environment variable) if not specified.
- `ingest_type` - (Optional) Filters the search to ingest keys of the given type. Valid options are `LICENSE` or `BROWSER` (case-sensitive). Only applies when `key_type` is `INGEST`.
- `user_id` - (Optional) Filters the search to user keys owned by the given New Relic user ID. Only applies when `key_type` is `USER`.
- `name` - (Optional) Filters the search to keys with the given name.

-> **NOTE:** When `key_id` is not specified, the search must resolve to exactly one key. If zero or more than one key matches the given criteria, the data source returns an error. In such a case, provide additional arguments (such as `name`, `ingest_type` or `user_id`) to identify a unique key, or specify `key_id` directly to bypass the search.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the API key. Identical in value to `key_id`.
- `key_id` - The ID of the API key.
- `key` - The value of the API key. This attribute is marked as **sensitive**.
  - <span style="color:tomato;">It is important to exercise caution when exporting the value of `key`, as it is sensitive information</span>. Avoid logging or exposing it inappropriately.
  - Not populated when the New Relic API returns a truncated key value; see the warning at the top of this page.
- `account_id` - The New Relic account ID the key belongs to.
- `key_type` - The type of the key (`INGEST` or `USER`).
- `ingest_type` - The type of the ingest key (`LICENSE` or `BROWSER`). Populated only when `key_type` is `INGEST`.
- `user_id` - The New Relic user ID that owns the key. Populated only when `key_type` is `USER`.
- `name` - The name of the API key.
- `notes` - Any notes attached to the API key.
