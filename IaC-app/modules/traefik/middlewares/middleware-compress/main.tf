/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
----------------------------------------------------------------------------------------
It compresses responses before sending them to the client; it uses the gzip compression.
----------------------------------------------------------------------------------------
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
variable minResponseBodyBytes {
  type = number
  default = 2048
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
      compress = {
        # Responses smaller than the specified values will not be compressed.
        minResponseBodyBytes = var.minResponseBodyBytes
      }
    }
  }
}
