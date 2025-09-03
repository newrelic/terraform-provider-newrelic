---
layout: "newrelic"
page_title: "NRQL Drop Rules EOL (Upcoming): Implications and Actions Needed 📢"
sidebar_current: "docs-newrelic-provider-drop-rules-eol-migration-guide"
description: |-
  Use this guide to find details on the end-of-life of NRQL Drop Rules, implications seen by customers maintaining NRQL Drop Rule resources via the New Relic Terraform Provider, and actions to be taken prior to the EOL to avoid consequences.
---
## NRQL Drop Rules EOL: Implications and Actions Needed 🚨

### About the EOL

As announced by New Relic ([see this announcement](https://docs.newrelic.com/eol/2025/05/drop-rule-filter/)), the <b style="color:red;">end-of-life (EOL)</b> for the **Drop Rules API** will take effect on <b style="color:red;">January 7, 2026</b>. Consequently, support for managing drop rules via the New Relic Terraform Provider's `newrelic_nrql_drop_rule` resource <b style="color:maroon;">will also officially end on January 7, 2026</b>. After the EOL is effective, all API requests made by the Terraform provider using the `newrelic_nrql_drop_rule` resource <span style="color:red;">will be blocked and result in an API error</span>.

In line with these changes, the `newrelic_nrql_drop_rule` resource <span style="color:red;">has been marked as <b>deprecated</b></span> starting with <b style="color:red;">v3.67.0</b> of the New Relic Terraform Provider. <span style="color:red;">It will be <b>removed</b> from the provider in a future release coinciding with the January 7, 2026 EOL</span>. This means the resource can no longer be used to create or manage drop rules after this date.

### Alternatives and Action Needed

NRQL Drop Rules are being replaced by **Pipeline Cloud Rules**. See [this article](https://docs.newrelic.com/docs/new-relic-control/pipeline-control/cloud-rules-api/) for an overview.

New Relic will handle the upstream migration of existing NRQL Drop Rules to Pipeline Cloud Rules. However, to continue managing these NRQL Drop Rules as Pipeline Cloud Rules, in their new definition (and any new Pipeline Cloud Rules) via Terraform, <span style="color:tomato;">customers using the `newrelic_nrql_drop_rule` resource must transition to the new `newrelic_pipeline_cloud_rule` resource.</span> Please see [this page](/providers/newrelic/newrelic/latest/docs/r/pipeline_cloud_rule) for documentation on using the `newrelic_pipeline_cloud_rule` resource.

To transition to the `newrelic_pipeline_cloud_rule` resource for rules already migrated upstream by New Relic, you will need to:
- _Import_ the existing Pipeline Cloud Rules (which were formerly NRQL Drop Rules) into your Terraform state as `newrelic_pipeline_cloud_rule` resources using the `terraform import` command with the ID(s) of Pipeline Cloud Rules. See [this page](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule#import) for details on how to import a Pipeline Cloud Rule. 
  - To simplify this process, you can pair the import (in a different form) with `terraform plan -generate-config-out` to automatically generate the corresponding resource configuration, as explained [here, in the Terraform docs](https://developer.hashicorp.com/terraform/language/import/generating-configuration).
  ```hcl
  import {
    to = newrelic_pipeline_cloud_rule.foo
    id = "MXxXX0XXxXXXXXXXXX5XX0XXX1XXX1XXXXX8XXX5XXXxX2XxXxxxXX03XXX2XXx0XXXxXXXxXxXxXXXxXXXx"
  }
  ```
 - _Remove_ all references to the `newrelic_nrql_drop_rule` resources from the Terraform state (after successfully importing them as `newrelic_pipeline_cloud_rule` resources) using the `terraform state rm` command. See [this page](https://developer.hashicorp.com/terraform/cli/commands/state/rm) for details on removing items from the Terraform state.

The process outlined above is our recommendation for migrating to Pipeline Cloud Rules. We are exploring ways to assist with automating this migration in certain scenarios and will share updates and resources here as they become available.