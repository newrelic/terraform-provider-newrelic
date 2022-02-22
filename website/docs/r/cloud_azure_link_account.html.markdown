---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_azure_link_account"
sidebar_current: "docs-newrelic-resource-cloud-azure-link-account"
description: |-
  Link an AWS account to New Relic.
---

# Resource: newrelic_cloud_azure_link_account

Use this resource to link an Azure account to New Relic.

## Prerequisite

Setup is required in Azure for this resource to work properly. The New Relic Azure integration can be set up to pull metrics from Azure services.

Using a metric stream to New Relic is the preferred way to integrate with Azure. Follow the [steps outlined here](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/get-started/azure-integration-metrics) to set up a metric stream.

To pull data from Azure instead, complete the [steps outlined here](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/get-started/activate-azure-integrations).

## Example Usage

```hcl

  resource "newrelic_cloud_azure_link_account" "foo"{
	application_id = "id of the application"
	client_secret_id = "secret value of clients Azure account"
	subscription_id = "%Subscription Id of Azure"
	tenant_id = "tenant id of the Azure"
	name  = "account name"
}
```

## Argument Reference

The following arguments are supported:

- `application_id` - (Required) - Application Id of the App.
- `client_secret_id` - (Required) - Secret Value of the client.
- `subscription_id` - (Required) - Subscription Id of the Azure cloud account.
- `tenant_id` - (Required) - Tenant Id of the Azure cloud account.
- `name` - (Required) - The name of the application in New Relic APM.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The application Id, subscription Id, clientsecret Id & tenant Id of the Azure linked account.

## Import
