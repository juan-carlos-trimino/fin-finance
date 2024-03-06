/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
-------------------------------------------------------
Define input variables to the module.
***/
variable namespace {
  type = string
}
variable service_name {
  type = string
}
variable chart_version {
  description = "Cert Manager Helm version."
  type = string
}
variable chart_name {
  description = "Cert Manager Helm name."
  default = "cert-manager"
  type = string
}
variable chart_repo {
  description = "Cert Manager Helm repository name."
  default = "https://charts.jetstack.io"
  type = string
}

resource "helm_release" "cert_manager" {
  name = var.service_name
  repository = var.chart_repo
  chart = var.chart_name
  version = var.chart_version
  namespace = var.namespace
  create_namespace = false
  # To automatically install and manage the CRDs as part of your Helm release, you must add the
  # --set installCRDs=true flag to your Helm installation command.
  set {
    name = "installCRDs"
    value = true
  }
}
