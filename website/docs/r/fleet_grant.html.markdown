---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet_grant"
sidebar_current: "docs-newrelic-resource-fleet-grant"
description: |-
  Create and manage fleet access grants in New Relic.
---

# Resource: newrelic\_fleet\_grant

Use this resource to grant fleet access to groups with specific roles in New Relic.

## Example Usage

### Single Grant

```hcl
data "newrelic_fleet_role" "fleet_manager" {
  name = "Fleet Manager"
  type = "STANDARD"
}

resource "newrelic_fleet_grant" "platform_team" {
  fleet_id = newrelic_fleet.production.id

  grant {
    group_id = "group-platform-123"
    role_id  = data.newrelic_fleet_role.fleet_manager.id
  }
}
```

### Multiple Grants

```hcl
data "newrelic_fleet_role" "fleet_manager" {
  name = "Fleet Manager"
  type = "STANDARD"
}

resource "newrelic_fleet_grant" "multi_team_access" {
  fleet_id = newrelic_fleet.production.id

  grant {
    group_id = "group-platform-123"
    role_id  = data.newrelic_fleet_role.fleet_manager.id
  }

  grant {
    group_id = "group-ops-456"
    role_id  = data.newrelic_fleet_role.fleet_manager.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `fleet_id` - (Required) The ID of the fleet. **Note**: This cannot be changed after creation.
* `grant` - (Required) One or more grant blocks (see below). At least one grant is required.

### Grant Block

The `grant` block supports:

* `group_id` - (Required) The group ID to grant access to.
* `role_id` - (Required) The role ID to assign to the group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A composite ID representing all grants.
* `organization_id` - The organization ID.

Each `grant` block also exports:

* `id` - The unique ID of the individual grant.

## Import

Fleet grants can be imported using the composite ID:

```
$ terraform import newrelic_fleet_grant.platform_team <composite_id>
```

**Note**: Import is complex due to the composite ID encoding. It's recommended to create new grants rather than importing existing ones.
