---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_secure_credential"
sidebar_current: "docs-newrelic-resource-synthetics-secure-credential"
description: |-
  Create and manage Synthetics secure credentials in New Relic.
---

# Resource: newrelic\_synthetics\_secure\_credential

Use this resource to create and manage New Relic Synthetic secure credentials.

## Example Usage

```hcl
resource "newrelic_synthetics_secure_credential" "foo" {
  key = "MY_KEY"
  value = "My value"
  description = "My description"
}
```

## Argument Reference

The following arguments are supported:

  * `key` - (Required) The secure credential's key name.  Regardless of the case used in the configuration, the provider will provide an upcased key to the underlying API.
  * `value` - (Required) The secure credential's value. 
  * `description` - (Optional) The secure credential's description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `created_at` - The time the secure credential was created.
  * `updated_at` - The time the secure credential was last updated.

## Import

A Synthetics secure credential can be imported using its `key`:

```
$ terraform import newrelic_synthetics_secure_credential.foo MY_KEY
```
