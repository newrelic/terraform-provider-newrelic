---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_oci_logging_integrations"
sidebar_current: "docs-newrelic-resource-cloud-oci-logging-integrations"
description: |-
Integrate Oracle Cloud Logging Services with New Relic.
---

# Resource: newrelic\_cloud\_oci\_logging\_integrations

This Terraform Module links Oracle Logs and Log Groups to New Relic.

## Prerequisites

- An Oracle Cloud Infrastructure (OCI) account with the necessary permissions to create resources.
- A New Relic account with the necessary permissions to create and manage cloud integrations.

## How the Module Works

This module creates several OCI resources including Functions, Service Connector Hub, IAM policies, and Networking Components to enable logs collection from Oracle to New Relic.

## Usage

> **Note:** Using this module requires a minimum version of `3.56.0` of the New Relic Terraform Provider.

```hcl
module "newrelic-aws-govcloud-integrations" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/cloud-integrations/oci/logging-integrations"

  # Variables
  tenancy_ocid = "ocid1.tenancy.oc1..***"
  compartment_ocid = "ocid1.tenancy.oc1..***"
  region = "us-ashburn-1"
  newrelic_logging_prefix = "nr_logging"
  create_vcn = false
  function_subnet_id = "ocid1.subnet.oc1.iad.***"
  debug_enabled = "FALSE"
  new_relic_region = "US"
  secret_ocid = "ocid1.vaultsecret.oc1.iad.***"
  log_sources_details = "[{\"display_name\":\"nr-service-connector-1\",\"description\":\"Service connector for logs from compartment A to New Relic\",\"log_sources\":[{\"compartment_id\":\"ocid1.tenancy.oc1..***\",\"log_group_id\":\"ocid1.loggroup.oc1.iad.***\"}]}]"
}
```

> **Note:** If you have cloned this repo and would like to deploy this configuration in `main.tf` in the `testing` folder, use the following path as the value of the `source` argument:

```hcl
  source                  = "../examples/modules/cloud-integrations/oci/logging-integrations"
```

### A Note on 'Applying' the Module :warning:

When applying this module, please use reduced parallelism (ideally `--parallelism=1`) with the `terraform apply` command. The volume of resources in the module sometimes leads to a race condition where AWS resources are applied and the ARN is made available for the `newrelic_aws_govcloud_link_account` resource, but not yet updated on the AWS backend. This could sometimes lead to validation issues with the ARN. To avoid this, reduced parallelism can keep the apply operation streamlined and ensure adequate time for the ARN to be available and valid on the AWS backend.

```sh
terraform apply --parallelism=1
```

## Variables
- `tenancy_id`: The OCID of your OCI tenancy.
- `compartment_id`: The OCID of the compartment where New Relic logging resources will be created.
- `newrelic_logging_prefix`: Prefix for naming New Relic logging resources.
- `region`: The OCI region where resources will be created.
- `create_vcn`: Boolean to determine if a new VCN should be created.
- `function_subnet_id`: The OCID of the subnet for the function if new VCN is not created.
- `debug_enabled`: Boolean to enable or disable debug logging.
- `new_relic_region`: The New Relic region (US or EU).
- `secret_ocid`: The OCID of the secret in OCI Vault containing New Relic License Key.
- `log_sources_details`: List of log sources to be integrated with New Relic. Use stringified json of below structure:
   ```json
     [
       {
        "display_name": "logging-connector-1",
        "description": "Service connector for logs from compartment A to New Relic",
        "log_sources": [
         {
          "compartment_id": "ocid1.tenancy.oc1..***",
          "log_group_id": "ocid1.loggroup.oc1.iad.***"
         }
        ]
       },
       {
        "display_name": "logging-connector-2",
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