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

variable empty_dir {
  type = bool
  description = "Deployment type: empty-dir."
  default = false
}

variable persistent_disk {
  type = bool
  description = "Deployment type: persistent disk."
  default = false
}

variable pprof {
  type = bool
  description = "Enable/disable profiling (pprof)."
  default = false
}

variable build_image {
  type = bool
  default = true
}

variable cluster_domain_suffix {
  type = string
  description = "FQDN: service-name.namespace.svc.cluster.local"
  default = ".svc.cluster.local"
}

/***
The limitations of the kubernetes_manifest resource
---------------------------------------------------
If you want to create arbitrary Kubernetes resources in a cluster using Terraform, particularly
CRDs (Custom Resource Definitions), you can use the kubernetes_manifest resource from the
Kubernetes provider, but with these limitations:
(1) This resource requires API access during the planning time. This means the cluster has to be
    accessible at plan time and thus cannot be created in the same apply operation. That is, it
    is required to use two (2) separate Terraform apply steps: (1) Provision the cluster;
    (2) Create the resource.
(2) Any CRD (Custom Resource Definition) must already exist in the cluster during the planning
    phase. That is, it is required to use two (2) separate Terraform apply steps: (1) Install the
    CRDs; (2) Install the resources that are using the CRDs.
This are Terraform limitations, not specific to Kubernetes.
See https://github.com/hashicorp/terraform-provider-kubernetes/issues/1782.
***/
variable k8s_crds {  # CustomResourceDefinitions
  type = bool
  default = false
}

# Phase 1:
# reverse_proxy = true
# k8s_crds = true
# Phase 2:
# reverse_proxy = true
# k8s_crds = false
variable reverse_proxy {
  type = bool
  default = false
}

variable db_postgres {
  type = bool
  default = false
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
# Storage #
###########
# To view the Object Storage namespace string, do the following:
# Open the Profile menu and click Tenancy: <tenancy_name>. The namespace string is listed under
# Object Storage Settings.
variable obj_storage_ns {
  type = string
  sensitive = true
}

# From the top navigation bar, find your region.
# From the table in "Regions and Availability Domains," find your region's <region-identifier>.
# Example: us-ashburn-1
# https://docs.oracle.com/en-us/iaas/Content/General/Concepts/regions.htm
variable region {
  description = "An OCI region."
  type = string
  sensitive = true
}

# Profile menu -> User settings -> My groups -> Customer secret keys
variable aws_access_key_id {
  type = string
  sensitive = true
}

variable aws_secret_access_key {
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
# Helm chart deployment can sometimes take longer than the default 5 minutes.
variable "helm_traefik_timeout_seconds" {
  type = number
  default = 600  # 10 minutes
}

variable traefik_dashboard_username {
  default = "<required>"
  sensitive = true
  type = string
}

variable traefik_dashboard_password {
  default = "<required>"
  sensitive = true
  type = string
}

variable traefik_gateway_username {
  default = "<required>"
  sensitive = true
  type = string
}

variable traefik_gateway_password {
  default = "<required>"
  sensitive = true
  type = string
}

# Digital Ocean
# From the left menu:
# (1) Select API
# (2) From the Tokens tab, select the Generate New Token button.
#     Scopes: read/write.
# (3) The token has an expiration date; ensure the token is valid.
variable traefik_dns_api_token {
  default = "<required>"
  sensitive = true
  type = string
}

variable traefik_le_email {
  default = "<required>"
  sensitive = true
  type = string
}

###########
# busybox #
###########
variable busybox {
  default = "arm64v8/busybox:musl"
  type = string
}

##############
# PostgreSQL #
##############
variable postgres_image_tag {  # https://www.postgresql.org/docs/release/
  description = "PostgreSQL Docker images (https://hub.docker.com/_/postgres)."
  type = string
  default = "postgres:17.6-alpine3.22"
}

variable postgres_db_label {
  type = string
  default = "postgres"
}

variable postgres_port {
  description = "PGPORT behaves the same as the port connection parameter."
  type = number
  default = 5432
}

# https://hub.docker.com/_/postgres#postgres_db
variable postgres_db {
  description = "POSTGRES_DB: The default database."
  default = "<required>"
  sensitive = true
  type = string
}

# https://hub.docker.com/_/postgres#postgres_user
variable postgres_user {
  description = "POSTGRES_USER: It defines the default superuser."
  default = "<required>"
  sensitive = true
  type = string
}

# https://hub.docker.com/_/postgres#postgres_password
variable postgres_password {
  description = "POSTGRES_PASSWORD: It sets the superuser password for PostgreSQL."
  default = "<required>"
  sensitive = true
  type = string
}

variable postgres_replication_password {
  description = "REPLICATION_PASSWORD: Secret for PostgreSQL Authentication."
  default = "<required>"
  sensitive = true
  type = string
}

# https://hub.docker.com/_/postgres#pgdata
variable postgres_data {
  description = "PGDATA: It defines another location for the database files."
  default = "/wsf_data_dir/data/pgdata"
  type = string
}

variable general_script_path {
  default = "./utilities/general/scripts"
  type = string
}

variable postgres_config_path {
  default = "./utilities/postgres/configs"
  type = string
}

variable postgres_script_path {
  default = "./utilities/postgres/scripts"
  type = string
}

variable postgres_sql_path {
  default = "./utilities/postgres/sql"
  type = string
}
