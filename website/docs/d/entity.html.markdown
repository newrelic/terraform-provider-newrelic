---
layout: "newrelic"
page_title: "New Relic: newrelic_entity"
sidebar_current: "docs-newrelic-datasource-entity"
description: |-
  Looks up the information about an entity in New Relic One.
---

# Data Source: newrelic\_entity

Use this data source to get information about a specific entity in New Relic One that already exists. More information on Terraform's data sources can be found [here](https://www.terraform.io/language/data-sources).

-> **IMPORTANT!** Version 2.0.0 of the New Relic Terraform Provider introduces some [additional requirements](/providers/newrelic/newrelic/latest/docs/guides/migration_guide_v2) for configuring the provider.
<br><br>
Before upgrading to version 2.0.0 or later, it is recommended to upgrade to the most recent 1.x version of the provider and ensure that your environment successfully runs `terraform plan` without unexpected changes.

## Example Usage

```hcl
data "newrelic_entity" "app" {
  name = "my-app"
  domain = "APM"
  type = "APPLICATION"
  tag {
    key = "my-tag"
    value = "my-tag-value"
  }
}

resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_nrql_alert_condition" "foo" {
  policy_id                    = newrelic_alert_policy.foo.id
  type                         = "static"
  name                         = "foo"
  description                  = "Alert when transactions are taking too long"
  runbook_url                  = "https://www.example.com"
  enabled                      = true
  violation_time_limit_seconds = 3600

  nrql {
    query             = "SELECT average(duration) FROM Transaction where appName = '${data.newrelic_entity.app.name}'"
  }

  critical {
    operator              = "above"
    threshold             = 5.5
    threshold_duration    = 300
    threshold_occurrences = "ALL"
  }
}

// Ignore name case
data "newrelic_entity" "app" {
  name = "mY-aPP"
  ignore_case = true
  domain = "APM"
  type = "APPLICATION"
}
```

### Example: Filter By Account ID

The default behaviour of this data source is to retrieve entities matching the [specified parameters](#argument-reference) (such as `name`, `domain`, `type`) from NerdGraph with the credentials specified in the configuration of the provider (account ID and API Key), filter them by the account ID specified in the configuration of the provider, and return the first match. 

This would mean, if no entity with the specified search parameters is found associated with the account ID in the configuration of the provider, i.e. `NEW_RELIC_ACCOUNT_ID`, an error is thrown, stating that no matching entity has been found.

```hcl
# The entity returned by this configuration would have to 
# belong to the account_id specified in the provider 
# configuration, i.e. NEW_RELIC_ACCOUNT_ID.
data "newrelic_entity" "app" {
  name   = "my-app"
  domain = "APM"
  type   = "APPLICATION"
}
```
However, in order to cater to scenarios in which it could be necessary to retrieve an entity belonging to a subaccount using the account ID and API Key of the parent account (for instance, when entities with identical names are present in both the parent account and subaccounts, since matching entities from subaccounts too are returned by NerdGraph), the `account_id` attribute of this data source may be availed. This ensures that the account ID in the configuration of the provider, used to filter entities returned by the API is now overridden by the `account_id` specified in the configuration; i.e., in the below example, the data source would now return an entity matching the specified `name`, belonging to the account with the ID `account_id`.
```hcl
# The entity returned by this configuration, unlike in 
# the above example, would have to belong to the account_id 
# specified in the configuration below, i.e. 654321.
data "newrelic_entity" "app" {
  name       = "my-app"
  account_id = 654321
  domain     = "APM"
  type       = "APPLICATION"
}
```
The following example explains a use case along the lines of the aforementioned; using the `account_id` argument in the data source to allow the filtering criteria to be the `account_id` specified (of the subaccount), and not the account ID in the provider configuration. 

In simpler terms, when entities are queried from the parent account, entities with matching names are returned from subaccounts too, hence, specifying the `account_id` of the subaccount in the configuration allows the entity returned to belong to the subaccount with `account_id`.
```hcl
# The `account_id` specified in the configuration of the
# provider is that of the parent account.
provider "newrelic" {
  account_id = "12345"
  ..
}

# A subaccount is created using the `newrelic_account_management` 
# resource.
resource "newrelic_account_management" "default" { 
  name   = "Sample Subaccount"
  region = "us01"
}

# The ID of the subaccount is specified in the configuration
# to allow the entity returned to belong to the subaccount.
data "newrelic_entity" "app" {
  account_id = newrelic_account_management.default.id 
  name       = "my-app"
  domain     = "APM"
  type       = "APPLICATION"
}
```

The `accountId` tag may also be added to the configuration of this data source as specified below. 

-> **NOTE:** Not to be confused with the `account_id` argument of this data source that helps filter entities retrieved from the API by the specified `account_id` and return a matching entity, adding the `accountId` tag adds the specified account to the NRQL Query that is sent to NerdGraph, i.e. it causes entities not matching `accountId` to be filtered out of the API response that is received by this data source. The entity that is finally returned by this data source, however, is the one that has an account ID matching the account ID specified in the provider configuration, or the `account_id` attribute, as specified in the examples above.

```hcl
# The `accountId` tag is automatically added to all entities by the platform.
data "newrelic_entity" "app" {
  name = "my-app"
  domain = "APM"
  type = "APPLICATION"
  tag {
    key = "accountID"
    value = "345211"
  }
}
```




## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the entity in New Relic One.  The first entity matching this name for the given search parameters will be returned.
* `account_id` - (Optional) The New Relic account ID the entity to be returned would be associated with, i.e. if specified, the data source would filter matching entities received by `account_id` and return the first match. If not, matching entities are filtered by the account ID specified in the configuration of the provider. See the **Example: Filter By Account ID** section above for more details.
* `ignore_case` - (Optional) Ignore case of the `name` when searching for the entity. Defaults to false.
* `type` - (Optional) The entity's type. Valid values are APPLICATION, DASHBOARD, HOST, MONITOR, WORKLOAD, AWSLAMBDAFUNCTION, SERVICE_LEVEL, and KEY_TRANSACTION. Note: Other entity types may also be queryable as the list of entity types may fluctuate over time.
* `domain` - (Optional) The entity's domain. Valid values are APM, BROWSER, INFRA, MOBILE, SYNTH, and EXT. If not specified, all domains are searched.
* `tag` - (Optional) A tag applied to the entity. See [Nested tag blocks](#nested-`tag`-blocks) below for details.
* `ignore_not_found`- (Optional) A boolean argument that, when set to true, prevents an error from being thrown when the queried entity is not found. Instead, a warning is displayed. Defaults to `false`.

-> **WARNING:** Setting the `ignore_not_found` argument to `true` will display an 'entity not found' warning instead of throwing an error. This can lead to downstream errors if the values of attributes exported by this data source are used elsewhere, as all of these values would be null. Please use this argument at your own risk. 

### Nested `tag` blocks

All nested `tag` blocks support the following common arguments:

  * `key` - (Required) The tag key.
  * `value` - (Required) The tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `guid` - The unique GUID of the entity.
* `application_id` - The domain-specific application ID of the entity. Only returned for APM and Browser applications.
* `serving_apm_application_id` - The browser-specific ID of the backing APM entity. Only returned for Browser applications.
* `entity_tags` - The `entity_tags` helps retrieve tags associated with the entity fetched by the data source, the tags are returned as a JSON-encoded string and not a conventional list or a map, owing to a couple of design considerations; which is why one would need to use the Terraform function `jsondecode()`, along with the attribute entity_tags in order to convert the JSON-encoded string into a map with key-value pairs.


## Additional Examples

-> If the entities are not found please try again without providing the `type` field.

### Entity Tags

* The following is an illustration of the aforementioned scenario. It may be observed that a key-value pair version of the JSON-encoded string exported by `entity_tags` is written to the variable `key_value_maps` , using the `jsondecode()` function.

```hcl
data "newrelic_entity" "foo" {
  name = "Sample Searchable Entity"
  domain = "EXT"
  type = "SERVICE_LEVEL"
}

locals {
  key_value_map = { for pair in jsondecode(data.newrelic_entity.foo.entity_tags) : pair.key => pair.values }
}

output "key_value_map" {
  value = local.key_value_map
}
```
The value of `local.key_value_map` would look like the following.

```hcl
{
  env = ["production"]
  team = ["ops"]
}
```

In order to prevent an error being thrown by the variable added if no tags are associated with the entity (consequently, the value of `entity_tags` is an empty string), the value of the variable may be conditionally defined in order to not throw an error, as shown below.

```hcl
locals {
  key_value_map = data.newrelic_entity.foo.entity_tags != null ? { for pair in jsondecode(data.newrelic_entity.foo.entity_tags) : pair.key => pair.values } : null
}
```

### Query for an OTEL entity

```hcl
data "newrelic_entity" "app" {
  name = "my-otel-app"
  domain = "EXT"
  type = "SERVICE"

  tag {
    key = "accountID"
    value = "12345"
  }
}
```

### Query for an entity by type (AWS Lambda entity in this example)

```hcl
data "newrelic_entity" "app" {
  name = "my_lambda_trace"
  type = "AWSLAMBDAFUNCTION"
}
```
