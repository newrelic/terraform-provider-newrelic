---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_aws_eu_sovereign_integrations"
sidebar_current: "docs-newrelic-resource-cloud-aws-eu-sovereign-integrations"
description: |-
    Integrate AWS EU Sovereign services with New Relic.
---

# Resource: newrelic\_cloud\_aws\_eu\_sovereign\_integrations

Use this resource to integrate AWS EU Sovereign services with New Relic.

## Prerequisite

Setup is required for this resource to work properly. This resource assumes you have [linked an AWS EU Sovereign account](cloud_aws_eu_sovereign_link_account.html) to New Relic.

The New Relic AWS EU Sovereign integration relies on two mechanisms to get data into New Relic:

* **CloudWatch Metric Streams (PUSH)**: This is the supported method for AWS EU Sovereign Cloud to get metrics into New Relic for the majority of AWS services. Follow the [steps outlined here](https://docs-preview.newrelic.com/docs/aws-eu-sovereign-cloud-integration) to set up a metric stream.

* **API Polling (PULL)**: Required for services that are **not supported** by CloudWatch Metric Streams. The following three services must be integrated via API Polling: **Billing**, **CloudTrail** and **X-Ray**. Follow the [steps outlined here](https://docs-preview.newrelic.com/docs/aws-eu-sovereign-cloud-integration).

This resource is used to configure API Polling integrations for those three services that are not available through AWS CloudWatch Metric Streams.

## Example Usage

The following example demonstrates the use of the `newrelic_cloud_aws_eu_sovereign_integrations` resource with multiple AWS EU Sovereign integrations supported by the resource.

To view a full example with all supported AWS EU Sovereign integrations, please see the [Additional Examples](#additional-examples) section. Integration blocks used in the resource may also be left empty to use the default configuration of the integration.

A full example, inclusive of setup of AWS resources (from the AWS Terraform Provider) associated with this resource, may be found in our [AWS EU Sovereign cloud integration guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/cloud_integrations_guide#aws-eu-sovereign).

```hcl
resource "newrelic_cloud_aws_eu_sovereign_link_account" "foo" {
  arn                    = "arn:aws-eusc:iam::123456789012:role/NewRelicInfrastructure-Integrations"
  metric_collection_mode = "PULL"
  name                   = "my-eu-sovereign-account"
}

resource "newrelic_cloud_aws_eu_sovereign_integrations" "bar" {
  linked_account_id = newrelic_cloud_aws_eu_sovereign_link_account.foo.id

  billing {
    metrics_polling_interval = 3600
  }

  cloudtrail {
    metrics_polling_interval = 300
    aws_regions              = ["eusc-de-east-1"]
  }

  x_ray {
    metrics_polling_interval = 300
    aws_regions              = ["eusc-de-east-1"]
  }
}
```

## Supported AWS EU Sovereign Integrations

-> **NOTE:** CloudWatch Metric Streams is the only supported method for AWS EU Sovereign Cloud. The following three integrations are for services **not supported by CloudWatch Metric Streams** and must be configured via API Polling using this resource.

<details>
  <summary>Expand this section to view all supported AWS EU Sovereign services that may be integrated via this resource.</summary>

| Block                  | Description                   |
|------------------------|-------------------------------|
| `billing`              | Billing Integration           |
| `cloudtrail`           | CloudTrail Integration        |
| `x_ray`                | X-Ray Integration             |

</details>

## Argument Reference

-> **WARNING:** Updating the `linked_account_id` of a `newrelic_cloud_aws_eu_sovereign_integrations` resource that has been applied would **force a replacement** of the resource (destruction of the resource, followed by the creation of a new resource). Please carefully review the output of `terraform plan`, which would clearly indicate a replacement of this resource, before performing a `terraform apply`.

* `account_id` - (Optional) The New Relic account ID to operate on. This allows the user to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.
* `linked_account_id` - (Required) The ID of the linked AWS EU Sovereign account in New Relic.

### Arguments to be Specified with Integration Blocks

The following arguments are intended to be used within certain ["integration blocks"](#integration-blocks) in the resource, i.e. these are supposed to be specified as nested arguments "within" an argument corresponding to a specific AWS EU Sovereign integration, unlike `account_id` and `linked_account_id` which are specified at the resource level. An exhaustive list of all of such arguments supported by the resource (and the integration blocks they would need to be specified with) has been given below.

In order to find the right set of arguments which go with each integration, and samples on the usage of these arguments "within" integration blocks, **check out the [Integration Blocks](#integration-blocks) section below.**

* `metrics_polling_interval` - (Optional) The data polling interval **in seconds**.
  * The following integration blocks support the usage of this argument:

        |                   |                   |                   |
        |-------------------|-------------------|-------------------|
        | `billing`         | `cloudtrail`      | `x_ray`           |

-> **NOTE** You may find the range of metric polling intervals of an integration under the [Integration Blocks](#integration-blocks) section.

* `aws_regions` - (Optional) Specify each AWS EU Sovereign region that includes the resources that you want to monitor.
  * The following integration blocks support the usage of this argument:

        |                   |                   |                   |
        |-------------------|-------------------|-------------------|
        | `cloudtrail`      | `x_ray`           |                   |

  * Valid regions: `eusc-de-east-1`

## Integration Blocks

The following section lists out arguments which may be used with each AWS EU Sovereign integration supported by this resource.

As specified above in the [Arguments to be Specified with Integration Blocks](#arguments-to-be-specified-with-integration-blocks) section, except for `linked_account_id` and `account_id`, all aforementioned arguments are to be specified within an integration block as they are supported by a specific set of integrations each; the following list of integration blocks elucidates the same with samples of what each integration block would look like.

<details>
  <summary>Expand this list to see a list of all integration blocks supported by this resource, the arguments which go with them and a sample of what the block would look like with these arguments.</summary>
  <details>
    <summary>billing</summary>

*  Supported Arguments: `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 3600, 21600 (seconds)
```hcl
     billing {
        metrics_polling_interval = 3600
     }
```
  </details>
  <details>
    <summary>cloudtrail</summary>

*  Supported Arguments: `aws_regions`, `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl
     cloudtrail {
        metrics_polling_interval = 300
        aws_regions              = ["eusc-de-east-1"]
     }
```
  </details>
  <details>
    <summary>x_ray</summary>

*  Supported Arguments: `aws_regions`, `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 60, 300, 900, 1800, 3600 (seconds)
```hcl
     x_ray {
        metrics_polling_interval = 300
        aws_regions              = ["eusc-de-east-1"]
     }
```
  </details>
</details>

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the AWS EU Sovereign linked account.

## Additional Examples

```hcl
resource "newrelic_cloud_aws_eu_sovereign_integrations" "bar" {
  linked_account_id = newrelic_cloud_aws_eu_sovereign_link_account.foo.id

  billing {
    metrics_polling_interval = 300
  }

  cloudtrail {
    metrics_polling_interval = 900
    aws_regions              = ["eusc-de-east-1"]
  }

  x_ray {
    metrics_polling_interval = 300
    aws_regions              = ["eusc-de-east-1"]
  }
}
```

## Import

Linked AWS EU Sovereign account integrations can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_cloud_aws_eu_sovereign_integrations.foo <id>
```