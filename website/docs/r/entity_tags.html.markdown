---
layout: "newrelic"
page_title: "New Relic: newrelic_entity_tags"
sidebar_current: "docs-newrelic-resource-entity-tags"
description: |-
  Create and manage tags for a New Relic One entity.
---

# Resource: newrelic\_entity\_tags

Use this resource to create, update, and delete tags for a New Relic One entity.

-> **IMPORTANT!** Version 2.0.0 of the New Relic Terraform Provider introduces some [additional requirements](/providers/newrelic/newrelic/latest/docs/guides/migration_guide_v2) for configuring the provider.
<br><br>
Before upgrading to version 2.0.0 or later, it is recommended to upgrade to the most recent 1.x version of the provider and ensure that your environment successfully runs `terraform plan` without unexpected changes.

## Example Usage

```hcl
data "newrelic_entity" "foo" {
  name   = "Example application"
  type   = "APPLICATION"
  domain = "APM"
}

resource "newrelic_entity_tags" "foo" {
  guid = data.newrelic_entity.foo.guid

  tag {
    key    = "my-key"
    values = ["my-value", "my-other-value"]
  }

  tag {
    key    = "my-key-2"
    values = ["my-value-2"]
  }
}
```

#### Example of applying multiple tags to multiple entities using a nested [`dynamic`](https://developer.hashicorp.com/terraform/language/expressions/dynamic-blocks) block
```hcl
locals {
  apps = toset([
    "Example App Name 1",
    "Example App Name 2",
  ])

  custom_tags = {
    "tag-key-1" = "tag-value-1"
    "tag-key-2" = "tag-value-2"
    "tag-key-3" = "tag-value-3"
  }
}

data "newrelic_entity" "foo" {
  for_each = local.apps
  name     = each.key # Note: each.key and each.value are the same for a set
  type     = "APPLICATION"
  domain   = "APM"
}

resource "newrelic_entity_tags" "foo" {
  for_each = local.apps
  guid     = data.newrelic_entity.foo[each.key].guid

  dynamic "tag" {
    for_each = local.custom_tags
    content {
      key    = tag.key
      values = [tag.value]
    }
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
