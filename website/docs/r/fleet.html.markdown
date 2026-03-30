---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet"
sidebar_current: "docs-newrelic-resource-fleet"
description: |-
  Create and manage New Relic fleets for centralized agent management.
---

# Resource: newrelic\_fleet

Use this resource to create and manage New Relic fleets for centralized agent management.

Fleets enable you to organize and manage New Relic agents across your infrastructure. You can create fleets for different types of infrastructure (HOST or KUBERNETESCLUSTER) and use them to centrally deploy agent configurations, manage updates, and organize your entities with tags.

## Example Usage

### Linux Host Fleet

```hcl
resource "newrelic_fleet" "linux_hosts" {
  name                = "Production Linux Hosts"
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
  description         = "Fleet for managing Linux production hosts"
  product             = "INFRA"

  tags = [
    "environment:production",
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
  description         = "Fleet for managing Kubernetes clusters"

  tags = [
    "environment:production"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the fleet. This can be changed after creation.
* `managed_entity_type` - (Required) The type of entities this fleet will manage. Valid values are `HOST` or `KUBERNETESCLUSTER`. **Note**: This cannot be changed after creation (forces new resource).
* `operating_system` - (Optional) The operating system type for HOST fleets. **Required when `managed_entity_type` is `HOST`**. Valid values are `LINUX` or `WINDOWS`. **Must not be set when `managed_entity_type` is `KUBERNETESCLUSTER`**. **Note**: This cannot be changed after creation (forces new resource).
* `description` - (Optional) A description of the fleet. This can be updated after creation.
* `product` - (Optional) The New Relic product associated with this fleet (e.g., `INFRA`).
* `tags` - (Optional) A list of tags for the fleet. Each tag should be in the format `"key:value1,value2"` where multiple values can be comma-separated.
* `organization_id` - (Optional) The organization ID. If not provided, it will be automatically fetched from your account. **Note**: This cannot be changed after creation (forces new resource).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the fleet.

## Import

Fleets can be imported using the fleet ID:

```
$ terraform import newrelic_fleet.linux_hosts <fleet_id>
```
