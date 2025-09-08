---
layout: "newrelic"
page_title: "New Relic Terraform Provider Logging Integration example for Oracle Cloud Integration"
sidebar_current: "docs-newrelic-provider-cloud-logging-integrations-guide"
description: |-
  Use this guide to set up the New Relic Logging Integrations fully automated through Terraform.
---

## New Relic Terraform Provider Logging Integration example for Oracle Cloud Integration

This guide describes the process of enabling the New Relic Logging Integration for your Oracle Cloud fully automated through Terraform.

> **NOTE:** This module assumes you've already set up the Oracle provider with the correct credentials. If you haven't done so, you can find the instructions here: [OCI instructions](https://registry.terraform.io/providers/oracle/oci/latest/docs).

If you encounter issues or bugs, please [report those on Github repository](https://github.com/newrelic/terraform-provider-newrelic/issues/new/choose).

## Usage
This module creates several OCI resources including Functions, Service Connector Hub, IAM policies, and Networking Components to enable logs collection from Oracle to New Relic.

```hcl
module "newrelic-oci-logging-integrations" {
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