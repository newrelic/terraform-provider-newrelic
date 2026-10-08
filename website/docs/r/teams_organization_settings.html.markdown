---
layout: "newrelic"
page_title: "New Relic: newrelic_teams_organization_settings"
sidebar_current: "docs-newrelic-resource-teams-organization-settings"
description: |-
  Manage New Relic Teams organisation-level settings.
---

# Resource: newrelic\_teams\_organization\_settings

Use this resource to manage the organisation-level settings that control how New Relic Teams works across your organisation:

- **Entity discovery**: Automatically assigns entities to teams based on tag values.
- **Hierarchy levels**: Define and rename the ordered tiers shown in the Teams UI.
- **Sync group rules**: Automatically creates teams from IdP groups that match name patterns.

-> **NOTE:** There is exactly one `TEAMS_ORGANIZATION_SETTINGS` entity per New Relic organisation. This resource is **import-optional** — if you apply it without importing first, the provider automatically locates the singleton, applies your configuration, and emits a warning stating that the existing settings have been overridden. `terraform import` is supported for an explicit workflow.

## Example Usage

```hcl
data "newrelic_teams_hierarchy_levels" "all" {}

resource "newrelic_teams_organization_settings" "org" {
  # Tag-based entity discovery
  discovery_enabled  = true
  discovery_tag_keys = ["team", "teamId"]

  # Ordered hierarchy levels — the order here controls the Teams UI display
  hierarchy_levels {
    id   = data.newrelic_teams_hierarchy_levels.all.levels[0].id
    name = "Division"
  }
  hierarchy_levels {
    id   = data.newrelic_teams_hierarchy_levels.all.levels[1].id
    name = "Squad"
  }

  # IdP group sync
  sync_groups_enabled = true

  sync_group_rules {
    conditions {
      type  = "STARTS_WITH"
      value = "team-"
    }
  }

  sync_group_rules {
    conditions {
      type  = "CONTAINS"
      value = "-engineering-"
    }
    conditions {
      type  = "ENDS_WITH"
      value = "-prod"
    }
  }
}
```

See additional [examples](#additional-examples) below.

## Argument Reference

### Discovery

  * `discovery_enabled` - (Optional, Computed) Whether tag-based entity discovery is enabled. When `true`, entities whose tags match a team's name or an alias are automatically assigned to that team.
  * `discovery_tag_keys` - (Optional, Computed) Tag keys used for automatic entity discovery (e.g. `["team"]`). An entity is auto-assigned to a team when one of these tag keys has a value equal to the team's name or an alias. Must contain at least one key when declared — the API rejects an empty list.

### Hierarchy Levels

  * `hierarchy_levels` - (Optional, Computed) Ordered list of organisation hierarchy levels. The order controls the visual display in the Teams UI. Use the [`newrelic_teams_hierarchy_levels`](../d/teams_hierarchy_levels.html) data source to obtain level GUIDs. Each block supports:
    * `id` - (Required) The entity GUID of the hierarchy level. Obtain from the `newrelic_teams_hierarchy_levels` data source or from the `hierarchy_level_id` attribute on a `newrelic_team` resource.
    * `name` - (Required) The display name of this level (e.g. `"Division"`, `"Squad"`). Changing this renames the level entity directly in the API.

### Sync Groups

  * `sync_groups_enabled` - (Optional, Computed) Whether automatic team creation from IdP groups is enabled.
  * `sync_group_rules` - (Optional, Computed) At most one rule controlling which IdP groups automatically create teams. Omitting this block preserves existing rules in state. The NGEP API currently enforces a single rule per organisation. Each rule has:
    * `conditions` - (Required) One or more conditions that a group name must satisfy. All conditions within a rule must match (AND logic). Each condition has:
      * `type` - (Required) The match type: `STARTS_WITH`, `ENDS_WITH`, or `CONTAINS`.
      * `value` - (Required) The string to match against the IdP group name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The entity GUID of the organisation settings singleton entity.

## Hierarchy Level IDs

Use the [`newrelic_teams_hierarchy_levels`](../d/teams_hierarchy_levels.html) data source to discover level GUIDs programmatically:

```hcl
data "newrelic_teams_hierarchy_levels" "all" {}

# Reference a level by index
output "division_level_id" {
  value = data.newrelic_teams_hierarchy_levels.all.levels[0].id
}
```

Alternatively, each `newrelic_team` resource exports a `hierarchy_level_id` attribute that reflects the level the team is currently assigned to. You can use this to reference a specific level without the data source:

```hcl
resource "newrelic_team" "division" {
  name = "Engineering Division"
}

resource "newrelic_teams_organization_settings" "org" {
  hierarchy_levels {
    id   = newrelic_team.division.hierarchy_level_id
    name = "Division"
  }
}
```

## Additional Examples

### Discovery only — no hierarchy configuration

```hcl
resource "newrelic_teams_organization_settings" "org" {
  discovery_enabled  = true
  discovery_tag_keys = ["team"]
}
```

### Disable all discovery

```hcl
resource "newrelic_teams_organization_settings" "org" {
  discovery_enabled = false
}
```

## Import

The organisation settings singleton must be imported using its entity GUID. Discover the GUID with NerdGraph:

```graphql
{
  actor {
    entityManagement {
      entitySearch(query: "type = 'TEAMS_ORGANIZATION_SETTINGS'") {
        entities { id name }
      }
    }
  }
}
```

Then import:

```bash
$ terraform import newrelic_teams_organization_settings.org <guid>
```
