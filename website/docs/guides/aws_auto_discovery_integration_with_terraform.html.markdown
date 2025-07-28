---
layout: "newrelic"
page_title: "New Relic Integration with AWS AutoDiscovery"
sidebar_current: "docs-newrelic-provider-aws-auto-discovery-integration-guide"
description: |-
  Use this guide to set up the New Relic AWS Auto Discovery fully automated through Terraform.
---

# NewRelic Integration with AWS AutoDiscovery
The below guide provides reference on how to integrate with AWS Auto-Discovery via terraform.

## Documentation Reference
* [AWS Auto Discovery](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/set-up-auto-discovery-of-aws-entities/)

## Pre-Requisites
The below integration process presumes that the client has already integrated with either of the New Relic's integration for monitoring their cloud-infrastructure on AWS, which is either of API Polling or Cloudwatch Metrics Streams integration.

## Variables -
* **New Relic Region** ```new_relic_region``` - The region where your New Relic account is hosted
 *  **Type** - String
 *  **Allowed Values** - US or EU

* **New Relic Account Id** ```new_relic_account_id``` - The account id associated with your organisation's New Relic Account.
 * **Type** - Integer
 * **Example** - 1234567

* **AWS IAM Role Name** ```aws_iam_role_name``` - The IAM role-name that has been used to integrate cloud monitoring with New Relic.
 * **Type** - String
 * **Example** - "NewRelicInfrastructureIntegrationRole"

* **AWS Account Id** ```aws_account_id``` - The AWS Account Id which has been integrated with New Relic on which Auto-Discovery has to be enabled.
 * **Type** - Integer
 * **Example** - 123456789012

* **Metric Polling Interval** ```metric_polling_interval``` - The interval at which the auto-discovery scans are to be triggered on the integrated AWS account in seconds
 * **Type** - Integer
 * **Allowed Values** - 28800, 43200, 86400

* **aws_regions** ```aws_regions``` - The list of regions at which the auto-discovery scans are to be triggered on the integrated AWS account.
 * **Type** - StringList
 * **Allowed Values** - ["af-south-1", "ap-east-1", "ap-northeast-1", "ap-northeast-2", "ap-northeast-3", "ap-south-1", "ap-southeast-1", "ap-southeast-2", "ca-central-1", "eu-north-1", "eu-south-1", "eu-west-1", "eu-west-2", "eu-west-3", "me-south-1", "sa-east-1",
  "us-east-1" "us-east-2", "us-west-1", "us-west-2"]

## Integrating with AWS Auto-Discovery via Terraform
```
provider "newrelic" {
  region = var.new_relic_region
  account_id = var.newrelic_account_id
}

resource "aws_iam_policy" "newrelic_aws_autodiscovery_permissions" {
  name        = "NewRelicAwsAutoDiscoveryPermissions"
  description = "IAM Policy for Auto-Discovery Permissions"
  policy      = << EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "cloudformation:StartResourceScan"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "newrelic_aws_policy_attach" {
  role       = "var.aws_iam_role_name"
  policy_arn = aws_iam_policy.newrelic_aws_autodiscovery_permissions.arn
}

resource "newrelic_cloud_aws_integrations" "auto_discovery_integrations" {
  linked_account_id = var.aws_account_id
  aws_auto_discovery = {
    metric_polling_interval = var.metric_polling_interval
    aws_regions = var.aws_regions
  }
}
```


