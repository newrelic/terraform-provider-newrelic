terraform {
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
    }
  }
}

provider "newrelic" {

  region = "US" # US or EU
}

resource "newrelic_cloud_azure_integrations" "foo" {
  account_id = 2520528
  linked_account_id = 111297
  azure_api_management {
    metrics_polling_interval = 1000
  }
}
