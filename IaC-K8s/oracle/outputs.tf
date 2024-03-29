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

###################################################################################
# Output the IP address of the Network Load Balancer (NLB) (module node-port-nlb) #
###################################################################################
# The special [*] symbol iterates over all of the elements of the list given to its left and
# accesses from each one the attribute name given on its right.
output "nlb_public_ip" {
  # The module node-port-nlb may have zero or one instance.
  value = [for ip in module.node-port-nlb[*].node_port_nlb_public_ip : ip[*].ip_address]
}

##########################################################################
# Output the IP address of the Load Balancer (LB) (module load-balancer) #
##########################################################################
output "lb_public_ip" {
  # The module load-balancer may have zero or one instance.
  value = module.load-balancer[*].lb_public_ip
}

########################################
# Output Internet Gateway (module igw) #
########################################
output "igw_id" {
  # The module load-balancer may have zero or one instance.
  value = module.igw[*].igw_id
}

output "igw" {
  value = module.igw[*].igw
}
