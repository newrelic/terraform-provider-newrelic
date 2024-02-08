---
layout: "newrelic"
page_title: "New Relic: newrelic_authentication_domain"
sidebar_current: "docs-newrelic-datasource-authentication-domain"
description: |-
  A data source that helps fetch authentication domains seen in the New Relic One UI, matching the name specified.
---

# Data Source: newrelic\_authentication\_domain

Use this data source to fetch the ID of an authentication domain belonging to your account, matching the specified name.

## Example Usage

```hcl
data "newrelic_authentication_domain" "foo" {
  name = "Test Authentication Domain"
}

output "foo" {
  value = data.newrelic_authentication_domain.foo.id
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the authentication domain to be searched for. An error is thrown, if no authentication domain is found with the specified name.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the matching authentication domain fetched.

