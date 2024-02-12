# Virtual Cloud Network (VCN) or Virtual Private Cloud (VPC).
module "vcn" {
  source = "oracle-terraform-modules/vcn/oci"
  version = "3.1.0"
  vcn_name = "fin-vcn"
  # The DNS Domain Name for your virtual cloud network is: <your-dns-label>.oraclevcn.com
  # Alphanumeric string that begins with a letter.
  vcn_dns_label = "findnslbl"
  vcn_cidrs = ["10.0.0.0/16"]
  # compartment_id = var.compartment_id
  compartment_id = oci_identity_compartment.tf-compartment.id
  region = var.region
  internet_gateway_route_rules = null
  local_peering_gateways = null
  nat_gateway_route_rules = null
  create_internet_gateway = true
  create_nat_gateway = true
  create_service_gateway = true
}

# Source from https://registry.terraform.io/providers/oracle/oci/latest/docs/resources/core_security_list
# It'll allow traffic to go out anywhere – this will be needed for the Kubernetes nodes to download updates – and it allows traffic for all ports within the range of the VCN – 10.0.0.0/16.

resource "oci_core_security_list" "private-security-list" {
  display_name = "fin-security-list-for-private-subnet"
  compartment_id = oci_identity_compartment.tf-compartment.id
  # For vcn_id, use the OCID of the basic virtual network. To assign the OCID, before knowing it, assign an
  # output from the module, as input for the security list resource:
  # * Get the module's output attribute from the module's Outputs page.
  # * Assign a value to the resource argument with the expression:
  #   * <resource argument> = module.<module-name>.<output-attribute>
  #   * Example: vcn_id = module.vcn.vcn_id
  #   * Both the oci_core_security_list resource and the oracle-terraform-modules/vcn use the same argument
  #     name for virtual cloud network OCID: vcn_id.
  #   * The vcn_id on the left side of the equation is the argument (required input) for the resource.
  #   * The vcn_id on the right side of the equation is the OCID of the VCN that you create with the module.
  #   * It doesn't matter if you have run the VCN module script and created the VCN or not. Either way,
  #     Terraform assigns the VCN OCID to the security list after the VCN module is created.
  vcn_id = module.vcn.vcn_id
  # Add an egress rule to the security list based on the following values:
  # Note: The Allows field in the table is automatically generated based on other fields. You don't add an
  #       argument for it in your script.
  egress_security_rules {
    stateless = false  # No
    destination = "0.0.0.0/0"
    destination_type = "CIDR_BLOCK"
    protocol = "all"  # All protocols
  }
  #
  ingress_security_rules {
    stateless = false
    source = "10.0.0.0/16"
    source_type = "CIDR_BLOCK"
    protocol = "all"
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
resource "oci_core_subnet" "vcn-private-subnet" {
  compartment_id = oci_identity_compartment.tf-compartment.id
  vcn_id = module.vcn.vcn_id
  cidr_block = "10.0.1.0/24"
  # Caution: For the route table id, use module.vcn.nat_route_id.
  # Do not use module.vcn.nat_gateway_id, because it is the OCID for the gateway and not the route table.
  route_table_id = module.vcn.nat_route_id
  security_list_ids = [
    oci_core_security_list.private-security-list.id
  ]
  display_name = "fin-private-subnet"
  # prohibit_public_ip_on_vnic = true
}

# Source from https://registry.terraform.io/providers/oracle/oci/latest/docs/resources/core_security_list

# This rule will again allow all traffic to go out to the internet and allow VCN traffic to come in as well as traffic from anywhere on port 6443 TCP. This is going to be important because we’ll use kubectl to manipulate the Kubernetes cluster.
resource "oci_core_security_list" "public-security-list" {
  display_name = "fin-security-list-for-public-subnet"
  compartment_id = oci_identity_compartment.tf-compartment.id
  # For vcn_id, use the OCID of the basic virtual network. To assign the OCID, before knowing it, assign an
  # output from the module, as input for the security list resource:
  # * Get the module's output attribute from the module's Outputs page.
  # * Assign a value to the resource argument with the expression:
  #   * <resource argument> = module.<module-name>.<output-attribute>
  #   * Example: vcn_id = module.vcn.vcn_id
  #   * Both the oci_core_security_list resource and the oracle-terraform-modules/vcn use the same argument
  #     name for virtual cloud network OCID: vcn_id.
  #   * The vcn_id on the left side of the equation is the argument (required input) for the resource.
  #   * The vcn_id on the right side of the equation is the OCID of the VCN that you create with the module.
  #   * It doesn't matter if you have run the VCN module script and created the VCN or not. Either way,
  #     Terraform assigns the VCN OCID to the security list after the VCN module is created.
  vcn_id = module.vcn.vcn_id
  # Add an egress rule to the security list based on the following values:
  # Note: The Allows field in the table is automatically generated based on other fields. You don't add an
  #       argument for it in your script.
  egress_security_rules {
    stateless = false  # No
    destination = "0.0.0.0/0"
    destination_type = "CIDR_BLOCK"
    protocol = "all"  # All protocols
  }
  #
  ingress_security_rules {
    stateless = false
    source = "10.0.0.0/16"
    source_type = "CIDR_BLOCK"
    protocol = "all"
  }
  #
  ingress_security_rules {
    stateless = false
    source = "0.0.0.0/0"
    source_type = "CIDR_BLOCK"
    # Get protocol numbers from https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml TCP is 6
    protocol = "6"
    tcp_options {
      min = 6443
      max = 6443
    }
  }
  #
  # ingress_security_rules {
  #   stateless = false
  #   source = "0.0.0.0/0"
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
resource "oci_core_subnet" "vcn-public-subnet" {
  compartment_id = oci_identity_compartment.tf-compartment.id
  vcn_id = module.vcn.vcn_id
  cidr_block = "10.0.0.0/24"
  route_table_id = module.vcn.ig_route_id
  security_list_ids = [
    oci_core_security_list.public-security-list.id
  ]
  display_name = "fin-public-subnet"
}
