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
* [Install Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli) and read the Terraform getting started guide that follows. This guide will assume a basic understanding of Terraform.
* Locate your User API key by following [New Relic's User API key docs][user_api_key] (previously referred to as "personal API key"). <br> <small>Note: To to locate or create your API key via our NerdGraph API, please follow these [instructions][user_api_key_via_nerdgraph].</small>

## Configuring the Provider

Please see the [latest provider configuration docs](provider_configuration.html) to get started with configuring the provider.

## Initialize Your Terraform Setup

Once the provider is configured, you should be able to initialize your Terraform configuration, so let's give it a try.

```bash
$ terraform init
```

-> <small>This is the first command that should be run for any new or existing Terraform configuration per machine. This sets up all the local data necessary to run Terraform that is typically not committed to version control. This command is always safe to run multiple times.</small>

Once you've successfully initialized your Terraform working directory, you'll want to create your first configuration file (.tf file). Let's start by creating `main.tf`.

**main.tf**

```hcl
# Configure terraform
terraform {
  required_version = "~> 1.0"
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
    }
  }
}

# Configure the New Relic provider
provider "newrelic" {
  account_id = <Your Account ID>
  api_key = <Your User API Key>    # usually prefixed with 'NRAK'
  region = "US"                    # Valid regions are US and EU
}
```

> <small>**Note:** You can also use [environment variables](provider_configuration.html#configuration-via-environment-variables) to configure the provider, which can simplify your `provider` block.</small>

Now let's try running the following command.

```bash
$ terraform plan
```

This command will output some information into your console regarding Terraform's execution plan. Running `terraform plan` is essentially a "dry run" and will not provision anything. We'll get to provisioning in the next steps.

## Add an alert condition to an alert policy

We started with a minimal configuration with an Alert Policy, but it doesn't contain any alert conditions. Let's add a NRQL alert condition to that policy which we'll associate the condition to an application.

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


Now let's add the alert condition so we can see an alert when a particular scenario occurs.

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

# NRQL alert condition
resource "newrelic_nrql_alert_condition" "foo" {
  policy_id                    = newrelic_alert_policy.alert_policy_name.id
  type                         = "static"
  name                         = "foo"
  description                  = "Alert when transactions are taking too long"
  runbook_url                  = "https://www.example.com"
  enabled                      = true
  violation_time_limit_seconds = 3600

  nrql {
    query             = "SELECT average(duration) FROM Transaction where appName = '${data.newrelic_entity.app_name.name}'"
  }

  critical {
    operator              = "above"
    threshold             = 5.5
    threshold_duration    = 300
    threshold_occurrences = "ALL"
  }
}
```

This alert condition will be triggered when the average response duration of your application rises above the threshold of `5.5` for 5 minutes. This alert is considered `critical` in priority based on the configuration.

But how will you actually be alerted if this scenario occurs? To control where to send alerts, we'll need to configure a Notification Channel for the alert.

## Add a Notification Channel

New Relic alerts are great, but they're even better when combined with good notifications. To wire up a notification to your previously configured alert condition and alert policy, add the following to your configuration file.

```hcl
# Notification channel
resource "newrelic_notification_destination" "sample_notification_destination" {
  name = "Sample Notification Destination"
  type = "EMAIL"
  
  property {
    key   = "email"
    value = "username@example.com"
  }
}

resource "newrelic_notification_channel" "sample_notification_channel" {
  name           = "Sample Notification Channel"
  type           = "EMAIL"
  destination_id = newrelic_notification_destination.sample_notification_destination.id
  product        = "IINT"
  
  property {
    key   = "subject"
    value = "Sample Email Subject"
  }
}

resource "newrelic_workflow" "sample_workflow" {
  name                  = "Sample Workflow"
  muting_rules_handling = "NOTIFY_ALL_ISSUES"
  
  issues_filter {
    name = "Issue Filter"
    type = "FILTER"
    predicate {
      attribute = "labels.policyIds"
      operator  = "EXACTLY_MATCHES"
      values    = [newrelic_alert_policy.alert_policy_name.id]
    }
  }
  
  destination {
    channel_id = newrelic_notification_channel.sample_notification_channel.id
  }
}
```

This example explains setting up a workflow that is linked to the alert policy (created in the first example) and to an 'email' notification channel, which would enable sending an email to the specified recipients whenever the associated alert condition is triggered. 

If you would like to send notifications via different modalities, such as Slack, the [notification destination](/providers/newrelic/newrelic/latest/docs/resources/notification_destination) resource supports multiple types of destinations, and so does the [notification channel](/providers/newrelic/newrelic/latest/docs/resources/notification_channel) resource. 

Please refer to these pages comprising [additional examples on notification destinations](/providers/newrelic/newrelic/latest/docs/resources/notification_destination#additional-examples) and [additional examples on notification channels](/providers/newrelic/newrelic/latest/docs/resources/notification_channel#additional-examples), for more examples on configuring different types of notification channels (via notification destinations).

## Apply Your Terraform Configuration

To summarize, so far we've configured an alert policy that contains an alert condition that is associated with a specific application, but we haven't actually provisioned these resources in our New Relic account. Let's do that now.

To apply your configuration and provision these resources in your New Relic account, run the following command.

```bash
$ terraform apply
```

Follow the prompt, which should involve you answering `yes` to apply the changes. Terraform will then provision the resources.

Once complete, you'll be able to navigate to your Alerts tab in your New Relic account and click on Alert Policies. You should see your newly created alert policy. Clicking on the alert policy should display the associated alert condition that we just configured as well. You should also be able to see the workflow and the email notification destination created, upon navigating to **Alerts > Workflows** and **Alerts > Destinations** respectively, in the UI.

If you ever need to make changes to your configuration, you can run `terraform apply` again after saving your latest configuration and Terraform will update the proper resources with your changes.

You can also run `terraform destroy` to tear down your resources if that's ever needed.

[user_api_key]: https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys#user-api-key
[user_api_key_via_nerdgraph]: https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-user-keys
[terraform_data_sources]: https://www.terraform.io/language/data-sources#data-sources
