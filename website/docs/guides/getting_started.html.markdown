---
layout: "newrelic"
page_title: "Getting Started with the New Relic provider"
sidebar_current: "docs-newrelic-provider-getting-started"
description: |-
  Getting started with the New Relic provider
---

# Getting Started with the New Relic Provider


## Before You Begin

* The examples below assume you already have a New Relic agent deployed. If you don't have New Relic integrated yet, check out [New Relic's introduction documentation](https://docs.newrelic.com/docs/using-new-relic/welcome-new-relic/get-started/introduction-new-relic) to get started there, then head back over here to get started with the New Relic Terraform Provider using the examples provided.
* [Install Terraform](https://www.terraform.io/intro/getting-started/install.html) and read the Terraform getting started guide that follows. This guide will assume a basic understanding of Terraform.
* Locate your Personal API key by following [New Relic's Personal API key docs][personal_api_key].
* Locate your Admin's API key by following [New Relic's Admin API key docs][admin_api_key].

## Configuring the Provider

Please see the [latest provider configuration docs](/guides/provider_configuration.html) to get started with configuring the provider.

## Initialize Your Terraform Setup

Once the provider is configured, you should be able to initialize your Terraform configuration, so let's give it a try.

```bash
$ terraform init
```

-> <small>This is the first command that should be run for any new or existing Terraform configuration per machine. This sets up all the local data necessary to run Terraform that is typically not committed to version control. This command is always safe to run multiple times.</small>

Once you've successfully initialized your Terraform working directory, you'll want to create your first configuration file (.tf file). Let's start by creating `main.tf`.

**main.tf**

```hcl
# Configure the New Relic provider
provider "newrelic" {
  account_id = <Your Account ID>
  api_key = <Your Personal API Key>    # usually prefixed with 'NRAK'
  admin_api_key = <Your Admin API Key> # usually prefixed with 'NRAA'
  region = "US"                        # Valid regions are US and EU
}
```

> <small>**Note:** You can also use [environment variables](/guides/provider_configuration.html#configuration-via-environment-variables) to configure the provider, which can simplify your `provider` block.</small>

Now let's try running the following command.

```bash
$ terraform plan
```

This command will output some information into your console regarding Terraform's execution plan. Running `terraform plan` is essentially a "dry run" and will not provision anything. We'll get to provisioning in the next steps.

## Add an Alert Condition to an Alert Policy

We started with a minimal configuration with an Alert Policy, but it doesn't contain any Alert Conditions. Let's add an Alert Condition to that policy which we'll associate the condition to an application.

First, let's add a data source by adding a `data` block. This will store your application's information for Terraform to use. Terraform data sources operate like a `GET` request - they fetch the data of the resource that matches the criteria you provide, in the example below it's a New Relic application entity named "my-app".

```hcl
provider "newrelic" {
  # ...your configuration from the previous step
}

# Data Source
data "newrelic_entity" "app_name" {
  name = "my-app" # Note: This must be an exact match of your app name in New Relic (Case sensitive)
  type = "APPLICATION"
  domain = "APM"
}

resource "newrelic_alert_policy" "alert_policy_name" {
  name = "My Alert Policy Name"
}
```

-> Terraform's data sources are read-only views into pre-existing data, or they compute new values on the fly within Terraform itself. More information on data sources can be found [here][terraform_data_sources].


Now let's add the Alert Condition so we can see an alert when a particular scenario occurs.

```hcl
provider "newrelic" {}

data "newrelic_entity" "app_name" {
  name = "my-app"
  type = "APPLICATION"
  domain = "APM"
}

resource "newrelic_alert_policy" "alert_policy_name" {
  name = "My Alert Policy Name"
}

# Alert Condition
resource "newrelic_alert_condition" "alert_condition_name" {
  policy_id = newrelic_alert_policy.alert_policy_name.id

  name            = "My Alert Condition Name"
  type            = "apm_app_metric"
  entities        = [data.newrelic_entity.app_name.application_id]
  metric          = "apdex"
  runbook_url     = "https://www.example.com"
  condition_scope = "application"

  term {
    duration      = 5
    operator      = "below"
    priority      = "critical"
    threshold     = "0.75"
    time_function = "all"
  }
}
```

This alert condition will be triggered when the [Apdex score](https://docs.newrelic.com/docs/apm/new-relic-apm/apdex/apdex-measure-user-satisfaction) of your application falls below the threshold of `0.75` for 5 minutes. This alert is considered `critical` in priority based on the configuration.

But how will you actually be alerted if this scenario occurs? To control where to send alerts, we'll need to configure a Notification Channel for the alert.

## Add a Notification Channel

New Relic alerts are great, but they're even better when combined with good notifications. To wire up a notification to your previously configured Alert Condition and Alert Policy, add the following to your configuration file.

```hcl
# Notification channel
resource "newrelic_alert_channel" "alert_notification_email" {
  name = "username@example.com"
  type = "email"

  config {
    recipients              = "username@example.com"
    include_json_attachment = "1"
  }
}

# Link the above notification channel to your policy
resource "newrelic_alert_policy_channel" "alert_policy_email" {
  policy_id  = newrelic_alert_policy.alert_policy_name.id
  channel_ids = [
    newrelic_alert_channel.alert_notification_email.id
  ]
}
```

This example will send an email to the specified recipients whenever the associated alert condition is triggered. If you would like to send notifications via different modalities, such as Slack, the [alert channel](/docs/providers/newrelic/r/alert_channel.html) resource supports multiple types of channels. Use the [additional alert channel examples](/docs/providers/newrelic/r/alert_channel.html#additional-examples) for assistance with configuring your different notification channels.

## A Note About Secrets

As part of a `newrelic` resource, there is often some amount of configuration
that is required in order for a resource to reach its full potential.  In some
cases, once a given entity is created, API results will obscure the values for
items that are deemed to be secret.  As a result, Terraform is unable to make
an accurate detection of a resource state, and so marks a resource as changed
for every run.

Consider the following example.

```hcl
resource "newrelic_alert_channel" "slack" {
  name = "slack"
  type = "slack"

  config {
    channel = "test"
    url     = "https://hooks.slack.com/services/xxxx/xxxxx"
  }
}
```

The resource above yields the following Terraform plan.

    -/+ newrelic_alert_channel.slack (new resource required)
          id:                    "2344397" => <computed> (forces new resource)
          config.%:              "1" => "2" (forces new resource)
          config.channel:        <sensitive> => <sensitive> (attribute changed)
          config.url:            <sensitive> => <sensitive> (forces new resource)
          name:                  "slack" => "slack"
          type:                  "slack" => "slack"

To avoid the resource being marked as changed every run, the following can be
implemented for the resource.

```hcl
resource "newrelic_alert_channel" "slack" {
  ...
  lifecycle {
    ignore_changes = ["config"]
  }
}
  ...
```

This should avoid any of the configuration items from causing a change to the
resource.


## Apply Your Terraform Configuration

To summarize, so far we've configured an Alert Policy that contains an Alert Condition that is associated with a specific application, but we haven't actually provisioned these resources in our New Relic account. Let's do that now.

To apply your configuration and provision these resources in your New Relic account, run the following command.

```bash
$ terraform apply
```

Follow the prompt, which should involve you answering `yes` to apply the changes. Terraform will then provision the resources.

Once complete, you'll be able to navigate to your Alerts tab in your New Relic account and click on Alert Policies. You should see your newly created alert policy. Clicking on the alert policy should display the associated alert condition that we just configured as well.

If you ever need to make changes to your configuration, you can run `terraform apply` again after saving your latest configuration and Terraform will update the proper resources with your changes.

You can also run `terraform destroy` to tear down your resources if that's ever needed.

[personal_api_key]: https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys#personal-api-key
[admin_api_key]: https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys#admin
[terraform_data_sources]: https://www.terraform.io/docs/configuration/data-sources.html
