locals {
  response = local.is_resource_created ? jsondecode(data.graphql_query.query_with_id[0].query_response): jsondecode(data.graphql_query.basic_query[0].query_response)
  key = local.response["data"]["actor"]["apiAccess"]["key"]["key"]
  name = local.response["data"]["actor"]["apiAccess"]["key"]["name"]
  type = local.response["data"]["actor"]["apiAccess"]["key"]["type"]
  ingestType = lookup(local.response["data"]["actor"]["apiAccess"]["key"],"ingestType",null)
  is_resource_created = var.newrelic_account_id != ""
}

variable "api_key" {
  type = string
  default = ""
}

variable "key_id" {
  type    = string
  default = "XXL"
}

variable "key_type" {
  type    = string
  default = "INGEST"
}

variable "graphiql_url" {
  type    = string
  default = "https://api.newrelic.com/graphql"
}

variable "newrelic_account_id" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
}

variable "ingest_type" {
  type    = string
  default = ""
}

variable "notes" {
  type    = string
  default = "API Key created using the newrelic_api_access_key Terraform resource"
}

variable "user_id" {
  type    = string
  default = null
}