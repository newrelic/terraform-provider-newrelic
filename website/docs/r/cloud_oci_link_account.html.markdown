---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_oci_link_account"
sidebar_current: "docs-newrelic-cloud-resource-oci-link-account"
description: |-
  Link an Oracle Cloud Infrastructure (OCI) account to New Relic.
---
# Resource: newrelic_cloud_oci_link_account

Use this resource to link an Oracle Cloud Infrastructure (OCI) account to New Relic.

## Prerequisite

To link an OCI account to New Relic, you need an Oracle Cloud Infrastructure tenancy with appropriate permissions. OCI provides enterprise-grade cloud infrastructure services including compute, storage, networking, and database services across multiple regions globally.

## Example Usage

```hcl
resource "newrelic_cloud_oci_link_account" "foo" {
  account_id = 1234567
  name       = "My New Relic - OCI Linked Account"
  tenant_id  = "ocid1.tenancy.oc1..aaaaaaaaexample"
}
```

## Argument Reference

The following arguments are supported:

- `account_id` - (Optional) The New Relic account ID to operate on. This allows the user to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`, if not specified in the configuration.
- `name` - (Required) The name/identifier of the OCI - New Relic 'linked' account.
- `tenant_id` - (Required) The Oracle Cloud Infrastructure (OCI) tenancy OCID.

-> **NOTE:** Altering the `account_id` of an already applied `newrelic_cloud_oci_link_account` resource shall trigger a recreation of the resource, instead of an update.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the OCI linked account.

## Import

Linked OCI accounts can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_cloud_oci_link_account.foo <id>
```
