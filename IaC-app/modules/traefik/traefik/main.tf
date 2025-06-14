/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
---------------------------------------------------------------------------------------------------
Traefik is an edge router or reverse proxy that works as an ingress controller (router in
OpenShift) for a Kubernetes cluster. It is responsible for ensuring that when traffic from a web
application hits the cluster, it will go to the right Service; furthermore, Traefik supports ACME
and Let's Encrypt thereby making it very easy to manage and assign an SSL/TLS certificate to a web
application.

Traefik will be deployed as the reverse proxy with Basic Authentication (credentials are base64
encoded but not encrypted, see RFC7617) over SSL/TLS (HTTPS). However, since you are using the
Traefik Community Edition, only a single instance of Traefik with Let's Encrypt is supported, and,
by extension, a single point of failure is introduced to the deployment. (Traefik Enterprise
Edition with Let's Encrypt supports High Availability (HA).)
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
variable api_auth_token {
  type = string
}
variable chart_name {
  type = string
  description = "Ingress Controller Helm chart name."
  default = "traefik"
}
variable chart_repo {
  type = string
  description = "Using the official Traefik helm chart (Ingress Controller)."
  default = "https://helm.traefik.io/traefik"
}
variable chart_version {
  type = string
  # To use the latest version, go to https://artifacthub.io/ and type "traefik" on the edit box.
  description = "Ingress Controller Helm repository version."
}
variable timeout {
  type = number
}

# The null_resource can be used to create Terraform resources that don't have any particular
# resource type, and the local-exec is used to invoke commands on the local computer. The timestamp
# forces the commands to always be invoked.
# resource "null_resource" "scc-traefik" {
#   triggers = {
#     always_run = timestamp()
#   }
#   #
#   provisioner "local-exec" {
#     command = "kubectl apply -f ./modules/traefik/traefik/util/traefik-scc.yaml"
#   }
#   #
#   provisioner "local-exec" {
#     when = destroy
#     command = "kubectl delete scc fin-traefik-scc"
#   }
# }

# See 'env:' in ./modules/traefik/traefik/util/values.yaml.
resource "kubernetes_secret" "secret" {
  metadata {
    name = "${var.service_name}-provider-secret"
    namespace = var.namespace
    labels = {
      app = var.app_name
    }
  }
  # Plain-text data.
  data = {
    api_auth_token = var.api_auth_token
  }
  type = "Opaque"
}

# A ServiceAccount is used by an application running inside a pod to authenticate itself with the
# API server. A default ServiceAccount is automatically created for each namespace; each pod is
# associated with exactly one ServiceAccount, but multiple pods can use the same ServiceAccount. A
# pod can only use a ServiceAccount from the same namespace.
#
# Using the least-privilege approach with the namespace-scoped RoleBindings. In general, this is a
# preferred approach if a cluster's namespaces do not change dynamically and if Traefik is not
# required to watch all cluster namespaces.
resource "kubernetes_service_account" "service_account" {
  metadata {
    name = "${var.service_name}-service-account"
    namespace = var.namespace
    labels = {
      app = var.app_name
    }
    # This annotation indicates that pods running as this service account may only reference Secret
    # API objects specified in the service account's secrets field.
    annotations = {
      "kubernetes.io/enforce-mountable-secrets" = true
    }
  }
  #
  secret {
    name = kubernetes_secret.secret.metadata[0].name
  }
}

# Roles define WHAT can be done; role bindings define WHO can do it.
# The distinction between a Role/RoleBinding and a ClusterRole/ClusterRoleBinding is that the Role/
# RoleBinding is a namespaced resource; ClusterRole/ClusterRoleBinding is a cluster-level resource.
# A Role resource defines what actions can be taken on which resources; i.e., which types of HTTP
# requests can be performed on which RESTful resources.
resource "kubernetes_role" "role" {
  metadata {
    name = "${var.service_name}-role"
    namespace = var.namespace
    labels = {
      app = var.app_name
    }
  }
  rule {
    # Resources in the core apiGroup, which has no name - hence the "".
    api_groups = [""]
    verbs = ["get", "watch", "list"]
    # The plural form must be used when specifying resources.
    resources = ["services", "endpoints", "secrets"]
  }
  rule {
    api_groups = ["traefik.containo.us/v1alpha1"]
    verbs = ["get", "watch", "list"]
    resources = [
      "middlewares",
      "ingressroutes",
      "traefikservices",
      "ingressroutetcps",
      "ingressrouteudps",
      "tlsoptions",
      "tlsstores",
      "serverstransports"
    ]
  }
  rule {
    api_groups = ["extensions", "networking.k8s.io"]
    verbs = ["get", "watch", "list"]
    resources = ["ingresses", "ingressclasses"]
  }
  rule {
    api_groups = ["extensions", "networking.k8s.io"]
    verbs = ["update"]
    resources = ["ingresses/status"]
  }
  # This rule adds the custom SCC to the Role.
  # rule {
  #   # The resource SecurityContextConstraints (SCC) is associated with the API group
  #   # security.openshift.io. SCCs in OpenShift are a security feature that allows cluster
  #   # administrators to control permissions and access to security features for pods within a
  #   # cluster. SCCs are cluster-level resources requiring administrators to have cluster-admin
  #   # privileges to manage them. While SCCs were historically exposed under the core Kubernetes API
  #   # group, this is deprecated, and the recommended approach is to use the security.openshift.io
  #   # group for management.
  #   api_groups = ["security.openshift.io"]
  #   verbs = ["use"]
  #   resources = ["securitycontextconstraints"]
  #   resource_names = ["fin-traefik-scc"]
  # }
}

# Bind the role to the service account.
resource "kubernetes_role_binding" "role_binding" {
  metadata {
    name = "${var.service_name}-role-binding"
    namespace = var.namespace
    labels = {
      app = var.app_name
    }
  }
  # A RoleBinding always references a single Role, but it can bind the Role to multiple subjects.
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind = "Role"
    # This RoleBinding references the Role specified below...
    name = kubernetes_role.role.metadata[0].name
  }
  # ... and binds it to the specified ServiceAccount in the specified namespace.
  subject {
    # The default permissions for a ServiceAccount don't allow it to list or modify any resources.
    kind = "ServiceAccount"
    name = kubernetes_service_account.service_account.metadata[0].name
    namespace = kubernetes_service_account.service_account.metadata[0].namespace
  }
}

resource "helm_release" "traefik" {
  chart = var.chart_name
  repository = var.chart_repo
  version = var.chart_version
  namespace = var.namespace
  name = var.service_name
  values = [file("./modules/traefik/traefik/util/values.yaml")]
  # timeout = var.timeout
}
