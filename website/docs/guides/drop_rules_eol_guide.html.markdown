---
layout: "newrelic"
page_title: "🚨 NRQL Drop Rules EOL (Upcoming): Implications and Actions Needed"
sidebar_current: "docs-newrelic-provider-drop-rules-eol-migration-guide"
description: |-
  Use this guide to find details on the end-of-life of NRQL Drop Rules, implications seen by customers maintaining NRQL Drop Rule resources via the New Relic Terraform Provider, and actions to be taken prior to the EOL to avoid consequences.
---
## NRQL Drop Rules EOL: Implications and Actions Needed 🚨

### About the EOL

As announced by New Relic ([see this announcement](https://docs.newrelic.com/eol/2025/05/drop-rule-filter/)), the <b style="color:red;">end-of-life (EOL)</b> for the **Drop Rules API** will take effect on <b style="color:red;">January 7, 2026</b>. Consequently, support for managing drop rules via the New Relic Terraform Provider's [`newrelic_nrql_drop_rule`](/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule) resource <b style="color:maroon;">will also officially end on January 7, 2026</b>. After the EOL is effective, all API requests made by the Terraform provider using the [`newrelic_nrql_drop_rule`](/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule) resource <span style="color:red;">will be blocked and result in an API error</span>.

In line with these changes, the [`newrelic_nrql_drop_rule`](/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule) resource <span style="color:red;">has been marked as <b>deprecated</b></span> starting with <b style="color:red;">v3.68.0</b> of the New Relic Terraform Provider. <span style="color:red;">It will be <b>removed</b> from the provider in a future release coinciding with the January 7, 2026 EOL</span>. This means the resource can no longer be used to create or manage drop rules after this date.

### Alternatives and Action Needed

NRQL Drop Rules are being replaced by **Pipeline Cloud Rules**. See [this article](https://docs.newrelic.com/docs/new-relic-control/pipeline-control/cloud-rules-api/) for an overview.

New Relic will handle the upstream migration of existing NRQL Drop Rules to Pipeline Cloud Rules. However, to continue managing these NRQL Drop Rules as Pipeline Cloud Rules, in their new definition (and any new Pipeline Cloud Rules) via Terraform, <span style="color:tomato;">customers using the [`newrelic_nrql_drop_rule`](/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule) resource must transition to the new [`newrelic_pipeline_cloud_rule`](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule) resource.</span> Please see [this page](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule) for documentation on using the [`newrelic_pipeline_cloud_rule`](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule) resource.

To transition to the [`newrelic_pipeline_cloud_rule`](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule) resource for rules already migrated upstream by New Relic, you will need to:
- **_Import_** the existing Pipeline Cloud Rules (which were formerly NRQL Drop Rules) into your Terraform state as [`newrelic_pipeline_cloud_rule`](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule) resources using the `terraform import` command with the ID(s) of Pipeline Cloud Rules. See [this page](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule#import) for details on how to import a Pipeline Cloud Rule. 
  - To simplify this process, you can pair the import (in a different form) with `terraform plan -generate-config-out` to automatically generate the corresponding resource configuration, as explained [here, in the Terraform docs](https://developer.hashicorp.com/terraform/language/import/generating-configuration).
  - <span style="color:tomato;">Note that the ID of a Pipeline Cloud Rule is _not the same_ as the ID of the corresponding NRQL Drop Rule.</span>
  - To import a Pipeline Cloud Rule (as a [`newrelic_pipeline_cloud_rule`](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule) resource) corresponding to an existing [`newrelic_nrql_drop_rule`](/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule) resource in your configuration, upgrade to **v3.68.0** or greater of the New Relic Terraform Provider, refresh and apply your configuration comprising the [`newrelic_nrql_drop_rule`](/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule) resources, which would allow an argument [`pipeline_cloud_rule_entity_id`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule#pipeline_cloud_rule_entity_id-1) to be exported by the updated [`newrelic_nrql_drop_rule`](/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule) resources. 
  - The value of this argument, e.g. `newrelic_nrql_drop_rule.foo.pipeline_cloud_rule_entity_id` can then be used as the ID to import the corresponding Pipeline Cloud Rule, to import and generate configuration as stated in the above steps and the example below, or with the `terraform import` command too.
  ```hcl
  import {
    to = newrelic_pipeline_cloud_rule.foo
    id = "MXxXX0XXxXXXXXXXXX5XX0XXX1XXX1XXXXX8XXX5XXXxX2XxXxxxXX03XXX2XXx0XXXxXXXxXxXxXXXxXXXx"
    
    # id = newrelic_nrql_drop_rule.foo.pipeline_cloud_rule_entity_id
    
    # The ID of a Pipeline Cloud Rule (to be imported
    # as a `newrelic_pipeline_cloud_rule` resource) 
    # corresponding to an existing `newrelic_nrql_drop_rule` resource
    # can be exported from the `newrelic_nrql_drop_rule` resource 
    # from the argument `pipeline_cloud_rule_entity_id`, 
    # when on version >= 3.68.0 of the provider.
    # See the note above (in the migration guide) for details.
  }
  ```
 - **_Remove_** all references to the [`newrelic_nrql_drop_rule`](/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule) resources from the Terraform state (after successfully importing them as [`newrelic_pipeline_cloud_rule`](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule) resources) using the `terraform state rm` command. See [this page](https://developer.hashicorp.com/terraform/cli/commands/state/rm) for details on removing items from the Terraform state.

The process outlined above is our recommended approach for migrating to Pipeline Cloud Rules. However, the following sections describe automation helpers that can simplify the migration in some specific scenarios or environments.


### Automated Migration for CI/CD Environments

For users managing `newrelic_nrql_drop_rule` resources in CI/CD environments (such as Atlantis, Grandcentral, or similar GitOps workflows), we provide automation helpers that could help streamline the migration process through a **three-phase approach** in some CI/CD environments:

**Phase 1 - Export Drop Rule Data (CI/CD Environment):**
- Modify your CI/CD workspace configuration to identify and export existing drop rules as JSON data
- Uses automation scripts that validate and extract drop rule information including their Pipeline Cloud Rule counterparts
- Requires New Relic Terraform Provider version 3.72.1 or above
- See the [Phase 1 Usage Guide](https://github.com/newrelic/terraform-provider-newrelic/blob/main/examples/drop_rule_migration_ci) for detailed instructions

**Phase 2 - Generate Pipeline Cloud Rule Configurations (Local Environment):**
- Process the exported JSON data using the New Relic CLI command `tf-importgen-ci`
- Automatically generates `newrelic_pipeline_cloud_rule` configurations and import scripts
- **Executed locally using the New Relic CLI**
- See the [Phase 2 Usage Guide](https://github.com/newrelic/newrelic-cli/blob/main/internal/migrate/tf_importgen_ci_guide.md) for detailed instructions

**Phase 3 - Apply Configurations and Complete Migration (CI/CD Environment):**
- Apply the generated configurations back in the CI/CD environment
- Import the new Pipeline Cloud Rules into Terraform state
- Remove legacy NRQL Drop Rules from management

This three-phase automation process is specifically designed for CI/CD workflows where direct local access to Terraform state may be limited. It provides a structured approach to migrate from `newrelic_nrql_drop_rule` to `newrelic_pipeline_cloud_rule` resources while maintaining the integrity of your GitOps processes.

For complete documentation and step-by-step instructions for Phase 1, refer to the [automation helper usage guide](https://github.com/newrelic/terraform-provider-newrelic/blob/main/examples/drop_rule_migration_ci/). Instructions for Phases 2 and 3 are provided through the New Relic CLI `tf-importgen-ci` [command documentation](https://github.com/newrelic/newrelic-cli/blob/main/internal/migrate/tf_importgen_ci_guide.md).