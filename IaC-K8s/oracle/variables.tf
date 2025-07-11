####################
# GLOBAL VARIABLES #
####################
variable app_name {
  type = string
  default = "finances"
  description = "The name of the application."
}

variable app_version {
  type = string
  default = "1.0.0"
  description = "The application version."
}

######################################
# CONFIDENTIAL/SENSITIVE INFORMATION #
#############################################################################################################
# Create a file with an extension of .auto.tfvars; e.g., tf_secrets.auto.tfvars. Next, add a line to this   #
# file for every variable below; e.g.,                                                                      #
# compartment_id = "xxx...xxx"                                                                              #
#                                                                                                           #
# IMPORTANT: Because this file contains confidential/sensitive information, do not push the file to a       #
#            version control system. This file is meant to be on your local system only.                    #
#############################################################################################################
# Compartments help you organize and control access to your resources. A compartment is a collection of     #
# related resources (such as cloud networks, compute instances, or block volumes) that can be accessed only #
# by those groups that have been given permission by an administrator in your organization.                 #
#############################################################################################################
# https://docs.oracle.com/en/cloud/foundation/cloud_architecture/governance/compartments.html#what-is-a-compartment
# In the top navigation bar, click the Profile menu, go to "Tenancy: <your-tenancy>"" and copy OCID.
variable tenancy_ocid {
  type = string
  sensitive = true
}

# Ensure "compartment_id" in outputs.tf has the correct name.
variable compartment_name {
  type = string
  sensitive = true
}

# From the Profile menu, go to "User settings" or "My profile" and copy OCID.
variable user_ocid {
  type = string
  sensitive = true
}

# From the Profile menu (user avatar), go to "My profile" and click "API Keys" (on the left side).
# Copy the fingerprint associated with the RSA public key. The format is:
# xx:xx:xx...xx
variable fingerprint {
  type = string
  sensitive = true
}

# From the top navigation bar, find your region.
# From the table in "Regions and Availability Domains," find your region's <region-identifier>.
# Example: us-ashburn-1
# https://docs.oracle.com/en-us/iaas/Content/General/Concepts/regions.htm
variable region {
  type = string
  sensitive = true
  description = "An OCI region."
}

# Collect the following information from your environment.
# Path to the RSA private key file to use for authentication against OCI. An API key can be created
# in the UI under Profile -> My Profile -> API keys.
# Example for Oracle Linux: /home/opc/.oci/<your-rsa-key-name>.pem
variable private_key_path {
  description = "Path to the RSA private key."
  type = string
  sensitive = true
}

variable kubeconfig_path {
  description = "Path to the kubectl config file."
  type = string
  sensitive = true
}

# OCI offers 4 OCPUs, 24GB RAM and 200GB of storage for free. These resources can be used to create
# up to 4 instances. Enter a value between [1,4] and the resources will be equally spread across
# the instance count.
variable nodes {
  description = "Count of nodes."
  type = number
  validation {
    condition = var.nodes >= 1 && var.nodes <= 4
    error_message = "Node count must be between 1 and 4."
  }
  sensitive = true
}

variable ocpus_per_node {
  description = "Ocpus per node."
  type = number
  # validation {
  #   condition = var.nodes >= 1 && var.nodes <= 4
  #   error_message = "Node count must be between 1 and 4."
  # }
  sensitive = true
}

variable memory_per_node {
  description = "Memory (in GB) per node."
  type = number
  # validation {
  #   condition = var.nodes >= 1 && var.nodes <= 4
  #   error_message = "Node count must be between 1 and 4."
  # }
  sensitive = true
}

# Specify the disk size in GB for the nodes in the cluster.
# variable boot_volume_size {
#   description = "Disk size in GB."
#   type = number
#   validation {
#     condition = var.boot_volume_size >= 10 && var.boot_volume_size <= 50
#     error_message = "Disk size Node must be between 10 and 50."
#   }
#   sensitive = true
# }

variable k8s_version {
  type = string
  default = "v1.33.1"
}

variable network_load_balancer {
  type = bool
  default = false
}

variable nlb_node_port {
  type = number
  default = 31600
}

variable load_balancer {
  type = bool
  default = false
}

variable nat {
  type = bool
  default = false
}

variable igw {  //Internet gateway.
  type = bool
  default = false
}
