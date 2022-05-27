/*

    Complete example to enable Azure integration with New Relic

*/

variable "NEW_RELIC_ACCOUNT_ID" {
    type = string
}

variable "NEW_RELIC_ACCOUNT_NAME" {
    type = string
    default = "Production"
}

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

data "azuread_client_config" "newrelic_client_config" {}
data "azurerm_subscription" "newrelic_subscription" {}
data "azuread_application_published_app_ids" "well_known" {}

resource "azuread_service_principal" "msgraph" {
  application_id = data.azuread_application_published_app_ids.well_known.result.MicrosoftGraph
  use_existing   = true
}

resource "azuread_application" "newrelic_application" {
    display_name     = "NewRelic-Integrations"
    owners           = [data.azuread_client_config.newrelic_client_config.object_id]
    sign_in_audience = "AzureADMyOrg"

    web {
        redirect_uris = ["https://www.newrelic.com/"]
    }
}

resource "azuread_service_principal" "newrelic_service_principal" {
  application_id = azuread_application.newrelic_application.application_id
}

resource "azurerm_role_assignment" "newrelic_role_assignment" {
  scope                = data.azurerm_subscription.newrelic_subscription.id
  role_definition_name = "Reader"
  principal_id         = azuread_service_principal.newrelic_service_principal.object_id
}

resource "azuread_application_password" "newrelic_application_password" {
  application_object_id = azuread_application.newrelic_application.object_id
}

resource "newrelic_cloud_azure_link_account" "newrelic_cloud_azure_integration" {
    account_id = var.NEW_RELIC_ACCOUNT_ID
    application_id = azuread_application.newrelic_application.application_id
    client_secret = azuread_application_password.newrelic_application_password.value
    subscription_id = data.azurerm_subscription.newrelic_subscription.subscription_id
    tenant_id = data.azurerm_subscription.newrelic_subscription.tenant_id
    name  = var.NEW_RELIC_ACCOUNT_NAME

    depends_on = [
        azurerm_role_assignment.newrelic_role_assignment
    ]
}

resource "newrelic_cloud_azure_integrations" "foo" {
    account_id        = var.NEW_RELIC_ACCOUNT_ID
    linked_account_id = newrelic_cloud_azure_link_account.newrelic_cloud_azure_integration.id

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
    maria_db {}
    mysql {}
    postgresql {}
    power_bi_dedicated {}
    redis_cache {}
    service_bus {}
    sql {}
    sql_managed{}
    storage {}
    virtual_machine {}
    virtual_networks {}
    vms {}
    vpn_gateway {}
}
