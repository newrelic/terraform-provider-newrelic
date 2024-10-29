---
layout: "newrelic"
page_title: "New Relic Terraform Provider Cloud integrations example for AWS, GCP, and Azure"
sidebar_current: "docs-newrelic-provider-cloud-integrations-guide"
description: |-
  Use this guide to set up the New Relic cloud integrations fully automated through Terraform.
---

## New Relic Terraform Provider Cloud integrations example for AWS, GCP and Azure

The [New Relic Cloud integrations](https://docs.newrelic.com/docs/infrastructure/infrastructure-integrations/get-started/introduction-infrastructure-integrations/) collect data from cloud platform services and accounts. There's no installation process for cloud integrations and they do not require the use of our infrastructure agent: you simply connect your New Relic account to your cloud provider account. This guide describes the process of enabling the New Relic cloud integrations fully automated through Terraform.

We have different instructions for each cloud provider, use the links below to go the relevant sections:

- [AWS](#aws)
- [Azure](#azure)
- [Google Cloud Platform](#gcp)

If you encounter issues or bugs, please [report those on Github repository](https://github.com/newrelic/terraform-provider-newrelic/issues/new/choose).

### AWS

The New Relic AWS integration relies on two mechanisms to get data in New Relic: [AWS Metric stream](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/aws-metric-stream/) and [AWS Polling integrations](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/connect-aws-new-relic-infrastructure-monitoring). For the majority of AWS services the AWS Metric stream is used as it [has many advantages compared to polling](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/aws-metric-stream#why-it-matters). AWS Polling integrations are also enabled because AWS does not yet [support all metrics through AWS Metric Stream](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/aws-metric-stream#integrations-not-replaced-streams). 

The following AWS services may be integrated via API Polling, using the New Relic Terraform Provider. A combination of these services may be added to the Terraform Configuration, to set up an AWS Integration comprising these services via API Polling. More [examples](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/cloud_aws_integrations#example-usage) and [details](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/cloud_aws_integrations#argument-reference) on the arguments needed to set up each service can be found in the documentation of the [`newrelic_cloud_aws_integrations`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/cloud_aws_integrations) resource.

|                     |                     |                    |
|---------------------|---------------------|--------------------|
| `ALB`               | `API Gateway`       | `AppSync`          |
| `Athena`            | `Auto Scaling`      | `Billing`          |
| `CloudFront`        | `CloudTrail`        | `Cognito`          |
| `Connect`           | `Direct Connect`    | `DocumentDB`       |
| `DynamoDB`          | `EBS`               | `EC2`              |
| `ECS`               | `EFS`               | `ElastiCache`      |
| `Elastic Beanstalk` | `Elasticsearch`     | `ELB`              |
| `EMR`               | `FSx`               | `Glue`             |
| `Health`            | `IAM`               | `IoT`              |
| `Kinesis`           | `Kinesis Analytics` | `Kinesis Firehose` |
| `Lambda`            | `Media Package VOD` | `MediaConvert`     |
| `MQ`                | `MSK`               | `Neptune`          |
| `QLDB`              | `RDS`               | `Redshift`         |
| `Route53`           | `Route53 Resolver`  | `S3`               |
| `SES`               | `SNS`               | `SQS`              |
| `States`            | `Transit Gateway`   | `Trusted Advisor`  |
| `VPC`               | `WAF`               | `WAFv2`            |
| `X-Ray`             | 



Check out our [introduction to AWS integrations](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/get-started/introduction-aws-integrations) and the [requirements](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/get-started/integrations-managed-policies/) documentation before continuing with the steps below.

The GitHub repository of the Terraform Provider also has an AWS Cloud Integrations 'module', that can be used to simplify setting up an AWS integration; both, via metric streams, and with API polling (inclusive of a few of the aforementioned services). To use this module, add the following to your Terraform code, and set the variables to your desired values.

-> **NOTE:** This module assumes you've already set up the New Relic and AWS provider with the correct credentials. If you haven't done so, you can find the instructions here: [New Relic instructions](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/getting_started), [AWS instructions](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#authentication-and-configuration).

```
module "newrelic-aws-cloud-integrations" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/cloud-integrations/aws"

  newrelic_account_id     = 1234567
  newrelic_account_region = "US"
  name                    = "production"

  include_metric_filters = {
    "AWS/EC2" = [], # include ALL metrics from the EC2 namespace
    "AWS/S3" = ["NumberOfObjects"]. # include just a specific metric from the S3 namespace
  }
}
```

[*You can find the sourcecode for the module on Github.*](https://github.com/newrelic/terraform-provider-newrelic/tree/main/examples/modules/cloud-integrations/aws)


Variables:

* `newrelic_account_id`: The New Relic account you want to link to AWS. This account will receive all the data observability from your AWS environment.
* `newrelic_account_region` (Optional): The region of your New Relic account, this can be `US` for United States or `EU` for Europe. (Default `US`)
* `name` (Optional): A unique name used throughout the module to name the resources. (Default `production`)
* `exclude_metric_filters` (Optional): a map of namespaces and metric names to exclude from the Cloudwatch metric stream. `Conflicts with include_metric_filters`.
* `include_metric_filters` (Optional): a map of namespaces and metric names to include in the Cloudwatch metric stream. `Conflicts with exclude_metric_filters`.

### Azure

The Microsoft Azure integrations reports data from various Azure platform services to your New Relic account.

The following Azure services may be integrated using the New Relic Terraform Provider. A combination of these services may be added to the Terraform Configuration, to set up an Azure Integration comprising the selected services. More [examples](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/cloud_azure_integrations#example-usage) and [details](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/cloud_azure_integrations#argument-reference) on the arguments needed to set up each service can be found in the documentation of the [`newrelic_cloud_azure_integrations`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/cloud_azure_integrations) resource.

|                    |                       |                      |
|--------------------|-----------------------|----------------------|
| `API Management`   | `Application Gateway` | `App Service`        |
| `Containers`       | `Cosmos DB`           | `Cost Management`    |
| `Data Factory`     | `Event Hub`           | `Express Route`      |
| `Firewalls`        | `Front Door`          | `Functions`          |
| `Key Vault`        | `Load Balancer`       | `Logic Apps`         |
| `Machine Learning` | `MariaDB`             | `Monitor`            |
| `MySQL`            | `PostgreSQL`          | `Power BI Dedicated` |
| `Redis Cache`      | `Service Bus`         | `SQL`                |
| `SQL Managed`      | `Storage`             | `Virtual Machine`    |
| `Virtual Network`  | `VMs`                 | `VPN Gateway`        |

Check out our [introduction to Azure integrations](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/get-started/introduction-azure-monitoring-integrations/) and the [requirements](https://docs.newrelic.com/docs/infrastructure/microsoft-azure-integrations/get-started/activate-azure-integrations/#reqs) documentation before continuing with the steps below.

The GitHub repository of the Terraform Provider also has an Azure Cloud Integrations 'module', comprising all the aforementioned Azure services, that can be used to simplify setting up an Azure Integration. To use this module, add the following to your Terraform code, and set the variables to your desired values. 

-> **NOTE:** This module assumes you've already set up the New Relic and Azure provider with the correct credentials. If you haven't done so, you can find the instructions here: [New Relic instructions](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/getting_started), [Azure instructions](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs).


```
module "newrelic-azure-cloud-integrations" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/cloud-integrations/azure"

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

The following GCP services may be integrated using the New Relic Terraform Provider. A combination of these services may be added to the Terraform Configuration, to set up an GCP Integration comprising the selected services. More [examples](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/cloud_gcp_integrations#example-usage) and [details](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/cloud_gcp_integrations#argument-reference) on the arguments needed to set up each service can be found in the documentation of the [`newrelic_cloud_gcp_integrations`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/cloud_gcp_integrations) resource.

|                      |                           |                       |
|----------------------|---------------------------|-----------------------|
| `Alloy DB`           | `App Engine`              | `BigQuery`            |
| `Cloud Bigtable`     | `Cloud Composer`          | `Cloud Functions`     |
| `Cloud Interconnect` | `Cloud Load Balancing`    | `Cloud Pub/Sub`       |
| `Cloud Router`       | `Cloud Run`               | `Cloud SQL`           |
| `Cloud Spanner`      | `Cloud Storage`           | `Dataflow`            |
| `Dataproc`           | `Datastore`               | `Firebase Database`   |
| `Firebase Hosting`   | `Firebase Storage`        | `Firestore`           |
| `Kubernetes Engine`  | `Memorystore (Memcached)` | `Memorystore (Redis)` |
| `VPC Access`         | `Virtual Machines`        |

Check out our [introduction to GCP integrations](https://docs.newrelic.com/docs/infrastructure/google-cloud-platform-integrations/get-started/introduction-google-cloud-platform-integrations/) and the [requirements](https://docs.newrelic.com/docs/infrastructure/google-cloud-platform-integrations/get-started/connect-google-cloud-platform-services-new-relic/#reqs) documentation before continuing with the steps below.

The GitHub repository of the Terraform Provider also has an GCP Cloud Integrations 'module', comprising a few of the aforementioned GCP services as a sample, that can be used to simplify setting up an GCP Integration. To use this module, add the following to your Terraform code, and set the variables to your desired values.

-> **NOTE:** This module assumes you've already set up the New Relic and GCP provider with the correct credentials. If you haven't done so, you can find the instructions here: [New Relic instructions](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/getting_started), [GCP instructions](https://registry.terraform.io/providers/hashicorp/google/latest/docs).

```
module "newrelic-gcp-cloud-integrations" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/cloud-integrations/gcp"

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
