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
# Create a file with an extension of .tfvars; e.g., tf_secrets.tfvars. Next, add a line to this file for    #
# every variable below; e.g.,
# compartment_id = "xxx...xxx"
#
# IMPORTANT: Because this file contains confidential/sensitive information, do not push the file to a       #
#            version control system. This file is meant to be on your local system only.                    #
#############################################################################################################
# Compartments help you organize and control access to your resources. A compartment is a collection of     #
# related resources (such as cloud networks, compute instances, or block volumes) that can be accessed only #
# by those groups that have been given permission by an administrator in your organization.                 #
#############################################################################################################
# https://docs.oracle.com/en/cloud/foundation/cloud_architecture/governance/compartments.html#what-is-a-compartment
variable compartment_id {
  type = string
  sensitive = true
  description = "Organize and control access to your resources."
}

# https://docs.oracle.com/en-us/iaas/Content/General/Concepts/regions.htm
variable region {
  type = string
  sensitive = true
  description = "xxxxxxxxxxThe region to provision the resources in"
}

variable ssh_public_key {
  type = string
  sensitive = true
  description = "xxxxxxxThe SSH public key to use for connecting to the worker nodes"
}
