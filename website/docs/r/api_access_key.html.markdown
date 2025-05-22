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

-> **Action Required for Key Retrieval:**
    To retrieve the actual API key value after creation, **it is essential to use this resource in conjunction with the `Create Access Keys` module** detailed in the "Extended Usage with Modules" section below. Please see the "Important Considerations" section for a full explanation.

## Important Considerations for Using This Resource

Before you begin, please take note of the following critical points:

1.  **Retrieving the API Key:**

    -> The newrelic_api_access_key resource will create an API key in New Relic, but it does not directly output the sensitive key string for use in your Terraform configuration (e.g., via the key attribute). This is to help prevent accidental exposure of sensitive credentials. To programmatically create an API key and retrieve its value for use in subsequent configurations or outputs, you must use the Create Access Keys module. This module is detailed in the "Extended Usage with Modules" section. It works by first creating the key with this resource and then immediately fetching the key's value using a secure API call.
2.  **Updating Existing Keys:**

    -> **IMPORTANT!** Exercise extreme caution when updating existing `newrelic_api_access_key` resources. Only the `name` and `notes` attributes are updatable in place. Modifying any other attribute will force the resource to be recreated, which **invalidates the previous API key(s)** and generates new ones.
3.  **Account Type Restrictions for Ingest Keys:**

    -> **WARNING:** Creating 'Ingest - License' and 'Ingest - Browser' keys using this resource is restricted to 'core' or 'full platform' New Relic user accounts. If you've signed up as a 'basic' user with New Relic, or have been added as a 'basic' user to your organization on New Relic, you would not be able to use your account to create 'Ingest' keys. If you see the message `"You do not have permission to create this key"` in the response of the API called by this resource, it could be owing to the aforementioned. For more insights into user account types on New Relic and associated privileges, please check out this [page](https://docs.newrelic.com/docs/accounts/accounts-billing/new-relic-one-user-management/user-type/#api-access).
    
## Example Usage (Resource Only - Key Not Retrievable Directly)

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
- `key` - The actual API key. This attribute is masked and not be visible in your terminal, CI, etc.Use the `Create Access Keys module` to retrieve this value.

## Import

Existing API access keys can be imported using a composite ID of `<api_access_key_id>:<key_type>`. `<key_type>`
will be either `INGEST` or `USER`.

For example:
```
$ terraform import newrelic_api_access_key.foobar "1234567:INGEST"
```
## Extended Usage with Modules

As highlighted in the "Important Considerations" section, to effectively create and retrieve API keys, or to fetch existing keys, the following modules are provided. They utilize NerdGraph queries behind the scenes to access the key values.

### Module: Create and Retrieve Access Keys

### Overview
This module allows you to create a new User or Ingest API key using the `newrelic_api_access_key` resource and and fetch the created key, by performing a NerdGraph query under the hood, using the ID of the key created via the resource to fetch the created key. This is the recommended method for creating new API keys with Terraform when you need to use the key value programmatically.

### Outputs
The following output values are provided by the module:

* `key`: The actual API key.
* `name`: The name of the key.
* `type`: The type of API key.
* `ingest_type`: The type of ingest (applicable only for key_type = INGEST).


### Example usage #1 (USER)
```terraform
module "create_access_keys" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/newrelic_api_access_key_extended/"

  create_access_keys_service = {
    api_key             = "NRAK-XXXXXXXXXX"
    newrelic_account_id = "12345678"
    name                = "Access key for DemoApp"
    key_type            = "USER"
    user_id             = 12345623445
  }
}

output "required_attributes" {
  value = module.create_access_keys.required_attributes
}
```
### Example usage #2 (INGEST-LICENSE)
```terraform
module "create_access_keys" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/newrelic_api_access_key_extended/"

  create_access_keys_service = {
    api_key             = "NRAK-XXXXXXXXXX"
    newrelic_account_id = "12345678"
    name                = "DemoApp"
    key_type            = "USER"
    ingest_type         = "LICENSE"
  }
}

output "required_attributes" {
  value = module.create_access_keys.required_attributes
}
```
### Example usage #3 (INGEST-BROWSER)
```terraform
module "create_access_keys" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/newrelic_api_access_key_extended/"

  create_access_keys_service = {
    api_key             = "NRAK-XXXXXXXXXX"
    newrelic_account_id = "12345678"
    name                = "DemoApp"
    key_type            = "USER"
    ingest_type         = "BROWSER"
  }
}

output "required_attributes" {
  value = module.create_access_keys.required_attributes
}
```
## Module: Fetch Access keys
### Overview
This module may be used to fetch a user or ingest key, using the ID of the key. Note that the ID of a key can be copied from the New Relic One UI, and is also exported by the `newrelic_api_access_key` resource in the New Relic Terraform Provider, if the key is created using this resource.

### Outputs
The following output values are provided by the module:

* `key`: The actual API key
* `name`: The name of the key.
* `type`: The type of API key
* `ingest_type`: The type of ingest (applicable only for key_type = INGEST).

### Example usage

```terraform
module "fetch_access_keys" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/newrelic_api_access_key_extended/"

  fetch_access_keys_service = {
        api_key = "NRAK-XXXXXXXXXXXXXXXX"
        key_id = "DWEGHFF327532576931786356532327538273"
        key_type = "INGEST"
  }
}

output "required_attributes" {
  value = module.fetch_access_keys.required_attributes
}
```