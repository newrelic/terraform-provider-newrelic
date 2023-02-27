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

Basic usage to create a standalone browser application.
```hcl
resource "newrelic_browser_application" "foo" {
  name = "example-browser-app"
  cookies_enabled = true
  distributed_tracing_enabled = true
  loader_type = "SPA"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the browser application.
* `cookies_enabled` - (Optional) Configure cookies. The default is enabled: true.
* `distributed_tracing_enabled` - (Optional) Configure distributed tracing in browser apps. The default is enabled: true.
* `loader_type` - (Optional) Determines which browser loader is configured. Valid values are `SPA`, `PRO`, and `LITE`. The default is `SPA`. See the [browser agent loader documentation](https://docs.newrelic.com/docs/browser/browser-monitoring/installation/install-browser-monitoring-agent/#agent-types) for a for information on the valid loader types.
* `account_id` - (Optional) The New Relic account ID of the account you wish to create the browser application. Defaults to the account ID set in your environment variable `NEW_RELIC_ACCOUNT_ID`.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the browser application.

## Import

Browser applications can be imported using the GUID of the browser application.

```bash
$ terraform import newrelic_browser_application.foo <GUID>
```
