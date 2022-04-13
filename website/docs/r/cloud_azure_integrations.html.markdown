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
* `azure_api_management` - (Optional) for Azure API Management refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-api-management-monitoring-integration).
* `azure_app_gateway` - (Optional) for Azure App Gateway refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-application-gateway-monitoring-integration).
* `azure_app_service` - (Optional) for Azure App Service refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-app-service-monitoring-integration).
* `azure_containers` - (Optional) for Azure Containers refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-containers-monitoring-integration).
* `azure_cosmos_db` - (Optional) for Azure CosmosDB refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-cosmos-db-document-db-monitoring-integration).
* `azure_cost_management` - (Optional) for Azure Cost Management refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-cost-management-monitoring-integration).
* `azure_data_factory` - (Optional) for Azure Data Factory refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-data-factory-integration).
* `azure_event_hub` - (Optional) for Azure Event Hub refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-event-hub-monitoring-integration).
* `azure_express_route` - (Optional) for Azure Express Route refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-express-route-monitoring-integration).
* `azure_firewalls` - (Optional) for Azure Firewalls refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-firewalls-monitoring-integration).
* `azure_front_door` - (Optional) for Azure Front Door refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-front-door-monitoring-integration).
* `azure_functions` - (Optional) for Azure Functions refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-functions-monitoring-integration).
* `azure_key_vault` - (Optional) for Azure Key Vault refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-key-vault-monitoring-integration).
* `azure_load_balancer` - (Optional) for Azure Load Balancer refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-load-balancer-monitoring-integration).
* `azure_logic_apps` - (Optional) for Azure Logic Apps refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-logic-apps-monitoring-integration).
* `azure_machine_learning` - (Optional) for Azure Machine Learning refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-machine-learning-integration).
* `azure_maria_db` - (Optional) for Azure MariaDB refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-database-mariadb-monitoring-integration).
* `azure_mysql` - (Optional) for Azure MySQL refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-database-mysql-monitoring-integration).
* `azure_postgresql` - (Optional) for Azure PostgreSQL refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-database-postgresql-monitoring-integration).
* `azure_power_bi_dedicated` - (Optional) for Azure Power BI Dedicated refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-power-bi-dedicated-capacities-monitoring-integration).
* `azure_redis_cache` - (Optional) for Azure Redis Cache refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-redis-cache-monitoring-integration).
* `azure_service_bus` - (Optional) for Azure Service Bus refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-service-bus-monitoring-integration).
* `azure_service_fabric` - (Optional) for Azure Service Fabric refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-service-fabric-monitoring-integration).
* `azure_sql` - (Optional) for Azure SQL refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-sql-database-monitoring-integration).
* `azure_sql_managed` - (Optional) for SQL Managed refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-sql-managed-instances-monitoring-integration).
* `azure_storage` - (Optional) for Azure Storage refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-storage-monitoring-integration).
* `azure_virtual_machine` - (Optional) for Azure Virtual machine refer  [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-virtual-machine-scale-sets-monitoring-integration).
* `azure_virtual_networks` - (Optional) for Azure Virtual Network refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-virtual-network-monitoring-integration).
* `azure_vms` - (Optional) for Azure VMs refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-vms-monitoring-integration).
* `azure_vpn_gateway` - (Optional) for Azure VPN Gateway refer [our documentation](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-vpn-gateway-integration).

### Integration blocks

All `integration` blocks support the following common arguments:

* `metrics_polling_interval` - (Optional) The data polling interval in seconds.
* `resource_groups` - (Optional) Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive

Other integration type support an additional argument:

* `azure_cost_management`
  * `tag_keys` - (Optional) Specify if additional cost data per tag should be collected.



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the linked Azure account in New Relic.

## Import

Linked Azure accounts can be imported using `id`, you can find the `id` of existing Azure linked accounts in Azure dashboard under Infrastructure in NewRelic

```bash
$ terraform import newrelic_cloud_azure_integrations.foo <id>

```