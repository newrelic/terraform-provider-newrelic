---
layout: 'newrelic'
page_title: 'New Relic: newrelic_api_access_key'
sidebar_current: 'docs-newrelic-resource-api-access-key'
description: |-
  Create and Manage New Relic API access keys
---

# Resource; newrelic_api_access_key

Use this resource to programmatically create and manage the following types of keys:
- [Personal API keys](https://docs.newrelic.co.jp/docs/apis/get-started/intro-apis/types-new-relic-api-keys#personal-api-key)
- License (or ingest) keys, including:
    - General [license key](https://docs.newrelic.co.jp/docs/accounts/install-new-relic/account-setup/license-key) used for APM
    - [Browser license key](https://docs.newrelic.co.jp/docs/browser/new-relic-browser/configuration/copy-browser-monitoring-license-key-app-id)

Please visit the New Relic article ['Use NerdGraph to manage license keys and personal API keys'](https://docs.newrelic.co.jp/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-personal-api-keys)
for more information.

-> **IMPORTANT!**
Please be very careful when updating existing `newrelic_api_access_key` resources as only `newrelic_api_access_key.name`
and `newrelic_api_access_key.notes` are updatable. All other resource attributes will force a resource recreation which will
invalidate the previous API key(s).

## Example Usage
```hcl-terraform
    resource "newrelic_api_access_key" "foobar" {
        account_id  = 1234567
        key_type    = "INGEST"
        ingest_type = "LICENSE"
        name        = "APM Ingest License Key (blue)"
        notes       = "To be used with service X"
    }
```

## Argument Reference

The following arguments are supported:

- `account_id` - (Required) The New Relic account ID of the account you wish to create the API access key.
- `key_type` - (Required) What type of API key to create. Valid options are `INGEST` or `USER`, case-sensitive.
- `ingest_type` - (Optional) Required if `key_type = INGEST`. Valid options are `BROWSER` or `LICENSE`, case-sensitive.
- `user_id` - (Optional) Required if `key_type = USER`. The New Relic user ID yous wish to create the API access key for in an account.
- `name` - (Optional) The name of the key.
- `notes` - (Optional) Any notes about this ingest key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the API key.
- `key` - The actual API key. This attribute is masked and not be visible in your terminal, CI, etc.

## Import

Existing API access keys can be imported using a composite ID of `<api_access_key_id>:<key_type>`. `<key_type>`
will be either `INGEST` or `USER`.

For example:
```
$ terraform import newrelic_api_access_key.foobar "1234567:INGEST"
```