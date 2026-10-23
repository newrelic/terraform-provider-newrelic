---
layout: "newrelic"
page_title: "New Relic: newrelic_teams_hierarchy_levels"
sidebar_current: "docs-newrelic-datasource-teams-hierarchy-levels"
description: |-
  Read all Teams hierarchy level entities in the organisation.
---

# Data Source: newrelic\_teams\_hierarchy\_levels

Use this data source to retrieve the GUIDs and display names of all Teams hierarchy level entities in your New Relic organisation.

Hierarchy levels define the organisational tiers shown in the Teams UI (for example, "Division" and "Squad"). Their GUIDs are required when configuring the ordered hierarchy in [`newrelic_teams_organization_settings`](../r/teams_organization_settings.html).

## Example Usage

```hcl
data "newrelic_teams_hierarchy_levels" "all" {}

resource "newrelic_teams_organization_settings" "org" {
  hierarchy_levels {
    id   = data.newrelic_teams_hierarchy_levels.all.levels[0].id
    name = "Division"
  }
  hierarchy_levels {
    id   = data.newrelic_teams_hierarchy_levels.all.levels[1].id
    name = "Squad"
  }
}
```

## Argument Reference

This data source has no configurable arguments.

## Attributes Reference

The following attributes are exported:

  * `levels` - A list of hierarchy level objects. Each object contains:
    * `id` - The entity GUID of the hierarchy level. Pass this to `newrelic_teams_organization_settings.hierarchy_levels[*].id` to reference and rename the level.
    * `name` - The current display name of this level (e.g. `"Division"`, `"Squad"`).
