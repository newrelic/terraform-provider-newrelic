---
layout: "newrelic"
page_title: "New Relic: newrelic_application"
sidebar_current: "docs-newrelic-resource-application"
description: |-
  Manage configuration for an existing application in New Relic.
---

# Resource: newrelic_application

Use this resource to manage configuration for an application that already exists in New Relic.

## Example Usage

```hcl
resource "newrelic_application" "app" {
  name = "my-app"
  app_apdex_threshold = "0.7"
  end_user_apdex_threshold = "0.8"
  enable_real_user_monitoring = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the application in New Relic APM.
* `app_apdex_threshold` - (Required) The appex threshold for the New Relic application.
* `end_user_apdex_threshold` - (Required) The user's apdex threshold for the New Relic application.
* `enable_real_user_monitoring` - (Required) Enable or disable real user monitoring for the New Relic application.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the application.

## Import

Applications can be imported using notation `application_id`, e.g.

```
$ terraform import newrelic_application.main 6789012345
```

## Notes

Applications that have reported data in the last twelve hours cannot be deleted.
