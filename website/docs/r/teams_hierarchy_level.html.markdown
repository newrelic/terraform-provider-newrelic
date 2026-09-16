---
layout: "newrelic"
page_title: "New Relic: newrelic_teams_hierarchy_level"
sidebar_current: "docs-newrelic-resource-teams-hierarchy-level"
description: |-
  Manage New Relic Teams hierarchy level names.
---

# Resource: newrelic\_teams\_hierarchy\_level

Use this resource to manage the name of a [New Relic Teams hierarchy level](https://docs.newrelic.com/docs/service-architecture-intelligence/teams/teams-intro/).

-> **NOTE:** Hierarchy level entities are created automatically by New Relic when a team's `parent_id` is set (see [`newrelic_team`](team.html)). This resource is **import-only** — it cannot create or delete hierarchy levels. Use `terraform import` to bring an existing level into Terraform management, then use this resource to rename it.

## Example Usage

```hcl
resource "newrelic_teams_hierarchy_level" "department" {
  name = "Department"
}

resource "newrelic_teams_hierarchy_level" "squad" {
  name = "Squad"
}
```

To import and manage existing hierarchy levels, first discover their GUIDs using NerdGraph:

```graphql
{
  actor {
    entityManagement {
      entitySearch(query: "type = 'TEAMS_HIERARCHY_LEVEL'") {
        entities {
          id
          name
        }
      }
    }
  }
}
```

Then import each level individually (see [Import](#import) below).

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The display name of this hierarchy level as shown in the New Relic UI (e.g. `"Department"`, `"Squad"`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The entity GUID of the hierarchy level.

## Import

Hierarchy levels must be imported before they can be managed. Use the entity GUID:

```bash
$ terraform import newrelic_teams_hierarchy_level.department <guid>
$ terraform import newrelic_teams_hierarchy_level.squad <guid>
```

After importing, run `terraform plan` to confirm the configuration matches the imported state.
