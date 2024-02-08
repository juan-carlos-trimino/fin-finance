#############################################################################################################
# In Terraform to fetch data, you use a data source. Fetching data from a data source is similar to the GET #
# method in REST APIs.                                                                                      #
#############################################################################################################
# Tenancy is the root or parent to all compartments.
# Use the value of <tenancy-ocid> for the compartment OCID.
# The data source gets a list of availability domains in your entire tenancy. The tenancy is the compartment
# OCID for the root compartment. Providing a specific "<compartment-ocid>" or the "<tenancy-ocid>" outputs
# the same list.
# https://registry.terraform.io/providers/oracle/oci/latest/docs/data-sources/identity_availability_domains
data "oci_identity_availability_domains" "avail_domains" {
  compartment_id = var.tenancy_ocid
}

locals {
  # Collect the list of availability domains.
  ads = data.oci_identity_availability_domains.avail_domains.availability_domains[*].name
}

resource "tls_private_key" "node_pool_ssh_key_pair" {
  algorithm = "RSA"
  rsa_bits = 4096
}

# Source from https://registry.terraform.io/providers/oracle/oci/latest/docs/resources/containerengine_node_pool
resource "oci_containerengine_node_pool" "k8s-node-pool" {
  name = "worker-pool"
  cluster_id = oci_containerengine_cluster.k8s-cluster.id
  compartment_id = oci_identity_compartment.tf-compartment.id
  kubernetes_version = var.k8s_version
  node_config_details {
    node_pool_pod_network_option_details {
      cni_type = oci_containerengine_cluster.k8s-cluster.cluster_pod_network_options[0].cni_type
    }
    dynamic placement_configs {
      for_each = local.ads
      content {
        availability_domain = placement_configs.value
        subnet_id = oci_core_subnet.vcn-private-subnet.id
      }
    }
    # The number of nodes that should be in the node pool.
    size = var.node_count
  }
  # Enhanced cluster feature.
  node_pool_cycling_details {
    is_node_cycling_enabled = false
  }
  # An ARM instance from Oracle using the VM.Standard.A1.Flex shape.
  # https://docs.oracle.com/en-us/iaas/Content/FreeTier/freetier_topic-Always_Free_Resources.htm
  node_shape = "VM.Standard.A1.Flex"
  # node_shape = "VM.Standard.E2.1.Micro"  # Previous generation.

  # Specify the configuration of the shape to launch nodes in the node pool.

  # Configure how much memory and OCPUs to use in each node.

  # Since the free tier allows at max 4 ARM instances with an overall 24GB memory and 4 OCPUs,
  # I chose a 6 GB memory and 1 OCPU setup for each node, meaning that if I want to, I can provision 4
  # nodes at max for free within this node pool.
  node_shape_config {
    # The total amount of memory available to each node, in gigabytes.
    memory_in_gbs = 6
    # memory_in_gbs = floor(24 / var.node_count)
    # The total number of OCPUs available to each node in the node pool. See
    # https://docs.oracle.com/en-us/iaas/api/#/en/iaas/20160918/Shape/ for details.
    ocpus = 1
    # ocpus = floor(4 / var.node_count)
  }
  # Using image Oracle-Linux-7.9-aarch64-2023.12.08-0  Oracle-Linux-8.8-aarch64-2023.12.13-0
  # Find image OCID for YOUR REGION from https://docs.oracle.com/iaas/images/
  # Note: Since ARM instances are being used (see node_shape_config above), you will need to search for ARM
  #       architecture compatible Oracle Linux images so search for the keyword aarch.
  node_source_details {
    boot_volume_size_in_gbs = 50 #floor(200 / var.node_count)
    # image_id = "ocid1.image.oc1.us-chicago-1.aaaaaaaa6ywtssrjn35yao2upseif62n3adevgqjvznilsoxvjxhn5mrwwsq"
    image_id = "ocid1.image.oc1.us-chicago-1.aaaaaaaahfmrkrms6vvrkop7c4ymjzhq7ltea62nr6wt2i3u26jysfhzigja"
    source_type = "image"
  }
  initial_node_labels {
    key = "name"
    value = "k8s-cluster"
  }
  ssh_public_key = tls_private_key.node_pool_ssh_key_pair.public_key_openssh
}
