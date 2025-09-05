# OCI Integration Module

This Terraform module links an Oracle Cloud Infrastructure (OCI) account to New Relic. It sets up the necessary infrastructure including OCI Functions, Service Connector Hub, IAM policies, and networking components to enable real-time metric streaming from OCI to New Relic.

## Prerequisites

- An OCI account with administrative privileges.
- OCI API key pair generated and configured.
- A New Relic account with the necessary permissions to create and manage cloud integrations.

## How the Module Works

This module creates several OCI resources, including IAM dynamic groups, policies, KMS vault and key, OCI Functions application and function, Service Connector Hub, and VCN networking infrastructure to facilitate the integration between OCI and New Relic. The Service Connector Hub streams OCI monitoring metrics to an OCI Function, which then forwards the metrics to New Relic's metric API.

It then uses the following resources based on cloud integrations from the New Relic Terraform Provider (the tenant ID used with the New Relic Terraform Provider's cloud integrations resources comes from the OCI resources deployed by the module, as stated above):

- `newrelic_cloud_oci_link_account`: Links the OCI tenancy to New Relic.
- `newrelic_cloud_oci_integrations`: Configures OCI metadata and tags collection to send to New Relic.

## Usage

> **Note:** Using this module requires a minimum version of the New Relic Terraform Provider that supports OCI cloud integrations.

```hcl
module "newrelic-oci-integrations" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/cloud-integrations/oci"
  
  # OCI Configuration
  tenancy_ocid      = "ocid1.tenancy.oc1..your_tenancy_ocid"
  current_user_ocid = "ocid1.user.oc1..your_user_ocid"
  compartment_ocid  = "ocid1.compartment.oc1..your_compartment_ocid"
  region           = "us-ashburn-1"
  
  # OCI Authentication
  fingerprint      = "your:api:key:fingerprint"
  private_key_path = "~/.oci/oci_api_key.pem"
  
  # New Relic Configuration
  newrelic_ingest_api_key = "your_newrelic_ingest_api_key"
  newrelic_user_api_key   = "your_newrelic_user_api_key"
  newrelic_account_id     = "your_newrelic_account_id"
  
  # Function Image (adjust prefix based on region)
  function_image = "iad.ocir.io/idms1yfytybe/public-newrelic-repo:latest"
  
  # Optional: Customize monitored services
  metrics_namespaces = [
    "oci_compute",
    "oci_database",
    "oci_autonomous_database",
    "oci_objectstorage"
  ]
}
```

> **Note:** If you have cloned this repo and would like to deploy this configuration in `main.tf` in the `testing` folder, use the following path as the value of the `source` argument:

```hcl
  source = "../examples/modules/cloud-integrations/oci"
```

### A Note on 'Applying' the Module :warning:

When applying this module, please be aware that some resources have dependencies that require sequential creation. The IAM policies and dynamic groups must be created in the home region before other resources can reference them. The module handles these dependencies automatically, but initial deployment may take several minutes as resources are created in the correct order.

```sh
terraform apply
```

## Variables

### Required Variables

- `tenancy_ocid` - (Required) The OCI tenant OCID.
- `current_user_ocid` - (Required) The OCID of the current user executing the terraform script.
- `compartment_ocid` - (Required) The OCID of the compartment where resources will be created.
- `region` - (Required) The OCI region for resource deployment.
- `fingerprint` - (Required) The fingerprint of the OCI API public key.
- `newrelic_ingest_api_key` - (Required) The New Relic Ingest API key for sending metrics.
- `newrelic_user_api_key` - (Required) The New Relic User API key for account linking.
- `newrelic_account_id` - (Required) The New Relic account ID.
- `function_image` - (Required) The container image URL for the New Relic metrics function. **Note:** The image prefix must match your OCI region:
  - For Ashburn region (us-ashburn-1): `iad.ocir.io/idms1yfytybe/public-newrelic-repo:latest`
  - For Phoenix region (us-phoenix-1): `phx.ocir.io/idms1yfytybe/public-newrelic-repo:latest`

### Optional Variables

- `private_key_path` - (Optional) Path to the OCI API private key file. Default: `""`.
- `private_key` - (Optional) The OCI API private key content (alternative to private_key_path). Default: `""`.
- `config_file_profile` - (Optional) OCI config file profile to use. Default: `""`.
- `dynamic_group_name` - (Optional) Name of the dynamic group for service connector access. Default: `"newrelic-metrics-dynamic-group"`.
- `newrelic_metrics_policy` - (Optional) Name of the IAM policy for metrics. Default: `"newrelic-metrics-policy"`.
- `newrelic_endpoint` - (Optional) New Relic metric API endpoint. Default: `"https://metric-api.newrelic.com/metric/v1"`.
- `newrelic_function_app` - (Optional) Name of the function application. Default: `"newrelic-metrics-function-app"`.
- `connector_hub_name` - (Optional) Name of the service connector hub. Default: `"newrelic-metrics-connector-hub"`.
- `function_app_shape` - (Optional) Shape of the function application. Default: `"GENERIC_X86"`.
- `metrics_namespaces` - (Optional) List of OCI service namespaces to monitor. Default includes 25+ OCI services.
- `vcn_name` - (Optional) Name of the VCN for New Relic metrics infrastructure. Default: `"newrelic-metrics-vcn"`.
- `kms_vault_name` - (Optional) Display name of the KMS vault for storing New Relic secrets. Default: `"newrelic-vault"`.

## Architecture

The module creates the following infrastructure:

1. **KMS Vault & Key**: Securely stores the New Relic API key
2. **VCN & Networking**: Virtual Cloud Network with public subnet, gateways for function connectivity
3. **IAM Resources**: Dynamic group and policies for service permissions (created only in home region)
4. **OCI Functions**: Application and function to process and forward metrics
5. **Service Connector Hub**: Streams OCI monitoring metrics to the function
6. **New Relic Integration**: Links OCI account and configures monitoring in New Relic

## Outputs

The module provides outputs for key resource identifiers that can be used for monitoring and troubleshooting:

- Function application OCID
- Service connector hub OCID  
- New Relic linked account ID
- KMS vault and key OCIDs

## Troubleshooting

### Common Issues

1. **Authentication Errors**: Verify OCI API key and fingerprint are correct
2. **Permission Errors**: Ensure user has administrative privileges in the target compartment
3. **Function Deployment**: Check that the function image is accessible from the specified region and that the image prefix matches your region (iad.ocir.io for Ashburn, phx.ocir.io for Phoenix)
4. **Network Connectivity**: Verify VCN and subnet configuration allows outbound internet access

### Validation

After deployment, verify the integration by:

1. Checking OCI Console for created resources (Service Connector Hub, Functions, etc.)
2. Monitoring New Relic Infrastructure for incoming OCI metrics
3. Reviewing function logs in OCI Console for any errors
