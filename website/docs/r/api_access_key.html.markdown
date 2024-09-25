---
layout: 'newrelic'
page_title: 'New Relic: newrelic_api_access_key'
sidebar_current: 'docs-newrelic-resource-api-access-key'
description: |-
  Create and Manage New Relic API access keys
---

# Resource: newrelic_api_access_key

Use this resource to programmatically create and manage the following types of keys:
- [User API keys](https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys#user-api-key)
- License (or ingest) keys, including:
    - General [license key](https://docs.newrelic.com/docs/accounts/install-new-relic/account-setup/license-key) used for APM
    - [Browser license key](https://docs.newrelic.com/docs/browser/new-relic-browser/configuration/copy-browser-monitoring-license-key-app-id)

Please visit the New Relic article ['Use NerdGraph to manage license keys and User API keys'](https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-user-keys)
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
  name        = "APM Ingest License Key"
  notes       = "To be used with service X"
}
```


-> **WARNING:** Creating 'Ingest - License' and 'Ingest - Browser' keys using this resource is restricted to 'core' or 'full platform' New Relic user accounts. If you've signed up as a 'basic' user with New Relic, or have been added as a 'basic' user to your organization on New Relic, you would not be able to use your account to create 'Ingest' keys. If you see the message `"You do not have permission to create this key"` in the response of the API called by this resource, it could be owing to the aforementioned. For more insights into user account types on New Relic and associated privileges, please check out this [page](https://docs.newrelic.com/docs/accounts/accounts-billing/new-relic-one-user-management/user-type/#api-access).


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
## Extended Usage
This module may be used to create a user or ingest key using the `create_access_keys_service` resource, and fetch the created key using `fetch_access_keys_service`, by performing a NerdGraph query under the hood, using the ID of the key created via the resource to fetch the created key.
Please refer  
[create access keys and fetch access keys](https://github.com/newrelic/terraform-provider-newrelic/blob/main/examples/modules/newrelic_api_access_key_extended/README.md) for more info.