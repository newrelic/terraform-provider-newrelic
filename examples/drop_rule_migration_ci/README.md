# New Relic Terraform Provider: `newrelic_nrql_drop_rule` -> `newrelic_pipeline_cloud_rule` Migration Guide for CI/CD Workflows Using Automation Helpers

This guide describes the **first phase** of a three-phase automation helper-process designed to assist in migrating `newrelic_nrql_drop_rule` resources (Drop Rules, managed via Terraform) to `newrelic_pipeline_cloud_rule` resources in CI-based environments, such as Atlantis and Grandcentral. This migration is necessary due to the upcoming end-of-life (EOL) of NRQL Drop Rules, scheduled for June 30, 2026. After this date, any Drop Rules managed via the `newrelic_nrql_drop_rule` Terraform resource will no longer function. For a general overview and more details on the EOL, its implications, and the required actions to replace `newrelic_nrql_drop_rule` resources, refer to [this detailed article](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/drop_rules_eol_guide).

For context, the three-phase migration process consists of (see an overview of this in the [NRQL Drop Rule EOL Guide in the documentation of the New Relic Terraform Provider](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/drop_rules_eol_guide)):
- **Phase 1** **(with the procedure outlined in this document)** - executed in the CI/CD environment with inputs on Terraform-managed `newrelic_nrql_drop_rule` resources added to a custom script, which is in turn added to the CI/CD environment and applied, to identify and export existing drop rules as JSON data; 
- **Phase 2** - executed locally using the New Relic CLI command `tf-importgen-ci` to process the JSON data and generate Pipeline Cloud Rule configurations and import scripts; and 
- **Phase 3** - executed back in the CI/CD environment to apply the generated configurations, import the new Pipeline Cloud Rules, and remove the legacy NRQL Drop Rules from management.

While more details on the working of the second and third phase may be found in the documentation and logs of the `tf-importgen-ci` command used in the second phase, **this document outlines the first phase of the three-phase automation process**, which involves modifying the configuration of `newrelic_nrql_drop_rule` resources in the CI environment according to the guidelines provided below. These modifications will enable the export of Drop Rule details in a structured format.

The data exported during this phase will be used in the second phase, which involves running a New Relic CLI command in a local environment. This command processes the exported data and generates Terraform configuration that, when integrated into the CI environment, facilitates the import of all identified Drop Rules as Pipeline Cloud Rules. Therefore, completing this first phase is a prerequisite for the subsequent phases, which automate the generation of Pipeline Cloud Rule-based configurations for existing Drop Rules. For instructions on the second and third phases, refer to the relevant sections of this guide.

## Overview

The `outputs.tf` file provides a validation and extraction system for drop rule resources, outputting their IDs in JSON format for use in migration procedures within CI/CD pipelines.

## Considerations and Prerequisites

- The procedure specified below **can be run on one workspace at a time; i.e., it may be applied to resources managed by a single state file at once.**

- The file `outputs.tf` in this directory must be copied, and **all of the instructions specified pertain to edits to be made to a copy of `outputs.tf`, after which it shall be added to the workspace in the CI.**

### Technical Prerequisites
- To allow Drop Rules to export IDs of their Pipeline Cloud Rule counterparts, **the procedure requires the New Relic Terraform Provider version 3.72.1 or above (preferably, the latest version available).** Hence, the version of the New Relic Terraform Provider will need to be upgraded in the CI environment.

## Configuration

As specified in the section above, make a copy of `outputs.tf` from this directory and add it to the workspace in the CI/CD environment where the `newrelic_nrql_drop_rule` resources are managed. Then, follow the steps below to configure the file according to your requirements.

### Step 1: Identify the Types of Drop Rules to Migrate and Specify Script Options Accordingly

The first step involves identifying drop rules present in the workspace (Terraform-managed via the resource `newrelic_nrql_drop_rule`) and listing them in designated locations.

First, configure which type of drop rules you want to process by setting the appropriate flags in the `locals` block:

```hcl  
locals {
   # Enable for standalone drop rules
	enable_standalone_drop_rules = true

   # Enable for modular drop rules
	enable_modular_drop_rules = true  
}  
```

To clarify, the difference between the two categories is straightforward:

- Any `newrelic_nrql_drop_rule` resource that is defined directly in its resource form and not wrapped in a module is a "standalone" resource. To enable operation on "standalone" drop rules using this Terraform script, set the value of `enable_standalone_drop_rules` in `locals{}` to `true`.
- Similarly, any `newrelic_nrql_drop_rule` resource that has its implementation wrapped inside a module, and is not a "standalone" `newrelic_nrql_drop_rule` resource, is a "modular" resource. To enable operation on "modular" drop rules using this Terraform script, set the value of `enable_modular_drop_rules` in `locals{}` to `true`.

The purpose of this step is to identify the types of drop rules in the state to which the procedure will be applied, and accordingly set the values of `enable_standalone_drop_rules` and `enable_modular_drop_rules` to `true` or `false`. With this configuration, you can proceed to the next step—2A, 2B, or both—based on your selection.

### Step 2A: Configure Standalone Drop Rules

If using standalone drop rules (remember that `enable_standalone_drop_rules` must be set to `true` in the previous step), populate the `standalone_rules` list with a collection of such resources, specifying the name and resource identifier of each resource.

For instance, if `newrelic_nrql_drop_rule.drop_health_checks` and `newrelic_nrql_drop_rule.drop_pii_data` are two standalone drop rules currently managed, add the following to the empty list in the variable `standalone_rules`, under the comment that specifies the location for such resources:

```hcl  
{
  name     = "drop_health_checks"
  resource = newrelic_nrql_drop_rule.drop_health_checks
},
{
  name     = "drop_pii_data"
  resource = newrelic_nrql_drop_rule.drop_pii_data
},
```

Note that the value of `resource` must correspond to the actual resource identifier of the `newrelic_nrql_drop_rule` resource as defined in the Terraform configuration. For drop rule resources created using the `count` or `for_each` meta-arguments, the resource identifiers must be explicitly indexed to reference a singular, well-defined, existing resource. For example: `newrelic_nrql_drop_rule.drop_health_checks[0]` or `newrelic_nrql_drop_rule.drop_health_checks_two["key_1"]`.

**Note**: Identifying all standalone `newrelic_nrql_drop_rule` resources in the state can be a significant effort, especially since this solution is intended for use in CI/CD environments. If state access is enabled and commands can be executed in the environment, the following command can help fetch a list of all standalone `newrelic_nrql_drop_rule` resources in the state:

- The `terraform` command may need to be replaced with an appropriate equivalent command, depending on the Terraform wrapper or GitOps environment being used. This is an experimental command intended to aid in the identification of standalone drop rules and may need adjustment based on the specific state structure and naming conventions in use.

```bash
(
  echo "#################### COPY FROM HERE ####################"
  terraform state list | \
    grep 'newrelic_nrql_drop_rule' | \
    grep -v '^module\.' | \
    sed -E \
      -e '/\.([^[]+)\["/ s/^.*\.([^[]+)\["([^"]+)"\].*$/  {\n    name     = "\1_\2"\n    resource = &\n  },/' \
      -e 't' \
      -e '/\.([^[]+)\[[0-9]/ s/^.*\.([^[]+)\[([0-9]+)\].*$/  {\n    name     = "\1_\2"\n    resource = &\n  },/' \
      -e 't' \
      -e 's/^.*\.([^ ]+)$/  {\n    name     = "\1"\n    resource = &\n  },/'
  echo "#################### COPY UNTIL HERE ###################"
)
```

### Step 2B: Configure Modular Drop Rules

If using modular drop rules (remember that `enable_modular_drop_rules` must be set to `true` in the previous step), populate the `_modular_rules_raw` list with a collection of such resources, specifying the name and resource identifier (module-based) of each resource. However, there is a mandatory, important prerequisite to consider to ensure this step works seamlessly.

**Prerequisite**: Before specifying the resources in the list, ensure that the modules we'd like to reference export an attribute `all_rules`, which holds references to all `newrelic_nrql_drop_rule` resources managed by the module. This is a mandatory requirement for the script to access the resources. For instance, if the `newrelic_nrql_drop_rule` resources in the module look like the following:

```hcl
resource "newrelic_nrql_drop_rule" "drop_metadata_drop_rule" {
  ..
  ..
}

resource "newrelic_nrql_drop_rule" "drop_sensitive_info_rule" {
  ..
  ..
}
```

An addition will need to be made to this part of the module to export an attribute `all_rules`, holding references to all drop rules comprised by the module; i.e., the following code would need to be added:

```hcl
output "all_rules" {
  description = "A map of all drop rule resource objects created by this module."
  value = {
    debug_logs     = newrelic_nrql_drop_rule.drop_metadata_drop_rule
    debug_logs_two = newrelic_nrql_drop_rule.drop_sensitive_info_rule
  }
}
```

After the export of `all_rules` with the desired syntax (as specified above) has been addressed, you can now populate the `_modular_rules_raw` list with a collection of such resources, specifying the name and "modular" resource identifier of each resource.

For example, if a module named `drop_rules` contains two `newrelic_nrql_drop_rule` resources, `drop_metadata_drop_rule` and `drop_sensitive_info_rule`, you would add the following entries to the `_modular_rules_raw` list under the designated section for modular resources (assuming the module has a `name = drop_metadata`):

```hcl
    {
      name     = "drop_metadata"
      resource = module.drop_rules["drop_metadata"].all_rules
    },
```

As stated in Step 2A, note that the value of `resource` must correspond to the actual resource identifier of the `newrelic_nrql_drop_rule`; since the resource identifier must be accessible by Terraform to find a valid resource and process further steps.

**Note**: Similar to the note in Step 2A, identifying all modular `newrelic_nrql_drop_rule` resources in the state can be a significant effort, especially since this solution is intended for use in CI/CD environments. If state access is enabled and commands can be executed in the environment, the following command can help fetch a list of all modular `newrelic_nrql_drop_rule` resources in the state:

- The `terraform` command may need to be replaced with an appropriate equivalent command, depending on the Terraform wrapper or GitOps environment being used. This is an experimental command intended to aid in the identification of modular drop rules and is likely to deviate from the expected result, based on the specific state structure and naming conventions in use.

```bash
(
  echo "#################### COPY FROM HERE ####################"
  terraform state list | \
    grep '^module\..*newrelic_nrql_drop_rule' | \
    sed -E 's/(\.newrelic_nrql_drop_rule.*)//' | \
    sort -u | \
    sed -E 's/^module\.([^[]+)\["([^"]+)"\]$/  {\n    name     = "\2"\n    resource = module.\1\["\2"].all_rules\n  },/'
  echo "#################### COPY UNTIL HERE ###################"
)
```

### Step 3: Plan and Apply Configuration

After configuring the standalone and/or modular drop rules as per the previous steps, add `outputs.tf` to the workspace in the CI/CD environment. Run `terraform plan` to validate the configuration. If the plan is successful and no validation errors are reported, proceed to apply the configuration using `terraform apply`. This will create or update resources as necessary and generate the required outputs for migration.

If you see a message with errors, refer to the "Troubleshooting" section at the end of this document for help resolving common issues. However, with all configurations done correctly, you should see a success message in the output upon `terraform plan`:

```
Changes to Outputs:
+ a_validation_success                          = "✅ All listed resources export pipeline_cloud_rule_entity_id"
```

This is confirmation that a `terraform apply` can now be performed. After the apply, the following will be generated as output:

```
Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

a_validation_success = "✅ All listed resources export pipeline_cloud_rule_entity_id"
experimental_drop_rule_resource_ids = "{\"drop_rule_resource_ids\":[{\"id\":\"1111111:222222222\",\"name\":\"drop_fake_logs\",\"pipeline_cloud_rule_entity_id\":\"MzgwNjUyNnxOR0VQfFBJUEVMSU5FX0NMT1VEX1JVTEV8MDE5OTRjZjctYmFmNy03MjU3LWE3M2MtZWY5OTkxYTQxMjgy\"},{\"id\":\"1111111:333333333\",\"name\":\"drop_fake_checks\",\"pipeline_cloud_rule_entity_id\":\"MzgwNjUyNnxOR0VQfFBJUEVMSU5FX0NMT1VEX1JVTEV8MDE5OTRjZmItMTQ0Yy03NDM5LWJhNDYtZjI4MTg0ODc5YmE2\"},{\"id\":\"1111111:444444444\",\"name\":\"drop_fake_data\",\"pipeline_cloud_rule_entity_id\":\"MzgwNjUyNnxOR0VQfFBJUEVMSU5FX0NMT1VEX1JVTEV8MDE5OTRjZmItMTQ4Ni03MDI4LWJlMDktZmYzOTM2NWQ4ODUw\"}]}"
experimental_drop_rule_resource_ids_formatted = <<EOT
{
  "drop_rule_resource_ids": [
    {
      "name": "drop_fake_logs",
      "id": "1111111:222222222",
      "pipeline_cloud_rule_entity_id": "MzgwNjUyNnxOR0VQfFBJUEVMSU5FX0NMT1VEX1JVTEV8MDE5OTRjZjctYmFmNy03MjU3LWE3M2MtZWY5OTkxYTQxMjgy"
    },
    {
      "name": "drop_fake_checks",
      "id": "1111111:333333333",
      "pipeline_cloud_rule_entity_id": "MzgwNjUyNnxOR0VQfFBJUEVMSU5FX0NMT1VEX1JVTEV8MDE5OTRjZmItMTQ0Yy03NDM5LWJhNDYtZjI4MTg0ODc5YmE2"
    },
    {
      "name": "drop_fake_data",
      "id": "1111111:444444444",
      "pipeline_cloud_rule_entity_id": "MzgwNjUyNnxOR0VQfFBJUEVMSU5FX0NMT1VEX1JVTEV8MDE5OTRjZmItMTQ4Ni03MDI4LWJlMDktZmYzOTM2NWQ4ODUw"
    }
    ]
}

EOT
```

Voilà! We now have `experimental_drop_rule_resource_ids_formatted`, which is the most important output and the desired outcome of this entire procedure. This enables us to supply the JSON value contained by `experimental_drop_rule_resource_ids_formatted` to the New Relic CLI command in the next phase of the overall migration process.

## What's Next?

The next phase of the migration process (which must be run locally) encompasses the JSON we've obtained from the procedure we've just performed. 

This JSON will be used with the New Relic CLI command `tf-importgen-ci`, which, with the help of Terraform on the local machine, will generate the configuration of `newrelic_pipeline_cloud_rule` resources corresponding to the Drop Rules identified and exported via the JSON. It will also generate an import script to import the resources into the state, followed by tailored recommendations to remove existing `newrelic_nrql_drop_rule` resources from the state and apply configuration to successfully complete the migration. 

For more details, refer to this page, which explains the operation and usage of the `tf-importgen-ci` command in the New Relic CLI.

## Validation System

The configuration includes built-in validation that checks:

1. **Rule Existence**: At least one drop rule must be configured
2. **Required Attributes**: Each drop rule must export `pipeline_cloud_rule_entity_id`

### Validation Outputs

- `a_validation_success`: Shows success message when all validations pass (the typo "a_" exists on purpose, to ensure this variable is shown alphabetically first in the outputs, when planned or applied)
- `validation_errors`: Lists specific validation errors when they occur

## Troubleshooting

### Common Issues

1. **No drop rules configured**
   - Error: "❌ No drop rules have been listed in the local variables above"
   - Solution: This error is thrown only when both of the lists in the local variables, `standalone_rules` and `_modular_rules_raw`, are empty. Add drop rules to the appropriate local variable lists: `standalone_rules` and/or `_modular_rules_raw`.

2. **Drop rules specified in variable lists, but not appearing in exported JSON**
   - This is a probable occurrence when the list in the local variable (standalone/modular) is populated with the rules in the correct format, but the variable needed to "unlock" the processing of these rules is not yet enabled and remains set to `false`.
   - Solution: Ensure that the `enable_` local variables, `enable_standalone_drop_rules` and/or `enable_modular_drop_rules`, are set to `true`, depending on which type of drop rules you want to process, when the corresponding lists (`standalone_rules` / `_modular_rules_raw` respectively) are updated.

3. **Missing pipeline_cloud_rule_entity_id / Provider version incompatibility**
   - Error: "❌ Drop rule 'rule_name' does not export pipeline_cloud_rule_entity_id" or attribute not found errors
   - Cause: Using New Relic Terraform Provider version below 3.72.1
   - Solution: Upgrade to New Relic Terraform Provider version 3.72.1 or above (preferably the latest version).

4. **Resource reference errors**
   - Error: "Resource 'newrelic_nrql_drop_rule.rule_name' does not exist" or similar Terraform reference errors
   - Cause: Incorrect resource identifiers in the configuration lists
   - Solution: Verify that all resource references match the actual resource names in your Terraform configuration. For resources using `count` or `for_each`, ensure proper indexing (e.g., `[0]` or `["key"]`).

5. **Module output not found**
   - Error: "Module 'module_name' does not have an output 'all_rules'" or similar module reference errors
   - Cause: Referenced modules do not export the required `all_rules` output
   - Solution: Add the `all_rules` output to your modules as specified in Step 2B prerequisites.

6. **Terraform plan/apply failures**
   - Error: Various Terraform execution errors during plan or apply
   - Cause: Configuration syntax errors, missing resources, or state inconsistencies
   - Solution: 
     - Ensure all referenced resources exist in your configuration
     - Ensure all elements of the populated lists, `standalone_rules` and/or `_modular_rules_raw`, comply with the syntax specified in this guide
     - Check that the workspace has proper permissions to access all referenced resources

7. **Empty or malformed JSON output**
   - Error: JSON output appears empty or contains null values despite validation passing
   - Cause: Issues with resource attribute access or conditional logic
   - Solution: 
     - Verify that drop rule resources are properly configured and applied
     - Check that the New Relic provider has successfully created the pipeline cloud rule counterparts
     - Ensure the workspace has the latest state after all resources have been applied