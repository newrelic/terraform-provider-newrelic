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
  name        = "New Relic widget histogram 9999"
  permissions = "public_read_only"

  page {
    name = "New Relic widget histogram 9999"

    widget_table {
      title  = "List of Transactions"
      row    = 1
      column = 4
      width  = 6
      height = 3

      nrql_query {
        account_id = 3806526
        query = "SELECT average(duration), max(duration), min(duration) FROM Transaction FACET name SINCE 1 day ago"
      }


    }

    widget_billboard {
      title  = "Requests per minute"
      row    = 1
      column = 1
      width  = 7
      height = 3

      nrql_query {
        query = "FROM Transaction SELECT rate(count(*), 1 minute)"
      }

    }

  }

  variable {
    is_multi_selection = true
    name = "select_country"
    nrql_query {
      account_ids = [3806526]
      query       = "SELECT uniques(ES)"
    }
    replacement_strategy = "string"
    title                = "Isa prueba"
    type                 = "nrql"
  }
}
