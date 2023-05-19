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
  billing {}
  cloudtrail {
    metrics_polling_interval = 6000
    aws_regions              = ["region-1", "region-2"]
  }
  health {
    metrics_polling_interval = 6000
  }
  trusted_advisor {
    metrics_polling_interval = 6000
  }
  vpc {
    metrics_polling_interval = 6000
    aws_regions              = ["region-1", "region-2"]
    fetch_nat_gateway        = true
    fetch_vpn                = false
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  x_ray {
    metrics_polling_interval = 6000
    aws_regions              = ["region-1", "region-2"]
  }
  s3 {
    metrics_polling_interval = 6000
  }
  doc_db {
    metrics_polling_interval = 6000
  }
  sqs {
    metrics_polling_interval = 6000
    aws_regions              = ["us-east-1"]
    tag_key                  = "test"
    tag_value                = "test"
  }
  ebs {
    metrics_polling_interval = 6000
    aws_regions              = ["us-east-1"]
    tag_key                  = "test"
    tag_value                = "test"
  }
  alb {
    metrics_polling_interval = 6000
    aws_regions              = ["us-east-1"]
  }
  elasticache {
    metrics_polling_interval = 6000
    aws_regions              = ["us-east-1"]
  }
  api_gateway {
    metrics_polling_interval = 6000
    aws_regions              = ["us-east-1"]
    stage_prefixes           = [""]
    tag_key                  = "test"
    tag_value                = "test"
  }
  auto_scaling {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_app_sync {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_athena {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_cognito {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_connect {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_direct_connect {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_fsx {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_glue {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_kinesis_analytics {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_media_convert {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_media_package_vod {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_meta_data {
    metrics_polling_interval = 6000
  }
  aws_mq {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_msk {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_neptune {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_qldb {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_route53resolver {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_states {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_tags_global {
    metrics_polling_interval = 6000
  }
  aws_transit_gateway {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_waf {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  aws_wafv2 {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 6000
  }
  cloudfront {
    fetch_lambdas_at_edge    = true
    fetch_tags               = true
    metrics_polling_interval = 6000
    tag_key                  = "test"
    tag_value                = "test"
  }
}
```
## Argument Reference

The following arguments are supported:

<details>
* `account_id` - (Optional) The New Relic account ID to operate on.  This allows the user to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.
* `linked_account_id` - (Required) The ID of the linked AWS account in New Relic.
* `billing` - (Optional) Billing integration. See [Integration blocks](#integration-blocks) below for details.
* `cloudtrail` - (Optional) Cloudtrail integration. See [Integration blocks](#integration-blocks) below for details.
* `health` - (Optional) Health integration. See [Integration blocks](#integration-blocks) below for details.
* `trusted_advisor` - (Optional) Trusted Advisor integration. See [Integration blocks](#integration-blocks) below for details.
* `vpc` - (Optional) VPC integration. See [Integration blocks](#integration-blocks) below for details.
* `x_ray` - (Optional) X-Ray integration. See [Integration blocks](#integration-blocks) below for details.
* `s3` - (Optional) S3 integration. See [Integration blocks](#integration-blocks) below for details.
* `doc_db` - (Optional) Doc_DB integration. See [Integration blocks](#integration-blocks) below for details.
* `sqs` - (Optional) SQS integration. See [Integration blocks](#integration-blocks) below for details.
* `ebs` - (Optional) EBS integration. See [Integration blocks](#integration-blocks) below for details.
* `alb` - (Optional) ALB integration. See [Integration blocks](#integration-blocks) below for details.
* `elasticache` - (Optional) Elasticache integration. See [Integration blocks](#integration-blocks) below for details.
* `api_gateway` - (Optional) ApiGateway integration. See [Integration blocks](#integration-blocks) below for details.
* `auto_scaling` - (Optional) AutoScaling integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_app_sync` - (Optional) AppSync integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_athena` - (Optional) Athena integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_cognito` - (Optional) Cognito integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_connect` - (Optional) Connect integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_direct_connect` - (Optional) DirectConnect integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_fsx` - (Optional) Fsx integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_glue` - (Optional) Glue integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_kinesis_analytics` - (Optional) Kinesis Analytics integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_media_convert` - (Optional) Media Convert integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_media_package_vod` - (Optional) Media Package vod integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_meta_data` - (Optional) Meta data integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_mq` - (Optional) Mq integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_msk` - (Optional) Msk integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_neptune` - (Optional) Neptune integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_qldb` - (Optional) Qldb integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_route53resolver` - (Optional) Route53resolver integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_states` - (Optional) States integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_tags_global` - (Optional) TagsGlobal integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_transit_gateway` - (Optional) TransitGateway integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_waf` - (Optional) Waf integration. See [Integration blocks](#integration-blocks) below for details.
* `aws_wafv2` - (Optional) Wafv2 integration. See [Integration blocks](#integration-blocks) below for details.
* `cloudfront` - (Optional) Cloudfront integration. See [Integration blocks](#integration-blocks) below for details.
</details>

### `Integration` blocks

All `integration` blocks support the following common arguments:

* `metrics_polling_interval` - (Optional) The data polling interval in seconds.

Some integration types support an additional set of arguments:

<details>
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

Below integration types supports the following common arguments.

* `auto_scaling`,`aws_app_sync`,`aws_athena`,`aws_cognito`,`aws_connect`,`aws_direct_connect`,`aws_fsx`,`aws_glue`,`aws_kinesis_analytics`,`aws_media_convert`,`aws_media_package_vod`,`aws_mq`,`aws_msk`,`aws_neptune`,`aws_qldb`,`aws_route53resolver`,`aws_states`,`aws_transit_gateway`,`aws_waf`,`aws_wafv2`.    
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
