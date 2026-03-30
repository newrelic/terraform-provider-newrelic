---
layout: "newrelic"
page_title: "New Relic: newrelic_fleet_deployment"
sidebar_current: "docs-newrelic-resource-fleet-deployment"
description: |-
  Create and manage New Relic fleet deployments.
---

# Resource: newrelic\_fleet\_deployment

Use this resource to create and manage New Relic fleet deployments for rolling out agent configurations.

## Example Usage

### Single Agent Deployment

```hcl
resource "newrelic_fleet_deployment" "infra_deployment" {
  fleet_id    = newrelic_fleet.linux_hosts.id
  name        = "Infrastructure Agent v1.70.0"
  description = "Deploy Infrastructure agent to production hosts"

  agent {
    agent_type                = "NRInfra"
    version                   = "1.70.0"
    configuration_version_ids = [
      "config-version-abc-123",
      "config-version-def-456"
    ]
  }

  tags = ["env:production"]
}
```

### Multiple Agent Deployment

```hcl
resource "newrelic_fleet_deployment" "multi_agent" {
  fleet_id    = newrelic_fleet.app_hosts.id
  name        = "Multi-Agent Deployment"
  description = "Deploy Infrastructure and .NET agents"

  agent {
    agent_type                = "NRInfra"
    version                   = "1.70.0"
    configuration_version_ids = ["config-infra-v1"]
  }

  agent {
    agent_type                = "NRDOT"
    version                   = "2.0.0"
    configuration_version_ids = ["config-dotnet-v1"]
  }

  tags = ["release:v2.0"]
}
```

## Argument Reference

The following arguments are supported:

* `fleet_id` - (Required) The ID of the fleet to deploy to. **Note**: This cannot be changed after creation.
* `name` - (Required) The name of the deployment.
* `agent` - (Required) One or more agent configuration blocks (see below).
* `description` - (Optional) The description of the deployment.
* `tags` - (Optional) A list of tags for the deployment in format `"key:value1,value2"`.

### Agent Configuration Block

The `agent` block supports:

* `agent_type` - (Required) The type of agent. Allowed values: `NRInfra`, `NRDOT`, `FluentBit`, `NRPrometheusAgent`.
* `version` - (Required) The agent version (e.g., `1.70.0`). Use `*` for wildcard version (only allowed for KUBERNETESCLUSTER fleets).
* `configuration_version_ids` - (Required) List of configuration version IDs to deploy with this agent.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the deployment.
* `phase` - The current phase of the deployment (e.g., CREATED, IN_PROGRESS, COMPLETED, FAILED).

## Import

Fleet deployments can be imported using the deployment ID:

```
$ terraform import newrelic_fleet_deployment.infra_deployment <deployment_id>
```
