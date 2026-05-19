---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet_members"
sidebar_current: "docs-newrelic-datasource-fleet-members"
description: |-
  Retrieves member entities of a New Relic Fleet ring.
---

# Data Source: newrelic\_fleet\_members

Use this data source to retrieve the entities (hosts, Kubernetes clusters) that are currently members of a New Relic Fleet, optionally filtered by ring.

## Example Usage

### All members of a fleet

```hcl
data "newrelic_fleet_members" "all" {
  fleet_id = newrelic_fleet.example.id
}

output "fleet_member_ids" {
  value = [for m in data.newrelic_fleet_members.all.members : m.id]
}
```

### Members in a specific ring

```hcl
data "newrelic_fleet_members" "canary" {
  fleet_id = newrelic_fleet.example.id
  ring     = "canary"
}
```

## Argument Reference

* `fleet_id` - (Required) The GUID of the fleet to query members for.
* `ring` - (Optional) Filter members by ring name. If omitted, members across all rings are returned.

## Attributes Reference

* `members` - List of member entities. Each element contains:
  * `id` - The entity GUID of the fleet member.
  * `name` - The name of the entity.
  * `type` - The entity type (e.g. `HOST`, `KUBERNETESCLUSTER`).
