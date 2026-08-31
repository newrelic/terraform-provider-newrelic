---
layout: "newrelic"
page_title: "New Relic: newrelic_change_tracking_deployment"
sidebar_current: "docs-newrelic-resource-change-tracking-deployment"
description: |-
  Create and manage a deployment marker for change tracking in New Relic.
---

# Resource: newrelic\_change\_tracking\_deployment

Use this resource to create deployment markers for change tracking in New Relic. Details regarding change tracking and deployment markers can be found [here](https://docs.newrelic.com/docs/change-tracking/change-tracking-introduction/).

## Example Usage

```hcl
resource "newrelic_change_tracking_deployment" "foo" {
  entity_guid     = "MXxBUE18QVBQTElDQVRJT058MTIzNDU2Nzg"
  version         = "1.0.0"
  deployment_type = "BASIC"
  changelog       = "Added new feature X"
  commit          = "abc123def456"
  deep_link       = "https://github.com/example/repo/commit/abc123def456"
  description     = "Production deployment of version 1.0.0"
  group_id        = "deployment-group-1"
  user            = "deployment-bot"
}
```

## Argument Reference

The following arguments are supported:

* `entity_guid` - (Required) The GUID of the entity the deployment is associated with.
* `version` - (Required) The version of the deployed software, for example, something like v1.1.
* `changelog` - (Optional) A URL to the changelog or, if not linkable, a list of changes.
* `commit` - (Optional) The commit identifier, for example, a Git commit SHA.
* `deep_link` - (Optional) A URL to the system that generated the deployment.
* `deployment_type` - (Optional) The type of deployment. One of: `BASIC`, `BLUE_GREEN`, `CANARY`, `OTHER`, `ROLLING`, `SHADOW`.
* `description` - (Optional) A description of the deployment.
* `group_id` - (Optional) An identifier used to correlate account-wide changes across entities. These changes are shown together in the `Changes in group` section of the change event details UI.
* `timestamp` - (Optional) The start time of the deployment as the number of milliseconds since the Unix epoch. Defaults to now.
* `user` - (Optional) The username of the deployer or bot.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `deployment_id` - The unique identifier of the deployment.

## Import

Change tracking deployments can be imported using the `deployment_id`:

```
$ terraform import newrelic_change_tracking_deployment.foo <deployment_id>
```

~> **NOTE:** Change tracking deployments are immutable — all fields require recreation (ForceNew). Importing a deployment will allow you to manage its lifecycle in Terraform, but no updates are supported. Use `ignore_changes` if needed:

```hcl
resource "newrelic_change_tracking_deployment" "foo" {
  lifecycle {
    ignore_changes = all
  }
  entity_guid = "MXxBUE18QVBQTElDQVRJT058MTIzNDU2Nzg"
  version     = "1.0.0"
}
```

## Additional Information

More information about change tracking can be found in the New Relic [documentation](https://docs.newrelic.com/docs/change-tracking/change-tracking-introduction/).
More details about the change tracking API can be found [here](https://docs.newrelic.com/docs/change-tracking/change-tracking-graphql/).