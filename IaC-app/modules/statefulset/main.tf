/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
-------------------------------------------------------
Define input variables to the module.
***/
variable affinity {
  default = {}
  type = object({
    # affinity_type = string
    pod_anti_affinity = optional(object({
      required_during_scheduling_ignored_during_execution = optional(list(object({
        topology_key = string
        namespaces = optional(set(string), [])
        label_selector = optional(object({
          match_labels = optional(map(string), {})
          match_expressions = optional(list(object({
            key = string
            # Valid operators are In, NotIn, Exists, and DoesNotExist.
            operator = string
            # If the operator is In or NotIn, the values array must be non-empty. If the operator is
            # Exists or DoesNotExist, the values array must be empty.
            values = set(string)
          })), [])
        }), {})
      })), [])
    }), {})
    #
    pod_affinity = optional(object({
      required_during_scheduling_ignored_during_execution = optional(list(object({
        topology_key = string
        namespaces = optional(set(string), [])
        label_selector = optional(object({
          match_labels = optional(map(string), {})
          match_expressions = optional(list(object({
            key = string
            # Valid operators are In, NotIn, Exists, and DoesNotExist.
            operator = string
            # If the operator is In or NotIn, the values array must be non-empty. If the operator is
            # Exists or DoesNotExist, the values array must be empty.
            values = set(string)
          })), [])
        }), {})
      })), [])
    }), {})
  })
}
variable app_name {
  type = string
}
variable app_version {
  type = string
}
variable args {
  default = []
  type = list(string)
}
variable automount_service_account_token {
  default = false
  type = bool
}
# When defined, it overrides the image's default command.
variable command {
  default = []
  type = list(string)
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
# https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/pod#nested-schema-for-speccontainersecurity_context
variable containers_security_context {  # spec.containers[x].securityContext
  default = {}
  type = object({
    allow_privilege_escalation = optional(bool)
    capabilities = optional(object({
      add = optional(list(string))
      drop = optional(list(string))
    }))
    privileged = optional(bool)
    read_only_root_filesystem = optional(bool)
    # Processes inside container will run as primary group "run_as_group".
    run_as_group = optional(number)
    run_as_non_root = optional(bool)
    # Processes inside container will run as user "run_as_user".
    run_as_user = optional(number)
    se_linux_options = optional(object({
      user = optional(string)
      role = optional(string)
      type = optional(string)
      level = optional(string)
    }))
    seccomp_profile = optional(object({
      type = optional(string)
      localhost_profile = optional(string)
    }))
  })
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
variable init_container {
  default = []
  type = list(object({
    name = string
    image = string
    image_pull_policy = optional(string)
    command = optional(list(string))
    env = optional(map(any), {})
    env_from_secrets = optional(list(string), [])
    security_context = optional(object({
      allow_privilege_escalation = optional(bool)
      capabilities = optional(object({
        add = optional(list(string))
        drop = optional(list(string))
      }))
      privileged = optional(bool)
      read_only_root_filesystem = optional(bool)
      # Processes inside container will run as primary group "run_as_group".
      run_as_group = optional(number)
      run_as_non_root = optional(bool)
      # Processes inside container will run as user "run_as_user".
      run_as_user = optional(number)
      se_linux_options = optional(object({
        user = string
        role = string
        type = string
        level = string
      }))
      seccomp_profile = optional(object({
        type = string
       localhost_profile = optional(string)
      }))
    }))
    volume_mounts = optional(list(object({
      name = string
      mount_path = string
      sub_path = optional(string)
      read_only = optional(bool)
    })), [])
  }))
}
variable labels {
  default = {}
  type = map(string)
}
variable liveness_probe {
  default = []
  type = list(object({
    initial_delay_seconds = optional(number)
    period_seconds = optional(number)
    timeout_seconds = optional(number)
    failure_threshold = optional(number)
    success_threshold = optional(number)
    http_get = optional(list(object({
      host = optional(string)
      path = optional(string)
      port = number
      scheme = optional(string)
      http_header = optional(list(object({
        name = string
        value = string
      })), [])
    })), [])
    exec = optional(object({
      command = list(string)
    }), null)
    tcp_socket = optional(object({
      port = number
    }), null)
  }))
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
  type = string
}
# https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/pod#nested-schema-for-specsecurity_context
variable pod_security_context {  # spec.securityContext
  default = {}
  type = object({
    # fs_group ensures that any volumes mounted by the Pod will have their ownership changed to
    # this specified group ID.
    # The "volumeMounts.mountPath" will have its group ownership set to "fs_group".
    # Any files created within "mountPath" by the container will be owned by user "run_as_user"
    # and group "fs_group" (due to "fsGroup").
    fs_group = optional(number)
    fs_group_change_policy = optional(string)
    # Processes inside container will run as primary group "run_as_group".
    run_as_group = optional(number)
    run_as_non_root = optional(bool)
    # Processes inside container will run as user "run_as_user".
    run_as_user = optional(number)
    se_linux_options = optional(object({
      user = optional(string)
      role = optional(string)
      type = optional(string)
      level = optional(string)
    }))
    seccomp_profile = optional(object({
      type = optional(string)
      localhost_profile = optional(string)
    }))
    supplemental_groups = optional(set(number))
    sysctl = optional(list(object({
      name = string
      value = string
    })), [])
    windows_options = optional(object({
      gmsa_credential_spec = optional(string)
      gmsa_credential_spec_name = optional(string)
      host_process = optional(bool)
      run_as_username = optional(string)
    }))
  })
}
variable readiness_probe {
  default = []
  type = list(object({
    /***
    Number of seconds after the container has started before liveness or readiness probes are
    initiated. Defaults to 0 seconds. Minimum value is 0.
    ***/
    initial_delay_seconds = optional(number)
    # How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
    period_seconds = optional(number)
    # Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1.
    timeout_seconds = optional(number)
    /***
    When a probe fails, Kubernetes will try failureThreshold times before giving up. Giving up in
    case of liveness probe means restarting the container. In case of readiness probe the Pod
    will be marked Unready. Defaults to 3. Minimum value is 1.
    ***/
    failure_threshold = optional(number)
    /***
    Minimum consecutive successes for the probe to be considered successful after having failed.
    Defaults to 1. Must be 1 for liveness and startup Probes. Minimum value is 1.
    ***/
    success_threshold = optional(number)
    http_get = optional(list(object({
      # Host name to connect to, defaults to the pod IP.
      host = optional(string)
      # Path to access on the HTTP server. Defaults to /.
      path = optional(string)
      /***
      Name or number of the port to access on the container. Number must be in the range 1 to
      65535.
      ***/
      port = number
      # Scheme to use for connecting to the host (HTTP or HTTPS). Defaults to HTTP.
      scheme = optional(string)
      http_header = optional(list(object({
        name = string
        value = string
      })), [])
    })), [])
    exec = optional(object({
      command = list(string)
    }), null)
    tcp_socket = optional(object({
      port = number
    }), null)
  }))
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
variable restart_policy {
  default = "Always"
  type = string
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
    immutable = optional(bool, true)
  }))
  sensitive = true
}
variable service {
  # default = null
  type = object({
    name = string
    namespace = string
    labels = optional(map(string), {})
    annotations = optional(map(string), {})
    # Only applies to types ClusterIP, NodePort, and LoadBalancer.
    selector = map(string)
    /***
    The service normally forwards each connection to a randomly selected backing pod. To ensure
    that connections from a particular client are passed to the same Pod each time, set the
    service's sessionAffinity property to ClientIP instead of None (default).
    Session affinity and Web Browsers (for LoadBalancer Services)
    Since the service is now exposed externally, accessing it with a web browser will hit the same
    pod every time. If the sessionAffinity is set to None, then why? The browser is using
    keep-alive connections and sends all its requests through a single connection. Services work at
    the connection level, and when a connection to a service is initially open, a random pod is
    selected and then all network packets belonging to that connection are sent to that single pod.
    Even with the sessionAffinity set to None, the same pod will always get hit (until the
    connection is closed).
    ***/
    session_affinity = optional(string, "None")
    /***
    The primary use case for setting this field is to use a StatefulSet's Headless Service to
    propagate SRV records for its Pods without respect to their readiness for purpose of peer
    discovery.
    ***/
    publish_not_ready_addresses = optional(bool, "false")
    /***
    The ServiceType allows to specify what kind of Service to use: ClusterIP (default), NodePort,
    LoadBalancer, and ExternalName.
    ***/
    type = optional(string, "ClusterIP")
    ports = optional(list(object({
      name = string
      service_port = number
      target_port = number
      node_port = optional(number)
      protocol = string
    })), [{
      name = "ports"
      service_port = 80
      target_port = 8080
      protocol = "TCP"
    }])
  })
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
variable statefulset_name {
  type = string
}
variable termination_grace_period_seconds {
  default = 30
  type = number
}
# https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/stateful_set#nestedblock--spec--update_strategy
variable update_strategy {
  default = null
  type = object({
    type = string
    partition = optional(number)
  })
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
    # default_mode = "0660" (-rw-rw----)
    # The default permission is "0644" (-rw-r--r--)
    default_mode = optional(string)
    # An array of keys from the ConfigMap to create as files.
    items = optional(list(object({
      # The configMap entry.
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
  pod_selector_label = "ps-${var.statefulset_name}"
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
  immutable = var.secrets[count.index].immutable
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
    name = var.statefulset_name
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
    service_name = var.service.name
    dynamic "update_strategy" {
      for_each = var.update_strategy == null ? [] : [1]
      content {
        type = var.update_strategy.type
        rolling_update {
          partition = var.update_strategy.partition
        }
      }
    }
    pod_management_policy = var.pod_management_policy
    # The desired number of pods that should be running.
    replicas = var.replicas
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
          svc_selector_label = var.service.selector["svc_selector_label"]
        }
      }
      #
      spec {
        dynamic "affinity" {
          for_each = var.affinity == {} ? [] : [1]
          content {
            dynamic "pod_anti_affinity" {
              for_each = var.affinity.pod_anti_affinity == {} ? [] : [1]
              content {
                dynamic "required_during_scheduling_ignored_during_execution" {
                  for_each = var.affinity.pod_anti_affinity.required_during_scheduling_ignored_during_execution
                  iterator = it1
                  content {
                    label_selector {
                      match_labels = it1.value["label_selector"].match_labels
                      dynamic "match_expressions" {
                        for_each = it1.value["label_selector"].match_expressions
                        iterator = it2
                        content {
                          key = it2.value["key"]
                          operator = it2.value["operator"]
                          values = it2.value["values"]
                        }
                      }
                    }
                    topology_key = "kubernetes.io/hostname"
                  }
                }
              }
            }
            dynamic "pod_affinity" {
              for_each = var.affinity.pod_affinity == {} ? [] : [1]
              content {
                dynamic "required_during_scheduling_ignored_during_execution" {
                  for_each = var.affinity.pod_affinity.required_during_scheduling_ignored_during_execution
                  iterator = it1
                  content {
                    label_selector {
                      match_labels = it1.value["label_selector"].match_labels
                      dynamic "match_expressions" {
                        for_each = it1.value["label_selector"].match_expressions
                        iterator = it2
                        content {
                          key = it2.value["key"]
                          operator = it2.value["operator"]
                          values = it2.value["values"]
                        }
                      }
                    }
                    topology_key = "kubernetes.io/hostname"
                  }
                }
              }
            }
          }
        }
        termination_grace_period_seconds = var.termination_grace_period_seconds
        service_account_name = var.service_account == null ? "default" : var.service_account.name
        /***
        By default, the default-token Secret is mounted into every container, but you can
        disable that in each pod by setting the automountServiceAccountToken field in the pod spec
        to false or by setting it to false on the service account the pod is using.
        ***/
        automount_service_account_token = var.automount_service_account_token
        restart_policy = var.restart_policy
        /***
        Security context options at the pod level serve as a default for all the pod's containers
        but can be overridden at the container level.
        ***/
        dynamic "security_context" {
          for_each = var.pod_security_context == {} ? [] : [1]
          content {
            /***
            Set the group that owns the pod volumes. This group will be used by K8s to change the
            permissions of all files/directories in the volumes, when the volumes are mounted by
            a pod.
            ***/
            fs_group = var.pod_security_context.fs_group
            /***
            By default, Kubernetes recursively changes ownership and permissions for the contents
            of each volume to match the fsGroup specified in a Pod's securityContext when that
            volume is mounted. For large volumes, checking and changing ownership and permissions
            can take a lot of time, slowing Pod startup. You can use the fsGroupChangePolicy
            field inside a securityContext to control the way that Kubernetes checks and manages
            ownership and permissions for a volume.
            ***/
            fs_group_change_policy = var.pod_security_context.fs_group_change_policy
            run_as_group = var.pod_security_context.run_as_group
            run_as_non_root = var.pod_security_context.run_as_non_root
            run_as_user = var.pod_security_context.run_as_user
            dynamic "se_linux_options" {
              for_each = var.pod_security_context.se_linux_options == null ? [] : [1]
              content {
                user = var.pod_security_context.se_linux_options.user
                role = var.pod_security_context.se_linux_options.role
                type = var.pod_security_context.se_linux_options.type
                level = var.pod_security_context.se_linux_options.level
              }
            }
            dynamic "seccomp_profile" {
              for_each = var.pod_security_context.seccomp_profile == null ? [] : [1]
              content {
                type = var.pod_security_context.seccomp_profile.type
                localhost_profile = var.pod_security_context.seccomp_profile.localhost_profile
              }
            }
            supplemental_groups = var.pod_security_context.supplemental_groups
            dynamic "sysctl" {
              for_each = var.pod_security_context.sysctl
              iterator = it
              content {
                name = it.name
                value = it.value
              }
            }
            dynamic "windows_options" {
              for_each = var.pod_security_context.windows_options == null ? [] : [1]
              content {
                gmsa_credential_spec = var.pod_security_context.windows_options.gmsa_credential_spec
                gmsa_credential_spec_name = var.pod_security_context.windows_options.gmsa_credential_spec_name
                host_process = var.pod_security_context.windows_options.host_process
                run_as_username = var.pod_security_context.windows_options.run_as_username
              }
            }
          }
        }
        # These containers are run during pod initialization.
        # An Init Container in K8s exists within the same Pod and therefore operates within the
        # same namespace as the Pod it belongs to.
        dynamic "init_container" {
          for_each = var.init_container
          iterator = it
          content {
            name = it.value["name"]
            image = it.value["image"]
            image_pull_policy = it.value["image_pull_policy"]
            dynamic "env" {
              for_each = it.value["env"]
              content {
                name = env.key
                value = env.value
              }
            }
            dynamic "env_from" {
              for_each = it.value["env_from_secrets"]
              iterator = it1
              content {
                secret_ref {
                  name = it1.value
                }
              }
            }
            command = it.value["command"]
            dynamic "security_context" {
              for_each = it.value.security_context == null ? [] : [1]
              # iterator = it1
              content {
                allow_privilege_escalation = it.value.security_context.allow_privilege_escalation
                dynamic "capabilities" {
                  for_each = it.value.security_context.capabilities == null ? [] : [1]
                  content {
                    add = it.value.security_context.capabilities.add
                    drop = it.value.security_context.capabilities.drop
                  }
                }
                privileged = it.value.security_context.privileged
                read_only_root_filesystem = it.value.security_context.read_only_root_filesystem
                run_as_group = it.value.security_context.run_as_group
                run_as_non_root = it.value.security_context.run_as_non_root
                run_as_user = it.value.security_context.run_as_user
                dynamic "se_linux_options" {
                  for_each = it.value.security_context.se_linux_options == null ? [] : [1]
                  content {
                    user = it.value.security_context.se_linux_options.user
                    role = it.value.security_context.se_linux_options.role
                    type = it.value.security_context.se_linux_options.type
                    level = it.value.security_context.se_linux_options.level
                  }
                }
                dynamic "seccomp_profile" {
                  for_each = it.value.security_context.seccomp_profile == null ? [] : [1]
                  content {
                    type = it.value.security_context.seccomp_profile.type
                    localhost_profile = it.value.security_context.seccomp_profile.localhost_profile
                  }
                }
              }
            }
            dynamic "volume_mount" {
              for_each = it.value["volume_mounts"]
              iterator = it1
              content {
                name = it1.value["name"]
                mount_path = it1.value["mount_path"]
                sub_path = it1.value["sub_path"]
                read_only = it1.value["read_only"]
              }
            }
          }
        }
        container {
          name = var.statefulset_name
          image = var.image_tag
          image_pull_policy = var.image_pull_policy
          /***
          Environment variables must be defined before they are used!!!
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
          /***
          In K8s, envFrom with secretRef is a method used to inject all key-value pairs from a
          specified Kubernetes Secret as environment variables into a container within a Pod.
          This differs from secretKeyRef which allows for the selection of specific keys from a
          Secret to be injected as environment variables.
          ***/
          dynamic "env_from" {
            for_each = var.secrets
            content {
              secret_ref {
                name = env_from.value["name"]
              }
            }
          }
          command = var.command
          args = var.args
          dynamic "security_context" {
            for_each = var.containers_security_context == {} ? [] : [1]
            content {
              allow_privilege_escalation = var.containers_security_context.allow_privilege_escalation
              dynamic "capabilities" {
                for_each = var.containers_security_context.capabilities == null ? [] : [1]
                content {
                  add = var.containers_security_context.capabilities.add
                  drop = var.containers_security_context.capabilities.drop
                }
              }
              privileged = var.containers_security_context.privileged
              read_only_root_filesystem = var.containers_security_context.read_only_root_filesystem
              run_as_group = var.containers_security_context.run_as_group
              run_as_non_root = var.containers_security_context.run_as_non_root
              run_as_user = var.containers_security_context.run_as_user
              dynamic "se_linux_options" {
                for_each = var.containers_security_context.se_linux_options == null ? [] : [1]
                content {
                  user = var.containers_security_context.se_linux_options.user
                  role = var.containers_security_context.se_linux_options.role
                  type = var.containers_security_context.se_linux_options.type
                  level = var.containers_security_context.se_linux_options.level
                }
              }
              dynamic "seccomp_profile" {
                for_each = var.containers_security_context.seccomp_profile == null ? [] : [1]
                content {
                  type = var.containers_security_context.seccomp_profile.type
                  localhost_profile = var.containers_security_context.seccomp_profile.localhost_profile
                }
              }
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
            for_each = var.service.ports
            content {
              name = port.value.name
              container_port = port.value.target_port  # The port the app is listening.
              protocol = port.value.protocol
            }
          }
          dynamic "liveness_probe" {
            for_each = var.liveness_probe
            iterator = it
            content {
              initial_delay_seconds = it.value["initial_delay_seconds"]
              period_seconds = it.value["period_seconds"]
              timeout_seconds = it.value["timeout_seconds"]
              failure_threshold = it.value["failure_threshold"]
              success_threshold = it.value["success_threshold"]
              /***
              K8s can probe a container using one of the three probes:
              The HTTP GET probe performs an HTTP GET request on the container. If the probe
              receives a response that doesn't represent an error (HTTP response code is 2xx or
              3xx), the probe is considered successful. If the server returns an error response
              code or it doesn't respond at all, the probe is considered a failure and the
              container will be restarted as a result.
              ***/
              dynamic "http_get" {
                for_each = it.value.http_get
                iterator = it1
                content {
                  host = it1.value["host"]
                  path = it1.value["path"]
                  port = it1.value["port"]
                  scheme = it1.value["scheme"]
                  dynamic "http_header" {
                    for_each = it1.value.http_header
                    iterator = it2
                    content {
                      name = it2.value["name"]
                      value = it2.value["value"]
                    }
                  }
                }
              }
              /***
              The Exec probe executes an arbitrary command inside the container and checks the
              command's exit status code. If the status code is 0, the probe is successful. All
              other codes are considered failures.
              ***/
              dynamic "exec" {
                for_each = it.value["exec"] != null ? [it.value["exec"]] : []
                content {
                  command = exec.value.command
                }
              }
              /***
              The TCP Socket probe tries to open a TCP connection to the specified port of the
              container. If the connection is established successfully, the probe is successful.
              Otherwise, the container is restarted.
              ***/
              dynamic "tcp_socket" {
                for_each = it.value["tcp_socket"] != null ? [it.value["tcp_socket"]] : []
                content {
                  port = tcp_socket.value.port
                }
              }
            }
          }
          /***
          Liveness probes keep pods healthy by killing unhealthy containers and replacing them
          with new healthy containers; readiness probes ensure that only pods with containers
          that are ready to serve requests receive them. Unlike liveness probes, if a container
          fails the readiness check, it won't be killed or restarted.
          ***/
          dynamic "readiness_probe" {
            for_each = var.readiness_probe
            content {
              initial_delay_seconds = readiness_probe.value["initial_delay_seconds"]
              period_seconds = readiness_probe.value["period_seconds"]
              timeout_seconds = readiness_probe.value["timeout_seconds"]
              failure_threshold = readiness_probe.value["failure_threshold"]
              success_threshold = readiness_probe.value["success_threshold"]
              /***
              K8s can probe a container using one of the three probes:
              The HTTP GET probe sends an HTTP GET request to the container, and the HTTP status
              code of the response determines whether the container is ready or not.
              ***/
              dynamic "http_get" {
                for_each = readiness_probe.value.http_get
                content {
                  host = http_get.value["host"]
                  path = http_get.value["path"]
                  port = http_get.value["port"]
                  scheme = http_get.value["scheme"]
                  dynamic "http_header" {
                    for_each = http_get.value.http_header
                    content {
                      name = http_headers.value["name"]
                      value = http_headers.value["value"]
                    }
                  }
                }
              }
              /***
              The Exec probe executes a process. The container's status is determined by the
              process' exit status code.
              ***/
              dynamic "exec" {
                # for_each = it.value["exec"] != null ? [it.value["exec"]] : []
                for_each = readiness_probe.value["exec"] != null ? [readiness_probe.value["exec"]] : []
                content {
                  command = exec.value.command
                }
              }
              /***
              The TCP Socket probe opens a TCP connection to a specified port of the container.
              If the connection is established, the container is considered ready.
              ***/
              dynamic "tcp_socket" {
                # for_each = it.value["tcp_socket"] != null ? [it.value["tcp_socket"]] : []
                for_each = readiness_probe.value["tcp_socket"] != null ? [readiness_probe.value["tcp_socket"]] : []
                content {
                  port = tcp_socket.value.port
                }
              }
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
          dynamic "volume_mount" {
            for_each = var.volume_mount
            content {
              name = volume_mount.value["name"]
              mount_path = volume_mount.value["mount_path"]
              /***
              When you mount a volume as a directory, you are hiding any files that are stored in
              the directory located inside the container image. In general, this is what happens in
              Linux when you mount a filesystem into a non-empty directory. The directory will only
              contain the files from the mounted filesystem, and the original files in that
              directory are inaccessible for as long as the filesystem is mounted. To add
              individual files into an existing directory without hiding existing files stored in
              it, you use the subPath property on the volumeMount as doing so allows you to mount a
              single file or a single directory from the volume instead of mounting the whole
              volume.
              ***/
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
        /***
        Volumes are defined as a part of a pod and share the same lifecycle as the pod. That is, a
        volume is created when the pod is started and is destroyed when the pod is deleted; a
        volume's contents will persist across container restarts. After a container is restarted,
        the new container can use all the files that were written to the volume by previous
        containers. Furthermore, if a pod contains multiple containers, the volume can be used by
        all of them at once.

        A configMap volume will expose each entry of the ConfigMap as a file. The process running
        in the container can obtain the entry's value by reading the contents of the file.
        ***/
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

/***
A headless service is a K8s service that does not have a cluster IP address. Instead, it provides
DNS records for the pods in the stateful set.
***/
resource "kubernetes_service" "headless_service" {
  metadata {
    name = var.service.name
    namespace = var.service.namespace
    labels = var.service.labels
  }
  #
  spec {
    # All pods with the svc_selector_lbl=local.svc_selector_label label belong to this service.
    selector = var.service.selector
    session_affinity = var.service.session_affinity
    dynamic "port" {
      for_each = var.service.ports
      iterator = it
      content {
        name = it.value["name"]
        port = it.value["service_port"]
        target_port = it.value["target_port"]
        node_port = it.value["node_port"]
        protocol = it.value["protocol"]
      }
    }
    type = var.service.type
    cluster_ip = "None"  # Headless Service.
    publish_not_ready_addresses = var.service.publish_not_ready_addresses
  }
}
