---
layout: "newrelic"
page_title: "New Relic: newrelic_awsGovcloud_ink_account"
sidebar_current: "docs-newrelic-resource-awsGovCloud-link-account"
description: |-
  Link an AwsGovCloud account to New Relic.
---
-> **IMPORTANT!** We do not have access to AWS GovCloud account and can't properly test this resource.


# Resource: newrelic_cloud_aws_govcloud_link_account

Use this resource to link an AWSGovCloud account to New Relic.

## Prerequisite

Obtain the AwsGovCloud account designed to address the specific regulatory needs of United States (federal, state, and local agencies), education institutions, and the supporting ecosystem.

It is an isolated AWS region designed to host sensitive data and regulated workloads in the cloud, helping customers support their US government compliance requirements.

To pull data from AWSGovCloud, complete the [steps outlined here](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/get-started/connect-aws-govcloud-new-relic).

## Example Usage

```hcl
resource "newrelic_awsGovcloud_link_account" "foo" {
  account_id = "The New Relic account ID where you want to link the AWSGovCloud account"
  access_key_id = "access-key-id of awsGovcloud account"
  aws_account_id = "wsGovcloud account id"
  metric_collection_mode = "PULL"
  name = "account name"
  secret_access_key = "secret access key of the awsGovcloud account"
}
```

## Argument Reference

The following arguments are supported:

- `account_id` - (Optional) The New Relic account ID to operate on. This allows the user to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.
- `access_key_id` - (Required) The access key of the AwsGovCloud.
- `aws_account_id` - (Required) The AwsGovCloud account ID.
- `secret_access_key` - (Required) The secret key of the AwsGovCloud.
- `metric_collection_mode` - (Optional) How metrics will be collected. Use `PUSH` for a metric stream or `PULL` to integrate with individual services.
- `name` - (Required) - The linked account name

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the AWSGovCloud linked account.

## Import

Linked AWSGovCloud accounts can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_cloud_aws_govcloud_link_account.foo <id>
```
