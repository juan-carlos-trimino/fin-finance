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

# The limitations of the kubernetes_manifest resource
# ---------------------------------------------------
# If you want to create arbitrary Kubernetes resources in a cluster using Terraform, particularly
# CRDs (Custom Resource Definitions), you can use the kubernetes_manifest resource from the
# Kubernetes provider, but with these limitations:
# (1) This resource requires API access during the planning time. This means the cluster has to be
#     accessible at plan time and thus cannot be created in the same apply operation. That is, it
#     is required to use two (2) separate Terraform apply steps: (1) Provision the cluster;
#     (2) Create the resource.
# (2) Any CRD (Custom Resource Definition) must already exist in the cluster during the planning
#     phase. That is, it is required to use two (2) separate Terraform apply steps: (1) Install the
#     CRDs; (2) Install the resources that are using the CRDs.
# This are Terraform limitations, not specific to Kubernetes.
variable k8s_manifest_crd {
  type = bool
  default = "true"
}

variable nlb_node_port {
  type = number
  default = 31600
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

##############
# Image repo #
##############
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

###########
# Traefik #
###########
variable reverse_proxy {
  description = "Use a reverse proxy."
  type = bool
  default = false
}

# Helm chart deployment can sometimes take longer than the default 5 minutes.
variable "helm_traefik_timeout_seconds" {
  type = number
  default = 600  # 10 minutes
}

variable traefik_dashboard_username {
  default = "<required>"
  sensitive = true
}

variable traefik_dashboard_password {
  default = "<required>"
  sensitive = true
}

variable traefik_gateway_username {
  default = "<required>"
  sensitive = true
}

variable traefik_gateway_password {
  default = "<required>"
  sensitive = true
}

# From the left menu:
# (1) Select API
# (2) From the Tokens tab, select the Generate New Token button.
variable traefik_dns_api_token {  # digital ocean - Valid for 90 days (2/29/2024).
  default = "<required>"
  sensitive = true
}

variable traefik_le_email {
  default = "<required>"
  sensitive = true
}
