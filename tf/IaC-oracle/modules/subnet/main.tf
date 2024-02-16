/***
Define input variables to the module.
***/
variable subnet_display_name {
  type = string
}
variable cidr_block {
  type = string
}
variable compartment_id {
  type = string
}
variable vcn_id {
  type = string
}
variable route_table_id {
  type = string
}
#
variable sl_display_name {
  type = string
}
variable sl_egress_security_rules {
  default = []
  type = list(object({
    stateless = bool
    destination = string
    destination_type = string
    protocol = string
  }))
}
variable sl_ingress_security_rules {
  default = []
  type = list(object({
    stateless = bool
    source = string
    source_type = string
    protocol = string
    tcp_options = optional(list(object({
      min = number
      max = number
    })), [])
  }))
}




# Source from https://registry.terraform.io/providers/oracle/oci/latest/docs/resources/core_security_list
# It'll allow traffic to go out anywhere – this will be needed for the Kubernetes nodes to download updates – and it allows traffic for all ports within the range of the VCN – 10.0.0.0/16.

resource "oci_core_security_list" "security-list" {
  display_name = var.sl_display_name
  compartment_id = var.compartment_id
  vcn_id = var.vcn_id
  # Add an egress rule to the security list based on the following values:
  # Note: The Allows field in the table is automatically generated based on other fields. You don't add an
  #       argument for it in your script.
  dynamic "egress_security_rules" {
    for_each = var.sl_egress_security_rules
    iterator = it
    content {
      stateless = it.value["stateless"]
      destination = it.value["destination"]
      destination_type = it.value["destination_type"]
      protocol = it.value["protocol"]
    }
  }
  #
  dynamic "ingress_security_rules" {
    for_each = var.sl_ingress_security_rules
    iterator = it
    content {
      stateless = it.value["stateless"]
      source = it.value["source"]
      source_type = it.value["source_type"]
      # Get protocol numbers from https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml TCP is 6
      protocol = it.value["protocol"]
      dynamic "tcp_options" {
        for_each = it.value["tcp_options"]
        iterator = itn
        content {
          min = itn.value["min"]
          max = itn.value["max"]
        }
      }
    }
  }
  #
  # ingress_security_rules {
  #   stateless = false
  #   source = "10.0.0.0/16"
  #   source_type = "CIDR_BLOCK"
  #   # Get protocol numbers from https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml TCP is 6
  #   protocol = "6"
  #   tcp_options {
  #     min = 22
  #     max = 22
  #   }
  # }
  #
  # ingress_security_rules {
  #   stateless = false
  #   source = "0.0.0.0/0"
  #   source_type = "CIDR_BLOCK"
  #   # Get protocol numbers from https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml ICMP is 1
  #   protocol = "1"
  #   # For ICMP type and code see: https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml
  #   icmp_options {
  #     type = 3
  #     code = 4
  #   }
  # }
  #
  # ingress_security_rules {
  #   stateless = false
  #   source = "10.0.0.0/16"
  #   source_type = "CIDR_BLOCK"
  #   # Get protocol numbers from https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml ICMP is 1
  #   protocol = "1"
  #   # For ICMP type and code see: https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml
  #   icmp_options {
  #     type = 3
  #   }
  # }
}

# Source from https://registry.terraform.io/providers/oracle/oci/latest/docs/resources/core_subnet
resource "oci_core_subnet" "subnet" {
  display_name = var.subnet_display_name
  compartment_id = var.compartment_id
  vcn_id = var.vcn_id
  cidr_block = var.cidr_block
  # Caution: For the route table id, use module.vcn.nat_route_id.
  # Do not use module.vcn.nat_gateway_id, because it is the OCID for the gateway and not the route table.
  route_table_id = var.route_table_id
  security_list_ids = [
    oci_core_security_list.security-list.id
  ]
  # prohibit_public_ip_on_vnic = true
}



#####################################
# Outputs for private security list #
#####################################
output "security-list-name" {
  value = oci_core_security_list.security-list.display_name
}

output "security-list-OCID" {
  value = oci_core_security_list.security-list.id
}

####################################
# Outputs for public security list #
####################################
# output "public-security-list-name" {
#   value = oci_core_security_list.subnet.display_name
# }

# output "public-security-list-OCID" {
#   # value = oci_core_security_list.public-security-list.id
#   value = oci_core_security_list.subnet.id
# }





##############################
# Outputs for private subnet #
##############################
# output "private-subnet-name" {
#   value = oci_core_subnet.private-subnet.display_name
# }

# output "private-subnet-OCID" {
#   value = oci_core_subnet.private-subnet.id
# }

#############################
# Outputs for public subnet #
#############################
output "subnet-name" {
  value = oci_core_subnet.subnet.display_name
}

output "subnet-id" {
  value = oci_core_subnet.subnet.id
}
