# Output the "list" of all availability domains.
# The data source oci_identity_availability_domains, fetches a list of availability domains.
# Declare an output block to print the fetched information.
output "all-availability-domains" {
  value = data.oci_identity_availability_domains.availability_domains.availability_domains
}

###########################
# Outputs for compartment #
###########################
output "compartment-name" {
  value = oci_identity_compartment.tf-compartment.name
}

output "compartment-OCID" {
  value = oci_identity_compartment.tf-compartment.id
}

##############################
# Outputs for the vcn module #
##############################
output "vcn_id" {
  description = "OCID of the VCN that is created"
  value = module.vcn.vcn_id
}

output "id-for-route-table-that-includes-the-internet-gateway" {
  description = "OCID of the internet-route table. This route table has an internet gateway to be used for public subnets"
  value = module.vcn.ig_route_id
}

output "nat-gateway-id" {
  description = "OCID for NAT gateway"
  value = module.vcn.nat_gateway_id
}

output "id-for-for-route-table-that-includes-the-nat-gateway" {
  description = "OCID of the nat-route table - This route table has a nat gateway to be used for private subnets. This route table also has a service gateway."
  value = module.vcn.nat_route_id
}

#####################################
# Outputs for private security list #
#####################################
output "private-security-list-name" {
  value = oci_core_security_list.private-security-list.display_name
}

output "private-security-list-OCID" {
  value = oci_core_security_list.private-security-list.id
}

####################################
# Outputs for public security list #
####################################
output "public-security-list-name" {
  value = oci_core_security_list.public-security-list.display_name
}

output "public-security-list-OCID" {
  value = oci_core_security_list.public-security-list.id
}

##############################
# Outputs for private subnet #
##############################
output "private-subnet-name" {
  value = oci_core_subnet.vcn-private-subnet.display_name
}

output "private-subnet-OCID" {
  value = oci_core_subnet.vcn-private-subnet.id
}

#############################
# Outputs for public subnet #
#############################
output "public-subnet-name" {
  value = oci_core_subnet.vcn-public-subnet.display_name
}

output "public-subnet-OCID" {
  value = oci_core_subnet.vcn-public-subnet.id
  # value = module.vcn-public-subnet.id
}

###########################
# Outputs for k8s cluster #
###########################
output "cluster-name" {
  value = oci_containerengine_cluster.k8s-cluster.name
}

output "cluster-OCID" {
  value = oci_containerengine_cluster.k8s-cluster.id
}

output "cluster-kubernetes-version" {
  value = oci_containerengine_cluster.k8s-cluster.kubernetes_version
}

output "cluster-state" {
  value = oci_containerengine_cluster.k8s-cluster.state
}

#############################
# Outputs for k8s node pool #
#############################
output "node-pool-name" {
  value = oci_containerengine_node_pool.k8s-node-pool.name
}

output "node-pool-OCID" {
  value = oci_containerengine_node_pool.k8s-node-pool.id
}

output "node-pool-kubernetes-version" {
  value = oci_containerengine_node_pool.k8s-node-pool.kubernetes_version
}

output "node-size" {
  value = oci_containerengine_node_pool.k8s-node-pool.node_config_details[0].size
}

output "node-shape" {
  value = oci_containerengine_node_pool.k8s-node-pool.node_shape
}

###################
# Outputs for k8s #
###################
output "cluster_id" {
  value = oci_containerengine_cluster.k8s-cluster.id
}

output "cluster_public_endpoint" {
  value = oci_containerengine_cluster.k8s-cluster.endpoints[0].public_endpoint
}




# output "k8s_services_subnet_id" {
#   value = oci_core_subnet.oke_subnets["oke-services"].id
# }

# output "node_ips" {
#   value = oci_containerengine_node_pool.oke_node_pool.nodes.*.public_ip
# }
