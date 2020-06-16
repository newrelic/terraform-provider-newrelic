---
layout: "newrelic"
page_title: "New Relic: newrelic_entity_tags"
sidebar_current: "docs-newrelic-resource-entity-tags"
description: |-
  Create and manage tags for a New Relic One entity.
---

# Resource: newrelic\_entity\_tags

Use this resource to create, update, and delete tags for a New Relic One entity.

A New Relic Personal API key is required to provision this resource.  Set the `api_key`
attribute in the `provider` block or the `NEW_RELIC_API_KEY` environment
variable with your Personal API key.

## Example Usage

```hcl
data "newrelic_entity" "foo" {
    name = "Example application"
    type = "APPLICATION"
    domain = "APM"
}

resource "newrelic_entity_tags" "foo" {
	guid = data.newrelic_entity.foo.guid

	tag {
        key = "my-key"
        values = ["my-value", "my-other-value"]
    }

	tag {
        key = "my-key-2"
        values = ["my-value-2"]
    }
}
```

## Argument Reference

The following arguments are supported:

  * `guid` - (Required) The guid of the entity to tag.
  * `tag` - (Optional) A nested block that describes an entity tag. See [Nested tag blocks](#nested-`tag`-blocks) below for details.

### Nested `tag` blocks

All nested `tag` blocks support the following common arguments:

  * `key` - (Required) The tag key.
  * `values` - (Required) The tag values.

## Import

New Relic One entity tags can be imported using a concatenated string of the format
 `<guid>`, e.g.

```bash
$ terraform import newrelic_entity_tags.foo MjUyMDUyOHxBUE18QVBRTElDQVRJT058MjE1MDM3Nzk1
```
