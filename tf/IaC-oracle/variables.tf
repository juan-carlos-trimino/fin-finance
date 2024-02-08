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

# See https://docs.oracle.com/en-us/iaas/Content/ContEng/Concepts/contengaboutk8sversions.htm for a list of
# supported versions.
variable k8s_version {  # Container Engine for Kubernetes (OKE).
  type = string
  default = "v1.28.2"
}

# OCI offers 4 OCPUs, 24GB RAM and 200GB of storage for free. These resources can be used to create up to 4
# instances. Enter a value between [1,4] and the resources will be equally spread across the instance count.
variable node_count {
  description = "Count of nodes."
  type = number
  default = 2
  validation {
    condition = var.node_count >= 1 && var.node_count <= 4
    error_message = "Node count must be between 1 and 4."
  }
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
# In the top navigation bar, click the Profile menu, go to "Tenancy: <your-tenancy>"" and copy OCID.
variable tenancy_ocid {
  type = string
  sensitive = true
}

# From the Profile menu, go to "My profile" and copy OCID.
variable user_ocid {
  type = string
  sensitive = true
}

#############################################################################################################
# Compartments help you organize and control access to your resources. A compartment is a collection of     #
# related resources (such as cloud networks, compute instances, or block volumes) that can be accessed only #
# by those groups that have been given permission by an administrator in your organization.                 #
#############################################################################################################
# https://docs.oracle.com/en/cloud/foundation/cloud_architecture/governance/compartments.html#what-is-a-compartment
# variable compartment_id {
#   type = string
#   sensitive = true
#   description = "Organize and control access to your resources."
# }

variable compartment_name {
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
# Path to the RSA private key file to use for authentication against OCI. An API key can be created in the UI
# under Profile -> My Profile -> API keys.
# Example for Oracle Linux: /home/opc/.oci/<your-rsa-key-name>.pem
variable private_key_path {
  type = string
  sensitive = true
  description = "Path to the RSA private key."
}
