
provider "graphql" {
  url = var.graphiql_url
  headers = {
    "Content-Type" = "application/json"
    "API-Key" = var.api_key != "" ? var.api_key : var.api_key
  }
}

data "graphql_query" "basic_query" {
  query_variables = {
    "id"        = var.key_id
    "key_type"  = var.key_type
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
  count  = var.newrelic_account_id != "" ? 1 : 0
  account_id  = var.newrelic_account_id
  key_type    = var.key_type
  name        = "${var.key_type != "USER" ? "APM " : "" }${var.key_type}${var.key_type != "USER" ? "-" : "" }${var.ingest_type} Key for ${var.name}"
  notes       = var.notes
  user_id     = var.key_type == "USER" ? var.user_id : null
  ingest_type = var.key_type == "INGEST" ? var.ingest_type : null
}

data "graphql_query" "query_with_id" {
  query_variables = {
    "id"        = newrelic_api_access_key.api_access_key[0].id
    "key_type"  = var.key_type
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





