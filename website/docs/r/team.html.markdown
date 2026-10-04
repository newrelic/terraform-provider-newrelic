---
layout: "newrelic"
page_title: "New Relic: newrelic_team"
sidebar_current: "docs-newrelic-resource-team"
description: |-
  Create and manage New Relic Teams.
---

# Resource: newrelic\_team

Use this resource to create, update, and delete [New Relic Teams](https://docs.newrelic.com/docs/service-architecture-intelligence/teams/teams-intro/).

Teams group people and associate owned entities — services, dashboards, scorecards, and more — under a single operational identity. They appear across alerts, service maps, and the Teams UI.

-> **NOTE:** Entity discovery can auto-assign entities to a team when their tags match the team name or an alias. See [`newrelic_teams_organization_settings`](teams_organization_settings.html) to configure which tag keys trigger discovery. Terraform tracks and diffs only entities that are declared in the `entities` attribute or manually added outside Terraform — tag-discovered entities are excluded from drift detection and surfaced as an informational warning instead.

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
  * `aliases` - (Optional) A set of alternate names for the team. Aliases are used alongside the primary team name for tag-based entity discovery and searching.
  * `tags` - (Optional) One or more `tags` blocks assigning key-value metadata to this team. Tags prefixed with `nr.` are managed by New Relic and are preserved automatically during updates. Each block supports:
    * `key` - (Required) The tag key.
    * `values` - (Required) One or more tag values.

### Hierarchy

  * `parent_id` - (Optional) The entity GUID of a parent team. Setting this places the team within the organisational hierarchy. Must reference another `newrelic_team` resource and cannot be set to the team's own GUID.

### Supplemental Resources

  * `resources` - (Optional) A list of supplemental links associated with this team, such as runbooks, wikis, or communication channels. Each block supports:
    * `type` - (Required) The resource category. Accepted values: `ATLASSIAN_CONFLUENCE`, `ATLASSIAN_JIRA`, `ATLASSIAN_JIRA_SCORECARDS`, `BASECAMP`, `BLAMELESS`, `EMAIL`, `FACEBOOK_WORKPLACE`, `GITHUB`, `GITLAB`, `GOOGLE_CHAT`, `GOOGLE_CLOUD_PLATFORM`, `GOOGLE_DRIVE`, `MICROSOFT_AZURE`, `MICROSOFT_SHAREPOINT`, `MICROSOFT_TEAMS`, `OPSGENIE`, `OTHER_CONTACT`, `OTHER_LINK`, `PAGERDUTY`, `ROCKET_CHAT`, `SERVICENOW`, `SKYPE`, `SLACK`, `ZENDESK`.
    * `content` - (Required) The resource URL or contact address.
    * `title` - (Optional) A human-readable label for this resource.

### Membership

  * `members` - (Optional) A set of New Relic user IDs (integers) to add as team members. Every manager must also be listed here.
  * `managers` - (Optional) A set of New Relic user IDs (integers) to designate as team managers. Each ID must also appear in `members`.

### Entity Ownership

  * `entity_management_mode` - (Optional) Controls how Terraform manages the team's owned entities. Accepted values: `managed` (default) and `unmanaged`. See [Entity Ownership Modes](#entity-ownership-modes) below for a full explanation of each mode and when to use them.
  * `entities` - (Optional, Computed) A set of entity GUIDs this team statically owns. Valid only when `entity_management_mode = "managed"`. Cannot be set when mode is `"unmanaged"`.

## Entity Ownership Modes

### `managed` (default)

Terraform is the authoritative source for the team's owned entities. The `entities` attribute declares exactly which entity GUIDs the team owns, and Terraform reconciles the ownership collection on every apply.

**When to use:** You want Terraform to be the single source of truth for entity ownership — for example, when the team owns a fixed set of APM applications or scorecards.

**Behaviour:**
- Entities in the `entities` block are added to the collection on apply.
- Entities removed from the `entities` block are removed from the collection on apply.
- Entities added to the collection outside Terraform (via the UI or API) are detected as drift on the next plan. A warning is emitted explaining what happened, and the entity will be removed on apply unless it is also added to the `entities` block.
- Entities auto-assigned by tag-based discovery are excluded from drift detection. A separate informational warning is emitted for them. To take declarative control of a tag-discovered entity, add its GUID to the `entities` block — the provider handles the "already in collection" response gracefully.

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
- The `entities` attribute must not be set in config — the provider returns an error at plan time if it is.
- The `entities` attribute is always empty in state.
- No drift is shown for the ownership collection, regardless of what the platform or UI adds or removes.
- No entity warnings are emitted.

```hcl
resource "newrelic_team" "example" {
  name                   = "Payments Service"
  entity_management_mode = "unmanaged"
  # entities must not be set in unmanaged mode
}
```

### Switching Between Modes

Switching from `managed` to `unmanaged` clears the `entities` attribute from state on the next apply. The ownership collection is **not modified** — entities remain in the collection and are now governed entirely outside Terraform.

Switching from `unmanaged` back to `managed` re-enables entity tracking. Any entities that were in the collection during the unmanaged phase (including any added out-of-band) will appear as drift on the next plan if they are not declared in the `entities` block.

## Attributes Reference

In addition to all arguments above, the following computed attributes are exported:

  * `id` - The entity GUID of the team.
  * `organization_id` - The organisation UUID, resolved automatically from the provider account.
  * `membership_collection_id` - The GUID of the auto-created team membership collection. Managed by the provider; do not modify directly.
  * `ownership_collection_id` - The GUID of the auto-created team ownership collection. Managed by the provider; do not modify directly.
  * `hierarchy_level_id` - The GUID of the hierarchy level this team is assigned to. Set automatically by the platform based on the team's position in the `parent_id` chain. Matches a level GUID from the `newrelic_teams_organization_settings` `hierarchy_levels` list. See [Hierarchy Levels](#hierarchy-levels) for usage.

## Hierarchy Levels

When teams are linked via `parent_id`, the platform automatically assigns each team a `hierarchy_level_id` corresponding to its depth in the hierarchy. You can reference this attribute in `newrelic_teams_organization_settings` to rename the level that a team belongs to, without needing a separate NerdGraph query.

```hcl
data "newrelic_teams_hierarchy_levels" "all" {}

resource "newrelic_team" "division" {
  name = "Engineering Division"
}

resource "newrelic_team" "squad" {
  name      = "Backend Squad"
  parent_id = newrelic_team.division.id
}

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
# Any entity tagged with "team: payments" is auto-assigned to this team
# by New Relic's tag discovery (configured in newrelic_teams_organization_settings).
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

## Import

Teams can be imported using the entity GUID:

```bash
$ terraform import newrelic_team.example <guid>
```
