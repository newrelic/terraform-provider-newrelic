
provider "graphql" {
  url = var.fetch_access_keys_service.graphiql_url
  headers = {
    "Content-Type" = "application/json"
    "API-Key" = var.fetch_access_keys_service.api_key != "" ? var.fetch_access_keys_service.api_key : var.newrelic_api_access_key_extended.api_key
  }
}

data "graphql_query" "basic_query" {
  query_variables = {
    "id"        = var.fetch_access_keys_service.key_id
    "key_type"  = var.fetch_access_keys_service.key_type
  }
  query = <<EOF
    query getUser($id: ID!, $key_type: ApiAccessKeyType!) {
      actor {
      apiAccess {
        key(id: $id, keyType: $key_type) {
        key
        name
        type
        ... on ApiAccessIngestKey {
          ingestType
        }
        }
      }
      }
    }
    EOF
  count = local.is_resource_created ? 0 : 1
}

resource "newrelic_api_access_key" "api_access_key" {
  count  = var.newrelic_api_access_key_extended.newrelic_account_id != "" ? 1 : 0
  account_id  = var.newrelic_api_access_key_extended.newrelic_account_id
  key_type    = var.newrelic_api_access_key_extended.key_type
  name        = "${var.newrelic_api_access_key_extended.key_type != "USER" ? "APM " : "" }${var.newrelic_api_access_key_extended.key_type}${var.newrelic_api_access_key_extended.key_type != "USER" ? "-" : "" }${var.newrelic_api_access_key_extended.ingest_type} Key for ${var.newrelic_api_access_key_extended.name}"
  notes       = var.newrelic_api_access_key_extended.notes
  user_id     = var.newrelic_api_access_key_extended.key_type == "USER" ? var.newrelic_api_access_key_extended.user_id : null
  ingest_type = var.newrelic_api_access_key_extended.key_type == "INGEST" ? var.newrelic_api_access_key_extended.ingest_type : null
}

data "graphql_query" "query_with_id" {
  query_variables = {
    "id"        = newrelic_api_access_key.api_access_key[0].id
    "key_type"  = var.newrelic_api_access_key_extended.key_type
  }
  query = <<EOF
    query getUser($id: ID!, $key_type: ApiAccessKeyType!) {
      actor {
      apiAccess {
        key(id: $id, keyType: $key_type) {
        key
        name
        type
        ... on ApiAccessIngestKey {
          ingestType
        }
        }
      }
      }
    }
    EOF
  depends_on = [newrelic_api_access_key.api_access_key]
  count = local.is_resource_created ? 1 : 0
}





