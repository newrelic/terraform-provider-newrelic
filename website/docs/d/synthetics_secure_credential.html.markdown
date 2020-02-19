---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_secure_credential"
sidebar_current: "docs-newrelic-datasource-synthetics-secure-credential"
description: |-
  Grabs a Synthetics secure credential by its key.
---

# Data Source: newrelic\_synthetics\_secure\_credential

Use this data source to get information about a specific Synthetics secure credential in New Relic that already exists.

Note that the secure credential's value is not returned as an attribute for security reasons.

## Example Usage

```hcl
data "newrelic_synthetics_secure_credential" "foo" {
  key = "MY_KEY"
}
```

## Argument Reference

The following arguments are supported:

  * `key` - (Required) The secure credential's key name.  Regardless of the case used in the configuration, the provider will provide an upcased key to the underlying API.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `description` - The secure credential's description.
  * `created_at` - The time the secure credential was created.
  * `updated_at` - The time the secure credential was last updated.
