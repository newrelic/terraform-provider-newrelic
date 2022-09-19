terraform {
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
    }
  }
}

provider "newrelic" {
  region = "US" # US or EU
}

resource "newrelic_workload" "foo" {
  name = "pkWorkloads"
  account_id = 2520528

  entity_guids = ["MjUyMDUyOHxTWU5USHxNT05JVE9SfDYyMjUwYzk0LWNlY2EtNGJjOS04NDQ1LTcxZDFjOTcwNTEyYg"]

  entity_search_query {
    query = "name like '%Example application%'"
  }

  scope_account_ids =  [2520528]

  description = "Something"

  status_config_automatic{
    enabled = true
#        remaining_entities_rule_rollup {
#          strategy = "BEST_STATUS_WINS"
##          threshold_type = "FIXED"
##          threshold_value = 100
#          group_by = "ENTITY_TYPE"
#        }
#    rules{
#      entity_guids = ["MjUyMDUyOHxTWU5USHxNT05JVE9SfGQyN2ExNmFhLWRjZmUtNDBlZi05YjgxLWFjZDI2ZTc4MmJkZg"]
#      nrql_query{
#        query = "name like '%Example application%'"
#      }
#      rollup{
#        strategy = "BEST_STATUS_WINS"
#        threshold_type = "FIXED"
#        threshold_value = 100
#      }
#    }
#    rules{
#      entity_guids = ["MjUyMDUyOHxTWU5USHxNT05JVE9SfGQyN2ExNmFhLWRjZmUtNDBlZi05YjgxLWFjZDI2ZTc4MmJkZg"]
#      nrql_query{
#        query = "name like '%Example application%'"
#      }
#      rollup{
#        strategy = "BEST_STATUS_WINS"
#        threshold_type = "PERCENTAGE"
#        threshold_value = 10
#      }
#    }
  }
  status_config_static{
    description = "test"
    enabled = true
    status = "OPERATIONAL"
    summary = "egetgykwesgksegkerh"
  }
}
