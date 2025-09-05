---
layout: "newrelic"
page_title: "Using modules to create and retrieve New Relic API Access Keys"
sidebar_current: "docs-newrelic-guide-api-access-key-usage"
description: |-
  Practical module-based workflows that create API access keys and retrieve their values, with important limitations and examples.
---

# Guide: API Access Key usage with modules

This guide shows module-based patterns for creating and then retrieving New Relic API Access Keys (User and Ingest/License), along with important limitations enforced by the underlying APIs.

-> Important: For providing your New Relic API key securely to Terraform itself (for provider authentication), see the API key usage guide: https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/api_key_usage_guide

-> Sensitive handling: The `newrelic_api_access_key` resource returns the `key` value only at creation time and marks it as Sensitive. You must capture and store it securely if you need the value later.

## Module: Create and Retrieve Access Keys

- Overview: Creates a key with the `newrelic_api_access_key` resource and then queries NerdGraph to fetch the key value by ID.
- Outputs: `key`, `name`, `type`, `ingest_type`.

Limitations

- For `key_type = USER`, the New Relic API only allows a user to fetch their own key value. This module works for USER keys created for oneself but not for keys created for other users. If you pass a different `user_id`, the module cannot retrieve the key value.

### Example: USER key

```terraform
module "create_access_keys" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/newrelic_api_access_key_extended/"

  api_key             = var.newrelic_api_key
  newrelic_account_id = "12345678"
  name                = "Access key for DemoApp"
  key_type            = "USER"
  user_id             = 12345623445
}

output "required_attributes" {
  value = module.create_access_keys.required_attributes
}
```

-> If `key_type = "USER"` and `user_id` is omitted, the user is inferred from the API key used to authenticate the request and the key is created for that user.

### Example: INGEST-LICENSE key

```terraform
module "create_access_keys" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/newrelic_api_access_key_extended/"

  api_key             = var.newrelic_api_key
  newrelic_account_id = "12345678"
  name                = "DemoApp"
  key_type            = "INGEST"
  ingest_type         = "LICENSE"
}

output "required_attributes" {
  value = module.create_access_keys.required_attributes
}
```

### Example: INGEST-BROWSER key

```terraform
module "create_access_keys" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/newrelic_api_access_key_extended/"

  api_key             = var.newrelic_api_key
  newrelic_account_id = "12345678"
  name                = "DemoApp"
  key_type            = "INGEST"
  ingest_type         = "BROWSER"
}

output "required_attributes" {
  value = module.create_access_keys.required_attributes
}
```

## Module: Fetch Access Keys

- Overview: Fetches an existing key value by ID using a NerdGraph query.
- Outputs: `key`, `name`, `type`, `ingest_type`.

Limitations

- For `key_type = USER`, only a user can fetch their own key value. Fetching another user’s key value is not supported by the API.

### Example

```terraform
module "fetch_access_keys" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/newrelic_api_access_key_extended/"

  api_key  = var.newrelic_api_key
  key_id   = "DWEGHFF327532576931786356532327538273"
  key_type = "INGEST"
}

output "required_attributes" {
  value = module.fetch_access_keys.required_attributes
}
```

