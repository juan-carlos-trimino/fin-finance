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
      errors = {
        status = [
          "400-599"
        ]
        query = "/{status}.html"
      }
    }
  }
}
