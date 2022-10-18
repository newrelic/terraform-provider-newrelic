---
layout: "newrelic"
page_title: "New Relic: newrelic_workload"
sidebar_current: "docs-newrelic-resource-workload"
description: |-
  Create and manage a New Relic One workload.
---

# Resource: newrelic\_workload

Use this resource to create, update, and delete a New Relic One workload.

A New Relic User API key is required to provision this resource.  Set the `api_key`
attribute in the `provider` block or the `NEW_RELIC_API_KEY` environment
variable with your User API key.

## Example Usage

Include entities with a certain string on the name.
```hcl
resource "newrelic_workload" "foo" {
	name = "Example workload"
	account_id = 12345678

	entity_guids = ["MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"]

	entity_search_query {
		query = "name like '%Example application%'"
	}

	scope_account_ids =  [12345678]
}
```

Include entities with a set of tags.
```hcl
resource "newrelic_workload" "foo" {
	name = "Example workload with tags"
	account_id = 12345678

	entity_guids = ["MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"]

	entity_search_query {
		query = "tags.accountId = '12345678' AND tags.environment='production' AND tags.language='java'"
	}

	scope_account_ids =  [12345678]
}
```

Include entities with a set of tags.
```hcl
resource "newrelic_workload" "foo" {
	name = "Example workload with tags"
	account_id = 12345678

	entity_guids = ["MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"]

	entity_search_query {
		query = "tags.accountId = '12345678' AND tags.environment='production' AND tags.language='java'"
	}

	scope_account_ids =  [12345678]
}
```

Include automatic status

-> The global status of your workload is a quick indicator of the workload health. You can configure it to be calculated automatically, and you can also set an alert and get a notification whenever the workload stops being operational. Alternatively, you can communicate a certain status of the workload by setting up a static value and a description. [See our docs](https://docs.newrelic.com/docs/workloads/use-workloads/workloads/workload-status)


```hcl
resource "newrelic_workload" "foo" {
  name       = "Example workload"
  account_id = 12345678

  entity_guids = ["MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"]

  entity_search_query {
    query = "name like '%Example application%'"
  }

  scope_account_ids = [12345678]

  description = "Description"

  status_config_automatic {
    enabled = true
    remaining_entities_rule {
      remaining_entities_rule_rollup {
        strategy        = "BEST_STATUS_WINS"
        threshold_type  = "FIXED"
        threshold_value = 100
        group_by        = "ENTITY_TYPE"
      }
    }
    rule {
      entity_guids = ["MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"]
      nrql_query {
        query = "name like '%Example application2%'"
      }
      rollup {
        strategy        = "BEST_STATUS_WINS"
        threshold_type  = "FIXED"
        threshold_value = 100
      }
    }
  }
}
```

Include static status

-> You can use this during maintenance tasks or any other time you want to provide a fixed status for your workload. This overrides all automatic rules. [See our docs](https://docs.newrelic.com/docs/workloads/use-workloads/workloads/workload-status#configure-static)

```hcl
resource "newrelic_workload" "foo" {
  name       = "Example workload"
  account_id = 12345678

  entity_guids = ["MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"]

  entity_search_query {
    query = "name like '%Example application%'"
  }

  scope_account_ids = [12345678]

  description = "Description"

  status_config_static {
    description = "test"
    enabled     = true
    status      = "OPERATIONAL"
    summary     = "summary of the status"
  }
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The workload's name.
  * `account_id` - (Required) The New Relic account ID where you want to create the workload.
  * `entity_guids` - (Optional) A list of entity GUIDs manually assigned to this workload.
  * `entity_search_query` - (Optional) A list of search queries that define a dynamic workload.  See [Nested entity_search_query blocks](#nested-entity_search_query-blocks) below for details.
  * `scope_account_ids` - (Optional) A list of account IDs that will be used to get entities from.
  * `description` - (Optional) Relevant information about the workload.
  * `status_config_automatic` - (Optional) An input object used to represent an automatic status configuration.See [Nested status_config_automatic blocks](#nested-status_config_automatic-blocks) below for details.
  * `status_config_static` - (Optional) A list of static status configurations. You can only configure one static status for a workload.See [Nested status_config_static blocks](#nested-status_config_static-blocks) below for details.

### Nested `entity_search_query` blocks

All nested `entity_search_query` blocks support the following common arguments:

  * `query` - (Required) The query.

### Nested `status_config_automatic` blocks

  * `enabled` - (Required) Whether the automatic status configuration is enabled or not.
  * `remaining_entities_rule` - (Optional) An additional meta-rule that can consider all entities that haven't been evaluated by any other rule. See [Nested remaining_entities_rule blocks](#nested-remaining_entities_rule-blocks) below for details.
  * `rule` - (Optional) The input object used to represent a rollup strategy. See [Nested rule blocks](#nested-rule-blocks) below for details.

### Nested `status_config_static` blocks
  
  * `description` - (Optional) A description that provides additional details about the status of the workload.
  * `enabled` - (Required) Whether the static status configuration is enabled or not.
  * `status` - (Required) The status of the workload.
  * `summary` - (Optional) A short description of the status of the workload.


### Nested `remaining_entities_rule` blocks

  * `remaining_entities_rule_rollup` - (Required) The input object used to represent a rollup strategy. See [Nested remaining_entities_rule_rollup blocks](#nested-remaining_entities_rule_rollup-blocks) below for details.

### Nested `rule` blocks

All nested `rule` blocks support the following common arguments:

  * `entity_guids` - (Optional) A list of entity GUIDs composing the rule.
  * `nrql_query` - (Optional) A list of entity search queries used to retrieve the entities that compose the rule. See [Nested nrql_query blocks](#nested-nrql_query-blocks) below for details.
  * `rollup` - (Required) The input object used to represent a rollup strategy. See [Nested rollup blocks](#nested-rollup-blocks) below for details.

### Nested `remaining_entities_rule_rollup` blocks

  * `group_by` - (Required) The grouping to be applied to the remaining entities.
  * `strategy` - (Required) The rollup strategy that is applied to a group of entities.
  * `threshold_type` - (Optional) Type of threshold defined for the rule. This is an optional field that only applies when strategy is WORST_STATUS_WINS. Use a threshold to roll up the worst status only after a certain amount of entities are not operational.
  * `threshold_value` - (Optional) Threshold value defined for the rule. This optional field is used in combination with thresholdType. If the threshold type is null, the threshold value will be ignored.

### Nested `nrql_query` blocks

All nested `nrql_query` blocks support the following common arguments:

  * `query` - The entity search query that is used to perform the search of a group of entities.

### Nested `rollup` blocks

  * `strategy` - (Required) The rollup strategy that is applied to a group of entities.
  * `threshold_type` - (Optional) Type of threshold defined for the rule. This is an optional field that only applies when strategy is WORST_STATUS_WINS. Use a threshold to roll up the worst status only after a certain amount of entities are not operational.
  * `threshold_value` - (Optional) Threshold value defined for the rule. This optional field is used in combination with thresholdType. If the threshold type is null, the threshold value will be ignored.

## Attributes Reference

The following attributes are exported:

  * `guid` - The unique entity identifier of the workload in New Relic.
  * `workload_id` - The unique entity identifier of the workload.
  * `permalink` - The URL of the workload.
  * `composite_entity_search_query` - The composite query used to compose a dynamic workload.

## Import

New Relic workloads can be imported using a concatenated string of the format
 `<account_id>:<workload_id>:<guid>`, e.g.

```bash
$ terraform import newrelic_workload.foo 12345678:1456:MjUyMDUyOHxBUE18QVBRTElDQVRJT058MjE1MDM3Nzk1
```
