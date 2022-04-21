
/*

    Complete example to enable Azure integration with New Relic

*/

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