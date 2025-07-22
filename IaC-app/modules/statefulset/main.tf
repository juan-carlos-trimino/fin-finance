/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
-------------------------------------------------------
Define input variables to the module.
***/
variable affinity {
  default = []
  type = list(object({
    pod_anti_affinity = optional(object({
      required_during_scheduling_ignored_during_execution = optional(list(object({
        topology_key = string
        namespaces = optional(set(string), [])
        match_labels = optional(map(string), {})
        match_expressions = optional(list(object({
        label_selector = object({
          key = string
          # Valid operators are In, NotIn, Exists, and DoesNotExist.
          operator = string
          # If the operator is In or NotIn, the values array must be non-empty. If the operator is
          # Exists or DoesNotExist, the values array must be empty.
          values = set(string)
        })
        })), [])
      })), [])
    }), {})
  }))
}
variable app_name {
  type = string
}
variable app_version {
  type = string
}
variable automount_service_account_token {
  default = false
  type = bool
}
variable config_map {
  default = []
  type = list(object({
    name = string
    namespace = string
    labels = optional(map(string), {})
    # Binary data need to be base64 encoded.
    binary_data = optional(map(string), {})
    data = optional(map(string), {})
    immutable = optional(bool, false)
  }))
}
variable env {
  default = {}
  type = map(any)
}
variable env_field {
  default = []
  type = list(object({
    name = string
    field_path = string
  }))
}
variable image_pull_policy {
  default = "Always"
  type = string
}
variable image_tag {
  default = ""
  type = string
}
variable labels {
  default = {}
  type = map(string)
}
variable namespace {
  default = "default"
  type = string
}
/***
To relax the StatefulSet ordering guarantee while preserving its uniqueness and identity
guarantee.
***/
variable pod_management_policy {
  default = "OrderedReady"
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
/***
The primary use case for setting this field is to use a StatefulSet's Headless Service to
propagate SRV records for its Pods without respect to their readiness for purpose of peer
discovery.
***/
variable publish_not_ready_addresses {
  default = "false"
  type = bool
}
variable replicas {
  default = 1
  type = number
}
variable resources {
  default = {}
  type = object({
    requests_cpu = optional(string)
    requests_memory = optional(string)
    limits_cpu = optional(string)
    limits_memory = optional(string)
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
variable secrets {
  default = []
  type = list(object({
    name = string
    namespace = string
    labels = optional(map(string), {})
    annotations = optional(map(string), {})
    data = optional(map(string), {})
    binary_data = optional(map(string), {})  # base64 encoding.
    type = optional(string, "Opaque")
  }))
  sensitive = true
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
variable service_name {
  type = string
}
variable service_session_affinity {
  default = "None"
  type = string
}
/***
The ServiceType allows to specify what kind of Service to use: ClusterIP (default), NodePort,
LoadBalancer, and ExternalName.
***/
variable service_type {
  default = "ClusterIP"
}
variable termination_grace_period_seconds {
  default = 30
  type = number
}
variable volume_claim_templates {
  default = []
  type = list(object({
    name = string
    namespace = string
    labels = optional(map(string), {})
    access_modes = list(string)
    # A volumeMode of Filesystem presents a volume as a directory within the Pod's filesystem while
    # a volumeMode of Block presents it as a raw block storage device. Filesystem is the default
    # and usually preferred mode, enabling standard file system operations on the volume. Block
    # mode is used for applications that need direct access to the block device, like databases
    # requiring low-latency access.
    volume_mode = optional(string, "Filesystem")
    storage = string
    # By specifying an empty string ("") as the storage class name, the PVC binds to a
    # pre-provisioned PV instead of dynamically provisioning a new one.
    storage_class_name = optional(string)
  }))
}
variable volume_config_map {
  default = []
  type = list(object({
    name = string
    # Name of the ConfigMap containing the files to add to the container.
    config_map_name = string
    # Although ConfigMaps should be used for non-sensitive configuration data, you may want to
    # make the file readable and writeble only to the user and group that owned the file; e.g.,
    # default_mode = "6600" (-rw-rw------)
    # The default permission is "6440" (-rw-r--r----)
    default_mode = optional(string)
    # An array of keys from the ConfigMap to create as files.
    items = optional(list(object({
      # Include the entry under this key.
      key = string
      # The entry's value should be stored in this file.
      path = string
    })), [])
  }))
}
variable volume_empty_dir {
  description = "(Optional) A temporary directory that shares a pod's lifetime."
  default = []
  type = list(object({
    name = string
    medium = optional(string)
    size_limit = optional(string)
  }))
}
variable volume_mount {
  default = []
  type = list(object({
    name = string
    mount_path = string
    sub_path = optional(string)
    read_only = optional(bool)
  }))
}
variable volume_secrets {
  default = []
  type = list(object({
    name = string
    # Name of the ConfigMap containing the files to add to the container.
    secret_name = string
    # Although ConfigMaps should be used for non-sensitive configuration data, you may want to
    # make the file readable and writeble only to the user and group that owned the file; e.g.,
    # default_mode = "6600" (-rw-rw------)
    # The default permission is "6440" (-rw-r--r----)
    default_mode = optional(string)
    # An array of keys from the ConfigMap to create as files.
    items = optional(list(object({
      # Include the entry under this key.
      key = string
      # The entry's value should be stored in this file.
      path = string
    })), [])
  }))
}

/***
Define local variables.
***/
locals {
  svc_name = "${var.service_name}-headless"
  pod_selector_label = "ps-${var.service_name}"
  svc_selector_label = "svc-${local.svc_name}"
}

resource "kubernetes_secret" "secrets" {
  count = length(var.secrets)
  metadata {
    name = var.secrets[count.index].name
    namespace = var.secrets[count.index].namespace
    labels = var.secrets[count.index].labels
    annotations = var.secrets[count.index].annotations
  }
  # Plain-text data.
  data = var.secrets[count.index].data
  /***
  ***/
  binary_data = var.secrets[count.index].binary_data
  # https://kubernetes.io/docs/concepts/configuration/secret/#secret-types
  # https://kubernetes.io/docs/concepts/configuration/secret/#serviceaccount-token-secrets
  type = var.secrets[count.index].type
}

resource "kubernetes_service_account" "service_account" {
  count = var.service_account == null ? 0 : 1
  metadata {
    name = var.service_account.name
    namespace = var.service_account.namespace
    labels = var.service_account.labels
    # https://kubernetes.io/docs/concepts/security/service-accounts/#enforce-mountable-secrets
    annotations = var.service_account.annotations
  }
  /***
  If you don't want the kubelet to automatically mount a ServiceAccount's API credentials, you
  can opt out of the default behavior. You can opt out of automounting API credentials on
  /var/run/secrets/kubernetes.io/serviceaccount/token for a service account by setting
  'automountServiceAccountToken: false' on the ServiceAccount.
  ***/
  automount_service_account_token = var.service_account.automount_service_account_token
  dynamic "secret" {
    for_each = var.service_account.secrets
    iterator = it
    content {
      name = it.value["name"]
    }
  }
}

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

resource "kubernetes_config_map" "config" {
  count = length(var.config_map)
  metadata {
    name = var.config_map[count.index].name
    namespace = var.namespace
    labels = var.config_map[count.index].labels
  }
  data = var.config_map[count.index].data
  binary_data = var.config_map[count.index].binary_data
  immutable = var.config_map[count.index].immutable
}

resource "kubernetes_stateful_set" "stateful_set" {
  metadata {
    name = var.service_name
    namespace = var.namespace
    # Labels attach to the Deployment.
    labels = var.labels
  }
  #
  spec {
    /***
    The name of the service that governs this StatefulSet.
    This service must exist before the StatefulSet and is responsible for the network identity of
    the set. Pods get DNS/hostnames that follow the pattern:
      pod-name.service-name.namespace.svc.cluster.local.
    ***/
    service_name = local.svc_name
    replicas = var.replicas
    pod_management_policy = var.pod_management_policy
    /***
    Pod Selector - You must set the .spec.selector field of a StatefulSet to match the labels of
    its .spec.template.metadata.labels. Failing to specify a matching Pod Selector will result in
    a validation error during StatefulSet creation.
    ***/
    selector {
      match_labels = {
        # It must match the labels in the Pod template (spec.template.metadata.labels).
        pod_selector_lbl = local.pod_selector_label
      }
    }
    # Pod template.
    template {
      metadata {
        # Labels attach to the Pod.
        labels = {
          app = var.app_name
          # It must match the label selector of spec.selector.match_labels.
          pod_selector_lbl = local.pod_selector_label
          # It must match the label selector of the Service.
          svc_selector_lbl = local.svc_selector_label
        }
      }
      #
      spec {
        /***
        By default, the default-token Secret is mounted into every container, but you can
        disable that in each pod by setting the automountServiceAccountToken field in the pod spec
        to false or by setting it to false on the service account the pod is using.
        ***/
        automount_service_account_token = var.automount_service_account_token
        service_account_name = var.service_account == null ? "default" : var.service_account.name
        dynamic "affinity" {
          for_each = var.affinity
          iterator = it1
          content {
            pod_anti_affinity {
              dynamic "required_during_scheduling_ignored_during_execution" {
                for_each = it1.value["required_during_scheduling_ignored_during_execution"]
                iterator = it2
                content {
                  label_selector {
                    match_labels = it2.match_labels
                    dynamic "match_expressions" {
                      for_each = it2.value["match_expressions"]
                      iterator = it3
                      content {
                        key = it3.key
                        operator = it3.operation
                        values = it3.values
                      }
                    }
                  }
                  topology_key = "kubernetes.io/hostname"
                }
              }
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
          /***
          Specifying ports in the pod definition is purely informational. Omitting them has no
          effect on whether clients can connect to the pod through the port or not. If the
          container is accepting connections through a port bound to the 0.0.0.0 address, other
          pods can always connect to it, even if the port isn't listed in the pod spec
          explicitly. Nonetheless, it is good practice to define the ports explicitly so that
          everyone using the cluster can quickly see what ports each pod exposes.
          ***/
          dynamic "port" {
            for_each = var.ports
            content {
              name = port.value.name
              container_port = port.value.target_port  # The port the app is listening.
              protocol = port.value.protocol
            }
          }
          dynamic "resources" {
            for_each = var.resources == {} ? [] : [1]
            content {
              requests = {
                cpu = var.resources.requests_cpu
                memory = var.resources.requests_memory
              }
              limits = {
                cpu = var.resources.limits_cpu
                memory = var.resources.limits_memory
              }
            }
          }
          /***
          To list all of the environment variables:
          Linux: $ printenv
          ***/
          dynamic "env" {
            for_each = var.env
            content {
              name = env.key
              value = env.value
            }
          }
          dynamic "env" {
            for_each = var.env_field
            content {
              name = env.value["name"]
              value_from {
                field_ref {
                  field_path = env.value["field_path"]
                }
              }
            }
          }
          dynamic "env_from" {
            for_each = var.config_map
            content {
              config_map_ref {
                name = env_from.value["name"]
              }
            }
          }
          dynamic "env_from" {
            for_each = var.secrets
            content {
              secret_ref {
                name = env_from.value["name"]
              }
            }
          }
          dynamic "volume_mount" {
            for_each = var.volume_mount
            content {
              name = volume_mount.value["name"]
              mount_path = volume_mount.value["mount_path"]
              sub_path = volume_mount.value["sub_path"]
              read_only = volume_mount.value["read_only"]
            }
          }
        }
        /***
        Set volumes at the Pod level, then mount them into containers inside that Pod.

        By default, K8s emptyDir volumes are created with root:root ownership and 750
        permissions. This means that the directory created by K8s for the emptyDir volume is
        owned by the root user and group, which translates to read-write-execute permissions for
        the owner (root), read-execute permissions for the group, and no permissions for others.
        (For directories, execute permission is required to access the contents of the
        directory.)
        In many cases, especially when running containers as non-root users, this default
        ownership can lead to permission issues when containers try to write to the emptyDir
        volume. To address this, you might need to adjust the ownership and permissions of the
        emptyDir volume or consider using other volume types or approaches.
        ***/
        dynamic "volume" {
          for_each = var.volume_empty_dir
          content {
            name = volume.value["name"]
            empty_dir {
              medium = volume.value["medium"]
              size_limit = volume.value["size_limit"]
            }
          }
        }
        dynamic "volume" {
          for_each = var.volume_config_map
          iterator = it
          content {
            name = it.value["name"]
            config_map {
              name = it.value["config_map_name"]
              default_mode = it.value["default_mode"]
              dynamic "items" {
                for_each = it.value["items"]
                iterator = itn
                content {
                  key = itn.value["key"]
                  path = itn.value["path"]
                }
              }
            }
          }
        }
        dynamic "volume" {
          for_each = var.volume_secrets
          iterator = it
          content {
            name = it.value["name"]
            secret {
              secret_name = it.value["secret_name"]
              default_mode = it.value["default_mode"]
              dynamic "items" {
                for_each = it.value["items"]
                iterator = itn
                content {
                  key = itn.value["key"]
                  path = itn.value["path"]
                }
              }
            }
          }
        }
      }
    }
    /***
    This template will be used to create a PersistentVolumeClaim for each pod.
    Since PersistentVolumes are cluster-level resources, they do not belong to any namespace, but
    PersistentVolumeClaims can only be created in a specific namespace; they can only be used by
    pods in the same namespace.
    ***/
    dynamic "volume_claim_template" {
      # for_each = var.volume_claim_template
      for_each = var.volume_claim_templates
      iterator = it
      content {
        metadata {
          name = it.value["name"]
          namespace = it.value["namespace"]
          labels = it.value["labels"]
        }
        spec {
          access_modes = it.value["access_modes"]
          volume_mode = it.value["volume_mode"]
          storage_class_name = it.value["storage_class_name"]
          resources {
            requests = {
              storage = it.value["storage"]
            }
          }
        }
      }
    }
  }
}

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
    type = "ClusterIP"  # Default.
    cluster_ip = "None"  # Headless Service.
    publish_not_ready_addresses = var.publish_not_ready_addresses
  }
}
