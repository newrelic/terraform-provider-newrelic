---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet_configuration"
sidebar_current: "docs-newrelic-resource-fleet-configuration"
description: |-
  Create and manage fleet agent configurations in New Relic.
---

# Resource: newrelic\_fleet\_configuration

Use this resource to create and manage fleet agent configurations in New Relic. Configurations define settings for New Relic agents deployed across your infrastructure.

## Example Usage

### Configuration from Inline Content

```hcl
resource "newrelic_fleet_configuration" "infra_config" {
  name                   = "Production Infrastructure Config"
  agent_type             = "NRInfra"
  managed_entity_type    = "HOST"
  configuration_content  = jsonencode({
    log_level = "info"
    custom_attributes = {
      environment = "production"
      team = "platform"
    }
  })
}
```

### Configuration from File

```hcl
resource "newrelic_fleet_configuration" "k8s_config" {
  name                  = "Kubernetes APM Config"
  agent_type            = "NRDOT"
  managed_entity_type   = "KUBERNETESCLUSTER"
  configuration_file    = "${path.module}/configs/apm-config.yaml"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the configuration. **Note**: Cannot be changed after creation.
* `agent_type` - (Required) The type of agent this configuration targets. Allowed values: `NRInfra`, `NRDOT`, `FluentBit`, `NRPrometheusAgent`. **Note**: Cannot be changed after creation.
* `managed_entity_type` - (Required) The type of entities this configuration applies to. Allowed values: `HOST`, `KUBERNETESCLUSTER`. **Note**: Cannot be changed after creation.
* `configuration_file` - (Optional) Path to the configuration file (JSON/YAML). Mutually exclusive with `configuration_content`. **Note**: Cannot be changed after creation.
* `configuration_content` - (Optional) Inline configuration content (JSON/YAML). Mutually exclusive with `configuration_file`. **Note**: Cannot be changed after creation.
* `organization_id` - (Optional) The organization ID. If not provided, it will be auto-fetched. **Note**: Cannot be changed after creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID (entity GUID) of the configuration.
* `entity_guid` - The entity GUID of the configuration.
* `blob_version_entity` - Information about the initial version:
  * `version` - The version number (always 1 for initial version).
  * `guid` - The version entity GUID.
  * `blob_id` - The blob ID.

## Import

Fleet configurations can be imported using the entity GUID:

```
$ terraform import newrelic_fleet_configuration.infra_config <entity_guid>
```

## Additional Information

### Configuration Format

The configuration content can be in JSON or YAML format. The structure depends on the agent type:

- **NRInfra**: Infrastructure agent configuration
- **NRDOT**: APM agent configuration
- **FluentBit**: Fluent Bit logging configuration
- **NRPrometheusAgent**: Prometheus agent configuration

### Versions

When a configuration is created, an initial version (version 1) is automatically created. To add more versions to an existing configuration, use the `newrelic_fleet_configuration_version` resource.
