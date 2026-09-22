---
layout: "newrelic"
page_title: "New Relic: newrelic_change_tracking_event"
sidebar_current: "docs-newrelic-resource-change-tracking-event"
description: |-
  Create a change tracking event in New Relic.
---

# Resource: newrelic\_change\_tracking\_event

Use this resource to create change tracking events in New Relic. Details regarding change tracking and supported event types can be found [here](https://docs.newrelic.com/docs/change-tracking/change-tracking-graphql/).

~> **NOTE:** This resource is write-only. Once created, the event cannot be read back or updated via Terraform. Destroying the resource only removes it from Terraform state.

## Example Usage

##### Deployment Event
```hcl
resource "newrelic_change_tracking_event" "deployment" {
  entity_search {
    query = "name = 'my-app' AND domain = 'APM'"
  }

  category_and_type_data {
    kind {
      category = "DEPLOYMENT"
      type     = "BASIC"
    }

    category_fields {
      deployment {
        version   = "v1.2.3"
        changelog = "https://github.com/my-org/my-app/blob/main/CHANGELOG.md"
        commit    = "abc123def456"
        deep_link = "https://ci.example.com/builds/42"
      }
    }
  }

  description       = "Deployed version v1.2.3 to production"
  short_description = "v1.2.3 deployment"
  user              = "deploy-bot"
  group_id          = "release-2024-01"
  timestamp         = 1700000000000
}
```

##### Feature Flag Event
```hcl
resource "newrelic_change_tracking_event" "feature_flag" {
  entity_search {
    query = "name = 'my-app' AND domain = 'APM'"
  }

  category_and_type_data {
    kind {
      category = "FEATURE_FLAG"
      type     = "BASIC"
    }

    category_fields {
      feature_flag {
        feature_flag_id = "my-new-feature"
      }
    }
  }

  description       = "Enabled my-new-feature flag"
  short_description = "Feature flag toggle"
  user              = "ops-team"
}
```

##### Generic Event with Custom Category
```hcl
resource "newrelic_change_tracking_event" "generic" {
  entity_search {
    query = "name = 'my-app' AND domain = 'APM'"
  }

  category_and_type_data {
    kind {
      category = "OPERATIONAL"
      type     = "SERVER_REBOOT"
    }
  }

  description       = "Scheduled server reboot for maintenance"
  short_description = "Server reboot"
  user              = "ops-team"
  validation_flags  = ["ALLOW_CUSTOM_CATEGORY_OR_TYPE"]
}
```

## Argument Reference

The following arguments are supported:

* `entity_search` - (Optional) A nested block used to specify the entity to associate with the change tracking event via query. See [Nested entity_search blocks](#nested-entity_search-blocks) below for details.
* `category_and_type_data` - (Optional) A nested block that defines the category and type of the change event, as well as category-specific fields. See [Nested category_and_type_data blocks](#nested-category_and_type_data-blocks) below for details.
* `description` - (Optional) A description of the event.
* `group_id` - (Optional) An identifier used to correlate account-wide changes across entities. Changes with the same `group_id` are shown together in the `Changes in group` section of the change event details UI.
* `short_description` - (Optional) A short description of the change suitable for showing in various parts of New Relic One, such as in marker flags and activity stream panels.
* `timestamp` - (Optional) The start time of the change tracking event as the number of milliseconds since the Unix epoch. Defaults to now.
* `user` - (Optional) The name or identifier of the person responsible for the change.
* `validation_flags` - (Optional) A list of flags for validation. Valid values are `ALLOW_CUSTOM_CATEGORY_OR_TYPE`, `FAIL_ON_FIELD_LENGTH`, and `FAIL_ON_REST_API_FAILURES`.

### Nested `entity_search` blocks

* `query` - (Required) Entity search query string that matches exactly one entity. Supports operators `=`, `AND`, `IN`, and `LIKE`. You can filter by default properties (e.g. `id`, `accountId`, `name`) and tags.

### Nested `category_and_type_data` blocks

* `kind` - (Optional) A nested block that specifies the category and type of the change event. See [Nested kind blocks](#nested-kind-blocks) below for details.
* `category_fields` - (Optional) A nested block that is a container for category-specific input fields. See [Nested category_fields blocks](#nested-category_fields-blocks) below for details.

### Nested `kind` blocks

* `category` - (Required) The category of the change event. For a list of supported categories, [view our docs](https://docs.newrelic.com/docs/change-tracking/change-tracking-events/#supported-types).
* `type` - (Required) The type of the change event. For a list of supported types, [view our docs](https://docs.newrelic.com/docs/change-tracking/change-tracking-events/#supported-types).

### Nested `category_fields` blocks

* `deployment` - (Optional) A nested block for deployment-related fields. Required when the category is `DEPLOYMENT`. See [Nested deployment blocks](#nested-deployment-blocks) below for details.
* `feature_flag` - (Optional) A nested block for feature flag-related fields. Required when the category is `FEATURE_FLAG`. See [Nested feature_flag blocks](#nested-feature_flag-blocks) below for details.

### Nested `deployment` blocks

* `version` - (Required) The version of the deployed software, for example `v1.1`.
* `changelog` - (Optional) A URL to the changelog or, if not linkable, a list of changes.
* `commit` - (Optional) The commit identifier, for example a Git commit SHA.
* `deep_link` - (Optional) A link to the system that generated the deployment.

### Nested `feature_flag` blocks

* `feature_flag_id` - (Required) The identifier of the feature flag.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `change_tracking_id` - A unique change tracking identifier assigned to the created event.

## Additional Information

More information about change tracking can be found in the New Relic [documentation](https://docs.newrelic.com/docs/change-tracking/change-tracking-graphql/).

More details about supported category and type combinations can be found [here](https://docs.newrelic.com/docs/change-tracking/change-tracking-events/#supported-types).