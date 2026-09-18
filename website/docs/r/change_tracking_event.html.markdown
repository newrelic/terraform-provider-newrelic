---
layout: "newrelic"
page_title: "New Relic: newrelic_change_tracking_event"
sidebar_current: "docs-newrelic-resource-change-tracking-event"
description: |-
  Create a change tracking event in New Relic.
---

# Resource: newrelic\_change\_tracking\_event

Use this resource to create change tracking events in New Relic. Details regarding change tracking and supported event types can be found [here](https://docs.newrelic.com/docs/change-tracking/change-tracking-introduction/).

~> **NOTE:** This resource is write-only. It creates a change tracking event on apply and removes it from state on destroy (the underlying event is not deleted). There is no read operation.

## Example Usage

##### Deployment Event
```hcl
resource "newrelic_change_tracking_event" "example" {
  description       = "Example deployment event"
  short_description = "Deployed v1.2.3"
  group_id          = "my-group-id"
  timestamp         = 1698000000000
  user              = "deployer-bot"

  entity_search {
    query = "name = 'My Application' AND domain = 'APM'"
  }

  category_and_type_data {
    kind {
      category = "DEPLOYMENT"
      type     = "BASIC"
    }

    category_fields {
      deployment {
        version   = "v1.2.3"
        changelog = "https://github.com/example/repo/blob/main/CHANGELOG.md"
        commit    = "abc123def456"
        deep_link = "https://ci.example.com/builds/42"
      }
    }
  }

  data_handling_rules {
    validation_flags = ["FAIL_ON_FIELD_LENGTH"]
  }
}
```

##### Feature Flag Event
```hcl
resource "newrelic_change_tracking_event" "feature_flag" {
  description       = "Enabled dark mode feature flag"
  short_description = "Feature flag toggled"
  user              = "ops-team"

  entity_search {
    query = "name = 'My Application' AND domain = 'APM'"
  }

  category_and_type_data {
    kind {
      category = "FEATURE_FLAG"
      type     = "BASIC"
    }

    category_fields {
      feature_flag {
        feature_flag_id = "dark-mode-flag"
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `entity_search` - (Optional) A nested block that specifies the entity to associate with the change tracking event via query. See [Nested entity_search blocks](#nested-entity_search-blocks) below for details.
* `category_and_type_data` - (Optional) A nested block that defines the category and type of the change event, as well as category-specific fields. See [Nested category_and_type_data blocks](#nested-category_and_type_data-blocks) below for details.
* `description` - (Optional) A description of the event.
* `group_id` - (Optional) An identifier used to correlate account-wide changes across entities. These changes are shown together in the `Changes in group` section of the change event details UI.
* `short_description` - (Optional) A concise description of the change suitable for showing in various parts of New Relic One, such as in marker flags and activity stream panels.
* `timestamp` - (Optional) The start time of the change tracking event as the number of milliseconds since the Unix epoch. Defaults to now.
* `user` - (Optional) The name or identifier of the person responsible for the change.
* `data_handling_rules` - (Optional) A nested block that describes validation and data handling rules to be applied to the event input data. See [Nested data_handling_rules blocks](#nested-data_handling_rules-blocks) below for details.

### Nested `entity_search` blocks

* `query` - (Required) Entity search query string that matches exactly one entity. Supports operators `=`, `AND`, `IN`, `LIKE`. You can filter by default properties (e.g. `id`, `accountId`, `name`) and tags.

### Nested `category_and_type_data` blocks

* `kind` - (Optional) A nested block that specifies the category and type of the change event. See [Nested kind blocks](#nested-kind-blocks) below for details.
* `category_fields` - (Optional) A nested block that is a container for the various category-related input fields. See [Nested category_fields blocks](#nested-category_fields-blocks) below for details.

### Nested `kind` blocks

* `category` - (Required) The category of the change event. For a list of supported categories, [view our docs](https://docs.newrelic.com/docs/change-tracking/change-tracking-events/#supported-types).
* `type` - (Required) The type of the change event. For a list of supported types, [view our docs](https://docs.newrelic.com/docs/change-tracking/change-tracking-events/#supported-types).

### Nested `category_fields` blocks

* `deployment` - (Optional) A nested block for deployment-related fields. This container is used when the category is `DEPLOYMENT`. See [Nested deployment blocks](#nested-deployment-blocks) below for details.
* `feature_flag` - (Optional) A nested block for feature flag-related fields. This container is used when the category is `FEATURE_FLAG`. See [Nested feature_flag blocks](#nested-feature_flag-blocks) below for details.

### Nested `deployment` blocks

* `version` - (Required) The version of the deployed software, for example, something like `v1.1`.
* `changelog` - (Optional) A URL to the changelog or, if not linkable, a list of changes.
* `commit` - (Optional) The commit identifier, for example, a Git commit SHA.
* `deep_link` - (Optional) A link to the system that generated the deployment.

### Nested `feature_flag` blocks

* `feature_flag_id` - (Required) The identifier of the feature flag.

### Nested `data_handling_rules` blocks

* `validation_flags` - (Optional) A list of flags for validation. Valid values are:
  * `ALLOW_CUSTOM_CATEGORY_OR_TYPE` - Allows passing custom categories or types.
  * `FAIL_ON_FIELD_LENGTH` - Will validate all string fields to be within max size limit. An error is returned and data is not saved if any of the fields exceeds the max size limit.
  * `FAIL_ON_REST_API_FAILURES` - For APM entities, a call is made to the legacy New Relic v2 REST API. When this flag is set, if the call fails for any reason, an error will be returned containing the failure message.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `change_tracking_id` - A unique change tracking identifier assigned to the created event.