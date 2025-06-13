/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
---------------------------------------------------------------------------------------------------
It uses basic authentication with SSL/TLS for the dashboard.
HTTP Basic Authentication is a simple challenge/response mechanism with which a server can request
authentication information (a user ID and password) from a client. The client passes the
authentication information to the server in an Authorization header. The authentication information
is base-64 encoded; it is not encrypted.

This scheme can be considered secure only when the connection between the web client and the server
is secure. If the connection is insecure, the scheme does not provide sufficient security to
prevent unauthorized users from discovering the authentication information for a server. If you
think that the authentication information might be intercepted, use basic authentication with
SSL/TLS encryption to protect the user ID and password.

Traefik supports passwords hashed with MD5, SHA1, or BCrypt. In the middlewares
middleware-dashboard-basic-auth and middleware-gateway-basic-auth, the secret resource uses the
bcrypt function, which is included as part of Terraform's built-in functions. You should keep in
mind that a bcrypt hash value includes a randomly selected salt, and, therefore, each call to this
function will return a different value, even if the given string and cost are the same.
---------------------------------------------------------------------------------------------------
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
