---
layout: "newrelic"
page_title: "New Relic: newrelic_current_user"
sidebar_current: "docs-newrelic-datasource-current-user"
description: |-
  This data source helps fetch the current user, i.e. the owning user of the credentials (API Key), using which the New Relic Terraform Provider has been initialized to perform operations.
---

# Data Source: newrelic_current_user

The `newrelic_current_user` data source helps fetch the current user, i.e. the owning user of the credentials (API Key), using which the New Relic Terraform Provider has been initialized to perform operations.

-> **NOTE:** If you would like to search for a specific user by `name` or `email_id` within a specific authentication domain, please head over to the documentation of the [`newrelic_user`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/data-sources/user) data source for more details and examples to do so.

## Example Usage

The below example illustrates fetching the current user and associated metadata (the `name` and `email_id` of the current user) using the `newrelic_current_user` data source. The data source does not require any arguments to be specified.
```hcl
data "newrelic_current_user" "foo" {
}

output "current_user_details" {
  value = {
    "user_id" : data.newrelic_current_user.foo.id
    "user_email_id" : data.newrelic_current_user.foo.email_id
    "user_name" : data.newrelic_current_user.foo.name
  }
}
```

The ID of the current user fetched may be applied in other use cases within the New Relic Terraform Provider; for instance, this may be furnished as the value of the argument [`user_id`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/api_access_key#user_id) if required, to create `USER` API Keys with the same user's credentials using the [`newrelic_api_access_key`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/api_access_key) resource.

## Argument Reference
This data source does not require any arguments to be specified.

## Attributes Reference
The following attributes are exported by this data source.
* `id` - The ID of the current user.
* `name` - The name of the current user.
* `email_id` - The email ID of the current user.

