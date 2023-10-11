---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_gcp_link_account"
sidebar_current: "docs-new relic-resource-cloud-gcp-link-account"
description: |-
Link a GCP account to New Relic.
---

# Resource: newrelic_cloud_gcp_link_account

Use this resource to link a GCP account to New Relic.

## Prerequisite

To start receiving Google Cloud Platform (GCP) data with New Relic GCP integrations, connect your Google project to New Relic infrastructure monitoring. If you don't have one already, create a New Relic account. It's free, forever.

Setup is required in GCP for this resource to work properly. The New Relic GCP integration can be done by creating a user account or a service account.

A user with Project IAM Admin role is needed to add the service account ID as a member in your GCP project.

In the GCP project IAM & admin, the service account must have the Project Viewer role and the Service Usage Consumer role or, alternatively, a custom role.

Follow the [steps outlined here](https://docs.newrelic.com/docs/infrastructure/google-cloud-platform-integrations/get-started/connect-google-cloud-platform-services-new-relic) to set up the integration.

## Example Usage

You can also use the [full example, including the GCP set up, found in our guides](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/cloud_integrations_guide#gcp).

```hcl
resource "newrelic_cloud_gcp_link_account" "foo" {
  account_id = "account id of newrelic account"
  project_id = "id of the Project"
  name  = "account name"
}
```

## Argument Reference

The following arguments are supported:

- `account_id` - (Optional) - Account ID of the New Relic account.
- `project_id` - (Required) - Project ID of the GCP account.
- `name` - (Required) - The name of the GCP account in New Relic.

-> **WARNING:** Starting with [v3.27.2](https://registry.terraform.io/providers/newrelic/newrelic/3.27.2) of the New Relic Terraform Provider, updating any of the aforementioned attributes (except `name`) of a `newrelic_cloud_gcp_link_account` resource that has been applied would **force a replacement** of the resource (destruction of the resource, followed by the creation of a new resource). Please carefully review the output of `terraform plan`, which would clearly indicate a replacement of this resource, before performing a `terraform apply`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The id of the GCP linked account.

## Import

Linked GCP accounts can be imported using `id`, you can find the `id` of an existing GCP linked accounts in GCP dashboard under Infrastructure in Newrelic Console.

```bash

  $  terraform import newrelic_cloud_gcp_link_account.foo <id>

```
