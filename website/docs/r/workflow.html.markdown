---
layout: "newrelic"
page_title: "New Relic: newrelic_workflow"
sidebar_current: "docs-newrelic-resource-workflow"
description: |-
Create and manage a workflow in New Relic.
---

# Resource: newrelic\_workflow

Use this resource to create and manage New Relic workflow.

## Example Usage

##### Workflow
```hcl
resource "newrelic_workflow" "foo" {
  name = "workflow-example"
  account_id = 12345678
  enrichments_enabled = true
  destinations_enabled = true
  enabled = true
  muting_rules_handling = "NOTIFY_ALL_ISSUES"

  enrichments {
    nrql {
      name = "Log"
      configuration {
        query = "SELECT count(*) FROM Log WHERE message like '%error%' since 10 minutes ago"
      }
    }

    nrql {
      name = "Metric"
      configuration {
        query = "SELECT count(*) FROM Metric WHERE metricName = 'myMetric'"
      }
    }
  }

  issues_filter {
    name = "filter-name"
    type = "FILTER"

    predicate {
      attribute = "accumulations.sources"
      operator = "EXACTLY_MATCHES"
      values = [ "newrelic", "pagerduty" ]
    }
  }

  destination {
    channel_id = "20d86999-169c-461a-9c16-3cf330f4b3aa"
  }

  destination {
    channel_id = "e6af0870-cabb-453f-bf0d-fb2b6a14e05c"
  }
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) Determines the New Relic account where the workflow will be created. Defaults to the account associated with the API key used.
* `name` - (Required) The name of the workflow.
* `enrichments_enabled` - (Optional) Whether enrichments are enabled..
* `destinations_enabled` - (Optional) Whether destinations are enabled..
* `enabled` - (Optional) Whether workflow is enabled.
* `muting_rules_handling` - (Required) Which muting rule handling this workflow has.
* `destination` - (Required) A nested block that contains a channel id.
* `issues_filter` - (Required) The issues filter.  See [Nested issues_filter blocks](#nested-issues_filter-blocks) below for details.
* `enrichments` - (Optional) A nested block that describes a workflow's enrichments. See [Nested enrichments blocks](#nested-enrichments-blocks) below for details.

### Nested `issues_filter` blocks

Each workflow type supports a set of arguments for the `issues_filter` block:

* `name` - the filter's name.
* `type` - the filter's type.   One of: `FILTER` or `VIEW`.
* `predicate`
  * `attribute` - A predicate's attribute.
  * `operator` - A predicate's operator. One of: `CONTAINS`, `DOES_NOT_CONTAIN`, `DOES_NOT_EQUAL`, `DOES_NOT_EXACTLY_MATCH`, `ENDS_WITH`, `EQUAL`, `EXACTLY_MATCHES`, `GREATER_OR_EQUAL`, `GREATER_THAN`, `IS`, `IS_NOT`, `LESS_OR_EQUAL`, `LESS_THAN` or `STARTS_WITH` (workflows).
  * `values` - A list of values.

### Nested `enrichments` blocks

Each workflow type supports a specific set of arguments for the `enrichments` block:

* `nrql`
  * `name` - A nrql enrichment name.
  * `configuration` - A list of nrql enrichments.
    * `query` - the nrql query.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the workflow.

## Full Scenario Example
Create a destination resource and reference that destination to the channel resource. Then create a workflow and reference the channel resource to it.

### Create a policy
```hcl
resource "newrelic_alert_policy" "collector-policy" {
  name = "my_policy"
}
```

### Create a destination
```hcl
resource "newrelic_notification_destination" "webhook-destination" {
  account_id = 12345678
  name = "destination-webhook"
  type = "WEBHOOK"

  property {
    key = "url"
    value = "https://webhook.mywebhook.com"
  }

  auth_basic {
    user = "username"
    password = "password"
  }
}
```

### Create a channel
```hcl
resource "newrelic_notification_channel" "webhook-channel" {
  account_id = 12345678
  name = "channel-webhook"
  type = "WEBHOOK"
  destination_id = newrelic_notification_destination.webhook-destination.id
  product = "IINT"

  property {
    key = "payload"
    value = "{name: {{ variable }} }"
    label = "Payload Template"
  }
}
```

### Create a workflow
```hcl
resource "newrelic_workflow" "workflow-example" {
  name = "workflow-example"
  account_id = 12345678
  muting_rules_handling = "NOTIFY_ALL_ISSUES"

  enrichments {
    nrql {
      name = "Log count"
      
      configuration {
       query = "SELECT count(*) FROM Log WHERE message like '%error%' since 10 minutes ago"
      }
    }
  }

  issues_filter {
    name = "Filter-name"
    type = "FILTER"

    predicate {
      attribute = "accumulations.policyName"
      operator = "EXACTLY_MATCHES"
      values = [ "my_policy" ]
    }

    predicate {
      attribute = "accumulations.sources"
      operator = "EXACTLY_MATCHES"
      values = [ "newrelic" ]
    }
  }

  destination {
    channel_id = newrelic_notification_channel.webhook-channel.id
  }
}
```

## Additional Information
More details about the workflows can be found [here](https://docs.newrelic.com/docs/alerts-applied-intelligence/applied-intelligence/incident-workflows/incident-workflows/).

## V1 changes:
In workflow resource V1 we renamed the following arguments:

- `workflow_enabled` changed to `enabled`.
- `destination_configuration` changed to `destination`.
- `predicates` changed to `predicate`.
- Enrichment's `configurations` changed to `configuration`.
