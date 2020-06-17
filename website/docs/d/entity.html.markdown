---
layout: "newrelic"
page_title: "New Relic: newrelic_entity"
sidebar_current: "docs-newrelic-datasource-entity"
description: |-
  Looks up the information about an entity in New Relic One.
---

# Data Source: newrelic\_entity

Use this data source to get information about a specific entity in New Relic One that already exists. More information on Terraform's data sources can be found [here](https://www.terraform.io/docs/configuration/data-sources.html).

-> **IMPORTANT!** Version 2.0.0 of the New Relic Terraform Provider introduces some [additional requirements](/docs/providers/newrelic/index.html) for configuring the provider.
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

resource "newrelic_alert_condition" "foo" {
  policy_id = newrelic_alert_policy.foo.id

  name        = "foo"
  type        = "apm_app_metric"
  entities    = [data.newrelic_application.app.application_id]
  metric      = "apdex"
  runbook_url = "https://www.example.com"

  term {
    duration      = 5
    operator      = "below"
    priority      = "critical"
    threshold     = "0.75"
    time_function = "all"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the entity in New Relic One.  The first entity matching this name for the given search parameters will be returned.
* `type` - (Optional) The entity's type. Valid values are APPLICATION, DASHBOARD, HOST, MONITOR, and WORRKLOAD.
* `domain` - (Optional) The entity's domain. Valid values are APM, BROWSER, INFRA, MOBILE, and SYNTH.
* `tags` - (Optional) A tag applied to the entity.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `guid` - The unique GUID of the entity.
* `account_id` - The New Relic account ID associated with this entity.
* `application_id` - The domain-specific application ID of the entity. Only returned for APM and Browser applications.
