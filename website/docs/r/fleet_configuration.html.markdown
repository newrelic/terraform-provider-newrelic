---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet_configuration"
sidebar_current: "docs-newrelic-resource-fleet-configuration"
description: |-
  Create and manage New Relic fleet configurations for agent management.
---

# Resource: newrelic\_fleet\_configuration

Use this resource to create and manage New Relic fleet configurations for centralized agent management.

Fleet configurations allow you to define and version agent settings that can be deployed to your fleets. Each configuration is specific to an agent type (INFRASTRUCTURE or KUBERNETES) and managed entity type (HOST or KUBERNETESCLUSTER).

## Example Usage

### Infrastructure Agent Configuration

```hcl
resource "newrelic_fleet_configuration" "infra_config" {
  name                = "Production Infrastructure Config"
  agent_type          = "INFRASTRUCTURE"
  managed_entity_type = "HOST"

  configuration_content = <<-EOT
    log:
      level: info
      file: /var/log/newrelic-infra/newrelic-infra.log
    metrics:
      enabled: true
      system_sample_rate: 15
    integrations:
      - name: nri-docker
        enabled: true
  EOT
}
```

### Kubernetes Agent Configuration

```hcl
resource "newrelic_fleet_configuration" "k8s_config" {
  name                = "Production K8s Config"
  agent_type          = "KUBERNETES"
  managed_entity_type = "KUBERNETESCLUSTER"

  configuration_content = <<-EOT
    cluster:
      enabled: true
      name: production-cluster
    prometheus:
      enabled: true
      scrape_interval: 30s
    logging:
      enabled: true
      level: info
  EOT
}
```

### Using Configuration from File

```hcl
resource "newrelic_fleet_configuration" "file_config" {
  name                     = "Infrastructure Config from File"
  agent_type               = "INFRASTRUCTURE"
  managed_entity_type      = "HOST"
  configuration_file_path  = "${path.module}/configs/infra-config.yml"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the configuration. This can be changed after creation.
* `agent_type` - (Required) The type of agent this configuration is for. Valid values are `INFRASTRUCTURE` or `KUBERNETES`. **Note**: This cannot be changed after creation (forces new resource).
* `managed_entity_type` - (Required) The type of entities this configuration manages. Valid values are `HOST` or `KUBERNETESCLUSTER`. **Note**: This cannot be changed after creation (forces new resource).
* `configuration_file_path` - (Optional) Path to a file containing the configuration content. **Mutually exclusive with `configuration_content`**.
* `configuration_content` - (Optional) Inline configuration content (YAML format). **Mutually exclusive with `configuration_file_path`**.
* `organization_id` - (Optional) The organization ID. If not provided, it will be automatically fetched from your account. **Note**: This cannot be changed after creation (forces new resource).

**Note**: You must provide either `configuration_file_path` or `configuration_content`, but not both.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the configuration entity.
* `configuration_id` - The ID of the configuration (same as `id`).
* `version` - The current version number of the configuration.

## Configuration Content Format

The configuration content must be valid YAML appropriate for the agent type:

### Infrastructure Agent
Refer to the [New Relic Infrastructure agent configuration documentation](https://docs.newrelic.com/docs/infrastructure/install-infrastructure-agent/configuration/infrastructure-agent-configuration-settings/) for available settings.

### Kubernetes Agent
Refer to the [New Relic Kubernetes integration configuration documentation](https://docs.newrelic.com/docs/kubernetes-pixie/kubernetes-integration/installation/kubernetes-integration-install-configure/) for available settings.

## Versioning

Fleet configurations support versioning. When you create a configuration with this resource, it creates version 1. To add new versions to an existing configuration, use the `newrelic_fleet_configuration_version` resource.

**Important**: Updating the `configuration_content` or `configuration_file_path` in this resource will replace the entire configuration. To preserve version history, use the separate version resource instead.

## Import

Fleet configurations can be imported using the configuration GUID:

```
$ terraform import newrelic_fleet_configuration.infra_config <configuration_guid>
```

**Note**: When importing, `configuration_content` and `configuration_file_path` will not be populated in state, as the API does not return the raw configuration content in the same format. You should update your configuration file to match the actual configuration after import.
