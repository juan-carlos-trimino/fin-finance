# Source from https://registry.terraform.io/providers/oracle/oci/latest/docs/resources/containerengine_cluster

resource "oci_containerengine_cluster" "k8s-cluster" {
  name = "k8s-cluster"
  type = "BASIC_CLUSTER"
  compartment_id = oci_identity_compartment.tf-compartment.id
  kubernetes_version = var.k8s_version
  vcn_id = module.vcn.vcn_id
  endpoint_config {
    is_public_ip_enabled = true
    subnet_id = oci_core_subnet.vcn-public-subnet.id
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
      oci_core_subnet.vcn-public-subnet.id
    ]
  }
}

data "oci_containerengine_cluster_kube_config" "kubeconfig" {
  cluster_id = oci_containerengine_cluster.k8s-cluster.id
  endpoint = "PUBLIC_ENDPOINT"  # LEGACY_KUBERNETES,PUBLIC_ENDPOINT,PRIVATE_ENDPOINT,VCN_HOSTNAME.
  token_version = "2.0.0"
}
