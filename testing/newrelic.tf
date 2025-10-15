terraform {
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
      version = "3.63.0"
    }
  }
}

provider "newrelic" {
  region = "US" # US or EU
}

variable "new_relic_account_id" {
  description = "New Relic Account ID"
  type        = string
  default     = 3806526
}

locals {

  drop_rules_set = toset(["a", "b", "c"])

  drop_rules = [
    {
      name        = "drop_debug_logs"
      description = "Modular: Filters out debug level logs from production environment to reduce data volume"
      action      = "drop_data"
      nrql        = "SELECT * FROM Log WHERE level = 'debug' AND environment = 'production'"
    },
    {
      name        = "drop_health_checks"
      description = "Modular: Removes userEmail and userName attributes from MyCustomEvent data"
      action      = "drop_attributes"
      nrql        = "SELECT userEmail, userName FROM MyCustomEvent"
    },
    {
      name        = "drop_pii_data"
      description = "Modular: Excludes containerId attribute from Metric aggregates"
      action      = "drop_attributes_from_metric_aggregates"
      nrql        = "SELECT containerId FROM Metric"
    }
  ]
}

module "drop_rules" {
  for_each = { for rule in local.drop_rules : rule.name => rule }
  source   = "./modules/newrelic_drop_rule"

  account_id  = var.new_relic_account_id
  description = each.value.description
  action      = each.value.action
  nrql        = each.value.nrql
}

resource "newrelic_nrql_drop_rule" "drop_health_checks" {
  account_id  = var.new_relic_account_id
  description = "Removes userEmail and userName attributes from MyCustomEvent data"
  action      = "drop_attributes"
  nrql        = "SELECT userEmail, userName FROM MyCustomEvent"
}

resource "newrelic_nrql_drop_rule" "drop_health_checks_two" {
  count       = 2
  account_id  = var.new_relic_account_id
  description = "Removes userEmail and userName attributes from MyCustomEvent data"
  action      = "drop_attributes"
  nrql        = "SELECT userEmail, userName FROM MyCustomEvent"
}

resource "newrelic_nrql_drop_rule" "drop_health_checks_three" {
  for_each    = local.drop_rules_set
  account_id  = var.new_relic_account_id
  description = "${each.value} Removes userEmail and userName attributes from MyCustomEvent data"
  action      = "drop_attributes"
  nrql        = "SELECT userEmail, userName FROM MyCustomEvent"
}
