# AWS EU Sovereign Cloud Integration Module

This Terraform module creates and configures New Relic AWS EU Sovereign Cloud integrations. It simplifies the setup of monitoring for AWS services running in the EU Sovereign Cloud (aws-eusc partition).

## Prerequisites

Before using this module, ensure you have:

1. An AWS EU Sovereign account with the necessary permissions
2. A New Relic account with cloud integration capabilities
3. An IAM role in AWS EU Sovereign with permissions for New Relic integrations
4. The New Relic Terraform provider configured

## Usage

### Basic Example

```hcl
module "aws_eu_sovereign_integration" {
  source = "./modules/cloud-integrations/aws-eu-sovereign"

  account_name  = "my-eu-sovereign-account"
  aws_role_arn  = "arn:aws-eusc:iam::123456789012:role/NewRelicInfrastructure-Integrations"

  # Enable specific integrations
  cloudtrail_integration = {
    metrics_polling_interval = 300
    aws_regions              = ["eusc-de-east-1"]
    fetch_extended_inventory = true
    fetch_tags               = true
  }

  health_integration = {
    metrics_polling_interval = 300
    fetch_extended_inventory = true
    fetch_tags               = true
  }
}
```

### Advanced Example with All Supported Integrations

```hcl
module "aws_eu_sovereign_integration" {
  source = "./modules/cloud-integrations/aws-eu-sovereign"

  account_name           = "production-eu-sovereign"
  aws_role_arn          = "arn:aws-eusc:iam::123456789012:role/NewRelicInfrastructure-Integrations"
  metric_collection_mode = "PUSH"

  # CloudTrail Integration
  cloudtrail_integration = {
    metrics_polling_interval = 300
    aws_regions              = ["eusc-de-east-1"]
    fetch_extended_inventory = true
    fetch_tags               = true
    tag_key                  = "Environment"
    tag_value                = "Production"
  }

  # AWS Health Integration
  health_integration = {
    metrics_polling_interval = 300
    fetch_extended_inventory = true
    fetch_tags               = true
    tag_key                  = "Application"
    tag_value                = "MyApp"
  }

  # AWS Trusted Advisor Integration
  trusted_advisor_integration = {
    metrics_polling_interval = 300
    fetch_extended_inventory = true
    fetch_tags               = true
    tag_key                  = "Environment"
    tag_value                = "Production"
  }

  # AWS X-Ray Integration
  xray_integration = {
    metrics_polling_interval = 300
    aws_regions              = ["eusc-de-east-1"]
    fetch_extended_inventory = true
    fetch_tags               = true
    tag_key                  = "Application"
    tag_value                = "MyApp"
  }
}
```

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.0 |
| newrelic | ~> 3.0 |

## Providers

| Name | Version |
|------|---------|
| newrelic | ~> 3.0 |

## Inputs

### Required Inputs

| Name | Description | Type |
|------|-------------|------|
| account_name | The name of the AWS EU Sovereign account in New Relic | `string` |
| aws_role_arn | The ARN of the IAM role in AWS EU Sovereign for New Relic integrations | `string` |

### Optional Inputs

| Name | Description | Type | Default |
|------|-------------|------|---------|
| metric_collection_mode | How metrics are collected. Either PULL or PUSH | `string` | `"PUSH"` |
| enable_integrations | Whether to enable AWS integrations | `bool` | `true` |

### Integration Configuration Inputs

Each integration can be configured with the following pattern. Set to `null` to disable an integration:

| Name | Description | Type | Default |
|------|-------------|------|---------|
| cloudtrail_integration | Configuration for CloudTrail integration | `object` | `null` |
| health_integration | Configuration for AWS Health integration | `object` | `null` |
| trusted_advisor_integration | Configuration for AWS Trusted Advisor integration | `object` | `null` |
| xray_integration | Configuration for AWS X-Ray integration | `object` | `null` |

## Outputs

| Name | Description |
|------|-------------|
| linked_account_id | The ID of the linked AWS EU Sovereign account |
| linked_account_name | The name of the linked AWS EU Sovereign account |
| integrations_id | The ID of the AWS EU Sovereign integrations configuration |

## Notes

1. **Polling-Only Services**: This module configures the four AWS services that require polling mode (PULL) because they are not supported by CloudWatch Metric Streams: CloudTrail, Health, Trusted Advisor, and X-Ray. All other AWS services in EU Sovereign use Metric Streams (PUSH) mode.

2. **EU Sovereign Regions**: AWS EU Sovereign Cloud operates in the `eusc-de-east-1` region. Ensure you're using the correct region identifier.

3. **AWS Partition**: This module is specifically for AWS EU Sovereign Cloud (aws-eusc partition). Use `arn:aws-eusc:iam::...` format for IAM role ARNs.

4. **Metric Collection Mode**: While the account-level `metric_collection_mode` is set to PUSH (for metric streams), these specific 4 services automatically use polling regardless of the account setting.

5. **IAM Permissions**: Your AWS EU Sovereign IAM role must have appropriate permissions for each service you want to monitor. Refer to New Relic's documentation for required permissions.

6. **Polling Intervals**: Consider the impact of polling intervals on AWS API rate limits and costs. Lower intervals provide more frequent data but consume more API calls.

7. **Tag-based Filtering**: Use `tag_key` and `tag_value` parameters to limit monitoring to specific resources, which can help reduce costs and improve performance.
