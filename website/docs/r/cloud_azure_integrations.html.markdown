---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_azure_integrations"
sidebar_current: "docs-newrelic-resource-cloud-azure-integrations"
description: |-
Integrate Azure services with New Relic.
---

# Resource: newrelic\_cloud\_azure\_integrations

Use this resource to integrate Azure services with New Relic.

## Prerequisite

To start receiving Azure data with New Relic Azure integrations, connect your Azure account to New Relic infrastructure monitoring. If you don't have one already, create a New Relic account. It's free, forever.

Setup is required for this resource to work properly. This resource assumes you have [linked an Azure account](cloud_azure_link_account.html) to New Relic.

You can find instructions on how to set up Azure on [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/get-started/activate-azure-integrations/).

## Example Usage

```hcl

  resource "newrelic_cloud_azure_link_account" "foo"{
    account_id = "The New Relic account ID where you want to link the Azure account"
	application_id = "ID of the application"
	client_secret = "Secret value of client's Azure account"
	subscription_id = "Subscription ID of Azure"
	tenant_id = "Tenant ID of the Azure"
	name  = "Name of the linked account"
}

  resource "newrelic_cloud_azure_integrations" "foo" {
    linked_account_id = newrelic_cloud_azure_link_account.foo.id
    account_id = "The New Relic account ID"
    azure_api_management {
      metrics_polling_interval = 1200
      resource_groups = "beyond"
    }

    azure_app_gateway {
      metrics_polling_interval = 1200
      resource_groups = "beyond"
    }
    azure_app_service {
      metrics_polling_interval = 1200
      resource_groups = "beyond"
    }
    azure_containers {
      metrics_polling_interval = 1200
      resource_groups = "beyond"
    }

    azure_cosmos_db {
      metrics_polling_interval = 1200
      resource_groups = "beyond"
    }
    azure_cost_management {
      metrics_polling_interval = 1200
      tag_keys = ""
    }

    azure_data_factory {
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }

    azure_event_hub {
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }

    azure_express_route {
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }
    azure_event_hub {
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }
    azure_firewalls {
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }
    azure_front_door {
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }
    azure_functions {
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }
    azure_key_vault {
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }
    azure_load_balancer {
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }
    azure_logic_apps {
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }
    azure_machine_learning {
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }
    azure_maria_db {
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }

    azure_mysql {
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }
    azure_postgresql{
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }
    azure_power_bi_dedicated{
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }
    
    azure_redis_cache{
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }

    azure_service_bus{
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }

    azure_sql{
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }

    azure_storage{
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }

    azure_virtual_machine{
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }

    azure_virtual_networks{
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }

    azure_vms{
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }

    azure_vpn_gateway{
      metrics_polling_interval = 1200
      resource_groups = "beyond"

    }
}
```
## Argument Reference


The following arguments are supported:

* `account_id` - (Optional) The New Relic account ID to operate on.  This allows the user to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.
* `linked_account_id` - (Required) The ID of the linked Azure account in New Relic.
* `azure_api_management` - (Optional) Azure API Management. See [Integration blocks](#integration-blocks) below for details.
* `azure_app_gateway` - (Optional) Azure App Gateway. See [Integration blocks](#integration-blocks) below for details.
* `azure_app_service` - (Optional) Azure App Service. See [Integration blocks](#integration-blocks) below for details.
* `azure_containers` - (Optional) Azure Containers. See [Integration blocks](#integration-blocks) below for details.
* `azure_cosmos_db` - (Optional) Azure CosmosDB. See [Integration blocks](#integration-blocks) below for details.
* `azure_cost_management` - (Optional) Azure Cost Management. See [Integration blocks](#integration-blocks) below for details.
* `azure_data_factory` - (Optional) for Azure Data Factory. See [Integration blocks](#integration-blocks) below for details.
* `azure_event_hub` - (Optional) for Azure Event Hub. See [Integration blocks](#integration-blocks) below for details.
* `azure_express_route` - (Optional) for Azure Express Route. See [Integration blocks](#integration-blocks) below for details.
* `azure_firewalls` - (Optional) for Azure Firewalls. See [Integration blocks](#integration-blocks) below for details.
* `azure_front_door` - (Optional) for Azure Front Door. See [Integration blocks](#integration-blocks) below for details.
* `azure_functions` - (Optional) for Azure Functions. See [Integration blocks](#integration-blocks) below for details.
* `azure_key_vault` - (Optional) for Azure Key Vault. See [Integration blocks](#integration-blocks) below for details.
* `azure_load_balancer` - (Optional) for Azure Load Balancer. See [Integration blocks](#integration-blocks) below for details.
* `azure_logic_apps` - (Optional) for Azure Logic Apps. See [Integration blocks](#integration-blocks) below for details.
* `azure_machine_learning` - (Optional) for Azure Machine Learning. See [Integration blocks](#integration-blocks) below for details.
* `azure_maria_db` - (Optional) for Azure MariaDB. See [Integration blocks](#integration-blocks) below for details.
* `azure_mysql` - (Optional) for Azure MySQL. See [Integration blocks](#integration-blocks) below for details.
* `azure_postgresql` - (Optional) for Azure PostgreSQL. See [Integration blocks](#integration-blocks) below for details.
* `azure_power_bi_dedicated` - (Optional) for Azure Power BI Dedicated. See [Integration blocks](#integration-blocks) below for details.
* `azure_redis_cache` - (Optional) for Azure Redis Cache. See [Integration blocks](#integration-blocks) below for details.
* `azure_service_bus` - (Optional) for Azure Service Bus. See [Integration blocks](#integration-blocks) below for details.
* `azure_service_fabric` - (Optional) for Azure Service Fabric. See [Integration blocks](#integration-blocks) below for details.
* `azure_sql` - (Optional) for Azure SQL. See [Integration blocks](#integration-blocks) below for details.
* `azure_sql_managed` - (Optional) for SQL Managed. See [Integration blocks](#integration-blocks) below for details.
* `azure_storage` - (Optional) for Azure Storage. See [Integration blocks](#integration-blocks) below for details.
* `azure_virtual_machine` - (Optional) for Azure Virtual machine. See [Integration blocks](#integration-blocks) below for details.
* `azure_vms` - (Optional) for Azure VMs. See [Integration blocks](#integration-blocks) below for details.
* `azure_vpn_gateway` - (Optional) for Azure VPN Gateway. See [Integration blocks](#integration-blocks) below for details.
* 
### `Integration` blocks

All `integration` blocks support the following common arguments:

* `metrics_polling_interval` - (Optional) The data polling interval in seconds.
* `resource_groups` - (Optional) Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive

Other integration type support an additional argument:

* `azure_cost_management`
  * `tag_keys` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the linked Azure account in New Relic.

## Import

Linked Azure accounts can be imported using `id`, you can find the `id` of existing Azure linked accounts in Azure dashboard under Infrastructure in NewRelic

```bash
$ terraform import newrelic_cloud_azure_integrations.foo <id>

```