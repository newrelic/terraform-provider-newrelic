locals {
  freeform_tags = {
    newrelic-terraform = "true"
  }

  # --- VCN Resource Names ---
  vcn_name        = "${var.newrelic_logging_prefix}-logging-vcn"
  nat_gateway     = "${local.vcn_name}-natgateway"
  service_gateway = "${local.vcn_name}-servicegateway"
  subnet          = "${local.vcn_name}-public-subnet"
  connector_name  = "${var.newrelic_logging_prefix}-logging-connector"
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
  subnet_ids                 = [data.oci_core_subnet.input_subnet.id]
}

# --- Function Resources ---
resource "oci_functions_function" "logging_function" {
  application_id  = oci_functions_application.logging_function_app.id
  display_name    = "${oci_functions_application.logging_function_app.display_name}-logging-function"
  memory_in_mbs   = "256"
  freeform_tags   = local.freeform_tags
  image           = "${var.region}.ocir.io/idfmbxeaoavl/testing-registry/oci-function-test:0.0.1" #TODO to change the actual function name
  provisioned_concurrency_config {
    strategy      = "CONSTANT"
    count         = 20
  }
}

# --- Service Connector Hub - Routes logs to New Relic function ---
resource "oci_sch_service_connector" "nr_logging_service_connector" {
  compartment_id = var.compartment_ocid
  display_name   = local.connector_name
  freeform_tags  = local.freeform_tags

  source {
    kind = "logging"
    log_sources {
      compartment_id = var.compartment_ocid
      log_group_id   = var.log_group_id
      log_id         = var.log_id
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
  create_internet_gateway       = true                       # Enable creation of Internet Gateway
  internet_gateway_display_name = "NRLoggingInternetGateway" # Name the Internet Gateway
}

# --- Route Table Resources ---
resource "oci_core_default_route_table" "default_internet_route" {
  manage_default_resource_id = data.oci_core_route_tables.default_vcn_route_table[0].route_tables[0].id
  depends_on = [
    module.vcn,
    data.oci_core_route_tables.default_vcn_route_table
  ]
  route_rules {
    destination       = "0.0.0.0/0"
    destination_type  = "CIDR_BLOCK"
    network_entity_id = module.vcn[0].internet_gateway_id # Reference the internet gateway created by the module
    description       = "Route to Internet Gateway for New Relic logging"
  }
}