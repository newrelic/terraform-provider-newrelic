---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_gcp_link_account"
sidebar_current: "docs-new relic-resource-cloud-gcp-link-account"
description: |-
  Link an GCP account to New Relic.
---

# Resource: newrelic_cloud_gcp_link_account

Use this resource to link an GCP account to New Relic.

## Prerequisite

Setup is required in GCP for this resource to work properly. The New Relic GCP integration can be set up to pull metrics from GCP services.

Using a metric stream to New Relic is the preferred way to integrate with Azure. Follow the [steps outlined here](https://docs.newrelic.com/docs/infrastructure/google-cloud-platform-integrations/get-started/gcp-integration-metrics) to set up a metric stream.

To pull data from GCP instead, complete the [steps outlined here](https://docs.newrelic.com/docs/infrastructure/google-cloud-platform-integrations/get-started/connect-google-cloud-platform-services-new-relic).

## Example Usage

```hcl

 resource "newrelic_cloud_azure_link_account" "foo"{
   project_Id = "id of the Project"
   name  = "account name"
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) - project Id of the gcp account.
- `name` - (Required) - The name of the application in New Relic APM.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The project Id of the GCP linked account.

## import
