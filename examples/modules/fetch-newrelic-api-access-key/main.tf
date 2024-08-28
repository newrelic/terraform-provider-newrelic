
provider "graphql" {
  url = var.newrelic_api_access_key_extended.graphiql_url
  headers = {
    "Content-Type" = "application/json"
    "API-Key" = var.newrelic_api_access_key_extended.api_key != "" ? var.newrelic_api_access_key_extended.api_key : var.create_access_keys_service.api_key
  }
}

data "graphql_query" "basic_query" {
  query_variables = {
    "id"        = var.newrelic_api_access_key_extended.key_id
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
  count = local.is_resource_created ? 0 : 1
}

resource "newrelic_api_access_key" "api_access_key" {
  count  = var.create_access_keys_service.newrelic_account_id != "" ? 1 : 0
  account_id  = var.create_access_keys_service.newrelic_account_id
  key_type    = var.create_access_keys_service.key_type
  name        = "${var.create_access_keys_service.key_type != "USER" ? "APM " : "" }${var.create_access_keys_service.key_type}${var.create_access_keys_service.key_type != "USER" ? "-" : "" }${var.create_access_keys_service.ingest_type} Key for ${var.create_access_keys_service.name}"
  notes       = var.create_access_keys_service.notes
  user_id     = var.create_access_keys_service.key_type == "USER" ? var.create_access_keys_service.user_id : null
  ingest_type = var.create_access_keys_service.key_type == "INGEST" ? var.create_access_keys_service.ingest_type : null
}

data "graphql_query" "query_with_id" {
  query_variables = {
    "id"        = newrelic_api_access_key.api_access_key[0].id
    "key_type"  = var.create_access_keys_service.key_type
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





