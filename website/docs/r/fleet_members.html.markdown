---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet_members"
sidebar_current: "docs-newrelic-resource-fleet-members"
description: |-
  Manages the set of entities assigned to one or more rings of a New Relic Fleet.
---

# Resource: newrelic\_fleet\_members

Use this resource to manage entity assignments across one or more rings of a New Relic Fleet. Each resource instance is scoped to a single fleet and may declare any number of `ring` blocks, one per ring to manage.

-> **Note** This resource requires access to the New Relic Fleet Control API. Ensure the API key provided to the provider has the necessary Fleet Control permissions.

## How membership management works

Entities can join a fleet ring through two independent paths:

1. **Agent Control instrumentation** — when an entity is instrumented with Agent Control and associated with a fleet, the platform assigns it to the fleet automatically. These entities are outside of Terraform's control unless explicitly declared (see [Adopting Agent Control entities](#adopting-agent-control-entities) below).

2. **Explicit assignment** — this resource manages entities via the `FleetControlAddFleetMembers` and `FleetControlRemoveFleetMembers` API mutations. Only entities listed in `entity_ids` are tracked.

This resource operates on an **opt-in management model**: it tracks only the entities explicitly declared in each `ring` block. Entities present in a ring through any other means—including Agent Control instrumentation—are not visible to Terraform, are never reported as drift, and are never removed by this resource.

## Warning reference

The following warnings may appear during `plan` or `apply` operations:

### Entities already assigned in fleet

```
Warning: Entities already assigned in fleet — skipped add for ring "default"
```

This warning fires during `create` or `update` when one or more entities listed in `entity_ids` are already assigned somewhere in the fleet (including in a different ring, or via Agent Control). The API rejects duplicate add requests, so the affected entities are skipped for the add mutation. However, if those entities are already in the target ring, they are confirmed by the subsequent read and adopted into Terraform state. From that point on, Terraform manages their lifecycle: removing them from `entity_ids` will remove them from the fleet ring on the next `apply`.

If an entity is in a *different* ring than the one declared, it must first be removed from that ring before it can be reassigned. The warning message identifies the affected entity GUIDs and indicates which ring they were targeted for.

### Fleet members removed outside Terraform

```
Warning: Fleet members removed outside Terraform in ring "default"
```

This warning fires during `plan` or `apply` when an entity listed in `entity_ids` is no longer present in the ring according to the API. This indicates out-of-band drift—the entity was removed from the ring by a process other than Terraform. The next `apply` will re-add the entity to restore the declared state. To accept the removal instead, remove the entity GUID from `entity_ids` and apply again.

A second `plan` after apply will confirm the state is clean once API-side changes have propagated.

## Example Usage

### Single ring

```hcl
resource "newrelic_fleet" "example" {
  name                = "my-fleet"
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
}

resource "newrelic_fleet_members" "example" {
  fleet_id = newrelic_fleet.example.id

  ring {
    name       = "default"
    entity_ids = [
      "MXxOR0VQfEhPU1R8MTIzNDU2Nzg",
      "MXxOR0VQfEhPU1R8ODc2NTQzMjE",
    ]
  }
}
```

### Multiple rings

Multiple `ring` blocks can be declared within a single resource instance to manage entity assignments across rings in one operation. Moving an entity from one ring to another is handled by updating the `entity_ids` lists in both ring blocks and applying in a single step.

```hcl
resource "newrelic_fleet_members" "example" {
  fleet_id = newrelic_fleet.example.id

  ring {
    name       = "default"
    entity_ids = [
      "MXxOR0VQfEhPU1R8MTIzNDU2Nzg",
    ]
  }

  ring {
    name       = "canary"
    entity_ids = [
      "MXxOR0VQfEhPU1R8ODc2NTQzMjE",
      "MXxOR0VQfEhPU1R8OTk5ODc2NTQ",
    ]
  }
}
```

### Adopting Agent Control entities

To bring an entity that joined via Agent Control under Terraform management, add its GUID to `entity_ids`. A warning will indicate that the entity is already assigned in the fleet and was skipped during the add mutation. The entity is then confirmed by the subsequent read and adopted into state. Subsequent plans will be no-ops, and destroying the resource (or removing the GUID from `entity_ids`) will remove the entity from the fleet ring.

```hcl
resource "newrelic_fleet_members" "example" {
  fleet_id = newrelic_fleet.example.id

  ring {
    name = "default"
    entity_ids = [
      # Previously managed by Agent Control — now under Terraform lifecycle management.
      "MXxOR0VQfEhPU1R8MTIzNDU2Nzg",
      # Explicitly assigned entity.
      "MXxOR0VQfEhPU1R8ODc2NTQzMjE",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `fleet_id` - (Required, Forces new resource) The GUID of the fleet to manage entity assignments for.

* `ring` - (Required) One or more ring blocks. Each block declares which entities Terraform should maintain in that ring. At least one `ring` block must be specified. The following arguments are supported within each `ring` block:
  * `name` - (Required) The name of the ring (e.g. `"default"`, `"canary"`).
  * `entity_ids` - (Required) An ordered list of entity GUIDs to assign to this ring. Only the entities listed here are managed by this resource; any other entities present in the ring through other means are not affected. Removing a GUID from this list will remove that entity from the fleet ring on the next `apply`.

## Drift Detection

On each `plan` or `refresh`, this resource queries the current membership of each declared ring and compares it against the entity GUIDs in state. If any declared entity is no longer present in the ring, a warning is emitted and the pending plan will include a change to re-add it. If the removal was intentional, remove the GUID from `entity_ids` to silence the warning and prevent re-addition.

Entities present in the ring but not declared in `entity_ids` are never surfaced as drift.

## Import

A `newrelic_fleet_members` resource can be imported using the fleet GUID:

```sh
terraform import newrelic_fleet_members.example <fleet_guid>
```

On import, the resource queries all current members of the fleet and populates a single `ring` block named `"default"` with all discovered entity GUIDs. For fleets with entities distributed across multiple rings, update the `ring` blocks in the configuration to match the actual ring topology after import and run `terraform plan` to confirm there are no unintended changes.

-> **Note on ordering** Because `entity_ids` is an ordered list, the entity GUIDs populated on import reflect the order returned by the API. If your configuration specifies a different order, `terraform plan` will show an in-place update to reorder the list. Reorder the `entity_ids` entries in your configuration to match the imported order, or apply once to let Terraform reorder the state.
