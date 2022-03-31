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


resource "newrelic_cloud_gcp_integrations" "foo1" {
  account_id = 2520528
  linked_account_id = 110602
  app_engine {
    metrics_polling_interval = 500
  }
}

