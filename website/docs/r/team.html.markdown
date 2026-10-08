---
layout: "newrelic"
page_title: "New Relic: newrelic_team"
sidebar_current: "docs-newrelic-resource-team"
description: |-
  Create and manage New Relic Teams.
---

# Resource: newrelic\_team

Use this resource to create, update, and delete [New Relic Teams](https://docs.newrelic.com/docs/service-architecture-intelligence/teams/teams-intro/).

Teams group people and associate owned entities (services, dashboards, scorecards, and more) under a single operational identity. They appear across alerts, service maps, and the Teams UI.

-> **NOTE:** Entity discovery can auto-assign entities to a team when their tags match the team name or an alias. See [`newrelic_teams_organization_settings`](teams_organization_settings.html) to configure which tag keys trigger discovery. Terraform tracks and diffs only entities that are declared in the `entities` attribute or added manually outside Terraform. Tag-discovered entities are excluded from drift detection and reported as an informational warning.

## Example Usage

```hcl
resource "newrelic_team" "platform" {
  name        = "Platform Engineering"
  description = "Owns core platform services and shared infrastructure"

  aliases = ["platform", "infra"]

  tags {
    key    = "department"
    values = ["engineering"]
  }
  tags {
    key    = "tier"
    values = ["1"]
  }

  members  = [1234567, 2345678]
  managers = [1234567]

  resources {
    type    = "GITHUB"
    content = "https://github.com/example-org/platform"
    title   = "GitHub Repository"
  }
  resources {
    type    = "SLACK"
    content = "https://example.slack.com/channels/platform-eng"
    title   = "Team Slack Channel"
  }

  entity_management_mode = "managed"
  entities = [
    "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1",
    newrelic_scorecard.engineering_standards.id,
  ]
}
```

See additional [examples](#additional-examples) below.

## Argument Reference

The following arguments are supported:

### Core Identity

  * `name` - (Required) The display name of the team. Must be unique within the organisation.
  * `description` - (Optional) A free-text description of the team's purpose. Can be cleared by setting to an empty string `""`.
  * `aliases` - (Optional) A set of alternate names for the team, used alongside the primary name for tag-based entity discovery and searching.
  * `tags` - (Optional) One or more `tags` blocks assigning key-value metadata to this team. Tags prefixed with `nr.` are managed by New Relic and preserved automatically during updates. Each block supports:
    * `key` - (Required) The tag key.
    * `values` - (Required) One or more tag values.

### Hierarchy

  * `parent_id` - (Optional) The entity GUID of a parent team. Setting this places the team within the organisational hierarchy. Must reference another `newrelic_team` resource and cannot be set to the team's own GUID.

### Supplemental Resources

  * `resources` - (Optional) A set of supplemental links for this team, such as runbooks, wikis, or communication channels. These appear in the team's Settings page in the New Relic UI under **Contacts** or **Links**, automatically categorised by the `type` field. Each block supports:
    * `type` - (Required) The resource category. Accepted values: `ATLASSIAN_CONFLUENCE`, `ATLASSIAN_JIRA`, `ATLASSIAN_JIRA_SCORECARDS`, `BASECAMP`, `BLAMELESS`, `EMAIL`, `FACEBOOK_WORKPLACE`, `GITHUB`, `GITLAB`, `GOOGLE_CHAT`, `GOOGLE_CLOUD_PLATFORM`, `GOOGLE_DRIVE`, `MICROSOFT_AZURE`, `MICROSOFT_SHAREPOINT`, `MICROSOFT_TEAMS`, `OPSGENIE`, `OTHER_CONTACT`, `OTHER_LINK`, `PAGERDUTY`, `ROCKET_CHAT`, `SERVICENOW`, `SKYPE`, `SLACK`, `ZENDESK`.
    * `content` - (Required) The resource URL or contact address.
    * `title` - (Optional) A short label shown next to the resource in the UI.

### Membership

  * `members` - (Optional) A set of New Relic user IDs (integers) to add as team members. Every manager must also be listed here.
  * `managers` - (Optional) A set of New Relic user IDs (integers) to designate as team managers. Each ID must also appear in `members`.

### Entity Ownership

  * `entity_management_mode` - (Optional) Controls how Terraform manages the team's owned entities. Accepted values: `managed` (default) and `unmanaged`. See [Entity Ownership Modes](#entity-ownership-modes) below.
  * `entities` - (Optional, Computed) A set of entity GUIDs this team statically owns. Valid only when `entity_management_mode = "managed"`. Cannot be set when mode is `"unmanaged"`.

-> **NOTE:** When the `entities` block is omitted from your configuration entirely, Terraform carries the current state forward (the attribute is `Computed`) and does not surface drift for entities added to the collection. To have Terraform detect additions and surface them as drift, declare `entities = []` explicitly. This design allows teams to coexist with tag-based discovery without generating spurious plan noise.

## Entity Ownership Modes

### `managed` (default)

Terraform is the authoritative source for the team's owned entities. The `entities` attribute declares exactly which entity GUIDs the team owns, and Terraform reconciles the ownership collection on every apply.

**When to use:** You want Terraform to be the single source of truth for entity ownership, such as when the team owns a fixed set of APM applications or scorecards.

**Behaviour:**
- Entities in the `entities` block are added to the collection on apply.
- Entities removed from the `entities` block are removed on apply.
- Entities added to the collection outside Terraform (via the UI or API) are detected as drift on the next plan. A warning explains what happened, and the entity is removed on apply unless you add it to the `entities` block.
- Entities auto-assigned by tag-based discovery are excluded from drift detection. A separate informational warning is emitted. To take declarative control of a tag-discovered entity, add its GUID to the `entities` block.

```hcl
resource "newrelic_team" "example" {
  name                   = "Payments Service"
  entity_management_mode = "managed"
  entities = [
    "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1",
  ]
}
```

### `unmanaged`

Terraform does not manage the team's entity ownership collection at all. The collection is left entirely to New Relic's tag-based discovery and/or manual management through the Teams UI.

**When to use:** Your team's entity ownership is governed by tagging conventions (for example, any service tagged `team: platform` is automatically assigned to the Platform team). Declaring entities in Terraform would conflict with the platform's dynamic assignment.

**Behaviour:**
- The `entities` attribute must not be set in config. The provider returns a plan-time error if it is.
- The `entities` attribute is always empty in state.
- No drift is shown for the ownership collection.
- No entity warnings are emitted.

```hcl
resource "newrelic_team" "example" {
  name                   = "Payments Service"
  entity_management_mode = "unmanaged"
}
```

### Switching Between Modes

Switching from `managed` to `unmanaged` clears `entities` from state on the next apply. The ownership collection itself is **not modified** — entities remain in the collection and are governed entirely outside Terraform from that point.

Switching from `unmanaged` back to `managed` re-enables entity tracking. Any entities that were added to the collection during the unmanaged phase will appear as drift on the next plan if they are not declared in the `entities` block.

## Attributes Reference

In addition to all arguments above, the following computed attributes are exported:

  * `id` - The entity GUID of the team.
  * `organization_id` - The organisation UUID, resolved automatically from the provider account.
  * `membership_collection_id` - The GUID of the auto-created team membership collection. Managed by the provider; do not modify directly.
  * `ownership_collection_id` - The GUID of the auto-created team ownership collection. Managed by the provider; do not modify directly.
  * `hierarchy_level_id` - The GUID of the hierarchy level this team is assigned to. The platform assigns this automatically based on the team's depth in the `parent_id` chain. It becomes available on the plan cycle **after** `parent_id` is first set (one refresh cycle of eventual consistency). Matches a level GUID from the `newrelic_teams_organization_settings` `hierarchy_levels` list. See [Hierarchy Levels](#hierarchy-levels).

## Hierarchy Levels

When teams are linked via `parent_id`, the platform automatically assigns each team a `hierarchy_level_id` based on its depth in the hierarchy. This value becomes available on the **next** plan after `parent_id` is first set.

You can reference `hierarchy_level_id` directly in `newrelic_teams_organization_settings` to rename the level, without requiring a separate NerdGraph query or data source lookup:

```hcl
resource "newrelic_team" "division" {
  name = "Engineering Division"
}

resource "newrelic_team" "squad" {
  name      = "Backend Squad"
  parent_id = newrelic_team.division.id
}

resource "newrelic_teams_organization_settings" "org" {
  hierarchy_levels {
    id   = newrelic_team.division.hierarchy_level_id
    name = "Division"
  }
  hierarchy_levels {
    id   = newrelic_team.squad.hierarchy_level_id
    name = "Squad"
  }
}
```

Alternatively, use the `newrelic_teams_hierarchy_levels` data source to reference levels by index:

```hcl
data "newrelic_teams_hierarchy_levels" "all" {}

resource "newrelic_teams_organization_settings" "org" {
  hierarchy_levels {
    id   = data.newrelic_teams_hierarchy_levels.all.levels[0].id
    name = "Division"
  }
}
```

## Import

Teams can be imported using the entity GUID:

```bash
$ terraform import newrelic_team.example <guid>
```

To find a team's GUID after creating it, use `terraform state show` or the following NerdGraph query:

```graphql
{
  actor {
    entityManagement {
      entitySearch(query: "type = 'TEAM' AND name = 'Your Team Name'") {
        entities { id name }
      }
    }
  }
}
```

## Additional Examples

### Minimal team

```hcl
resource "newrelic_team" "backend" {
  name = "Backend Engineering"
}
```

### Team with a parent (hierarchy)

```hcl
resource "newrelic_team" "engineering" {
  name        = "Engineering"
  description = "Top-level engineering division"
}

resource "newrelic_team" "backend" {
  name      = "Backend"
  parent_id = newrelic_team.engineering.id
}
```

### Team owning a scorecard

```hcl
resource "newrelic_scorecard" "standards" {
  name     = "Engineering Standards"
  rule_ids = [newrelic_scorecard_rule.alert_coverage.id]
}

resource "newrelic_team" "platform" {
  name                   = "Platform"
  entity_management_mode = "managed"
  entities               = [newrelic_scorecard.standards.id]
}
```

### Tag-discovered team (unmanaged entity ownership)

When entity ownership is governed by tag-based discovery, use `unmanaged` mode to prevent Terraform from interfering with the collection.

```hcl
resource "newrelic_team" "payments" {
  name                   = "Payments"
  aliases                = ["payments-service"]
  entity_management_mode = "unmanaged"
}
```

### Team with members, managers, and supplemental resources

```hcl
resource "newrelic_team" "platform" {
  name        = "Platform Engineering"
  description = "Owns shared infrastructure and developer tooling"

  members  = [1234567, 2345678, 3456789]
  managers = [1234567]

  resources {
    type    = "GITHUB"
    content = "https://github.com/example-org/platform"
    title   = "Source Repository"
  }
  resources {
    type    = "SLACK"
    content = "https://example.slack.com/channels/platform-eng"
    title   = "Team Channel"
  }
  resources {
    type    = "PAGERDUTY"
    content = "https://example.pagerduty.com/teams/platform"
    title   = "On-Call Schedule"
  }
}
```
