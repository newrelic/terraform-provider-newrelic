---
layout: "newrelic"
page_title: "New Relic Terraform Provider Cloud integrations example for AWS, GCP, Azure, and OCI"
sidebar_current: "docs-newrelic-provider-cloud-integrations-guide"
description: |-
  Use this guide to set up the New Relic cloud integrations fully automated through Terraform.
---

## New Relic Terraform Provider Cloud integrations example for AWS, GCP, Azure, and OCI

The [New Relic Cloud integrations](https://docs.newrelic.com/docs/infrastructure/infrastructure-integrations/get-started/introduction-infrastructure-integrations/) collect data from cloud platform services and accounts. There's no installation process for cloud integrations and they do not require the use of our infrastructure agent: you simply connect your New Relic account to your cloud provider account. This guide describes the process of enabling the New Relic cloud integrations fully automated through Terraform.

We have different instructions for each cloud provider, use the links below to go the relevant sections:

- [AWS](#aws)
- [Azure](#azure)
- [Google Cloud Platform](#gcp)
- [Oracle Cloud Infrastructure](#oci)

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
  output_format           = "opentelemetry0.7"

  include_metric_filters = {
    "AWS/EC2" = [], # include ALL metrics from the EC2 namespace
    "AWS/S3" = ["NumberOfObjects"]. # include just a specific metric from the S3 namespace
  }
  enable_config_recorder = true # Set to true to enable AWS Config Configuration Recorder
}
```

[*You can find the sourcecode for the module on Github.*](https://github.com/newrelic/terraform-provider-newrelic/tree/main/examples/modules/cloud-integrations/aws)


Variables:

* `newrelic_account_id`: The New Relic account you want to link to AWS. This account will receive all the data observability from your AWS environment.
* `newrelic_account_region` (Optional): The region of your New Relic account, this can be `US` for United States or `EU` for Europe. (Default `US`)
* `name` (Optional): A unique name used throughout the module to name the resources. (Default `production`)
* `output_format` (Optional): The output format for telemetry data. Supported values are `opentelemetry0.7` and `opentelemetry1.0`. (Default `opentelemetry0.7`)
* `exclude_metric_filters` (Optional): a map of namespaces and metric names to exclude from the Cloudwatch metric stream. `Conflicts with include_metric_filters`.
* `include_metric_filters` (Optional): a map of namespaces and metric names to include in the Cloudwatch metric stream. `Conflicts with exclude_metric_filters`.
* `enable_config_recorder` (Optional): Set to `true` to enable creation of an [AWS Config Configuration Recorder](https://docs.aws.amazon.com/config/latest/developerguide/stop-start-recorder.html) in your AWS account. Only one recorder is allowed per region per account. Default is `false`.

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

<a id="oci"></a>
### Oracle Cloud Infrastructure

The New Relic OCI integration collects metrics, logs, and metadata from supported OCI services and sends them to your New Relic account. This integration uses a combination of:

* Service Connector Hub pipelines (for metrics / logs export)
* Functions for data transformation and payload enrichment
* API polling to supplement metadata and tags

#### Supported OCI service categories

The following OCI namespaces are supported by New Relic for metrics/logs collection via Service Connector Hub. These namespaces correspond to different OCI services and can be used when configuring your Service Connector Hub to export metrics to New Relic.

|                            |                                    |                             |
|----------------------------|------------------------------------|-----------------------------| 
| `API Gateway`              | `Block Storage`                    | `Compute`                   |
| `Compute Agent`            | `Compute Infrastructure Health`    | `Compute Instance Health`   |
| `Container Engine (OKE)`   | `Container Instance`               | `Functions (FaaS)`          |
| `Health Checks`            | `Instance Pools`                   | `Load Balancer (LBaaS)`     |
| `Logging`                  | `Network Load Balancer`            | `Object Storage`            |
| `PostgreSQL`               | `Queue`                            | `Service Connector Hub`     |
| `Streaming`                |                                    |                             |


#### Modular OCI setup

The modular approach allows you to deploy only the specific OCI integration components you need, with clear separation between policy setup and data collection integrations. This design enables flexible deployment strategies where policy configuration can be managed independently from metrics and logging integrations.

The following composable modules are available under `examples/modules/cloud-integrations/oci/` so you can provision only what you need:

* `policy-setup` – Creates IAM policies and identity trust / configuration prerequisites (including workload identity federation inputs) required to link an OCI tenancy to New Relic.
* `metrics-integration` – Creates Service Connector Hub resources, optional networking (VCN / subnets), and supporting artifacts that export metrics (and optionally logs) to New Relic.
* `logs-integration` – Creates connector hubs, function and function app to export logs from Oracle Cloud to New Relic.

Use them independently or combine them in the same configuration. In all cases, the `policy-setup` module must be applied successfully before the `metrics-integration` or `logging-integration` module, because the latter depends on IAM policies, dynamic groups / identity trust, and (if configured) workload identity federation artifacts created by the former.

#### Example: Policy setup module

```hcl
module "oci_policy_setup" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/cloud-integrations/oci/policy-setup"

  tenancy_ocid      = "ocid1.tenancy.oc1..aaaaaaaaexampletenancy"
  region            = "iad"
  fingerprint       = "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0"
  private_key       = "USER_PVT_KEY"

  # New Relic linkage / API keys (dummy values)
  newrelic_account_id      = 1234567
  newrelic_ingest_api_key  = "NRII-INGEST-API-KEY-EXAMPLE"
  newrelic_user_api_key    = "NRAA-USER-API-KEY-EXAMPLE"
  newrelic_provider_region = "US" # or "EU"

  # Workload Identity Federation / OAuth2 (sample values)
  client_id      = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  client_secret  = "super-secret-client-value"
  oci_domain_url = "https://idcs-abcdef1234567890.identity.oraclecloud.com"
  svc_user_name  = "svc-newrelic-wif"

  # Optional: Use existing vault secrets (leave empty to create new ones)
  # user_key_secret_ocid   = "ocid1.vaultsecret.oc1..existing-user-secret"
  # ingest_key_secret_ocid = "ocid1.vaultsecret.oc1..existing-ingest-secret"

  # Enable metrics & logs policies (example)
  instrumentation_type = "METRICS,LOGS"
}
```

Key variables:

* `instrumentation_type` – Comma‑separated list of any of `METRICS`, `LOGS`, `METRICS,LOGS` controlling which policy sets are deployed.
* `client_id`, `client_secret`, `oci_domain_url`, `svc_user_name` – Workload identity federation (OAuth2) inputs (see the [OCI link account](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/cloud_oci_link_account) resource docs for guidance).
* `newrelic_provider_region` – Region context for New Relic provider operations (for example, `US` or `EU`).
* `user_key_secret_ocid` / `ingest_key_secret_ocid` (Optional) – OCIDs of existing vault secrets containing New Relic API keys. Leave empty to create new vault secrets.

#### Example: Metrics integration module

```hcl
module "oci_metrics_integration" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/cloud-integrations/oci/metrics-integration"

  tenancy_ocid     = "ocid1.tenancy.oc1..aaaaaaaaexampletenancy"
  compartment_ocid = "ocid1.compartment.oc1..bbbbbbbbexamplecmp" # or module.oci_policy_setup.compartment_ocid
  region           = "iad"
  fingerprint      = "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0"
  private_key      = "USER_PVT_KEY"

  # New Relic account configuration
  newrelic_account_id   = "1234567"
  provider_account_id   = "1234567"

  # Endpoint selection (validated internally)
  newrelic_endpoint = "US" # or EU

  # Networking
  create_vcn         = true
  function_subnet_id = "" # ignored when create_vcn = true

  # Vault secret OCIDs (dummy)
  ingest_api_secret_ocid = "ocid1.vaultsecret.oc1..dddddddigingestsecret" # or module.oci_policy_setup.ingest_vault_ocid
  user_api_secret_ocid   = "ocid1.vaultsecret.oc1..eeeeeeeeusersecret123" # or module.oci_policy_setup.user_vault_ocid

  # Docker image configuration (optional)
  image_version = "latest"
  image_bucket  = "idptojlonu4e"

  connector_hubs_data = "[{\"compartments\":[{\"compartment_id\":\"ocid1.tenancy.oc1..aaaaaaaaexampletenancy\",\"namespaces\":[\"oci_faas\"]}],\"description\":\"[DO NOT DELETE] New Relic Metrics Connector Hub\",\"name\":\"newrelic-metrics-connector-hub-us-ashburn\"}]"
}
```

Key variables:

* `newrelic_account_id` / `provider_account_id` – New Relic account identifiers for linking the OCI integration.
* `create_vcn` / `function_subnet_id` – Networking control. Set `create_vcn=false` and provide an existing `function_subnet_id` to reuse existing infrastructure.
* `connector_hubs_data` – A JSON *string* (must be valid, stringified JSON) whose root is an array of connector hub definition objects. 
  
  Each object supports:
  * `compartments` (array of objects with `compartment_id` and `namespaces` (array of strings))
  * `description` (string)
  * `name` (string)
  The example above shows a single‑element JSON array wrapped in quotes to satisfy Terraform's string input expectation. Example object structure:

  ```json
  [
    {
      "compartments": [
        {
          "compartment_id": "ocid1.tenancy.oc1..aaaaaaaaexampletenancy",
          "namespaces": ["oci_faas"]
        }
      ],
      "description": "[DO NOT DELETE] New Relic Metrics Connector Hub",
      "name": "newrelic-metrics-connector-hub-us-ashburn"
    }
  ]
  ```
* `ingest_api_secret_ocid` / `user_api_secret_ocid` – Vault secret OCIDs for ingest and user API keys (avoid embedding plain‑text keys).
* `newrelic_endpoint` – Logical endpoint selector; the module maps this value to the actual metric ingest URL (use the EU variant for EU accounts).
* `region` – OCI region key (short code) where resources for this module are created (for example: `iad`, `phx`, `fra`). Provide ONLY the region key, not the full region identifier (so use `iad` instead of `us-ashburn-1`).
* `image_version` / `image_bucket` – Docker image configuration for the New Relic function (optional, defaults to latest version).

#### Example: Logs integration module

```hcl
module "oci_logs_integration" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/cloud-integrations/oci/logs-integration"

  # oci configuration
  tenancy_ocid     = "ocid1.tenancy.oc1..***"
  compartment_ocid = module.oci_policy_setup.compartment_ocid
  region           = "us-ashburn-1"
  
  # New Relic account configuration
  newrelic_account_id = "1234567"
  provider_account_id = "1234567"
  
  # new relic logging prefix
  newrelic_logging_identifier = "logs"
  
  # network components
  create_vcn         = true # set to false to reuse existing VCN/subnet created from metrics module
  function_subnet_id = ""   # ignored when create_vcn = true
  
  # function application environment variables configuration
  image_version        = "latest" # latest image version for the logging function
  debug_enabled        = "FALSE"
  new_relic_region     = "US"
  secret_ocid          = module.oci_policy_setup.ingest_vault_ocid
  user_api_secret_ocid = module.oci_policy_setup.user_vault_ocid
  
  # connector hub configuration (Optional)
  # Don't add the following variables if you want to skip log export.
  connector_hub_details = "[{\"display_name\":\"newrelic-logs-connector\",\"description\":\"Service connector for logs from compartment A to New Relic\",\"log_sources\":[{\"compartment_id\":\"ocid1.tenancy.oc1..***\",\"log_group_id\":\"ocid1.loggroup.oc1.iad.***\"}]}]"
  batch_size_in_kbs     = 6000 # max payload size in KBs (default 6000)
  batch_time_in_sec     = 60   # max wait time in seconds before sending batch (default 60)
}
```

Key variables:

- New Relic account configuration:
  - `newrelic_account_id` / `provider_account_id`: New Relic account identifiers for linking the OCI integration.
- network components:
  - `create_vcn`: set to false to reuse existing VCN/subnet created from metrics module. 
  - `function_subnet_id`: subnet OCID for the function to be created in. Ignored if create_vcn is true.
> If you want to use an existing private subnet, make sure it has required route rules and gateways with internet and all OCI services access. 
- function application environment variables configuration:
  - `debug_enabled`: Boolean to enable or disable function debug logs.
  - `new_relic_region`: The New Relic region (US or EU).
  - `secret_ocid`: The OCID of the secret in OCI Vault containing New Relic License Key.
  - `user_api_secret_ocid`: The OCID of the secret in OCI Vault containing New Relic User API Key.
  - `image_version`: Docker image version for the logging function (defaults to "latest").
- connector hub configuration: A JSON *string* (must be valid, stringified JSON) whose root is an array of connector hub definition objects. Each object supports:
  * `display_name` (string) : name of the connector hub - must have prefix `newrelic-logs`
  * `description` (string) (optional): connector hub description
  * `log_sources`: 
    * list of compartment OCID and log group OCID

The example above shows a single‑element JSON array wrapped in quotes to satisfy Terraform's string input expectation. Example object structure:
     
```json
[
  {
    "display_name": "newrelic-logs-connector",
    "description": "Service connector for logs from compartment A to New Relic",
    "log_sources": [
      {
        "compartment_id": "ocid1.compartment.oc1..***",
        "log_group_id": "ocid1.loggroup.oc1.iad.***"
      },
      {
        "compartment_id": "ocid1.compartment.oc1..***",
        "log_group_id": "ocid1.loggroup.oc1.iad.***"
      }
    ]
  }
]
```

> When implementing the New Relic OCI integration, the `policy-setup` module must always be applied before the `metrics-integration` or `logging-integration` modules. These modules can be run together in a single Terraform configuration only if the dependency graph can be successfully resolved. For example, this can be achieved by referencing outputs from the `policy-setup` module in the other integration modules. Failure to apply the necessary policies first will result in authorization errors when creating Service Connector Hub resources or invoking functions.

[*Browse the OCI module source code on GitHub*](https://github.com/newrelic/terraform-provider-newrelic/tree/main/examples/modules/cloud-integrations/oci)
