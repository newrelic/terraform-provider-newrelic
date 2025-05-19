---
layout: "newrelic"
page_title: "New Relic: newrelic_entity_relationship"
sidebar_current: "docs-newrelic-resource-entity-relationship"
description: |-
  Create and manage relationships between New Relic entities.
---

# Resource: newrelic\_entity\_relationship

Use this resource to create, update, and delete relationships between New Relic entities.

## Example Usage

```hcl
resource "newrelic_entity_relationship" "example" {
  source_entity_guid = "MzgwNjUyNnxFWFR8U0VSVklDRV9MRVZFTHw1ODA4MDM"
  target_entity_guid = "MzgwNjUyNnxFWFR8U0VSVklDRV9MRVZFTHw1NzE0Nzk"
  relation_type      = "CONTAINS"
}
```
## Argument Reference

The following arguments are supported:

* `source_entity_guid` - (Required) The GUID of the source entity in the relationship.
* `target_entity_guid` - (Required) The GUID of the target entity in the relationship.
* `relation_type` - (Required) The type of relationship to create between the source and target entities. Valid values are: `BUILT_FROM`, `BYPASS_CALLS`, `CALLS`, `CONNECTS_TO`, `CONSUMES`, `CONTAINS`, `HOSTS`, `IS`, `MANAGES`, `MEASURES`, `MONITORS`, `OPERATES_IN`, `OWNS`, `PRODUCES`, `SERVES`, `TRIGGERS`.

## Import

New Relic entity relationships can be imported using a concatenated string of the format <source_entity_guid>:<target_entity_guid>, e.g.

```bash
$ terraform import newrelic_entity_relationship.example MSHJHBKLJDKD:NSBVGHJLHKB
```