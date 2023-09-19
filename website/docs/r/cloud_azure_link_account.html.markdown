---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_azure_link_account"
sidebar_current: "docs-newrelic-resource-cloud-azure-link-account"
description: |-
  Link an Azure account to New Relic.
---

# Resource: newrelic_cloud_azure_link_account

Use this resource to link an Azure account to New Relic.

## Prerequisite

Some configuration is required in Azure for the New Relic Azure cloud integrations to be able to pull data.

To start receiving Azure data with New Relic Azure integrations, connect your Azure account to New Relic infrastructure monitoring. If you don't have one already, create a New Relic account. It's free, forever.

Setup is required in Azure for this resource to work properly. You can find instructions on how to set up Azure on [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/get-started/activate-azure-integrations/).

## Example Usage

You can also use the [full example, including the Azure set up, found in our guides](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/cloud_integrations_guide#azure).

```hcl
resource "newrelic_cloud_azure_link_account" "foo"{
  account_id = "The New Relic account ID where you want to link the Azure account"
  application_id = "ID of the application"
  client_secret = "Secret value of client's Azure account"
  subscription_id = "Subscription ID of Azure"
  tenant_id = "Tenant ID of the Azure"
  name  = "Name of the linked account"
}
```

## Argument Reference

The following arguments are supported:

- `account_id` - (Required) - Account ID of the New Relic.
- `application_id` - (Required) - Application ID of the App.
- `client_secret` - (Required) - Secret Value of the client.
- `subscription_id` - (Required) - Subscription ID of the Azure cloud account.
- `tenant_id` - (Required) - Tenant ID of the Azure cloud account.
- `name` - (Required) - The name of the application in New Relic APM.

-> **WARNING:** Starting with [v3.27.2](https://registry.terraform.io/providers/newrelic/newrelic/3.27.2) of the New Relic Terraform Provider, updating any of the aforementioned attributes (except `name`) of a `newrelic_cloud_azure_link_account` resource that has been applied would **force a replacement** of the resource (destruction of the resource, followed by the creation of a new resource). Please carefully review the output of `terraform plan`, which would clearly indicate a replacement of this resource, before performing a `terraform apply`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the linked Azure account in New Relic.

## Import

Linked Azure accounts can be imported using `id`, you can find the `id` of existing Azure linked accounts in Azure dashboard under Infrastructure in NewRelic

```bash
$ terraform import newrelic_cloud_azure_link_account.foo <id>

```
