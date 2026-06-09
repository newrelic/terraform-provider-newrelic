---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet_configuration"
sidebar_current: "docs-newrelic-resource-fleet-configuration"
description: |-
  Create and manage New Relic fleet configurations for centralized agent management.
---

# Resource: newrelic\_fleet\_configuration

Use this resource to create and manage New Relic fleet configurations for centralized agent management.

A fleet configuration holds versioned agent settings. The configuration content is immutable — each change to `configuration_content` creates a new version on the API automatically, similar to how AWS launch templates work. The resource ID (the configuration entity GUID) never changes across updates. Use the `newrelic_fleet_configuration` data source to access the content of a specific historical version.

## Example Usage

### Basic Infrastructure Agent Configuration

```hcl
resource "newrelic_fleet_configuration" "infra" {
  name                  = "Production Infrastructure Config"
  agent_type            = "NRInfra"
  managed_entity_type   = "HOST"
  operating_system      = "LINUX"
  configuration_content = <<-EOT
    log:
      level: info
      file: /var/log/newrelic-infra/newrelic-infra.log
    metrics:
      enabled: true
      system_sample_rate: 15
  EOT
}
```

### Loading Configuration from a File

```hcl
resource "newrelic_fleet_configuration" "infra" {
  name                  = "Production Infrastructure Config"
  agent_type            = "NRInfra"
  managed_entity_type   = "HOST"
  operating_system      = "LINUX"
  configuration_content = file("${path.module}/configs/infra.yaml")
}
```

### Kubernetes Configuration

```hcl
resource "newrelic_fleet_configuration" "k8s" {
  name                  = "Production K8s Config"
  agent_type            = "NRDOT"
  managed_entity_type   = "KUBERNETESCLUSTER"
  configuration_content = file("${path.module}/configs/k8s-config.yaml")
}
```

### Referencing Version Metadata in Outputs

```hcl
resource "newrelic_fleet_configuration" "infra" {
  name                  = "my-infra-config"
  agent_type            = "NRInfra"
  managed_entity_type   = "HOST"
  operating_system      = "LINUX"
  configuration_content = file("${path.module}/config.yaml")
}

output "latest_version_guid" {
  value = newrelic_fleet_configuration.infra.latest_version_entity_id
}

output "latest_version_number" {
  value = newrelic_fleet_configuration.infra.latest_version_number
}
```

### Accessing a Previous Version via the Data Source

Use `version_entity_ids` to reference an older version with the data source:

```hcl
resource "newrelic_fleet_configuration" "infra" {
  name                  = "my-infra-config"
  agent_type            = "NRInfra"
  managed_entity_type   = "HOST"
  operating_system      = "LINUX"
  configuration_content = file("${path.module}/config.yaml")
}

# Fetch the content of the first version that was ever created.
data "newrelic_fleet_configuration" "infra_v1" {
  version_entity_id = newrelic_fleet_configuration.infra.version_entity_ids[0]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the configuration. **Changing this forces resource recreation** — the API does not support renaming a configuration in place.
* `agent_type` - (Required, ForceNew) The type of agent this configuration is for. Valid values: `NRInfra`, `NRDOT`, `FluentBit`, `NRPrometheusAgent`. **Cannot be changed after creation.**
* `managed_entity_type` - (Required, ForceNew) The type of entities this configuration manages. Valid values: `HOST`, `KUBERNETESCLUSTER`. **Cannot be changed after creation.**
* `operating_system` - (Optional, ForceNew) The operating system this configuration targets. Valid values: `LINUX`, `WINDOWS`. Applicable to `HOST` configurations only — must not be set when `managed_entity_type` is `KUBERNETESCLUSTER`. **Cannot be changed after creation.**
* `configuration_content` - (Required) The YAML or JSON content for this configuration. Use `file()` to load content from a file. Each change to this field creates a new immutable version on the API; the resource ID remains constant.
* `organization_id` - (Optional, ForceNew) The organization ID. Auto-fetched from the account when not provided. **Cannot be changed after creation.**

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The entity GUID of the configuration (same as `configuration_id`).
* `configuration_id` - The entity GUID of the configuration.
* `latest_version_number` - The highest version number across all versions created so far.
* `latest_version_entity_id` - The entity GUID of the highest-numbered version.
* `total_versions` - Total number of versions currently in the configuration.
* `version_entity_ids` - A list of entity GUIDs for all versions, sorted oldest-first. Use with the `newrelic_fleet_configuration` data source to retrieve the content of a specific historical version.

## How Versioning Works

Every time `configuration_content` changes, the provider creates a new immutable version on the API. The resource itself (identified by `id` / `configuration_id`) is never recreated — only a new version record is appended. This mirrors how AWS launch templates work: you edit the content and Terraform increments the version counter automatically.

To roll back to a previous configuration, simply set `configuration_content` to the older content — the provider will create a new version with that content, and `latest_version_entity_id` will point to it.

Previous versions are accessible via the `version_entity_ids` list and the `newrelic_fleet_configuration` data source. Versions are never deleted on update; they accumulate until the configuration resource itself is destroyed.

## Out-of-band drift warnings

If a version is deleted outside of Terraform (UI, API, or another tool), the next `plan` or `refresh` will surface a warning so you understand why state changed:

```
Warning: Fleet configuration drift: 1 version(s) deleted out-of-band
  The following version entity GUIDs were tracked in Terraform state but no longer
  exist in the New Relic API: [...]
  State has been synced to reflect the API.
```

If the previously-tracked **latest** version was the one deleted, an additional, stronger warning fires explaining that `configuration_content` has been refreshed from the new latest version on the API. If your declared content differs from that new latest, the next `apply` will create a new version restoring your declared content — this is the expected, self-healing behavior.

## Import

Fleet configurations can be imported using a composite ID of `<configuration_guid>:<managed_entity_type>`:

```shell
terraform import newrelic_fleet_configuration.infra <configuration_guid>:HOST
terraform import newrelic_fleet_configuration.infra <configuration_guid>:KUBERNETESCLUSTER
```

The `managed_entity_type` portion is required because the New Relic API does not return it via the entity lookup query (a GraphQL schema constraint). All other attributes — `name`, `agent_type`, `operating_system`, `organization_id`, `configuration_content` — are resolved automatically from the API.
