# OCI Logs Integration Module

Ships OCI logs to New Relic. Creates (or reuses) a VCN, a log-forwarder Function, a private Container Repository the Function's image is mirrored into, and one Service Connector Hub per entry in `connector_hub_details` to route log sources to that Function.

## Where this module fits in the full integration

A complete New Relic + OCI integration via Terraform composes four modules plus one resource. This module is **step 4b** (logs):

| Order | Component | What it provisions | Required? |
|-------|-----------|--------------------|-----------|
| 1 | [`wif-setup`](../wif-setup) | OAuth apps, IAM, identity propagation trust — the auth foundation | Yes |
| 2 | [`policy-setup`](../policy-setup) | IAM policies + (optionally) Vault secrets for NR keys | Yes |
| 3 | [`newrelic_cloud_oci_link_account`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/cloud_oci_link_account) | Registers the OCI tenancy with New Relic | Yes |
| 4a | [`metrics-integration`](../metrics-integration) | Service Connector Hub + Function for metrics shipping | Optional (need metrics) |
| 4b | `logs-integration` (this module) | Service Connector Hub + Function for logs shipping | Optional (need logs) |
| — | [`hub-calculator`](../hub-calculator) | Script that computes `connector_hubs_data`-style payloads across multiple compartments | Recommended for step 4a/4b |

See the [Cloud integrations guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/cloud_integrations_guide#oci) for the full end-to-end example wiring all four modules together.

## Inputs

| Variable | Type | Default | Required | Notes |
|----------|------|---------|----------|-------|
| `tenancy_ocid` | string | — | yes | Your OCI tenancy OCID |
| `compartment_ocid` | string | — | yes | Compartment where resources are created. Use the compartment OCID already created by `policy-setup` |
| `region` | string | — | yes | OCI region to deploy into |
| `newrelic_logging_identifier` | string | `"logs"` | no | Label used in resource names; leave blank if not needed |
| `create_vcn` | bool | `true` | no | Set `false` to use an existing subnet via `function_subnet_id` instead |
| `function_subnet_id` | string | — | only if `create_vcn = false` | Existing subnet OCID for the Function app |
| `connector_hub_details` | string (JSON) | `null` | no | JSON array of connectors to create — see [Connector Hub format](#connector-hub-format) below |
| `batch_size_in_kbs` | number | `6000` | no | Max batch size per Service Connector Hub delivery (max `6000`) |
| `batch_time_in_sec` | number | `60` | no | Max wait before flushing a batch (max `300`) |
| `user_ocid` | string | — | yes | OCID of the user Terraform authenticates as; used for the Container Registry push token |
| `function_image` | string | `"docker.io/newrelic/oci-log-forwarder:latest"` | no | Public source image, mirrored into a private repo in your tenancy. Re-apply after a new push under the same tag to pick up the update |
| `registry_username` | string | `""` | no | Container Registry username (no tenancy namespace). Leave empty to derive from `user_ocid` |
| `registry_auth_token` | string (sensitive) | `""` | no | Existing registry auth token. Leave empty to let the module create one (OCI allows 2 per user) |
| `debug_enabled` | string | `"FALSE"` | no | Enables Function debug mode |
| `new_relic_region` | string | `"US"` | no | `US`, `EU`, or `JP` |
| `secret_ocid` | string | — | yes | Vault Secret OCID holding the New Relic license key. Use the one already created by `policy-setup` |
| `user_api_secret_ocid` | string | — | yes | Vault Secret OCID holding the New Relic user API key |
| `newrelic_account_id` | string (sensitive) | — | yes | New Relic account ID to send logs to |
| `provider_account_id` | string (sensitive) | — | yes | Linked OCI provider account ID in New Relic |
| `metrics_tier` | string | `"basic"` | no | Tier of `forwarder.*` self-observability metrics the Function emits about its own health, in addition to shipping your logs: `none`, `basic` (invocations, records received/delivered/dropped, delivery duration, pipeline lag), or `advanced` (adds byte volumes, decode/serialize errors, batching behavior, delivery error classes, run duration, secret-fetch failures, client-cache hit rate). Billed on ingest like any other metric. See [OCI log forwarder observability](https://docs.newrelic.com/docs/logs/forward-logs/oci-log-forwarder-observability/) |

`TENANCY_NAME` and `COMPARTMENT_NAME` are also set on the Function automatically (not separate inputs) — human-readable labels, derived via `data.oci_identity_tenancy`/`data.oci_identity_compartment`, used to tag the `forwarder.*` metrics above so dashboards can facet by tenancy/compartment rather than by raw OCID. For a tenancy deployed at its root compartment, `COMPARTMENT_NAME` falls back to the tenancy's own name, since OCI's `GetCompartment` API can't resolve the root compartment by OCID.

### Connector Hub format

`connector_hub_details` is a JSON-encoded array; each entry becomes one `oci_sch_service_connector`:

```json
[
  {
    "display_name": "newrelic-logs-connector",
    "description": "Routes application logs to New Relic",
    "log_sources": [
      { "compartment_id": "ocid1.compartment.oc1..aaaaaaaa...", "log_group_id": "ocid1.loggroup.oc1.iad..aaaaaaaa..." }
    ]
  }
]
```

Leave as `null` (the default) to skip creating any connectors — useful if you only want the Function/registry infrastructure provisioned ahead of a separate step that wires up log sources.

## Outputs

```hcl
vcn_network_details = {
  vcn_id             = "ocid1.vcn.oc1..."          # empty string if create_vcn = false
  nat_gateway_id      = "ocid1.natgateway.oc1..."
  nat_route_id        = "ocid1.route.oc1..."
  service_gateway_id  = "ocid1.servicegateway.oc1..."
  sgw_route_id        = "ocid1.route.oc1..."
  subnet_id           = "ocid1.subnet.oc1..."       # = var.function_subnet_id if create_vcn = false
}
```

## What gets created

```
┌──────────────────────────────────────────────────────────┐
│  Your OCI Compartment                                     │
│                                                            │
│  VCN + subnet + NAT/service/internet gateways (optional,  │
│  skipped if create_vcn = false)                            │
│                                                            │
│  ┌────────────────────────┐   ┌─────────────────────────┐ │
│  │ Container Repository   │◄──│ function_image mirrored │ │
│  │ (private)               │   │ in via image_mirror.py  │ │
│  └────────────────────────┘   └─────────────────────────┘ │
│              │                                             │
│              ▼                                             │
│  ┌────────────────────────┐                                │
│  │ OCI Function App +     │  env: FORWARDER_METRICS_TIER,  │
│  │ Function (log forwarder)│      TENANCY_NAME,            │
│  └────────────────────────┘      COMPARTMENT_NAME          │
│              ▲                                             │
│              │                                             │
│  ┌────────────────────────┐                                │
│  │ Service Connector Hub  │  one per connector_hub_details │
│  │ (one or more)           │  entry; routes log_sources to  │
│  └────────────────────────┘  the Function above             │
└──────────────────────────────────────────────────────────┘
```

After the Function and connectors are created, a `null_resource` calls New Relic's `cloudUpdateAccount` GraphQL mutation to finish linking the account.

## Calling pattern

```hcl
module "newrelic_oci_logs" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/cloud-integrations/oci/logs-integration"

  tenancy_ocid      = var.tenancy_ocid
  compartment_ocid  = var.compartment_ocid
  region            = var.region
  user_ocid         = var.user_ocid
  secret_ocid       = var.secret_ocid
  user_api_secret_ocid = var.user_api_secret_ocid
  newrelic_account_id  = var.newrelic_account_id
  provider_account_id  = var.provider_account_id

  connector_hub_details = jsonencode([{
    display_name = "newrelic-logs-connector"
    description  = "Routes application logs to New Relic"
    log_sources = [{
      compartment_id = var.compartment_ocid
      log_group_id   = var.log_group_id
    }]
  }])

  # Opt in to forwarder self-observability metrics (defaults to "basic")
  metrics_tier = "basic"
}
```
