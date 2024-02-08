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
