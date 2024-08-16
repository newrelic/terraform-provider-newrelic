locals {
  response = jsondecode(data.graphql_query.basic_query.query_response)
  key_name = local.response["data"]["actor"]["apiAccess"]["key"]["key"]
  name = local.response["data"]["actor"]["apiAccess"]["key"]["name"]
  type = local.response["data"]["actor"]["apiAccess"]["key"]["type"]
}

variable "service" {
  description = "The service is to get api keys"
  type = object({
    api_key                    = string
    key_id                     = string
    key_type                   = string
    graphiql_url               = optional(string,"https://api.newrelic.com/graphql")
  })
  default = {
    api_key  = "XXXXX-XXXXXXXXXX"
    key_id   = "XXXX"
    key_type = "XXXX"
  }
}

