---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_aws_integrations"
sidebar_current: "docs-newrelic-resource-cloud-aws-integrations"
description: |-
    Integrate AWS services with New Relic.
---

# Resource: newrelic\_cloud\_aws\_integrations

Use this resource to integrate AWS services with New Relic.

## Prerequisite

Setup is required for this resource to work properly. This resource assumes you have [linked an AWS account](cloud_aws_link_account.html) to New Relic and configured it to push metrics using CloudWatch Metric Streams.

New Relic doesn't automatically receive metrics from AWS for some services so this resource can be used to configure integrations to those services.

Using a metric stream to New Relic is the preferred way to integrate with AWS. Follow the [steps outlined here](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/aws-integrations-list/aws-metric-stream/#set-up-metric-stream) to set up a metric stream. This resource supports any integration that's not available through AWS metric stream.

## Example Usage

Leave an integration block empty to use its default configuration. You can also use the [full example, including the AWS set up, found in our guides](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/cloud_integrations_guide#aws).

```hcl
resource "newrelic_cloud_aws_link_account" "foo" {
  arn = aws_iam_role.newrelic_aws_role.arn
  metric_collection_mode = "PULL"
  name = "foo"
}

resource "newrelic_cloud_aws_integrations" "bar" {
  linked_account_id = newrelic_cloud_aws_link_account.foo.id
  billing {
    metrics_polling_interval = 3600
  }
  cloudtrail {
    metrics_polling_interval = 300
    aws_regions              = ["us-east-1", "us-east-2"]
  }
  health {
    metrics_polling_interval = 300
  }
  trusted_advisor {
    metrics_polling_interval = 3600
  }
  vpc {
    metrics_polling_interval = 900
    aws_regions              = ["us-east-1", "us-east-2"]
    fetch_nat_gateway        = true
    fetch_vpn                = false
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  x_ray {
    metrics_polling_interval = 300
    aws_regions              = ["us-east-1", "us-east-2"]
  }
  s3 {
    metrics_polling_interval = 3600
  }
  doc_db {
    metrics_polling_interval = 300
  }
  sqs {
    fetch_extended_inventory = true
    fetch_tags               = true
    queue_prefixes           = ["queue prefix"]
    metrics_polling_interval = 300
    aws_regions              = ["us-east-1"]
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  ebs {
    metrics_polling_interval = 900
    fetch_extended_inventory = true
    aws_regions              = ["us-east-1"]
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  alb {
    fetch_extended_inventory = true
    fetch_tags               = true
    load_balancer_prefixes   = ["load balancer prefix"]
    metrics_polling_interval = 300
    aws_regions              = ["us-east-1"]
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  elasticache {
    aws_regions              = ["us-east-1"]
    fetch_tags               = true
    metrics_polling_interval = 300
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  api_gateway {
    metrics_polling_interval = 300
    aws_regions              = ["us-east-1"]
    stage_prefixes           = ["stage prefix"]
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  auto_scaling {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_app_sync {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_athena {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_cognito {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_connect {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_direct_connect {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_fsx {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_glue {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_kinesis_analytics {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_media_convert {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_media_package_vod {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_mq {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_msk {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_neptune {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_qldb {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_route53resolver {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_states {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_transit_gateway {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_waf {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  aws_wafv2 {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  cloudfront {
    fetch_lambdas_at_edge    = true
    fetch_tags               = true
    metrics_polling_interval = 300
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  dynamodb {
    aws_regions              = ["us-east-1"]
    fetch_extended_inventory = true
    fetch_tags               = true
    metrics_polling_interval = 300
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  ec2 {
    aws_regions              = ["us-east-1"]
    duplicate_ec2_tags       = true
    fetch_ip_addresses       = true
    metrics_polling_interval = 300
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  ecs {
    aws_regions              = ["us-east-1"]
    fetch_tags               = true
    metrics_polling_interval = 300
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  efs {
    aws_regions              = ["us-east-1"]
    fetch_tags               = true
    metrics_polling_interval = 300
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  elasticbeanstalk {
    aws_regions              = ["us-east-1"]
    fetch_extended_inventory = true
    fetch_tags               = true
    metrics_polling_interval = 300
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  elasticsearch {
    aws_regions              = ["us-east-1"]
    fetch_nodes              = true
    metrics_polling_interval = 300
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  elb {
    aws_regions              = ["us-east-1"]
    fetch_extended_inventory = true
    fetch_tags               = true
    metrics_polling_interval = 300
  }
  emr {
    aws_regions              = ["us-east-1"]
    fetch_tags               = true
    metrics_polling_interval = 300
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  iam {
    metrics_polling_interval = 300
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  iot {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  kinesis {
    aws_regions              = ["us-east-1"]
    fetch_shards             = true
    fetch_tags               = true
    metrics_polling_interval = 900
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  kinesis_firehose {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  lambda {
    aws_regions              = ["us-east-1"]
    fetch_tags               = true
    metrics_polling_interval = 300
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  rds {
    aws_regions              = ["us-east-1"]
    fetch_tags               = true
    metrics_polling_interval = 300
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  redshift {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  route53 {
    fetch_extended_inventory = true
    metrics_polling_interval = 300
  }
  ses {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 300
  }
  sns {
    aws_regions              = ["us-east-1"]
    fetch_extended_inventory = true
    metrics_polling_interval = 300
  }
}
```
## Argument Reference

-> **WARNING:** Starting with [v3.27.2](https://registry.terraform.io/providers/newrelic/newrelic/3.27.2) of the New Relic Terraform Provider, updating the `linked_account_id` of a `newrelic_cloud_aws_integrations` resource that has been applied would **force a replacement** of the resource (destruction of the resource, followed by the creation of a new resource). Please carefully review the output of `terraform plan`, which would clearly indicate a replacement of this resource, before performing a `terraform apply`.

* `account_id` - (Optional) The New Relic account ID to operate on.  This allows the user to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.
* `linked_account_id` - (Required) The ID of the linked AWS account in New Relic.


The following arguments are supported with minimum metrics polling interval of 300 seconds

* `alb` - (Optional) AWS ALB. See [Integration blocks](#integration-blocks) below for details.
* `api_gateway` - (Optional) AWS API Gateway. See [Integration blocks](#integration-blocks) below for details.
* `auto_scaling` - (Optional) AWS Auto Scaling. See [Integration blocks](#integration-blocks) below for details.
* `aws_app_sync` - (Optional) AWS AppSync. See [Integration blocks](#integration-blocks) below for details.
* `aws_athena` - (Optional) AWS Athena. See [Integration blocks](#integration-blocks) below for details.
* `aws_cognito` - (Optional) AWS Cognito. See [Integration blocks](#integration-blocks) below for details.
* `aws_connect` - (Optional) AWS Connect. See [Integration blocks](#integration-blocks) below for details.
* `aws_direct_connect` - (Optional) AWS Direct Connect. See [Integration blocks](#integration-blocks) below for details.
* `aws_fsx` - (Optional) AWS FSx. See [Integration blocks](#integration-blocks) below for details.
* `aws_glue` - (Optional) AWS Glue. See [Integration blocks](#integration-blocks) below for details.
* `aws_kinesis_analytics` - (Optional) AWS Kinesis Data Analytics. See [Integration blocks](#integration-blocks) below for details.
* `aws_media_convert` - (Optional) AWS Media Convert. See [Integration blocks](#integration-blocks) below for details.
* `aws_media_package_vod` - (Optional) AWS MediaPackage VOD. See [Integration blocks](#integration-blocks) below for details.
* `aws_mq` - (Optional) AWS MQ. See [Integration blocks](#integration-blocks) below for details.
* `aws_msk` - (Optional) Amazon Managed Kafka (MSK). See [Integration blocks](#integration-blocks) below for details.
* `aws_neptune` - (Optional) AWS Neptune. See [Integration blocks](#integration-blocks) below for details.
* `aws_route53resolver` - (Optional) AWS Route53 Resolver. See [Integration blocks](#integration-blocks) below for details.
* `aws_qldb` - (Optional) Amazon QLDB. See [Integration blocks](#integration-blocks) below for details.
* `aws_transit_gateway` - (Optional) Amazon Transit Gateway. See [Integration blocks](#integration-blocks) below for details.
* `aws_waf` - (Optional) AWS WAF. See [Integration blocks](#integration-blocks) below for details.
* `aws_wafv2` - (Optional) AWS WAF V2. See [Integration blocks](#integration-blocks) below for details.
* `cloudfront` - (Optional) AWS CloudFront. See [Integration blocks](#integration-blocks) below for details.
* `cloudtrail` - (Optional) AWS CloudTrail. See [Integration blocks](#integration-blocks) below for details.
* `doc_db` - (Optional) AWS DocumentDB. See [Integration blocks](#integration-blocks) below for details.
* `dynamodb` - (Optional) Amazon DynamoDB. See [Integration blocks](#integration-blocks) below for details.
* `ec2` - (Optional) Amazon EC2. See [Integration blocks](#integration-blocks) below for details.
* `ecs` - (Optional) Amazon ECS. See [Integration blocks](#integration-blocks) below for details.
* `efs` - (Optional) Amazon EFS. See [Integration blocks](#integration-blocks) below for details.
* `elasticache` - (Optional) AWS ElastiCache. See [Integration blocks](#integration-blocks) below for details.
* `elasticbeanstalk` - (Optional) AWS Elastic Beanstalk. See [Integration blocks](#integration-blocks) below for details.
* `elasticsearch` - (Optional) AWS ElasticSearch. See [Integration blocks](#integration-blocks) below for details.
* `elb` - (Optional) AWS ELB (Classic). See [Integration blocks](#integration-blocks) below for details.
* `emr` - (Optional) AWS EMR. See [Integration blocks](#integration-blocks) below for details.
* `health` - (Optional) AWS Health. See [Integration blocks](#integration-blocks) below for details.
* `iam` - (Optional) AWS IAM. See [Integration blocks](#integration-blocks) below for details.
* `iot` - (Optional) AWS IoT. See [Integration blocks](#integration-blocks) below for details.
* `kinesis_firehose` - (Optional) Amazon Kinesis Data Firehose. See [Integration blocks](#integration-blocks) below for details.
* `lambda` - (Optional) AWS Lambda. See [Integration blocks](#integration-blocks) below for details.
* `rds` - (Optional) Amazon RDS. See [Integration blocks](#integration-blocks) below for details.
* `redshift` - (Optional) Amazon Redshift. See [Integration blocks](#integration-blocks) below for details.
* `route53` - (Optional) Amazon Route 53. See [Integration blocks](#integration-blocks) below for details.
* `s3` - (Optional) Amazon S3. See [Integration blocks](#integration-blocks) below for details.
* `ses` - (Optional) Amazon SES. See [Integration blocks](#integration-blocks) below for details.
* `sns` - (Optional) AWS SNS. See [Integration blocks](#integration-blocks) below for details.
* `sqs` - (Optional) AWS SQS. See [Integration blocks](#integration-blocks) below for details.
* `x_ray` - (Optional) AWS X-Ray. See [Integration blocks](#integration-blocks) below for details.
x

The following arguments are supported with minimum metrics polling interval of 900 seconds

* `ebs` - (Optional) Amazon EBS. See [Integration blocks](#integration-blocks) below for details.
* `kinesis` - (Optional) AWS Kinesis. See [Integration blocks](#integration-blocks) below for details.

The following arguments are supported with minimum metrics polling interval of 3600 seconds

* `billing` - (Optional) AWS Billing. See [Integration blocks](#integration-blocks) below for details.
* `trusted_advisor` - (Optional) AWS Trusted Advisor. See [Integration blocks](#integration-blocks) below for details.

All other arguments are dependent on the services to be integrated, which have been listed in the collapsible section below. All of these are **optional** blocks that can be added in any required combination. **For details on arguments that can be used with each service, check out the [`Integration` blocks](#integration-blocks) section below.**
<details>
  <summary>Expand this section to view all supported AWS services supported, that may be integrated via this resource.</summary>

| Block                   | Description                   |
|-------------------------|-------------------------------|
| `alb`                   | ALB Integration               |
| `api_gateway`           | API Gateway Integration       |
| `auto_scaling`          | Auto Scaling Integration      |
| `aws_app_sync`          | AppSync Integration           |
| `aws_athena`            | Athena Integration            |
| `aws_cognito`           | Cognito Integration           |
| `aws_connect`           | Connect Integration           |
| `aws_direct_connect`    | Direct Connect Integration    |
| `aws_fsx`               | FSx Integration               |
| `aws_glue`              | Glue Integration              |
| `aws_kinesis_analytics` | Kinesis Analytics Integration |
| `aws_media_convert`     | MediaConvert Integration      |
| `aws_media_package_vod` | Media Package VOD Integration |
| `aws_mq`                | MQ Integration                |
| `aws_msk`               | MSK Integration               |
| `aws_neptune`           | Neptune Integration           |
| `aws_qldb`              | QLDB Integration              |
| `aws_route53resolver`   | Route53 Resolver Integration  |
| `aws_states`            | States Integration            |
| `aws_transit_gateway`   | Transit Gateway Integration   |
| `aws_waf`               | WAF Integration               |
| `aws_wafv2`             | WAFv2 Integration             |
| `billing`               | Billing Integration           |
| `cloudfront`            | CloudFront Integration        |
| `cloudtrail`            | CloudTrail Integration        |
| `doc_db`                | DocumentDB Integration        |
| `dynamodb`              | DynamoDB Integration          |
| `ebs`                   | EBS Integration               |
| `ec2`                   | EC2 Integration               |
| `ecs`                   | ECS Integration               |
| `efs`                   | EFS Integration               |
| `elasticache`           | ElastiCache Integration       |
| `elasticbeanstalk`      | Elastic Beanstalk Integration |
| `elasticsearch`         | Elasticsearch Integration     |
| `elb`                   | ELB Integration               |
| `emr`                   | EMR Integration               |
| `health`                | Health Integration            |
| `iam`                   | IAM Integration               |
| `iot`                   | IoT Integration               |
| `kinesis`               | Kinesis Integration           |
| `kinesis_firehose`      | Kinesis Firehose Integration  |
| `lambda`                | Lambda Integration            |
| `rds`                   | RDS Integration               |
| `redshift`              | Redshift Integration          |
| `route53`               | Route53 Integration           |
| `s3`                    | S3 Integration                |
| `ses`                   | SES Integration               |
| `sns`                   | SNS Integration               |
| `sqs`                   | SQS Integration               |
| `trusted_advisor`       | Trusted Advisor Integration   |
| `vpc`                   | VPC Integration               |
| `x_ray`                 | X-Ray Integration             |

</details>

### `Integration` blocks

All `integration` blocks support the following common arguments:

* `metrics_polling_interval` - (Optional) The data polling interval **in seconds**.

-> **NOTE** For more information on the ranges of metric polling intervals of each of these integrations, head over to [this page](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/get-started/introduction-aws-integrations/)


<details>
  <summary> Some integration types support an additional set of arguments. Expand this section to take a look at these supported arguments. </summary>
* `cloudtrail`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
* `vpc`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `fetch_nat_gateway` - (Optional) Specify if NAT gateway should be monitored. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `fetch_vpn` - (Optional) Specify if VPN should be monitored. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `x_ray`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
* `s3`
  * `fetch_extended_inventory` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `doc_db`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
* `sqs`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `fetch_extended_inventory` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `queue_prefixes` - (Optional) Specify each name or prefix for the Queues that you want to monitor. Filter values are case-sensitive.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `ebs`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `fetch_extended_inventory` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `alb`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `fetch_extended_inventory` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `load_balancer_prefixes` - (Optional) Specify each name or prefix for the LBs that you want to monitor. Filter values are case-sensitive.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `elasticache`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `api_gateway`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `stage_prefixes` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.
* `cloudfront`
  * `fetch_lambdas_at_edge` - (Optional) Specify if Lambdas@Edge should be monitored. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `dynamodb`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `fetch_extended_inventory` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `ec2`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `duplicate_ec2_tags` - (Optional) Specify if the old legacy metadata and tag names have to be kept, it will consume more ingest data size.
  * `fetch_ip_addresses` - (Optional) Specify if IP addresses of ec2 instance should be collected.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `ecs`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `efs`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `elasticbeanstalk`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `fetch_extended_inventory` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `elasticsearch`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `fetch_nodes` - (Optional) Specify if metrics should be collected for nodes. Turning it on will increase the number of API calls made to CloudWatch.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `elb`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `fetch_extended_inventory` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
* `emr`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `iam`
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `kinesis`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `fetch_shards` - (Optional) Specify if Shards should be monitored. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `lambda`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `rds`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `redshift`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `route53`
  * `fetch_extended_inventory` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.
* `sns`
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * `fetch_extended_inventory` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.

Furthermore, below integration types supports the following common arguments.

* `auto_scaling`,`aws_app_sync`,`aws_athena`,`aws_cognito`,`aws_connect`,`aws_direct_connect`,`aws_fsx`,`aws_glue`,`aws_kinesis_analytics`,`aws_media_convert`,`aws_media_package_vod`,`aws_mq`,`aws_msk`,`aws_neptune`,`aws_qldb`,`aws_route53resolver`,`aws_states`,`aws_transit_gateway`,`aws_waf`,`aws_wafv2`,`iot`,`kinesis_firehose` and `ses`.    
  * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.  
</details>


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the AWS linked account.

## Import

Linked AWS account integrations can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_cloud_aws_integrations.foo <id>
```