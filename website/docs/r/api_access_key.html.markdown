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
  - General (Ingest) [license keys](https://docs.newrelic.com/docs/accounts/install-new-relic/account-setup/license-key) used for APM
  - [Browser license keys](https://docs.newrelic.com/docs/browser/new-relic-browser/configuration/copy-browser-monitoring-license-key-app-id)

Refer to the New Relic article ['Use NerdGraph to manage license keys and User API keys'](https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-user-keys) for detailed information.

## Example Usage

```
resource "newrelic_api_access_key" "apm_license_key" {
  account_id  = 1234321
  key_type    = "INGEST"
  ingest_type = "LICENSE"
  name        = "User X's Ingest License Key"
  notes       = "A Terraform-managed Ingest License Key, used by apps reporting to New Relic"
}
```

## Argument Reference

The following arguments are supported:

- `account_id` - (Required) The New Relic account ID where the API access key will be created.
- `key_type` - (Required) The type of API key to create. Valid options are `INGEST` or `USER` (case-sensitive).
  - If `key_type` is `INGEST`, then `ingest_type` must be specified.
  - If `key_type` is `USER`, then `user_id` must be specified.
- `ingest_type` - (Optional) Required if `key_type` is `INGEST`. Valid options are `BROWSER` or `LICENSE` (case-sensitive).
- `user_id` - (Optional) Required if `key_type` is `USER`. The New Relic user ID for which the API access key will be created.
- `name` - (Optional) The name of the API key.
  - **Note**: While `name` is optional, it is <b style="color:red;">\*\*strongly recommended\*\*</b> to provide a meaningful name for easier identification and management of keys. If a `name` is not provided, the API will assign a default name when processing the request to create the API key, which may cause unexpected drift in your Terraform state. To prevent this, it is best practice to always specify a `name`.
- `notes` - (Optional) Additional notes about the API access key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the API key.
- `key` - The actual API key.
  - <span style="color:tomato;">It is important to exercise caution when exporting the value of `key`, as it is sensitive information</span>. Avoid logging or exposing it inappropriately.

## Important Considerations
#### Updating Existing Keys
- Only `name` and `notes` can be updated in place. Changes to other attributes will recreate the key (the `newrelic_api_access_key` resource), invalidating the existing one.

#### Creating API Keys for Other Users
- If an API key is created for a user other than the owner of the API key used to run Terraform, the full key value will not be returned by the API for security reasons. Instead, a truncated version of the key will be provided. To retrieve the full key, ensure the necessary capabilities and access management settings are applied to the user running Terraform. For more details, contact New Relic Support.

#### Importing Existing Keys into Terraform State
- A key may be imported with its ID using the syntax described in the [Import](#import) section below. However, the actual value of the key _cannot be imported_ if the key being fetched was created by a user other than the one whose API key is being used to run Terraform. In such cases, the API returns a truncated key for security reasons. For more details, see [Use NerdGraph to manage license keys and User API keys](https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-user-keys/#query-keys).

#### Account Type Restrictions for Ingest Keys
- Creating `INGEST` keys requires a New Relic user with core or full platform access. See [user types](https://docs.newrelic.com/docs/accounts/accounts-billing/new-relic-one-user-management/user-type/#api-access).

## Import

Existing API access keys can be imported using a composite ID of `<api_access_key_id>:<key_type>`, where `<key_type>` is either `INGEST` or `USER`. Refer to the considerations listed in the [Important Considerations](#importing-existing-keys-into-terraform-state) section above regarding limitations on importing the actual key value.

For example:
```
$ terraform import newrelic_api_access_key.foobar "131313133A331313130B5F13DF01313FDB13B13133EE5E133D13EAAB3A3C13D3:INGEST"
```

For customers using Terraform v1.5 and above, it is recommended to use the `import {}` block in your Terraform configuration. This allows Terraform to [generate the resource configuration automatically](https://developer.hashicorp.com/terraform/language/import/generating-configuration#workflow) during the import process by running a `terraform plan -generate-config-out=<filename>.tf`, reducing manual effort and ensuring accuracy.

For example:
```hcl
import {
  id = "131313133A331313130B5F13DF01313FDB13B13133EE5E133D13EAAB3A3C13D3:INGESTT"
  to = newrelic_api_access_key.foobar
}
```

This approach simplifies the import process and ensures that the resource configuration aligns with the imported state.