---
layout: "newrelic"
page_title: "New Relic: newrelic_notification_destination"
sidebar_current: "docs-newrelic-resource-notification-destination"
description: |-
Create and manage a notification destination for notifications in New Relic.
---

# Resource: newrelic\_notification\_destination

Use this resource to create and manage New Relic notification destinations.

## Example Usage

##### Webhook
```hcl
resource "newrelic_notification_destination" "foo" {
  name = "foo"
  type = "WEBHOOK"

  properties {
    key = "url"
    value = "https://webhook.site/"
  }

  auth = {
    type = "BASIC"
    user = "user"
    password = "1234"
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) Determines the New Relic account where the notification destination will be created. Defaults to the account associated with the API key used.
* `name` - (Required) The name of the destination.
* `type` - (Required) The type of destination.  One of: `EMAIL`, `SERVICE_NOW`, `WEBHOOK`, `PAGERDUTY_ACCOUNT_INTEGRATION` or `PAGERDUTY_SERVICE_INTEGRATION`.
* `auth` - (Required) A nested block that describes a notification destination authentication. Only one auth block is permitted per notification destination definition.  See [Nested auth blocks](#nested-auth-blocks) below for details.
* `properties` - (Required) A nested block that describes a notification destination properties.  Only one properties block is permitted per notification destination definition.  See [Nested properties blocks](#nested-properties-blocks) below for details.

### Nested `auth` blocks

* `type` - (Required) The type of the auth.  One of: `TOKEN` or `BASIC`.

Each authentication type supports a specific set of arguments:

* `basic`
  * `user` - (Required) The username of the basic auth.
  * `password` - (Optional) Specifies an authentication password for use with a destination.
* `token`
  * `prefix` - (Required) The prefix of the token auth.
  * `token` - (Required) Specifies the token for integrating.

### Nested `properties` blocks

Each notification destination type supports a specific set of arguments for the `properties` block:

* `EMAIL`
  * `email` - (Required) A map of key/value pairs that represents the email addresses.
* `WEBHOOK`
  * `url` - (Required) A map of key/value pairs that represents the webhook url.
* `SERVICE_NOW`
  * `url` - (Required) A map of key/value pairs that represents the service now url.
  * `two_way_integration` - (Optional) A map of key/value pairs that represents the two-way integration on/off flag.
* `PAGERDUTY_SERVICE_INTEGRATION`
  * `two_way_integration` - (Optional) A map of key/value pairs that represents the two-way integration on/off flag.
* `PAGERDUTY_ACCOUNT_INTEGRATION`
  * `two_way_integration` - (Optional) A map of key/value pairs that represents the two-way integration on/off flag.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the destination.

## Additional Examples

##### ServiceNow
# oauth2 not possible

```hcl
resource "newrelic_notification_destination" "foo" {
  name = "servicenow-example"
  type = "SERVICE_NOW"

  properties {
    key = "url"
    value = "https://service-now.com/"
  }

  properties {
    key = "two_way_integration"
    value = "true"
  }

  auth = {
    type = "BASIC"
    user = "user"
    password = "pass"
  }
}
```

##### Email
```hcl
resource "newrelic_notification_destination" "foo" {
  name = "email-example"
  type = "EMAIL"
  
  properties {
    key = "email"
    value = "email@email.com,email2@email.com"
  }
  
  auth = {
    type = "TOKEN"
    prefix = "prefix"
    token = "bearer"
  }
}
```

##### PagerDuty with service integration
```hcl
resource "newrelic_notification_destination" "foo" {
  name = "pagerduty-service-example"
  type = "PAGERDUTY_SERVICE_INTEGRATION"

  properties {
    key = "two_way_integration"
    value = "true"
  }

  auth = {
    type   = "TOKEN"
    prefix = "prefix"
    token  = "bearer"
  }
}
```

##### PagerDuty with account integration
```hcl
resource "newrelic_notification_destination" "foo" {
  name = "pagerduty-account-example"
  type = "PAGERDUTY_ACCOUNT_INTEGRATION"

  auth = {
    type   = "TOKEN"
    prefix = "prefix"
    token  = "bearer"
  }
}
``` 

~> **NOTE:** Sensitive data such as destination API keys, service keys, etc are not returned from the underlying API for security reasons and may not be set in state when importing.