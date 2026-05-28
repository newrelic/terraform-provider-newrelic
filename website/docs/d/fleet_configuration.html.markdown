---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet_configuration"
sidebar_current: "docs-newrelic-datasource-fleet-configuration"
description: |-
  Fetches the content of a New Relic fleet configuration and its version metadata.
---

# Data Source: newrelic\_fleet\_configuration

Use this data source to fetch the content and version metadata of an existing New Relic fleet configuration. Three mutually exclusive lookup modes are supported: by **configuration GUID**, by **version entity GUID**, or by **name**.

## Example Usage

### Look Up by Configuration GUID

Fetches the content of the **latest** version of the configuration identified by its entity GUID. Also returns the GUIDs of all versions.

```hcl
data "newrelic_fleet_configuration" "by_id" {
  configuration_id = "NjQyNTg2NXxOR0VQ..."
}

output "latest_content" {
  value = data.newrelic_fleet_configuration.by_id.configuration_content
}

output "all_version_guids" {
  value = data.newrelic_fleet_configuration.by_id.version_entity_ids
}
```

### Look Up by Version Entity GUID

Fetches the content of a **specific version** identified by its version entity GUID. Also resolves and returns the parent configuration GUID.

```hcl
data "newrelic_fleet_configuration" "by_version" {
  version_entity_id = "NjQyNTg2NXxOR0VQfEFHRU5UX0NPTkZJR1VSQVRJT05fVkVSU0lPTnw..."
}

output "version_content" {
  value = data.newrelic_fleet_configuration.by_version.configuration_content
}

output "parent_config_guid" {
  value = data.newrelic_fleet_configuration.by_version.configuration_id
}
```

### Look Up by Name

Fetches the content of the **latest** version of the configuration matching the given name. Also returns the GUIDs of all versions. The first matching configuration is returned if multiple share the same name.

```hcl
data "newrelic_fleet_configuration" "by_name" {
  name = "Production Infrastructure Config"
}

output "latest_content" {
  value = data.newrelic_fleet_configuration.by_name.configuration_content
}
```

### Reference from a Managed Resource

A common pattern is to use this data source alongside a `newrelic_fleet_configuration` resource to read version metadata:

```hcl
resource "newrelic_fleet_configuration" "infra" {
  name                = "my-infra-config"
  agent_type          = "NRInfra"
  managed_entity_type = "HOST"

  version {
    configuration_content = file("${path.module}/config.yaml")
  }
}

data "newrelic_fleet_configuration" "infra" {
  configuration_id = newrelic_fleet_configuration.infra.configuration_id
  depends_on       = [newrelic_fleet_configuration.infra]
}

output "all_version_guids" {
  value = data.newrelic_fleet_configuration.infra.version_entity_ids
}
```

## Argument Reference

Exactly one of the following lookup inputs must be specified:

* `configuration_id` - (Optional) The entity GUID of the fleet configuration. Returns the content of its **latest** version. This attribute is also populated as an output when looking up by `version_entity_id`.
* `version_entity_id` - (Optional) The entity GUID of a specific configuration version. Returns the content of that exact version.
* `name` - (Optional) The name of the fleet configuration. Returns the content of its **latest** version. The first matching configuration is returned if multiple share the same name.

The following optional argument is supported in all modes:

* `organization_id` - (Optional) The organization ID. Resolved automatically from the provider account when omitted.

## Attributes Reference

The following attributes are exported. Availability varies by lookup mode — see the table below.

* `configuration_content` - The raw YAML/JSON content of the resolved version.
* `configuration_id` - The entity GUID of the fleet configuration.
* `organization_id` - The organization ID the configuration belongs to.
* `latest_version_entity_id` - The entity GUID of the latest (highest-numbered) version.
* `version_entity_ids` - Entity GUIDs of all versions ordered by version number, oldest first.

### Attribute Availability by Lookup Mode

| Attribute | `configuration_id` | `version_entity_id` | `name` |
|---|---|---|---|
| `configuration_content` | ✓ (latest version) | ✓ (exact version) | ✓ (latest version) |
| `configuration_id` | ✓ (input) | ✓ (resolved from API) | ✓ (resolved from API) |
| `organization_id` | ✓ | ✓ | ✓ |
| `latest_version_entity_id` | ✓ | — | ✓ |
| `version_entity_ids` | ✓ | — | ✓ |

-> **NOTE:** `latest_version_entity_id` and `version_entity_ids` are not populated when looking up by `version_entity_id`, because a single version GUID does not carry the full version history of its parent configuration.
