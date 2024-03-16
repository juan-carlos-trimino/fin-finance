/***
Define input variables to the module.
***/
variable name {
  type = string
}
variable compartment_id {
  type = string
}
variable enabled {
  type = bool
  default = false
}
variable vcn_id {
  type = string
}

resource "oci_core_internet_gateway" "igw" {
  display_name = var.name
  compartment_id = var.compartment_id
  # Whether the gateway is enabled upon creation.
  enabled = var.enabled
  vcn_id = var.vcn_id
}

output "igw_id" {
  value = oci_core_internet_gateway.igw.id
}

output "igw" {
  value = oci_core_internet_gateway.igw
}
