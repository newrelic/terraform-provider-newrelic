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

```hcl

resource "newrelic_log_parsing_rule" "foo"{
	account_id = %[1]d
	name = "%[2]s"
	attribute = "%[3]s"
	enabled     = true
    grok        = "sampleattribute='%%{NUMBER:test:int}'"
    lucene      = "logtype:linux_messages"
    nrql        = "SELECT * FROM Log WHERE logtype = 'linux_messages'"
}

```


## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The account id associated with the obfuscation rule.
* `name` - (Required) Name of rule.
* `grok` - (Required) The Grok of what to parse.
* `lucene` - (Required) The Lucene to match events to the parsing rule.
* `enabled` - (Required) Whether the rule should be applied or not to incoming data.
* `matched` - (Optional) Whether the Grok pattern matched.
* `nrql` - (Required) The NRQL to match events to the parsing rule.




## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the log parsing rule.

## Import

New Relic log parsing rule can be imported using the rule ID, e.g.

```bash
$ terraform import newrelic_log_parsing_rule.foo 3456789
```