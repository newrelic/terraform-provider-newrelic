---
layout: "newrelic"
page_title: "New Relic: newrelic_notification_destination"
sidebar_current: "docs-newrelic-resource-notification-destination"
description: |-
Create and manage a notification destination for notifications in New Relic.
---

# Resource: newrelic\_notification\_destination

Use this resource to create and manage New Relic notification destinations. Details regarding supported products and permissions can be found [here](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/destinations).

## Example Usage

##### [Webhook](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/notification-integrations/#webhook)
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
    user = "username"
    password = "password"
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) Determines the New Relic account where the notification destination will be created. Defaults to the account associated with the API key used.
* `name` - (Required) The name of the destination.
* `type` - (Required) The type of destination.  One of: `EMAIL`, `SERVICE_NOW`, `WEBHOOK`, `JIRA`, `PAGERDUTY_ACCOUNT_INTEGRATION` or `PAGERDUTY_SERVICE_INTEGRATION`.
* `auth` - A nested block that describes a notification destination authentication. Only one auth block is permitted per notification destination definition.  See [Nested auth blocks](#nested-auth-blocks) below for details.
* `properties` - A nested block that describes a notification destination properties. See [Nested properties blocks](#nested-properties-blocks) below for details.

### Nested `auth` blocks

* `type` - (Required) The type of the auth.  One of: `TOKEN` or `BASIC`.

Each authentication type supports a specific set of arguments:

* `BASIC`
  * `user` - (Required) The username of the basic auth.
  * `password` - (Required) Specifies an authentication password for use with a destination.
* `TOKEN`
  * `prefix` - (Required) The prefix of the token auth.
  * `token` - (Required) Specifies the token for integrating.

~> **NOTE:** OAuth2 authentication type is not available via terraform for notifications destinations.

### Nested `properties` blocks

Each notification destination type supports a specific set of arguments for the `properties` block:

* `EMAIL`
  * `email` - (Required) A map of key/value pairs that represents the email addresses.
* `WEBHOOK`
  * `url` - (Required) A map of key/value pairs that represents the webhook url.
* `SERVICE_NOW`
  * `url` - (Required) A map of key/value pairs that represents the service now destination url.
  * `two_way_integration` - (Optional) A map of key/value pairs that represents the two-way integration on/off flag.
* `JIRA`
  * `url` - (Required) A map of key/value pairs that represents the jira url.
  * `two_way_integration` - (Optional) A map of key/value pairs that represents the two-way integration on/off flag.
* `PAGERDUTY_ACCOUNT_INTEGRATION`
  * `two_way_integration` - (Optional) A map of key/value pairs that represents the two-way integration on/off flag.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the destination.

## Additional Examples

##### [ServiceNow](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/notification-integrations/#servicenow)

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
    user = "username"
    password = "password"
  }
}
```

##### [Email](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/notification-integrations/#email)
```hcl
resource "newrelic_notification_destination" "foo" {
  name = "email-example"
  type = "EMAIL"

  properties {
    key = "email"
    value = "email@email.com,email2@email.com"
  }
}
```

##### [Jira](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/notification-integrations/#jira)
```hcl
resource "newrelic_notification_destination" "foo" {
  name = "jira-example"
  type = "JIRA"

  properties {
    key = "url"
    value = "https://example.atlassian.net"
  }
  
  auth = {
    type = "BASIC"
    user = "example@email.com"
    password = "password"
  }
}
```

##### [PagerDuty with service integration](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/notification-integrations/#pagerduty-sli)
```hcl
resource "newrelic_notification_destination" "foo" {
  name = "pagerduty-service-example"
  type = "PAGERDUTY_SERVICE_INTEGRATION"

  auth = {
    type   = "TOKEN"
    prefix = "Token token="
    token  = "10567a689d984d03c021034b22a789e2"
  }
}
```

##### [PagerDuty with account integration](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/notification-integrations/#pagerduty-ali)
```hcl
resource "newrelic_notification_destination" "foo" {
  name = "pagerduty-account-example"
  type = "PAGERDUTY_ACCOUNT_INTEGRATION"

  properties {
    key = "two_way_integration"
    value = "true"
  }

  auth = {
    type   = "TOKEN"
    prefix = "Token token="
    token  = "u+E8EU3MhsZwLfZ1ic1A"
  }
}
``` 


~> **NOTE:** Sensitive data such as destination API keys, service keys, auth object, etc are not returned from the underlying API for security reasons and may not be set in state when importing.

## Additional Information
More information about destinations integrations can be found in NewRelic [documentation](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/notification-integrations/).
More details about the destinations API can be found [here](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-api-notifications-destinations).
