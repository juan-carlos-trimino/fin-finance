/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
---------------------------------------------------------------------------------------------------
cert-manager is a powerful and extensible X.509 certificate management for Kubernetes. It will
obtain certificates from a variety of issuers; it ensures the certificates are valid and
up-to-date; and it will attempt to renew certificates at a configured time before expiry. On the
other hand, Traefik is capable of handling certificates in your cluster, but only when there is a
single pod of Traefik running. This, of course, is not acceptable because this pod becomes a single
point of failure in the infrastructure. To solve this issue, you’ll use cert-manager to request,
issue, renew, and store your certificates.
---------------------------------------------------------------------------------------------------
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
    name = "crds.enabled"
    value = true
  }
}
