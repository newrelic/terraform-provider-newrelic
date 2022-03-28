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

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The workload's name.
  * `account_id` - (Required) The New Relic account ID where you want to create the workload.
  * `entity_guids` - (Optional) A list of entity GUIDs manually assigned to this workload.
  * `entity_search_query` - (Optional) A list of search queries that define a dynamic workload.  See [Nested entity_search_query blocks](#nested-entity_search_query-blocks) below for details.
  * `scope_account_ids` - (Optional) A list of account IDs that will be used to get entities from.

### Nested `entity_search_query` blocks

All nested `entity_search_query` blocks support the following common arguments:

  * `query` - (Required) The query.

## Attributes Reference

The following attributes are exported:

  * `guid` - The unique entity identifier of the workload in New Relic.
  * `workload_id` - The unique entity identifier of the workload.
  * `permalink` - The URL of the workload.
  * `composite_entity_search_query` - The composite query used to compose a dynamic workload.

## Import

New Relic One workloads can be imported using a concatenated string of the format
 `<account_id>:<workload_id>:<guid>`, e.g.

```bash
$ terraform import newrelic_workload.foo 12345678:1456:MjUyMDUyOHxBUE18QVBRTElDQVRJT058MjE1MDM3Nzk1
```
