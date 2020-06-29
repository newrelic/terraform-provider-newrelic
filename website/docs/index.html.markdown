---
layout: 'newrelic'
page_title: 'Provider: New Relic'
sidebar_current: 'docs-newrelic-index'
description: |-
  New Relic offers a performance management solution enabling developers to
  diagnose and fix application performance problems in real time.
---

# New Relic Provider

[New Relic](https://newrelic.com/) offers tools that help you fix problems
quickly, maintain complex systems, improve your code, and accelerate your
digital transformation.

Use the navigation to the left to read about the available resources.

## Quick Links

- [**Configure the Provider**](guides/provider_configuration.html)
- [**Getting Started Guide**](guides/getting_started.html)
- [**Migration Guide: Upgrading to v2.x**](guides/migration_guide_v2.html)

## Argument Reference

The following arguments are supported.

| Argument | Required?         | Description                                                                                                                                                                            |
| ---------------------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `account_id`           | Required | Your New Relic account ID. The `NEW_RELIC_ACCOUNT_ID` environment variable can also be used.                                                                                |
| `api_key`              | Required | Your New Relic Personal API key (starts with `NRAK`). The `NEW_RELIC_API_KEY` environment variable can also be used.                                                        |
| `admin_api_key`        | Required | Your New Relic Admin API key (starts with `NRAA`). The `NEW_RELIC_ADMIN_API_KEY` environment variable can also be used.                                                     |
| `region`               | Required | The region for the data center for which your New Relic account is configured. The `NEW_RELIC_REGION` environment variable can also be used. Valid values are `US` or `EU`. |
| `insecure_skip_verify` | Optional | Trust self-signed SSL certificates. If omitted, the `NEW_RELIC_API_SKIP_VERIFY` environment variable is used.                                                               |
| `insights_insert_url`  | Optional | Your Insights insert key used when inserting Insights events via the `newrelic_insights_event` resource. Can also use `NEW_RELIC_INSIGHTS_INSERT_KEY` environment variable. |
| `cacert_file`          | Optional | A path to a PEM-encoded certificate authority used to verify the remote agent's certificate. The `NEW_RELIC_API_CACERT` environment variable can also be used.              |


## Example Usage

```hcl
# Configure the New Relic provider
provider "newrelic" {
  account_id = <Your Account ID>
  api_key = <Your Personal API Key>    # starts with 'NRAK'
  admin_api_key = <Your Admin API Key> # starts with 'NRAA'
  region = "US"                        # Valid regions are US and EU
}

# Read an APM application resource
data "newrelic_entity" "foo" {
  name = "Your App Name"
  domain = "APM"
  type = "APPLICATION"
}

# Create an alert policy
resource "newrelic_alert_policy" "alert" {
  name = "Your Concise Alert Name"
}

# Add a condition
resource "newrelic_alert_condition" "foo" {
  policy_id = newrelic_alert_policy.alert.id

  name        = "foo"
  type        = "apm_app_metric"
  entities    = [data.newrelic_entity.foo.application_id]
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

  config {
    recipients              = "username@example.com"
    include_json_attachment = "1"
  }
}

# Link the channel to the policy
resource "newrelic_alert_policy_channel" "alert_email" {
  policy_id  = newrelic_alert_policy.alert.id
  channel_ids = [
    newrelic_alert_channel.email.id
  ]
}
```

## Support for v1.x

While the sun rises on the `2.x` release, the sunset of the `1.x` approaches.
We intend to support minor bug fixes through the end of 2020, but we don't plan
to merge any new features into `release/1.x` branch.  Please see the section
below about upgrading the provider.  All new feature work and focus will be
directed at the newer provider version.

-> <small>**Deprecation notice:** 2020-06-12<br>
-> **End of support:** 2020-01-15</small>

If you wish to pin your environment to a specific release, you can do so with a `required_providers` statement in your Terraform manifest. You can also pin the version within your `provider` block.

Using the `provider` block:

```hcl
provider "newrelic" {
  version = "~> 2.0.0"
}
```

Using the `required_providers` block:

```hcl
required_providers {
  newrelic = "~> 2.0.0"
}
```

See the [Terraform docs][provider_version_configuration] for more information on pinning versions.

## Upgrading to v2.x

Upgrading to v2 of the provider involves some changes to your provider configuration. Please view our [**migration guide**](guides/migration_guide_v2.html) for more information and assistance.

Please see the [latest provider configuration docs](/docs/providers/newrelic/guides/provider_configuration.html) for the current recommended configuration settings.

## Resource endpoint authentication

This provider is in the midst of migrating away from our older REST based APIs to a newer GraphQL based API that we lovingly call NerdGraph.  During this transition, the provider will be using different endpoints depending on which resource is in use.  Below is a table that reflects the current state of the resources compared to which endpoint is in use.

| Resource                                       | RESTv2 | NerdGraph |
| ---------------------------------------------- | ------ | --------- |
| resource_newrelic_alert_channel                | yes    | no        |
| resource_newrelic_alert_condition              | yes    | no        |
| resource_newrelic_alert_policy                 | no     | yes       |
| resource_newrelic_alert_policy_channel         | yes    | no        |
| resource_newrelic_application_settings         | yes    | no        |
| resource_newrelic_dashboard                    | yes    | no        |
| resource_newrelic_infra_alert_condition        | yes    | no        |
| resource_newrelic_insights_event               | yes    | no        |
| resource_newrelic_nrql_alert_condition         | no     | yes       |
| resource_newrelic_plugins_alert_condition      | yes    | no        |
| resource_newrelic_synthetics_alert_condition   | yes    | no        |
| resource_newrelic_synthetics_label             | yes    | no        |
| resource_newrelic_synthetics_monitor           | yes    | no        |
| resource_newrelic_synthetics_monitor_script    | yes    | no        |
| resource_newrelic_synthetics_secure_credential | yes    | no        |
| resource_newrelic_workload                     | no     | yes       |


## Debugging

Additional debugging information can be generated by exporting the `TF_LOG` environment variable when running Terraform commands. See [Debugging Terraform](https://www.terraform.io/docs/internals/debugging.html) for more information.

### HTTP Request logging

Setting `TF_LOG` to a value of `DEBUG` will generate request log messages from the underlying HTTP client, and a value of `TRACE` will add additional context to these messages, including request and response body and headers.

## Community

New Relic hosts and moderates an online forum where customers can interact with New Relic employees as well as other customers to get help and share best practices.

* [Roadmap](https://newrelic.github.io/developer-toolkit/roadmap/) - As part of the Developer Toolkit, the roadmap for this project follows the same RFC process
* [Issues or Enhancement Requests](https://github.com/newrelic/terraform-provider-newrelic/issues) - Issues and enhancement requests can be submitted in the Issues tab of this repository. Please search for and review the existing open issues before submitting a new issue.
* [Community discussion board](https://discuss.newrelic.com/c/build-on-new-relic/developer-toolkit) - Like all official New Relic open source projects, there's a related Community topic in the New Relic Explorers Hub.

Keep in mind that when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. If you'd like to execute our corporate CLA, or if you have any questions, please drop us an email at opensource@newrelic.com.

## Support

New Relic has open-sourced this project. This project is provided AS-IS WITHOUT WARRANTY OR SUPPORT, although you can report issues and contribute to the project here on GitHub.

_Please do not report issues with this software to New Relic Global Technical Support._

[provider_version_configuration]: https://www.terraform.io/docs/configuration/providers.html#provider-versions
