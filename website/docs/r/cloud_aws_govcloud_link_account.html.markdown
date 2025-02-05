---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_aws_govcloud_link_account"
sidebar_current: "docs-newrelic-cloud-resource-aws-govcloud-link-account"
description: |-
  Link an AWS GovCloud account to New Relic.
---
-> **IMPORTANT!** This resource is in alpha state, and could still contain issues and missing functionality. If you encounter any issue please create a ticket on [Github](https://github.com/newrelic/terraform-provider-newrelic/issues/new/choose) with all the required information.

# Resource: newrelic_cloud_aws_govcloud_link_account

Use this resource to link an AWS GovCloud account to New Relic.

## Prerequisite

Obtain the AwsGovCloud account designed to address the specific regulatory needs of United States (federal, state, and local agencies), education institutions, and the supporting ecosystem.

It is an isolated AWS region designed to host sensitive data and regulated workloads in the cloud, helping customers support their US government compliance requirements.

To pull data from AWSGovCloud, complete the [steps outlined here](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/get-started/connect-aws-govcloud-new-relic).

## Example Usage

```hcl
resource "newrelic_cloud_aws_govcloud_link_account" "foo" {
  account_id             = 1234567
  name                   = "My New Relic - AWS GovCloud Linked Account"
  metric_collection_mode = "PUSH"
  aws_account_id         = "<Your AWS GovCloud Account's ID>"
  access_key_id          = "<Your AWS GovCloud Account's Access Key ID>"
  secret_access_key      = "<Your AWS GovCloud Account's Secret Access Key>"
}
```

## Argument Reference

The following arguments are supported:

- `account_id` - (Optional) The New Relic account ID to operate on. This allows the user to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.
- `name` - (Required) - The name/identifier of the AWS GovCloud - New Relic 'linked' account.
- `metric_collection_mode` - (Optional) The mode by which metric data is to be collected from the linked AWS GovCloud account. Use `PUSH` for Metric Streams and `PULL` for API Polling based metric collection respectively.
  - Note: Altering the `metric_collection_mode` of an already applied `newrelic_cloud_aws_govcloud_link_account` resource shall trigger a recreation of the resource, instead of an update.
- `aws_account_id` - (Required) The ID of the AWS GovCloud account.
- `access_key_id` - (Required) The Access Key used to programmatically access the AWS GovCloud account.
- `secret_access_key` - (Required) The Secret Access Key used to programmatically access the AWS GovCloud account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the AWS GovCloud linked account.

## Import

Linked AWS GovCloud accounts can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_cloud_aws_govcloud_link_account.foo <id>
```
