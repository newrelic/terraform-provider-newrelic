---
layout: "newrelic"
page_title: "New Relic: newrelic_change_tracking_event"
sidebar_current: "docs-newrelic-resource-change-tracking-event"
description: |-
  Create and manage a change tracking event in New Relic.
---

# Resource: newrelic\_change\_tracking\_event

Use this resource to create and manage New Relic change tracking events. Details regarding change tracking and supported event types can be found [here](https://docs.newrelic.com/docs/change-tracking/change-tracking-introduction/).

## Example Usage

```hcl
resource "newrelic_change_tracking_event" "foo" {
  entity_guid     = "MXxBUE18QVBQTElDQVRJT058MQ"
  version         = "1.0.0"
  changelog       = "https://github.com/example/repo/CHANGELOG.md"
  commit          = "abc123def456"
  deep_link       = "https://ci.example.com/builds/123"
  deployment_type = "BASIC"
  description     = "A production deployment of v1.0.0"
  group_id        = "deployment-group-1"
  timestamp       = 1700000000000
  user            = "deployment-bot"
}
```

## Argument Reference

The following arguments are supported:

* `entity_guid` - (Required) The entity GUID to associate the change tracking event with.
* `version` - (Required) The version of the deployed software.
* `changelog` - (Optional) A URL to the changelog or a list of changes.
* `commit` - (Optional) The commit identifier, for example, a Git commit SHA.
* `deep_link` - (Optional) A link to the system that generated the deployment.
* `deployment_type` - (Optional) The type of deployment. One of: `BASIC`, `BLUE_GREEN`, `CANARY`, `OTHER`, `ROLLING`, `SHADOW`.
* `description` - (Optional) A description of the deployment.
* `group_id` - (Optional) An identifier used to correlate account-wide changes across entities.
* `timestamp` - (Optional) The start time of the deployment as the number of milliseconds since the Unix epoch. Defaults to the current time if not specified.
* `user` - (Optional) The username of the deployer or bot.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique change tracking identifier assigned to the created event.
* `timestamp` - The start time of the deployment as the number of milliseconds since the Unix epoch. Populated from the API response when not explicitly set.