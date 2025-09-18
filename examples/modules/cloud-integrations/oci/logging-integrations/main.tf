locals {
  freeform_tags = {
    newrelic-terraform = "true"
  }

  # --- VCN Resource Names ---
  vcn_name        = "${var.newrelic_logging_prefix}-${var.region}-logging-vcn"
  nat_gateway     = "${local.vcn_name}-natgateway"
  service_gateway = "${local.vcn_name}-servicegateway"
  subnet          = "${local.vcn_name}-public-subnet"
}

# --- Function App Resources ---
resource "oci_functions_application" "logging_function_app" {
  compartment_id = var.compartment_ocid
  config = {
    "VAULT_REGION"      = var.region
    "DEBUG_ENABLED"     = var.debug_enabled
    "NEW_RELIC_REGION"  = var.new_relic_region
    "SECRET_OCID"       = var.secret_ocid
    "CLIENT_TTL"        = 30
  }
  display_name               = "${var.newrelic_logging_prefix}-logging-function-app"
  freeform_tags              = local.freeform_tags
  shape                      = "GENERIC_X86"
  subnet_ids                 = [var.create_vcn ? module.vcn[0].subnet_id[local.subnet] : var.function_subnet_id]
}

# --- Function Resources ---
resource "oci_functions_function" "logging_function" {
  application_id  = oci_functions_application.logging_function_app.id
  display_name    = "${oci_functions_application.logging_function_app.display_name}-logging-function"
  memory_in_mbs   = "128"
  freeform_tags   = local.freeform_tags
  image           = "${var.region}.ocir.io/idfmbxeaoavl/testing-registry/oci-function-test:0.0.1" #TODO to change the actual function name
}

# --- Service Connector Hub - Routes logs to New Relic function ---
resource "oci_sch_service_connector" "nr_logging_service_connector" {
  for_each = {
    for connector in jsondecode(var.log_sources_details) : connector.display_name => connector
  }

  compartment_id = var.compartment_ocid
  display_name   = each.value.display_name
  description    = each.value.description
  freeform_tags  = local.freeform_tags

  source {
    kind = "logging"
    dynamic "log_sources" {
      for_each = each.value.log_sources
      content {
        compartment_id = log_sources.value.compartment_id
        log_group_id   = log_sources.value.log_group_id
      }
    }
  }

  target {
    kind              = "functions"
    batch_size_in_kbs = 100
    batch_time_in_sec = 60
    compartment_id    = var.compartment_ocid
    function_id       = oci_functions_function.logging_function.id
  }
}

# --- VCN Resources ---
module "vcn" {
  source                   = "oracle-terraform-modules/vcn/oci"
  version                  = "3.6.0"
  count                    = var.create_vcn ? 1 : 0
  compartment_id           = var.compartment_ocid
  defined_tags             = {}
  freeform_tags            = local.freeform_tags
  vcn_cidrs                = ["10.0.0.0/16"]
  vcn_dns_label            = "nrlogging"
  vcn_name                 = local.vcn_name
  lockdown_default_seclist = false
  subnets = {
    public = {
      cidr_block = "10.0.0.0/16"
      type       = "public"
      name       = local.subnet
    }
  }
  create_nat_gateway            = true
  nat_gateway_display_name      = local.nat_gateway
  create_service_gateway        = true
  service_gateway_display_name  = local.service_gateway
  create_internet_gateway       = true
  internet_gateway_display_name = "NRLoggingInternetGateway"
}

# --- Route Table Resources ---
resource "oci_core_default_route_table" "default_internet_route" {
  for_each = var.create_vcn ? { "default" = true } : {}

  manage_default_resource_id = data.oci_core_route_tables.default_vcn_route_table["default"].route_tables[0].id
  route_rules {
    destination       = "0.0.0.0/0"
    destination_type  = "CIDR_BLOCK"
    network_entity_id = module.vcn[0].internet_gateway_id
    description       = "Route to Internet Gateway for New Relic logging"
  }
}