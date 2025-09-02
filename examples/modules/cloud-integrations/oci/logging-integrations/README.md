# OCI Logging Integration Example

This Terraform Module links Oracle Logs and Log Groups to New Relic. It sets up the necessary infrastructure including Functions, Service Connector Hub, IAM policies, and Networking Components to enable logs collection from Oracle to New Relic.

## Prerequisites

- Ensure you have the [Terraform CLI](https://learn.hashicorp.com/tutorials/terraform/install-cli) installed.
- An Oracle Cloud Infrastructure (OCI) account with the necessary permissions to create resources.
- A New Relic account with License Key and User API Key.
- OCI CLI installed and configured in your local. If not, follow the [OCI CLI installation guide](https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm).

## Steps to Deploy the Terraform Script
1. Configure OCI CLI with your credentials:
   ```sh
   oci setup config
   ```
2. Clone the repository containing the Terraform example:
   ```sh
   git clone https://github.com/newrelic/terraform-provider-newrelic.git
   ```
3. Navigate to the directory containing the Terraform example:
   ```sh
   cd terraform-provider-newrelic/examples/modules/cloud-integrations/oci/logging-integrations
    ```
4. Create `connectors.json` file in the following format based on the number of log groups to be instrumented or based on preference. This file will be used to create the Service Connector Hub resources.
   ```json
   [
       {
           "display_name": "logging-connector-1",
           "description": "Service connector for logs from compartment A to New Relic",
           "log_sources": [
               {
                   "compartment_id": "ocid1.tenancy.oc1..****",
                   "log_group_id": "ocid1.loggroup.oc1.iad.****"
               }
           ]
       },
       {
           "display_name": "logging-connector-2",
           "description": "Service connector for logs from compartment A to New Relic",
           "log_sources": [
               {
                   "compartment_id": "ocid1.compartment.oc1..****",
                   "log_group_id": "ocid1.loggroup.oc1.iad.****"
               },
               {
                   "compartment_id": "ocid1.compartment.oc1..****",
                   "log_group_id": "ocid1.loggroup.oc1.iad.****"
               }
           ]
       }
   ]
   ```
5. Create a `terraform.tfvars` file in the same directory and add your variables.
6. Initialize the Terraform configuration. Review the Terraform plan to see the resources that will be created. Apply the Terraform configuration to create the resources.
   ```sh
   terraform init
   terraform plan
   terraform apply
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