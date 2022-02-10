---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_aws_link_account"
sidebar_current: "docs-newrelic-resource-cloud-aws-link-account"
description: |-
    Link an AWS account to New Relic.
---

# Resource: newrelic\_cloud\_aws\_link\_account

Use this resource to link an AWS account to New Relic.

## Prerequisite

Setup is required in AWS for this resource to work properly. The New Relic AWS integration can be set up to pull metrics from AWS services or AWS can push metrics to New Relic using CloudWatch Metric Streams. 

Using a metric stream to New Relic is the preferred way to integrate with AWS. Follow the [steps outlined here](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/aws-integrations-list/aws-metric-stream/#set-up-metric-stream) to set up a metric stream.

To pull data from AWS instead, complete the [steps outlined here](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/get-started/connect-aws-new-relic-infrastructure-monitoring#connect).

## Example Usage

```hcl
resource "newrelic_cloud_aws_link_account" "foo" {
  arn = "arn:aws:service:region:account-id:resource-id"
  metric_collection_mode = "PUSH"
  name = "account name"
}
```
## Argument Reference

The following arguments are supported:

* `arn` - (Required) The Amazon Resource Name (ARN) of the IAM role.
* `metric_collection_mode` - (Optional) How metrics will be collected. Use `PUSH` for a metric stream or `PULL` to integrate with individual services. 
* `name` - (Required) - The linked account name

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the AWS linked account.

## Import

Linked AWS accounts can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_cloud_aws_link_account.foo <id>
```