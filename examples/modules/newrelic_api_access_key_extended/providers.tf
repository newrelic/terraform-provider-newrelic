terraform {
  required_providers {
    newrelic = {
      source = "newrelic/newrelic"
    }
    graphql = {
      source = "sullivtr/graphql"
    }
  }
}

provider "newrelic" {
  region = "US" # US or EU
}


