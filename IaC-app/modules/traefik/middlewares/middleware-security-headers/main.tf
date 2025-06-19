/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
-------------------------------------------------------
It manages the headers of requests and responses.
-------------------------------------------------
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
    apiVersion = "traefik.io/v1alpha1"
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
      headers = {
        # Set frameDeny to true to add the X-Frame-Options header with the value of DENY.
        frameDeny = true
        # The sslRedirect only allow HTTPS requests when set to true.
        # Deprecated in favor of EntryPoint redirection
        # (https://doc.traefik.io/traefik/routing/entrypoints/#redirection) or the RedirectScheme
        # middleware (https://doc.traefik.io/traefik/middlewares/http/redirectscheme/).
        # sslRedirect = true
        browserXssFilter = true
        # Set contentTypeNosniff to true to add the X-Content-Type-Options header with the value
        # nosniff.
        contentTypeNosniff = true
        # If the stsIncludeSubdomains is set to true, the includeSubDomains directive is appended
        # to the Strict-Transport-Security header.
        stsIncludeSubdomains = true
        # Set stsPreload to true to have the preload flag appended to the Strict-Transport-Security
        # header.
        stsPreload = true
        # The stsSeconds is the max-age of the Strict-Transport-Security header. If set to 0, the
        # header is not set.
        stsSeconds = 31536000  # 365 days
      }
    }
  }
}
