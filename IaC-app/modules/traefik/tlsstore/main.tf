/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
---------------------------------------------------------------------------------------------------
In Traefik, certificates are grouped together in certificates stores, and the TLS Store allows you
to configure the default TLS store. For more information, see
https://doc.traefik.io/traefik/routing/providers/kubernetes-crd/#kind-tlsstore and
https://doc.traefik.io/traefik/https/tls/#tls-options
---------------------------------------------------------------------------------------------------
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
    apiVersion = "traefik.io/v1alpha1"
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
