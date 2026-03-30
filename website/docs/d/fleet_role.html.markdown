---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet_role"
sidebar_current: "docs-newrelic-datasource-fleet-role"
description: |-
  Looks up fleet-scoped roles in New Relic.
---

# Data Source: newrelic\_fleet\_role

Use this data source to look up fleet-scoped roles (standard or custom) in New Relic.

## Example Usage

### Get Standard Fleet Manager Role

```hcl
data "newrelic_fleet_role" "fleet_manager" {
  name = "Fleet Manager"
  type = "STANDARD"
}

resource "newrelic_fleet_grant" "example" {
  fleet_id = newrelic_fleet.example.id

  grant {
    group_id = "group-123"
    role_id  = data.newrelic_fleet_role.fleet_manager.id
  }
}
```

### Get Custom Fleet Role

```hcl
data "newrelic_fleet_role" "custom_deployer" {
  name = "Custom Fleet Deployer"
  type = "CUSTOM"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the fleet-scoped role. If both `name` and `type` are omitted, defaults to the standard "Fleet Manager" role.
* `type` - (Optional) The type of the fleet-scoped role. Allowed values: `STANDARD`, `CUSTOM`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the role.
* `scope` - The scope of the role (always "fleet" for fleet roles).
