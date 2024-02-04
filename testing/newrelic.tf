terraform {
  required_providers {
    newrelic = {
      source = "newrelic/newrelic"
    }
  }
}

provider "newrelic" {
  region = "US" # US or EU
}

resource "newrelic_group_management" "foo" {
  name = "Some Random Pro Max Group"
  authentication_domain_id = "84cb286a-8eb0-4478-b469-cdf2ccfef553"
  users = ["1004823231", "1005720113"]
}