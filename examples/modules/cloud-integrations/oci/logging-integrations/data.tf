data "oci_core_route_tables" "default_vcn_route_table" {
  for_each = var.create_vcn ? { "default" = true } : {}

  compartment_id = var.compartment_ocid
  vcn_id         = module.vcn[0].vcn_id

  filter {
    name   = "display_name"
    values = ["Default Route Table for ${local.vcn_name}"]
    regex  = false
  }
}