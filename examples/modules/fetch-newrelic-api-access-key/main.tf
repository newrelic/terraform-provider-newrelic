
provider "graphql" {
  url = var.fetch_access_keys_service.graphiql_url
  headers = {
    "Content-Type" = "application/json"
    "API-Key" = var.fetch_access_keys_service.api_key != "" ? var.fetch_access_keys_service.api_key : var.create_access_keys_service.api_key
  }
}

data "graphql_query" "basic_query" {
  query_variables = {
    "id"        = var.fetch_access_keys_service.key_id
    "key_type"  = var.fetch_access_keys_service.key_type
  }
  query  = file("query.gql")
  count = local.is_resource_created ? 0 : 1
}

resource "newrelic_api_access_key" "api_access_key" {
  count  = var.create_access_keys_service.newrelic_account_id != "" ? 1 : 0
  account_id  = var.create_access_keys_service.newrelic_account_id
  key_type    = var.create_access_keys_service.key_type
  ingest_type = "LICENSE"
  name        = "APM ${var.create_access_keys_service.key_type} License Key for ${var.create_access_keys_service.name}"
  notes       = "To be used with service XXXX"
}

data "graphql_query" "query_with_id" {
  query_variables = {
    "id"        = newrelic_api_access_key.api_access_key[0].id
    "key_type"  = var.create_access_keys_service.key_type
  }
  query = file("query.gql")
  depends_on = [newrelic_api_access_key.api_access_key]
  count = local.is_resource_created ? 1 : 0
}





