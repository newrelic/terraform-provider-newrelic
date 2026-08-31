---
layout: "newrelic"
page_title: "New Relic: newrelic_installation_recipe_event"
sidebar_current: "docs-newrelic-resource-installation-recipe-event"
description: |-
  Create and manage an installation recipe event in New Relic.
---

# Resource: newrelic\_installation\_recipe\_event

Use this resource to create and manage New Relic installation recipe events. These events are created on behalf of the newrelic-cli whenever the CLI attempts to install a recipe (e.g., the infrastructure-agent). Details regarding guided install can be found [here](https://docs.newrelic.com/docs/full-stack-observability/observe-everything/get-started/new-relic-guided-install-overview/).

## Example Usage

```hcl
resource "newrelic_installation_recipe_event" "foo" {
  account_id       = 12345678
  cli_version      = "1.0.0"
  complete         = true
  display_name     = "Infrastructure Agent"
  entity_guid      = "ABC123"
  host_name        = "my-host"
  kernel_arch      = "x86_64"
  kernel_version   = "5.4.0"
  log_file_path    = "/var/log/newrelic-cli.log"
  name             = "infrastructure-agent-linux-installer"
  os               = "linux"
  platform         = "debian"
  platform_family  = "debian"
  platform_version = "10"
  status           = "INSTALLED"
  targeted_install = false
  timestamp        = 1609459200
  validation_duration_milliseconds = 5000

  error {
    details           = ""
    message           = ""
    optimized_message = ""
  }
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) Determines the New Relic account where the installation recipe event will be created. Defaults to the account associated with the API key used.
* `cli_version` - (Required) The version of the newrelic-cli that was used for a given recipe.
* `complete` - (Required) Whether or not the recipe has been installed and all steps have been completed.
* `display_name` - (Required) The display name for a given recipe.
* `entity_guid` - (Required) The entity GUID for a given recipe.
* `error` - (Optional) A nested block that describes the error returned for a given recipe. Only one `error` block is permitted per resource definition. See [Nested error blocks](#nested-error-blocks) below for details.
* `host_name` - (Required) The host name of the customer's machine.
* `install_id` - (Optional) The unique ID that corresponds to an install event.
* `install_library_version` - (Optional) The version of the open-install-library that is being used.
* `kernel_arch` - (Required) The kernel architecture of the customer's machine.
* `kernel_version` - (Required) The kernel version of the customer's machine.
* `log_file_path` - (Required) The path to the log file on the customer's host.
* `name` - (Required) The unique name for a given recipe.
* `os` - (Required) The OS of the customer's machine.
* `platform` - (Required) The platform name provided by the open-install-library.
* `platform_family` - (Required) The platform family name provided by the open-install-library.
* `platform_version` - (Required) The platform version provided by the open-install-library.
* `redirect_url` - (Optional) The redirect URL created by the CLI used for redirecting to a particular entity.
* `status` - (Required) The status for a given recipe. One of: `AVAILABLE`, `CANCELED`, `DETECTED`, `FAILED`, `INSTALLED`, `INSTALLING`, `RECOMMENDED`, `SKIPPED`, `UNSUPPORTED`.
* `targeted_install` - (Required) Whether or not the recipe being installed is a targeted install.
* `task_path` - (Optional) The path to the installation task as defined in the open-install-library.
* `timestamp` - (Required) The timestamp for when the recipe event occurred, in epoch seconds.
* `validation_duration_milliseconds` - (Required) The number of milliseconds it took to validate the recipe.

### Nested `error` blocks

* `details` - (Optional) Error details, if any.
* `message` - (Optional) The actual error message.
* `optimized_message` - (Optional) An optimised message for the error.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The install ID of the recipe event.

## Import

New Relic installation recipe events can be imported using the install ID:

```
terraform import newrelic_installation_recipe_event.foo <install_id>
```

~> **NOTE:** Installation recipe events are write-only operations in the New Relic API. The Read operation is a no-op for this resource, meaning that after import, the state will only reflect what was last written via Terraform. Sensitive data is not returned from the underlying API.

## Additional Information

More information about the New Relic guided install and the newrelic-cli can be found in the [New Relic documentation](https://docs.newrelic.com/docs/full-stack-observability/observe-everything/get-started/new-relic-guided-install-overview/) and the [newrelic-cli GitHub repository](https://github.com/newrelic/newrelic-cli).