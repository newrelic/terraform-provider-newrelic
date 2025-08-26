terraform {
  required_version = ">= 1.2.0"
  required_providers {
    oci = {
      source  = "oracle/oci"
      version = "7.12.0"
    }
  }
}

# Variables
provider "oci" {
  alias        = "home"
  tenancy_ocid = var.tenancy_ocid
  user_ocid    = data.oci_identity_user.current_user.user_id
  region       = var.region
}


locals {

  freeform_tags = {
    newrelic-terraform = "true"
  }
  # Names for the network infra
  vcn_name        = "${var.newrelic_logging_prefix}-logging-vcn"
  nat_gateway     = "${local.vcn_name}-natgateway"
  service_gateway = "${local.vcn_name}-servicegateway"
  subnet          = "${local.vcn_name}-public-subnet"

  connectors = jsondecode(data.external.connector_payload.result.connectors)
  connectors_map = {
    for conn in local.connectors : conn.display_name => conn
  }
}





#Resource for the logging function application
resource "oci_functions_application" "logging_function_app" {
  compartment_id = var.compartment_ocid
  config = {
    "VAULT_REGION"  = var.region
    "DEBUG_ENABLED" = var.debug_enabled
  }
  defined_tags               = {}
  display_name               = "${var.newrelic_logging_prefix}-logging-function-app"
  freeform_tags              = local.freeform_tags
  network_security_group_ids = []
  shape                      = "GENERIC_X86"
  subnet_ids = [
    data.oci_core_subnet.input_subnet.id,
  ]
}

# Resource for the function
resource "oci_functions_function" "logging_function" {
  depends_on = [oci_functions_application.logging_function_app]

  application_id = oci_functions_application.logging_function_app.id
  display_name   = "${oci_functions_application.logging_function_app.display_name}-logging-function"
  memory_in_mbs  = "256"

  defined_tags  = {}
  freeform_tags = local.freeform_tags
  image         = "${var.region}.ocir.io/idms1yfytybe/oci-testing-registry/oci-function-x86:0.0.1" #TODO to change the actual function name 
  provisioned_concurrency_config {
    strategy = "CONSTANT"
    count    = 20
  }
}


# Service Connector Hub - Routes logs from multiple log groups to New Relic function
resource "oci_sch_service_connector" "nr_logging_service_connector" {
  for_each = local.connectors_map

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

  depends_on = [oci_functions_function.logging_function]
}


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

data "oci_core_route_tables" "default_vcn_route_table" {
  depends_on     = [module.vcn] # Ensure VCN is created before attempting to find its route tables
  count          = var.create_vcn ? 1 : 0
  compartment_id = var.compartment_ocid
  vcn_id         = module.vcn[0].vcn_id

  filter {
    name   = "display_name"
    values = ["Default Route Table for ${local.vcn_name}"]
    regex  = false
  }
}

# Resource to manage the VCN's default route table and add your rule.
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

output "vcn_network_details" {
  depends_on  = [module.vcn]
  description = "Output of the created network infra"
  value = var.create_vcn && length(module.vcn) > 0 ? {
    vcn_id             = module.vcn[0].vcn_id
    nat_gateway_id     = module.vcn[0].nat_gateway_id
    nat_route_id       = module.vcn[0].nat_route_id
    service_gateway_id = module.vcn[0].service_gateway_id
    sgw_route_id       = module.vcn[0].sgw_route_id
    subnet_id          = module.vcn[0].subnet_id[local.subnet]
    } : {
    vcn_id             = ""
    nat_gateway_id     = ""
    nat_route_id       = ""
    service_gateway_id = ""
    sgw_route_id       = ""
    subnet_id          = var.function_subnet_id
  }
}

output "stack_id" {
  value = data.oci_resourcemanager_stacks.current_stack.stacks[0].id
}