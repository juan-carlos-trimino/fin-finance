/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
-------------------------------------------------------
Define input variables to the module.
***/
variable app_name {
  type = string
}
variable app_version {
  type = string
}
variable image_tag {
  default = ""
  type = string
}
variable namespace {
  default = "default"
  type = string
}
# Be aware that the default imagePullPolicy depends on the image tag. If a container refers to the
# latest tag (either explicitly or by not specifying the tag at all), imagePullPolicy defaults to
# Always, but if the container refers to any other tag, the policy defaults to IfNotPresent.
#
# When using a tag other that latest, the imagePullPolicy property must be set if changes are made
# to an image without changing the tag. Better yet, always push changes to an image under a new
# tag.
variable image_pull_policy {
  default = "Always"
  type = string
}
variable security_context {
  default = [{
    run_as_non_root = false
    run_as_user = 0
    run_as_group = 0
    read_only_root_filesystem = false
  }]
  type = list(object({
    run_as_non_root = bool
    run_as_user = number
    run_as_group = number
    read_only_root_filesystem = bool
  }))
}
variable env {
  default = {}
  type = map(any)
}
# Quality of Service (QoS) classes for pods:
# (1) BestEffort (lowest priority) - It's assigned to pods that do not have any requests or limits
#     set at all (in any of their containers).
# (2)
# (3) Guaranteed (highest priority) - It's assigned to pods whose containers' requests are equal to
#     the limits for all resources (for each container in the pod). For a pod's class to be
#     Guaranteed, three things need to be true:
#     * Requests and limits need to be set for both CPU and memory.
#     * They need to be set for each container.
#     * They need to be equal; the limit needs to match the request for each resource in each
#       container.
# If a Container specifies its own memory limit, but does not specify a memory request, Kubernetes
# automatically assigns a memory request that matches the limit. Similarly, if a Container
# specifies its own CPU limit, but does not specify a CPU request, Kubernetes automatically assigns
# a CPU request that matches the limit.
variable resources {
  default = {}
  type = object({
    requests_cpu = optional(string)
    requests_memory = optional(string)
    limits_cpu = optional(string)
    limits_memory = optional(string)
  })
}
variable replicas {
  default = 1
  type = number
}
variable termination_grace_period_seconds {
  default = 30
  type = number
}
variable service_name {
  type = string
}
# The ServiceType allows to specify what kind of Service to use: ClusterIP (default),
# NodePort, LoadBalancer, and ExternalName.
variable service_type {
  default = "ClusterIP"
  type = string
}
variable secrets {
  default = []
  type = list(object({
    name = string
    annotations = optional(map(string), {})
    data = optional(map(string), {})
    binary_data = optional(map(string), {})
    type = optional(string, "Opaque")
  }))
  sensitive = true
}
variable service_account {
  default = null
  type = object({
    name = string
    namespace = string
    # Note: The keys and the values in the map must be strings. In other words, you cannot use
    #       numeric, boolean, list or other types for either the keys or the values.
    labels = optional(map(string), {})
    # https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
    annotations = optional(map(string), {})
    automount_service_account_token = optional(bool, true)
    secrets = optional(list(object({
      name = string
    })), [])
  })
}
variable role {
  default = null
  type = object({
    name = string
    namespace = string
    labels = optional(map(string), {})
    annotations = optional(map(string), {})
    rules = optional(list(object({
      api_groups = set(string)
      resources = set(string)
      resource_names = optional(set(string))
      verbs = set(string)
    })), [])
  })
}
variable role_binding {
  default = null
  type = object({
    # Name of the role binding, must be unique.
    name = string
    namespace = string
    labels = optional(map(string), {})
    annotations = optional(map(string), {})
    # A RoleBinding always references a single Role, but it can bind the Role to multiple subjects.
    # The Role to bind Subjects to.
    role_ref = object({
      kind = string
      # 'name' must match the name of the Role or ClusterRole you wish to bind to.
      name = string
      # The API group to drive authorization decisions. This value must be and defaults to
      # 'rbac.authorization.k8s.io'.
      api_group = string
    })
    # The Users, Groups, or ServiceAccounts to grand permissions to.
    # More than one 'subject' is allowed.
    subjects = list(object({
      # The type of binding to use. This value must be ServiceAccount, User or Group.
      kind = string
      # The name of this Role to bind Subjects to.
      # 'name' is case sensitive.
      name = string
      # Namespace defines the namespace of the ServiceAccount to bind to. This value only applies
      # to kind ServiceAccount.
      namespace = optional(string)
      # The API group to drive authorization decisions. This value only applies to kind User and
      # Group. It must be 'rbac.authorization.k8s.io'.
      api_group = optional(string)
    }))
  })
}
# The service normally forwards each connection to a randomly selected backing pod. To ensure that
# connections from a particular client are passed to the same Pod each time, set the service's
# sessionAffinity property to ClientIP instead of None (default).
# Session affinity and Web Browsers (for LoadBalancer Services)
# Since the service is now exposed externally, accessing it with a web browser will hit the same
# pod every time. If the sessionAffinity is set to None, then why? The browser is using keep-alive
# connections and sends all its requests through a single connection. Services work at the
# connection level, and when a connection to a service is initially open, a random pod is selected
# and then all network packets belonging to that connection are sent to that single pod. Even with
# the sessionAffinity set to None, the same pod will always get hit (until the connection is
# closed).
variable service_session_affinity {
  default = "None"
  type = string
}
variable ports {
  default = [{
    name = "ports"
    service_port = 80
    target_port = 8080
    protocol = "TCP"
  }]
  type = list(object({
    name = string
    service_port = number
    target_port = number
    node_port = optional(number)
    protocol = string
  }))
}


########################
# To relax the StatefulSet ordering guarantee while preserving its uniqueness and identity
# guarantee.
variable pod_management_policy {
  default = "OrderedReady"
}
# The primary use case for setting this field is to use a StatefulSet's Headless Service to
# propagate SRV records for its Pods without respect to their readiness for purpose of peer
# discovery.
variable publish_not_ready_addresses {
  default = "false"
  type = bool
}



variable pvc_access_modes {
  default = []
  type = list(any)
}
variable pvc_storage_class_name {
  default = ""
}
variable pvc_storage_size {
  default = "20Gi"
}

locals {
  svc_name = "${var.service_name}-headless"
  pod_selector_label = "ps-${var.service_name}"
  svc_selector_label = "svc-${local.svc_name}"
  rmq_label = "mem-rmq-cluster"
}

resource "null_resource" "scc-rabbitmq" {
  triggers = {
    always_run = timestamp()
  }
  provisioner "local-exec" {
    command = "oc apply -f ${var.path_rabbitmq_files}/mem-rabbitmq-scc.yaml"
  }
  #
  provisioner "local-exec" {
    when = destroy
    command = "oc delete scc mem-rabbitmq-scc"
  }
}

# The maximum size of a Secret is limited to 1MB.
# K8s helps keep Secrets safe by making sure each Secret is only distributed to the nodes that run
# the pods that need access to the Secret.
# On the nodes, Secrets are always stored in memory and never written to physical storage. (The
# secret volume uses an in-memory filesystem (tmpfs) for the Secret files.)
# From K8s version 1.7, etcd stores Secrets in encrypted form.
resource "kubernetes_secret" "secrets" {
  count = length(var.secrets)
  metadata {
    name = var.secrets[count.index].name
    namespace = var.namespace
    labels = {
      app = var.app_name
    }
    annotations = var.secrets[count.index].annotations
  }
  # Plain-text data.
  data = var.secrets[count.index].data
  binary_data = var.secrets[count.index].binary_data
  # https://kubernetes.io/docs/concepts/configuration/secret/#secret-types
  # https://kubernetes.io/docs/concepts/configuration/secret/#serviceaccount-token-secrets
  type = var.secrets[count.index].type
}

# A ServiceAccount is used by an application running inside a pod to authenticate itself with the
# API server. A default ServiceAccount is automatically created for each namespace; each pod is
# associated with exactly one ServiceAccount, but multiple pods can use the same ServiceAccount. A
# pod can only use a ServiceAccount from the same namespace.
resource "kubernetes_service_account" "service_account" {
  count = var.service_account == null ? 0 : 1
  metadata {
    name = var.service_account.name
    namespace = var.service_account.namespace
    labels = var.service_account.labels
    # https://kubernetes.io/docs/concepts/security/service-accounts/#enforce-mountable-secrets
    annotations = var.service_account.annotations
  }
  # If you don't want the kubelet to automatically mount a ServiceAccount's API credentials, you
  # can opt out of the default behavior. You can opt out of automounting API credentials on
  # /var/run/secrets/kubernetes.io/serviceaccount/token for a service account by setting
  # 'automountServiceAccountToken: false' on the ServiceAccount.
  automount_service_account_token = var.service_account.automount_service_account_token
  dynamic "secret" {
    for_each = var.service_account.secrets
    iterator = it
    content {
      name = it.value["name"]
    }
  }
}

# Roles define WHAT can be done; role bindings define WHO can do it.
# The distinction between a Role/RoleBinding and a ClusterRole/ClusterRoleBinding is that the Role/
# RoleBinding is a namespaced resource; ClusterRole/ClusterRoleBinding is a cluster-level resource.
# A Role resource defines what actions can be taken on which resources; i.e., which types of HTTP
# requests can be performed on which RESTful resources.
resource "kubernetes_role" "role" {
  count = var.role == null ? 0 : 1
  metadata {
    name = var.role.name
    namespace = var.role.namespace
    labels = var.role.labels
    annotations = var.role.annotations
  }
  dynamic "rule" {
    for_each = var.role.rules
    iterator = it
    content {
      api_groups = it.value["api_groups"]
      verbs = it.value["verbs"]
      resources = it.value["resources"]
      resource_names = it.value["resource_names"]
    }
  }
}

# Bind the role to the service account.
resource "kubernetes_role_binding" "role_binding" {
  count = var.role_binding == null ? 0 : 1
  metadata {
    name = var.role_binding.name
    namespace = var.role_binding.namespace
    labels = var.role_binding.labels
    annotations = var.role_binding.annotations
  }
  # A RoleBinding always references a single Role, but it can bind the Role to multiple subjects.
  role_ref {
    kind = var.role_binding.role_ref.kind
    # This RoleBinding references the Role specified below...
    name = var.role_binding.role_ref.name
    api_group = var.role_binding.role_ref.api_group
  }
  # ... and binds it to the specified ServiceAccount in the specified namespace.
  dynamic "subject" {
    # The default permissions for a ServiceAccount don't allow it to list or modify any resources.
    for_each = var.role_binding.subjects
    iterator = it
    content {
      kind = it.value["kind"]
      name = it.value["name"]
      namespace = it.value["namespace"]
      api_group = it.value["api_group"]
    }
  }
}




# The ConfigMap passes to the rabbitmq daemon a bootstrap configuration which mainly defines peer
# discovery and connectivity settings.
resource "kubernetes_config_map" "config" {
  metadata {
    name = "${var.service_name}-config"
    namespace = var.namespace
    labels = {
      app = var.app_name
    }
  }
  data = {
    # The enabled_plugins file is usually located in the node data directory or under /etc,
    # together with configuration files. The file contains a list of plugin names ending with
    # a dot.
    "enabled_plugins" = "[rabbitmq_federation, rabbitmq_management, rabbitmq_peer_discovery_k8s]."
    "rabbitmq.conf" = "${file("${var.path_rabbitmq_files}/rabbitmq.conf")}"
  }
}

# RabbitMQ requires using a StatefulSet to deploy a RabbitMQ cluster to Kubernetes. The StatefulSet
# ensures that the RabbitMQ nodes are deployed in order, one at a time. This avoids running into a
# potential peer discovery race condition when deploying a multi-node RabbitMQ cluster.
#
# There are other, equally important reasons for using a StatefulSet instead of a Deployment:
# sticky identity, simple network identifiers, stable persistent storage and the ability to perform
# ordered rolling upgrades.
#
# $ kubectl get sts -n memories
#
# Storage for stateful pods needs to be persistent and decoupled from the pods.
# A StatefulSet must guarantee at-most-one semantics for stateful pod instances.
resource "kubernetes_stateful_set" "stateful_set" {
  metadata {
    name = var.service_name
    namespace = var.namespace
    labels = {
      app = var.app_name
    }
  }
  #
  spec {
    replicas = var.replicas
    # The name of the service that governs this StatefulSet.
    # This service must exist before the StatefulSet and is responsible for the network identity of
    # the set. Pods get DNS/hostnames that follow the pattern:
    #   pod-name.service-name.namespace.svc.cluster.local.
    service_name = local.svc_name
    pod_management_policy = var.pod_management_policy
    # Pod Selector - You must set the .spec.selector field of a StatefulSet to match the labels of
    # its .spec.template.metadata.labels. Failing to specify a matching Pod Selector will result in
    # a validation error during StatefulSet creation.
    selector {
      match_labels = {
        # It must match the labels in the Pod template (.spec.template.metadata.labels).
        pod_selector_lbl = local.pod_selector_label
      }
    }
    # Pod template.
    template {
      metadata {
        # Labels attach to the Pod.
        labels = {
          app = var.app_name
          # It must match the label for the pod selector (.spec.selector.matchLabels).
          pod_selector_lbl = local.pod_selector_label
          # It must match the label selector of the Service.
          svc_selector_lbl = local.svc_selector_label
          rmq_lbl = local.rmq_label
        }
      }
      #
      spec {
        service_account_name = kubernetes_service_account.service_account.metadata[0].name
        affinity {
          pod_anti_affinity {
            required_during_scheduling_ignored_during_execution {
              label_selector {
                match_expressions {
                  # Description of the pod label that determines when the anti-affinity rule
                  # applies. Specifies a key and value for the label.
                  key = "rmq_lbl"
                  # The operator represents the relationship between the label on the existing
                  # pod and the set of values in the matchExpression parameters in the
                  # specification for the new pod. Can be In, NotIn, Exists, or DoesNotExist.
                  operator = "In"
                  values = ["${local.rmq_label}"]
                }
              }
              topology_key = "kubernetes.io/hostname"
            }
          }
        }
        termination_grace_period_seconds = var.termination_grace_period_seconds
        container {
          name = var.service_name
          image = var.image_tag
          image_pull_policy = var.image_pull_policy
          dynamic "security_context" {
            for_each = var.security_context
            iterator = item
            content {
              run_as_non_root = item.value.run_as_non_root
              run_as_user = item.value.run_as_user
              run_as_group = item.value.run_as_group
              read_only_root_filesystem = item.value.read_only_root_filesystem
            }
          }
          # Specifying ports in the pod definition is purely informational. Omitting them has no
          # effect on whether clients can connect to the pod through the port or not. If the
          # container is accepting connections through a port bound to the 0.0.0.0 address, other
          # pods can always connect to it, even if the port isn't listed in the pod spec
          # explicitly. Nonetheless, it is good practice to define the ports explicitly so that
          # everyone using the cluster can quickly see what ports each pod exposes.
          dynamic "port" {
            for_each = var.ports
            content {
              name = port.value.name
              container_port = port.value.target_port  # The port the app is listening.
              protocol = port.value.protocol
            }
          }
          resources {
            requests = {
              # If a Container specifies its own memory limit, but does not specify a memory
              # request, Kubernetes automatically assigns a memory request that matches the limit.
              # Similarly, if a Container specifies its own CPU limit, but does not specify a CPU
              # request, Kubernetes automatically assigns a CPU request that matches the limit.
              cpu = var.qos_requests_cpu == "" ? var.qos_limits_cpu : var.qos_requests_cpu
              memory = var.qos_requests_memory == "" ? var.qos_limits_memory : var.qos_requests_memory
            }
            limits = {
              cpu = var.qos_limits_cpu
              memory = var.qos_limits_memory
            }
          }
          # Using the Pod field as a value for the environment variable; pass RABBIT_POD_NAME to
          # build the FQDN.
          env {
            name = "RABBIT_POD_NAME"
            value_from {
              field_ref {
                field_path = "metadata.name"
              }
            }
          }
          # Using the Pod field as a value for the environment variable; pass RABBIT_POD_NAMESPACE
          # to build the FQDN.
          env {
            name = "RABBIT_POD_NAMESPACE"
            value_from {
              field_ref {
                field_path = "metadata.namespace"
              }
            }
          }
          # The name of the headless service needs to be provided to the discovery plugin via this
          # environment variable. It uses the name to query the K8s API for information on all pods
          # selected by the service.
          env {
            name = "K8S_SERVICE_NAME"
            value = local.svc_name
          }
          # When a node starts up, it checks whether it has been assigned a node name. If no value
          # was explicitly configured, the node resolves its hostname and prepends rabbit to it to
          # compute its node name.
          # Build the rabbitmq host FQDN.
          env {
            name = "RABBITMQ_NODENAME"
            value = "rabbit@$(RABBIT_POD_NAME).$(K8S_SERVICE_NAME).$(RABBIT_POD_NAMESPACE).svc.cluster.local"
          }
          # Build the cluster DNS domain name.
          # Suffix to match FQDN of rabbitmq instances in the K8s namespace.
          env {
            name = "K8S_HOSTNAME_SUFFIX"
            value = ".$(K8S_SERVICE_NAME).$(RABBIT_POD_NAMESPACE).svc.cluster.local"
          }
          # This environment variable is only mean to be used in development and CI environments.
          # This has the same meaning as default_user in rabbitmq.conf but higher priority. This
          # option may be more convenient in cases where providing a config file is impossible, and
          # environment variables is the only way to seed a user.
          env {
            name = "RABBITMQ_DEFAULT_PASS"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.secret.metadata[0].name
                key = "pass"
              }
            }
          }
          # This environment variable is only mean to be used in development and CI environments.
          # This has the same meaning as default_pass in rabbitmq.conf but higher priority. This
          # option may be more convenient in cases where providing a config file is impossible, and
          # environment variables is the only way to seed a user.
          env {
            name = "RABBITMQ_DEFAULT_USER"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.secret.metadata[0].name
                key = "user"
              }
            }
          }
          # dynamic "env" {
          #   for_each = var.env_field_ref
          #   content {
          #     name = env.value["name"]
          #     value_from {
          #       field_ref {
          #         field_path = env.value["field_path"]
          #       }
          #     }
          #   }
          # }
          dynamic "env" {
            for_each = var.env
            content {
              name = env.key
              value = env.value
            }
          }
          # liveness_probe {
          #   exec {
          #     command = ["rabbitmq-diagnostics", "status", "--erlang-cookie", "$(RABBITMQ_ERLANG_COOKIE)"]
          #   }
          #   initial_delay_seconds = 60
          #   # See https://www.rabbitmq.com/monitoring.html for monitoring frequency recommendations.
          #   period_seconds = 60
          #   timeout_seconds = 15
          #   failure_threshold = 3
          #   success_threshold = 1
          # }
          # readiness_probe {
          #   exec {
          #     command = ["rabbitmq-diagnostics", "status", "--erlang-cookie", "$(RABBITMQ_ERLANG_COOKIE)"]
          #   }
          #   initial_delay_seconds = 20
          #   period_seconds = 60
          #   timeout_seconds = 10
          # }
          volume_mount {
            name = "rabbitmq-storage"
            mount_path = "/var/lib/rabbitmq/mnesia"
            read_only = false
          }
          volume_mount {
            name = "rabbitmq-storage"
            mount_path = "/var/lib/rabbitmq/mnesia/$(RABBITMQ_NODENAME).pid"
            sub_path = "$(RABBITMQ_NODENAME).pid"
            read_only = false
          }
          # In Linux when a filesystem is mounted into a non-empty directory, the directory will
          # only contain the files from the newly mounted filesystem. The files in the original
          # directory are inaccessible for as long as the filesystem is mounted. In cases when the
          # original directory contains crucial files, mounting a volume could break the container.
          # To overcome this limitation, K8s provides an additional subPath property on the
          # volumeMount; this property mounts a single file or a single directory from the volume
          # instead of mounting the whole volume, and it does not hide the existing files in the
          # original directory.
          volume_mount {
            name = "erlang-cookie"
            # Mounting into a file, not a directory.
            mount_path = "/var/lib/rabbitmq/mnesia/.erlang.cookie"
            # Instead of mounting the whole volume, only mounting the given entry.
            sub_path = ".erlang.cookie"
            read_only = false
          }
          volume_mount {
            name = "config"
            mount_path = "/config/rabbitmq"
            read_only = true
          }
        }
        volume {
          name = "erlang-cookie"
          secret {
            secret_name = kubernetes_secret.secret.metadata[0].name
            default_mode = "0600" # Octal
            # Selecting which entries to include in the volume by listing them.
            items {
              # Include the entry under this key.
              key = "cookie"
              # The entry's value will be stored in this file.
              path = ".erlang.cookie"
            }
          }
        }
        volume {
          name = "config"
          config_map {
            name = kubernetes_config_map.config.metadata[0].name
            # By default, the permissions on all files in a configMap volume are set to 644
            # (rw-r--r--).
            default_mode = "0600" # Octal
            items {
              key = "enabled_plugins"
              path = "enabled_plugins" #File name.
            }
            items {
              key = "rabbitmq.conf"
              path = "rabbitmq.conf" #File name.
            }
          }
        }
      }
    }
    # This template will be used to create a PersistentVolumeClaim for each pod.
    # Since PersistentVolumes are cluster-level resources, they do not belong to any namespace, but
    # PersistentVolumeClaims can only be created in a specific namespace; they can only be used by
    # pods in the same namespace.
    #
    # In order for RabbitMQ nodes to retain data between Pod restarts, node's data directory must
    # use durable storage. A Persistent Volume must be attached to each RabbitMQ Pod.
    #
    # If a transient volume is used to back a RabbitMQ node, the node will lose its identity and
    # all of its local data in case of a restart. This includes both schema and durable queue data.
    # Syncing all of this data on every node restart would be highly inefficient. In case of a loss
    # of quorum during a rolling restart, this will also lead to data loss.
    volume_claim_template {
      metadata {
        name = "rabbitmq-storage"
        namespace = var.namespace
        labels = {
          app = var.app_name
        }
      }
      spec {
        access_modes = var.pvc_access_modes
        storage_class_name = var.pvc_storage_class_name
        resources {
          requests = {
            storage = var.pvc_storage_size
          }
        }
      }
    }
  }
}



# Before deploying a StatefulSet, you will need to create a headless Service, which will be used
# to provide the network identity for your stateful pods.
resource "kubernetes_service" "headless_service" {
  metadata {
    name = local.svc_name
    namespace = var.namespace
    labels = {
      app = var.app_name
    }
  }
  #
  spec {
    selector = {
      # All pods with the svc_selector_lbl=local.svc_selector_label label belong to this service.
      svc_selector_lbl = local.svc_selector_label
    }
    session_affinity = var.service_session_affinity
    dynamic "port" {
      for_each = var.ports
      content {
        name = port.value.name
        port = port.value.service_port
        target_port = port.value.target_port
        protocol = port.value.protocol
      }
    }
    type = var.service_type
    cluster_ip = "None" # Headless Service.
    publish_not_ready_addresses = var.publish_not_ready_addresses
  }
}
