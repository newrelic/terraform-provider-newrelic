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

Leave an integration block empty to use its default configuration. You can also use the [full example, including the Azure set up, found in our guides](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/cloud_integrations_guide#azure).

```hcl

resource "newrelic_cloud_azure_link_account" "foo" {
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
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  app_gateway {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  app_service {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  containers {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  cosmos_db {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  cost_management {
    metrics_polling_interval = 3600
    tag_keys = ["tag_keys"]
  }

  data_factory {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  event_hub {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  express_route {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  firewalls {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  front_door {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  functions {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  key_vault {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  load_balancer {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  logic_apps {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  machine_learning {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  maria_db {
    metrics_polling_interval = 3600
    resource_groups = ["resource_groups"]
  }

  monitor {
    metrics_polling_interval = 60
    resource_groups          = ["resource_groups"]
    include_tags             = ["env:production"]
    exclude_tags             = ["env:staging", "env:testing"]
    enabled                  = true
    resource_types           = ["microsoft.datashare/accounts"]
  }
  
  mysql {
    metrics_polling_interval = 3600
    resource_groups = ["resource_groups"]
  }

  mysql_flexible {
    metrics_polling_interval = 3600
    resource_groups = ["resource_groups"]
  }

  postgresql {
    metrics_polling_interval = 3600
    resource_groups = ["resource_groups"]
  }

  postgresql_flexible {
    metrics_polling_interval = 3600
    resource_groups = ["resource_groups"]
  }

  power_bi_dedicated {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  redis_cache {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  service_bus {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  sql {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  sql_managed {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  storage {
    metrics_polling_interval = 1800
    resource_groups = ["resource_groups"]
  }

  virtual_machine {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  virtual_networks {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  vms {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }

  vpn_gateway {
    metrics_polling_interval = 300
    resource_groups = ["resource_groups"]
  }
}
```
## Argument Reference

-> **WARNING:** Starting with [v3.27.2](https://registry.terraform.io/providers/newrelic/newrelic/3.27.2) of the New Relic Terraform Provider, updating the `linked_account_id` of a `newrelic_cloud_azure_integrations` resource that has been applied would **force a replacement** of the resource (destruction of the resource, followed by the creation of a new resource). When such an update is performed, please carefully review the output of `terraform plan`, which would clearly indicate a replacement of this resource, before performing a `terraform apply`.

The following arguments are supported with minimum metric polling interval of 300 seconds

* `account_id` - (Optional) The New Relic account ID to operate on.  This allows the user to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.
* `linked_account_id` - (Required) The ID of the linked Azure account in New Relic.
* `api_management` - (Optional) Azure API Management. See [Integration blocks](#integration-blocks) below for details.
* `app_gateway` - (Optional) Azure App Gateway. See [Integration blocks](#integration-blocks) below for details. 
* `app_service` - (Optional) Azure App Service. See [Integration blocks](#integration-blocks) below for details.
* `containers` - (Optional) Azure Containers. See [Integration blocks](#integration-blocks) below for details.
* `cosmos_db` - (Optional) Azure CosmosDB. See [Integration blocks](#integration-blocks) below for details.
* `data_factory` - (Optional) Azure Data Factory. See [Integration blocks](#integration-blocks) below for details.
* `event_hub` - (Optional) Azure Event Hub. See [Integration blocks](#integration-blocks) below for details.
* `express_route` - (Optional) Azure Express Route. See [Integration blocks](#integration-blocks) below for details.
* `firewalls` - (Optional) Azure Firewalls. See [Integration blocks](#integration-blocks) below for details.
* `front_door` - (Optional) Azure Front Door. See [Integration blocks](#integration-blocks) below for details.
* `functions` - (Optional) Azure Functions. See [Integration blocks](#integration-blocks) below for details.
* `key_vault` - (Optional) Azure Key Vault. See [Integration blocks](#integration-blocks) below for details.
* `load_balancer` - (Optional) Azure Load Balancer. See [Integration blocks](#integration-blocks) below for details.
* `logic_apps` - (Optional) Azure Logic Apps. See [Integration blocks](#integration-blocks) below for details.
* `machine_learning` - (Optional) Azure Machine Learning. See [Integration blocks](#integration-blocks) below for details.
* `maria_db` - (Optional) Azure MariaDB. See [Integration blocks](#integration-blocks) below for details.
* `monitor` - (Optional) Azure Monitor. See [Integration blocks](#integration-blocks) below for details.
* `mysql` - (Optional) Azure MySQL. See [Integration blocks](#integration-blocks) below for details.
* `mysql_flexible` - (Optional) Azure MySQL Flexible Server. See [Integration blocks](#integration-blocks) below for details.
* `postgresql` - (Optional) Azure PostgreSQL. See [Integration blocks](#integration-blocks) below for details.
* `postgresql_flexible` - (Optional) Azure PostgreSQL Flexible Server. See [Integration blocks](#integration-blocks) below for details.
* `power_bi_dedicated` - (Optional) Azure Power BI Dedicated. See [Integration blocks](#integration-blocks) below for details.
* `redis_cache` - (Optional) Azure Redis Cache. See [Integration blocks](#integration-blocks) below for details.
* `service_bus` - (Optional) Azure Service Bus. See [Integration blocks](#integration-blocks) below for details.
* `sql` - (Optional) Azure SQL. See [Integration blocks](#integration-blocks) below for details.
* `sql_managed` - (Optional) Azure SQL Managed. See [Integration blocks](#integration-blocks) below for details.
* `virtual_machine` - (Optional) Azure Virtual machine. See [Integration blocks](#integration-blocks) below for details.
* `vms` - (Optional) Azure VMs. See [Integration blocks](#integration-blocks) below for details.
* `vpn_gateway` - (Optional) Azure VPN Gateway. See [Integration blocks](#integration-blocks) below for details.

Below arguments supports the minimum metric polling interval of 900 seconds

* `storage` - (Optional) for Azure Storage. See [Integration blocks](#integration-blocks) below for details.
* `virtual_networks` - (Optional) for Azure Virtual networks. See [Integration blocks](#integration-blocks) below for details.

Below argument supports the minimum metric polling interval of 3600 seconds

* `cost_management` - (Optional) Azure Cost Management. See [Integration blocks](#integration-blocks) below for details.

### `Integration` blocks

All `integration` blocks support the following common arguments:

* `metrics_polling_interval` - (Optional) The data polling interval **in seconds**.
* `resource_groups` - (Optional) Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive

Other integration type support an additional argument:

* `cost_management`
  * `tag_keys` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.

* `monitor`
  * `resource_types` - (Optional) A list of Azure resource types that need to be monitored.
  * `include_tags` - (Optional) A list of resource tags associated with the resources that need to be monitored, in a "key:value" format. If this is not specified, all resources will be monitored.
  * `exclude_tags` - (Optional) A list of resource tags associated with the resources that need to be excluded from monitoring.
  * `enabled` - (Optional) A boolean value, that specifies if the integration needs to be active. Defaults to 'true' if not specified.

-> **IMPORTANT!** Using the `monitor` integration along with other polling integrations in this resource might lead to duplication of metrics. More information about this scenario may be found in the note in [this section](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/azure-integrations-list/azure-monitor/#migration-from-polling) of New Relic's documentation on the Azure Monitor integration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the linked Azure account in New Relic.

## Import

Linked Azure accounts can be imported using `id`, you can find the `id` of existing Azure linked accounts in Azure dashboard under Infrastructure in NewRelic.

```bash
$ terraform import newrelic_cloud_azure_integrations.foo <id>

```
