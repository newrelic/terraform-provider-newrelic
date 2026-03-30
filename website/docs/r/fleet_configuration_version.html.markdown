---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet_configuration_version"
sidebar_current: "docs-newrelic-resource-fleet-configuration-version"
description: |-
  Add new versions to existing fleet configurations in New Relic.
---

# Resource: newrelic\_fleet\_configuration\_version

Use this resource to add new versions to existing fleet configurations. Configuration versions are immutable snapshots of agent settings that can be deployed to fleets.

## Example Usage

### Add Version with Inline Content

```hcl
resource "newrelic_fleet_configuration" "base" {
  name                   = "Infrastructure Config"
  agent_type             = "NRInfra"
  managed_entity_type    = "HOST"
  configuration_content  = jsonencode({
    log_level = "info"
  })
}

resource "newrelic_fleet_configuration_version" "v2" {
  configuration_id      = newrelic_fleet_configuration.base.id
  configuration_content = jsonencode({
    log_level = "debug"
    custom_attributes = {
      version = "2"
    }
  })
}
```

### Add Version from File

```hcl
resource "newrelic_fleet_configuration_version" "v3" {
  configuration_id   = newrelic_fleet_configuration.base.id
  configuration_file = "${path.module}/configs/v3-config.yaml"
}
```

## Argument Reference

The following arguments are supported:

* `configuration_id` - (Required) The configuration entity ID to add a version to. **Note**: Cannot be changed after creation.
* `configuration_file` - (Optional) Path to the configuration file (JSON/YAML). Mutually exclusive with `configuration_content`. **Note**: Cannot be changed after creation.
* `configuration_content` - (Optional) Inline configuration content (JSON/YAML). Mutually exclusive with `configuration_file`. **Note**: Cannot be changed after creation.
* `organization_id` - (Optional) The organization ID. If not provided, it will be auto-fetched. **Note**: Cannot be changed after creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The version entity GUID.
* `version` - The version number (auto-incremented).
* `blob_id` - The blob ID.

## Import

Fleet configuration versions can be imported using the version entity GUID:

```
$ terraform import newrelic_fleet_configuration_version.v2 <version_guid>
```

## Additional Information

### Version Numbering

Versions are automatically numbered starting from 1 (created with the base configuration) and incrementing with each new version added.

### Immutability

Configuration versions are immutable once created. To make changes, create a new version rather than modifying an existing one.

### Deployments

To deploy a specific configuration version to a fleet, reference the version ID in a `newrelic_fleet_deployment` resource.
