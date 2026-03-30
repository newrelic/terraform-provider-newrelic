---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet"
sidebar_current: "docs-newrelic-resource-fleet"
description: |-
  Create and manage New Relic fleets for centralized agent management.
---

# Resource: newrelic\_fleet

Use this resource to create and manage New Relic fleets for centralized agent management.

## Example Usage

### Linux Host Fleet

```hcl
resource "newrelic_fleet" "linux_hosts" {
  name                = "Production Linux Hosts"
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
  description         = "Fleet for managing Linux production hosts"
  product             = "Infrastructure"

  tags = [
    "env:production",
    "team:platform,ops"
  ]
}
```

### Windows Host Fleet

```hcl
resource "newrelic_fleet" "windows_hosts" {
  name                = "Production Windows Hosts"
  managed_entity_type = "HOST"
  operating_system    = "WINDOWS"
  description         = "Fleet for managing Windows production hosts"
}
```

### Kubernetes Cluster Fleet

```hcl
resource "newrelic_fleet" "k8s_clusters" {
  name                = "Production Kubernetes Clusters"
  managed_entity_type = "KUBERNETESCLUSTER"
  description         = "Fleet for managing K8s clusters"

  tags = [
    "env:production"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the fleet.
* `managed_entity_type` - (Required) The type of entities this fleet will manage. Allowed values: `HOST`, `KUBERNETESCLUSTER`. **Note**: This cannot be changed after creation.
* `operating_system` - (Optional) The operating system type. **Required for HOST fleets**. Allowed values: `LINUX`, `WINDOWS`. **Must not be set for KUBERNETESCLUSTER fleets**. **Note**: This cannot be changed after creation.
* `description` - (Optional) The description of the fleet.
* `product` - (Optional) The New Relic product associated with this fleet.
* `tags` - (Optional) A list of tags for the fleet in format `"key:value1,value2"`. Each tag can have multiple comma-separated values.
* `organization_id` - (Optional) The organization ID. If not provided, it will be auto-fetched from the account. **Note**: This cannot be changed after creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the fleet.

## Import

Fleets can be imported using the fleet ID:

```
$ terraform import newrelic_fleet.linux_hosts <fleet_id>
```
