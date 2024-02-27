/***
Define input variables to the module.
***/
variable name {
  type = string
}
variable compartment_id {
  type = string
}
variable cluster_id {
  type = string
}
variable tenancy_ocid {
  type = string
}
variable cluster_cni_type {
  type = string
}
# See https://docs.oracle.com/en-us/iaas/Content/ContEng/Concepts/contengaboutk8sversions.htm for a
# list of supported versions.
variable k8s_version {
  type = string
  default = "v1.28.2"
}
variable nodes {
  type = number
}
variable memory_per_node {
  type = number
}
variable ocpus_per_node {
  type = number
}
variable subnet_id {
  type = string
}

###################################################################################################
# In Terraform to fetch data, you use a data source. Fetching data from a data source is similar  #
# to the GET method in REST APIs.                                                                 #
###################################################################################################
# Tenancy is the root or parent to all compartments.
# Use the value of <tenancy-ocid> for the compartment OCID.
# The data source gets a list of availability domains in your entire tenancy. The tenancy is the
# compartment OCID for the root compartment. Providing a specific "<compartment-ocid>" or the
# "<tenancy-ocid>" outputs the same list.
# https://registry.terraform.io/providers/oracle/oci/latest/docs/data-sources/identity_availability_domains
data "oci_identity_availability_domains" "avail_domains" {
  compartment_id = var.tenancy_ocid
}

locals {
  # The data source oci_identity_availability_domains, fetches a list of availability domains.
  ads = data.oci_identity_availability_domains.avail_domains.availability_domains[*].name
}

resource "tls_private_key" "node_pool_ssh_key_pair" {
  algorithm = "RSA"
  rsa_bits = 4096
}

data "oci_core_images" "use_image" {
  compartment_id = var.compartment_id
  operating_system = "Oracle Linux"
  operating_system_version = "8"
  shape = "VM.Standard.A1.Flex"
  sort_by = "TIMECREATED"
  sort_order = "DESC"
}

# https://registry.terraform.io/providers/oracle/oci/latest/docs/resources/containerengine_node_pool
resource "oci_containerengine_node_pool" "node-pool" {
  name = var.name
  cluster_id = var.cluster_id
  compartment_id = var.compartment_id
  kubernetes_version = var.k8s_version
  node_config_details {
    # node_pool_pod_network_option_details {
    #   cni_type = var.cluster_cni_type
    # }
    dynamic "placement_configs" {
      for_each = local.ads
      content {
        availability_domain = placement_configs.value
        subnet_id = var.subnet_id
      }
    }
    # The number of nodes that should be in the node pool.
    size = var.nodes
  }
  # Enhanced cluster feature.
  # node_pool_cycling_details {
  #   is_node_cycling_enabled = false
  #   maximum_surge = 1
  #   maximum_unavailable = 0
  # }
  # https://docs.oracle.com/en-us/iaas/Content/FreeTier/freetier_topic-Always_Free_Resources.htm
  node_shape = "VM.Standard.A1.Flex"
  # Configure how much memory and OCPUs to use in each node.
  node_shape_config {
    # The total amount of memory available to each node, in gigabytes.
    memory_in_gbs = var.memory_per_node
    # The total number of OCPUs available to each node in the node pool. See
    # https://docs.oracle.com/en-us/iaas/api/#/en/iaas/20160918/Shape/ for details.
    ocpus = var.ocpus_per_node
  }
  # Find the image OCID for YOUR REGION from https://docs.oracle.com/iaas/images/
  # Note: Since ARM instances are being used (see node_shape_config above), you will need to search
  #       for an ARM architecture compatible Oracle Linux images so search for the keyword aarch.
  node_source_details {
    boot_volume_size_in_gbs = 50
    image_id = data.oci_core_images.use_image.images.0.id
    source_type = "image"
  }
  initial_node_labels {
    key = "name"
    value = "k8s-cluster"
  }
  ssh_public_key = tls_private_key.node_pool_ssh_key_pair.public_key_openssh
}

#############################
# Outputs for public subnet #
#############################
output "node-pool-id" {
  value = oci_containerengine_node_pool.node-pool.id
}
