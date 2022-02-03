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

Setup is required in AWS for this resource to work properly. Complete the [steps outlined here](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/get-started/connect-aws-new-relic-infrastructure-monitoring#connect) before using this resource.

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
* `metric_collection_mode` - (Optional) How metrics will be collected. One of `PUSH` or `PULL`
* `name` - (Required) - The linked account name

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the AWS linked account.

## Import

Linked AWS accounts can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_cloud_aws_link_account.foo <id>
```