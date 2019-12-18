---
layout: "newrelic"
page_title: "Getting Started with the New Relic provider"
sidebar_current: "docs-newrelic-provider-getting-started"
description: |-
  Getting started with the New Relic provider
---

# Getting Started with the New Relic Provider

## Before You Begin

* This guide assumes you already have a New Relic agent deployed. If you don't have New Relic integrated yet, check out [New Relic's introduction documentation](https://docs.newrelic.com/docs/using-new-relic/welcome-new-relic/get-started/introduction-new-relic) to get started there, then head back over here to get started with the New Relic Terraform Provider.
* [Install Terraform](https://www.terraform.io/intro/getting-started/install.html)
and read the Terraform getting started guide that follows. This guide will
assume a basic understanding of Terraform.
* Locate your Admin's API key by following [New Relic's Admin API key docs](https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys#admin).

## First things first: Credentials

The environment variable `NEWRELIC_API_KEY` is automatically detected by the New Relic Terraform Provider when running `terraform` commands, so we recommend adding this environment variable to your machine's startup file, such as your `.bash_profile`.

This guide assumes your API key has been set with an environment variable.

-> <sup>You can set the environment variable `NEWRELIC_API_KEY` in your `.bash_profile` or `.bashrc` file (on UNIX machines). Or you can set the variable inline with the `terraform plan` or `terraform apply` commands (see examples below).</sup>

**.bash_profile**

```bash
# Add this to your .bash_profile
export NEWRELIC_API_KEY=abc123
```

Example inline with `terraform` command

```bash
$ NEWRELIC_API_KEY=abc123 terraform apply
```

## Configuring the Provider

Let's start with a minimal Terraform config file to create an Alert Policy.

**main.tf**

```hcl
provider "newrelic" {}

resource "newrelic_alert_policy" "my_alert_policy_name" {
  name = "My Alert Policy Name"
}
```
We'll add an Alert Condition under this policy as we move through this guide.


## Initialize Your Terraform Setup

At this point you should be able to initialize your Terraform setup, so let's give it a try.

```bash
$ terraform init
```

-> <sup>This is the first command that should be run for any new or existing Terraform configuration per machine. This sets up all the local data necessary to run Terraform that is typically not committed to version control. This command is always safe to run multiple times.</sup>

Once you've successfully initialized your Terraform working directory, try running the following command.

```bash
$ terraform plan
```

This command will output some information into your console regarding Terraform's execution plan. Running `terraform plan` is essentially a "dry run" and will not provision anything. We'll get to provisioning in the next steps.

## Add an Alert Condition to an Alert Policy

We started with a minimal configuration with an Alert Policy, but it doesn't contain any Alert Conditions. Let's add an Alert Condition to that policy which we'll associate the condition to an application.

First, let's add a data source by adding a `data` block. This will store your application's information for Terraform to use.

```hcl
provider "newrelic" {}

# Data Source
data "newrelic_application" "app" {
  name = "my-app" # Note: This must be an exact match of your app name in New Relic
}

resource "newrelic_alert_policy" "my_alert_policy_name" {
  name = "My Alert Policy Name"
}
```

-> Terraform's data sources are read-only views into pre-existing data, or they compute new values on the fly within Terraform itself. More information on data sources can be found [here](https://www.terraform.io/docs/configuration-0-11/data-sources.html).


Now let's add the Alert Condition so we can see an alert when a particular scenario occurs.

```hcl
provider "newrelic" {}

data "newrelic_application" "app" {
  name = "my-app"
}

resource "newrelic_alert_policy" "alert_policy_name" {
  name = "My Alert Policy Name"
}

# Alert Condition
resource "newrelic_alert_condition" "alert_condition_name" {
  policy_id = newrelic_alert_policy.my_alert_policy_name.id

  name            = "My Alert Condition Name"
  type            = "apm_app_metric"
  entities        = [data.newrelic_application.app_name.id]
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
  name = "paul@example.com"
  type = "email"

  configuration = {
    recipients              = "paul@example.com"
    include_json_attachment = "1"
  }
}

# Link the above notification channel to your policy
resource "newrelic_alert_policy_channel" "alert_policy_email" {
  policy_id  = newrelic_alert_policy.my_alert_policy_name.id
  channel_id = newrelic_alert_channel.alert_notification_email.id
}
```

This example will send an email to the specified recipients whenever the associated alert condition is triggered. If you would like to send notifications via different modalities, such as Slack, you can configure updating the `type` in your [alert channel](https://www.terraform.io/docs/providers/newrelic/r/alert_channel.html).


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
