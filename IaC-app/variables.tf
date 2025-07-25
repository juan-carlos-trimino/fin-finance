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

variable db_mysql {
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

################
# MySQL Server #
################
variable mysql_image_tag {
  description = "MySQL version (https://hub.docker.com/_/mysql)"
  type = string
  default = "mysql:9"
}

/***
MYSQL_DATABASE: This variable allows you to specify the name of a database to be created on image
startup. If a user name and a password are supplied with MYSQL_USER and MYSQL_PASSWORD, the user is
created and granted superuser access to this database (corresponding to GRANT ALL). The specified
database is created by a CREATE DATABASE IF NOT EXIST statement, so that the variable has no effect
if the database already exists.
***/
variable mysql_database {
  default = "mydb"
  type = string
}

/***
MYSQL_USER, MYSQL_PASSWORD: These variables are used in conjunction to create a user and set that
user's password, and the user is granted superuser permissions for the database specified by the
MYSQL_DATABASE variable. Both MYSQL_USER and MYSQL_PASSWORD are required for a user to be created
-- if any of the two variables is not set, the other is ignored. If both variables are set but
MYSQL_DATABASE is not, the user is created without any privileges.
***/
variable mysql_user {
  default = "<required>"
  sensitive = true
  type = string
}

variable mysql_password {
  default = "<required>"
  sensitive = true
  type = string
}

/***
MYSQL_ROOT_HOST: By default, MySQL creates the 'root'@'localhost' account. This account can only be
connected to from inside the container. To allow root connections from other hosts, set this
environment variable. For example, the value 172.17.0.1, which is the default Docker gateway IP,
allows connections from the host machine that runs the container. The option accepts only one
entry, but wildcards are allowed (for example, MYSQL_ROOT_HOST=172.*.*.* or MYSQL_ROOT_HOST=%).
***/
variable mysql_root_host {
  default = "%"
  type = string
}

# MYSQL_ROOT_PASSWORD: This variable specifies a password that is set for the MySQL root account.
variable mysql_root_password {
  default = "<required>"
  sensitive = true
  type = string
}

################
# MySQL Router #
################
variable mysql_router_image_tag {
  description = "MySQL Router version (https://hub.docker.com/r/mysql/mysql-router)"
  type = string
  default = "mysql:9"
}

variable mysql_router_host {
  description = "Required. MySQL host to connect to."
  type = string
  default = "%"  #???????
}

variable mysql_router_port {
  description = "Required. MySQL server listening port."
  type = string
  default = "3306"
}

variable mysql_router_user {
  description = "Required. MySQL user to connect with."
  type = string
  default = "ruser"
}

variable mysql_router_password {
  description = "Required. MySQL user's password."
  type = string
  default = "ruser"
}

variable mysql_router_cluster_members {
  description = "Required. Wait for this number of cluster instances to be online."
  type = number
  default = 3
}

variable mysql_router_bootstrap {
  description = "Optional. List of additional command line options to apply during bootstrapping."
  type = string
  default = "--directory /fin-router/router"
}
