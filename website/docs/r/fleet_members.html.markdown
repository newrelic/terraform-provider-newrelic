---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet_members"
sidebar_current: "docs-newrelic-resource-fleet-members"
description: |-
  Manages the set of member entities in a ring of a New Relic Fleet.
---

# Resource: newrelic\_fleet\_members

Use this resource to manage which entities (hosts, Kubernetes clusters) are members of a specific ring within a New Relic Fleet.

Each resource instance is authoritative over exactly one `(fleet, ring)` pair. Entities in rings not declared by any `newrelic_fleet_members` resource are left untouched.

-> **Note** This resource requires access to the New Relic Fleet Control API. Ensure the API key used has the necessary permissions.

## Example Usage

### Basic membership

```hcl
resource "newrelic_fleet" "example" {
  name                = "my-fleet"
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
}

resource "newrelic_fleet_members" "default" {
  fleet_id   = newrelic_fleet.example.id
  ring       = "default"
  entity_ids = [
    "MXxOR0VQfEhPU1R8MTIzNDU2Nzg",
    "MXxOR0VQfEhPU1R8ODc2NTQzMjE",
  ]
}
```

### Multiple rings (canary + production rollout)

```hcl
resource "newrelic_fleet_members" "canary" {
  fleet_id   = newrelic_fleet.example.id
  ring       = "canary"
  entity_ids = ["MXxOR0VQfEhPU1R8MTIzNDU2Nzg"]
}

resource "newrelic_fleet_members" "production" {
  fleet_id   = newrelic_fleet.example.id
  ring       = "production"
  entity_ids = [
    "MXxOR0VQfEhPU1R8ODc2NTQzMjE",
    "MXxOR0VQfEhPU1R8OTk5ODc2NTQ",
  ]
}
```

## Argument Reference

* `fleet_id` - (Required, Forces new resource) The GUID of the fleet to manage members for.
* `ring` - (Required, Forces new resource) The ring name within the fleet (e.g. `"default"`, `"canary"`, `"production"`).
* `entity_ids` - (Required) Set of entity GUIDs to add as members of the fleet ring. Changes to this set are applied as incremental add/remove operations against the API.

## Drift Detection

This resource is **authoritative** over the ring it manages. If entities are added to the ring outside Terraform, they will appear as a diff on the next `terraform plan` and will be removed on `terraform apply`.

To adopt existing ring membership into Terraform state, use `terraform import` before your first apply.

## Import

Fleet ring membership can be imported using `fleet_id:ring`:

```sh
terraform import newrelic_fleet_members.default <fleet_guid>:default
```

After import, Terraform will read the current ring membership from the API and populate `entity_ids` in state. Subsequent plans will be no-ops as long as the declared `entity_ids` match what is in the ring.
