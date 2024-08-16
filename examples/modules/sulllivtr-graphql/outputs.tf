
output "required_attributes" {
  value = {
    "key_name": local.key_name,
    "name": local.name,
    "type": local.type
  }
}