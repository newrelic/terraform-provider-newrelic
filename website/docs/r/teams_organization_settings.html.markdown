---
layout: "newrelic"
page_title: "New Relic: newrelic_teams_organization_settings"
sidebar_current: "docs-newrelic-resource-teams-organization-settings"
description: |-
  Manage New Relic Teams organization-level discovery settings.
---

# Resource: newrelic\_teams\_organization\_settings

Use this resource to manage the organization-level settings that control how New Relic Teams automatically discovers and assigns entities to teams.

-> **NOTE:** There is exactly one `TEAMS_ORGANIZATION_SETTINGS` entity per New Relic organization. This resource is **import-only** — it cannot create or delete the settings entity. Use `terraform import` to bring the existing settings into Terraform management (see [Import](#import) below).

## Example Usage

```hcl
resource "newrelic_teams_organization_settings" "org" {
  discovery_enabled  = true
  discovery_tag_keys = ["team", "teamId"]
}
```

After importing, run `terraform plan`. If the plan shows no changes, the configuration already matches the current settings.

## Argument Reference

The following arguments are supported:

  * `discovery_enabled` - (Required) Whether automatic entity discovery is enabled for the organization. When `true`, entities with tags matching a team's name or aliases are automatically assigned to that team's ownership.
  * `discovery_tag_keys` - (Required) A list of tag keys used for automatic entity discovery. Entities whose tags have any of these keys — with a value matching a team's name or one of its aliases — are automatically assigned to that team. Common values: `["team"]`, `["team", "teamId"]`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The entity GUID of the organization settings entity.

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

After importing, run `terraform plan` to confirm the configuration matches the current state.
