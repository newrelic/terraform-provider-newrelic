---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_aws_govcloud_integrations"
sidebar_current: "docs-newrelic-cloud-resource-aws-govcloud-integrations"
description: |-
Integrating an AwsGovCloud account to New Relic.
---
-> **IMPORTANT!** This resource is in alpha state, and could still contain issues and missing functionality. If you encounter any issue please create a ticket on [Github](https://github.com/newrelic/terraform-provider-newrelic/issues/new/choose) with all the required information.

# Resource: newrelic_cloud_aws_govcloud_integrations

Use this resource to integrate an AWSGovCloud account to New Relic.

## Prerequisite

Obtain the AwsGovCloud account designed to address the specific regulatory needs of United States (federal, state, and local agencies), education institutions, and the supporting ecosystem.

It is an isolated AWS region designed to host sensitive data and regulated workloads in the cloud, helping customers support their US government compliance requirements.

To pull data from AWSGovCloud, complete the [steps outlined here](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/connect-aws-govcloud-new-relic).

## Example Usage

```hcl
resource "newrelic_cloud_aws_govcloud_link_account" "account" {
  access_key_id          = "%[1]s"
  aws_account_id         = "%[2]s"
  metric_collection_mode = "PULL"
  name                   = "%[4]s"
  secret_access_key      = "%[3]s"
}
resource "newrelic_cloud_aws_govcloud_integrations" "foo" {
  account_id        = 1234321
  linked_account_id = newrelic_cloud_awsGovcloud_link_account.account.id
  alb {
    metrics_polling_interval = 1000
    aws_regions              = ["us-east-1"]
    fetch_extended_inventory = true
    fetch_tags               = true
    load_balancer_prefixes   = [""]
    tag_key                  = ""
    tag_value                = ""
  }
  api_gateway {
    metrics_polling_interval = 1000
    aws_regions              = [""]
    stage_prefixes           = [""]
    tag_key                  = ""
    tag_value                = ""
  }
  auto_scaling {
    metrics_polling_interval = 1000
    aws_regions              = [""]
  }
  aws_direct_connect {
    metrics_polling_interval = 1000
    aws_regions              = [""]
  }
  aws_states {
    metrics_polling_interval = 1000
    aws_regions              = [""]
  }
  cloudtrail {
    metrics_polling_interval = 1000
    aws_regions              = [""]
  }
  dynamo_db {
    metrics_polling_interval = 1000
    aws_regions              = [""]
    fetch_extended_inventory = true
    fetch_tags               = true
    tag_key                  = ""
    tag_value                = ""
  }
  ebs {
    metrics_polling_interval = 1000
    aws_regions              = [""]
    fetch_extended_inventory = true
    tag_key                  = ""
    tag_value                = ""
  }
  ec2 {
    metrics_polling_interval = 1000
    aws_regions              = [""]
    fetch_ip_addresses       = true
    tag_key                  = ""
    tag_value                = ""
  }
  elastic_search {
    metrics_polling_interval = 1000
    aws_regions              = [""]
    fetch_nodes              = true
    tag_key                  = ""
    tag_value                = ""
  }
  elb {
    metrics_polling_interval = 1000
    aws_regions              = [""]
    fetch_extended_inventory = true
    fetch_tags               = true
  }
  emr {
    metrics_polling_interval = 1000
    aws_regions              = [""]
    fetch_tags               = true
    tag_key                  = ""
    tag_value                = ""
  }
  iam {
    metrics_polling_interval = 1000
    tag_key                  = ""
    tag_value                = ""
  }
  lambda {
    metrics_polling_interval = 1000
    aws_regions              = [""]
    fetch_tags               = true
    tag_key                  = ""
    tag_value                = ""
  }
  rds {
    metrics_polling_interval = 1000
    aws_regions              = [""]
    fetch_tags               = true
    tag_key                  = ""
    tag_value                = ""
  }
  red_shift {
    metrics_polling_interval = 1000
    aws_regions              = [""]
    tag_key                  = ""
    tag_value                = ""
  }
  route53 {
    metrics_polling_interval = 1000
    fetch_extended_inventory = true
  }
  s3 {
    metrics_polling_interval = 1000
    fetch_extended_inventory = true
    fetch_tags               = true
    tag_key                  = ""
    tag_value                = ""
  }
  sns {
    metrics_polling_interval = 1000
    aws_regions              = [""]
    fetch_extended_inventory = true
  }
  sqs {
    metrics_polling_interval = 1000
    aws_regions              = [""]
    fetch_extended_inventory = true
    fetch_tags               = true
    queue_prefixes           = [""]
    tag_key                  = ""
    tag_value                = ""
  }
}
```
## Argument Reference

The following arguments are supported:

- `account_id` - (Optional) The New Relic account ID to operate on. This allows the user to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.
- `linked_account_id` - (Required) The access key of the AwsGovCloud.
- `alb` - (Optional) Application load balancer AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `api_gateway` - (Optional) Api Gateway AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `auto_scaling` - (Optional) Autoscaling AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `aws_direct_connect` - (Optional) Aws Direct Connect AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `aws_states` - (Optional) Aws States AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `cloudtrail` - (Optional) Cloudtrail AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `dynamo_db` - (Optional) Dynamo DB AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `ebs` - (Optional) Elastic Beanstalk AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `ec2` - (Optional) EC2 AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `elastic_search` - (Optional) Elastic search AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `elb` - (Optional) Elb AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `emr` - (Optional) Emr AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `iam` - (Optional) IAM AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `lambda` - (Optional) Lambda AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `rds` - (Optional) RDS AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `red_shift` - (Optional) Redshift AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `route53` - (Optional) Route53 AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `sns` - (Optional) SNS AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.
- `sqs` - (Optional) SQS AwsGovCloud integration.See [Integration blocks](#integration-blocks) below for details.

### `Integration` blocks

All `integration` blocks support the following common arguments:

* `metrics_polling_interval` - (Optional) The data polling interval in seconds.

Some integration types support an additional set of arguments:

* `alb`
    * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
    * `fetch_extended_inventory` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.
    * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
    * `load_balancer_prefixes` - (Optional) Specify each name or prefix for the LBs that you want to monitor. Filter values are case-sensitive.
    * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
    * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `api Gateway`
    * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
    * `stage_prefixes` - (Optional) Specify each name or prefix for the Stages that you want to monitor. Filter values are case-sensitive.
    * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
    * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `auto scaling`
    * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
* `direct connect`
    * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
* `aws states`
    * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
* `cloudtrail`
    * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
* `dynamoDB`
    * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
    * `fetch_extended_inventory` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.
    * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
    * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
    * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `ebs`
    * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
    * `fetch_extended_inventory` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.
    * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
    * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `ec2`
    * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
    * `fetch_ip_addresses` - (Optional) Specify if IP addresses of ec2 instance should be collected
    * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
    * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `elastic search`
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
* `s3`
    * `fetch_extended_inventory` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.
    * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
    * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
    * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
* `sns`
    * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
    * `fetch_extended_inventory` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.
* `sqs`
    * `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
    * `fetch_extended_inventory` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.
    * `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
    * `queue_prefixes` - (Optional) Specify each name or prefix for the Queues that you want to monitor. Filter values are case-sensitive.
    * `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
    * `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the AWSGovCloud linked account.

## Import

Integrate AWSGovCloud accounts can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_cloud_aws_govcloud_integrations.foo <id>
```
