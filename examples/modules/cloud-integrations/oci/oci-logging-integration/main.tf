locals {
  freeform_tags = {
    newrelic-terraform = "true"
  }
  connectors = jsondecode(data.external.connector_payload.result.connectors)
  connectors_map = {
    for conn in local.connectors : conn.display_name => conn
  }
}

# --- VCN Resources ---
resource "oci_core_vcn" "logging_vcn" {
  count          = var.create_vcn ? 1 : 0
  display_name   = "${var.newrelic_logging_prefix}-logging-vcn"
  compartment_id = var.compartment_ocid
  cidr_block     = "10.0.0.0/16"
  freeform_tags  = local.freeform_tags
}

# --- Gateway Resources ---
resource "oci_core_nat_gateway" "nat_gateway" {
  count          = var.create_vcn ? 1 : 0
  compartment_id = var.compartment_ocid
  vcn_id         = oci_core_vcn.logging_vcn[0].id
  display_name   = "${var.newrelic_logging_prefix}-logging-nat-gtw"
  freeform_tags  = local.freeform_tags
}

resource "oci_core_service_gateway" "service_gateway" {
  count          = var.create_vcn ? 1 : 0
  compartment_id = var.compartment_ocid
  vcn_id         = oci_core_vcn.logging_vcn[0].id
  display_name   = "${var.newrelic_logging_prefix}-logging-service-gtw"
  freeform_tags  = local.freeform_tags

  services {
    service_id   = data.oci_core_services.all_services.services[0].id
  }
}

# --- Route Table Resources ---
resource "oci_core_route_table" "private_route_table" {
  count          = var.create_vcn ? 1 : 0
  compartment_id = var.compartment_ocid
  vcn_id         = oci_core_vcn.logging_vcn[0].id
  display_name   = "${var.newrelic_logging_prefix}-logging-private-route-table"
  freeform_tags = local.freeform_tags

  route_rules {
    destination       = "0.0.0.0/0"
    destination_type  = "CIDR_BLOCK"
    network_entity_id = oci_core_nat_gateway.nat_gateway[0].id
  }

  route_rules {
    destination       = data.oci_core_services.all_services.services[0].cidr_block
    destination_type  = "SERVICE_CIDR_BLOCK"
    network_entity_id = oci_core_service_gateway.service_gateway[0].id
  }
}

# --- Security Lists ---
resource "oci_core_security_list" "private_subnet_security_list" {
  count          = var.create_vcn ? 1 : 0
  compartment_id = var.compartment_ocid
  vcn_id         = oci_core_vcn.logging_vcn[0].id
  display_name   = "private-subnet-security-list"
  freeform_tags = local.freeform_tags

  # Ingress Rules
  ingress_security_rules {
    stateless   = false
    protocol    = "all"
    source      = oci_core_vcn.logging_vcn[0].cidr_block
    description = "Allow all traffic from within the VCN"
  }

  # Egress Rules
  egress_security_rules {
    stateless   = false
    protocol    = "all"
    destination = "0.0.0.0/0"
    description = "Allow outbound traffic to the internet via NAT Gateway"
  }

  egress_security_rules {
    stateless   = false
    protocol    = "all"
    destination = "all-service-cidr-in-oracle-services-network"
    description = "Allow outbound traffic to Oracle Services via Service Gateway"
  }
}

# --- Subnet Resources ---
resource "oci_core_subnet" "private_subnet" {
  count          = var.create_vcn ? 1 : 0
  compartment_id     = var.compartment_ocid
  vcn_id             = oci_core_vcn.logging_vcn[0].id
  cidr_block         = "10.0.2.0/24"
  display_name       = "${var.newrelic_logging_prefix}-logging-private-subnet"
  prohibit_public_ip_on_vnic = true
  route_table_id     = oci_core_route_table.private_route_table[0].id
  security_list_ids  = [oci_core_security_list.private_subnet_security_list[0]]
  freeform_tags = local.freeform_tags
}

# --- Function Resources ---
resource "oci_functions_application" "logging_function_app" {
  compartment_id             = var.compartment_ocid
  config = {
    "VAULT_REGION"           = var.newrelic_region
    "DEBUG_ENABLED"          = var.debug_enabled
  }
  display_name               = "${var.newrelic_logging_prefix}-logging-function-app"
  freeform_tags              = local.freeform_tags
  shape                      = "GENERIC_X86"
  subnet_ids                 = [var.create_vcn ? oci_core_subnet.private_subnet[0].id : var.function_subnet_id]
}

resource "oci_functions_function" "logging_function" {
  application_id = oci_functions_application.logging_function_app.id
  display_name   = "${oci_functions_application.logging_function_app.display_name}-logging-function"
  memory_in_mbs  = "256"
  freeform_tags = local.freeform_tags
  image         = "${var.newrelic_region}.ocir.io/idms1yfytybe/oci-testing-registry/oci-function-x86:0.0.1" #TODO to change the actual function name
  provisioned_concurrency_config {
    strategy = "CONSTANT"
    count = 20
  }
}

# --- Connector Hub Resources ---
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
        log_id         = log_sources.value.log_id
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


output "vcn_network_details" {
  depends_on  = [oci_core_vcn.logging_vcn]
  description = "Output of the created network infra"
  value = var.create_vcn > 0 ? {
    vcn_id             = oci_core_vcn.logging_vcn[0].id
    nat_gateway_id     = oci_core_nat_gateway.nat_gateway[0].id
    nat_route_id       = oci_core_route_table.private_route_table[0].id
    service_gateway_id = oci_core_service_gateway.service_gateway[0].id
    sgw_route_id       = oci_core_route_table.private_route_table[0].id
    subnet_id          = oci_core_subnet.private_subnet[0].id
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