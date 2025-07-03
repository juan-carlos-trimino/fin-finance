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
variable labels {
  default = {}
  type = map(string)
}
variable image_tag {
  default = ""
  type = string
}
variable namespace {
  default = "default"
  type = string
}
/***
Be aware that the default imagePullPolicy depends on the image tag. If a container refers to the
latest tag (either explicitly or by not specifying the tag at all), imagePullPolicy defaults to
Always, but if the container refers to any other tag, the policy defaults to IfNotPresent.

When using a tag other that latest, the imagePullPolicy property must be set if changes are made
to an image without changing the tag. Better yet, always push changes to an image under a new
tag.
***/
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
variable env_secret {
  default = []
  type = list(object({
    env_name = string
    secret_name = string
    secret_key = string
  }))
}
variable env_field {
  default = []
  type = list(object({
    env_name = string
    field_path = string
  }))
}
/***
Quality of Service (QoS) classes for pods:
(1) BestEffort (lowest priority) - It's assigned to pods that do not have any requests or limits
    set at all (in any of their containers).
(2) Burstable - Pods have some lower-bound resource guarantees based on the request, but do not
    require a specific limit. A Pod is given a QoS class of Burstable if:
    * The Pod does not meet the criteria for QoS class Guaranteed.
    * At least one Container in the Pod has a memory or CPU request or limit.
(3) Guaranteed (highest priority) - It's assigned to pods whose containers' requests are equal to
    the limits for all resources (for each container in the pod). For a pod's class to be
    Guaranteed, three things need to be true:
    * Requests and limits need to be set for both CPU and memory.
    * They need to be set for each container.
    * They need to be equal; the limit needs to match the request for each resource in each
      container.
If a Container specifies its own memory limit, but does not specify a memory request, Kubernetes
automatically assigns a memory request that matches the limit. Similarly, if a Container
specifies its own CPU limit, but does not specify a CPU request, Kubernetes automatically assigns
a CPU request that matches the limit.
***/
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
/***
The ServiceType allows to specify what kind of Service to use: ClusterIP (default),
NodePort, LoadBalancer, and ExternalName.
***/
variable service_type {
  default = "None"  # Headless service.
  type = string
}
variable secrets {
  default = []
  type = list(object({
    name = string
    labels = optional(map(string), {})
    annotations = optional(map(string), {})
    data = optional(map(string), {})
    # base64 encoding.
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
/***
The service normally forwards each connection to a randomly selected backing pod. To ensure that
connections from a particular client are passed to the same Pod each time, set the service's
sessionAffinity property to ClientIP instead of None (default).
Session affinity and Web Browsers (for LoadBalancer Services)
Since the service is now exposed externally, accessing it with a web browser will hit the same
pod every time. If the sessionAffinity is set to None, then why? The browser is using keep-alive
connections and sends all its requests through a single connection. Services work at the
connection level, and when a connection to a service is initially open, a random pod is selected
and then all network packets belonging to that connection are sent to that single pod. Even with
the sessionAffinity set to None, the same pod will always get hit (until the connection is
closed).
***/
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
/***
To relax the StatefulSet ordering guarantee while preserving its uniqueness and identity
guarantee.
***/
variable pod_management_policy {
  default = "OrderedReady"
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
variable config_map {
  default = []
  type = list(object({
    name = string
    labels = optional(map(string), {})
    # Binary data need to be base64 encoded.
    binary_data = optional(map(string), {})
    data = optional(map(string), {})
    immutable = optional(bool, false)
  }))
}
variable persistent_volume_claims {
  default = []
  type = list(object({
    name = string
    labels = optional(map(string), {})
    access_modes = list(string)
    # A volumeMode of Filesystem presents a volume as a directory within the Pod's filesystem while
    # a volumeMode of Block presents it as a raw block storage device. Filesystem is the default
    # and usually preferred mode, enabling standard file system operations on the volume. Block
    # mode is used for applications that need direct access to the block device, like databases
    # requiring low-latency access.
    volume_mode = optional(string)
    storage_size = string
    # By specifying an empty string ("") as the storage class name, the PVC binds to a
    # pre-provisioned PV instead of dynamically provisioning a new one.
    storage_class_name = optional(string)
  }))
}
/***
In Linux when a filesystem is mounted into a non-empty directory, the directory will only contain
the files from the newly mounted filesystem. The files in the original directory are inaccessible
for as long as the filesystem is mounted. In cases when the original directory contains crucial
files, mounting a volume could break the container. To overcome this limitation, K8s provides an
additional subPath property on the volumeMount; this property mounts a single file or a single
directory from the volume instead of mounting the whole volume, and it does not hide the existing
files in the original directory.
***/
variable volume_mount {
  default = []
  type = list(object({
    name = string
    mount_path = string
    sub_path = optional(string)
    read_only = optional(bool)
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
variable volume_config_map {
  default = []
  type = list(object({
    volume_name = string
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
variable volume_pv {  # PersistentVolumeClaim
  default = []
  type = list(object({
    pv_name = string
    claim_name = string
  }))
}
variable volume_claim_template {
  default = []
  type = list(object({
    name = string
    labels = optional(map(string), {})
    access_modes = list(string)
    storage_class_name = string
    storage = string
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

/***
The maximum size of a Secret is limited to 1MB.
K8s helps keep Secrets safe by making sure each Secret is only distributed to the nodes that run
the pods that need access to the Secret.
On the nodes, Secrets are always stored in memory and never written to physical storage. (The
secret volume uses an in-memory filesystem (tmpfs) for the Secret files.)
From K8s version 1.7, etcd stores Secrets in encrypted form.
***/
resource "kubernetes_secret" "secrets" {
  count = length(var.secrets)
  metadata {
    name = var.secrets[count.index].name
    namespace = var.namespace
    labels = var.secrets[count.index].labels
    annotations = var.secrets[count.index].annotations
  }
  # Plain-text data.
  data = var.secrets[count.index].data
  binary_data = var.secrets[count.index].binary_data
  # https://kubernetes.io/docs/concepts/configuration/secret/#secret-types
  # https://kubernetes.io/docs/concepts/configuration/secret/#serviceaccount-token-secrets
  type = var.secrets[count.index].type
}

/***
A ServiceAccount is used by an application running inside a pod to authenticate itself with the
API server. A default ServiceAccount is automatically created for each namespace; each pod is
associated with exactly one ServiceAccount, but multiple pods can use the same ServiceAccount. A
pod can only use a ServiceAccount from the same namespace.
***/
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

/***
Roles define WHAT can be done; role bindings define WHO can do it.
The distinction between a Role/RoleBinding and a ClusterRole/ClusterRoleBinding is that the Role/
RoleBinding is a namespaced resource; ClusterRole/ClusterRoleBinding is a cluster-level resource.
A Role resource defines what actions can be taken on which resources; i.e., which types of HTTP
requests can be performed on which RESTful resources.
***/
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

/***
PersistentVolumeClaims can only be created in a specific namespace; they can then only be used by
pods in the same namespace.
***/
resource "kubernetes_persistent_volume_claim" "pvc" {
  count = length(var.persistent_volume_claims)
  metadata {
    name = var.persistent_volume_claims[count.index].name
    namespace = var.namespace
    labels = var.persistent_volume_claims[count.index].labels
  }
  spec {
    # https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes
    access_modes = var.persistent_volume_claims[count.index].access_modes
    volume_mode = var.persistent_volume_claims[count.index].volume_mode
    resources {
      requests = {
        storage = var.persistent_volume_claims[count.index].storage_size
      }
    }
    # If a value for storageClassName isn't explicitly specify, the cluster's default storage class
    # is used.
    storage_class_name = var.persistent_volume_claims[count.index].storage_class_name
  }
}

/***
The contents of the ConfigMap are passed to containers as either environment variables or as files
in a volume.
***/
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
    replicas = var.replicas
    /***
    The name of the service that governs this StatefulSet.
    This service must exist before the StatefulSet and is responsible for the network identity of
    the set. Pods get DNS/hostnames that follow the pattern:
      pod-name.service-name.namespace.svc.cluster.local.
    ***/
    service_name = local.svc_name
    pod_management_policy = var.pod_management_policy
    /***
    Pod Selector - You must set the .spec.selector field of a StatefulSet to match the labels of
    its .spec.template.metadata.labels. Failing to specify a matching Pod Selector will result in
    a validation error during StatefulSet creation.
    ***/
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
        }
      }
      #
      spec {
        service_account_name = kubernetes_service_account.service_account[0].metadata[0].name
        affinity {
          pod_anti_affinity {
            required_during_scheduling_ignored_during_execution {
              label_selector {
                match_expressions {
                  /***
                  Description of the pod label that determines when the anti-affinity rule
                  applies. Specifies a key and value for the label.
                  key = "rmq_lbl"
                  The operator represents the relationship between the label on the existing
                  pod and the set of values in the matchExpression parameters in the
                  specification for the new pod. Can be In, NotIn, Exists, or DoesNotExist.
                  ***/
                  operator = "In"
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
            for_each = var.env_secret
            content {
              name = env.value["env_name"]
              value_from {
                secret_key_ref {
                  name = env.value["secret_name"]
                  key = env.value["secret_key"]
                }
              }
            }
          }
          dynamic "env" {
            for_each = var.env_field
            content {
              name = env.value["env_name"]
              value_from {
                field_ref {
                  field_path = env.value["field_path"]
                }
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
            name = it.value["volume_name"]
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
        # Pods access storage by using the claim as a volume. Claims must exist in the same
        # namespace as the Pod using the claim. The cluster finds the claim in the Pod's namespace
        # and uses it to get the PersistentVolume backing the claim. The volume is then mounted to
        # the host and into the Pod.
        # dynamic "volume" {
        #   for_each = var.volume_pv
        #   content {
        #     name = volume.value["pv_name"]
        #     persistent_volume_claim {
        #       claim_name = volume.value["claim_name"]
        #     }
        #   }
        # }
      }
    }
    /***
    This template will be used to create a PersistentVolumeClaim for each pod.
    Since PersistentVolumes are cluster-level resources, they do not belong to any namespace, but
    PersistentVolumeClaims can only be created in a specific namespace; they can only be used by
    pods in the same namespace.
    ***/
    dynamic "volume_claim_template" {
      for_each = var.volume_claim_template
      iterator = it
      content {
        metadata {
          name = it.value["name"]
          namespace = var.namespace
          labels = it.value["lables"]
        }
        spec {
          access_modes = it.value["access_modes"]
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
Before deploying a StatefulSet, you will need to create a headless Service, which will be used
to provide the network identity for your stateful pods.
***/
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
