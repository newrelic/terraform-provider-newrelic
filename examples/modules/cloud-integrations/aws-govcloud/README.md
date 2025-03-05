# AWS GovCloud Integration Module

This Terraform module links an AWS GovCloud account to New Relic. It sets up the necessary IAM roles, policies, and configurations to enable metric collection from AWS GovCloud to New Relic.

## Prerequisites

- An AWS GovCloud account.
- A New Relic account with the necessary permissions to create and manage cloud integrations.

## How the Module Works

This module creates several AWS resources, including IAM roles, policies, S3 buckets, and Kinesis Firehose delivery streams, to facilitate the integration between AWS GovCloud and New Relic. It also configures CloudWatch metric streams to send metrics to New Relic.

It then uses the following resources based on cloud integrations from the New Relic Terraform Provider (the ARN used with the New Relic Terraform Provider's cloud integrations resources comes from the AWS resources deployed by the module, as stated above):

- `newrelic_cloud_aws_govcloud_link_account`: Links the AWS GovCloud account to New Relic.
- `newrelic_cloud_aws_govcloud_integrations`: Configures various AWS (GovCloud) services (e.g., EC2, S3, RDS) to send metrics to New Relic.

## Usage

> **Note:** Using this module requires a minimum version of `3.56.0` of the New Relic Terraform Provider.

```hcl
module "newrelic-aws-govcloud-integrations" {
  source                  = "github.com/newrelic/terraform-provider-newrelic//examples/modules/cloud-integrations/aws-govcloud"
  newrelic_account_id     = 1234321
  name                    = "Test AWS GovCloud Integrations"

  include_metric_filters = {
    # "AWS/EC2" = [], # include ALL metrics from the EC2 namespace
  }
}
```

> **Note:** If you have cloned this repo and would like to deploy this configuration in `main.tf` in the `testing` folder, use the following path as the value of the `source` argument:

```hcl
  source                  = "../examples/modules/cloud-integrations/aws-govcloud"
```

### A Note on 'Applying' the Module :warning:

When applying this module, please use reduced parallelism (ideally `--parallelism=1`) with the `terraform apply` command. The volume of resources in the module sometimes leads to a race condition where AWS resources are applied and the ARN is made available for the `newrelic_aws_govcloud_link_account` resource, but not yet updated on the AWS backend. This could sometimes lead to validation issues with the ARN. To avoid this, reduced parallelism can keep the apply operation streamlined and ensure adequate time for the ARN to be available and valid on the AWS backend.

```sh
terraform apply --parallelism=1
```

## Variables

- `newrelic_account_id` - (Required) The New Relic account ID.
- `name` - (Required) The name/identifier for the integration.
- `exclude_metric_filters` - (Optional) A map of metric namespaces and metric names to exclude from the metric stream.
- `include_metric_filters` - (Optional) A map of metric namespaces and metric names to include in the metric stream.