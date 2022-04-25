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

Setup is required for this resource to work properly. This resource assumes you have [linked an Azure account](cloud_azure_link_account.html.markdown) to New Relic.

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
    
    api_management {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    app_gateway {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    app_service {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    containers {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    cosmos_db {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }
    
    cost_management {
      metrics_polling_interval = 3600
      tag_keys = ["tag_keys"]
    }

    data_factory {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    event_hub {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    express_route {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    firewalls {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    front_door {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    functions {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    key_vault {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    load_balancer {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    logic_apps {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    machine_learning {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    maria_db {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    mysql {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    postgresql {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    power_bi_dedicated {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }
    
    redis_cache {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    service_bus {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    sql {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    sql_managed {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    storage {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    virtual_machine {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    virtual_networks {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    vms {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }

    vpn_gateway {
      metrics_polling_interval = 1200
      resource_groups = ["resource_groups"]
    }
}
```
## Argument Reference


The following arguments are supported:

* `account_id` - (Optional) The New Relic account ID to operate on.  This allows the user to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.
* `linked_account_id` - (Required) The ID of the linked Azure account in New Relic.
* `api_management` - (Optional) Azure API Management. See [Integration blocks](#integration-blocks) below for details.
* `app_gateway` - (Optional) Azure App Gateway. See [Integration blocks](#integration-blocks) below for details.
* `app_service` - (Optional) Azure App Service. See [Integration blocks](#integration-blocks) below for details.
* `containers` - (Optional) Azure Containers. See [Integration blocks](#integration-blocks) below for details.
* `cosmos_db` - (Optional) Azure CosmosDB. See [Integration blocks](#integration-blocks) below for details.
* `cost_management` - (Optional) Azure Cost Management. See [Integration blocks](#integration-blocks) below for details.
* `data_factory` - (Optional) for Azure Data Factory. See [Integration blocks](#integration-blocks) below for details.
* `event_hub` - (Optional) for Azure Event Hub. See [Integration blocks](#integration-blocks) below for details.
* `express_route` - (Optional) for Azure Express Route. See [Integration blocks](#integration-blocks) below for details.
* `firewalls` - (Optional) for Azure Firewalls. See [Integration blocks](#integration-blocks) below for details.
* `front_door` - (Optional) for Azure Front Door. See [Integration blocks](#integration-blocks) below for details.
* `functions` - (Optional) for Azure Functions. See [Integration blocks](#integration-blocks) below for details.
* `key_vault` - (Optional) for Azure Key Vault. See [Integration blocks](#integration-blocks) below for details.
* `load_balancer` - (Optional) for Azure Load Balancer. See [Integration blocks](#integration-blocks) below for details.
* `logic_apps` - (Optional) for Azure Logic Apps. See [Integration blocks](#integration-blocks) below for details.
* `machine_learning` - (Optional) for Azure Machine Learning. See [Integration blocks](#integration-blocks) below for details.
* `maria_db` - (Optional) for Azure MariaDB. See [Integration blocks](#integration-blocks) below for details.
* `mysql` - (Optional) for Azure MySQL. See [Integration blocks](#integration-blocks) below for details.
* `postgresql` - (Optional) for Azure PostgreSQL. See [Integration blocks](#integration-blocks) below for details.
* `power_bi_dedicated` - (Optional) for Azure Power BI Dedicated. See [Integration blocks](#integration-blocks) below for details.
* `redis_cache` - (Optional) for Azure Redis Cache. See [Integration blocks](#integration-blocks) below for details.
* `service_bus` - (Optional) for Azure Service Bus. See [Integration blocks](#integration-blocks) below for details.
* `sql` - (Optional) for Azure SQL. See [Integration blocks](#integration-blocks) below for details.
* `sql_managed` - (Optional) for SQL Managed. See [Integration blocks](#integration-blocks) below for details.
* `storage` - (Optional) for Azure Storage. See [Integration blocks](#integration-blocks) below for details.
* `virtual_machine` - (Optional) for Azure Virtual machine. See [Integration blocks](#integration-blocks) below for details.
* `vms` - (Optional) for Azure VMs. See [Integration blocks](#integration-blocks) below for details.
* `vpn_gateway` - (Optional) for Azure VPN Gateway. See [Integration blocks](#integration-blocks) below for details.

### `Integration` blocks

All `integration` blocks support the following common arguments:

* `metrics_polling_interval` - (Optional) The data polling interval in seconds.
* `resource_groups` - (Optional) Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive

Other integration type support an additional argument:

* `cost_management`
  * `tag_keys` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the linked Azure account in New Relic.

## Import

Linked Azure accounts can be imported using `id`, you can find the `id` of existing Azure linked accounts in Azure dashboard under Infrastructure in NewRelic.

```bash
$ terraform import newrelic_cloud_azure_integrations.foo <id>

```