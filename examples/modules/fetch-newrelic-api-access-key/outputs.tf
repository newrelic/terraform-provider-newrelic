
output "required_attributes" {
  value = {
    "key": local.key,
    "name": local.name,
    "type": local.type
  }
}

output "key_id" {
  value = length(newrelic_api_access_key.api_access_key) > 0 ? newrelic_api_access_key.api_access_key[0].id : null
}

output "key" {
  value = length(newrelic_api_access_key.api_access_key) > 0 ? newrelic_api_access_key.api_access_key[0].key : null
}
