---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_oci_link_account"
sidebar_current: "docs-newrelic-cloud-resource-oci-link-account"
description: |-
  Link an Oracle Cloud Infrastructure (OCI) account to New Relic.
---
# Resource: newrelic_cloud_oci_link_account

Use this resource to link an Oracle Cloud Infrastructure (OCI) account to New Relic.

This setup is used to create a provider account with OCI credentials, establishing a relationship between Oracle and New Relic. Additionally, as part of this integration, we store WIF (Workload Identity Federation) credentials which are further used for fetching data and validations, and vault OCIDs corresponding to the vault resource where the New Relic ingest and user keys are stored in the OCI console.

## Prerequisites

For the `newrelic_cloud_oci_link_account` resource to work properly, you need an OCI tenancy with IAM permissions to create and manage the identity artifacts (client/application, secrets, compartments, and service user) referenced below. OCI provides enterprise-grade cloud services across multiple global regions.

> NOTE: Before using this resource, ensure the New Relic provider is configured with valid credentials.  
> See Getting Started: [New Relic provider guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/getting_started)

If you encounter issues or bugs, please [open an issue in the GitHub repository](https://github.com/newrelic/terraform-provider-newrelic/issues/new/choose).

### Workload Identity Federation (WIF) Attributes

The following arguments rely on an OCI Identity Domain OAuth2 client set up for workload identity federation (identity propagation): `oci_client_id`, `oci_client_secret` and `oci_domain_url`.

To create and retrieve these values, follow Oracle's guidance for configuring identity propagation / JWT token exchange:

[Oracle documentation: Create an identity propagation trust (JWT token exchange)](https://docs.oracle.com/en-us/iaas/Content/Identity/api-getstarted/json_web_token_exchange.htm#jwt_token_exchange__create-identity-propagation-trust)

WIF configuration steps:
1. Create (or identify) an Identity Domain and register an OAuth2 confidential application (client) to represent New Relic ingestion.
2. Generate / record the client ID (`oci_client_id`) and client secret (`oci_client_secret`). Store the secret securely (e.g., in OCI Vault; reference its OCID via `ingest_vault_ocid` / `user_vault_ocid` if desired).
3. Use the Identity Domain base URL as `oci_domain_url` (format: `https://idcs-<hash>.identity.oraclecloud.com`).
4. Ensure the client has the required scopes and the tenancy policies allow the token exchange.

> TIP: Rotating the OAuth2 client secret only requires updating `oci_client_secret`; it does not force resource replacement.

## Example Usage

Minimal example (required arguments for creation):

```hcl
resource "newrelic_cloud_oci_link_account" "example" {
  # Optional if set via the provider block or NEW_RELIC_ACCOUNT_ID environment variable
  account_id        = 1234567

  # Changing this forces replacement (ForceNew)
  tenant_id         = "ocid1.tenancy.oc1..aaaaaaaaexample"

  name              = "my-oci-link"
  compartment_ocid  = "ocid1.compartment.oc1..bbbbbbbbexample"
  oci_client_id     = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"                     # OCI Identity Domain OAuth2 client (WIF)
  oci_client_secret = var.oci_client_secret                                        # Sensitive
  oci_domain_url    = "https://idcs-1234567890abcdef.identity.oraclecloud.com"    # Identity domain base URL
  oci_home_region   = "us-ashburn-1"
}
```

Example including optional secret references and update-only fields:

```hcl
resource "newrelic_cloud_oci_link_account" "full" {
  name              = "my-oci-link-full"
  tenant_id         = "ocid1.tenancy.oc1..aaaaaaaaexample"
  compartment_ocid  = "ocid1.compartment.oc1..bbbbbbbbexample"
  oci_client_id     = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  oci_client_secret = var.oci_client_secret
  oci_domain_url    = "https://idcs-1234567890abcdef.identity.oraclecloud.com"
  oci_home_region   = "us-ashburn-1"

  # Vault secret OCIDs (these may point to secrets that store rotated values)
  ingest_vault_ocid = "ocid1.vaultsecret.oc1..ccccccccexample"
  user_vault_ocid   = "ocid1.vaultsecret.oc1..ddddddddexample"

  # Integration configuration
  instrumentation_type = "METRICS,LOGS"

  # Update-only fields (ignored during initial create, applied on update)
  oci_region        = "us-phoenix-1"
  metric_stack_ocid = "ocid1.stack.oc1..eeeeeeeeexample"
  logging_stack_ocid = "ocid1.stack.oc1..ffffffloggingstack"
}
```

## Argument Reference

The following arguments are supported (current provider schema):

- `account_id` - (Optional, ForceNew) New Relic account to operate on. Overrides the provider-level `account_id`. If omitted, use the provider default or `NEW_RELIC_ACCOUNT_ID`.
- `tenant_id` - (Required, ForceNew) OCI tenancy OCID (root tenancy). Changing forces a new linked account.
- `name` - (Required) Display name for the linked account.
- `compartment_ocid` - (Required) OCI compartment OCID representing (or containing) the monitored resources/newrelic compartment.
- `oci_client_id` - (Required) OCI Identity Domain (IDCS) OAuth2 client ID used for workload identity federation.
- `oci_client_secret` - (Required, Sensitive) OAuth2 client secret. Not displayed in plans or state outputs.
- `oci_domain_url` - (Required) Base URL of the OCI Identity Domain (e.g. `https://idcs-<hash>.identity.oraclecloud.com`).
- `oci_home_region` - (Required) Home region of the tenancy (e.g. `us-ashburn-1`).
- `ingest_vault_ocid` - (Required) Vault secret OCID containing an ingest secret.
- `user_vault_ocid` - (Required) Vault secret OCID containing a user or auxiliary secret.
- `instrumentation_type` - (Optional) Specifies the type of integration, such as metrics, logs, or a combination of logs and metrics (e.g., `METRICS`, `LOGS`, `METRICS,LOGS`).
- `oci_region` - (Optional, Update-only) OCI region for the linkage (ignored on create, applied on update).
- `metric_stack_ocid` - (Optional, Update-only) Metric stack OCID (ignored on create, applied on update).
- `logging_stack_ocid` - (Optional) The Logging stack identifier for the OCI account.

### ForceNew & Update-only Behavior

- Changing `account_id` or `tenant_id` forces resource replacement.
- Update-only fields (`oci_region`, `metric_stack_ocid`) are ignored at initial creation and only sent on update operations.

### Sensitive Data Handling

- `oci_client_secret` is stored as a sensitive value in state and excluded from plan/apply output. Rotate as needed and re-apply to update; this performs an in-place update (no replacement) unless another ForceNew attribute changed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the OCI linked account.

> NOTE: Only a subset of arguments may currently be returned by the read operation (`account_id`, `tenant_id`, `name`). Other write-only, sensitive, or create-time fields may not round-trip during `terraform refresh` or `terraform plan` until backend API read support is expanded. This is expected.

## Import

Linked OCI accounts can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_cloud_oci_link_account.foo <id>
```
