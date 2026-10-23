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

-> **NOTE:** Entity ownership is managed non-authoritatively. Terraform tracks only the entities you declare in the `entities` block and ensures they remain in the collection. Entities added to the collection outside Terraform — through tag-based discovery or manual additions via the UI or API — are not surfaced as drift and are not removed on apply. To disable entity management entirely and suppress all entity-related plan output, set `entity_management_mode = "unmanaged"`.

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
  * `entities` - (Optional, Computed) A set of entity GUIDs that Terraform tracks within the team's ownership collection. Valid only when `entity_management_mode = "managed"`. Cannot be set when mode is `"unmanaged"`.

-> **NOTE:** When the `entities` attribute is omitted from your configuration, Terraform carries the current state forward (`Computed`) and surfaces no drift. To explicitly clear all tracked entities, declare `entities = []`.

## Entity Ownership Modes

### `managed` (default)

Terraform manages a declared set of entities within the team's ownership collection using a non-authoritative approach. You declare which entity GUIDs Terraform should track; Terraform ensures those entities are present in the collection and surfaces a plan change only when one of them is missing.

Entities added to the collection outside Terraform — through tag-based discovery or manual additions via the UI or API — are not surfaced as drift and are not removed on apply. This means the ownership collection can contain more entities than are declared in `entities`, and Terraform will not interfere with those additional entries.

**When to use:** You want Terraform to maintain a specific set of owned entities (for example, a scorecard or a key service) while allowing the platform to populate the rest of the collection through its own discovery mechanisms.

**Behaviour:**
- Entities declared in the `entities` block are added to the collection on apply.
- Entities removed from the `entities` block are removed from the collection on apply.
- If a declared entity is removed from the collection outside Terraform, the next plan surfaces a change to re-add it.
- Entities present in the collection but not declared in `entities` are silently ignored — no drift, no warnings.

```hcl
resource "newrelic_team" "example" {
  name                   = "Payments Service"
  entity_management_mode = "managed"
  entities = [
    "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1",
    newrelic_scorecard.platform_standards.id,
  ]
}
```

### `unmanaged`

Terraform does not manage the team's entity ownership collection at all. The collection is left entirely to the New Relic platform — tag-based discovery or manual management through the Teams UI or API.

**When to use:** Entity ownership for this team is governed entirely outside Terraform (for example, by tagging conventions). You want to manage the team resource itself — name, aliases, members — without Terraform interfering with what the platform assigns to the collection.

**Behaviour:**
- The `entities` attribute cannot be set in config. The provider returns a plan-time error if it is.
- The `entities` attribute is always empty in state.
- No drift is surfaced for the ownership collection, regardless of what is added or removed externally.

```hcl
resource "newrelic_team" "example" {
  name                   = "Payments Service"
  entity_management_mode = "unmanaged"
}
```

### Switching Between Modes

Switching from `managed` to `unmanaged` clears `entities` from state on the next apply. The ownership collection itself is **not modified** — entities remain in the collection and are governed entirely outside Terraform from that point.

Switching from `unmanaged` back to `managed` allows entity tracking to resume. Entities already in the collection are unaffected; you can begin declaring any of them in the `entities` block to bring them under Terraform management.

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

### Team with externally managed entity ownership

Use `unmanaged` mode when entity ownership is managed outside Terraform — through tag-based discovery, relationship rules, or the Teams UI — and you do not want Terraform to produce any drift for the ownership collection.

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
