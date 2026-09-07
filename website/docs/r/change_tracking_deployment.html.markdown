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
  entity_guid     = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM4Nzk2"
  version         = "v1.0.0"
  deployment_type = "BASIC"
  changelog       = "https://github.com/example/repo/CHANGELOG.md"
  commit          = "abc123def456"
  deep_link       = "https://ci.example.com/builds/1234"
  description     = "Deployment of version v1.0.0"
  group_id        = "deployment-group-1"
  user            = "deployer-bot"
}
```

## Argument Reference

The following arguments are supported:

* `entity_guid` - (Required) The GUID of the New Relic entity that was deployed.
* `version` - (Required) The version of the deployed software, for example, something like v1.1.
* `account_id` - (Optional) The account ID of the entity the deployment belongs to. Defaults to the account associated with the API key used.
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

* `id` - The unique deployment identifier.
* `deployment_id` - The unique deployment identifier.

## Additional Examples

#### Basic Deployment

```hcl
resource "newrelic_change_tracking_deployment" "basic" {
  entity_guid = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM4Nzk2"
  version     = "1.0.0"
}
```

#### Canary Deployment

```hcl
resource "newrelic_change_tracking_deployment" "canary" {
  entity_guid     = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM4Nzk2"
  version         = "2.0.0"
  deployment_type = "CANARY"
  user            = "ci-pipeline"
  description     = "Canary release of version 2.0.0"
  commit          = "deadbeef"
}
```

#### Rolling Deployment with Group ID

```hcl
resource "newrelic_change_tracking_deployment" "rolling" {
  entity_guid     = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM4Nzk2"
  version         = "3.1.0"
  deployment_type = "ROLLING"
  group_id        = "release-2024-01"
  user            = "release-bot"
  changelog       = "https://github.com/example/repo/releases/tag/v3.1.0"
  deep_link       = "https://jenkins.example.com/job/deploy/42"
}
```

## Additional Information

More information about change tracking can be found in the New Relic [documentation](https://docs.newrelic.com/docs/change-tracking/change-tracking-introduction/).

More details about the change tracking API can be found [here](https://docs.newrelic.com/docs/change-tracking/change-tracking-graphql/).

~> **NOTE:** This resource supports create operations only. Deployment markers cannot be updated or deleted through the New Relic API. The resource will be removed from Terraform state on destroy without making any API calls.