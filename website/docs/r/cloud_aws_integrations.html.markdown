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

The following example demonstrates the use of the `newrelic_cloud_aws_integrations` resource with multiple AWS integrations supported by the resource.

To view a full example with all supported AWS integrations, please see the [Additional Examples](#additional-examples) section. Integration blocks used in the resource may also be left empty to use the default configuration of the integration.

A full example, inclusive of setup of AWS resources (from the AWS Terraform Provider) associated with this resource, may be found in our [AWS cloud integration guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/cloud_integrations_guide#aws).

```hcl
resource "newrelic_cloud_aws_link_account" "foo" {
  arn = aws_iam_role.newrelic_aws_role.arn
  metric_collection_mode = "PULL"
  name = "foo"
}

resource "newrelic_cloud_aws_integrations" "bar" {
  linked_account_id = newrelic_cloud_aws_link_account.foo.id
  cloudtrail {
    metrics_polling_interval = 300
    aws_regions              = ["us-east-1", "us-east-2"]
  }
  vpc {
    metrics_polling_interval = 900
    aws_regions              = ["us-east-1", "us-east-2"]
    fetch_nat_gateway        = true
    fetch_vpn                = false
    tag_key                  = "tag key"
    tag_value                = "tag value"
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
  api_gateway {
    metrics_polling_interval = 300
    aws_regions              = ["us-east-1"]
    stage_prefixes           = ["stage prefix"]
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  cloudfront {
    fetch_lambdas_at_edge    = true
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
  elasticsearch {
    aws_regions              = ["us-east-1"]
    fetch_nodes              = true
    metrics_polling_interval = 300
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
  kinesis {
    aws_regions              = ["us-east-1"]
    fetch_shards             = true
    fetch_tags               = true
    metrics_polling_interval = 900
    tag_key                  = "tag key"
    tag_value                = "tag value"
  }
}
```

## Supported AWS Integrations
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
| `security_hub`          | Security Hub Integration      |
| `sns`                   | SNS Integration               |
| `sqs`                   | SQS Integration               |
| `trusted_advisor`       | Trusted Advisor Integration   |
| `vpc`                   | VPC Integration               |
| `x_ray`                 | X-Ray Integration             |

</details>

## Argument Reference

-> **WARNING:** Starting with [v3.27.2](https://registry.terraform.io/providers/newrelic/newrelic/3.27.2) of the New Relic Terraform Provider, updating the `linked_account_id` of a `newrelic_cloud_aws_integrations` resource that has been applied would **force a replacement** of the resource (destruction of the resource, followed by the creation of a new resource). Please carefully review the output of `terraform plan`, which would clearly indicate a replacement of this resource, before performing a `terraform apply`.

* `account_id` - (Optional) The New Relic account ID to operate on.  This allows the user to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.
* `linked_account_id` - (Required) The ID of the linked AWS account in New Relic.


### Arguments to be Specified with Integration Blocks 

The following arguments are intended to be used within certain ["integration blocks"](#integration-blocks) in the resource, i.e. these are supposed to be specified as nested arguments "within" an argument corresponding to a specific AWS integration, unlike `account_id` and `linked_account_id` which are specified at the resource level. An exhaustive list of all of such arguments supported by the resource (and the integration blocks they would need to be specified with) has been given below.

In order to find the right set of arguments which go with each integration, and samples on the usage of these arguments "within" integration blocks, **check out the [Integration Blocks](#integration-blocks) section below.**

* `metrics_polling_interval` - (Optional) The data polling interval **in seconds**.
  * The following integration blocks support the usage of this argument: 

        |                        |                      |                        |
        |------------------------|----------------------|------------------------|
        | `billing`              | `cloudtrail`         | `health`               |
        | `trusted_advisor`      | `vpc`                | `x_ray`                |
        | `s3`                   | `doc_db`             | `sqs`                  |
        | `ebs`                  | `alb`                | `elasticache`          |
        | `api_gateway`          | `auto_scaling`       | `aws_app_sync`         |
        | `aws_athena`           | `aws_cognito`        | `aws_connect`          |
        | `aws_direct_connect`   | `aws_fsx`            | `aws_glue`             |
        | `aws_kinesis_analytics`| `aws_media_convert`  | `aws_media_package_vod`|
        | `aws_mq`               | `aws_msk`            | `aws_neptune`          |
        | `aws_qldb`             | `aws_route53resolver`| `aws_states`           |
        | `aws_transit_gateway`  | `aws_waf`            | `aws_wafv2`            |
        | `cloudfront`           | `dynamodb`           | `ec2`                  |
        | `ecs`                  | `efs`                | `elasticbeanstalk`     |
        | `elasticsearch`        | `elb`                | `emr`                  |
        | `iam`                  | `iot`                | `security_hub`         |


-> **NOTE** For more information on the ranges of metric polling intervals of each of these integrations, head over to [this page](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/get-started/introduction-aws-integrations/). You may also find the range of metric polling intervals of an integration under the [Integration Blocks](#integration-blocks) section.

* `aws_regions` - (Optional) Specify each AWS region that includes the resources that you want to monitor.
  * The following integration blocks support the usage of this argument: 

        |                        |                      |                        |
        |------------------------|----------------------|------------------------|
        | `cloudtrail`           | `vpc`                | `x_ray`                |
        | `ebs`                  | `alb`                | `elasticache`          |
        | `api_gateway`          | `auto_scaling`       | `aws_app_sync`         |
        | `aws_athena`           | `aws_cognito`        | `aws_connect`          |
        | `aws_direct_connect`   | `aws_fsx`            | `aws_glue`             |
        | `aws_kinesis_analytics`| `aws_media_convert`  | `aws_media_package_vod`|
        | `aws_mq`               | `aws_msk`            | `aws_neptune`          |
        | `aws_qldb`             | `aws_route53resolver`| `aws_states`           |
        | `aws_transit_gateway`  | `aws_waf`            | `aws_wafv2`            |
        | `cloudfront`           | `dynamodb`           | `ec2`                  |
        | `ecs`                  | `efs`                | `elasticbeanstalk`     |
        | `elasticsearch`        | `elb`                | `emr`                  |
        | `lambda`               | `iot`                | `kinesis`              |
        | `rds`                  | `ses`                | `redshift`             |
        | `sns`                  | `sqs`                | `security_hub`         |   


* `fetch_nat_gateway` - (Optional) Specify if NAT gateway should be monitored. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * The following integration blocks support the usage of this argument: `vpc`
* `fetch_vpn` - (Optional) Specify if VPN should be monitored. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * The following integration blocks support the usage of this argument: `vpc`
* `tag_key` - (Optional) Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.
  * The following integration blocks support the usage of this argument: `vpc`, `sqs`, `ebs`, `alb`, `elasticache`, `api_gateway`, `cloudfront`, `dynamodb`, `ec2`, `ecs`, `efs`, `elasticbeanstalk`, `elasticsearch`, `emr`, `iam`, `kinesis`, `lambda`, `rds`, `redshift`
* `tag_value` - (Optional) Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.
  * The following integration blocks support the usage of this argument: `vpc`, `sqs`, `ebs`, `alb`, `elasticache`, `api_gateway`, `cloudfront`, `dynamodb`, `ec2`, `ecs`, `efs`, `elasticbeanstalk`, `elasticsearch`, `emr`, `iam`, `kinesis`, `lambda`, `rds`, `redshift`
* `fetch_extended_inventory` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * The following integration blocks support the usage of this argument: `sqs`, `ebs`, `alb`, `cloudfront`, `dynamodb`, `elasticbeanstalk`, `elb`, `route53`, `sns`
* `fetch_tags` - (Optional) Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * The following integration blocks support the usage of this argument: `sqs`, `alb`, `elasticache`, `cloudfront`, `dynamodb`, `ecs`, `efs`, `elasticbeanstalk`, `elasticsearch`, `emr`, `kinesis`, `lambda`, `rds`
* `queue_prefixes` - (Optional) Specify each name or prefix for the Queues that you want to monitor. Filter values are case-sensitive.
  * The following integration blocks support the usage of this argument: `sqs`
* `load_balancer_prefixes` - (Optional) Specify each name or prefix for the LBs that you want to monitor. Filter values are case-sensitive.
  * The following integration blocks support the usage of this argument: `alb`
* `stage_prefixes` - (Optional) Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * The following integration blocks support the usage of this argument: `api_gateway`
* `fetch_lambdas_at_edge` - (Optional) Specify if Lambdas@Edge should be monitored. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * The following integration blocks support the usage of this argument: `cloudfront`
* `duplicate_ec2_tags` - (Optional) Specify if the old legacy metadata and tag names have to be kept, it will consume more ingest data size.
  * The following integration blocks support the usage of this argument: `ec2`
* `fetch_ip_addresses` - (Optional) Specify if IP addresses of ec2 instance should be collected.
  * The following integration blocks support the usage of this argument: `ec2`
* `fetch_shards` - (Optional) Specify if Shards should be monitored. May affect total data collection time and contribute to the Cloud provider API rate limit.
  * The following integration blocks support the usage of this argument: `kinesis`


## Integration Blocks

The following section lists out arguments which may be used with each AWS integration supported by this resource. 

As specified above in the [Arguments to be Specified with Integration Blocks](#arguments-to-be-specified-with-integration-blocks) section, except for `linked_account_id` and `account_id`, all aforementioned arguments are to be specified within an integration block as they are supported by a specific set of integrations each; the following list of integration blocks elucidates the same with samples of what each integration block would look like.

<details>
  <summary> Expand this list to see a list of all integration blocks supported by this resource, the arguments which go with them and a sample of what the block would look like with these arguments. </summary>
  <details>
    <summary>cloudtrail</summary>
*  Supported Arguments: `aws_regions` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl 
     cloudtrail {
        metrics_polling_interval = 300
        aws_regions              = ["us-east-1", "us-east-2"]
     }
``` 
  </details>
  <details>
    <summary>vpc</summary>
*  Supported Arguments: `aws_regions` `fetch_nat_gateway` `fetch_vpn` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl 
     vpc {
      metrics_polling_interval = 900
      aws_regions              = ["us-east-1", "us-east-2"]
      fetch_nat_gateway        = true
      fetch_vpn                = false
      tag_key                  = "tag key"
      tag_value                = "tag value"
    }
```
  </details>
  <details>
    <summary>x_ray</summary>
*  Supported Arguments: `aws_regions` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 60,300, 900, 1800, 3600 (seconds)
```hcl 
     x_ray {
      metrics_polling_interval = 300
      aws_regions              = ["us-east-1", "us-east-2"]
    }
```
  </details>
  <details>
    <summary>s3</summary>
*  Supported Arguments: `fetch_extended_inventory` `fetch_tags` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl 
     s3 { 
        metrics_polling_interval = 3600
        fetch_extended_inventory = true
        fetch_tags               = true
        tag_key                  = "tag key"
        tag_value                = "tag value"
     }
```
</details>
  <details>
    <summary>doc_db</summary>
*  Supported Arguments: `aws_regions` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl 
       doc_db {
          metrics_polling_interval = 300
          aws_regions              = ["us-east-1", "us-east-2"]
       }
```
  </details>
  <details>
    <summary>sqs</summary>
*  Supported Arguments: `aws_regions` `fetch_extended_inventory` `fetch_tags` `queue_prefixes` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl 
       sqs {
          fetch_extended_inventory = true
          fetch_tags               = true
          queue_prefixes           = ["queue prefix"]
          metrics_polling_interval = 300
          aws_regions              = ["us-east-1"]
          tag_key                  = "tag key"
          tag_value                = "tag value"
        }
```
  </details>
  <details>
    <summary>ebs</summary>
*  Supported Arguments: `aws_regions` `fetch_extended_inventory` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 900, 1800, 3600 (seconds)
```hcl  
       ebs {
        metrics_polling_interval = 900
        fetch_extended_inventory = true
        aws_regions              = ["us-east-1"]
        tag_key                  = "tag key"
        tag_value                = "tag value"
      }
```
  </details>
  <details>
    <summary>alb</summary>
*  Supported Arguments: `aws_regions` `fetch_extended_inventory` `fetch_tags` `load_balancer_prefixes` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl 
        alb {
          fetch_extended_inventory = true
          fetch_tags               = true
          load_balancer_prefixes   = ["load balancer prefix"]
          metrics_polling_interval = 300
          aws_regions              = ["us-east-1"]
          tag_key                  = "tag key"
          tag_value                = "tag value"
        }
``` 
  </details>
  <details>
    <summary>elasticache</summary>
*  Supported Arguments: `aws_regions` `fetch_tags` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl 
       elasticache {
        aws_regions              = ["us-east-1"]
        fetch_tags               = true
        metrics_polling_interval = 300
        tag_key                  = "tag key"
        tag_value                = "tag value"
      }
```
  </details>
  <details>
    <summary>api_gateway</summary>
*  Supported Arguments: `aws_regions` `tag_key` `tag_value` `stage_prefixes`
```hcl 
       api_gateway {
        metrics_polling_interval = 300
        aws_regions              = ["us-east-1"]
        stage_prefixes           = ["stage prefix"]
        tag_key                  = "tag key"
        tag_value                = "tag value"
      }
```
  </details>
  <details>
    <summary>cloudfront</summary>
*  Supported Arguments: `fetch_lambdas_at_edge` `fetch_tags` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl 
       cloudfront {
        fetch_lambdas_at_edge    = true
        fetch_tags               = true
        metrics_polling_interval = 300
        tag_key                  = "tag key"
        tag_value                = "tag value"
      }
```
  </details>
  <details>
    <summary>dynamodb</summary>
*  Supported Arguments: `aws_regions` `fetch_extended_inventory` `fetch_tags` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl 
        dynamodb {
          aws_regions              = ["us-east-1"]
          fetch_extended_inventory = true
          fetch_tags               = true
          metrics_polling_interval = 300
          tag_key                  = "tag key"
          tag_value                = "tag value"
        }
``` 
  </details>
  <details>
    <summary>ec2</summary>
*  Supported Arguments: `aws_regions` `duplicate_ec2_tags` `fetch_ip_addresses` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl 
       ec2 {
        aws_regions              = ["us-east-1"]
        duplicate_ec2_tags       = true
        fetch_ip_addresses       = true
        metrics_polling_interval = 300
        tag_key                  = "tag key"
        tag_value                = "tag value"
      }
```
  </details>
  <details>
    <summary>ecs</summary>
*  Supported Arguments: `aws_regions` `fetch_tags` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl
        ecs {
          aws_regions              = ["us-east-1"]
          fetch_tags               = true
          metrics_polling_interval = 300
          tag_key                  = "tag key"
          tag_value                = "tag value"
        }
``` 
  </details>
  <details>
    <summary>efs</summary>
*  Supported Arguments: `aws_regions` `fetch_tags` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl 
        efs {
          aws_regions              = ["us-east-1"]
          fetch_tags               = true
          metrics_polling_interval = 300
          tag_key                  = "tag key"
          tag_value                = "tag value"
        }
``` 
  </details>
  <details>
    <summary>elasticbeanstalk</summary>
*  Supported Arguments: `aws_regions` `fetch_extended_inventory` `fetch_tags` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl 
       elasticbeanstalk {
        aws_regions              = ["us-east-1"]
        fetch_extended_inventory = true
        fetch_tags               = true
        metrics_polling_interval = 300
        tag_key                  = "tag key"
        tag_value                = "tag value"
      }
``` 
  </details>
  <details>
    <summary>elasticsearch</summary>
*  Supported Arguments: `aws_regions` `fetch_nodes` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl 
       elasticsearch {
        aws_regions              = ["us-east-1"]
        fetch_nodes              = true
        metrics_polling_interval = 300
        tag_key                  = "tag key"
        tag_value                = "tag value"
      }
``` 
  </details>
  <details>
    <summary>elb</summary>
*  Supported Arguments: `aws_regions` `fetch_extended_inventory` `fetch_tags` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl 
        elb {
          aws_regions              = ["us-east-1"]
          fetch_extended_inventory = true
          fetch_tags               = true
          metrics_polling_interval = 300
        }
 ```
  </details>
  <details>
    <summary>emr</summary>
*  Supported Arguments: `aws_regions` `fetch_tags` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl 
        emr {
          aws_regions              = ["us-east-1"]
          fetch_tags               = true
          metrics_polling_interval = 300
          tag_key                  = "tag key"
          tag_value                = "tag value"
        }
``` 
  </details>
  <details>
    <summary>iam</summary>
*  Supported Arguments: `tag_key` `tag_value` `metrics_polling_interval`
```hcl  
        iam {
          metrics_polling_interval = 300
          tag_key                  = "tag key"
          tag_value                = "tag value"
        }
``` 
  </details>
  <details>
    <summary> kinesis </summary>
*  Supported Arguments: `aws_regions` `fetch_shards` `fetch_tags` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 900, 1800, 3600 (seconds)
```hcl
        kinesis {
          aws_regions              = ["us-east-1"]
          fetch_shards             = true
          fetch_tags               = true
          metrics_polling_interval = 900
          tag_key                  = "tag key"
          tag_value                = "tag value"
        }
 ``` 
  </details>
  <details>
    <summary>lambda</summary>
*  Supported Arguments: `aws_regions` `fetch_tags` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl  
        lambda {
          aws_regions              = ["us-east-1"]
          fetch_tags               = true
          metrics_polling_interval = 300
          tag_key                  = "tag key"
          tag_value                = "tag value"
        }
 ``` 
  </details>
  <details>
    <summary>rds</summary>
*  Supported Arguments: `aws_regions` `fetch_tags` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl  
        rds {
          aws_regions              = ["us-east-1"]
          fetch_tags               = true
          metrics_polling_interval = 300
          tag_key                  = "tag key"
          tag_value                = "tag value"
        }
``` 
  </details>
  <details>
    <summary>redshift</summary>
*  Supported Arguments: `aws_regions` `tag_key` `tag_value` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl
        redshift {
          aws_regions              = ["us-east-1"]
          metrics_polling_interval = 300
          tag_key                  = "tag key"
          tag_value                = "tag value"
        }
``` 
  </details>
  <details>
    <summary>route53</summary>
*  Supported Arguments: `fetch_extended_inventory` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl
        route53 {
          fetch_extended_inventory = true
          metrics_polling_interval = 300
        }
``` 
  </details>
  <details>
    <summary>sns</summary>
*  Supported Arguments: `aws_regions` `fetch_extended_inventory` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 300, 900, 1800, 3600 (seconds)
```hcl  
        sns {
          aws_regions              = ["us-east-1"]
          fetch_extended_inventory = true
          metrics_polling_interval = 300
        }
``` 
  </details>
    <details>
    <summary>security hub</summary>
*  Supported Arguments: `aws_regions` `metrics_polling_interval`
*  Valid `metrics_polling_interval` values: 21600, 43200, 86400 (seconds)
```hcl  
        security_hub {
          aws_regions              = ["us-east-1"]
          metrics_polling_interval = 86400
        }
``` 
  </details>
</details>


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the AWS linked account.

## Additional Examples

```hcl
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
    metrics_polling_interval = 300
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
  security_hub {
    aws_regions              = ["us-east-1"]
    metrics_polling_interval = 86400
  }
}
```
## Import

Linked AWS account integrations can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_cloud_aws_integrations.foo <id>
```