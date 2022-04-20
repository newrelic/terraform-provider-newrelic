---
layout: "newrelic"
page_title: "New Relic Terraform Provider Cloud integrations example for AWS, GCP, and Azure"
sidebar_current: "docs-newrelic-provider-cloud-integrations-guide"
description: |-
  Use this guide to set up the New Relic cloud integrations fully automated through Terraform.
---

## New Relic Terraform Provider Cloud integrations example for AWS, GCP and Azure

The New Relic cloud integrations collect data from cloud platform services and accounts. There's no installation process for cloud integrations and they do not require the use of our infrastructure agent: you simply connect your New Relic account to your cloud provider account. This guide describes the process of enabling the New Relic cloud integrations fully automated through Terraform.

### Requirements

* AWS, Azure or GCP account with administrator permissions
* New Relic account with admin permissions

### Documentation

* [New Relic Cloud integrations](https://docs.newrelic.com/docs/infrastructure/infrastructure-integrations/get-started/introduction-infrastructure-integrations/)

AWS:
* [Introduction to AWS integrations](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/get-started/introduction-aws-integrations)
* [New Relic Terraform newrelic_cloud_aws_link_account resource](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/cloud_aws_link_account)
* [New Relic Terraform newrelic_cloud_aws_integrations resource](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/cloud_aws_integrations)

Azure:
**We are finishing up the development of Azure integrations, and the documentation will be updated once those are available**

GCP:
* [Introduction to gcp integrations](https://docs.newrelic.com/docs/infrastructure/google-cloud-platform-integrations/get-started/introduction-google-cloud-platform-integrations)
* [New Relic Terraform newrelic_cloud_gcp_link_account resource](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/cloud_gcp_link_account)
* [New Relic Terraform newrelic_cloud_gcp_integrations resource](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/cloud_gcp_integrations)

### Examples

Below you can find the full code examples, and required variables for each of the full end to end examples.

#### AWS

The AWS integration relies on two mechanisms to get data in New Relic: [AWS Metric stream](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/aws-metric-stream/) and [AWS Polling integrations](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/connect-aws-new-relic-infrastructure-monitoring). For the majority of AWS services the AWS Metric stream is used as it [has many advantages compared to polling](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/aws-metric-stream#why-it-matters). The AWS Polling integrations are also enabled because AWS does not yet [support all metrics through AWS Metric Stream](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/aws-metric-stream#integrations-not-replaced-streams). We enable `AWS Billing`, `AWS CloudTrail`, `AWS Health`, `AWS Trusted Advisor` and `AWS X-Ray`. Feel free to adapt the example to your needs.

Link: https://github.com/newrelic/terraform-provider-newrelic/blob/main/examples/cloud-integrations-aws.tf

Variables:

* NEW_RELIC_ACCOUNT_ID: The New Relic account you want to link to AWS
* NEW_RELIC_CLOUDWATCH_ENDPOINT: The datacenter where your New Relic account is located.

#### Azure

**We are finishing up the development of Azure integrations, and an example will be greated once those are available**

#### GCP

The AWS integration relies on two mechanisms to get data in New Relic: [AWS Metric stream](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/aws-metric-stream/) and [AWS Polling integrations](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/connect-aws-new-relic-infrastructure-monitoring). For the majority of AWS services the AWS Metric stream is used as it [has many advantages compared to polling](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/aws-metric-stream#why-it-matters). The AWS Polling integrations are also enabled because AWS does not yet [support all metrics through AWS Metric Stream](https://docs.newrelic.com/docs/infrastructure/amazon-integrations/connect/aws-metric-stream#integrations-not-replaced-streams). 

We can enable `App Engine`, `Big Query`, `BigTable`, `Composer`, `Dataflow`, `Dataproc`, `Datastore`, `Firebase Database`, `Firebase Hosting`, `Firebase Storage`, `Firestore`, `Functions`, `Interconnect`, `Kubernetes`, `Load balancing`, `Memcache`, `Pubsub`, `Redis`, `Router`, `Run`, `Spanner`, `Sql`, `Storage`, `Vms`, `Vpc access` Feel free to adapt the example to your needs.

Link: https://github.com/newrelic/terraform-provider-newrelic/blob/main/examples/cloud-integrations-aws.tf

Variables:

* NEW_RELIC_ACCOUNT_ID: The New Relic account you want to link to GCP

