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
# The maximum rate, by default in requests per second, allowed from a given source.
variable average {
  type = number
  default = 0  # No rate limiting.
}
# Rate = average / period.
variable period {
  type = string
  default = "1s"  # 1 second.
}
# The maximum number of requests allowed to go through in the same arbitrarily small period of
# time.
variable burst {
  type = number
  default = 1
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
      rateLimit = {
        average = var.average
        period = var.period
        burst = var.burst
      }
    }
  }
}
