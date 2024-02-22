####################
# GLOBAL VARIABLES #
####################
variable app_name {
  type = string
  description = "The name of the application."
  default = "finances"
}

variable app_version {
  type = string
  description = "The application version."
  default = "1.0.0"
}

######################################
# CONFIDENTIAL/SENSITIVE INFORMATION #
###################################################################################################
# Create a file with an extension of .auto.tfvars; e.g., tf_secrets.auto.tfvars. Next, add a line #
# to this file for every variable below; e.g.,                                                    #
# variable_name = "xxx...xxx"                                                                     #
#                                                                                                 #
# IMPORTANT: Because this file contains confidential/sensitive information, do not push the file  #
#            to a version control system. This file is meant to be on your local system only.     #
###################################################################################################
variable compartment_id {
  description = "Organize and control access to your resources."
  type = string
  sensitive = true
}

variable public_subnet_id {
  type = string
  sensitive = true
}

# Navigation menu->Developer Services->Containers & Artifacts [Under]->Kubernetes Clusters (OKE)->
# Compartment that contains the cluster->Clusters page->Resources [Under]->Node pools
variable node_pool_id {
  type = string
  sensitive = true
}

variable region {
  description = "An OCI region."
  type = string
  sensitive = true
}

variable cr_username {
  description = "Username for dockerhub."
  type = string
  sensitive = true
}

variable cr_password {
  description = "Password for dockerhub."
  type = string
  sensitive = true
}
