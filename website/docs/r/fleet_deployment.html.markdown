---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet_deployment"
sidebar_current: "docs-newrelic-resource-fleet-deployment"
description: |-
  Create and manage New Relic fleet deployments.
---

# Resource: newrelic\_fleet\_deployment

Use this resource to create and manage New Relic fleet deployments.

A fleet deployment defines the agent versions and configuration versions to roll out to a fleet. Each deployment belongs to a fleet and may contain zero or more `agent` blocks describing which agent type and version to deploy and which configuration version (from `newrelic_fleet_configuration`) to apply.

~> **Note: Phase-gate immutability.** Deployments can only be modified while in the `CREATED` phase. Once the fleet backend begins executing the deployment (phase advances to `IN_PROGRESS`, `FAILED`, or `COMPLETED`), any attempt to change `name`, `description`, `agent`, or `tags` is **blocked at plan time** with a clear error. The recommended recovery path is `terraform state rm <resource_address>` (or a [`removed` block](https://developer.hashicorp.com/terraform/language/state/remove)) to drop the executed deployment from Terraform state, then re-declare a fresh deployment with the desired configuration. `terraform destroy` works only if you have no pending changes to the deployment in your HCL — otherwise the same plan-time gate fires during the destroy plan.

## Example Usage

### Basic Deployment

```hcl
resource "newrelic_fleet_configuration" "infra_cfg" {
  name                  = "Production Infra Config"
  agent_type            = "NRInfra"
  managed_entity_type   = "HOST"
  operating_system      = "LINUX"
  configuration_content = <<-EOT
    log:
      level: info
  EOT
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

Each `agent_type` may appear at most once per deployment.

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
    agent_type               = "FluentBit"
    version                  = "3.2.0"
    configuration_version_id = newrelic_fleet_configuration.fb_cfg.latest_version_entity_id
  }
}
```

### Zero-Agent Deployment

A deployment can be created or updated with zero `agent` blocks — for example, to drain all agent assignments from an existing deployment, or to seed a deployment record that will have agents added later (while still in `CREATED` phase).

```hcl
resource "newrelic_fleet_deployment" "drained" {
  fleet_id = newrelic_fleet.prod.id
  name     = "Drained deployment"
  # No agent blocks — uninstalls all agent assignments.
}
```

### Pinning to a Specific Configuration Version

By default, referencing `newrelic_fleet_configuration.<name>.latest_version_entity_id` ties the deployment to whichever version is current at plan time. Updating the configuration's `configuration_content` will change `latest_version_entity_id`, which in turn proposes an update to any deployment referencing it. If the deployment has already left `CREATED`, that update will be **blocked by the phase-gate** — see the note above.

For long-lived deployments that should remain stable, pin to a specific historical version instead:

```hcl
resource "newrelic_fleet_deployment" "pinned" {
  fleet_id = newrelic_fleet.prod.id
  name     = "Stable v1 rollout"

  agent {
    agent_type = "NRInfra"
    version    = "1.58.0"
    # First-ever version of this configuration — won't drift when the config is updated.
    configuration_version_id = newrelic_fleet_configuration.infra_cfg.version_entity_ids[0]
  }
}
```

### With Tags

```hcl
resource "newrelic_fleet_deployment" "infra" {
  fleet_id = newrelic_fleet.prod.id
  name     = "Production Deployment"

  agent {
    agent_type               = "NRInfra"
    version                  = "1.58.0"
    configuration_version_id = newrelic_fleet_configuration.infra_cfg.latest_version_entity_id
  }

  tags = ["environment:production", "team:platform"]
}
```

## Argument Reference

The following arguments are supported:

* `fleet_id` - (Required, ForceNew) The entity GUID of the fleet this deployment belongs to. **Cannot be changed after creation.**
* `name` - (Required) The name of the deployment. Updatable while the deployment is in `CREATED` phase.
* `description` - (Optional) A description of the deployment.
* `agent` - (Optional) Zero or more `agent` blocks. An empty list is accepted on both create and update — useful to drain agent assignments. Each `agent_type` may appear at most once per deployment. See [Nested `agent` blocks](#nested-agent-blocks) below.
* `tags` - (Optional) A list of tags in `key:value1,value2` format.
* `organization_id` - (Optional, ForceNew) The organization ID. Auto-fetched from the account when not provided. **Cannot be changed after creation.**

### Nested `agent` blocks

Each `agent` block supports:

* `agent_type` - (Required) The agent type. Valid values: `NRInfra`, `NRDOT`, `FluentBit`, `NRPrometheusAgent`.
* `version` - (Required) The agent version string to deploy (e.g. `"1.58.0"`).
* `configuration_version_id` - (Required) The entity GUID of the configuration version (from `newrelic_fleet_configuration`) to associate with this agent. Reference `latest_version_entity_id` to follow the current version, or `version_entity_ids[N]` to pin to a specific historical version.

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
