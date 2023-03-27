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

// Filter by account ID.
// The `accountId` tag is automatically added to all entities by the platform.
data "newrelic_entity" "app" {
  name = "my-app"
  domain = "APM"
  type = "APPLICATION"
  tag {
    key = "accountID"
    value = "12345"
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

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the entity in New Relic One.  The first entity matching this name for the given search parameters will be returned.
* `ignore_case` - (Optional) Ignore case of the `name` when searching for the entity. Defaults to false.
* `type` - (Optional) The entity's type. Valid values are APPLICATION, DASHBOARD, HOST, MONITOR, WORKLOAD, AWSLAMBDAFUNCTION, SERVICE_LEVEL, and KEY_TRANSACTION. Note: Other entity types may also be queryable as the list of entity types may fluctuate over time.
* `domain` - (Optional) The entity's domain. Valid values are APM, BROWSER, INFRA, MOBILE, SYNTH, and EXT. If not specified, all domains are searched.
* `tag` - (Optional) A tag applied to the entity. See [Nested tag blocks](#nested-`tag`-blocks) below for details.

### Nested `tag` blocks

All nested `tag` blocks support the following common arguments:

  * `key` - (Required) The tag key.
  * `value` - (Required) The tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `guid` - The unique GUID of the entity.
* `account_id` - The New Relic account ID associated with this entity.
* `application_id` - The domain-specific application ID of the entity. Only returned for APM and Browser applications.
* `serving_apm_application_id` - The browser-specific ID of the backing APM entity. Only returned for Browser applications.


## Additional Examples

-> If the entities are not found please try again without providing the `types` field.

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
