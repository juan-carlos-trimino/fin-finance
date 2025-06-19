/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
---------------------------------------------------------------------------------------------------
An Ingress (entryway) exposes HTTP and HTTPS routes from outside the cluster to services within the
cluster. Traffic routing is controlled by rules defined on the Ingress resource. An Ingress
operates at the application layer of the network stack (for example HTTP) and can provide features
that a service cannot. But in order to make an Ingress resource work, an Ingress controller needs
to be running in the cluster. Different Kubernetes environments use different implementations of
the controller, but several do not provide a controller at all; e.g., OpenShift uses an Ingress
controller that is based on HAProxy. Here we are using Traefik.
---------------------------------------------------------------------------------------------------
Define input variables to the module.
***/
variable app_name {
  type = string
}
variable namespace {
  type = string
}
variable svc_finances {
  type = string
}
# variable svc_gateway {
#   type = string
# }
# variable svc_rabbitmq {
#   type = string
# }
# variable svc_kibana {
#   type = string
# }
variable middleware_rate_limit {
  type = string
}
variable middleware_compress {
  type = string
}
variable middleware_gateway_basic_auth {
  type = string
}
variable middleware_dashboard_basic_auth {
  type = string
}
# variable middleware_rabbitmq_basic_auth {
#   type = string
# }
# variable middleware_kibana_basic_auth {
#   type = string
# }
variable middleware_security_headers {
  type = string
}
variable tls_store {
  type = string
}
variable tls_options {
  type = string
}
variable secret_name {
  type = string
}
variable issuer_name {
  type = string
}
variable host_name {
  type = string
}
variable service_name {
  type = string
}

resource "kubernetes_manifest" "ingress-route" {
  manifest = {
    apiVersion = "traefik.io/v1alpha1"
    # This CRD is Traefik-specific.
    kind = "IngressRoute"
    metadata = {
      name = var.service_name
      namespace = var.namespace
      labels = {
        app = var.app_name
      }
    }
    #
    spec = {
      entryPoints = [  # Listening ports.
        "web",
        "websecure"
      ]
      routes = [
        {
          kind = "Rule"
          match = "(Host(`${var.host_name}`) || Host(`www.${var.host_name}`)) && (PathPrefix(`/dashboard`) || PathPrefix(`/api`))"
          priority = 40
          middlewares = [
            {
              name = var.middleware_dashboard_basic_auth
              namespace = var.namespace
            },
            {
              name = var.middleware_security_headers
              namespace = var.namespace
            }
          ]
          services = [
            {
              kind = "TraefikService"
              # If you enable the API, a new special service named api@internal is created and can
              # then be referenced in a router.
              name = "api@internal"
              port = 9000  # K8s service.
              # (default 1) A weight used by the weighted round-robin strategy (WRR).
              weight = 1
              # (default true) PassHostHeader controls whether to leave the request's Host Header
              # as it was before it reached the proxy, or whether to let the proxy set it to the
              # destination (backend) host.
              passHostHeader = true
              responseForwarding = {
                # (default 100ms) Interval between flushes of the buffered response body to the
                # client.
                flushInterval = "100ms"
              }
              strategy = "RoundRobin"
            }
          ]
        },
        {
          kind = "Rule"
          match = "Host(`${var.host_name}`, `www.${var.host_name}`) && PathPrefix(`/ping`)"
          priority = 40
          middlewares = [
            {
              name = var.middleware_dashboard_basic_auth
              namespace = var.namespace
            },
            {
              name = var.middleware_security_headers
              namespace = var.namespace
            }
          ]
          services = [
            {
              kind = "TraefikService"
              # If you enable the API, a new special service named api@internal is created and can
              # then be referenced in a router.
              name = "ping@internal"
              port = 9000  # K8s service.
              # (default 1) A weight used by the weighted round-robin strategy (WRR).
              weight = 1
              # (default true) PassHostHeader controls whether to leave the request's Host Header
              # as it was before it reached the proxy, or whether to let the proxy set it to the
              # destination (backend) host.
              passHostHeader = true
              responseForwarding = {
                # (default 100ms) Interval between flushes of the buffered response body to the
                # client.
                flushInterval = "100ms"
              }
              strategy = "RoundRobin"
            }
          ]
        },
        {
          kind = "Rule"
          # match = "Host(`169.46.98.220.nip.io`) && PathPrefix(`/`)"
          # match = "Host(`memories.mooo.com`) && (PathPrefix(`/`) || Path(`/upload`) || Path(`/api/upload`))"
          match = "(Host(`${var.host_name}`) || Host(`www.${var.host_name}`)) && PathPrefix(`/`)"
          # See https://doc.traefik.io/traefik/v2.0/routing/routers/#priority
          priority = 20
          # The rule is evaluated 'before' any middleware has the opportunity to work, and 'before'
          # the request is forwarded to the service.
          # Middlewares are applied in the same order as their declaration in router.
          middlewares = [
            {
              name = var.middleware_gateway_basic_auth
              namespace = var.namespace
            },
            {
              name = var.middleware_rate_limit
              namespace = var.namespace
            },
            {
              name = var.middleware_security_headers
              namespace = var.namespace
            }
          ]
          observability = {
            accesslogs = true
            metrics = true
            tracing = true
          }
          services = [
            {
              kind = "Service"
              name = var.svc_finances
              namespace = var.namespace
              port = 80  # K8s service.
              weight = 1
              passHostHeader = true
              responseForwarding = {
                flushInterval = "100ms"
              }
              strategy = "RoundRobin"
            }
          ]
        },

/***
        {
          kind = "Rule"
          match = "Host(`${var.host_name}`, `www.${var.host_name}`) && (Path(`/upload`) || Path(`/api/upload`))"
          priority = 50
          middlewares = [
            {
              name = var.middleware_rate_limit
              namespace = var.namespace
            },
            {
              name = var.middleware_security_headers
              namespace = var.namespace
            }
          ]
          services = [
            {
              kind = "Service"
              name = var.svc_gateway
              namespace = var.namespace
              port = 80  # K8s service.
              weight = 1
              passHostHeader = true
              responseForwarding = {
                flushInterval = "100ms"
              }
              strategy = "RoundRobin"
            }
          ]
        },
        {
          kind = "Rule"
          match = "Host(`${var.host_name}`, `www.${var.host_name}`) && (Path(`/video`) || Path(`/api/video`))"
          priority = 50
          middlewares = [
            {
              name = var.middleware_rate_limit
              namespace = var.namespace
            },
            {
              name = var.middleware_compress
              namespace = var.namespace
            },
            {
              name = var.middleware_security_headers
              namespace = var.namespace
            }
          ]
          services = [
            {
              kind = "Service"
              name = var.svc_gateway
              namespace = var.namespace
              port = 80  # K8s service.
              weight = 1
              passHostHeader = true
              responseForwarding = {
                flushInterval = "100ms"
              }
              strategy = "RoundRobin"
            }
          ]
        },
        {
          kind = "Rule"
          match = "Host(`${var.host_name}`, `www.${var.host_name}`) && Path(`/history`)"
          priority = 50
          middlewares = [
            {
              name = var.middleware_rate_limit
              namespace = var.namespace
            },
            {
              name = var.middleware_security_headers
              namespace = var.namespace
            }
          ]
          services = [
            {
              kind = "Service"
              name = var.svc_gateway
              namespace = var.namespace
              port = 80  # K8s service.
              weight = 1
              passHostHeader = true
              responseForwarding = {
                flushInterval = "100ms"
              }
              strategy = "RoundRobin"
            }
          ]
        },
        {
          kind = "Rule"
          # match = "Host(`169.46.98.220.nip.io`) && PathPrefix(`/`)"
          # match = "Host(`memories.mooo.com`) && (PathPrefix(`/`) || Path(`/upload`) || Path(`/api/upload`))"
          match = "Host(`${var.host_name}`, `www.${var.host_name}`) && PathPrefix(`/`)"
          # See https://doc.traefik.io/traefik/v2.0/routing/routers/#priority
          priority = 20
          # The rule is evaluated 'before' any middleware has the opportunity to work, and 'before'
          # the request is forwarded to the service.
          # Middlewares are applied in the same order as their declaration in router.
          middlewares = [
            {
              name = var.middleware_gateway_basic_auth
              namespace = var.namespace
            },
            {
              name = var.middleware_rate_limit
              namespace = var.namespace
            },
            {
              name = var.middleware_security_headers
              namespace = var.namespace
            }
          ]
          services = [
            {
              kind = "Service"
              name = var.svc_gateway
              namespace = var.namespace
              port = 80  # K8s service.
              weight = 1
              passHostHeader = true
              responseForwarding = {
                flushInterval = "100ms"
              }
              strategy = "RoundRobin"
            }
          ]
        },
        {
          kind = "Rule"
          match = "Host(`${var.host_name}`, `www.${var.host_name}`) && PathPrefix(`/kibana`)"
          priority = 40
          middlewares = [
            {
              name = var.middleware_kibana_basic_auth
              namespace = var.namespace
            },
            {
              name = var.middleware_security_headers
              namespace = var.namespace
            }
          ]
          services = [
            {
              kind = "Service"
              name = var.svc_kibana
              namespace = var.namespace
              port = 5601  # K8s service.
              # (default 1) A weight used by the weighted round-robin strategy (WRR).
              weight = 1
              # (default true) PassHostHeader controls whether to leave the request's Host Header
              # as it was before it reached the proxy, or whether to let the proxy set it to the
              # destination (backend) host.
              passHostHeader = true
              responseForwarding = {
                # (default 100ms) Interval between flushes of the buffered response body to the
                # client.
                flushInterval = "100ms"
              }
              strategy = "RoundRobin"
            }
          ]
        },
***/
        /*** Testing...
        {
          kind = "Rule"
          # match = "Host(`169.46.98.220.nip.io`) && PathPrefix(`/`)"
          # match = "Host(`memories.mooo.com`) && (PathPrefix(`/`) || Path(`/upload`) || Path(`/api/upload`))"
          match = "Host(`${var.host_name}`, `www.${var.host_name}`) && PathPrefix(`/fin`)"
          # See https://doc.traefik.io/traefik/v2.0/routing/routers/#priority
          priority = 20
          # The rule is evaluated 'before' any middleware has the opportunity to work, and 'before'
          # the request is forwarded to the service.
          # Middlewares are applied in the same order as their declaration in router.
          middlewares = [
            {
              # ???????????????????????
              name = var.middleware_gateway_basic_auth
              namespace = var.namespace
            },
            {
              name = var.middleware_rate_limit
              namespace = var.namespace
            },
            {
              name = var.middleware_security_headers
              namespace = var.namespace
            }
          ]
          services = [
            {
              kind = "Service"
              name = var.svc_finance
              namespace = var.namespace
              port = 80  # K8s service.
              weight = 1
              passHostHeader = true
              responseForwarding = {
                flushInterval = "100ms"
              }
              strategy = "RoundRobin"
            }
          ]
        },
        ***/ # Testing...

      ]
      tls = {
        certResolver = "le"
        domains = [
          {
            main = var.host_name
            sans = [  # URI Subject Alternative Names
              "www.${var.host_name}"
            ]
          }
        ]
        secretName = var.secret_name
        store = {
          name = var.tls_store
        }
        options = {
          name = var.tls_options
          namespace = var.namespace
        }
      }
    }
  }
}
