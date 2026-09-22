---
layout: "newrelic"
page_title: "New Relic: newrelic_pipeline_cloud_rule"
sidebar_current: "docs-newrelic-pipeline-cloud-rule"
description: |-
  Use this resource to create and manage a New Relic Pipeline Cloud Rule.
---

# Resource: newrelic\_pipeline\_cloud\_rule

Use this resource to create and manage a New Relic Pipeline Cloud Rule.

-> **NOTE:** **Starting v3.68.0 of the New Relic Terraform Provider**, <b style="color:green;">Pipeline Cloud Rules can be managed using the [`newrelic_pipeline_cloud_rule`](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule) resource.</b> This resource is the designated replacement for the <span style="color:red;">now end-of-life [`newrelic_nrql_drop_rule`](/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule) resource</span>, which reached its end-of-life on <b style="color:red;">August 31, 2026</b> and will be removed from the provider in an upcoming release. <br><br>If you are currently using the [`newrelic_nrql_drop_rule`](/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule) resource, please migrate to [`newrelic_pipeline_cloud_rule`](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule) at the earliest. Refer to the [Drop Rules EOL Migration Guide](/providers/newrelic/newrelic/latest/docs/guides/drop_rules_eol_guide) for step-by-step instructions.<br><br>

## Example Usage

```hcl
resource "newrelic_pipeline_cloud_rule" "foo" {
  account_id  = 1000100
  name        = "Test Pipeline Cloud Rule"
  description = "This rule deletes all DEBUG logs from the dev environment."
  nrql        = "DELETE FROM Log WHERE logLevel = 'DEBUG' AND environment = 'dev'"
}
```

## Argument Reference

The following arguments are supported:

*   `account_id` - (Optional) The account ID where the Pipeline Cloud Rule will be created.
*   `name` - (Required) The name of the rule. This must be unique within an account.
*   `nrql` - (Required) The NRQL query that defines the data to be processed by this Pipeline Cloud Rule.
*   `description` - (Optional) Additional information about the rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   `id` - The ID of the Pipeline Cloud Rule.

## Import

Pipeline Cloud Rules can be imported using the `id`. For example:

```bash
$ terraform import newrelic_pipeline_cloud_rule.foo <id>
```

-> **NOTE:** If you'd like to import a `newrelic_pipeline_cloud_rule` resource corresponding to an existing `newrelic_nrql_drop_rule` resource in your configuration in light of the aforementioned EOL, please head over to the [instructions in our Drop Rules EOL Migration Guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/drop_rules_eol_guide#alternatives-and-action-needed).