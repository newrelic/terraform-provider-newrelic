variable "tenancy_ocid" {
  type        = string
  description = "OCI tenant OCID, more details can be found at https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/apisigningkey.htm#five"
}

variable "current_user_ocid" {
  type        = string
  description = "The OCID of the current user executing the terraform script. Do not modify."
}

variable "compartment_ocid" {
  type        = string
  description = "The OCID of the compartment where resources will be created."
}

variable "region" {
  type        = string
  description = "OCI Region as documented at https://docs.cloud.oracle.com/en-us/iaas/Content/General/Concepts/regions.htm"
}

variable "private_key_path" {
  type        = string
  description = "The path to the private key file used for OCI API authentication. Generate using: openssl genrsa -out ~/.oci/oci_api_key.pem 2048"
  default     = ""
}

variable "private_key" {
  type        = string
  sensitive   = true
  description = "The private key content for OCI API authentication (alternative to private_key_path). Use this if you want to pass the key content directly instead of a file path."
  default     = ""
}

variable "fingerprint" {
  type        = string
  description = "The fingerprint of the public key. Get this from OCI Console -> User Settings -> API Keys"
}

variable "config_file_profile" {
  type        = string
  description = "The profile to use from the OCI config file (~/.oci/config). Leave empty to use environment variables or direct authentication."
  default     = ""
}

variable "dynamic_group_name" {
  type        = string
  description = "The name of the dynamic group for giving access to service connector"
  default     = "newrelic-metrics-dynamic-group"
}

variable "newrelic_metrics_policy" {
  type        = string
  description = "The name of the policy for metrics"
  default     = "newrelic-metrics-policy"
}


variable "newrelic_endpoint" {
  type        = string
  default     = "https://metric-api.newrelic.com/metric/v1"
  description = "The endpoint to hit for sending the metrics. Varies by region [US|EU]"
  validation {
    condition     = contains(["https://metric-api.newrelic.com/metric/v1", "https://metric-api.newrelic.com/metric/v1"], var.newrelic_endpoint)
    error_message = "Valid values for var: newrelic_endpoint are (metric-api.newrelic.com, metric-api.eu.newrelic.com)."
  }
}

variable "newrelic_ingest_api_key" {
  type        = string
  sensitive   = true
  description = "The Ingest API key for sending metrics to New Relic endpoints"
}

variable "newrelic_user_api_key" {
  type        = string
  sensitive   = true
  description = "The User API key for Linking the OCI Account to the New Relic account"
}

variable "newrelic_account_id" {
  type        = string
  sensitive   = true
  description = "The New Relic account ID for sending metrics to New Relic endpoints"
}

variable "newrelic_function_app" {
  type        = string
  description = "The name of the function application"
  default     = "newrelic-metrics-function-app"
}


variable "connector_hub_name" {
  type        = string
  description = "The prefix for the name of all of the resources"
  default     = "newrelic-metrics-connector-hub"
}

variable "function_app_shape" {
  type        = string
  default     = "GENERIC_X86"
  description = "The shape of the function application. The docker image should be built accordingly. Use ARM if using Oracle Resource manager stack"
}

variable "function_image" {
  type        = string
  description = "The container image URL for the New Relic metrics function"
  default     = ""
}

variable "metrics_namespaces" {
  type        = list(string)
  description = "The list of namespaces to send metrics for, within their respective compartments. Remove any namespaces where metrics should not be sent."
  default = [
    "oci_apigateway",
    "oci_autonomous_database",
    "oci_blockstore",
    "oci_compute",
    "oci_compute_infrastructure_health",
    "oci_compute_instance_health",
    "oci_computeagent",
    "oci_database",
    "oci_database_cluster",
    "oci_faas",
    "oci_healthchecks",
    "oci_internet_gateway",
    "oci_lbaas",
    "oci_logging",
    "oci_nat_gateway",
    "oci_nlb",
    "oci_nlb_extended",
    "oci_nosql",
    "oci_objectstorage",
    "oci_oke",
    "oci_postgresql",
    "oci_service_connector_hub",
    "oci_service_gateway",
    "oci_vcn",
    "oci_vcnip",
    "oci_vmi_resource_utilization"
  ]
}
