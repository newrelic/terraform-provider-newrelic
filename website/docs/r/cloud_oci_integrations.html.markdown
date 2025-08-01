---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_oci_integrations"
sidebar_current: "docs-newrelic-resource-cloud-oci-integrations"
description: |-
Integrate OCI services with New Relic.
---

# Resource: newrelic\_cloud\_oci\_integrations

Use this resource to integrate Oracle Cloud Infrastructure (OCI) services with New Relic.

## Prerequisite

Setup is required for this resource to work properly. This resource assumes you have [linked an OCI account](cloud_oci_link_account.html) to New Relic and configured it to pull metrics from OCI.

New Relic doesn't automatically receive metrics from OCI services, so this resource can be used to configure integrations to those services.

## Example Usage

Leave an integration block empty to use its default configuration. You can also use the [full example, including the OCI set up, found in our guides](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/cloud_integrations_guide#oci).

```hcl
resource "newrelic_cloud_oci_link_account" "foo" {
  name      = "example"
  tenant_id = "ocid1.tenancy.oc1..aaaaaaaaexample"
}

resource "newrelic_cloud_oci_integrations" "foo1" {
  linked_account_id = newrelic_cloud_oci_link_account.foo.id
  oci_metadata_and_tags {}
}
```
## Argument Reference

-> **WARNING:** Starting with [v3.27.2](https://registry.terraform.io/providers/newrelic/newrelic/3.27.2) of the New Relic Terraform Provider, updating the `linked_account_id` of a `newrelic_cloud_oci_integrations` resource that has been applied would **force a replacement** of the resource (destruction of the resource, followed by the creation of a new resource). When such an update is performed, please carefully review the output of `terraform plan`, which would clearly indicate a replacement of this resource, before performing a `terraform apply`.

* `account_id` - (Optional) The New Relic account ID to operate on.  This allows the user to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.
* `linked_account_id` - (Required) The ID of the linked OCI account in New Relic.

The following arguments/integration blocks are intended to be used.

* `oci_metadata_and_tags` - OCI metadata and tags integration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the OCI linked account.

## Import

Linked OCI account integrations can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_cloud_oci_integrations.foo <id>
```
