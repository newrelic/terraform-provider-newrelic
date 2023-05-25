---
layout: "newrelic"
page_title: "New Relic Terraform Provider Cloud integrations example for AWS, GCP, and Azure"
sidebar_current: "docs-newrelic-provider-cloud-integrations-guide"
description: |-
  Use this guide to set up the New Relic cloud integrations fully automated through Terraform.
---

## New Relic Terraform Provider Cloud integrations example for AWS, GCP and Azure

The [New Relic Cloud integrations](https://docs.newrelic.com/docs/infrastructure/infrastructure-integrations/get-started/introduction-infrastructure-integrations/) collect data from cloud platform services and accounts. There's no installation process for cloud integrations and they do not require the use of our infrastructure agent: you simply connect your New Relic account to your cloud provider account. This guide describes the process of enabling the New Relic cloud integrations fully automated through Terraform.

We have different instruqtions for each cloud provider, use the links below to go the relevant sections:

- [AWS](#aws)
- [Azure](#azure)
- [Google Cloud Platform](#gcp)

If you encounter issues or bugs, please [report those on Github repository](https://github.com/newrelic/terraform-provider-newrelic/issues/new/choose).

### AWS

The New Relic AWS integration relies on two mechanisms to get data in New Relic: [AWS Metric stream](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/aws-metric-stream/) and [AWS Polling integrations](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/connect-aws-new-relic-infrastructure-monitoring). For the majority of AWS services the AWS Metric stream is used as it [has many advantages compared to polling](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/aws-metric-stream#why-it-matters). The AWS Polling integrations are also enabled because AWS does not yet [support all metrics through AWS Metric Stream](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/aws-metric-stream#integrations-not-replaced-streams). We enable `AWS Billing`, `AWS CloudTrail`, `AWS Health`, `AWS Trusted Advisor`, `AWS X-Ray`, `AWS S3`, `AWS Doc DB`, `AWS SQS`, `AWS EBS`, `AWS ALB`, `AWS Elasticache`, `AWS ApiGateway`, `AWS AutoScaling`, `AWS AppSync`, `AWS Athena`, `AWS Cognito`, `AWS Connect`, `AWS DirectConnect`, `AWS Fsx`, `AWS Glue`, `AWS Kinesis Analytics`, `AWS MediaConvert`, `AWS Media Package vod`, `AWS Mq`, `AWS Msk`, `AWS Neptune`, `AWS Qldb`, `AWS Route53resolver`, `AWS States`, `AWS TransitGateway`, `AWS Waf`, `AWS Wafv2`, `Cloudfront`, `DynamoDB`, `Ec2`, `Ecs`, `Efs`, `Elasticbeanstalk`, `Elasticsearch`, `Elb`, `Emr`, `Iam`, `Iot`, `Kinesis`, `Kinesis Firehose`, `Lambda`, `Rds`, `Redshift`, `Route53`, `Ses` and `Sns`.

Check out our [introduction to AWS integrations](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/get-started/introduction-aws-integrations) and the [requirements](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/get-started/integrations-managed-policies/) documentation before continuing with the steps below.

Add the following module to your Terraform code, and set the variables to your desired values. This module assumes you've already set up the New Relic and AWS provider with the correct credentials. If you haven't done so, you can find the instructions here: [New Relic instructions](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/getting_started), [AWS instructions](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#authentication-and-configuration).

```
module "newrelic-aws-cloud-integrations" {
  source = "github.com/newrelic/terraform-provider-newrelic/tree/main/examples/modules/cloud-integrations/aws"

  newrelic_account_id     = 1234567
  newrelic_account_region = "US
  name                    = "production"
}
```
[*You can find the sourcecode for the module on Github.*](https://github.com/newrelic/terraform-provider-newrelic/tree/main/examples/modules/cloud-integrations/aws)


Variables:

* newrelic_account_id: The New Relic account you want to link to AWS. This account will receive all the data observability from your AWS environment.
* newrelic_account_region: The region of your New Relic account, this can be `US` for United States or `EU` for Europe. (Default `US`)
* name: A unique name used throughout the module to name the resources.


### Azure

The Microsoft Azure integrations reports data from various Azure platform services to your New Relic account.

This module will integrate the following resources: `API Management`, `App Gateway`, `App Service`, `Containers`, `Cosmos DB`, `Cost Management`, `Data Factory`, `Eventhub`, `Express Route`, `Firewalls`, `FrontDoor`, `Functions`, `KeyVault`, `Load Balancer`, `Logic apps`, `Machine learning`, `MariaDB`, `Mysql`, `Postgresql`, `PowerBI Dedicated`, `Redis cache`, `Service Bus`, `Sql`, `Sql Managed`, `Storage`, `Virtual Machine`, `Virtual Networks`, `Vms`, and `VPN gateway`.

Check out our [introduction to Azure integrations](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/get-started/introduction-azure-monitoring-integrations/) and the [requirements](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/get-started/activate-azure-integrations/#reqs) documentation before continuing with the steps below.

Add the following module to your Terraform code, and set the variables to your desired values. This module assumes you've already set up the New Relic and Azure provider with the correct credentials. If you haven't done so, you can find the instructions here: [New Relic instructions](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/getting_started), [Azure instructions](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs).

```
module "newrelic-azure-cloud-integrations" {
  source = "github.com/newrelic/terraform-provider-newrelic/tree/main/examples/modules/cloud-integrations/azure"

  newrelic_account_id     = 1234567
  name                    = "production"
}
```
[*You can find the sourcecode for the module on Github.*](https://github.com/newrelic/terraform-provider-newrelic/tree/main/examples/modules/cloud-integrations/azure)

Variables:

* newrelic_account_id: The New Relic account you want to link to Azure. This account will receive all the data observability from your Azure environment.
* name: A unique name used throughout the module to name the resources.


### GCP

The Google Cloud Platform integrations reports data from various GCP services to your New Relic account.

This module will integrate the following resources: `App Engine`, `BigQuery`, `Cloud Functions`, `Cloud Load Balancing`, `Cloud Pub/Sub`, `Cloud Spanner`, `Cloud SQL`, `Cloud Storage`, `Compute Engine`, and `Kubernetes Engine`.

Check out our [introduction to GCP integrations](https://docs.newrelic.com/docs/infrastructure/google-cloud-platform-integrations/get-started/introduction-google-cloud-platform-integrations/) and the [requirements](https://docs.newrelic.com/docs/infrastructure/google-cloud-platform-integrations/get-started/connect-google-cloud-platform-services-new-relic/#reqs) documentation before continuing with the steps below.

Add the following module to your Terraform code, and set the variables to your desired values. This module assumes you've already set up the New Relic and GCP provider with the correct credentials. If you haven't done so, you can find the instructions here: [New Relic instructions](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/getting_started), [GCP instructions](https://registry.terraform.io/providers/hashicorp/google/latest/docs).

```
module "newrelic-gcp-cloud-integrations" {
  source = "github.com/newrelic/terraform-provider-newrelic/tree/main/examples/modules/cloud-integrations/gcp"

  newrelic_account_id     = 1234567
  name                    = "production"
  service_account_id      = 1234567
  project_id              = 1234567
}
```
[*You can find the sourcecode for the module on Github.*](https://github.com/newrelic/terraform-provider-newrelic/tree/main/examples/modules/cloud-integrations/gcp)

Variables:

* newrelic_account_id: The New Relic account you want to link to GCP. This account will receive all the data observability from your Azure environment.
* name: A unique name used throughout the module to name the resources.
* service_account_id: The ID of the New Relic GCP [Service Account](https://cloud.google.com/iam/docs/service-accounts) with [Viewer and Service Usage Consumer roles](https://cloud.google.com/iam/docs/understanding-roles). You can find this ID in the New Relic UI by going to `Infrastructure > GCP > Add a GCP project`. For more information [check out the New Relic docs](https://docs.newrelic.com/docs/infrastructure/google-cloud-platform-integrations/get-started/connect-google-cloud-platform-services-new-relic/).
* project_id: The ID of the project you want to receive data from in GCP.
