variable "tenancy_ocid" {
  type        = string
  description = "OCI tenant OCID, more details can be found at https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/apisigningkey.htm#five.Do not modify."
}

variable "compartment_ocid" {
  type        = string
  description = "The OCID of the compartment where resources will be created. Do not modify."
}

variable "nr_prefix" {
  type        = string
  description = "The prefix for naming resources in this module."
  default     = "metrics"
}

variable "region" {
  type        = string
  description = "The name of the OCI region where these resources will be deployed."
}

variable "newrelic_endpoint" {
  type        = string
  default     = "US"
  description = "The endpoint to hit for sending the metrics. Varies by region [US|EU|JP]"
  validation {
    condition     = contains(["US", "EU", "JP"], var.newrelic_endpoint)
    error_message = "Valid values for var: newrelic_endpoint are (US, EU, JP)."
  }
}

variable "create_vcn" {
  type        = bool
  default     = true
  description = "Variable to create virtual network for the setup. True by default"
}

variable "function_subnet_id" {
  type        = string
  default     = ""
  description = "The OCID of the subnet to be used for the function app. If create_vcn is set to true, that will take precedence"
}

variable "connector_hubs_data" {
  type        = string
  description = "JSON-encoded string describing the connector hub layout. OCI caps each hub at 5 compartments and 50 namespaces total. For tenancies with multiple compartments or dense namespace selections, use the hub-calculator script (examples/modules/cloud-integrations/oci/hub-calculator) to compute this value automatically."
}

variable "ingest_api_secret_ocid" {
  type        = string
  description = "The OCID of the vault storing the ingest key for secure access."
}

variable "user_api_secret_ocid" {
  type        = string
  description = "The OCID of the vault storing the user key for secure access."
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

variable "user_ocid" {
  type        = string
  description = "OCID of the user Terraform authenticates as. Unless registry_username and registry_auth_token are set, the module uses this user's name to log in to Container Registry and creates an auth token for it to push the function image."
}

variable "function_image" {
  type        = string
  default     = "docker.io/newrelic/beyond-oci-metric-function:latest"
  description = "Public image the metrics function runs. The module copies it into a private Container Registry repository in this region and runs the function from that copy. Re-applying picks up a new image pushed under the same tag; use a version tag (for example :1.5) to choose when to upgrade."
}

variable "registry_username" {
  type        = string
  default     = ""
  description = "Container Registry username, without the tenancy namespace. Leave empty to use the name of user_ocid. Set it to the user that owns registry_auth_token when that is a different user, or to <identity_domain_name>/<username> for a user in a non-default identity domain."
}

variable "registry_auth_token" {
  type        = string
  default     = ""
  sensitive   = true
  description = "Existing auth token for pushing to Container Registry. Leave empty to have the module create one for user_ocid, which it keeps for later upgrades and deletes on destroy. Set it if that user already has two auth tokens (the OCI maximum), your organization doesn't allow creating tokens, or pushes should use a dedicated user."
}

variable "newrelic_account_id" {
  type        = string
  sensitive   = true
  description = "The New Relic account ID for sending metrics to New Relic endpoints"
}

variable "provider_account_id" {
  type        = string
  sensitive   = true
  description = "The Provider Account ID that has been linked with New Relic"
}
