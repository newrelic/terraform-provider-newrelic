locals {
  response = local.is_resource_created ? jsondecode(data.graphql_query.query_with_id[0].query_response): jsondecode(data.graphql_query.basic_query[0].query_response)
  key = local.response["data"]["actor"]["apiAccess"]["key"]["key"]
  name = local.response["data"]["actor"]["apiAccess"]["key"]["name"]
  type = local.response["data"]["actor"]["apiAccess"]["key"]["type"]
  ingestType = lookup(local.response["data"]["actor"]["apiAccess"]["key"],"ingestType",null)
  is_resource_created = var.newrelic_api_access_key_extended.newrelic_account_id != ""
}

variable "fetch_access_keys_service" {
  description = "The service is to get api keys"
  type = object({
    api_key                    = string
    key_id                     = string
    key_type                   = string
    graphiql_url               = optional(string,"https://api.newrelic.com/graphql")
  })
  default = {
    api_key  = ""
    key_id   = "XXXX"
    key_type = "XXXX"
  }
}

variable "newrelic_api_access_key_extended" {
  description = "The service is to create api keys"
  type = object({
    api_key                    = string
    newrelic_account_id        = string
    name                       = optional(string,"New API Key")
    key_type                   = string
    ingest_type                = optional(string,"")
    notes                      = optional(string,"API Key created using the newrelic_api_access_key Terraform resource")
    user_id                    = optional(string,null)
  })
  default = {
    api_key  = ""
    newrelic_account_id  = ""
    key_type             = "INGEST"
  }
}
