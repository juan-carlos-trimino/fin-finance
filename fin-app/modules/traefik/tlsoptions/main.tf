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

resource "kubernetes_manifest" "tlsoptions" {
  manifest = {
    apiVersion = "traefik.containo.us/v1alpha1"
    kind = "TLSOption"
    metadata = {
      name = var.service_name
      namespace = var.namespace
      labels = {
        app = var.app_name
      }
    }
    spec = {
      minVersion = "VersionTLS12"
      maxVersion = "VersionTLS13"
      # Cipher suites defined for TLS 1.2 and below cannot be used in TLS 1.3, and vice versa.
      cipherSuites = [
        "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256",
        "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
        "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
        "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
        "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
        "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
        "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305",
        "TLS_AES_256_GCM_SHA384",
        "TLS_AES_128_GCM_SHA256",
        "TLS_CHACHA20_POLY1305_SHA256",
        "TLS_FALLBACK_SCSV"
      ]
      # List of the elliptic curves references that will be used in an ECDHE handshake, in
      # preference order.
      curvePreferences = [
        "CurveP521",
        "CurveP384"
      ]
      # With strict SNI checking enabled (true), Traefik won't allow connections from clients that
      # do not specify a server_name extension or don't match any certificate configured on the
      # tlsOption.
      sniStrict = true
    }
  }
}
