---
layout: "newrelic"
page_title: "New Relic: newrelic_test_grok_pattern"
sidebar_current: "docs-newrelic-datasource-test-grok"
description: |-
Looks up if the given Grok pattern is matched against the log lines in New Relic.
---

# Data Source: newrelic\_test\_grok\_pattern

Use this data source to validate a grok pattern.  More information on Terraform's data sources can be found [here](https://www.terraform.io/language/data-sources).

## Example Usage

```hcl
# Data source
data "newrelic_test_grok_pattern" "foo" {
  grok = "%%{IP:host_ip}"
  log_lines = ["host_ip: 43.3.120.2","bytes_received: 2048"]
}


```

## Argument Reference

The following arguments are supported:

* `grok` - (Required) The Grok pattern to test.
* `log_lines` - (Required) The log lines to test the Grok pattern against.
* `account_id` - (Optional) The New Relic account ID to operate on.  This allows you to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `test_grok` - Nested attribute containing information about the test of Grok pattern against a list of log lines.
  * `matched` - Whether the Grok pattern matched.
  *  `log_line` - The log line that was tested against.
  * `attributes` - Nested list containing information about any attributes that were extracted.
      * `name` - The attribute name.
      * `value` - A string representation of the extracted value (which might not be a String).

```
Warning: This data source will use the account ID linked to your API key. At the moment it is not possible to dynamically set the account ID.
```
