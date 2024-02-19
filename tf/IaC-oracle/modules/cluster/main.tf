/***
Define input variables to the module.
***/
variable name {
  type = string
}
variable type {
  type = string
}
variable compartment_id {
  type = string
}
# See https://docs.oracle.com/en-us/iaas/Content/ContEng/Concepts/contengaboutk8sversions.htm for a list of
# supported versions.
variable k8s_version {  # Container Engine for Kubernetes (OKE).
  type = string
  default = "v1.28.2"
}
variable vcn_id {
  type = string
}
variable subnet_ids {
  type = list
  default = []
}



# https://registry.terraform.io/providers/oracle/oci/latest/docs/resources/containerengine_cluster
resource "oci_containerengine_cluster" "k8s-cluster" {
  name = var.name
  type = var.type
  compartment_id = var.compartment_id
  kubernetes_version = var.k8s_version
  vcn_id = var.vcn_id
  endpoint_config {
    # Whether the cluster should be assigned a public IP address. Defaults to false. If set to true
    # on a private subnet, the cluster provisioning will fail.
    is_public_ip_enabled = true
    subnet_id = var.subnet_ids[0]
  }
  options {
    add_ons {
      is_kubernetes_dashboard_enabled = false
      is_tiller_enabled = false
    }
    # The CIDR block for the pods must not overlap with the worker node and load balancer subnet
    # CIDR blocks.
    #
    # The CIDR block for the Kubernetes service must not overlap with the VCN CIDR block.
    #
    # For more explanation, see
    # https://docs.oracle.com/en-us/iaas/Content/ContEng/Concepts/contengcidrblocks.htm
    kubernetes_network_config {
      pods_cidr = "10.244.0.0/16"
      services_cidr = "10.96.0.0/16"
    }
    service_lb_subnet_ids = [
      var.subnet_ids[0]
    ]
  }
}

###########################
# Outputs for k8s cluster #
###########################
output "cluster-id" {
  value = oci_containerengine_cluster.k8s-cluster.id
}

output "cluster-cni-type" {
  value = oci_containerengine_cluster.k8s-cluster.cluster_pod_network_options[0].cni_type
}
