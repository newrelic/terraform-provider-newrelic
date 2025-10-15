locals {
  enable_standalone_drop_rules = false
  enable_modular_drop_rules    = false

  # 1. Prepare the standalone rules (conditionally)
  standalone_rules = local.enable_standalone_drop_rules ? [
    {
      name     = "drop_health_checks"
      resource = newrelic_nrql_drop_rule.drop_health_checks
    },
    {
      name     = "drop_health_checks_three_a"
      resource = newrelic_nrql_drop_rule.drop_health_checks_three["a"]
    },
    {
      name     = "drop_health_checks_three_b"
      resource = newrelic_nrql_drop_rule.drop_health_checks_three["b"]
    },
    {
      name     = "drop_health_checks_three_c"
      resource = newrelic_nrql_drop_rule.drop_health_checks_three["c"]
    },
    {
      name     = "drop_health_checks_two_0"
      resource = newrelic_nrql_drop_rule.drop_health_checks_two[0]
    },
    {
      name     = "drop_health_checks_two_1"
      resource = newrelic_nrql_drop_rule.drop_health_checks_two[1]
    }

  ] : []

  # 2. Prepare the modular rules (conditionally)
  # This list is only used as an intermediate step for the flatten logic below.
  _modular_rules_raw = local.enable_modular_drop_rules ? [
    {
      name     = "drop_rules"
      resource = module.drop_rules["drop_debug_logs"].all_rules
    }
  ] : []

  # 3. Process the modular rules into a flat list (conditionally)
  modular_rules_flat = local.enable_modular_drop_rules ? flatten([
    for module_source in local._modular_rules_raw : [
      for rule_name, rule_object in module_source.resource : {
        name     = "${module_source.name}_${rule_name}"
        resource = rule_object
      }
    ]
  ]) : []

  # 4. Combine them using concat()
  # This is much cleaner and avoids the type error.
  all_individual_rules = concat(local.standalone_rules, local.modular_rules_flat)


  # Validation logic compatible with Terraform 1.2.x
  validation_results = [
    for drop_rule in local.all_individual_rules : {
      name      = drop_rule.name
      entity_id = try(drop_rule.resource.pipeline_cloud_rule_entity_id, null)
      is_valid  = try(drop_rule.resource.pipeline_cloud_rule_entity_id, null) != null
    }
  ]

  # Create error message if any validation fails
  validation_errors = length(local.all_individual_rules) == 0 ? [
    "❌ No drop rules have been listed in the local variables above. Please update the list as needed."
    ] : [
    for result in local.validation_results :
    "❌ Drop rule '${result.name}' does not export `pipeline_cloud_rule_entity_id` or it is null"
    if !result.is_valid
  ]
  #
  #   # Pre-computed JSON object to avoid complex expressions in output
  json_data_object = {
    drop_rule_resource_ids = [
      for drop_rule in local.all_individual_rules : {
        name                          = drop_rule.name
        id                            = try(drop_rule.resource.id, drop_rule.resource.drop_rule_id)
        pipeline_cloud_rule_entity_id = length(local.validation_errors) == 0 ? drop_rule.resource.pipeline_cloud_rule_entity_id : null
      }
    ]
  }
  #
  #   # Pre-computed formatted JSON string to avoid heredoc issues in conditional expressions
  formatted_json_string = <<-EOT
  {
    "drop_rule_resource_ids": [
  ${join(",\n", [for drop_rule in local.all_individual_rules : "    {\n      \"name\": \"${drop_rule.name}\",\n      \"id\": \"${try(drop_rule.resource.id, drop_rule.resource.drop_rule_id)}\",\n      \"pipeline_cloud_rule_entity_id\": ${length(local.validation_errors) == 0 ? "\"${drop_rule.resource.pipeline_cloud_rule_entity_id}\"" : "null"}\n    }"])}
      ]
  }
  EOT
}
#
# # Success validation output - only shows meaningful content when all validations pass
output "a_validation_success" {
  description = "Confirmation that all resources export pipeline_cloud_rule_entity_id"
  value       = length(local.validation_errors) == 0 ? "✅ All listed resources export pipeline_cloud_rule_entity_id" : null
}
#
# # JSON output with all experimental drop rule IDs - only when validation passes
output "experimental_drop_rule_resource_ids" {
  description = "JSON containing all drop rule resource names and IDs from experimental folder"
  value       = length(local.validation_errors) == 0 ? jsonencode(local.json_data_object) : null
}
#
# # Formatted/Pretty JSON output for better readability - only when validation passes
output "experimental_drop_rule_resource_ids_formatted" {
  description = "Pretty-formatted JSON containing all drop rule resource names and IDs from experimental folder"
  value       = length(local.validation_errors) == 0 ? local.formatted_json_string : null
}
#
# # Validation output - only shows meaningful content when there are errors
output "validation_errors" {
  description = "Lists any drop rule resources that do not export pipeline_cloud_rule_entity_id"
  value       = length(local.validation_errors) > 0 ? local.validation_errors : null
}


