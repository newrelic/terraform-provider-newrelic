---
layout: "newrelic"
page_title: "New Relic: newrelic_pipeline_cloud_rule"
sidebar_current: "docs-newrelic-pipeline-cloud-rule"
description: |-
  Use this resource to create and manage a New Relic Pipeline Cloud Rule.
---

# Resource: newrelic\_pipeline\_cloud\_rule

Use this resource to create and manage a New Relic Pipeline Cloud Rule.

-> **❗<b style="color:green;">\*NEW\*</b>** **Starting v3.68.0 of the New Relic Terraform Provider**, <b style="color:green;">Pipeline Cloud Rules can be managed using the resource [`newrelic_pipeline_cloud_rule`](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule).</b> This resource replaces the <span style="color:red;">deprecated [`newrelic_nrql_drop_rule`](/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule) resource</span>. <br><br><b>For customers currently managing Drop Rules with the deprecated [`newrelic_nrql_drop_rule`](/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule) resource:</b> Please see our [migration guide](/providers/newrelic/newrelic/latest/docs/guides/drop_rules_eol_guide) for instructions on switching to the [`newrelic_pipeline_cloud_rule`](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule) resource. <span style="color:red;">The resource [`newrelic_nrql_drop_rule`](/providers/newrelic/newrelic/latest/docs/resources/nrql_drop_rule) is <b>deprecated</b> and will be removed on <b>January 7, 2026</b></span>. While New Relic has automatically migrated your Drop Rules to Pipeline Cloud Rules upstream, <span style="color:tomato;">you must update your Terraform configuration to continue managing Drop Rules as Pipeline Cloud Rules</span>, using the <b style="color:green;">new</b> [`newrelic_pipeline_cloud_rule`](/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule) resource.<br><br>

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