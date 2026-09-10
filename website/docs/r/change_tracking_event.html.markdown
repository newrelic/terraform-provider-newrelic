---
layout: "newrelic"
page_title: "New Relic: newrelic_change_tracking_event"
sidebar_current: "docs-newrelic-resource-change-tracking-event"
description: |-
  Create a change tracking event for an entity in New Relic.
---

# Resource: newrelic\_change\_tracking\_event

Use this resource to create change tracking events in New Relic. Details regarding change tracking and supported event types can be found [here](https://docs.newrelic.com/docs/change-tracking/change-tracking-introduction/).

## Example Usage

```hcl
resource "newrelic_change_tracking_event" "foo" {
  entity_guid     = "MXxBUE18QVBQTElDQVRJT058MTIzNDU2Nzg"
  version         = "1.0.0"
  changelog       = "https://github.com/example/repo/blob/main/CHANGELOG.md"
  commit          = "abc123def456"
  deep_link       = "https://ci.example.com/builds/42"
  deployment_type = "BASIC"
  description     = "Deployed version 1.0.0 of the example service."
  group_id        = "deployment-group-1"
  timestamp       = 1698350400000
  user            = "deployer-bot"
}
```

## Argument Reference

The following arguments are supported:

* `entity_guid` - (Required) The entity GUID to associate the change tracking event with.
* `version` - (Required) The version of the deployed software, for example, something like `v1.1`.
* `changelog` - (Optional) A URL to the changelog or, if not linkable, a list of changes.
* `commit` - (Optional) The commit identifier, for example, a Git commit SHA.
* `deep_link` - (Optional) A URL to the system that generated the deployment.
* `deployment_type` - (Optional) The type of deployment. One of: `BASIC`, `BLUE_GREEN`, `CANARY`, `OTHER`, `ROLLING`, or `SHADOW`.
* `description` - (Optional) A description of the deployment.
* `group_id` - (Optional) An identifier used to correlate account-wide changes across entities. These changes are shown together in the `Changes in group` section of the change event details UI.
* `timestamp` - (Optional) The start time of the deployment as the number of milliseconds since the Unix epoch. Defaults to now.
* `user` - (Optional) The username of the deployer or bot.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `deployment_id` - A unique deployment identifier assigned by New Relic after the change tracking event is created.