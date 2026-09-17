---
layout: "newrelic"
page_title: "New Relic: newrelic_change_tracking_event"
sidebar_current: "docs-newrelic-resource-change-tracking-event"
description: |-
  Create and manage a change tracking event in New Relic.
---

# Resource: newrelic\_change\_tracking\_event

Use this resource to create and manage New Relic change tracking events. Details regarding supported products and permissions can be found [here](https://docs.newrelic.com/docs/change-tracking/change-tracking-graphql/).

## Example Usage

##### [Deployment Event](https://docs.newrelic.com/docs/change-tracking/change-tracking-events/#supported-types)
```hcl
resource "newrelic_change_tracking_event" "foo" {
  description       = "Example deployment change tracking event"
  short_description = "Deployment v1.2.3"
  group_id          = "my-group-id"
  user              = "deployer-bot"
  timestamp         = 1698000000000

  entity_search {
    query = "name = 'MyApp' AND domain = 'APM'"
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

## Argument Reference

The following arguments are supported:

* `entity_search` - (Optional) A nested block that specifies the entity to associate with the change tracking event via query. See [Nested entity_search blocks](#nested-entity_search-blocks) below for details.
* `category_and_type_data` - (Optional) A nested block that defines the category and type of the change event as well as category-specific fields. See [Nested category_and_type_data blocks](#nested-category_and_type_data-blocks) below for details.
* `description` - (Optional) A description of the event.
* `group_id` - (Optional) An identifier used to correlate account-wide changes across entities. These changes are shown together in the `Changes in group` section of the change event details UI.
* `short_description` - (Optional) A short description of the change suitable for showing in various parts of New Relic One, such as in the marker flag and activity stream panels.
* `timestamp` - (Optional) The start time of the change tracking event as the number of milliseconds since the Unix epoch. Defaults to now.
* `user` - (Optional) The name or identifier of the person responsible for the change.
* `data_handling_rules` - (Optional) A nested block that describes validation and data handling rules to be applied to the event input data. See [Nested data_handling_rules blocks](#nested-data_handling_rules-blocks) below for details.

### Nested `entity_search` blocks

* `query` - (Required) Entity search query string that matches exactly one entity. Operators available: `=`, `AND`, `IN`, `LIKE`.

### Nested `category_and_type_data` blocks

* `kind` - (Optional) A nested block that describes the category and type of the change event. See [Nested kind blocks](#nested-kind-blocks) below for details.
* `category_fields` - (Optional) A nested block that is a container for the various category related input fields. See [Nested category_fields blocks](#nested-category_fields-blocks) below for details.

### Nested `kind` blocks

* `category` - (Required) The category of the change event. For a list of supported categories, [view our docs](https://docs.newrelic.com/docs/change-tracking/change-tracking-events/#supported-types).
* `type` - (Required) The type of the change event. For a list of supported types, [view our docs](https://docs.newrelic.com/docs/change-tracking/change-tracking-events/#supported-types).

### Nested `category_fields` blocks

* `deployment` - (Optional) A nested block for deployment-related fields. This container is mandatory for a category of `DEPLOYMENT`. See [Nested deployment blocks](#nested-deployment-blocks) below for details.
* `feature_flag` - (Optional) A nested block for feature flag-related fields. This container is mandatory for a category of `FEATURE_FLAG`. See [Nested feature_flag blocks](#nested-feature_flag-blocks) below for details.

### Nested `deployment` blocks

* `version` - (Required) The version of the deployed software, for example, something like v1.1.
* `changelog` - (Optional) A URL to the changelog or, if not linkable, a list of changes.
* `commit` - (Optional) The commit identifier, for example, a Git commit SHA.
* `deep_link` - (Optional) A link to the system that generated the deployment.

### Nested `feature_flag` blocks

* `feature_flag_id` - (Required) The identifier of the feature flag.

### Nested `data_handling_rules` blocks

* `validation_flags` - (Optional) A list of flags for validation. Valid values are `ALLOW_CUSTOM_CATEGORY_OR_TYPE`, `FAIL_ON_FIELD_LENGTH`, and `FAIL_ON_REST_API_FAILURES`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `change_tracking_id` - A unique change tracking identifier assigned to the created event.