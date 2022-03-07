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

Setup is required in Azure for this resource to work properly. New Relic Azure integration can be done by creating Azure account

Create a New Relic application and key in Azure

Grant this application access to the Azure services you want to monitor

To pull data from Azure instead, complete the [steps outlined here](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/get-started/activate-azure-integrations).

## Example Usage

```hcl

  resource "newrelic_cloud_azure_link_account" "foo"{
    account_id = "The New Relic account ID where you want to link the Azure account"
	application_id = "id of the application"
	client_secret_id = "secret value of clients Azure account"
	subscription_id = "%Subscription Id of Azure"
	tenant_id = "tenant id of the Azure"
	name  = "name of the linked account"
}
```

## Argument Reference

The following arguments are supported:
- `account_id` - (Required) - Account Id of the New Relic.
- `application_id` - (Required) - Application Id of the App.
- `client_secret_id` - (Required) - Secret Value of the client.
- `subscription_id` - (Required) - Subscription Id of the Azure cloud account.
- `tenant_id` - (Required) - Tenant Id of the Azure cloud account.
- `name` - (Required) - The name of the application in New Relic APM.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The application Id, subscription Id, clientsecret Id & tenant Id of the Azure linked account.

## Import

Linked Azure accounts can be imported using `id`, you can find the `id` of existing Azure linked accounts in Azure dashboard under Infrastructure in NewRelic

```bash
$ terraform import newrelic_cloud_azure_link_account.foo <id>

```
