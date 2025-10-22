locals {
  # ==========================================================================
  # CONFIGURATION FLAGS - CHANGE VALUE TO TRUE TO ENABLE PROCESSING
  # ==========================================================================

  # Enable for standalone drop rules
  # SET TO TRUE if you have standalone drop rules to process in the list `standalone_rules` below
  enable_standalone_drop_rules = false

  # Enable for modular drop rules
  # SET TO TRUE if you have modular drop rules to process in the list `_modular_rules_raw` below
  enable_modular_drop_rules = false

  # ==========================================================================
  # DROP RULE CONFIGURATION LISTS
  # ==========================================================================

  # Configure standalone drop rules here
  standalone_rules = local.enable_standalone_drop_rules ? [
    # ==================== PASTE STANDALONE DROP RULES HERE ====================
    # NOTE: Set enable_standalone_drop_rules = true above to process these rules
    # REPLACE THESE EXAMPLE COMMENTS WITH YOUR ACTUAL DROP RULE CONFIGURATIONS
    # Example format (remove these example lines and add your own):
    # {
    #   name     = "drop_health_checks"
    #   resource = newrelic_nrql_drop_rule.drop_health_checks
    # },
    # {
    #   name     = "drop_pii_data"
    #   resource = newrelic_nrql_drop_rule.drop_pii_data
    # },
    # REMOVE ALL EXAMPLE LINES ABOVE AND ADD YOUR ACTUAL DROP RULES BELOW
    # ========================================================================

  ] : []

  # Configure modular drop rules here
  _modular_rules_raw = local.enable_modular_drop_rules ? [
    # ==================== PASTE MODULAR DROP RULES HERE ====================
    # NOTE: Set enable_modular_drop_rules = true above to process these rules
    # REPLACE THESE EXAMPLE COMMENTS WITH YOUR ACTUAL MODULE CONFIGURATIONS
    # Example format (remove these example lines and add your own):
    # {
    #   name     = "drop_metadata"
    #   resource = module.drop_rules["drop_metadata"].all_rules
    # },
    # REMOVE ALL EXAMPLE LINES ABOVE AND ADD YOUR ACTUAL MODULE REFERENCES BELOW
    # =====================================================================

  ] : []

  # ===============================================================================
  # ===============================================================================
  # ========= INTERNAL PROCESSING - DO NOT EDIT ANYTHING BELOW THIS LINE! =========
  # ===============================================================================
  # ===============================================================================

  # Process modular rules into a flat list
  modular_rules_flat = local.enable_modular_drop_rules ? flatten([
    for module_source in local._modular_rules_raw : [
      for rule_name, rule_object in module_source.resource : {
        name     = "${module_source.name}_${rule_name}"
        resource = rule_object
      }
    ]
  ]) : []

  # Combine standalone and modular rules
  all_individual_drop_rules = concat(local.standalone_rules, local.modular_rules_flat)

  # Validate that all drop rules export pipeline_cloud_rule_entity_id
  validation_results = [
    for drop_rule in local.all_individual_drop_rules : {
      name      = drop_rule.name
      entity_id = try(drop_rule.resource.pipeline_cloud_rule_entity_id, null)
      is_valid  = try(drop_rule.resource.pipeline_cloud_rule_entity_id, null) != null
    }
  ]

  # Generate error messages for validation failures
  validation_errors = length(local.all_individual_drop_rules) == 0 ? [
    "❌ No drop rules have been listed in the local variables above. Please update the list as needed."
    ] : [
    for result in local.validation_results :
    "❌ Drop rule '${result.name}' does not export `pipeline_cloud_rule_entity_id` or it is null"
    if !result.is_valid
  ]

  # Prepare JSON data object for output
  json_data_object = {
    drop_rule_resource_ids = [
      for drop_rule in local.all_individual_drop_rules : {
        name                          = drop_rule.name
        id                            = try(drop_rule.resource.id, drop_rule.resource.drop_rule_id)
        pipeline_cloud_rule_entity_id = length(local.validation_errors) == 0 ? drop_rule.resource.pipeline_cloud_rule_entity_id : null
      }
    ]
  }

  # Prepare formatted JSON string for pretty output
  formatted_json_string = <<-EOT
  {
    "drop_rule_resource_ids": [
  ${join(",\n", [for drop_rule in local.all_individual_drop_rules : "    {\n      \"name\": \"${drop_rule.name}\",\n      \"id\": \"${try(drop_rule.resource.id, drop_rule.resource.drop_rule_id)}\",\n      \"pipeline_cloud_rule_entity_id\": ${length(local.validation_errors) == 0 ? "\"${drop_rule.resource.pipeline_cloud_rule_entity_id}\"" : "null"}\n    }"])}
      ]
  }
  EOT
}

# ==========================================================================
# TERRAFORM OUTPUTS
# ==========================================================================

# Validation success confirmation
output "a_validation_success" {
  description = "Confirmation that all resources export pipeline_cloud_rule_entity_id"
  value       = length(local.validation_errors) == 0 ? "✅ All listed resources export pipeline_cloud_rule_entity_id" : null
}

# Compact JSON output with drop rule IDs
output "experimental_drop_rule_resource_ids" {
  description = "JSON containing all drop rule resource names and IDs from experimental folder"
  value       = length(local.validation_errors) == 0 ? jsonencode(local.json_data_object) : null
}

# Pretty-formatted JSON output for better readability
output "experimental_drop_rule_resource_ids_formatted" {
  description = "Pretty-formatted JSON containing all drop rule resource names and IDs from experimental folder"
  value       = length(local.validation_errors) == 0 ? local.formatted_json_string : null
}

# Validation errors (only shown when there are errors)
output "validation_errors" {
  description = "Lists any drop rule resources that do not export pipeline_cloud_rule_entity_id"
  value       = length(local.validation_errors) > 0 ? local.validation_errors : null
}

