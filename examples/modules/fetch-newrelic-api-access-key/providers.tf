terraform {
  required_providers {
    newrelic = {
      source = "newrelic/newrelic"
    }
    graphql = {
      source = "sullivtr/graphql"
       version = "2.5.4"
    }
  }
}

provider "newrelic" {
  region = "US" # US or EU
}


