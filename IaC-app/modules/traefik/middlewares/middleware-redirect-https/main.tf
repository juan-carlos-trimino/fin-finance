/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
--------------------------------------------------------------------------------------------
It redirects the request if the request scheme is different from the configured scheme; this
implementation uses the EntryPoint redirection.
--------------------------------------------------------------------------------------------
Define input variables to the module.
***/
variable app_name {
  type = string
}
variable namespace {
  type = string
}
variable service_name {
  type = string
}

resource "kubernetes_manifest" "middleware" {
  manifest = {
    apiVersion = "traefik.containo.us/v1alpha1"
    kind = "Middleware"
    metadata = {
      name = var.service_name
      namespace = var.namespace
      labels = {
        app = var.app_name
      }
    }
    #
    spec = {
      redirectScheme = {
        scheme = "https"
        port = 443
        permanent = true
      }
    }
  }
}
