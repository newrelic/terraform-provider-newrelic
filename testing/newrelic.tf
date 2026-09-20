terraform {
  required_providers {
    newrelic = {
      source = "newrelic/newrelic"
    }
  }
}

provider "newrelic" {
  region = "US"
}

variable "account_id" {
  type    = number
  default = 3806526
}
