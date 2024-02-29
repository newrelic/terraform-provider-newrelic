---
layout: "newrelic"
page_title: "New Relic: newrelic_group"
sidebar_current: "docs-newrelic-datasource-group"
description: |-
  A data source that helps fetch group(s) seen in the New Relic One UI, matching the name specified.
---

# Data Source: newrelic\_group

The `newrelic_group` data source helps search for a group by its name and retrieve the ID of the matching group and other associated attributes.

## Example Usage

The below example illustrates fetching the ID of a group (and IDs of users who belong to the group, if any) using the required arguments.
    
```hcl
data "newrelic_authentication_domain" "foo" {
  name = "Test Authentication Domain"
}

data "newrelic_group" "foo" {
  authentication_domain_id = data.newrelic_authentication_domain.foo.id
  name                     = "Test Group"
}
```

## Argument Reference

The following arguments are supported:

* `authentication_domain_id` - (Required) The ID of the authentication domain the group to be searched for belongs to.
* `name` - (Required) The name of the group to search for.

-> **NOTE** The ID of an authentication domain can be retrieved using its name, via the data source `newrelic_authentication_domain`, as shown in the example above. Head over to the documentation of this data source for more details and examples.

## Attributes Reference

In addition to the attributes listed above, the following attributes are also exported by this data source:

* `id` - The ID of the fetched matching group.
* `user_ids` - IDs of users who belong to the group. In the absence of any users in the group, the value of this attribute would be an empty list.

## Additional Examples

The following example demonstrates utilizing attributes exported by this data source.

In order to directly reference the attributes `id` and `user_ids` from this data source, you can use the syntax `data.newrelic_group.foo.id` and `data.newrelic_group.foo.user_ids`, respectively. However, if you need to assign these values to local variables and perform further processing (such as conditionally formatting the `user_ids` attribute as shown in the example below), consider using the provided configuration. These variables can then be accessed elsewhere using the syntax `local.id` and `local.user_id`, respectively.

```hcl
locals {
  id       = data.newrelic_group.foo.id
  user_ids = length(data.newrelic_group.foo.user_ids) > 0 ? join(", ", data.newrelic_group.foo.user_ids) : ""
}

data "newrelic_authentication_domain" "foo" {
  name = "Test Authentication Domain"
}

data "newrelic_group" "foo" {
  authentication_domain_id = data.newrelic_authentication_domain.foo.id
  name                     = "Test Group"
}
```



