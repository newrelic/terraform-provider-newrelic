---
layout: "newrelic"
page_title: "Provider: New Relic"
sidebar_current: "docs-newrelic-index"
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

| Argument               | Required? | Description                                                                                                                                                                 |
| ---------------------- | --------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `account_id`           | Required  | Your New Relic account ID. The `NEW_RELIC_ACCOUNT_ID` environment variable can also be used.                                                                                |
| `api_key`              | Required  | Your New Relic Personal API key (usually prefixed with `NRAK`). The `NEW_RELIC_API_KEY` environment variable can also be used.                                              |
| `region`               | Required  | The region for the data center for which your New Relic account is configured. The `NEW_RELIC_REGION` environment variable can also be used. Valid values are `US` or `EU`. |
| `insecure_skip_verify` | Optional  | Trust self-signed SSL certificates. If omitted, the `NEW_RELIC_API_SKIP_VERIFY` environment variable is used.                                                               |
| `insights_insert_key`  | Optional  | Your Insights insert key used when inserting Insights events via the `newrelic_insights_event` resource. Can also use `NEW_RELIC_INSIGHTS_INSERT_KEY` environment variable. |
| `cacert_file`          | Optional  | A path to a PEM-encoded certificate authority used to verify the remote agent's certificate. The `NEW_RELIC_API_CACERT` environment variable can also be used.              |

## Authentication Requirements

This provider is in the midst of migrating away from our older REST based APIs
to a newer GraphQL based API that we lovingly call NerdGraph. During this
transition, the provider will be using different endpoints depending on which
resource is in use. Below is a table that reflects the current state of the
resources compared to which endpoint is in use.

### Resources

| Resource                                            | Endpoint                | Authentication        |
| --------------------------------------------------- | ----------------------- | --------------------- |
| `newrelic_alert_channel`                            | RESTv2                  | `api_key`             |
| `newrelic_alert_condition`                          | RESTv2                  | `api_key`             |
| `newrelic_alert_muting_rule`                        | NerdGraph               | `api_key`             |
| `newrelic_alert_policy`                             | NerdGraph               | `api_key`             |
| `newrelic_alert_policy_channel`                     | RESTv2                  | `api_key`             |
| `newrelic_api_access_key`                           | NerdGraph               | `api_key`             |
| `newrelic_application_settings`                     | RESTv2                  | `api_key`             |
| `newrelic_entity_tags`                              | NerdGraph               | `api_key`             |
| `newrelic_events_to_metrics_rule`                   | NerdGraph               | `api_key`             |
| `newrelic_infra_alert_condition`                    | Infrastructure REST API | `api_key`             |
| `newrelic_insights_event`                           | Insights API            | `insights_insert_key` |
| `newrelic_nrql_alert_condition`                     | NerdGraph               | `api_key`             |
| `newrelic_nrql_drop_rule`                           | NerdGraph               | `api_key`             |
| `newrelic_one_dashboard`                            | NerdGraph               | `api_key`             |
| `newrelic_service_level`                            | NerdGraph               | `api_key`             |
| `newrelic_synthetics_alert_condition`               | RESTv2                  | `api_key`             |
| `newrelic_synthetics_monitor`                       | Synthetics REST API     | `api_key`             |
| `newrelic_synthetics_monitor_script`                | Synthetics REST API     | `api_key`             |
| `newrelic_synthetics_multilocation_alert_condition` | RESTv2                  | `api_key`             |
| `newrelic_synthetics_secure_credential`             | Synthetics REST API     | `api_key`             |
| `newrelic_workload`                                 | NerdGraph               | `api_key`             |

### Data Sources

| Data Source                             | Endpoint            | Authentication |
| --------------------------------------- | ------------------- | -------------- |
| `newrelic_account`                      | NerdGraph           | `api_key`      |
| `newrelic_alert_channel`                | RESTv2              | `api_key`      |
| `newrelic_alert_policy`                 | NerdGraph           | `api_key`      |
| `newrelic_application`                  | RESTv2              | `api_key`      |
| `newrelic_entity`                       | NerdGraph           | `api_key`      |
| `newrelic_key_transaction`              | RESTv2              | `api_key`      |
| `newrelic_synthetics_monitor`           | Synthetics REST API | `api_key`      |
| `newrelic_synthetics_secure_credential` | Synthetics REST API | `api_key`      |

## Example Usage

```hcl
# Configure the New Relic provider
provider "newrelic" {
  account_id = <Your Account ID>
  api_key = <Your Personal API Key>    # usually prefixed with 'NRAK'
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
resource "newrelic_nrql_alert_condition" "foo" {
  policy_id                    = newrelic_alert_policy.alert.id
  type                         = "static"
  name                         = "foo"
  description                  = "Alert when transactions are taking too long"
  runbook_url                  = "https://www.example.com"
  enabled                      = true
  value_function               = "single_value"
  violation_time_limit_seconds = 3600

  nrql {
    query             = "SELECT average(duration) FROM Transaction where appName = '${data.newrelic_entity.foo.name}'"
    evaluation_offset = 3
  }

  critical {
    operator              = "above"
    threshold             = 5.5
    threshold_duration    = 300
    threshold_occurrences = "ALL"
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

## Support for v2.x

While the sun rises on the `3.x` release, the sunset of the `2.x` approaches.
We intend to support minor bug fixes through the end of 2021, but we don't plan
to merge any new features into `release/2.x` branch. Please see the section
below about upgrading the provider. All new feature work and focus will be
directed at the newer provider version.

-> <small>**Deprecation notice:** 2021-06-16<br>
-> **End of support:** 2022-01-07</small>

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

Please see the [latest provider configuration docs](guides/provider_configuration.html) for the current recommended configuration settings.

## Support for v1.x

Support for v1.x ended on January 15th, 2021.

## Debugging

Additional debugging information can be generated by exporting the `TF_LOG` environment variable when running Terraform commands. See [Debugging Terraform](https://www.terraform.io/internals/debugging) for more information.

### HTTP Request logging

Setting `TF_LOG` to a value of `DEBUG` will generate request log messages from the underlying HTTP client, and a value of `TRACE` will add additional context to these messages, including request and response body and headers.

## Community

New Relic hosts and moderates an online forum where customers can interact with New Relic employees as well as other customers to get help and share best practices.

- [Roadmap](https://newrelic.github.io/developer-toolkit/roadmap/) - As part of the Developer Toolkit, the roadmap for this project follows the same RFC process
- [Issues or Enhancement Requests](https://github.com/newrelic/terraform-provider-newrelic/issues) - Issues and enhancement requests can be submitted in the Issues tab of this repository. Please search for and review the existing open issues before submitting a new issue.
- [Community discussion board](https://discuss.newrelic.com/c/build-on-new-relic/developer-toolkit) - Like all official New Relic open source projects, there's a related Community topic in the New Relic Explorers Hub.

Keep in mind that when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. If you'd like to execute our corporate CLA, or if you have any questions, please drop us an email at opensource@newrelic.com.

## Support

New Relic has open-sourced this project. This project is provided AS-IS WITHOUT WARRANTY OR SUPPORT, although you can report issues and contribute to the project here on GitHub.

_Please do not report issues with this software to New Relic Global Technical Support._

[provider_version_configuration]: https://www.terraform.io/language/providers/requirements#requiring-providers
