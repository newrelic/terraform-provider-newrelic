---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet_configuration"
sidebar_current: "docs-newrelic-resource-fleet-configuration"
description: |-
  Create and manage New Relic fleet configurations for centralized agent management.
---

# Resource: newrelic\_fleet\_configuration

Use this resource to create and manage New Relic fleet configurations for centralized agent management.

A fleet configuration defines versioned agent settings deployable to your fleets. Each configuration is specific to an agent type and managed entity type. Versions are immutable - their content cannot be modified after creation. To add a new configuration, add a `version` block; to remove one, delete its block.

## Example Usage

### Basic Infrastructure Agent Configuration

```hcl
resource "newrelic_fleet_configuration" "infra" {
  name                = "Production Infrastructure Config"
  agent_type          = "NRInfra"
  managed_entity_type = "HOST"

  version {
    configuration_content = <<-EOT
      log:
        level: info
        file: /var/log/newrelic-infra/newrelic-infra.log
      metrics:
        enabled: true
        system_sample_rate: 15
    EOT
  }
}
```

### Loading Configuration from a File

```hcl
resource "newrelic_fleet_configuration" "infra" {
  name                = "Production Infrastructure Config"
  agent_type          = "NRInfra"
  managed_entity_type = "HOST"

  version {
    configuration_content = file("${path.module}/configs/infra-v1.yaml")
  }
}
```

### Multiple Versions

Add multiple `version` blocks to maintain several versions of the configuration simultaneously. Each version must have unique content.

```hcl
resource "newrelic_fleet_configuration" "infra" {
  name                = "Production Infrastructure Config"
  agent_type          = "NRInfra"
  managed_entity_type = "HOST"

  version {
    configuration_content = file("${path.module}/configs/infra-v1.yaml")
  }

  version {
    configuration_content = file("${path.module}/configs/infra-v2.yaml")
  }
}
```

### Kubernetes Configuration

```hcl
resource "newrelic_fleet_configuration" "k8s" {
  name                = "Production K8s Config"
  agent_type          = "NRDOT"
  managed_entity_type = "KUBERNETESCLUSTER"

  version {
    configuration_content = file("${path.module}/configs/k8s-config.yaml")
  }
}
```

### Referencing Version Metadata in Outputs

```hcl
resource "newrelic_fleet_configuration" "infra" {
  name                = "my-infra-config"
  agent_type          = "NRInfra"
  managed_entity_type = "HOST"

  version {
    configuration_content = file("${path.module}/config.yaml")
  }
}

output "latest_version_guid" {
  value = newrelic_fleet_configuration.infra.latest_version_entity_id
}

output "latest_version_number" {
  value = newrelic_fleet_configuration.infra.latest_version_number
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the configuration.
* `agent_type` - (Required, ForceNew) The type of agent this configuration is for. Valid values: `NRInfra`, `NRDOT`, `FluentBit`, `NRPrometheusAgent`. **Cannot be changed after creation.**
* `managed_entity_type` - (Required, ForceNew) The type of entities this configuration manages. Valid values: `HOST`, `KUBERNETESCLUSTER`. **Cannot be changed after creation.**
* `operating_system` - (Optional, ForceNew) The operating system this configuration targets. Valid values: `LINUX`, `WINDOWS`. Applicable to `HOST` configurations only — must not be set when `managed_entity_type` is `KUBERNETESCLUSTER`. **Cannot be changed after creation.**
* `version` - (Required) One or more version blocks. At least one is required. See [Nested `version` blocks](#nested-version-blocks) below.
* `organization_id` - (Optional, ForceNew) The organization ID. Auto-fetched from the account when not provided. **Cannot be changed after creation.**

### Nested `version` blocks

Each `version` block supports the following argument:

* `configuration_content` - (Required) The YAML or JSON content for this version. Must be unique across all `version` blocks in the resource. Use `file()` to load content from a file: `file("${path.module}/config.yaml")`.

The following attributes are exported from each `version` block:

* `version_number` - The version number assigned by the API (1, 2, 3, …).
* `version_entity_id` - The entity GUID of this version.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The entity GUID of the configuration (same as `configuration_id`).
* `configuration_id` - The entity GUID of the configuration.
* `latest_version_number` - The highest version number across all versions.
* `latest_version_entity_id` - The entity GUID of the highest-numbered version.
* `total_versions` - Total number of versions currently in the configuration.

## Working with Versions

### Version Immutability

Version content is **immutable** - the API does not support updating the content of an existing version. If you attempt to modify `configuration_content` of an already-applied `version` block, Terraform will catch this at plan time and surface an error before any API call is made:

```
configuration_content cannot be modified in place - versions are immutable:
  - index 0: content was changed (edit detected)
```

To update the configuration in use:
- **Add** a new `version` block with the updated content.
- **Remove** the old `version` block whose content you no longer need.

Terraform applies removals (API deletes) before creates, so if you add and remove a block in the same `apply`, the old version is deleted first and the new one is created after.

### Unique Content Requirement

All `version` blocks within a resource must have distinct `configuration_content` values. Duplicate content is caught at plan time before any changes are applied:

```
duplicate configuration_content detected across version blocks
```

This also applies to rollback scenarios. If you previously had versions A → B and want to roll back by reintroducing A's content as a new version, add a new `version` block with A's content rather than restoring an old block - the new version will get a new version number from the API.

### Version Numbering

Version numbers are assigned sequentially by the API and are never reused or renumbered. When you remove a `version` block, the remaining versions keep their original numbers. For example, if you have versions 1, 2, and 3 and remove version 2, the configuration will have versions 1 and 3 - the API does not compact the sequence.

`latest_version_number` and `latest_version_entity_id` always reflect the highest-numbered version, regardless of how many versions exist.

### Externally Deleted Versions

If a version is deleted outside of Terraform (for example, via the API or the New Relic UI), the next `terraform plan` will show a warning for the affected version:

```
Warning: version entity not found
  version block at index N (version_entity_id = "...") no longer exists in the API.
  Remove the block from your configuration if the deletion was intentional,
  or run terraform apply to recreate it.
```

The warning indicates that Terraform will recreate the missing version on the next `apply`. If the deletion was intentional, remove the corresponding `version` block from your configuration before applying.

## Import

Fleet configurations can be imported using the configuration entity GUID:

```shell
terraform import newrelic_fleet_configuration.infra <configuration_guid>
```

Because `agent_type`, `managed_entity_type`, `operating_system`, and `name` are not returned by the configuration read API, a plain GUID import leaves those fields empty. Use the compound import ID to reconstruct them in a single step:

```shell
terraform import newrelic_fleet_configuration.infra \
  <configGUID>:<orgID>:<agentType>:<managedEntityType>:<operatingSystem>:<name>
```

For `KUBERNETESCLUSTER` configurations (where `operating_system` is not set), leave that segment empty:

```shell
terraform import newrelic_fleet_configuration.infra \
  <configGUID>:<orgID>:NRInfra:KUBERNETESCLUSTER::<name>
```