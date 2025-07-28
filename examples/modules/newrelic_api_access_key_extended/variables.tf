locals {
  response = local.is_resource_created ? jsondecode(data.graphql_query.query_with_id[0].query_response): jsondecode(data.graphql_query.basic_query[0].query_response)
  key = local.response["data"]["actor"]["apiAccess"]["key"]["key"]
  name = local.response["data"]["actor"]["apiAccess"]["key"]["name"]
  type = local.response["data"]["actor"]["apiAccess"]["key"]["type"]
  ingestType = lookup(local.response["data"]["actor"]["apiAccess"]["key"],"ingestType",null)
  fetch_user_id_response = (var.user_id == null ||  var.user_id == "") && var.key_type == "USER" ? jsondecode(data.graphql_query.fetch_user_id[0].query_response) : null
  user_id_from_graphql_response = local.fetch_user_id_response != null ? local.fetch_user_id_response["data"]["actor"]["user"]["id"] : null
  is_resource_created = var.newrelic_account_id != ""
}

variable "api_key" {
  description = "New Relic API Key for authentication. Set via TF_VAR_newrelic_api_key."
  type        = string
  sensitive   = true
}

check "check_api_key_source" {
  assert {
    # The condition will check if the variable 'newrelic_api_key' has a non-empty value. If the user hardcoded it as an empty string, or didn't set the env var, then this check will fail and error out.
    condition = var.api_key != null && var.api_key != ""
    error_message = "A valid 'api_key' is mandatory and should be provided exclusively through the TF_VAR_newrelic_api_key environment variable and not hardcoded directly in the configuration to prevent exposure. Providing an empty value (\"\") or omitting this environment variable will cause API calls to fail with 'authentication required' errors. Please set TF_VAR_newrelic_api_key with your actual New Relic API key."
  }
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