---
layout: "newrelic"
page_title: "New Relic: newrelic_teams_organization_settings"
sidebar_current: "docs-newrelic-resource-teams-organization-settings"
description: |-
  Manage New Relic Teams organization-level settings: entity discovery, hierarchy level ordering, and IdP group sync rules.
---

# Resource: newrelic\_teams\_organization\_settings

Use this resource to manage the organization-level settings that control how New Relic Teams works across your organization, including:

- **Entity discovery** — automatically assign entities to teams based on tag values
- **Hierarchy levels** — define the ordered list of hierarchy levels shown in the Teams UI, and rename individual levels
- **Sync group rules** — automatically create teams from IdP groups that match name patterns

-> **NOTE:** There is exactly one `TEAMS_ORGANIZATION_SETTINGS` entity per New Relic organization. This resource is **import-only** — it cannot create or delete the settings entity. Use `terraform import` to bring the existing settings into Terraform management (see [Import](#import) below).

## Example Usage

```hcl
# First, discover the hierarchy level IDs
data "newrelic_teams_hierarchy_levels" "levels" {}

resource "newrelic_teams_organization_settings" "org" {
  # Tag-based entity discovery
  discovery_enabled  = true
  discovery_tag_keys = ["team", "teamId"]

  # Ordered hierarchy levels (order defines the Teams UI display order)
  hierarchy_levels {
    id   = data.newrelic_teams_hierarchy_levels.levels.levels[0].id
    name = "Division"
  }
  hierarchy_levels {
    id   = data.newrelic_teams_hierarchy_levels.levels.levels[1].id
    name = "Squad"
  }

  # IdP group sync rules
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

## Argument Reference

The following arguments are supported:

### Discovery

* `discovery_enabled` - (Optional, Computed) Whether tag-based entity discovery is enabled. When `true`, entities with tags matching a team's name or aliases are automatically assigned to that team.
* `discovery_tag_keys` - (Optional, Computed) A list of tag keys used for automatic entity discovery (e.g. `["team"]`). Entities whose tags match any of these keys — with a value equal to a team's name or alias — are auto-assigned to that team.

### Hierarchy Levels

* `hierarchy_levels` - (Optional, Computed) Ordered list of organization hierarchy levels. The order here defines the visual order in the Teams UI. Each entry has:
  * `id` - (Required) The NGEP GUID of the hierarchy level entity. Use the `newrelic_teams_hierarchy_levels` data source to obtain these IDs.
  * `name` - (Required) The display name of this hierarchy level (e.g. `"Division"`, `"Squad"`). Changing this renames the level entity directly via the API.

### Sync Groups

* `sync_groups_enabled` - (Optional, Computed) Whether automatic team creation from IdP groups is enabled.
* `sync_group_rules` - (Optional) Rules that control which IdP groups automatically create teams. Each rule has:
  * `conditions` - (Required) One or more conditions that a group name must satisfy. Each condition has:
    * `type` - (Required) The match type: `STARTS_WITH`, `ENDS_WITH`, or `CONTAINS`.
    * `value` - (Required) The string to match against the IdP group name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The entity GUID of the organization settings singleton entity.

## Import

The organization settings entity must be imported using its entity GUID. Discover the GUID using NerdGraph:

```graphql
{
  actor {
    entityManagement {
      entitySearch(query: "type = 'TEAMS_ORGANIZATION_SETTINGS'") {
        entities {
          id
          name
        }
      }
    }
  }
}
```

Then import:

```bash
$ terraform import newrelic_teams_organization_settings.org <guid>
```

After importing, run `terraform plan` to confirm the configuration matches the current state. The singleton entity ID for most organizations is:

```
NDgyOTY3M3xOR0VQfFRFQU1TX09SR0FOSVpBVElPTl9TRVRUSU5HU3wwMTllMjZlZi05YmYyLTc0MGUtYTQ3OS0wNmVlOTNmYjBjN2E
```
