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

resource "newrelic_cloud_azure_link_account" "foo"{
#  application_id = "1b88964c-9b73-46fc-8a60-d735600d4514"
#  client_secret = "gsx7Q~TjYAUJSWsT6eNQSVEW6bFHYVBo23D75"
#  subscription_id = "7bfb2949-e2fa-4bde-8705-067f14a947db"
#  tenant_id = "ad4c92da-afeb-42df-8548-96642d5183e0"
#  name = "qwerty"
}