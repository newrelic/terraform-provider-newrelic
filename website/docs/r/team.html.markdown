---
layout: "newrelic"
page_title: "New Relic: newrelic_team"
sidebar_current: "docs-newrelic-resource-team"
description: |-
  Create and manage New Relic Teams.
---

# Resource: newrelic\_team

Use this resource to create, update, and delete [New Relic Teams](https://docs.newrelic.com/docs/service-architecture-intelligence/teams/teams-intro/).

Teams let you group people and associate owned entities — services, dashboards, scorecards, and more — under a single operational identity within your New Relic organization. Teams are visible across alerts, service maps, and the Teams UI.

## Example Usage

```hcl
resource "newrelic_team" "platform" {
  name        = "Platform Engineering"
  description = "Owns the core platform services and shared infrastructure"

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

### Core Identity

* `name` - (Required) The display name of the team. Must be unique within the organization.
* `description` - (Optional) A free-text description of the team's purpose. Can be cleared by setting to an empty string `""`.
* `aliases` - (Optional) A set of alternate names the team is known by. Aliases are used alongside the primary team name for tag-based entity discovery.
* `tags` - (Optional) One or more `tags` blocks assigning key-value metadata to this team. Tags prefixed with `nr.` are managed by New Relic and are preserved automatically during updates. Each block supports:
  * `key` - (Required) The tag key.
  * `values` - (Required) One or more tag values.

### Hierarchy

* `parent_id` - (Optional) The entity GUID of a parent team, establishing this team's position in the organization hierarchy. Must reference another `newrelic_team` resource. Cannot be set to the team's own GUID.

### Supplemental Resources (Links)

* `resources` - (Optional) A list of supplemental links associated with this team — for example, runbooks, wikis, or communication channels. Each block supports:
  * `type` - (Required) The resource category. Accepted values: `ATLASSIAN_CONFLUENCE`, `ATLASSIAN_JIRA`, `ATLASSIAN_JIRA_SCORECARDS`, `BASECAMP`, `BLAMELESS`, `EMAIL`, `FACEBOOK_WORKPLACE`, `GITHUB`, `GITLAB`, `GOOGLE_CHAT`, `GOOGLE_CLOUD_PLATFORM`, `GOOGLE_DRIVE`, `MICROSOFT_AZURE`, `MICROSOFT_SHAREPOINT`, `MICROSOFT_TEAMS`, `OPSGENIE`, `OTHER_CONTACT`, `OTHER_LINK`, `PAGERDUTY`, `ROCKET_CHAT`, `SERVICENOW`, `SKYPE`, `SLACK`, `ZENDESK`.
  * `content` - (Required) The resource URL or contact address.
  * `title` - (Optional) A human-readable label for the resource.

### Membership

* `members` - (Optional) A set of New Relic user IDs (integers) to add as team members. Every manager must also be listed as a member.
* `managers` - (Optional) A set of New Relic user IDs (integers) to designate as team managers. Each ID must also appear in `members`.

### Entity Ownership

* `entity_management_mode` - (Optional) Controls how Terraform manages the team's owned entities. Accepted values: `managed` (default), `unmanaged`. See [Entity Ownership Modes](#entity-ownership-modes) below.
* `entities` - (Optional, Computed) A set of entity GUIDs that this team statically owns. Only valid when `entity_management_mode = "managed"`. Cannot be set when mode is `"unmanaged"`.

## Entity Ownership Modes

Teams in New Relic support two approaches to entity ownership. Choose the mode that matches how your organization manages entity tagging.

### `managed` (default)

In managed mode, Terraform is the authoritative source for the team's owned entities. The `entities` attribute declares exactly which entity GUIDs the team owns, and Terraform reconciles the ownership collection with that declaration on every apply.

**When to use:** You want Terraform to be the single source of truth for entity ownership. This is the recommended approach when your team owns a known, stable set of entities (for example, a specific set of APM applications or scorecards).

**What you see:**
- Entities declared in the `entities` block are added to the collection on apply.
- Entities removed from the `entities` block are removed from the collection on apply.
- If an entity is added to the collection outside Terraform (for example, via the New Relic UI), the next plan will show that entity as a planned removal, and a warning is surfaced explaining what happened. Applying reconciles the collection back to the declared state.
- Entities auto-assigned by New Relic's tag-based discovery (see [`newrelic_teams_organization_settings`](teams_organization_settings.html)) are excluded from drift detection — they are managed by the platform, not by Terraform. A separate warning is surfaced when discovery-assigned entities are present so you are aware of them.

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

In unmanaged mode, Terraform does not manage the team's entity ownership collection at all. The collection is left entirely to New Relic's tag-based discovery and/or manual management through the Teams UI.

**When to use:** Your team's entity ownership is governed by tagging conventions (for example, every service tagged `team: platform` is automatically assigned to the Platform team). In this model, declaring entities in Terraform would conflict with the platform's dynamic assignment and create unnecessary noise in your plans.

**What you see:**
- The `entities` attribute is always empty in state and cannot be set in config.
- No drift is shown for the ownership collection, regardless of what the platform or UI adds or removes.
- No entity warnings are surfaced.

```hcl
resource "newrelic_team" "example" {
  name                   = "Payments Service"
  entity_management_mode = "unmanaged"
  # entities must not be set in unmanaged mode
}
```

### Switching Between Modes

Switching from `managed` to `unmanaged` clears the `entities` attribute from Terraform state on the next apply. The ownership collection itself is **not modified** — any entities already in the collection remain there and are now managed entirely outside Terraform.

Switching from `unmanaged` back to `managed` re-enables entity tracking. Any entities already in the collection (added out-of-band or by tag discovery) will appear as out-of-band additions on the next plan if they are not present in the `entities` block.

## Attributes Reference

In addition to all arguments above, the following computed attributes are exported:

* `id` - The entity GUID of the team.
* `organization_id` - The organization UUID, resolved automatically from the provider account.
* `membership_collection_id` - The GUID of the auto-created team membership collection. Managed by the provider; do not modify directly.
* `ownership_collection_id` - The GUID of the auto-created team ownership collection. Managed by the provider; do not modify directly.

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
  name = "Engineering Standards"
  rule_ids = [newrelic_scorecard_rule.alert_coverage.id]
}

resource "newrelic_team" "platform" {
  name                   = "Platform"
  entity_management_mode = "managed"
  entities               = [newrelic_scorecard.standards.id]
}
```

### Tag-discovered team (unmanaged entity ownership)

```hcl
# Any entity tagged with "team: payments" is automatically assigned to this
# team by New Relic's tag discovery (configured in newrelic_teams_organization_settings).
# Terraform does not track or modify the ownership collection.
resource "newrelic_team" "payments" {
  name                   = "Payments"
  aliases                = ["payments-service"]
  entity_management_mode = "unmanaged"
}
```

## Import

Teams can be imported using the entity GUID:

```bash
$ terraform import newrelic_team.example <guid>
```
