---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_step_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-step-monitor"
description: |-
Create and manage a Synthetics Step monitor in New Relic.
---

# Resource: newrelic\_synthetics\_step\_monitor

Use this resource to create, update, and delete a Synthetics Step monitor in New Relic.

## Example Usage

```hcl
resource "newrelic_synthetics_step_monitor" "monitor" {
  name                                    = "step_monitor"
  enable_screenshot_on_failure_and_script = true
  locations_public                        = ["US_EAST_1", "US_EAST_2"]
  period                                  = "EVERY_6_HOURS"
  status                                  = "ENABLED"
  steps {
    ordinal = 0
    type    = "NAVIGATE"
    values  = ["https://www.newrelic.com"]
  }
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following are the common arguments supported for `STEP` monitor:

* `account_id`- (Optional) The account in which the Synthetics monitor will be created.
* `name` - (Required) The name for the monitor.
* `uri` - (Required) The uri the monitor runs against.
* `locations_public` - (Required) The location the monitor will run from. Valid public locations are https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/administration/synthetic-public-minion-ips/. You don't need the `AWS_` prefix as the provider uses NerdGraph. At least one of either `locations_public` or `location_private` is required.
* `location_private` - (Required) The location the monitor will run from. At least one of `locations_public` or `location_private` is required. See [Nested locations_private blocks](#nested-locations-private-blocks) below for details.
* `period` - (Required) The interval at which this monitor should run. Valid values are EVERY_MINUTE, EVERY_5_MINUTES, EVERY_10_MINUTES, EVERY_15_MINUTES, EVERY_30_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_12_HOURS, or EVERY_DAY.
* `status` - (Required) The run state of the monitor. (i.e. `ENABLED`, `DISABLED`, `MUTED`).
* `steps` - (Required) The steps that make up the script the monitor will run. See [Nested steps blocks](#nested-steps-blocks) below for details.
* `tag` - (Optional) The tags that will be associated with the monitor. See [Nested tag blocks](#nested-tag-blocks) below for details.

### Nested `location private` blocks

All nested `location_private` blocks support the following common arguments:

* `guid` - (Required) The unique identifier for the Synthetics private location in New Relic.
* `vse_password` - (Optional) The location's Verified Script Execution password, only necessary if Verified Script Execution is enabled for the location.

### Nested `steps` blocks

All nested `steps` blocks support the following common arguments:

* `ordinal` - (Required) The position of the step within the script ranging from 0-100.
* `type` - (Required) Name of the tag key. Valid values are ASSERT_ELEMENT, ASSERT_MODAL, ASSERT_TEXT, ASSERT_TITLE, CLICK_ELEMENT, DISMISS_MODAL, DOUBLE_CLICK_ELEMENT, HOVER_ELEMENT, NAVIGATE, SECURE_TEXT_ENTRY, SELECT_ELEMENT, TEXT_ENTRY.
* `values` - (Optional) The metadata values related to the step.

### Nested `tag` blocks

All nested `tag` blocks support the following common arguments:

* `key` - (Required) Name of the tag key.
* `values` - (Required) Values associated with the tag key.

## Additional Examples

### Create a monitor with a private location

The below example shows how you can define a private location and attach it to a monitor.

-> **NOTE:** It can take up to 10 minutes for a private location to become available.

```hcl
resource "newrelic_synthetics_private_location" "location" {
  description               = "Test Description"
  name                      = "private-location"
  verified_script_execution = true
}

resource "newrelic_synthetics_step_monitor" "bar" {
  name = "step_monitor"
  uri  = "https://www.one.example.com"
  location_private {
    guid         = newrelic_synthetics_private_location.location.id
    vse_password = "secret"
  }
  period = "EVERY_6_HOURS"
  status = "ENABLED"
  steps {
    ordinal = 0
    type    = "NAVIGATE"
    values  = ["https://google.com"]
  }
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```

### Creating a monitor with multiple steps

The following set of examples explain configuring a Step Monitor, comprising steps of various types, with details on `values` (a list of arguments needed by these steps), specific to each of them.

-> **NOTE:** Since the values required by each step type (which have been specified in the examples below) are to be added to the list `values`, and validating lists is currently [not supported](https://github.com/hashicorp/terraform-plugin-sdk/issues/156) by Terraform - if there is a mismatch in the values added corresponding to the step type, validation errors thrown by the NerdGraph API would be seen.

#### NAVIGATE

Use the `NAVIGATE` step to navigate to the specified URL. The values to be added are as follows (in the following order):
* Address of the website
* (Optional) User Agent String

The following example demonstrates the `steps` block used to create a step of the type `NAVIGATE`.

-> **INFO:** The `NAVIGATE` step is expected to be the first step in the configuration of the `newrelic_synthetics_step_monitor`, since other step types can be implemented only after a website is navigated to. Please ensure the step corresponding to the `ordinal` 0 is of the type `NAVIGATE` to avoid a 'BAD_REQUEST' error being thrown.

```hcl
  steps {
    ordinal = 0
    type    = "NAVIGATE"
    values  = ["https://www.example.com/", "user-agent-string"]
  }
```

#### ASSERT_ELEMENT

Use the `ASSERT_ELEMENT` step to assert on an element accessed by ID, CSS, or X-path. The values to be added are as follows (in the following order):
* Locator of the element (ID, CSS or XPath).
* `present` or `visible` (the HTML DOM property).
* `true` or `false`, corresponding to `present` or `visible` added before this. (For instance, `present` and `false` together translate to the condition "is not present", while `visible` and `true` together translate to "is visible").

The following example demonstrates the `steps` block used to create a step of the type `ASSERT_ELEMENT`.

```hcl
  steps {
    ordinal = 2
    type    = "ASSERT_ELEMENT"
    values  = ["/html/body/div/div/div/div[1]/div/div/nav/div/ul", "present", "false"]
  }
```


#### ASSERT_MODAL

Use the `ASSERT_MODAL` step to assert if the modal exists. The values to be added are as follows (in the following order):
* `true` or `false` (which translate to "is present" or "is not present" respectively)
* (Optional) Modal Text
* `accept` or `dismiss`
  The following example demonstrates the `steps` block used to create a step of the type `ASSERT_MODAL`.

```hcl
  steps {
    ordinal = 5
    type    = "ASSERT_MODAL"
    values  = ["true", "some-modal-text", "accept"]
  }

    // without the optional modal text
   steps {
    ordinal = 7
    type    = "ASSERT_MODAL"
    values  = ["false", "", "accept"]
  }
```

-> **WARNING:** The list of values is validated on the basis of the number of its contents too - hence, as shown in the second example above, if you have no optional modal text to add, please add an empty string "".

#### ASSERT_TEXT

Use the `ASSERT_TEXT` step to assert on text accessed by ID, CSS, or X-Path. The values to be added are as follows (in the following order):
* Locator of the element (ID, CSS or XPath).
* An operator to compare retrieved text with the input. Supported operators -
    * `==`  (Equals)
    * `!=`  (Not equal to)
    * `>=`  (Greater than or equal to)
    * `<=`  (Less than or equal to)
    * `>`   (Greater than)
    * `<`   (Less than)
    * `%=`  (Contains)
    * `!%=` (Does not contain)
* Input text, to compare the retrieved text from the element with.

The following example demonstrates the `steps` block used to create a step of the type `ASSERT_TEXT`.

```hcl
  steps {
      ordinal = 5
      type    = "ASSERT_TEXT"
      values  = ["//div[@class='boxed']//li[2]", "%=", "some-text-here"]
    }
```

#### ASSERT_TITLE

Use the `ASSERT_TITLE` step to assert on the title of a page. The values to be added are as follows (in the following order):
* An operator to compare the title on the page with the asserted title.
    * For a list of supported operators, please check the `ASSERT_TEXT` step type.
* Asserted text, to compare the retrieved title against.

The following example demonstrates the `steps` block used to create a step of the type `ASSERT_TITLE`.

```hcl
steps {
    ordinal = 3
    type    = "ASSERT_TITLE"
    values  = ["%=", "asserted-title"]
  }
```

#### CLICK, DOUBLE_CLICK, HOVER_ELEMENT

Use the `CLICK`, `DOUBLE_CLICK`, `HOVER_ELEMENT` step types to click, double click and hover on an element on a webpage respectively, identified by its locator (CSS, X-Path, etc). The only value to be added to use these resources is the locator (ID/CSS/X-Path) of the element.

The following example demonstrates the `steps` block used to create a step of the types `CLICK` and `DOUBLE_CLICK`.

```hcl
  steps {
    ordinal = 3
    type    = "HOVER_ELEMENT"
    values  = ["//div[@class='boxed']//li[1]"]
  }

  steps {
    ordinal = 4
    type    = "CLICK_ELEMENT"
    values  = ["//div[@class='boxed']//li[1]"]
  }

  steps {
    ordinal = 5
    type    = "DOUBLE_CLICK_ELEMENT"
    values  = ["//div[@class='boxed']//li[2]"]
  }
```

#### DISMISS_MODAL

Use the `DISMISS_MODAL` step to perform actions on a modal to dismiss it. The values to be added are as follows (in the following order):
* (Optional) Modal Text
* `accept` or `dismiss`

The following example demonstrates the `steps` block used to create a step of the type `DISMISS_MODAL`.

```hcl
  steps {
    ordinal = 4
    type    = "DISMISS_MODAL"
    values  = ["modal-text", "dismiss"]
  }
```

-> **WARNING:** The list of values is validated on the basis of the number of its contents too - hence, as previously stated under the `ASSERT_MODAL` section, if you have no optional modal text to add, please add an empty string "".


#### TEXT_ENTRY and SECURE_TEXT_ENTRY

Use the `TEXT_ENTRY` step to input text into an element accessed by ID, CSS, or X-Path. The values to be added are as follows (in the following order):
* Locator of the element (ID, CSS or XPath)
* Input text to be typed into the element

The following example demonstrates the `steps` block used to create a step of the type `TEXT_ENTRY`.

```hcl
  steps {
    ordinal = 2
    type    = "TEXT_ENTRY"
    values  = ["//div[@class='boxed']//li[1]", "text-to-be-typed"]
  }
```

Similar to `TEXT_ENTRY`, you may use the `SECURE_TEXT_ENTRY` step to input secure credentials into an element accessed by ID, CSS, or X-Path. The values to be added are as follows (in the following order):
* Locator of the element (ID, CSS or XPath)
* Label corresponding to the secure credentials created.
    * Visit [this page](https://one.newrelic.com/synthetics-nerdlets/secure-credential-list) to view the secure credentials linked to your New Relic account - or [this guide](https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/using-monitors/store-secure-credentials-scripted-browsers-api-tests/) on creating secure credentials if you haven't created already.

The following example demonstrates the `steps` block used to create a step of the type `SECURE_TEXT_ENTRY`.

```hcl
  steps {
    ordinal = 5
    type    = "SECURE_TEXT_ENTRY"
    values  = ["//div[@class='boxed']//li[1]", "TEST_CREDENTIALS"]
  }
```

#### SELECT_ELEMENT

Use the `SELECT_ELEMENT` step to select a dropdown element by value, text, ID, CSS, or X-Path. The values to be added are as follows (in the following order):
* Locator of the element (ID, CSS or XPath)
* Text or value on the basis of which the option to be selected may be identified

The following example demonstrates the `steps` block used to create a step of the type `SELECT_ELEMENT`.

```hcl
  steps {
    ordinal = 2
    type    = "SELECT_ELEMENT"
    values  = ["//div[@class='boxed']//li[1]", "text-to-identify-option"]
  }
```

## Attributes Reference

The following attributes are exported:

* `id` - The ID (GUID) of the synthetics step monitor.

## Import

Synthetics step monitor scripts can be imported using the `guid`, e.g.

```bash
$ terraform import newrelic_synthetics_step_monitor.monitor <guid>
```
