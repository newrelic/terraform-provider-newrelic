---
layout: "newrelic"
page_title: "New Relic: newrelic_change_tracking_deployment"
sidebar_current: "docs-newrelic-resource-change-tracking-deployment"
description: |-
  Create and manage a deployment marker for change tracking in New Relic.
---

# Resource: newrelic\_change\_tracking\_deployment

Use this resource to create deployment markers for change tracking in New Relic. Details regarding change tracking and supported deployment types can be found [here](https://docs.newrelic.com/docs/change-tracking/change-tracking-graphql/).

## Example Usage

```hcl
resource "newrelic_change_tracking_deployment" "example" {
  entity_guid     = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM4Mjk2"
  version         = "1.0.0"
  deployment_type = "BASIC"
  changelog       = "https://github.com/example/repo/blob/main/CHANGELOG.md"
  commit          = "abc123def456"
  deep_link       = "https://ci.example.com/builds/123"
  description     = "Deploying version 1.0.0 to production"
  group_id        = "deployment-group-1"
  user            = "deployer-bot"
}
```

## Argument Reference

The following arguments are supported:

* `entity_guid` - (Required) The GUID of the New Relic entity the deployment belongs to.
* `version` - (Required) The version of the deployed software, for example, something like v1.1.
* `deployment_type` - (Optional) The type of deployment. One of: `BASIC`, `BLUE_GREEN`, `CANARY`, `OTHER`, `ROLLING`, `SHADOW`.
* `changelog` - (Optional) A URL for the changelog or, if not linkable, a list of changes.
* `commit` - (Optional) The commit identifier, for example, a Git commit SHA.
* `deep_link` - (Optional) A URL to the system that generated the deployment.
* `description` - (Optional) A description of the deployment.
* `group_id` - (Optional) An identifier used to correlate account-wide changes across entities. These changes are shown together in the `Changes in group` section of the change event details UI.
* `timestamp` - (Optional) The start time of the deployment as the number of milliseconds since the Unix epoch. Defaults to now.
* `user` - (Optional) The username of the deployer or bot.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `deployment_id` - The unique deployment identifier assigned by New Relic.

## Additional Examples

#### Blue-Green Deployment

```hcl
resource "newrelic_change_tracking_deployment" "blue_green" {
  entity_guid     = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM4Mjk2"
  version         = "2.0.0"
  deployment_type = "BLUE_GREEN"
  commit          = "deadbeef1234"
  description     = "Blue-green deployment of version 2.0.0"
  user            = "ci-bot"
}
```

#### Canary Deployment with Timestamp

```hcl
resource "newrelic_change_tracking_deployment" "canary" {
  entity_guid     = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM4Mjk2"
  version         = "1.5.0-canary"
  deployment_type = "CANARY"
  description     = "Canary release for version 1.5.0"
  timestamp       = 1700000000000
  user            = "release-engineer"
}
```

## Import

~> **NOTE:** Change tracking deployments are immutable — they cannot be updated after creation. This resource does not support the standard Terraform import workflow, as the provider uses `schema.RemoveFromState` for Read and `schema.NoopContext` for updates; importing an existing deployment is not supported.

## Additional Information

More information about change tracking can be found in the New Relic [documentation](https://docs.newrelic.com/docs/change-tracking/change-tracking-graphql/).

More details about the change tracking API can be found [here](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-api-change-tracking/).