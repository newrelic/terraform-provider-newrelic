data "oci_resourcemanager_stacks" "current_stack" {
  compartment_id = var.compartment_ocid

  filter {
    name   = "display_name"
    values = [".*newrelic-logging-setup.*"]
    regex  = true
  }
}

data "oci_core_services" "all_services" {
  filter {
    name = "name"
    values = ["All services in ${var.region}"]
  }
}

data "external" "connector_payload" {
  program = ["python", "${path.module}/connector.py"]
  query = {
    "payload_link" = var.payload_link
  }
}