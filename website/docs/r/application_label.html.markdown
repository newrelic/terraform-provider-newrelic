---
layout: "newrelic"
page_title: "New Relic: newrelic_application_label"
sidebar_current: "docs-newrelic-resource-application-label"
description: |-
  Create and manage an Application label in New Relic.
---

# Resource: newrelic\_application\_label

#### DEPRECATED! Use at your own risk. This feature may be removed in the next major release.
Use this resource to create, update, and delete an Application label in New Relic.

## Example Usage

```hcl
data "newrelic_application" "app1" {
  name="myapp1"
}

data "newrelic_application" "app2" {
  name="myapp2"
}

resource "newrelic_application_label" "foo" {
  category = "Team"
  name = "MyTeam"
  links {
    applications = [
      data.newrelic_application.app1.id,
      data.newrelic_application.app2.id
    ]
    servers = []
  }
}
```

## Argument Reference

The following arguments are supported:

  * `category` - (Required) A string representing the label key/category.
  * `name` - (Required) A string that will be assigned to the label.
  * `links` - (Required) The resources to which label should be assigned to. At least one of the following attributes must be set.
    * `applications` - An array of application IDs.
    * `servers` - An array of server IDs.
