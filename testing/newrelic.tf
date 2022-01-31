terraform {
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
    }
  }
}

provider "newrelic" {
  api_key = ""
  account_id = 0
  region = "US" # US or EU
}
