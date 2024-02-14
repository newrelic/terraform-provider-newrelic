---
layout: "newrelic"
page_title: "New Relic: newrelic_user"
sidebar_current: "docs-newrelic-datasource-user"
description: |-
  A data source that helps fetch authentication domains seen in the New Relic One UI, matching the name specified.
---

# Data Source: newrelic\_user

The `newrelic_user` data source may be used to search for a user by their name and/or email ID, and accordingly, fetch the ID of the matching user.

## Example Usage

The below example illustrates fetching a user's ID (and other arguments) using the ID of the authentication domain the user belongs to, as well as a name and/or email ID, which can be used as criteria to search for a user who matches these specified parameters.
```hcl
data "newrelic_authentication_domain" "foo" {
  name = "Test Authentication Domain"
}

data "newrelic_user" "user_one" {
  authentication_domain_id = data.newrelic_authentication_domain.foo.id
  name                     = "Test User"
}

data "newrelic_user" "user_two" {
  authentication_domain_id = data.newrelic_authentication_domain.foo.id
  email_id                 = "test_user@random.com"
}
```

## Argument Reference

The following arguments are supported:

* `authentication_domain_id` - (Required) The ID of the authentication domain the user to be searched for belongs to.
* `name` - (Optional) The name of the user to search for.
* `email_id` - (Optional) The email ID of the user to search for.

It should be noted that either `name` or `email_id` must be specified in order to retrieve a matching user.

-> **NOTE** If the specified `name` matches, or is contained in the names of multiple users in the account, the data source will return the first match from the list of all matching users retrieved from the API. However, when using the `email_id` argument as the search criterion, only the user with the specified email ID will be returned, as each user has a unique email ID and multiple users cannot have the same email ID.

-> **NOTE** The ID of an authentication domain can be retrieved using its name, via the data source `newrelic_authentication_domain`, as shown in the example above. Head over to the documentation of this data source for more details and examples.

## Attributes Reference

In addition to the attributes listed above, the following attribute is also exported by this resource:

* `id` - The ID of the matching user fetched.

