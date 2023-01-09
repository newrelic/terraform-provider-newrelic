---
layout: "newrelic"
page_title: "New Relic: newrelic_log_parsing_rule"
sidebar_current: "docs-newrelic-resource-log-parsing-rule"
description: |-
Create and manage Log Parsing Rule.
---

# Resource: newrelic\_log\_parsing\_rule

Use this resource to create, update and delete New Relic Log Parsing Rule.

## Example Usage

Use this example to create the log parse rule.
```hcl

resource "newrelic_log_parsing_rule" "foo"{
    account_id  = 12345
    name        = "log_parse_rule"
    attribute   = "message"
    enabled     = true
    grok        = "sampleattribute='%%{NUMBER:test:int}'"
    lucene      = "logtype:linux_messages"
    nrql        = "SELECT * FROM Log WHERE logtype = 'linux_messages'"
}

```

## Additional Example
Use this example to validate a grok pattern and create the log parse rule.  More
information on grok pattern can be found [here](https://docs.newrelic.com/docs/logs/ui-data/parsing/#grok)
```hcl
data "newrelic_test_grok_pattern" "grok"{
    account_id  = 12345
    grok        = "%%{IP:host_ip}"
    log_lines   = ["host_ip: 43.3.120.2"]
}
resource "newrelic_log_parsing_rule" "foo"{
    account_id  = 12345
    name        = "log_parse_rule"
    attribute   = "message"
    enabled     = true
    grok        = data.newrelic_test_grok_pattern.grok.grok
    lucene      = "logtype:linux_messages"
    nrql        = "SELECT * FROM Log WHERE logtype = 'linux_messages'"
    matched     = data.newrelic_test_grok_pattern.grok.test_grok[0].matched
}

```


## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of rule.
* `grok` - (Required) The Grok of what to parse.
* `lucene` - (Required) The Lucene to match events to the parsing rule.
* `enabled` - (Required) Whether the rule should be applied or not to incoming data.
* `nrql` - (Required) The NRQL to match events to the parsing rule.
* `account_id` - (Optional) The account id associated with the obfuscation rule.
* `matched` - (Optional) Whether the Grok pattern matched.




## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the log parsing rule.

## Import

New Relic log parsing rule can be imported using the rule ID, e.g.

```bash
$ terraform import newrelic_log_parsing_rule.foo 3456789
```