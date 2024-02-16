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
variable subnet_id {
  type = string
}




# Source from https://registry.terraform.io/providers/oracle/oci/latest/docs/resources/containerengine_cluster

resource "oci_containerengine_cluster" "cluster" {
  name = var.name
  type = var.type
  compartment_id = var.compartment_id
  kubernetes_version = var.k8s_version
  vcn_id = var.vcn_id
  endpoint_config {
    is_public_ip_enabled = true
    # subnet_id = oci_core_subnet.public-subnet.id
    subnet_id = var.subnet_id
  }
  options {
    add_ons {
      is_kubernetes_dashboard_enabled = false
      is_tiller_enabled = false
    }
    # Note
    # The CIDR block for the pods must not overlap with the worker node and load balancer subnet CIDR blocks.
    #
    # The CIDR block for the Kubernetes service must not overlap with the VCN CIDR block.
    #
    #### The example code in this tutorial uses the same CIDR blocks as the Quick Create option in the Console.
    #
    # For more explanation, see
    # https://docs.oracle.com/en-us/iaas/Content/ContEng/Concepts/contengcidrblocks.htm
    kubernetes_network_config {
      pods_cidr = "10.244.0.0/16"
      services_cidr = "10.96.0.0/16"
    }
    service_lb_subnet_ids = [
      # module.public-subnet.oci_core_subnet.id
      var.subnet_id
    ]
  }
}

data "oci_containerengine_cluster_kube_config" "kubeconfig" {
  cluster_id = oci_containerengine_cluster.cluster.id
  endpoint = "PUBLIC_ENDPOINT"  # LEGACY_KUBERNETES,PUBLIC_ENDPOINT,PRIVATE_ENDPOINT,VCN_HOSTNAME.
  token_version = "2.0.0"
}

###########################
# Outputs for k8s cluster #
###########################
output "cluster-name" {
  value = oci_containerengine_cluster.cluster.name
}

output "cluster-OCID" {
  value = oci_containerengine_cluster.cluster.id
}

output "cluster-kubernetes-version" {
  value = oci_containerengine_cluster.cluster.kubernetes_version
}

output "cluster-state" {
  value = oci_containerengine_cluster.cluster.state
}

###################
# Outputs for k8s #
###################
output "cluster_id" {
  value = oci_containerengine_cluster.cluster.id
}

output "cluster_public_endpoint" {
  value = oci_containerengine_cluster.cluster.endpoints[0].public_endpoint
}

output "cluster_cni_type" {
  value = oci_containerengine_cluster.cluster.cluster_pod_network_options[0].cni_type
}
