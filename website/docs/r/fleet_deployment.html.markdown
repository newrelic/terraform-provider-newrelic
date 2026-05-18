---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet_deployment"
sidebar_current: "docs-newrelic-resource-fleet-deployment"
description: |-
  Create and manage New Relic fleet deployments.
---

# Resource: newrelic\_fleet\_deployment

Use this resource to create and manage New Relic fleet deployments.

A fleet deployment defines the agent versions and optional configuration versions to roll out to a fleet. Each deployment belongs to a fleet and contains one or more `agent` blocks describing which agent type and version to deploy, and optionally which configuration version (from `newrelic_fleet_configuration`) to apply.

~> **Note:** Deployments can only be updated while in the `CREATED` phase. Once the fleet backend begins executing the deployment (phase advances to `IN_PROGRESS`, `FAILED`, or `COMPLETED`), any attempt to change `name`, `description`, `agent`, or `tags` will be **blocked at plan time** with an error. Run `terraform destroy` to remove the deployment from state and then re-create it with the desired configuration. If `terraform destroy` itself fails because the deployment is actively executing, the resource will be removed from Terraform state with a warning — once the deployment reaches a terminal phase (`COMPLETED` or `FAILED`) you can clean it up manually in the New Relic UI.

## Example Usage

### Basic Deployment

```hcl
resource "newrelic_fleet_deployment" "infra" {
  fleet_id    = newrelic_fleet.prod.id
  name        = "Production Infra Deployment"
  description = "Deploys NRInfra v1.58.0 to the production fleet"

  agent {
    agent_type = "NRInfra"
    version    = "1.58.0"
  }
}
```

### Deployment Linked to a Configuration Version

```hcl
resource "newrelic_fleet_configuration" "infra_cfg" {
  name                = "Production Infra Config"
  agent_type          = "NRInfra"
  managed_entity_type = "HOST"

  version {
    configuration_content = <<-EOT
      log:
        level: info
    EOT
  }
}

resource "newrelic_fleet_deployment" "infra" {
  fleet_id    = newrelic_fleet.prod.id
  name        = "Production Infra Deployment"
  description = "Deploys NRInfra v1.58.0 with the production config"

  agent {
    agent_type               = "NRInfra"
    version                  = "1.58.0"
    configuration_version_id = newrelic_fleet_configuration.infra_cfg.latest_version_entity_id
  }
}
```

### Multiple Agents

```hcl
resource "newrelic_fleet_deployment" "full_stack" {
  fleet_id = newrelic_fleet.prod.id
  name     = "Full Stack Deployment"

  agent {
    agent_type               = "NRInfra"
    version                  = "1.58.0"
    configuration_version_id = newrelic_fleet_configuration.infra_cfg.latest_version_entity_id
  }

  agent {
    agent_type = "FluentBit"
    version    = "3.2.0"
  }
}
```

### With Tags

```hcl
resource "newrelic_fleet_deployment" "infra" {
  fleet_id = newrelic_fleet.prod.id
  name     = "Production Deployment"

  agent {
    agent_type = "NRInfra"
    version    = "1.58.0"
  }

  tags = ["environment:production", "team:platform"]
}
```

## Argument Reference

The following arguments are supported:

* `fleet_id` - (Required, ForceNew) The entity GUID of the fleet this deployment belongs to. **Cannot be changed after creation.**
* `name` - (Optional) The name of the deployment.
* `description` - (Optional) A description of the deployment.
* `agent` - (Required on create, may be empty on update) One or more agent blocks. At least one is required when creating a deployment. On update, the list may be set to empty (`agent = []`) to uninstall all agent assignments from the deployment. Each `agent_type` may appear at most once per deployment. See [Nested `agent` blocks](#nested-agent-blocks) below.
* `tags` - (Optional) A list of tags in `key:value1,value2` format.
* `organization_id` - (Optional, ForceNew) The organization ID. Auto-fetched from the account when not provided. **Cannot be changed after creation.**

### Nested `agent` blocks

Each `agent` block supports:

* `agent_type` - (Required) The agent type. Valid values: `NRInfra`, `NRDOT`, `FluentBit`, `NRPrometheusAgent`.
* `version` - (Required) The agent version string to deploy (e.g. `"1.58.0"`).
* `configuration_version_id` - (Required) A configuration version entity GUID (from `newrelic_fleet_configuration`) to associate with this agent in the deployment.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The entity GUID of the deployment (same as `deployment_id`).
* `deployment_id` - The entity GUID of the deployment.
* `phase` - The current phase of the deployment. Possible values: `CREATED`, `IN_PROGRESS`, `FAILED`, `COMPLETED`.

## Import

Fleet deployments can be imported using the deployment entity GUID:

```shell
terraform import newrelic_fleet_deployment.infra <deployment_guid>
```
