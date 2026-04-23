---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet_configuration_version"
sidebar_current: "docs-newrelic-resource-fleet-configuration-version"
description: |-
  Add new versions to existing New Relic fleet configurations.
---

# Resource: newrelic\_fleet\_configuration\_version

Use this resource to add new versions to existing New Relic fleet configurations.

Fleet configuration versioning allows you to maintain multiple versions of agent settings. Each version is immutable once created, providing a clear audit trail and the ability to roll back deployments to previous configurations.

**Important**: This resource adds a new version to an **existing** configuration created by `newrelic_fleet_configuration`. The initial configuration automatically creates version 1.

## Example Usage

### Adding Version 2 to Existing Configuration

```hcl
# Create the initial configuration (version 1)
resource "newrelic_fleet_configuration" "infra" {
  name                = "Production Infrastructure Config"
  agent_type          = "INFRASTRUCTURE"
  managed_entity_type = "HOST"

  configuration_content = <<-EOT
    log:
      level: info
    metrics:
      enabled: true
  EOT
}

# Add version 2 with updated settings
resource "newrelic_fleet_configuration_version" "infra_v2" {
  configuration_id = newrelic_fleet_configuration.infra.configuration_id

  configuration_content = <<-EOT
    log:
      level: debug
    metrics:
      enabled: true
      system_sample_rate: 10
  EOT
}
```

### Multiple Versions with Progressive Rollout

```hcl
resource "newrelic_fleet_configuration" "k8s" {
  name                = "Kubernetes Config"
  agent_type          = "KUBERNETES"
  managed_entity_type = "KUBERNETESCLUSTER"

  configuration_content = <<-EOT
    cluster:
      enabled: true
    version: 1
  EOT
}

# Version 2: Add prometheus
resource "newrelic_fleet_configuration_version" "k8s_v2" {
  configuration_id = newrelic_fleet_configuration.k8s.configuration_id

  configuration_content = <<-EOT
    cluster:
      enabled: true
    prometheus:
      enabled: true
    version: 2
  EOT
}

# Version 3: Enhanced logging
resource "newrelic_fleet_configuration_version" "k8s_v3" {
  configuration_id = newrelic_fleet_configuration.k8s.configuration_id

  configuration_content = <<-EOT
    cluster:
      enabled: true
    prometheus:
      enabled: true
    logging:
      enabled: true
      level: info
    version: 3
  EOT
}
```

### Using Configuration from File

```hcl
resource "newrelic_fleet_configuration" "base" {
  name                = "Base Config"
  agent_type          = "INFRASTRUCTURE"
  managed_entity_type = "HOST"
  configuration_content = file("${path.module}/configs/v1.yml")
}

resource "newrelic_fleet_configuration_version" "v2" {
  configuration_id         = newrelic_fleet_configuration.base.configuration_id
  configuration_file_path  = "${path.module}/configs/v2.yml"
}

resource "newrelic_fleet_configuration_version" "v3" {
  configuration_id         = newrelic_fleet_configuration.base.configuration_id
  configuration_file_path  = "${path.module}/configs/v3.yml"
}
```

## Argument Reference

The following arguments are supported:

* `configuration_id` - (Required) The ID (GUID) of the configuration to add this version to. **Note**: This cannot be changed after creation (forces new resource).
* `configuration_file_path` - (Optional) Path to a file containing the configuration content. **Mutually exclusive with `configuration_content`**. **Note**: This cannot be changed after creation (forces new resource).
* `configuration_content` - (Optional) Inline configuration content (YAML format). **Mutually exclusive with `configuration_file_path`**. **Note**: This cannot be changed after creation (forces new resource).
* `organization_id` - (Optional) The organization ID. If not provided, it will be automatically fetched from your account. **Note**: This cannot be changed after creation (forces new resource).

**Note**: You must provide either `configuration_file_path` or `configuration_content`, but not both.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the configuration version entity.
* `version_number` - The version number of this configuration version (e.g., 2, 3, 4...).

## Version Management Strategy

### How Versioning Works

1. **Initial Configuration**: When you create a `newrelic_fleet_configuration`, it automatically creates version 1
2. **Adding Versions**: Use `newrelic_fleet_configuration_version` to add versions 2, 3, 4, etc.
3. **Immutability**: Once created, version content cannot be changed (all fields are ForceNew)
4. **Deployment**: Use fleet deployments to specify which version to deploy to which fleet ring

### Version Lifecycle

```hcl
# Day 1: Create config with v1
resource "newrelic_fleet_configuration" "app" {
  name = "App Config"
  # ... v1 content
}

# Day 2: Add v2 for testing
resource "newrelic_fleet_configuration_version" "app_v2" {
  configuration_id = newrelic_fleet_configuration.app.configuration_id
  # ... v2 content with new features
}

# Day 3: Add v3 for production rollout
resource "newrelic_fleet_configuration_version" "app_v3" {
  configuration_id = newrelic_fleet_configuration.app.configuration_id
  # ... v3 content with fixes
}

# Versions exist independently - can deploy different versions to different rings
```

### Best Practices

1. **Semantic Content**: Include version info in configuration content for easier troubleshooting:
   ```yaml
   # In configuration_content
   metadata:
     version: 2
     description: "Added prometheus metrics"
   ```

2. **Progressive Rollout**: Test new versions on canary/test rings before production
3. **Version Tracking**: Use meaningful resource names (e.g., `config_v2_prometheus`, `config_v3_enhanced_logging`)
4. **Immutable Versions**: Never try to "update" a version - create a new one instead
5. **Cleanup**: Old versions can be deleted when no longer needed (ensure no deployments reference them)

## Relationship with Other Resources

```
newrelic_fleet                          # The fleet to deploy to
newrelic_fleet_configuration            # Creates base config (v1)
newrelic_fleet_configuration_version    # Adds v2, v3, v4...
newrelic_fleet_deployment              # Deploys specific versions to fleet rings
```

## Import

Fleet configuration versions can be imported using the version entity GUID:

```
$ terraform import newrelic_fleet_configuration_version.v2 <version_entity_guid>
```

**Note**: When importing, you must also set the `configuration_id` in your Terraform configuration, as this is required for the resource to function correctly. The `configuration_content` and `configuration_file_path` will not be populated, as the API does not return the raw content.

## Deletion Behavior

Deleting a `newrelic_fleet_configuration_version` resource removes that specific version from the configuration. The base configuration and other versions remain intact. Ensure no active deployments reference the version before deletion.
