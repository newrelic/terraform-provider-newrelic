
output "required_attributes" {
  value = {
    "key": local.key,
    "name": local.name,
    "key_type": local.type,
    "ingest_type": local.ingestType
  }
}

output "key_id" {
  value = length(newrelic_api_access_key.api_access_key) > 0 ? newrelic_api_access_key.api_access_key[0].id : null
}

output "key" {
  value = length(newrelic_api_access_key.api_access_key) > 0 ? newrelic_api_access_key.api_access_key[0].key : null
}
