---
layout: "newrelic"
page_title: "New Relic: newrelic_browser_application"
sidebar_current: "docs-newrelic-browser-application"
description: |-
    Create, update, and delete a standalone New Relic browser application.
---

# Resource: newrelic\_browser\_application

Use this resource to create, update, and delete a standalone New Relic browser application.

## Example Usage

The following Terraform configuration is an example that illustrates the basic use case of creating a standalone browser application.
```hcl
resource "newrelic_browser_application" "foo" {
  name                        = "example-browser-app"
  cookies_enabled             = true
  distributed_tracing_enabled = true
  loader_type                 = "SPA"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the browser application.
* `cookies_enabled` - (Optional) Configures cookies. Defaults to `true`, if not specified.
* `distributed_tracing_enabled` - (Optional) Configures distributed tracing in browser apps. Defaults to `true`, if not specified.
* `loader_type` - (Optional) Determines the browser loader configured. Valid values are `SPA`, `PRO`, and `LITE`. The default is `SPA`. Refer to the [browser agent loader documentation](https://docs.newrelic.com/docs/browser/browser-monitoring/installation/install-browser-monitoring-agent/#agent-types) for more information on valid loader types.
* `account_id` - (Optional) The account ID of the New Relic account you wish to create the browser application in. Defaults to the value of the environment variable `NEW_RELIC_ACCOUNT_ID` if not specified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the browser application.
* `application_id` - The application ID of the browser application (not to be confused with GUID).
* `js_config` - The JavaScript configuration of the browser application, encoded into a string.

## Import

A browser application can be imported using its GUID, i.e.

```bash
$ terraform import newrelic_browser_application.foo <GUID>
```
