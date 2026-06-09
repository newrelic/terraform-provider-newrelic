---
layout: "newrelic"
page_title: "Getting Started with New Relic Fleet Control"
sidebar_current: "docs-newrelic-provider-fleet-getting-started"
description: |-
  Use this guide to manage New Relic Fleet Control end-to-end with Terraform: create a fleet, define agent configurations, assign member entities, and trigger deployments.
---

# Getting Started with New Relic Fleet Control

New Relic Fleet Control lets you centrally manage agents across your infrastructure — define what configuration to run, which hosts or clusters to target, and roll out agent updates — all from a single control plane.

This guide walks through the full lifecycle using Terraform: creating a fleet, authoring a versioned agent configuration, assigning member entities to the fleet's rings, and triggering a deployment that pushes the configuration to those members.

-> **Note** Fleet Control requires a New Relic account with Agent Control enabled. The API key used by the provider must have Fleet Control permissions.

## Resources covered

| Resource / Data Source | What it does |
|---|---|
| `newrelic_fleet` | Creates the fleet — the top-level organizational container |
| `newrelic_fleet_configuration` | Defines versioned agent configuration content (YAML/JSON) |
| `newrelic_fleet_members` | Assigns entities (hosts, clusters) to the fleet's rings |
| `newrelic_fleet_deployment` | Triggers an agent version + configuration rollout to the fleet |
| `data.newrelic_fleet_members` | Reads which entities are currently in the fleet or a specific ring |
| `data.newrelic_fleet_configuration` | Reads configuration content and version GUIDs from an existing config |

## Step 1 — Create a fleet

A fleet is the top-level container that ties everything together. It declares what kind of infrastructure it manages (`HOST` or `KUBERNETESCLUSTER`) and, for host fleets, the operating system.

```hcl
resource "newrelic_fleet" "prod_linux" {
  name                = "Production Linux Hosts"
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
  description         = "Production fleet for Linux infrastructure agents"

  tags = [
    "environment:production",
    "team:platform",
  ]
}
```

`managed_entity_type` and `operating_system` are immutable after creation — if they ever need to change, Terraform will destroy and recreate the fleet. `name`, `description`, and `tags` can be updated in place at any time.

## Step 2 — Define an agent configuration

A fleet configuration holds the YAML or JSON settings for an agent type. The content is declared as a single top-level `configuration_content` attribute. Each change to that attribute creates a new immutable version on the API automatically — much like how AWS launch templates work. The resource ID never changes across updates.

```hcl
resource "newrelic_fleet_configuration" "infra" {
  name                  = "Production Infra Config"
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

For larger configurations, load the content from a file instead of an inline heredoc:

```hcl
  configuration_content = file("${path.module}/configs/infra.yaml")
```

After apply, the resource exports:

- `latest_version_number` — the current version number (1 on first create, increments on every content change).
- `latest_version_entity_id` — the entity GUID of the current latest version.
- `version_entity_ids` — every version GUID, oldest first. Use `version_entity_ids[N]` to pin a deployment to a specific historical version.
- `total_versions` — total number of versions accumulated.

### Rolling out a config update

To update the agent settings, simply edit `configuration_content` and run `terraform apply`. The provider creates a new version on the API; `latest_version_entity_id` and `latest_version_number` advance to point at it. Older versions are retained — they're accessible via `version_entity_ids` and the `data.newrelic_fleet_configuration` data source.

```hcl
resource "newrelic_fleet_configuration" "infra" {
  name                  = "Production Infra Config"
  agent_type            = "NRInfra"
  managed_entity_type   = "HOST"
  operating_system      = "LINUX"
  configuration_content = file("${path.module}/configs/infra.yaml")
}
```

### Rolling back

To revert to a previous configuration, set `configuration_content` back to the older content and apply. The provider creates a new version with that content (it does **not** resurrect the old version GUID — version numbers are never reused). `latest_version_entity_id` will point at the new version.

### Renaming a configuration

`name` is **immutable**. Changing it forces resource recreation: Terraform will destroy the existing configuration and create a new one with the new name. Plan output will explicitly show `# forces replacement` for `name` so this is never a surprise at apply time.

## Step 3 — Assign member entities to the fleet

`newrelic_fleet_members` controls which entities (hosts or Kubernetes clusters) belong to the fleet and which **ring** they are in. Rings are deployment tiers — a common pattern is `default` for the bulk of your fleet and `canary` for a small set of early-adopter nodes.

```hcl
resource "newrelic_fleet_members" "prod_linux" {
  fleet_id = newrelic_fleet.prod_linux.id

  ring {
    name       = "default"
    entity_ids = [
      "MXxOR0VQfEhPU1R8MTIzNDU2Nzg",   # host-1
      "MXxOR0VQfEhPU1R8ODc2NTQzMjE",   # host-2
      "MXxOR0VQfEhPU1R8OTk5ODc2NTQ",   # host-3
    ]
  }
}
```

### Multiple rings

Split your fleet into rings to stage rollouts — deploy to `canary` first, verify, then expand to `default`:

```hcl
resource "newrelic_fleet_members" "prod_linux" {
  fleet_id = newrelic_fleet.prod_linux.id

  ring {
    name       = "canary"
    entity_ids = [
      "MXxOR0VQfEhPU1R8MTIzNDU2Nzg",   # canary host
    ]
  }

  ring {
    name       = "default"
    entity_ids = [
      "MXxOR0VQfEhPU1R8ODc2NTQzMjE",
      "MXxOR0VQfEhPU1R8OTk5ODc2NTQ",
    ]
  }
}
```

### Opt-in management model

This resource uses an opt-in model: it tracks only the entities explicitly listed in `entity_ids`. Entities that joined the fleet through Agent Control instrumentation or any other out-of-band means are invisible to Terraform and are never removed. To bring such an entity under Terraform management, add its GUID to `entity_ids` — Terraform will adopt it into state on the next `apply`.

### Drift detection

On every `plan`, Terraform compares the declared `entity_ids` against the live membership returned by the API. If an entity was removed from a ring outside of Terraform, the plan will flag it and the next `apply` will re-add it. To accept the removal instead, delete the GUID from `entity_ids`.

## Step 4 — Create a deployment

A deployment triggers the actual rollout: it tells the fleet which agent version and which configuration version to push to its members.

```hcl
resource "newrelic_fleet_deployment" "infra_v1" {
  fleet_id    = newrelic_fleet.prod_linux.id
  name        = "Infra v1.58.0 — initial rollout"
  description = "Rolls out NRInfra 1.58.0 with the production configuration"

  agent {
    agent_type               = "NRInfra"
    version                  = "1.58.0"
    configuration_version_id = newrelic_fleet_configuration.infra.latest_version_entity_id
  }
}
```

A deployment may also be created or updated with **zero** `agent` blocks — useful to drain all agent assignments from a deployment, or to seed a deployment record before adding agents while it's still in `CREATED` phase. Within a deployment, each `agent_type` may appear at most once.

### Pinning to a specific configuration version

Referencing `latest_version_entity_id` ties the deployment to whichever version of the configuration is current at plan time. If you later edit the configuration's `configuration_content`, `latest_version_entity_id` changes — and any deployment using that reference will show a planned update. For deployments that have already started executing (phase != `CREATED`), that planned update will be **blocked by the phase-gate** described below.

For long-lived deployments that should remain stable, pin to a specific historical version GUID via `version_entity_ids[N]` instead:

```hcl
resource "newrelic_fleet_deployment" "stable" {
  fleet_id = newrelic_fleet.prod_linux.id
  name     = "Stable v1 rollout"

  agent {
    agent_type               = "NRInfra"
    version                  = "1.58.0"
    # First-ever version — won't drift when the config is later updated.
    configuration_version_id = newrelic_fleet_configuration.infra.version_entity_ids[0]
  }
}
```

### Deployment lifecycle and immutability

A deployment can only be modified while its `phase` is `CREATED`. Once the Fleet Control backend begins executing it (`IN_PROGRESS`, `COMPLETED`, or `FAILED`), the `name`, `description`, `agent`, and `tags` fields are locked. Any `plan` that proposes to change a locked field is rejected at plan time with a clear error.

To recover from a stuck plan after the deployment has executed:

1. **Preferred:** drop the executed deployment from Terraform state without touching the API:
   ```shell
   terraform state rm newrelic_fleet_deployment.infra_v1
   ```
   Or use a [`removed` block](https://developer.hashicorp.com/terraform/language/state/remove). Then re-declare a fresh deployment with the desired configuration.

2. `terraform destroy` works only if your HCL has no pending changes for the deployment — otherwise the same plan-time gate fires during the destroy plan.

To issue a follow-up rollout (e.g. roll out `v1.59.0`), declare a new deployment resource. The previous one stays in the API as a historical record:

```hcl
resource "newrelic_fleet_deployment" "infra_v2" {
  fleet_id    = newrelic_fleet.prod_linux.id
  name        = "Infra v1.59.0 — follow-up rollout"

  agent {
    agent_type               = "NRInfra"
    version                  = "1.59.0"
    configuration_version_id = newrelic_fleet_configuration.infra.latest_version_entity_id
  }
}
```

The `phase` attribute is exported from the resource and can be referenced in outputs or used with `terraform state show` to inspect the current status of a deployment.

## Using data sources

### Inspect live fleet membership

`data.newrelic_fleet_members` returns the entities currently in the fleet — useful for auditing or feeding member GUIDs into other resources.

```hcl
# All members across all rings.
data "newrelic_fleet_members" "all" {
  fleet_id   = newrelic_fleet.prod_linux.id
  depends_on = [newrelic_fleet_members.prod_linux]
}

# Members in a specific ring only.
data "newrelic_fleet_members" "canary_only" {
  fleet_id   = newrelic_fleet.prod_linux.id
  ring       = "canary"
  depends_on = [newrelic_fleet_members.prod_linux]
}

output "all_member_count" {
  value = length(data.newrelic_fleet_members.all.members)
}
```

### Read configuration content and version GUIDs

`data.newrelic_fleet_configuration` lets you look up a configuration that was created outside of the current Terraform workspace — for example, in a shared infrastructure module.

```hcl
# Look up by name.
data "newrelic_fleet_configuration" "shared_infra" {
  name = "Production Infra Config"
}

# Reference the latest version GUID in a deployment without owning the config resource.
resource "newrelic_fleet_deployment" "rollout" {
  fleet_id = newrelic_fleet.prod_linux.id
  name     = "Consume shared config"

  agent {
    agent_type               = "NRInfra"
    version                  = "1.58.0"
    configuration_version_id = data.newrelic_fleet_configuration.shared_infra.latest_version_entity_id
  }
}
```

## Putting it all together

The following is a self-contained configuration that combines all four resources into a complete fleet setup:

```hcl
# ── Fleet ───────────────────────────────────────────────────────────────────

resource "newrelic_fleet" "prod_linux" {
  name                = "Production Linux Hosts"
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
  description         = "Production fleet for Linux infrastructure agents"
  tags                = ["environment:production", "team:platform"]
}

# ── Agent configuration ──────────────────────────────────────────────────────

resource "newrelic_fleet_configuration" "infra" {
  name                  = "Production Infra Config"
  agent_type            = "NRInfra"
  managed_entity_type   = "HOST"
  operating_system      = "LINUX"
  configuration_content = file("${path.module}/configs/infra.yaml")
}

# ── Fleet membership ─────────────────────────────────────────────────────────

resource "newrelic_fleet_members" "prod_linux" {
  fleet_id = newrelic_fleet.prod_linux.id

  ring {
    name       = "canary"
    entity_ids = ["MXxOR0VQfEhPU1R8MTIzNDU2Nzg"]
  }

  ring {
    name       = "default"
    entity_ids = [
      "MXxOR0VQfEhPU1R8ODc2NTQzMjE",
      "MXxOR0VQfEhPU1R8OTk5ODc2NTQ",
    ]
  }
}

# ── Deployment ───────────────────────────────────────────────────────────────

resource "newrelic_fleet_deployment" "infra_v1" {
  fleet_id    = newrelic_fleet.prod_linux.id
  name        = "Infra v1.58.0 — initial rollout"
  description = "Rolls out NRInfra 1.58.0 to all production Linux hosts"

  agent {
    agent_type               = "NRInfra"
    version                  = "1.58.0"
    configuration_version_id = newrelic_fleet_configuration.infra.latest_version_entity_id
  }

  depends_on = [newrelic_fleet_members.prod_linux]
}

# ── Outputs ──────────────────────────────────────────────────────────────────

output "fleet_id" {
  description = "GUID of the production fleet."
  value       = newrelic_fleet.prod_linux.id
}

output "config_latest_version" {
  description = "Version number of the most recently created configuration version."
  value       = newrelic_fleet_configuration.infra.latest_version_number
}

output "deployment_phase" {
  description = "Current phase of the deployment (CREATED, IN_PROGRESS, COMPLETED, FAILED)."
  value       = newrelic_fleet_deployment.infra_v1.phase
}
```

Run `terraform apply` to provision all resources. Terraform resolves the dependency order automatically — the fleet is created first, followed by the configuration and members in parallel, and the deployment last (enforced by `depends_on` to ensure members are registered before the rollout begins).

## Importing existing resources

If you have fleet resources that were created outside of Terraform, you can bring them under management using `terraform import`.

```shell
# Fleet — plain fleet GUID
terraform import newrelic_fleet.prod_linux <fleet_guid>

# Fleet configuration — composite ID: <config_guid>:<managed_entity_type>
# managed_entity_type must be supplied because the API does not return it on lookup.
terraform import newrelic_fleet_configuration.infra <config_guid>:HOST

# Fleet members — plain fleet GUID (imports all current members into a single "default" ring)
terraform import newrelic_fleet_members.prod_linux <fleet_guid>

# Fleet deployment — plain deployment GUID
terraform import newrelic_fleet_deployment.infra_v1 <deployment_guid>
```

After importing fleet members, review the `ring` blocks that Terraform generates and update them to reflect your actual ring topology before running `terraform plan`.

## Next steps

- [newrelic_fleet](/providers/newrelic/newrelic/latest/docs/resources/fleet) — full argument reference for fleet creation
- [newrelic_fleet_configuration](/providers/newrelic/newrelic/latest/docs/resources/fleet_configuration) — flat `configuration_content` model, version numbering, rollback, and out-of-band drift warnings
- [newrelic_fleet_members](/providers/newrelic/newrelic/latest/docs/resources/fleet_members) — opt-in management model, adoption of Agent Control entities, and multi-ring moves
- [newrelic_fleet_deployment](/providers/newrelic/newrelic/latest/docs/resources/fleet_deployment) — deployment lifecycle, phase-gate immutability, and multi-agent deployments
- [data.newrelic_fleet_configuration](/providers/newrelic/newrelic/latest/docs/data-sources/fleet_configuration) — three lookup modes (by GUID, version GUID, or name)
- [data.newrelic_fleet_members](/providers/newrelic/newrelic/latest/docs/data-sources/fleet_members) — unfiltered and ring-filtered membership reads
