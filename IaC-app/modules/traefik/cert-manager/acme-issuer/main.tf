/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
---------------------------------------------------------------------------------------------------
The Automated Certificate Management Environment (ACME) protocol is a communications protocol for
automating interactions between certificate authorities and their users' web servers; it allows the
automated deployment of public key infrastructure at very low cost. It was designed by the Internet
Security Research Group (ISRG) for its Let's Encrypt service. The ACME Issuer type represents a
single account registered with the ACME Certificate Authority server. When you create a new ACME
Issuer, cert-manager will generate a private key, which is used to identify you with the ACME
server. Just keep in mind that since an Issuer is a namespaced resource, it is not possible to
issue certificates from an Issuer in a different namespace.
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
variable acme_email {
  type = string
}
variable acme_server {
  type = string
}
variable dns_names {
  default = []
  type = list
}
variable traefik_dns_api_token {
  sensitive = true
  type = string
}

resource "kubernetes_secret" "secret" {
  metadata {
    name = "${var.issuer_name}-api-token-secret"
    namespace = var.namespace
    labels = {
      app = var.app_name
    }
  }
  # Plain-text data.
  data = {
    access-token = var.traefik_dns_api_token
  }
  type = "Opaque"
}

resource "kubernetes_manifest" "acme-issuer" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind = "Issuer"
    metadata = {
      name = var.issuer_name
      namespace = var.namespace
      labels = {
        app = var.app_name
      }
    }
    spec = {
      acme = {
        # Email address used for ACME registration.
        email = var.acme_email
        # The ACME server URL; it will issue the certificates.
        server = var.acme_server
        # Name of the K8s secret use to store the ACME account private key.
        privateKeySecretRef = {
          name = "le-acme-private-key"
        }
        solvers = [
          # ACME DNS-01 provider configurations.
          {
            # (Optional) An empty 'selector' means that this solver matches all domains.
            # Only use digitalocean to solve challenges for trimino.xyz and www.trimino.xyz.
            selector = {
              dnsNames = var.dns_names
            }
            dns01 = {
              digitalocean = {
                tokenSecretRef = {
                  name = kubernetes_secret.secret.metadata[0].name
                  key = "access-token"
                }
              }
            }
          }
        ]
      }
    }
  }
}
