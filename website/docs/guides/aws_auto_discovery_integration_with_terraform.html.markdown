---
layout: "newrelic"
page_title: "New Relic Integration with AWS Infra 360"
sidebar_current: "docs-newrelic-provider-aws-auto-discovery-integration-guide"
description: |-
  Use this guide to set up New Relic AWS Infra 360 (formerly AWS Auto Discovery) fully automated through Terraform.
---

# New Relic Integration with AWS Infra 360

AWS Infra 360 (formerly known as AWS Auto Discovery) discovers the resources in your AWS account and the relationships between them. The guide below shows how to enable it with Terraform, whether your AWS integration is managed by Terraform or was created some other way.

~> **NOTE:** The Terraform block is still named `aws_auto_discovery` for backwards compatibility. It configures AWS Infra 360.

## Documentation Reference
* [AWS Auto Discovery](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/set-up-auto-discovery-of-aws-entities/)

## Pre-Requisites
AWS Infra 360 is added on top of an existing New Relic AWS integration, which is either API Polling or CloudWatch Metric Streams. If you have not integrated your AWS account yet, start with the [Cloud Integrations guide](/providers/newrelic/newrelic/latest/docs/guides/cloud_integrations_guide) and the `newrelic_cloud_aws_link_account` resource.

You need:
* The name of the IAM role that the integration uses (the role name, not its ARN).
* The New Relic linked account ID of the integration. This is New Relic's ID for the linked AWS account, **not** your AWS account ID. You can find it in the New Relic UI under the integration's details, or with NerdGraph:

```graphql
{
  actor {
    account(id: YOUR_NR_ACCOUNT_ID) {
      cloud {
        linkedAccounts {
          id
          name
          externalId
        }
      }
    }
  }
}
```
Match `externalId` against your AWS account ID.

## Variables
* **New Relic Region** `new_relic_region` - The region where your New Relic account is hosted.
  * **Type** - String
  * **Allowed Values** - US, EU, or JP
* **New Relic Account Id** `new_relic_account_id` - The ID of your New Relic account.
  * **Type** - Integer
  * **Example** - 1234567
* **AWS IAM Role Name** `aws_iam_role_name` - The name of the IAM role used to integrate AWS with New Relic.
  * **Type** - String
  * **Example** - "NewRelicInfrastructure-Integrations"
* **Linked Account Id** `linked_account_id` - The New Relic linked account ID of the AWS integration on which AWS Infra 360 is to be enabled.
  * **Type** - Integer
  * **Example** - 123456
* **Scan Interval** `metrics_polling_interval` - How often, in seconds, AWS Infra 360 scans the AWS account.
  * **Type** - Integer
  * **Allowed Values** - 28800 (8 hours), 43200 (12 hours), 86400 (24 hours)
* **AWS Regions** `aws_regions` - The AWS regions to be scanned.
  * **Type** - List of strings
  * **Allowed Values** - Any of the following AWS commercial regions:
    * US: `us-east-1`, `us-east-2`, `us-west-1`, `us-west-2`
    * Canada and Mexico: `ca-central-1`, `ca-west-1`, `mx-central-1`
    * South America: `sa-east-1`
    * Europe: `eu-central-1`, `eu-central-2`, `eu-north-1`, `eu-south-1`, `eu-south-2`, `eu-west-1`, `eu-west-2`, `eu-west-3`
    * Asia Pacific: `ap-east-1`, `ap-east-2`, `ap-northeast-1`, `ap-northeast-2`, `ap-northeast-3`, `ap-south-1`, `ap-south-2`, `ap-southeast-1`, `ap-southeast-2`, `ap-southeast-3`, `ap-southeast-4`, `ap-southeast-5`, `ap-southeast-6`, `ap-southeast-7`
    * Middle East and Africa: `af-south-1`, `il-central-1`, `me-central-1`
  * **Example** - ["us-east-1", "us-west-2"]
  * Regions that are opt-in in AWS (for example `af-south-1`, `ap-east-1` and `ap-southeast-3`) must be enabled in your AWS account before they can be scanned. AWS GovCloud, China, the AWS European Sovereign Cloud and `me-south-1` are not supported.

## IAM Permissions
AWS Infra 360 uses CloudFormation resource scans, so the IAM role of the integration needs the following permissions in addition to the ones it already has:

* `cloudformation:StartResourceScan`
* `cloudformation:ListResourceScans`
* `cloudformation:DescribeResourceScan`
* `cloudformation:ListResourceScanResources`
* `cloudformation:ListResourceScanRelatedResources`

## Enabling AWS Infra 360 with Terraform

### Integration created outside of Terraform
Use this when the integration was created from the New Relic UI, CloudFormation, the CLI, or anything other than this Terraform configuration. It works for both API Polling and CloudWatch Metric Streams integrations.

```hcl
provider "newrelic" {
  region     = var.new_relic_region
  account_id = var.new_relic_account_id
}

data "aws_iam_policy_document" "newrelic_aws_infra_360" {
  statement {
    effect = "Allow"
    actions = [
      "cloudformation:StartResourceScan",
      "cloudformation:ListResourceScans",
      "cloudformation:DescribeResourceScan",
      "cloudformation:ListResourceScanResources",
      "cloudformation:ListResourceScanRelatedResources",
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "newrelic_aws_infra_360" {
  name        = "NewRelicAwsInfra360Permissions"
  description = "IAM policy for New Relic AWS Infra 360"
  policy      = data.aws_iam_policy_document.newrelic_aws_infra_360.json
}

resource "aws_iam_role_policy_attachment" "newrelic_aws_infra_360" {
  role       = var.aws_iam_role_name
  policy_arn = aws_iam_policy.newrelic_aws_infra_360.arn
}

resource "newrelic_cloud_aws_integrations" "infra_360" {
  linked_account_id = var.linked_account_id

  aws_auto_discovery {
    aws_regions              = var.aws_regions
    metrics_polling_interval = var.metrics_polling_interval
  }

  # New Relic validates the permissions when the integration is configured.
  depends_on = [aws_iam_role_policy_attachment.newrelic_aws_infra_360]
}
```

### Integration managed by Terraform
If the integration is already managed by Terraform, with a `newrelic_cloud_aws_link_account` resource and an `aws_iam_role` for it, attach the same policy to that role and reference the resources directly instead of using variables:

```hcl
resource "aws_iam_role_policy_attachment" "newrelic_aws_infra_360" {
  role       = aws_iam_role.newrelic_aws_role.name
  policy_arn = aws_iam_policy.newrelic_aws_infra_360.arn
}

resource "newrelic_cloud_aws_integrations" "infra_360" {
  linked_account_id = newrelic_cloud_aws_link_account.newrelic_cloud_integration_pull.id

  aws_auto_discovery {
    aws_regions              = ["us-east-1", "us-west-2"]
    metrics_polling_interval = 86400
  }

  depends_on = [aws_iam_role_policy_attachment.newrelic_aws_infra_360]
}
```

~> **NOTE:** If you already use a `newrelic_cloud_aws_integrations` resource for the same linked account, you can add the `aws_auto_discovery` block to that resource instead of creating a new one.

### Scanning every region
To scan all the regions that are enabled in your AWS account, build the list from the AWS provider:

```hcl
data "aws_regions" "enabled" {
  all_regions = false # only regions enabled by default
}

resource "newrelic_cloud_aws_integrations" "infra_360" {
  linked_account_id = var.linked_account_id

  aws_auto_discovery {
    aws_regions              = data.aws_regions.enabled.names
    metrics_polling_interval = 86400
  }

  depends_on = [aws_iam_role_policy_attachment.newrelic_aws_infra_360]
}
```

## Notes
* CloudWatch Metric Streams are regional. `aws_regions` only controls which regions AWS Infra 360 scans; it is independent of the region of your Metric Stream.
* The scan interval cannot be shorter than 8 hours.
