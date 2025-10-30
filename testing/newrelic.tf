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

resource "newrelic_one_dashboard" "exampledash" {
  name        = "New Relic Terraform Example"
  permissions = "public_read_only"

  page {
    name = "New Relic Terraform Example"
  }

  variable {
    is_multi_selection = true
    item {
      title = "item"
      value = "ITEM"
    }
    name = "variable"
    nrql_query {
      account_ids = [3806526]
      query       = "FROM Transactions SELECT uniques(duration)"
    }
    replacement_strategy = "default"
    title                = "title"
    type                 = "nrql"

    # options {
    #   excluded = true
    #   ignore_time_range = true
    #   # show_apply_action = true
    # }

  }
}