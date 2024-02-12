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

data "newrelic_authentication_domain" "foo" {
  name = "Test Authentication Domain"
}

resource "newrelic_user_management" "foo" {
  name                     = "Test New User"
  email                    = "test_user@test.com"
  authentication_domain_id = data.newrelic_authentication_domain.foo.id
  user_type                = "CORE_USER_TIER"
}