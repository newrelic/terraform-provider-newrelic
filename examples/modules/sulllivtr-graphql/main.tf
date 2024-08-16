
provider "graphql" {
  url = var.service.graphiql_url
  headers = {
    "Content-Type" = "application/json"
    "API-Key" = var.service.api_key
  } 
}

resource "newrelic_browser_application" "access_data" {
  name = "New relic access data info"
}

data "graphql_query" "basic_query" {
  query_variables = {
    "id" = var.service.key_id
    "key_type" = var.service.key_type
  }
  query  = file("query.gql")
}




