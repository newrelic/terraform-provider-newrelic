---
layout: "newrelic"
page_title: "New Relic: newrelic_change_tracking_event"
sidebar_current: "docs-newrelic-resource-change-tracking-event"
description: |-
  Create a change tracking event in New Relic.
---

# Resource: newrelic\_change\_tracking\_event

Use this resource to create change tracking events in New Relic. Details regarding change tracking and supported event types can be found [here](https://docs.newrelic.com/docs/change-tracking/change-tracking-introduction/).

~> **NOTE:** This resource is write-only. It creates a change tracking event on apply but does not support read or update operations. Destroying the resource only removes it from Terraform state.

## Example Usage

##### Basic Deployment Event
```hcl
resource "newrelic_change_tracking_event" "example" {
  entity_guid       = "MXxBUE18QVBQTElDQVRJT058MTIzNDU2Nzg"
  version           = "1.0.0"
  category_type     = "DEPLOYMENT__BASIC"
  deployment_type   = "BASIC"
  changelog         = "https://github.com/example/repo/blob/main/CHANGELOG.md"
  commit            = "abc123def456"
  deep_link         = "https://ci.example.com/builds/42"
  description       = "Deploying version 1.0.0 to production"
  short_description = "v1.0.0 production deploy"
  group_id          = "deploy-group-001"
  user              = "deploy-bot"
  timestamp         = 1698000000000
}
```

##### Feature Flag Event
```hcl
resource "newrelic_change_tracking_event" "feature_flag" {
  entity_guid      = "MXxBUE18QVBQTElDQVRJT058MTIzNDU2Nzg"
  category_type    = "FEATURE_FLAG__BASIC"
  feature_flag_id  = "my-feature-flag"
  description      = "Enabled dark mode feature flag"
  user             = "platform-team"
}
```

##### Custom Category Event with Validation Flags
```hcl
resource "newrelic_change_tracking_event" "custom" {
  entity_guid   = "MXxBUE18QVBQTElDQVRJT058MTIzNDU2Nzg"
  category_type = "CUSTOMER_DEFINED__CUSTOM"
  custom_type   = "MY_CUSTOM_TYPE"
  description   = "A custom change tracking event"
  user          = "ops-team"

  validation_flags = ["ALLOW_CUSTOM_CATEGORY_OR_TYPE"]
}
```

##### Using Entity Search Query
```hcl
resource "newrelic_change_tracking_event" "by_search" {
  entity_search_query = "name = 'my-service' AND domain = 'APM'"
  version             = "2.3.1"
  category_type       = "DEPLOYMENT__CANARY"
  deployment_type     = "CANARY"
  description         = "Canary deployment of v2.3.1"
  user                = "release-bot"
}
```

## Argument Reference

The following arguments are supported:

* `entity_guid` - (Optional) The entity GUID to associate with the change tracking event. Exactly one of `entity_guid` or `entity_search_query` should be provided.
* `entity_search_query` - (Optional) Entity search query string that matches exactly one entity. Exactly one of `entity_guid` or `entity_search_query` should be provided.
* `version` - (Optional) The version of the deployed software, for example, `v1.1`. Required for deployment category events.
* `changelog` - (Optional) A URL to the changelog or, if not linkable, a list of changes.
* `commit` - (Optional) The commit identifier, for example, a Git commit SHA.
* `deep_link` - (Optional) A link to the system that generated the change event.
* `description` - (Optional) A description of the event.
* `short_description` - (Optional) A short description of the change suitable for showing in various parts of New Relic One, such as marker flags and activity stream panels.
* `group_id` - (Optional) An identifier used to correlate account-wide changes across entities. Changes with the same `group_id` are shown together in the `Changes in group` section of the change event details UI.
* `user` - (Optional) The name or identifier of the person responsible for the change.
* `timestamp` - (Optional) The start time of the change tracking event as the number of milliseconds since the Unix epoch. Defaults to now.
* `category` - (Optional) The category of the change event (e.g., `DEPLOYMENT`, `FEATURE_FLAG`). This is derived automatically from `category_type` when `category_type` is provided.
* `category_type` - (Optional) The combined category and type of the change event. One of: `BUSINESS_EVENT__CONVENTION`, `BUSINESS_EVENT__MARKETING_CAMPAIGN`, `BUSINESS_EVENT__OTHER`, `CUSTOMER_DEFINED__CUSTOM`, `DEPLOYMENT_LIFECYCLE__ARTIFACT_COPY`, `DEPLOYMENT_LIFECYCLE__ARTIFACT_DELETION`, `DEPLOYMENT_LIFECYCLE__ARTIFACT_DEPLOYMENT`, `DEPLOYMENT_LIFECYCLE__ARTIFACT_MOVE`, `DEPLOYMENT_LIFECYCLE__BUILD_DELETION`, `DEPLOYMENT_LIFECYCLE__BUILD_PROMOTION`, `DEPLOYMENT_LIFECYCLE__BUILD_UPLOAD`, `DEPLOYMENT_LIFECYCLE__IMAGE_DELETION`, `DEPLOYMENT_LIFECYCLE__IMAGE_PROMOTION`, `DEPLOYMENT_LIFECYCLE__IMAGE_PUSH`, `DEPLOYMENT_LIFECYCLE__RELEASE_BUNDLE_CREATION`, `DEPLOYMENT_LIFECYCLE__RELEASE_BUNDLE_DELETION`, `DEPLOYMENT_LIFECYCLE__RELEASE_BUNDLE_SIGN`, `DEPLOYMENT__BASIC`, `DEPLOYMENT__BLUE_GREEN`, `DEPLOYMENT__CANARY`, `DEPLOYMENT__OTHER`, `DEPLOYMENT__ROLLING`, `DEPLOYMENT__SHADOW`, `FEATURE_FLAG__BASIC`, `OPERATIONAL__CRASH`, `OPERATIONAL__OTHER`, `OPERATIONAL__SCHEDULED_MAINTENANCE_PERIOD`, `OPERATIONAL__SERVER_REBOOT`.
* `custom_type` - (Optional) A custom type override. Required when `category_type` is set to `CUSTOMER_DEFINED__CUSTOM`. You must also include `ALLOW_CUSTOM_CATEGORY_OR_TYPE` in `validation_flags`.
* `deployment_type` - (Optional) The type of deployment. One of: `BASIC`, `BLUE_GREEN`, `CANARY`, `OTHER`, `ROLLING`, `SHADOW`. Applicable when the category is `DEPLOYMENT`.
* `feature_flag_id` - (Optional) The identifier of the feature flag. Required when `category_type` is `FEATURE_FLAG__BASIC`.
* `validation_flags` - (Optional) A list of validation flags to control how input data is handled. Valid values are: `ALLOW_CUSTOM_CATEGORY_OR_TYPE`, `FAIL_ON_FIELD_LENGTH`, `FAIL_ON_REST_API_FAILURES`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique change tracking identifier (`changeTrackingId`) returned by the API after the event is created. If the API does not return an ID, a fallback value derived from the `timestamp` is used.