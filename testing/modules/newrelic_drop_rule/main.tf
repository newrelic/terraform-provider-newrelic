resource "newrelic_nrql_drop_rule" "drop_rule" {
  account_id  = var.account_id
  description = var.description
  action      = var.action
  nrql        = var.nrql
}

output "all_rules" {
  description = "A map of all drop rule resource objects created by this module."
  value = {
    debug_logs     = newrelic_nrql_drop_rule.drop_rule
  }
}