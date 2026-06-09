# OCI Workload Identity Federation Setup Module

Sets up the **authentication foundation** New Relic needs to securely access your OCI tenancy. This is one of four composable modules that together build a complete OCI integration; this one creates:

- 2 Confidential OAuth applications (admin + token exchange)
- IAM resources (group + policy) — shape varies by trust type
- Service user (UPST only)
- Identity Propagation Trust (UPST or RPST)

## Where this module fits in the full integration

A complete New Relic + OCI integration via Terraform composes four modules plus one resource. This module is **step 1**:

| Order | Component | What it provisions | Required? |
|-------|-----------|--------------------|-----------|
| 1 | [`wif-setup`](.) (this module) | OAuth apps, IAM, identity propagation trust — the auth foundation | Yes |
| 2 | [`policy-setup`](../policy-setup) | IAM policies + (optionally) Vault secrets for NR keys | Yes |
| 3 | [`newrelic_cloud_oci_link_account`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/cloud_oci_link_account) | Registers the OCI tenancy with New Relic | Yes |
| 4a | [`metrics-integration`](../metrics-integration) | Service Connector Hub + Function for metrics shipping | Optional (need metrics) |
| 4b | [`logs-integration`](../logs-integration) | Service Connector Hub + Function for logs shipping | Optional (need logs) |

This module **must run first** — its OAuth credentials and IAM domain URL feed into both `policy-setup` and the `newrelic_cloud_oci_link_account` resource. See the [Cloud integrations guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/cloud_integrations_guide#oci) for the full end-to-end example wiring all four modules together.

## Trust types

| Type | When to use | What it creates |
|------|-------------|-----------------|
| **UPST** (default) | Most customers | Service user, IAM group, group-based IAM policy, trust with `subjectType=User` |
| **RPST** | Multi-account customers, customers hitting OCI's IAM policy statement limit, or anyone who wants claim-based authorization | No service user, claim-based IAM policy, trust with `subjectType=Resource` + `claimPropagations` + `impersonatingResource` |

Both flows are backward-compatible. Existing UPST customers can adopt RPST by setting `trust_type = "RPST"` in their tfvars; old setups continue to work unchanged.

## Inputs

| Variable | Type | Default | Required | Notes |
|----------|------|---------|----------|-------|
| `tenancy_ocid` | string | — | yes | Your OCI tenancy OCID |
| `identity_domain_name` | string | `"Default"` | no | Identity Domain to use |
| `newrelic_region` | string | `"US"` | yes | `US` or `EU` |
| `home_region` | string | — | yes | OCI home region (e.g. `us-ashburn-1`) |
| `fingerprint` | string | — | yes | OCI API key fingerprint |
| `private_key` | string (sensitive) | — | yes | OCI API key PEM content |
| `resource_prefix` | string | `"newrelic"` | no | Prefix for created resources |
| `compartment_id` | string | tenancy root | no | Compartment for IAM policy |
| `trust_type` | string | `"UPST"` | no | `"UPST"` or `"RPST"` |
| `newrelic_account_id` | string | `""` | required when `trust_type = "RPST"` | Used in the IAM policy claim match for RPST |
| `resource_tag` | string | `""` | no (RPST only) | When set, NR stamps `ext_resource_tag` on every JWT and the IAM policy scopes by it. Lets you bind one integration to a specific tag value (e.g. `prod`). |

## Outputs

```hcl
newrelic_integration_details = {
  iam_domain_url               = "https://idcs-...identity.oraclecloud.com"
  token_exchange_client_id     = "<oauth client id>"
  token_exchange_client_secret = "<oauth client secret>"  # sensitive
  trust_type                   = "UPST" | "RPST"
}

impersonating_resource = "newrelic-integration"   # paste this into your trust config (RPST only)
```

## Calling pattern

```hcl
provider "newrelic" {
  account_id = var.newrelic_account_id
  api_key    = var.newrelic_user_api_key
  region     = var.newrelic_region
}

provider "oci" {
  tenancy_ocid = var.tenancy_ocid
  region       = var.home_region
  fingerprint  = var.fingerprint
  private_key  = var.private_key
}

module "newrelic_wif" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/cloud-integrations/oci/wif-setup"

  tenancy_ocid         = var.tenancy_ocid
  identity_domain_name = var.identity_domain_name
  newrelic_region      = var.newrelic_region
  home_region          = var.home_region
  fingerprint          = var.fingerprint
  private_key          = var.private_key

  # NR-562518 (RPST): set both for RPST. Omit for default UPST behavior.
  trust_type          = var.trust_type
  newrelic_account_id = var.newrelic_account_id
  # Optional (RPST only): if set, scopes the IAM policy on `ext_resource_tag` too.
  resource_tag = var.resource_tag
}

resource "newrelic_cloud_oci_link_account" "link" {
  tenant_id         = var.tenancy_ocid
  name              = var.linked_account_name
  compartment_ocid  = var.compartment_ocid
  oci_client_id     = module.newrelic_wif.newrelic_integration_details.token_exchange_client_id
  oci_client_secret = module.newrelic_wif.newrelic_integration_details.token_exchange_client_secret
  oci_domain_url    = module.newrelic_wif.newrelic_integration_details.iam_domain_url
  oci_home_region   = var.home_region
  trust_type        = module.newrelic_wif.newrelic_integration_details.trust_type
  resource_tag      = module.newrelic_wif.newrelic_integration_details.resource_tag
}
```

See `terraform.tfvars.example` for a populated example.

## What gets created

```
                  ┌────────────────────────────────────┐
                  │  Your OCI Identity Domain          │
                  │                                    │
                  │  ┌──────────────────────────────┐  │
                  │  │ Admin App (high privilege)   │  │
                  │  │ Disable after first apply.   │  │
                  │  └──────────────────────────────┘  │
                  │  ┌──────────────────────────────┐  │
                  │  │ Token Exchange App           │  │
                  │  │ (no roles)                   │  │
                  │  └──────────────────────────────┘  │
                  │                                    │
                  │  IAM Policy:                       │
                  │   UPST: "Allow group X to..."      │
                  │   RPST: "where ext_account_id=..." │
                  │                                    │
                  │  Service User (UPST only)          │
                  │  IAM Group (UPST only)             │
                  │                                    │
                  │  Identity Propagation Trust:       │
                  │   UPST → subjectType=User          │
                  │   RPST → subjectType=Resource      │
                  │           + impersonatingResource  │
                  │           + claimPropagations      │
                  └────────────────────────────────────┘
```

## After `terraform apply`

For **RPST customers only**: Confirm the static `impersonating_resource` value matches what's in your trust config (it's the same `"newrelic-integration"` string everywhere — module output + trust config + NR's runtime).

For **all customers**: the `newrelic_cloud_oci_link_account` resource registers the integration with New Relic. After that, NR begins polling your tenancy on the configured intervals.

For **security**: disable the admin app in OCI after the first successful `terraform apply`. The token-exchange app (no privileges) is what NR uses at runtime — keep it active.
