---
layout: "newrelic"
page_title: "New Relic Terraform Provider: Continuous Validation with Terraform Cloud"
sidebar_current: "docs-newrelic-provider-continuous-validation-guide"
description: |-
  Use this guide to set up the New Relic Terraform provider with Continuous
  Validation in Terraform Cloud.
---

# Continuous Validation with Terraform Cloud

The Continuous Validation feature in Terraform Cloud (TFC) allows users to make assertions about their infrastructure between applied runs. This helps users to identify issues at the time they first appear and avoid situations where a change is only identified during a future terraform plan/apply or once it causes a user-facing problem.

Checks can be added to Terraform configuration in Terraform Cloud (TFC) using check blocks. Check blocks contain assertions that are defined with a custom condition expression and an error message. When the condition expression evaluates to true during the check passes, but when the expression evaluates to false Terraform will show a warning message that includes the user-defined error message.

Custom conditions can be created using data from Terraform providers’ resources and data sources. Data can also be combined from multiple sources; for example, you can use checks to monitor expirable resources by comparing a resource’s expiration date attribute to the current time returned by Terraform’s built-in time functions.

Below, this guide shows examples of how data returned from the New Relic provider can be used to define checks in your Terraform configuration. For more information about continuous validation visit the [Workspace Health](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/health#continuous-validation) page in the Terraform Cloud documentation.

## Example - Check the configuration of a Notification Destination (`newrelic_notification_destination`)

[Notification destinations](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/destinations/#destinations-and-notifications-statuses) configured with the New Relic provider expose multiple attribute that indicate the configuration and status of the destination.

Continuous validation can be used to assert the configuration and status of a notification
destination and detect if there are any unexpected status changes that occur out-of-band.

The example below shows how a check block can be used to assert that a
notification destination is configured and working as expected.

```hcl
data "newrelic_notification_destination" "example" {
  id = "1e543419-0c25-456a-9057-fb0eb310e60b"
}

check "check_notification_destination_status" {
  assert {
    condition = data.newrelic_notification_destination.example.type == "SLACK"
    error_message = format("Notification destination (%s) type should be set to 'SLACK', instead type is '%s'",
      data.newrelic_notification_destination.example.name,
      data.newrelic_notification_destination.example.type
    )
  }
  assert {
    condition = data.newrelic_notification_destination.example.status == "DEFAULT"
    error_message = format("Notification destination (%s) should be in a 'DEFAULT' status, instead status is '%s'",
      data.newrelic_notification_destination.example.name,
      data.newrelic_notification_destination.example.status
    )
  }
}
```

## Example - Check the last updated time for Synthetic Secure Credential (`newrelic_synthetics_secure_credential`)

[Secure credentials](https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/using-monitors/store-secure-credentials-scripted-browsers-api-tests/) configured with the New Relic provider expose a `last_updated` attribute that indicates when the credential was last updated.

You can use this attribute to validate the age of the secure credential and make sure it's either set to latest, or not updated without permission.

```hcl
data "newrelic_synthetics_secure_credential" "example" {
  key = "MY_KEY"
}

check "check_synthetics_secure_credential" {
  assert {
    condition = timecmp(timestamp(), timeadd(data.newrelic_synthetics_secure_credential.example.last_updated, "-720h")) < 0
    error_message = format("Synthetics secure credential (%s) is outdated",
      data.newrelic_synthetics_secure_credential.example.key
    )
  }
}
```
