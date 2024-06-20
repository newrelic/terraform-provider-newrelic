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
  account_id = 12345678
  name = "foo"
  type = "WEBHOOK"

  secure_url {
    prefix = "https://webhook.mywebhook.com/"
    secure_suffix = "service_id/123456"
  }

  property {
    key = "source"
    value = "terraform"
  }

  auth_custom_header {
    key = "API_KEY"
    value = "test-api-key"
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) Determines the New Relic account where the notification destination will be created. Defaults to the account associated with the API key used.
* `name` - (Required) The name of the destination.
* `type` - (Required) The type of destination.  One of: `EMAIL`, `SERVICE_NOW`, `SERVICE_NOW_APP`, `WEBHOOK`, `JIRA`, `MOBILE_PUSH`, `EVENT_BRIDGE`, `PAGERDUTY_ACCOUNT_INTEGRATION` or `PAGERDUTY_SERVICE_INTEGRATION`. The types `SLACK` and `SLACK_COLLABORATION` can only be imported, updated and destroyed (cannot be created via terraform).
* `auth_basic` - (Optional) A nested block that describes a basic username and password authentication credentials. Only one auth_basic block is permitted per notification destination definition.  See [Nested auth_basic blocks](#nested-auth_basic-blocks) below for details.
* `auth_token` - (Optional) A nested block that describes a token authentication credentials. Only one auth_token block is permitted per notification destination definition.  See [Nested auth_token blocks](#nested-auth_token-blocks) below for details.
* `auth_custom_header` - (Optional) A nested block that describes a custom header authentication credentials. Multiple blocks are permitted per notification destination definition. [Nested auth_custom_header blocks](#nested-authcustomheader-blocks) below for details.
* `secure_url` - (Optional) A nested block that describes a URL that contains sensitive data at the path or parameters. Only one secure_url block is permitted per notification destination definition. See [Nested secure_url blocks](#nested-secureurl-blocks) below for details.
* `property` - (Required) A nested block that describes a notification destination property. See [Nested property blocks](#nested-property-blocks) below for details.

### Nested `auth_basic` blocks

* `user` - (Required) The username of the basic auth.
* `password` - (Required) Specifies an authentication password for use with a destination.

### Nested `auth_token` blocks

* `prefix` - (Required) The prefix of the token auth.
* `token` - (Required) Specifies the token for integrating.

### Nested `auth_custom_header` blocks

* `key` - (Required) The key of the header.
* `value` - (Required) The secret value of the header.

~> **NOTE:** OAuth2 authentication type is not available via terraform for notifications destinations.

### Nested `secure_url` blocks

* `prefix` - (Required) The prefix of the URL.
* `secure_suffix` - (Required) The suffix of the URL, which contains sensitive data.

### Nested `property` blocks

* `key` - (Required) The notification property key.
* `value` - (Required) The notification property value.
* `label` - (Optional) The notification property label.
* `display_value` - (Optional) The notification property display value.

Each notification destination type supports a specific set of arguments for the `property` block. See [Additional Examples](#additional-examples) below for details:

* `EMAIL`
  * `email` - (Required) A list of email addresses.
* `WEBHOOK`
  * `url` - (Required) The webhook url.
* `SERVICE_NOW`
  * `url` - (Required) The service now destination url (only base url).
  * `two_way_integration` - (Optional) A boolean that represents the two-way integration on/off flag.
* `SERVICE_NOW_APP`
  * `url` - (Required) The service now destination url (only base url).
* `JIRA`
  * `url` - (Required) The jira url (only base url).
  * `two_way_integration` - (Optional) A boolean that represents the two-way integration on/off flag.
* `PAGERDUTY_ACCOUNT_INTEGRATION`
  * `two_way_integration` - (Optional) A boolean that represents the two-way integration on/off flag.
* `MOBILE_PUSH`
  * `userId` - (Required) The new relic user id.
* `EVENT_BRIDGE`
  * `AWSAccountId` - (Required) The account id to integrate to.
  * `AWSRegion` - (Required) The AWS region this account is in.
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the destination.
* `guid` - The unique entity identifier of the destination in New Relic.

## Additional Examples

~> **NOTE:** We support all properties. The mentioned properties are just an example.


##### [ServiceNow](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/notification-integrations/#servicenow)

```hcl
resource "newrelic_notification_destination" "foo" {
  account_id = 12345678
  name = "servicenow-example"
  type = "SERVICE_NOW"

  property {
    key = "url"
    value = "https://service-now.com/"
  }

  property {
    key = "two_way_integration"
    value = "true"
  }

  auth_basic {
    user = "username"
    password = "password"
  }
}
```

##### [Email](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/notification-integrations/#email)
```hcl
resource "newrelic_notification_destination" "foo" {
  account_id = 12345678
  name = "email-example"
  type = "EMAIL"

  property {
    key = "email"
    value = "email@email.com,email2@email.com"
  }
}
```

##### [Jira](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/notification-integrations/#jira)
```hcl
resource "newrelic_notification_destination" "foo" {
  account_id = 12345678
  name = "jira-example"
  type = "JIRA"

  property {
    key = "url"
    value = "https://example.atlassian.net"
  }

  auth_basic {
    user = "example@email.com"
    password = "password"
  }
}
```

##### [PagerDuty with service integration](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/notification-integrations/#pagerduty-sli)
```hcl
resource "newrelic_notification_destination" "foo" {
  account_id = 12345678
  name = "pagerduty-service-example"
  type = "PAGERDUTY_SERVICE_INTEGRATION"

  property {
    key = ""
    value = ""
  }

  auth_token {
    prefix = "Token token="
    token  = "10567a689d984d03c021034b22a789e2"
  }
}
```

##### [PagerDuty with account integration](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/notification-integrations/#pagerduty-ali)
```hcl
resource "newrelic_notification_destination" "foo" {
  account_id = 12345678
  name = "pagerduty-account-example"
  type = "PAGERDUTY_ACCOUNT_INTEGRATION"

  property {
    key = "two_way_integration"
    value = "true"
  }

  auth_token {
    prefix = "Token token="
    token  = "u+E8EU3MhsZwLfZ1ic1A"
  }
}
```

#### Mobile Push
```hcl
resource "newrelic_notification_destination" "foo" {
  account_id = 12345678
  name = "mobile-push-example"
  type = "MOBILE_PUSH"

  property {
    key = "userId"
    value = "12345678"
  }
}
```

#### [AWS Event Bridge](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/notification-integrations/#eventBridge)
```hcl
resource "newrelic_notification_destination" "foo" {
  account_id = 12345678
  name = "event-bridge-example"
  type = "EVENT_BRIDGE"

  property {
    key = "AWSAccountId"
    value = "123456789123456"
  }

  property {
    key = "AWSRegion"
    value = "us-east-2"
  }
}
```

#### [Slack](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/notification-integrations/#slack)

In order to create a Slack destination, you have to grant our application access to your workspace. This process is [based on OAuth](https://api.slack.com/authentication/oauth-v2) and can only be done through a browser.
As a result, you cannot set up a Slack destination purely with Terraform code.
However, if you would like to use Slack-based destinations with other resources in the New Relic Terraform Provider, the [data source `newrelic_notification_destination`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/data-sources/notification_destination) may be used to fetch the ID of the destination; alternatively, you might want to source the ID of the destination from  NerdGraph, or from the New Relic One UI.

## Import

~> **WARNING:** Slack-based destinations can only be imported and destroyed; this resource **does not** support creating and updating Slack-based destinations, owing to the reasons stated above, under the **Slack** section.

Destination id can be found in the Destinations page -> three dots at the right of the chosen destination -> copy destination id to clipboard.
This example is especially useful for slack destinations which *must* be imported.

1. Add an empty resource to your terraform file:
```terraform
resource "newrelic_notification_destination" "foo" {
}
```
2. Run import command: `terraform import newrelic_notification_destination.foo <destination_id>`
3. Run the following command after the import successfully done and copy the information to your resource:
   `terraform state show newrelic_notification_destination.foo`
4. Add `ignore_changes` attribute on `all` in your imported resource:
```terraform
lifecycle {
    ignore_changes = all
  }
```

Your imported destination should look like that:
```terraform
resource "newrelic_notification_destination" "foo" {
  lifecycle {
    ignore_changes = all
  }
  name = "*********"
  type = "SLACK"
  auth_token {
    prefix = "Bearer"
  }
  
  property {
      key   = "teamName"
      label = "Team Name"
      value = "******"
  }
}
```

~> **NOTE:** Sensitive data such as destination API keys, service keys, auth object etc. are not returned from the underlying API for security reasons and may not be set in state when importing.

## Additional Information
More information about destinations integrations can be found in NewRelic [documentation](https://docs.newrelic.com/docs/alerts-applied-intelligence/notifications/notification-integrations/).
More details about the destinations API can be found [here](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-api-notifications-destinations).

### Moving from Legacy Alert Channels to Notification Channels
As stated in the documentation of this resource and [`newrelic_notification_channel`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/notification_channel), destinations, created using the resource [`newrelic_notification_destination`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/notification_destination) can be paired with [`newrelic_notification_channel`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/notification_channel) to set up channels. These resources combined, are an alternative to the legacy resource [`newrelic_alert_channel`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/alert_channel), which is **deprecated** and will be **removed in a future major release**, as stated in the documentation of the resource.

If you're currently using `newrelic_alert_channel` to manage channels, we **strongly recommend** migrating to these notifications-based resources at the earliest.

Please refer to the examples in this page, or [this example](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/getting_started#add-a-notification-channel) for illustrations on setting up channels with these resources.