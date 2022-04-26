
/*

    Complete example to enable Azure integration with New Relic

*/

terraform {
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
    }
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 2.15.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=3.0.0"
    }
  }
}

provider "azurerm" {
  features {}
}
provider "newrelic" {

  region = "US" # US or EU
}

data "azuread_client_config" "current" {}
data "azurerm_subscription" "primary" {}

resource "azuread_application" "example" {
  display_name = "NewRelic-Integrations-201"
  owners           = [data.azuread_client_config.current.object_id]
  web {
    redirect_uris = ["https://newrelic.com/"]
  }
}

resource "time_rotating" "example" {
  rotation_months = 6
}

resource "azuread_application_password" "example" {
  application_object_id = azuread_application.example.object_id
  rotate_when_changed = {
    rotation = time_rotating.example.id
  }
}

resource "azurerm_role_assignment" "example" {
  scope                = data.azurerm_subscription.primary.id
  role_definition_name = "Reader"
  principal_id         = data.azuread_client_config.current.object_id
}

resource "newrelic_cloud_azure_link_account" "newrelic_cloud_azure_integration_bar" {
  account_id = var.NEW_RELIC_ACCOUNT_ID
    application_id = "%[1]s"
    client_secret = "%[2]s"
    subscription_id = "%[3]s"
    tenant_id = "%[4]s"
    name  = "production-pull"
  }


resource "newrelic_cloud_azure_integrations" "foo" {
  account_id        = var.NEW_RELIC_ACCOUNT_ID
  linked_account_id = newrelic_cloud_azure_link_account.newrelic_cloud_azure_integration_bar.id

  api_management {}
  app_gateway {}
  app_service {}
  containers {}
  cosmos_db {}
  cost_management {}
  data_factory {}
  event_hub {}
  express_route {}
  firewalls {}
  front_door {}
  functions {}
  key_vault {}
  load_balancer {}
  logic_apps {}
  machine_learning {}
  maria_db {
  mysql {}
  postgresql {}
  power_bi_dedicated {}
  redis_cache {}
  service_bus {}
  service_fabric{}
  sql {}
  sql_managed{}
  storage {}
  virtual_machine {}
  virtual_networks {}
  vms {}
  vpn_gateway {}

}