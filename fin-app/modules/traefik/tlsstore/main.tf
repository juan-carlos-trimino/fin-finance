/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
-------------------------------------------------------
Define input variables to the module.
***/
variable app_name {
  type = string
}
variable namespace {
  type = string
}
variable secret_name {
  type = string
}
variable service_name {
  type = string
}

resource "kubernetes_manifest" "tlsstore" {
  manifest = {
    apiVersion = "traefik.containo.us/v1alpha1"
    kind = "TLSStore"
    metadata = {
      name = var.service_name
      namespace = var.namespace
      labels = {
        app = var.app_name
      }
    }
    spec = {
      defaultCertificate = {
        secretName = var.secret_name
      }
    }
  }
}
