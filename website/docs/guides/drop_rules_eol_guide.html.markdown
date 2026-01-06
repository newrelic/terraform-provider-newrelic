---
layout: "newrelic"
page_title: "🚨 NRQL Drop Rules EOL (Upcoming): Implications and Actions Needed"
sidebar_current: "docs-newrelic-provider-drop-rules-eol-migration-guide"
description: |-
  Use this guide to find details on the end-of-life of NRQL Drop Rules, implications seen by customers maintaining NRQL Drop Rule resources via the New Relic Terraform Provider, and actions to be taken prior to the EOL to avoid consequences.
---
## NRQL Drop Rules EOL: Implications and Actions Needed 🚨

### About the EOL

As announced by New Relic ([see this announcement](https://docs.newrelic.com/eol/2025/05/drop-rule-filter/)), the <b style="color:red;">end-of-life (EOL)</b> for the **Drop Rules API** will take effect on <b style="color:red;">June 30, 2026</b>. Consequently, support for managing drop rules via the New Relic Terraform Provider's [`newrelic_nrql_drop_rule`](/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule) resource <b style="color:maroon;">will also officially end on June 30, 2026</b>. After the EOL is effective, all API requests made by the Terraform provider using the [`newrelic_nrql_drop_rule`](/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule) resource <span style="color:red;">will be blocked and result in an API error</span>.

In line with these changes, the [`newrelic_nrql_drop_rule`](/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule) resource <span style="color:red;">has been marked as <b>deprecated</b></span> starting with <b style="color:red;">v3.68.0</b> of the New Relic Terraform Provider. <span style="color:red;">It will be <b>removed</b> from the provider in a future release coinciding with the June 30, 2026 EOL</span>. This means the resource can no longer be used to create or manage drop rules after this date.

### Alternatives and Action Needed

NRQL Drop Rules are being replaced by **Pipeline Cloud Rules**. See [this article](https://docs.newrelic.com/docs/new-relic-control/pipeline-control/cloud-rules-api/) for an overview.

New Relic will handle the upstream migration of existing NRQL Drop Rules to Pipeline Cloud Rules. However, to continue managing these rules via Terraform after the EOL, <span style="color:tomato;">customers using the [`newrelic_nrql_drop_rule`](/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule) resource must transition to the new [`newrelic_pipeline_cloud_rule`](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule) resource.</span> Please see [this page](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule) for documentation on using the [`newrelic_pipeline_cloud_rule`](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule) resource.

<hr/>

### Manual Migration Process Overview

To transition from `newrelic_nrql_drop_rule` to `newrelic_pipeline_cloud_rule` resources, follow this step-by-step process:

**Step 1: Upgrade Your New Relic Provider**

Update your Terraform configuration to use New Relic Terraform Provider **version 3.68.0 or higher**. This version adds support for the `pipeline_cloud_rule_entity_id` attribute needed for migration.

```hcl
terraform {
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
      version = "~> 3.68"  # Ensure version 3.68.0 or higher
    }
  }
}
```

**Step 2: Update Your Drop Rule Resources**

Run `terraform apply` to refresh your existing `newrelic_nrql_drop_rule` resources. This operation updates the Terraform state by adding the `pipeline_cloud_rule_entity_id` attribute against each `newrelic_nrql_drop_rule` resource targeted during the apply, which contains the ID of the corresponding Pipeline Cloud Rule that New Relic automatically created during their backend migration. This attribute is only added to the state file and doesn't modify your actual infrastructure or configuration files.

To observe this state update in your Terraform operation logs, use the `-refresh-only` flag to perform a state refresh without making any infrastructure changes. For a safer approach when applying changes, use the `-target` flag to limit the apply operation specifically to drop rule resources, avoiding unintended modifications to other resources in your configuration.

```bash
# Apply to refresh and update state with new pipeline_cloud_rule_entity_id attribute
terraform apply

# For safer approach: use refresh-only to see the pipeline_cloud_rule_entity_id 
# attribute values without making infrastructure changes (recommended for visibility)
terraform apply -refresh-only

# For targeted approach: limit the operation to specific drop rule resources
# to avoid unintended modifications to other resources in your configuration
terraform apply -refresh-only -target=newrelic_nrql_drop_rule.example
```

After this step, your existing drop rule resource will have the new attribute in state:

```hcl
resource "newrelic_nrql_drop_rule" "example" {
  account_id = var.account_id
  name       = "Drop high volume logs"
  action     = "drop_data"
  nrql       = "SELECT * FROM Log WHERE severity = 'DEBUG'"
  
  # This attribute gets populated automatically in state after refresh
  # pipeline_cloud_rule_entity_id = "MXxXXX18XXXXXXlXXXXXX058Xxx4XXxxXXXx"
}
```

<span style="color:tomato;">**Note**: The Pipeline Cloud Rule ID is different from the original NRQL Drop Rule ID.</span>

**Step 3: Import the Pipeline Cloud Rules**

Use the `pipeline_cloud_rule_entity_id` values exported by `newrelic_nrql_drop_rule` resources to import the corresponding Pipeline Cloud Rules into Terraform. Create import blocks and generate resource configurations using `terraform plan -generate-config-out` (as explained [here](https://developer.hashicorp.com/terraform/language/import/generating-configuration)).

```hcl
# Create import block using the pipeline_cloud_rule_entity_id from state
# For each NRQL Drop Rule to be migrated to a Pipeline Cloud Rule,
# an import{} block would be needed
import {
  to = newrelic_pipeline_cloud_rule.example
  # reference the "id" from the drop rule resource
  id = newrelic_nrql_drop_rule.example.pipeline_cloud_rule_entity_id 
  
  # or, use the actual value from your state
  # id = "MXxXXX18XXXXXXlXXXXXX058Xxx4XXxxXXXx"  
}
```

Then run Terraform commands to generate configuration and import:

```bash
# Generate Pipeline Cloud Rule configuration automatically
terraform plan -generate-config-out=generated_pipeline_rules.tf

# Apply to import the Pipeline Cloud Rules
terraform apply
```

**Step 4: Remove Legacy Drop Rules from State**

After successfully importing Pipeline Cloud Rules, remove the old `newrelic_nrql_drop_rule` resources from Terraform state and comment out their configuration blocks.

**Important:** <b style="color:red;">Do not use `terraform destroy` operations on `newrelic_nrql_drop_rule` resources, as this will also destroy the corresponding Pipeline Cloud Rules in the backend</b>. Instead, remove these resources from Terraform state using `terraform state rm` operations, then comment out or delete all `newrelic_nrql_drop_rule` resource configurations from your Terraform files.

```bash
# Remove the old drop rule from state (use your actual resource name)
terraform state rm newrelic_nrql_drop_rule.example
```

Then comment out or remove the `newrelic_nrql_drop_rule` resource blocks from your `.tf` files:

```hcl
# Comment out the old drop rule resource
# resource "newrelic_nrql_drop_rule" "example" {
#   account_id = var.account_id
#   name       = "Drop high volume logs"
#   action     = "drop_data"
#   nrql       = "SELECT * FROM Log WHERE severity = 'DEBUG'"
# }

# Your new Pipeline Cloud Rule resource (generated in Step 3)
resource "newrelic_pipeline_cloud_rule" "example" {
  account_id = var.account_id
  name       = "Drop high volume logs"
  action     = "drop_data"
  nrql       = "SELECT * FROM Log WHERE severity = 'DEBUG'"
  # Additional Pipeline Cloud Rule specific attributes...
}
```

The process outlined above is the recommended manual approach for migrating to Pipeline Cloud Rules. However, the following sections describe automation helpers that can simplify this migration process in specific scenarios or environments.

<hr/>

### Automated Migration for CI/CD Environments

For users managing `newrelic_nrql_drop_rule` resources in CI/CD environments (such as Atlantis, Grandcentral, or similar GitOps workflows **where direct local access to Terraform state may be limited**), we provide automation helpers that streamline the migration process through a **three-phase approach** specifically designed for CI/CD workflows. This approach utilizes specialized automation scripts and the **New Relic CLI** `tf-importgen-ci` command to handle the complexities of migrating drop rules in environments where state access is restricted or managed through automated pipelines.

**Prerequisites:**
- **Required**: Install the latest version of the **New Relic CLI** locally for Phase 2 processing. For installation instructions, see the [New Relic CLI documentation](https://docs.newrelic.com/docs/new-relic-solutions/tutorials/new-relic-cli/)
- **Required**: Ensure your CI/CD environment uses New Relic Terraform Provider version 3.73.0 or higher for Phase 1 data export
- **Required**: Verify CI/CD environment uses Terraform version 1.5 or higher for `import` block support in Phase 3

**Phase 1 - Export Drop Rule Data (CI/CD Environment):**

Execute the Phase 1 automation script in your CI/CD environment as documented in the Phase 1 guide.
<span style="color:green;">Upon successful completion of Phase 1, JSON data containing drop rule information will be exported and ready for Phase 2 processing.</span>

**Important Note:** For this phase to function correctly, ensure that both the New Relic Terraform Provider and Terraform versions used in your workspace meet the version requirements specified in the Prerequisites section above.

Briefly, this is the sequence of operations one would need to perform during this phase:
- Modify your CI/CD workspace configuration to identify and export existing drop rules as JSON data
- Copy the specified automation scripts that validate and extract drop rule information including their Pipeline Cloud Rule counterparts

See the [Phase 1 Usage Guide](https://github.com/newrelic/terraform-provider-newrelic/blob/main/examples/drop_rule_migration_ci) for detailed instructions.

**Phase 2 - Generate Pipeline Cloud Rule Configurations (Local Environment):**

Run the following **New Relic CLI** command locally using the exported JSON data from Phase 1.
<span style="color:green;">Upon successful command execution, generated Terraform configurations will be ready for Phase 3 deployment.</span>
```bash
newrelic migrate nrqldroprules tf-importgen-ci --file drop_rules.json
```

The `tf-importgen-ci` command processes the JSON data exported from Phase 1 and automatically generates `newrelic_pipeline_cloud_rule` configurations along with their corresponding import scripts for the Drop Rules identified in the exported data.

**Important**: This phase must be executed in a **local, empty workspace** and not in the CI/CD environment, unlike Phases 1 and 3.

For detailed instructions, see the [Phase 2 Usage Guide](https://github.com/newrelic/newrelic-cli/blob/main/internal/migrate/tf_importgen_ci_guide.md).

**Phase 3 - Apply Configurations and Complete Migration (CI/CD Environment):**

Deploy the generated configurations from Phase 2 back to your CI/CD environment following the provided integration instructions.
<span style="color:green;">Upon successful completion of Phase 3, the migration from NRQL Drop Rules to Pipeline Cloud Rules will be complete.</span>

Briefly, this is the sequence of operations one would need to perform during this phase:
- Apply the generated configurations back in the CI/CD environment
- Import the new Pipeline Cloud Rules into Terraform state
- Remove legacy NRQL Drop Rules from management

Detailed instructions for Phase 3 are provided at the end of the logs when the `tf-importgen-ci` command is executed in Phase 2.

**Congratulations!** 🎉 With all three phases executed successfully according to the recommended steps and all post-execution tasks completed (including removal of legacy Drop Rule resource configurations from your Terraform files), your NRQL Drop Rules have been fully migrated to Pipeline Cloud Rules. Your CI/CD environment should now be managing the migrated drop rules as `newrelic_pipeline_cloud_rule` resources, ensuring continued functionality beyond the June 30, 2026 EOL date.

This three-phase automation process is specifically designed for CI/CD workflows where direct local access to Terraform state may be limited. It provides a structured approach to migrate from `newrelic_nrql_drop_rule` to `newrelic_pipeline_cloud_rule` resources while maintaining the integrity of your GitOps processes.

For complete documentation and step-by-step instructions for Phase 1, refer to the [New Relic Terraform Provider: `newrelic_nrql_drop_rule` -> `newrelic_pipeline_cloud_rule` Migration Guide for CI/CD Workflows Using Automation Helpers](https://github.com/newrelic/terraform-provider-newrelic/blob/main/examples/drop_rule_migration_ci/). Instructions for Phases 2 and 3 are provided through the **New Relic CLI** `tf-importgen-ci` [command documentation](https://github.com/newrelic/newrelic-cli/blob/main/internal/migrate/tf_importgen_ci_guide.md). This guide provides detailed command usage, examples, troubleshooting, and best practices for each phase of the CI/CD migration process.

<hr/>

### Interactive Automated Migration for Local Terraform Workspaces

For users managing `newrelic_nrql_drop_rule` resources in **local Terraform workspaces with direct state access** (as opposed to CI/CD environments where state access may be restricted), we provide interactive command-line tools that streamline the migration process through a **step-by-step approach** for local development environments. This approach utilizes three specialized **New Relic CLI** commands: `tf-update`, `tf-importgen`, and `tf-delist`, each designed to handle a specific phase of the migration process with real-time feedback and validation.

**Prerequisites:**
- **Required**: Install the latest version of the **New Relic CLI** and ensure it's accessible in your PATH. For installation instructions, see the [New Relic CLI documentation](https://docs.newrelic.com/docs/new-relic-solutions/tutorials/new-relic-cli/)
- **Required**: Verify your Terraform workspace containing `newrelic_nrql_drop_rule` resources uses New Relic Terraform Provider version 3.73.0 or higher and Terraform version 1.5 or higher to ensure compatibility with the migration commands
- **Recommended**: Navigate to your Terraform workspace directory where your `newrelic_nrql_drop_rule` resources are defined before running these commands
- Ensure you have direct access to your Terraform state files and configuration

-> **NOTE** The following commands perform `terraform apply` operations with the `--auto-approve` flag enabled by default. This means that Terraform will automatically proceed with applying changes without requiring manual confirmation. When prompted by the command's user interface to proceed with a `terraform apply` operation, carefully review the proposed changes before confirming with a "Y", as the subsequent apply will execute automatically without additional approval prompts.

**Step 1 - Update Drop Rule Resources:**

Run the following **New Relic CLI** command in your Terraform workspace. <span style="color:green;">Upon successful command execution, you may proceed to the next step.</span>
```bash
newrelic migrate nrqldroprules tf-update
```
Briefly, this is the sequence of operations executed by the command:
- Refresh existing NRQL drop rule resources in Terraform state to populate `pipeline_cloud_rule_entity_id` values
- Uses `terraform apply -refresh-only` commands to update state with Pipeline Cloud Rule counterparts
- Requires New Relic Terraform Provider version 3.73.0 or above

**Step 2 - Generate and Execute Import Configuration:**

Run the following **New Relic CLI** command in your Terraform workspace. <span style="color:green;">Upon successful command execution, you may proceed to the next step.</span>
```bash
newrelic migrate nrqldroprules tf-importgen
```
Briefly, this is the sequence of operations executed by the command:
- Generate Terraform import blocks for Pipeline Cloud Rules using the **New Relic CLI** command `tf-importgen`
- Automatically creates import configuration and executes `terraform plan -generate-config-out` and `terraform apply`

**Step 3 - Safely Remove Legacy Resources:**

Run the following **New Relic CLI** command in your Terraform workspace. <span style="color:green;">Upon successful command execution, you may proceed to the next step.</span>
```bash
newrelic migrate nrqldroprules tf-delist
```
Briefly, this is the sequence of operations executed by the command:
- Remove NRQL drop rule resources from Terraform state without destroying actual resources in New Relic
- Uses `terraform state rm` commands to delist resources while keeping drop rules active
- Provides comprehensive post-migration cleanup instructions for configuration files

**Congratulations!** 🎉 With the three steps above successfully completed and the post-command actions taken (such as cleaning up configuration files containing Drop Rule resources), all of your Drop Rule resources have been seamlessly migrated to Pipeline Cloud Rule resources!

This step-by-step interactive process is specifically designed for local development workflows where you have direct access to Terraform state and configuration files. It provides granular control over each migration step with real-time feedback and the ability to validate changes before proceeding to the next phase.

For complete documentation and step-by-step instructions, refer to the [New Relic CLI: `newrelic_nrql_drop_rule` -> `newrelic_pipeline_cloud_rule` Migration Guide for Local Terraform Workspaces](https://github.com/newrelic/newrelic-cli/blob/main/internal/migrate/tf_importgen_guide.md). This guide provides detailed command usage, examples, troubleshooting, and best practices for each step of the local migration process. Additionally, [this section in the guide](https://github.com/newrelic/newrelic-cli/blob/main/internal/migrate/tf_importgen_guide.md#quick-start-three-command-workflow) is useful to find some quick details on the usage of this command with key options, while [this section](https://github.com/newrelic/newrelic-cli/blob/main/internal/migrate/tf_importgen_guide.md#complete-migration-workflow) could be useful to obtain an understanding of the migration procedure step by step with these commands.