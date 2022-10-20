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

resource "newrelic_workload" "foo" {
  name = "tf-test-workload-mbazhlek"

  entity_search_query {
    query = "name like 'App'"
  }
}