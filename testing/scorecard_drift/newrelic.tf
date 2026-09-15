terraform {
  required_providers {
    newrelic = {
      source = "newrelic/newrelic"
    }
  }
}

provider "newrelic" {
  account_id = var.account_id
  api_key    = var.api_key
  region     = "US"
}

variable "account_id" {}
variable "api_key" {}
