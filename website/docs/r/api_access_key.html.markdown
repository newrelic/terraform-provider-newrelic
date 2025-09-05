---
layout: 'newrelic'
page_title: 'New Relic: newrelic_api_access_key'
sidebar_current: 'docs-newrelic-resource-api-access-key'
description: |-
  Create and manage New Relic API access keys
---

# Resource: newrelic_api_access_key

Use this resource to programmatically create and manage the following types of keys:
- [User API keys](https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys#user-api-key)
- License (or ingest) keys, including:
    - General [license key](https://docs.newrelic.com/docs/accounts/install-new-relic/account-setup/license-key) used for APM
    - [Browser license key](https://docs.newrelic.com/docs/browser/new-relic-browser/configuration/copy-browser-monitoring-license-key-app-id)

Please visit the New Relic article ['Use NerdGraph to manage license keys and User API keys'](https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-user-keys)
for more information.

-> Key exposure behavior: The `key` attribute is returned only at creation time and is marked Sensitive. Subsequent refreshes (plan/apply) do not re-read or overwrite this value. Store it securely when you create the key.

## Example Usage (returns key at create time)

```terraform
resource "newrelic_api_access_key" "apm_license_key" {
  account_id  = 1234567
  key_type    = "INGEST"
  ingest_type = "LICENSE"
  name        = "APM Ingest License Key for Service X"
  notes       = "Managed by Terraform, used for Service X APM agent."
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
- `key` - The actual API key. This attribute is Sensitive and is only populated on creation.

## Important Considerations

- Updating existing keys: Only `name` and `notes` are updatable in place. Changing other attributes will recreate the key and invalidate the old one.
- Account type restrictions for ingest keys: Creating `INGEST` keys requires a New Relic user with core or full platform access. See [user types](https://docs.newrelic.com/docs/accounts/accounts-billing/new-relic-one-user-management/user-type/#api-access).

## Import

Existing API access keys can be imported using a composite ID of `<api_access_key_id>:<key_type>`. `<key_type>`
will be either `INGEST` or `USER`.

For example:
```
$ terraform import newrelic_api_access_key.foobar "1234567:INGEST"
```

## Extended Usage With Modules

Looking for module-based workflows to create and then retrieve key values? See the dedicated guide:

- API Access Key usage with modules: https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/api_access_key_usage

Note: For `key_type = USER`, New Relic only allows a user to fetch their own key value. The module patterns work for USER keys created for oneself, not for other users.
