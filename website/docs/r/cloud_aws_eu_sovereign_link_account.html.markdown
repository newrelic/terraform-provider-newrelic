---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_aws_eu_sovereign_link_account"
sidebar_current: "docs-newrelic-cloud-resource-aws-eu-sovereign-link-account"
description: |-
  Link an AWS EU Sovereign account to New Relic.
---

# Resource: newrelic\_cloud\_aws\_eu\_sovereign\_link\_account

Use this resource to link an AWS EU Sovereign account to New Relic.

## Prerequisite

Setup is required in AWS EU Sovereign for this resource to work properly. To link an AWS EU Sovereign account to New Relic, you need an AWS EU Sovereign Cloud account.

Using a metric stream to New Relic is the only supported method for AWS EU Sovereign Cloud to get metrics into New Relic for the majority of AWS services. Follow the [steps outlined here](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/connect-aws-eu-sovereign-new-relic/) to set up a metric stream.

To pull data from AWS EU Sovereign for services not supported by CloudWatch Metric Streams (Billing, CloudTrail and X-Ray), complete the [steps outlined here](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/connect-aws-eu-sovereign-new-relic/).

## Example Usage

You can also use the [full example, including the AWS set up, found in our guides](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/cloud_integrations_guide#aws-eu-sovereign).

```hcl
resource "newrelic_cloud_aws_eu_sovereign_link_account" "foo" {
  account_id             = 1234567
  name                   = "My New Relic - AWS EU Sovereign Linked Account"
  metric_collection_mode = "PUSH" # Options: "PUSH", "PULL", or "BOTH"
  arn                    = "arn:aws-eusc:iam::123456789012:role/NewRelicInfrastructure-Integrations"
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The New Relic account ID to operate on. This allows the user to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`, if not specified in the configuration.
* `arn` - (Required) The Amazon Resource Name (ARN) of the IAM role.
* `metric_collection_mode` - (Optional) How metrics will be collected. Use `PUSH` for metric stream, `PULL` for API polling of the 3 services not supported by metric streams (Billing, CloudTrail and X-Ray), or `BOTH` for both methods. Defaults to `PUSH`, if not specified in the configuration.
* `name` - (Required) The name/identifier of the AWS EU Sovereign - New Relic 'linked' account.

-> **WARNING:** Updating any of the aforementioned attributes (except `name`) of a `newrelic_cloud_aws_eu_sovereign_link_account` resource that has been applied would **force a replacement** of the resource (destruction of the resource, followed by the creation of a new resource). Please carefully review the output of `terraform plan`, which would clearly indicate a replacement of this resource, before performing a `terraform apply`.

-> **NOTE:** This resource requires the New Relic provider to be configured with `region = "EU"` or the `NEW_RELIC_REGION=EU` environment variable.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the AWS EU Sovereign linked account.

## Import

Linked AWS EU Sovereign accounts can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_cloud_aws_eu_sovereign_link_account.foo <id>
```