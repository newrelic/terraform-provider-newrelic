data "oci_core_subnet" "input_subnet" {
  depends_on = [module.vcn]
  subnet_id  = var.create_vcn ? module.vcn[0].subnet_id[local.subnet] : var.function_subnet_id
}

data "oci_resourcemanager_stacks" "current_stack" {
  compartment_id = var.compartment_ocid

  filter {
    name   = "display_name"
    values = [".*newrelic-logging-setup.*"]
    regex  = true
  }
}

data "external" "connector_payload" {
  program = ["python", "${path.module}/connector.py"]
  query = {
    "payload_link" = var.payload_link
  }
}

data "oci_core_route_tables" "default_vcn_route_table" {
  count          = var.create_vcn ? 1 : 0
  compartment_id = var.compartment_ocid
  vcn_id         = module.vcn[0].vcn_id

  filter {
    name   = "display_name"
    values = ["Default Route Table for ${local.vcn_name}"]
    regex  = false
  }
}