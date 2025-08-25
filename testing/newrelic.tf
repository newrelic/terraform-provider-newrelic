terraform {
  required_providers {
    newrelic = {
      source = "newrelic/newrelic"
      version = "3.65.0"
    }
  }
}

# resource "newrelic_nrql_drop_rule" "foo" {
#   account_id  = 3806526
#   description = "Drop Rule New 12 August"
#   action      = "drop_data"
#   nrql        = "SELECT * FROM Pranav WHERE appName='Vinay' AND environment='Shashank'"
# }
#
# resource "newrelic_nrql_drop_rule" "foo2" {
#   account_id  = 3806526
#   description = "Drop Rule New 14 August"
#   action      = "drop_data"
#   nrql        = "SELECT * FROM Pranav WHERE appName='Vinay' AND environment='Shashank'"
# }

# locals {
#   drop_rules = [
#     newrelic_nrql_drop_rule.foo,
#     newrelic_nrql_drop_rule.foo2,
#   ]
# }
#
# output "drop_rules" {
#   value = local.drop_rules_with_identifiers
# }
#
# data "newrelic_drop_rule_pipeline_cloud_rule_relationship" "foo" {
#   for_each = { for rule in local.drop_rules : rule.id => rule }
#   drop_rule_id = each.value.id
# }

# output "x" {
#   value = data.newrelic_drop_rule_pipeline_cloud_rule_relationship.foo
# }
#
# data "external" "example" {
#   program = ["bash", "something.sh"]
#
#   query = {
#     # arbitrary map from strings to strings, passed
#     # to the external program as the data query.
#     foo = "abc123"
#     baz = "def456"
#   }
# }
#
# output "y" {
#   value = data.external.example.result
# }

# import {
#   count = length(data.newrelic_drop_rule_pipeline_cloud_rule_relationship.foo)
#   to       = newrelic_pipeline_cloud_rule[count.index]
#   id       = each.value.id
# }



#
# output "foo" {
#   value = data.newrelic_drop_rule_pipeline_cloud_rule_relationship.foo.id
# }
# #
# import {
#   id = data.newrelic_drop_rule_pipeline_cloud_rule_relationship.foo.id
#   to = newrelic_pipeline_cloud_rule.foo
# }
