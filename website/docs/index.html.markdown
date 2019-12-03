---
layout: 'newrelic'
page_title: 'Provider: New Relic'
sidebar_current: 'docs-newrelic-index'
description: |-
  New Relic offers a performance management solution enabling developers to
  diagnose and fix application performance problems in real time.
---

# New Relic Provider

[New Relic](https://newrelic.com/) offers a performance management solution
enabling developers to diagnose and fix application performance problems in real time.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the New Relic provider
provider "newrelic" {
  api_key = "${var.newrelic_api_key}"
}

# Create an alert policy
resource "newrelic_alert_policy" "alert" {
  name = "Alert"
}

# Add a condition
resource "newrelic_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.alert.id}"

  name        = "foo"
  type        = "apm_app_metric"
  entities    = ["12345"]                             # You can look this up in New Relic
  metric      = "apdex"
  runbook_url = "https://docs.example.com/my-runbook"

  term {
    duration      = 5
    operator      = "below"
    priority      = "critical"
    threshold     = "0.75"
    time_function = "all"
  }
}

# Add a notification channel
resource "newrelic_alert_channel" "email" {
  name = "email"
  type = "email"

  configuration = {
    recipients              = "paul@example.com"
    include_json_attachment = "1"
  }
}

# Link the channel to the policy
resource "newrelic_alert_policy_channel" "alert_email" {
  policy_id  = "${newrelic_alert_policy.alert.id}"
  channel_id = "${newrelic_alert_channel.email.id}"
}
```

## Argument Reference

The following arguments are supported:

- `api_key` - (Required except for `newrelic_insights_event` resource) Your New Relic API key. Can also use `NEWRELIC_API_KEY` environment variable.
- `api_url` - (Optional) Possibility to change the API url (default is https://api.newrelic.com/v2). This is for instance relevant if the New Relic account is in the EU. The API url must then be set to https://api.eu.newrelic.com/v2. Can also use `NEWRELIC_API_URL` environment variable.
- `infra_api_url` - (Optional) Possibility to change the Infra API url (default is https://infra-api.newrelic.com/v2). This is for instance relevant if the New Relic account is in the EU. The Infra API url must then be set to https://infra-api.eu.newrelic.com/v2. Can also use `NEWRELIC_INFRA_API_URL` environment variable.
- `insecure_skip_verify` - (Optional) Trust self-signed SSL certificates. If omitted, the `NEWRELIC_API_SKIP_VERIFY` environment variable is used.
- `insights_account_id` - (Optional) Your New Relic Account ID used when inserting Insights events via the `newrelic_insights_event` resource. Can also use `NEWRELIC_INSIGHTS_ACCOUNT_ID` environment variable.
- `insights_insert_key` - (Optional) Your Insights insert key used when inserting Insights events via the `newrelic_insights_event` resource. Can also use `NEWRELIC_INSIGHTS_INSERT_KEY` environment variable.
- `insights_insert_url` - (Optional) Possibility to change the Insights insert url (default is https://insights-collector.newrelic.com/v1/accounts). This is for instance relevant if the New Relic account is in the EU. The Insights API url must then be set to https://insights-collector.eu.newrelic.com/v1. Can also use `NEWRELIC_INSIGHTS_INSERT_URL` environment variable.
- `cacert_file` - (Optional) A path to a PEM-encoded certificate authority used to verify the remote agent's certificate. The `NEWRELIC_API_CACERT` environment variable can also be used.
