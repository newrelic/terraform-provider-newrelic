---
layout: "newrelic"
page_title: "New Relic: newrelic_team"
sidebar_current: "docs-newrelic-resource-team"
description: |-
  Create and manage New Relic Teams.
---

# Resource: newrelic\_team

Use this resource to create, update, and delete [New Relic Teams](https://docs.newrelic.com/docs/service-architecture-intelligence/teams/teams-intro/).

Teams let you group people and associate owned entities — services, dashboards, scorecards, and more — under a single operational identity. Teams integrate with the broader New Relic Entity Platform (NGEP) and appear across alerts, service maps, and the Teams UI.

-> **NOTE:** Entity ownership follows a **non-authoritative model** for dynamically discovered entities. New Relic can auto-assign entities to a team based on tag-matching rules configured in [`newrelic_teams_organization_settings`](teams_organization_settings.html). Terraform only tracks entities declared in the `entities` block and will not remove entities that NGEP assigned automatically via tag discovery.

-> **NOTE:** To assign a hierarchy level to a team, set the `parent_id` attribute. The hierarchy level entities themselves are managed via [`newrelic_teams_hierarchy_level`](teams_hierarchy_level.html) (import-only).

## Example Usage

```hcl
resource "newrelic_team" "platform" {
  name        = "Platform Engineering"
  description = "Owns the core platform services and infrastructure"
  aliases     = ["platform", "infra"]
  tags        = ["team:platform", "tier:1"]

  entities {
    guid = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"
  }
  entities {
    guid = newrelic_scorecard.engineering.id
  }
}
```

See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The display name of the team. Must be unique within the organization.
  * `description` - (Optional) A description of the team's role or responsibilities. Can be cleared by setting to an empty string `""`.
  * `aliases` - (Optional) A list of alternate names the team is known by. Used for tag-based entity discovery alongside the primary team name.
  * `tags` - (Optional) A list of tags in `"key:value"` format to assign to the team entity. Tags managed by New Relic (prefixed with `nr.`) are preserved automatically and must not be included here.
  * `parent_id` - (Optional) The entity GUID of a parent team. Setting this establishes the team's position within a hierarchy. The parent must be another `newrelic_team` resource.
  * `managers` - (Optional) A list of New Relic user account IDs to designate as team managers.
  * `entities` - (Optional) One or more nested blocks each specifying an entity GUID that this team statically owns. See [Nested `entities` blocks](#nested-entities-blocks) below. Entities auto-assigned by tag discovery rules are ignored by Terraform.

### Nested `entities` blocks

Each `entities` block supports the following argument:

  * `guid` - (Required) The entity GUID to add to the team's static ownership collection. Can be any NGEP entity type including APM applications, hosts, dashboards, scorecards, and other teams.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The entity GUID of the team.

## Additional Examples

### Minimal team

```hcl
resource "newrelic_team" "infra" {
  name = "Infrastructure"
}
```

### Team with a parent (hierarchy)

```hcl
resource "newrelic_team" "parent" {
  name = "Engineering"
}

resource "newrelic_team" "child" {
  name      = "Backend"
  parent_id = newrelic_team.parent.id
}
```

### Team owning a scorecard

```hcl
resource "newrelic_scorecard" "standards" {
  name     = "Engineering Standards"
  rule_ids = [newrelic_scorecard_rule.alert_coverage.id]
}

resource "newrelic_team" "platform" {
  name = "Platform"

  entities {
    guid = newrelic_scorecard.standards.id
  }
}
```

## Import

Teams can be imported using the entity GUID:

```bash
$ terraform import newrelic_team.example <guid>
```
