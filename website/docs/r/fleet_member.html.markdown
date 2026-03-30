---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet_member"
sidebar_current: "docs-newrelic-resource-fleet-member"
description: |-
  Manage members (entities) in a New Relic fleet ring.
---

# Resource: newrelic\_fleet\_member

Use this resource to add and manage entities in a New Relic fleet ring for controlled deployment rollouts.

## Example Usage

```hcl
resource "newrelic_fleet_member" "canary_ring" {
  fleet_id = newrelic_fleet.linux_hosts.id
  ring     = "canary"

  entity_ids = [
    "entity-abc-123",
    "entity-def-456",
    "entity-ghi-789"
  ]
}

resource "newrelic_fleet_member" "production_ring" {
  fleet_id = newrelic_fleet.linux_hosts.id
  ring     = "production"

  entity_ids = [
    "entity-prod-001",
    "entity-prod-002"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `fleet_id` - (Required) The ID of the fleet. **Note**: This cannot be changed after creation.
* `ring` - (Required) The ring name within the fleet. **Note**: This cannot be changed after creation.
* `entity_ids` - (Required) A set of entity IDs to add to the fleet ring. Entities can be added or removed by updating this set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The composite ID in format `fleet_id:ring`.

## Import

Fleet members can be imported using the format `fleet_id:ring`:

```
$ terraform import newrelic_fleet_member.canary_ring fleet-abc-123:canary
```

**Note**: When importing, you'll need to manually specify the `entity_ids` in your configuration to match the current state.
