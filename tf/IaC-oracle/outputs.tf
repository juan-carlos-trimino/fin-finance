##############################
# Outputs for the vcn module #
##############################
output "vcn_id" {
  description = "OCID of the VCN that is created"
  value = module.vcn.vcn_id
}

###########################
# Outputs for compartment #
###########################
output "compartment-id" {
  value = oci_identity_compartment.fin-compartment.id
}
