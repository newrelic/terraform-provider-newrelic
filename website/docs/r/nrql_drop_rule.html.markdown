---
layout: "newrelic"
page_title: "New Relic: newrelic_nrql_drop_rule"
sidebar_current: "docs-newrelic-resource-nrql-drop-rule"
description: |-
  This resource has reached end-of-life and is no longer supported. Please refer to the migration guide for alternatives.
---
# Resource: newrelic\_nrql\_drop\_rule

!> **This resource is no longer supported.** <span style="color:red;"><b>The `newrelic_nrql_drop_rule` resource reached its end-of-life on August 31, 2026</b> and is no longer functional.</span> It will be removed from the New Relic Terraform Provider in an upcoming release.<br><br>If you are looking for documentation on the arguments, attributes, or usage examples for this resource, please refer to an older release of the New Relic Terraform Provider (v3.97.0 or earlier) in the [Terraform Registry release history](https://registry.terraform.io/providers/newrelic/newrelic/).<br><br>New Relic has completed the upstream migration of all existing Drop Rules to **Pipeline Cloud Rules**. To continue managing these rules via Terraform, update your configurations to use the [`newrelic_pipeline_cloud_rule`](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule) resource. For `newrelic_nrql_drop_rule` resources with `action = "drop_attributes_from_metric_aggregates"`, the recommended alternative is the [`newrelic_metric_pruning_rule`](/providers/newrelic/newrelic/latest/docs/resources/metric_pruning_rule) resource.<br><br>Please refer to the [Drop Rules EOL Migration Guide](/providers/newrelic/newrelic/latest/docs/guides/drop_rules_eol_guide) for detailed migration instructions.
