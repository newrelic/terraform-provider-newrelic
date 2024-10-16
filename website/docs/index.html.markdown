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

### Getting started

- [**Configure the Provider**](guides/provider_configuration.html)
- [**Getting Started Guide**](guides/getting_started.html)

### Additional Tooling

Use the **New Relic CodeStream** IDE extension to simplify your workflow. Access the official docs right inside your IDE and add resource templates via auto-complete.

Install New Relic CodeStream for [VS Code][codestream_vscode], [JetBrains][codestream_jetbrains], or [Visual Studio][codestream_visualstudio], and then look for the wrench icon at the top of the CodeStream pane.

### Data sources

- [**New Relic Cloud integrations example for AWS, GCP, Azure**](guides/cloud_integrations_guide.html)

### Advanced

#### Dashboards

- [**Part 1: Creating dashboards with Terraform and JSON templates**](https://newrelic.com/blog/how-to-relic/create-nr-dashboards-with-terraform-part-1)
- [**Part 2: Dynamically creating New Relic dashboards with Terraform**](https://newrelic.com/blog/how-to-relic/create-nr-dashboards-with-terraform-part-2)
- [**Part 3: Using Terraform to generate New Relic dashboards from NRQL queries**](https://newrelic.com/blog/how-to-relic/create-nr-dashboards-with-terraform-part-3)

## Argument Reference

The following arguments are supported.

| Argument               | Required? | Description                                                                                                                                                                                        |
| ---------------------- | --------- |----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `account_id`           | Required  | Your New Relic account ID. The `NEW_RELIC_ACCOUNT_ID` environment variable can also be used.                                                                                                       |
| `api_key`              | Required  | Your New Relic Personal API key (usually prefixed with `NRAK`). The `NEW_RELIC_API_KEY` environment variable can also be used.                                                                     |
| `region`               | Optional  | The region for the data center for which your New Relic account is configured. The `NEW_RELIC_REGION` environment variable can also be used. Valid values are `US` or `EU`. Default value is `US`. |
| `insecure_skip_verify` | Optional  | Trust self-signed SSL certificates. If omitted, the `NEW_RELIC_API_SKIP_VERIFY` environment variable is used.                                                                                      |
| `insights_insert_key`  | Optional  | Your Insights insert key used when inserting Insights events via the `newrelic_insights_event` resource. Can also use `NEW_RELIC_INSIGHTS_INSERT_KEY` environment variable.                        |
| `cacert_file`          | Optional  | A path to a PEM-encoded certificate authority used to verify the remote agent's certificate. The `NEW_RELIC_API_CACERT` environment variable can also be used.                                     |

## Authentication Requirements

This provider is in the midst of migrating away from our older REST based APIs
to a newer GraphQL based API that we lovingly call NerdGraph. During this
transition, the provider will be using different endpoints depending on which
resource is in use. Below is a table that reflects the current state of the
resources compared to which endpoint is in use.

### Resources

| Resource                                            | Endpoint                | Authentication        |
|-----------------------------------------------------|-------------------------|-----------------------|
| `newrelic_account_management`                       | NerdGraph               | `api_key`             |
| `newrelic_alert_channel`                            | RESTv2                  | `api_key`             |
| `newrelic_alert_condition`                          | RESTv2                  | `api_key`             |
| `newrelic_alert_muting_rule`                        | NerdGraph               | `api_key`             |
| `newrelic_alert_policy`                             | NerdGraph               | `api_key`             |
| `newrelic_alert_policy_channel`                     | RESTv2                  | `api_key`             |
| `newrelic_api_access_key`                           | NerdGraph               | `api_key`             |
| `newrelic_application_settings`                     | RESTv2                  | `api_key`             |
| `newrelic_browser_application`                      | NerdGraph               | `api_key`             |
| `newrelic_cloud_aws_govcloud_integrations`          | NerdGraph               | `api_key`             |
| `newrelic_cloud_aws_govcloud_link_account`          | NerdGraph               | `api_key`             |
| `newrelic_cloud_aws_integrations`                   | NerdGraph               | `api_key`             |
| `newrelic_cloud_aws_link_account`                   | NerdGraph               | `api_key`             |
| `newrelic_cloud_azure_integrations`                 | NerdGraph               | `api_key`             |
| `newrelic_cloud_azure_link_account`                 | NerdGraph               | `api_key`             |
| `newrelic_cloud_gcp_integrations`                   | NerdGraph               | `api_key`             |
| `newrelic_cloud_gcp_link_account`                   | NerdGraph               | `api_key`             |
| `newrelic_data_partition_rule`                      | NerdGraph               | `api_key`             |
| `newrelic_entity_tags`                              | NerdGraph               | `api_key`             |
| `newrelic_events_to_metrics_rule`                   | NerdGraph               | `api_key`             |
| `newrelic_group`                                    | NerdGraph               | `api_key`             |
| `newrelic_infra_alert_condition`                    | Infrastructure REST API | `api_key`             |
| `newrelic_insights_event`                           | Insights API            | `insights_insert_key` |
| `newrelic_key_transaction`                          | NerdGraph               | `api_key`             |
| `newrelic_log_parsing_rule`                         | NerdGraph               | `api_key`             |
| `newrelic_notification_channel`                     | NerdGraph               | `api_key`             |
| `newrelic_notification_destination`                 | NerdGraph               | `api_key`             |
| `newrelic_nrql_alert_condition`                     | NerdGraph               | `api_key`             |
| `newrelic_nrql_drop_rule`                           | NerdGraph               | `api_key`             |
| `newrelic_obfuscation_expression`                   | NerdGraph               | `api_key`             |
| `newrelic_obfuscation_rule`                         | NerdGraph               | `api_key`             |
| `newrelic_one_dashboard`                            | NerdGraph               | `api_key`             |
| `newrelic_one_dashboard_json`                       | NerdGraph               | `api_key`             |
| `newrelic_one_dashboard_raw`                        | NerdGraph               | `api_key`             |
| `newrelic_service_level`                            | NerdGraph               | `api_key`             |
| `newrelic_synthetics_alert_condition`               | RESTv2                  | `api_key`             |
| `newrelic_synthetics_broken_links_monitor`          | NerdGraph               | `api_key`             |
| `newrelic_synthetics_cert_check_monitor`            | NerdGraph               | `api_key`             |
| `newrelic_synthetics_monitor`                       | NerdGraph               | `api_key`             |
| `newrelic_synthetics_multilocation_alert_condition` | RESTv2                  | `api_key`             |
| `newrelic_synthetics_private_location`              | NerdGraph               | `api_key`             |
| `newrelic_synthetics_script_monitor`                | NerdGraph               | `api_key`             |
| `newrelic_synthetics_secure_credential`             | NerdGraph               | `api_key`             |
| `newrelic_synthetics_step_monitor`                  | NerdGraph               | `api_key`             |
| `newrelic_user`                                     | NerdGraph               | `api_key`             |
| `newrelic_workflow`                                 | NerdGraph               | `api_key`             |
| `newrelic_workload`                                 | NerdGraph               | `api_key`             |

### Data Sources

| Data Source                             | Endpoint  | Authentication |
| --------------------------------------- |-----------| -------------- |
| `newrelic_account`                      | NerdGraph | `api_key`      |
| `newrelic_alert_channel`                | RESTv2    | `api_key`      |
| `newrelic_alert_policy`                 | NerdGraph | `api_key`      |
| `newrelic_application`                  | RESTv2    | `api_key`      |
| `newrelic_cloud_account`                | NerdGraph | `api_key`      |
| `newrelic_entity`                       | NerdGraph | `api_key`      |
| `newrelic_key_transaction`              | NerdGraph | `api_key`      |
| `newrelic_notification_destination`     | NerdGraph | `api_key`      |
| `newrelic_obfuscation_expression`       | NerdGraph | `api_key`      |
| `newrelic_synthetics_private_location`  | NerdGraph | `api_key`      |
| `newrelic_synthetics_secure_credential` | NerdGraph | `api_key`      |
| `newrelic_test_grok_pattern`            | NerdGraph | `api_key`      |


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
  violation_time_limit_seconds = 3600

  nrql {
    query             = "SELECT average(duration) FROM Transaction where appName = '${data.newrelic_entity.foo.name}'"
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

## Support for v3.x

The v3.x version of the New Relic Terraform provider will get continued support from the New Relic Observability as Code team. We advise to always upgrade to the latest versions of the v3.x branch as support is only given for the most recent versions. Older versions of the v3.x branch could get deprecated when new product releases are made. In those cases a notice will be put on the repository and communication will be sent out to you by New Relic.

Please see the section below about upgrading to the latest versions of the provider. All new feature work and focus will be directed at the newer provider version. If you wish to pin your environment to a specific release, you can do so with a `required_providers` statement in your Terraform manifest.

Using the `required_providers` block:

```hcl
terraform {
  required_version = "~> 1.0"
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
    }
  }
}
```

See the [Terraform docs][provider_version_configuration] for more information on pinning versions.

## Upgrading to v3.x

Upgrading to v3 of the provider involves some changes to your provider configuration. Please view our [**migration guide**](guides/migration_guide_v3.html) for more information and assistance.

Please see the [latest provider configuration docs](guides/provider_configuration.html) for the current recommended configuration settings.

## Support for v2.x

The v2.x version of the New Relic Terraform provider is supported as is and will not receive new features or bugfixes. The v2 release or older versions of the v2.x branch could get deprecated when new product releases are made. In those cases a notice will be put on the repository and communication will be sent out to you by New Relic.

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

- [Issues or Enhancement Requests](https://github.com/newrelic/terraform-provider-newrelic/issues) - Issues and enhancement requests can be submitted in the Issues tab of this repository. Please search for and review the existing open issues before submitting a new issue.
- [Community discussion board](https://discuss.newrelic.com/c/build-on-new-relic/developer-toolkit) - Like all official New Relic open source projects, there's a related Community topic in the New Relic Explorers Hub.

Keep in mind that when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. If you'd like to execute our corporate CLA, or if you have any questions, please drop us an email at opensource@newrelic.com.

## Support

New Relic has open-sourced this project. This project is provided AS-IS WITHOUT WARRANTY OR SUPPORT, although you can report issues and contribute to the project here on GitHub.

[provider_version_configuration]: https://www.terraform.io/language/providers/requirements#requiring-providers
[codestream_vscode]: https://marketplace.visualstudio.com/items?itemName=CodeStream.codestream
[codestream_jetbrains]: https://plugins.jetbrains.com/plugin/12206-new-relic-codestream
[codestream_visualstudio]: https://marketplace.visualstudio.com/items?itemName=CodeStream.codestream-vs-22
