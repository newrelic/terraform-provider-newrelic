---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_private_location"
sidebar_current: "docs-newrelic-resource-synthetics-private-location"
description: |-
Create and manage Synthetics private location in New Relic.
---

# Resource: newrelic\_synthetics\_private\_location

Use this resource to create and manage New Relic Synthetic private location.

## Example Usage

```hcl
resource "newrelic_synthetics_private_location" "location" {
  description = "The private location description"
  name        = "The name of the private location"
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The account in which the private location will be created.
* `description` - (Required) The private location description.
* `name` - (Required) The name of the private location.
* `verified_script_execution` - (Optional) The private location requires a password to edit if value is true. Defaults to `false`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the private location. Same as `guid`.
* `domain_id` - The private location globally unique identifier.
* `guid` - The unique client identifier for the private location in New Relic. Same as `id`.
* `location_id` - An alternate identifier based on name.
* `key` - The private locations key.

## Import

A Synthetics private location can be imported using the `guid`

```
$ terraform import newrelic_synthetics_private_location.location GUID
```