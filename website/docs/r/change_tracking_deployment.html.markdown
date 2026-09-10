---
layout: "newrelic"
page_title: "New Relic: newrelic_change_tracking_deployment"
sidebar_current: "docs-newrelic-resource-change-tracking-deployment"
description: |-
  Create and manage a change tracking deployment marker in New Relic.
---

# Resource: newrelic\_change\_tracking\_deployment

Use this resource to create and manage New Relic change tracking deployment markers. Details regarding change tracking and supported deployment types can be found [here](https://docs.newrelic.com/docs/change-tracking/change-tracking-introduction/).

## Example Usage

```hcl
resource "newrelic_change_tracking_deployment" "foo" {
  entity_guid     = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3ODc2"
  version         = "v1.2.3"
  deployment_type = "BASIC"
  changelog       = "https://github.com/example/repo/blob/main/CHANGELOG.md"
  commit          = "a1b2c3d4e5f6"
  deep_link       = "https://ci.example.com/builds/1234"
  description     = "Deployed new feature release"
  group_id        = "deployment-group-abc"
  timestamp       = 1709600000000
  user            = "deployer-bot"
}
```

## Argument Reference

The following arguments are supported:

* `version` - (Required) The version of the deployed software, for example, something like v1.1.
* `entity_guid` - (Required) The GUID of the entity that was deployed.
* `deployment_type` - (Optional) The type of deployment. One of: (`BASIC`, `BLUE_GREEN`, `CANARY`, `OTHER`, `ROLLING`, `SHADOW`).
* `changelog` - (Optional) A URL for the changelog or, if not linkable, a list of changes.
* `commit` - (Optional) The commit identifier, for example, a Git commit SHA.
* `deep_link` - (Optional) A URL to the system that generated the deployment.
* `description` - (Optional) A description of the deployment.
* `group_id` - (Optional) An identifier used to correlate account-wide changes across entities.
* `timestamp` - (Optional) The start time of the deployment as the number of milliseconds since the Unix epoch. Defaults to now.
* `user` - (Optional) The username of the deployer or bot.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `deployment_id` - A unique deployment identifier generated after the deployment is created.