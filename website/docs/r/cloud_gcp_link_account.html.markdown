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

```hcl
 
 resource "newrelic_cloud_gcp_link_account" "foo"{
   account_id = "account id of newrelic account"
   project_id = "id of the Project"
   name  = "account name"
}
```

## Argument Reference

The following arguments are supported:

- `account_id` - (Required) - account id of the newrelic account.
- `project_id` - (Required) - project id of the gcp account.
- `name` - (Required) - The name of the application in New Relic APM.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The id of the GCP linked account.

## Import

Linked GCP accounts can be imported using `id`, you can find the `id` of an existing GCP linked accounts in GCP dashboard under Infrastructure in Newrelic Console.

```bash

  $  terraform import newrelic_cloud_gcp_link_account.foo <id>

```
