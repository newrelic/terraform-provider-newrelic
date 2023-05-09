variable "newrelic_account_id" {
  type = string
}

variable "name" {
  type    = string
  default = "production"
}

variable "gcp_service_account_id" {
  type = string
}

variable "gcp_project_id" {
  type = string
}
