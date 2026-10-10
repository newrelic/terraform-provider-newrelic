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

### NRQL-based rule

```hcl
resource "newrelic_pipeline_cloud_rule" "nrql_example" {
  account_id  = 1000100
  name        = "Drop debug logs"
  description = "Drops all DEBUG logs from the dev environment."
  nrql        = "DELETE FROM Log WHERE logLevel = 'DEBUG' AND environment = 'dev'"
}
```

### OTTL-based rule

~> **NOTE:** OTTL rule creation requires backend support that is currently being rolled out. This example shows the intended configuration once the feature is fully available.

```hcl
resource "newrelic_pipeline_cloud_rule" "ottl_example" {
  account_id  = 1000100
  name        = "Redact PII from logs"
  description = "Uses OTTL to redact email addresses from log bodies."

  ottl_transform {
    log_statements = [
      "replace_pattern(body, \"[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\\\.[a-zA-Z]{2,}\", \"[REDACTED]\")",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

*   `account_id` - (Optional) The account ID where the Pipeline Cloud Rule will be created.
*   `name` - (Required) The name of the rule. This must be unique within an account.
*   `description` - (Optional) Additional information about the rule.
*   `nrql` - (Optional) The NRQL query that defines the data to be processed by this Pipeline Cloud Rule. Mutually exclusive with `ottl_transform`; exactly one of `nrql` or `ottl_transform` must be set.
*   `ottl_transform` - (Optional) OTTL transformation statements for non-NRQL pipeline cloud rules. Mutually exclusive with `nrql`; exactly one of `nrql` or `ottl_transform` must be set. See [OTTL Transform](#ottl-transform) below for details.

### OTTL Transform

The `ottl_transform` block supports the following arguments. Exactly one of the statement fields must be specified, scoped to a single telemetry type:

*   `log_statements` - (Optional) List of OTTL statements applied to log data.
*   `event_statements` - (Optional) List of OTTL statements applied to event data.
*   `metric_statements` - (Optional) List of OTTL statements applied to metric data.
*   `trace_statements` - (Optional) List of OTTL statements applied to trace data.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   `id` - The ID of the Pipeline Cloud Rule.

## Import

Pipeline Cloud Rules can be imported using the `id`. For example:

```bash
$ terraform import newrelic_pipeline_cloud_rule.foo <id>
```

-> **NOTE:** If you'd like to import a `newrelic_pipeline_cloud_rule` resource corresponding to an existing `newrelic_nrql_drop_rule` resource in your configuration in light of the aforementioned EOL, please head over to the [instructions in our Drop Rules EOL Migration Guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/drop_rules_eol_guide#alternatives-and-action-needed).