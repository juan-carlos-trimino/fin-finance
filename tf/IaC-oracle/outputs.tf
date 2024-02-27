#############################
# Output for the vcn module #
#############################
output "vcn_id" {
  description = "OCID of the VCN that is created"
  value = module.vcn.vcn_id
}

###########################################
# Output for compartment (compartment.tf) #
###########################################
output "compartment_id" {
  value = oci_identity_compartment.fin-compartment.id
}

#######################################################################
# Output for Node Port Network Load Balancer (NLB) (node-port-nlb.tf) #
#######################################################################
output "node_port_nlb_public_ip" {
  value = [for ip in oci_network_load_balancer_network_load_balancer.node-port-nlb.ip_addresses : ip if ip.is_public == true]
}
