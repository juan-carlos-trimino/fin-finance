/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
---------------------------------------------------------------------------------------------------
cert-manager can be used to obtain certificates from a CA using the ACME protocol. A certificate is
a namespaced resource that references an Issuer or ClusterIssuer and defines a desired X.509
certificate that will be renewed and kept up to date.
---------------------------------------------------------------------------------------------------
Define input variables to the module.
***/
variable app_name {
  type = string
}
variable namespace {
  type = string
}
variable issuer_name {
  type = string
}
variable certificate_name {
  type = string
}
variable dns_names {
  default = []
  type = list
}
variable secret_name {
  type = string
}

# Create a Let's Encrypt TLS Certificate for the domain and inject it into K8s secrets.
resource "kubernetes_manifest" "certificate" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind = "Certificate"
    metadata = {
      name = var.certificate_name
      namespace = var.namespace
      labels = {
        app = var.app_name
      }
    }
    spec = {
      isCA = null
      privateKey = {
        rotationPolicy = "Always"
        size = 4096
        algorithm = "RSA"
        encoding = "PKCS1"
      }
      dnsNames = var.dns_names  # Add subdomains.
      duration = "2160h0m0s"  # 90 days.
      renewBefore = "720h0m0s" # 30 days
      secretName = var.secret_name
      # The Certificate will be issued using the issuer named 'var.issuer_name' in the
      # 'var.namespace' namespace (the same namespace as the Certificate resource).
      issuerRef = {
        kind = "Issuer"
        name = var.issuer_name
        # This is optional since cert-manager will default to this value; however, if you are using
        # an external issuer, change this to that issuer group.
        group = "cert-manager.io"
      }
    }
  }
}
