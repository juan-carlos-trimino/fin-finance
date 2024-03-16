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
variable traefik_dashboard_username {
  type = string
  sensitive = true
}
variable traefik_dashboard_password {
  type = string
  sensitive = true
}

resource "kubernetes_secret" "secret" {
  metadata {
    name = "${var.service_name}-secret"
    namespace = var.namespace
    labels = {
      app = var.app_name
    }
  }
  # Plain-text data.
  data = {
    # The second argument is optional and will default to 10 if unspecified. Since a bcrypt hash
    # value includes a randomly selected salt, each call to this function will return a different
    # value, even if the given string and cost are the same.
    # Traefik supports passwords hashed with MD5, SHA1, or BCrypt.
    users = "${var.traefik_dashboard_username}:${bcrypt(var.traefik_dashboard_password, 10)}"
  }
  type = "Opaque"
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
      basicAuth = {
        headerField = "X-WebAuth-User"
        removeHeader = true
        # The users option is an array of authorized users. Each user will be declared using the
        # username:encoded-password format.
        secret = kubernetes_secret.secret.metadata[0].name
      }
    }
  }
}
