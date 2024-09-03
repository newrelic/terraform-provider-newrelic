# Module: Create Access Keys and Fetch Access keys:

## Overview
This module may be used to create a user or ingest key using the `newrelic_api_access_key` resource, and fetch the created key, by performing a NerdGraph query under the hood, using the ID of the key created via the resource to fetch the created key.

### Outputs
The following output values are provided by the module:

* `key`: The actual API key.
* `name`: The name of the key.
* `type`: The type of API key.
* `ingest_type`: The type of ingest (applicable only for key_type = INGEST).


### Example usage #1 (USER)
```terraform
module "create_access_keys" {
  source = "../examples/modules/newrelic_api_access_key_extended"

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
  source = "../examples/modules/newrelic_api_access_key_extended"

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
  source = "../examples/modules/newrelic_api_access_key_extended"

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

## Overview
This module may be used to fetch a user or ingest key, using the ID of the key. Note that the ID of a key can be copied from the New Relic One UI, and is also exported by the newrelic_api_access_key resource in the New Relic Terraform Provider, if the key is created using this resource.

### Outputs
The following output values are provided by the module:

* `key`: The actual API key
* `name`: The name of the key.
* `type`: The type of API key
* `ingest_type`: The type of ingest (applicable only for key_type = INGEST).


### Example usage
```terraform
module "fetch_access_keys" {
  source = "../examples/modules/newrelic_api_access_key_extended"

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
