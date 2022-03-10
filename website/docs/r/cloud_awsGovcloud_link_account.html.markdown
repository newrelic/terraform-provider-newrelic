
---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_awsGovcloud_link_account"
sidebar_current: "docs-newrelic-resource-cloud-awsGovcloud-link-account"
description: |-
Link an AWSGovCloud account to New Relic.
---

# Resource: newrelic\_cloud\_awsGovCloud\_link\_account

Use this resource to link an AWSGovCloud account to New Relic.

## Prerequisite

Obtain the AwsGovCloud account designed to address the specific regulatory needs of United States (federal, state, and local agencies), education institutions, and the supporting ecosystem.

It is an isolated AWS region designed to host sensitive data and regulated workloads in the cloud, helping customers support their US government compliance requirements.

To pull data from AWSGovCloud, complete the [steps outlined here](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/get-started/connect-aws-govcloud-new-relic).

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

* `account_id` - (Optional) The New Relic account ID to operate on.  This allows the user to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.
* `access_key_id` - (Required) The access key of the AwsGovCloud.
* `aws_account_id` - (Required) The AwsGovCloud account id.
* `secret_access_key` - (Required) The secret key of the AwsGovCloud.
* `metric_collection_mode` - (Optional) How metrics will be collected. Use `PUSH` for a metric stream or `PULL` to integrate with individual services.
* `name` - (Required) - The linked account name

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the AWSGovCloud linked account.

## Import

Linked AWS accounts can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_cloud_awsGovcloud_link_account.foo <id>
```
