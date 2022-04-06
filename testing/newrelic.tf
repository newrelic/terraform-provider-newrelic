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

resource "newrelic_cloud_azure_integrations" "san" {

  account_id = 2520528
  linked_account_id = 110763

  azure_api_management {
    metrics_polling_interval = 1200

  }
}