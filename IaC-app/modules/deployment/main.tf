/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
-------------------------------------------------------
Define input variables to the module.
***/
variable app_version {
  type = string
}
variable build_image {
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
variable cr_login_server {
  default = ""
  type = string
}
variable cr_password {
  default = ""
  type = string
  sensitive = true
}
variable cr_username {
  default = ""
  type = string
  sensitive = true
}
variable deployment_name {
  type = string
}
variable dir_path {
  type = string
}
variable dockerfile_name {
  default = "Dockerfile-prod"
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
/***
https://registry.terraform.io/providers/hashicorp/kubernetes/2.25.2/docs/resources/deployment#min_ready_seconds-1
Until the pod is available, the rollout process will not continue. A pod is ready when readiness
probes of all its containers return a success. If a new pod isn't functioning properly and its
readiness probe starts failing before minReadySeconds have passed, the rollout of the new version
will effectively be blocked. Set minReadySeconds to some high value thereby ensuring pods keep
reporting they're ready after they've already started receiving actual traffic.

If there is no explicit readiness probe defined, the container and the pod are always considered
ready, even if the app is not truly ready or is returning errors.
***/
variable min_ready_seconds {
  default = 0
  type = number
}
variable namespace {
  default = "default"
  type = string
}
variable persistent_volume_claims {
  default = []
  type = list(object({
    name = string
    namespace = string
    labels = optional(map(string), {})
    /***
    ReadWriteOnce (RWO) - Only a single NODE can mount the volume for reading and writing.
    ReadOnlyMany (ROX) - Multiple NODES can mount the volume for reading.
    ReadWriteMany (RWX) - Multiple NODES can mount the volume for both reading and writing.
    ***/
    access_modes = list(string)
    /***
    A volumeMode of Filesystem presents a volume as a directory within the Pod's filesystem while
    a volumeMode of Block presents it as a raw block storage device. Filesystem is the default
    and usually preferred mode, enabling standard file system operations on the volume. Block
    mode is used for applications that need direct access to the block device, like databases
    requiring low-latency access.
    ***/
    volume_mode = optional(string, "Filesystem")
    storage_size = string
    /***
    By specifying an empty string ("") as the storage class name, the PVC binds to a
    pre-provisioned PV instead of dynamically provisioning a new one.
    ***/
    storage_class_name = optional(string)
  }))
}
variable pod {
  type = object({
    affinity = optional(object({
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
    }), {})
    automount_service_account_token = optional(bool, false)
    container = optional(list(object({
      name = string
      args = optional(list(string), [])
      # When defined, it overrides the image's default command.
      command = optional(list(string))
      env = optional(map(any), {})
      /***
      A Pod can use environment variables to expose information about itself to containers running in the
      Pod.
      ***/
      env_field = optional(list(object({
        name = string
        field_path = string
      })), [])
      env_from_secrets = optional(list(string), [])
      # image = string jct
      /***
      Be aware that the default imagePullPolicy depends on the image tag. If a container refers to the
      latest tag (either explicitly or by not specifying the tag at all), imagePullPolicy defaults to
      Always, but if the container refers to any other tag, the policy defaults to IfNotPresent.

      When using a tag other than latest, the imagePullPolicy property must be set if changes are made
      to an image without changing the tag. Better yet, always push changes to an image under a new
      tag.
      ***/
      image_pull_policy = optional(string, "Always")
      liveness_probe = optional(list(object({
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
      })), [])
      readiness_probe = optional(list(object({
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
      })), [])
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
      resources = optional(object({
        requests_cpu = optional(string)
        requests_memory = optional(string)
        limits_cpu = optional(string)
        limits_memory = optional(string)
      }), {})
      # https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/pod#nested-schema-for-speccontainersecurity_context
      security_context = optional(object({
        allow_privilege_escalation = optional(bool)
        capabilities = optional(object({
          add = optional(list(string))
          drop = optional(list(string))
        }))
        privileged = optional(bool)
        /***
        For security reasons, you want to prevent processes running in a container from writing to
        the container's filesystem. If you make the container's filesystem read-only, you will need
        to mount a volume in every directory the app writes information; e.g., logs.
        ***/
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
      /***
      In Linux when a filesystem is mounted into a non-empty directory, the directory will only contain
      the files from the newly mounted filesystem. The files in the original directory are inaccessible
      for as long as the filesystem is mounted. In cases when the original directory contains crucial
      files, mounting a volume could break the container. To overcome this limitation, K8s provides an
      additional subPath property on the volumeMount; this property mounts a single file or a single
      directory from the volume instead of mounting the whole volume, and it does not hide the existing
      files in the original directory.
      ***/
      volume_mounts = optional(list(object({
        name = string
        mount_path = string
        sub_path = optional(string)
        read_only = optional(bool)
      })), [])
    })), [])
    init_container = optional(list(object({
      name = string
      args = optional(list(string), [])
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
        /***
        For security reasons, you want to prevent processes running in a container from writing to
        the container's filesystem. If you make the container's filesystem read-only, you will need
        to mount a volume in every directory the app writes information; e.g., logs.
        ***/
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
    })), [])
    labels = optional(map(string), {})
    # https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/pod#nested-schema-for-specsecurity_context
    security_context = optional(object({  # spec.securityContext
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
    }), {})
    termination_grace_period_seconds = optional(number, 30)
    volume_config_map = optional(list(object({
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
    })), [])
    # A temporary directory that shares a pod's lifetime.
    volume_empty_dir = optional(list(object({
      name = string
      medium = optional(string)
      size_limit = optional(string)
    })), [])
    volume_secrets = optional(list(object({
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
    })), [])
    volume_pv = optional(list(object({
      name = string
      claim_name = string
    })), [])
  })
}
# https://registry.terraform.io/providers/hashicorp/kubernetes/2.25.2/docs/resources/deployment#progress_deadline_seconds-1
variable progress_deadline_seconds {
  default = 600
  type = number
}
variable replicas {
  default = 1
  type = number
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
    /***
    A RoleBinding always references a single Role, but it can bind the Role to multiple subjects.
    The Role to bind Subjects to.
    ***/
    role_ref = object({
      kind = string
      # 'name' must match the name of the Role or ClusterRole you wish to bind to.
      name = string
      /***
      The API group to drive authorization decisions. This value must be and defaults to
      'rbac.authorization.k8s.io'.
      ***/
      api_group = string
    })
    /***
    The Users, Groups, or ServiceAccounts to grand permissions to.
    More than one 'subject' is allowed.
    ***/
    subjects = list(object({
      # The type of binding to use. This value must be ServiceAccount, User or Group.
      kind = string
      /***
      The name of this Role to bind Subjects to.
      The name field is case sensitive.
      ***/
      name = string
      /***
      Namespace defines the namespace of the ServiceAccount to bind to. This value only applies
      to kind ServiceAccount.
      ***/
      namespace = optional(string)
      /***
      The API group to drive authorization decisions. This value only applies to kind User and
      Group. It must be 'rbac.authorization.k8s.io'.
      ***/
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
    # Only apply to types ClusterIP, NodePort, and LoadBalancer.
    # selector = map(string)
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
# https://registry.terraform.io/providers/hashicorp/kubernetes/2.25.2/docs/resources/deployment#strategy-1
variable strategy {
  default = null
  type = object({
    type = string
    max_surge = number
    max_unavailable = number
  })
}

/***
Define local variables.
***/
locals {
  pod_selector_label = "ps-${var.deployment_name}"
  service_selector_label = "svc-${var.deployment_name}"
  tmp_pod_lbls = {
    # It must match the label selector of spec.selector.match_labels.
    pod_selector_lbl = local.pod_selector_label
    # It must match the label selector of the Service.
    # Only apply to types ClusterIP, NodePort, and LoadBalancer.
    svc_selector_lbl = local.service_selector_label
  }
  service_selector = {
    svc_selector_lbl = local.service_selector_label
  }
  pod_labels = merge(var.pod.labels, local.tmp_pod_lbls)
  #
  image_tag = (
    var.build_image == true ?
    "${var.cr_login_server}/${var.cr_username}/${var.deployment_name}:${var.app_version}" :
    var.image_tag
  )
}

/***
Build the Docker image.
Use null_resource to create Terraform resources that do not have any particular resourse type.
Use local-exec to invoke commands on the local workstation.
Use timestamp to force the Docker image to build.
***/
resource "null_resource" "docker_build" {
  count = var.build_image == true ? 1 : 0
  triggers = {
    always_run = timestamp()
  }
  #
  provisioner "local-exec" {
    # Multi-platform build.
    # command = "docker build --platform linux/amd64,linux/arm64 --tag ${local.image_tag} --file ${var.dir_path}/${var.dockerfile_name} ${var.dir_path}"
    command = "docker build --platform linux/arm64 --tag ${local.image_tag} --file ${var.dir_path}/${var.dockerfile_name} ${var.dir_path}"
  }
}

/***
Login to the Container Registry.
***/
resource "null_resource" "docker_login" {
  count = var.build_image == true ? 1 : 0
  depends_on = [
    null_resource.docker_build
  ]
  triggers = {
    always_run = timestamp()
  }
  #
  provisioner "local-exec" {
    command = "docker login ${var.cr_login_server} --username ${var.cr_username} --password ${var.cr_password}"
  }
}

/***
Push the image to the Container Registry.
***/
resource "null_resource" "docker_push" {
  count = var.build_image == true ? 1 : 0
  depends_on = [
    null_resource.docker_login
  ]
  triggers = {
    always_run = timestamp()
  }
  #
  provisioner "local-exec" {
    command = "docker push ${local.image_tag}"
  }
}

/***
(1) The maximum size of a Secret is limited to 1MB.
(2) K8s helps keep Secrets safe by making sure each Secret is only distributed to the nodes that
    run the pods that need access to the Secret.
(3) On the nodes, Secrets are always stored in memory and never written to physical storage. (The
    secret volume uses an in-memory filesystem (tmpfs) for the Secret files.)
(4) From K8s version 1.7, etcd stores Secrets in encrypted form.
(5) A Secret's entries can contain binary values, not only plain-text. Base64 encoding allows you
    to include the binary data in YAML or JSON, which are both plaint-text formats.
(6) Even though Secrets can be exposed through environment variables, you may want to avoid doing
    so because they may get exposed inadvertently. For example,
    *	Apps usually dump environment variables in error reports or even write them to the app log at
      startup.
    *	Children processes inherit all the environment variables of the parent process thereby you
      have no way of knowing what happens with your secret data.
    To be safe, always use secret volumes for exposing Secrets.
***/
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
    namespace = var.persistent_volume_claims[count.index].namespace
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
    /***
    If a value for storageClassName isn't explicitly specify, the cluster's default storage class
    is used.
    ***/
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
    namespace = var.config_map[count.index].namespace
    labels = var.config_map[count.index].labels
  }
  data = var.config_map[count.index].data
  binary_data = var.config_map[count.index].binary_data
  immutable = var.config_map[count.index].immutable
}

/***
The random resources generate randomness only when they are created; the results produced are
stored in the Terraform state and re-used until the inputs change, prompting the resource to be
recreated. The resources all provide a map argument called 'keepers' that can be populated with
arbitrary key/value pairs that should be selected such that they remain the same until new random
values are desired.
***/
resource "random_string" "unique_string" {
  length = 5
  special = false  # Exclude special characters.
  upper = true  # Include uppercase letters.
  lower = true  # Include lowercase letters.
  numeric = true # Include numbers.
}

# Deployment -> Stateless.
resource "kubernetes_deployment" "stateless" {
  depends_on = [
    null_resource.docker_push
  ]
  metadata {
    name = var.deployment_name
    namespace = var.namespace
    # Labels attach to the Deployment.
    labels = var.labels
  }
  # The Deployment's specification.
  spec {
    min_ready_seconds = var.min_ready_seconds
    progress_deadline_seconds = var.progress_deadline_seconds
    dynamic "strategy" {
      for_each = var.strategy == null ? [] : [1]
      content {
        type = var.strategy.type
        rolling_update {
          max_surge = var.strategy.max_surge
          max_unavailable = var.strategy.max_unavailable
        }
      }
    }
    # The desired number of pods that should be running.
    replicas = var.replicas
    # revision_history_limit = var.revision_history_limit
    # The label selector determines the pods the ReplicaSet manages.
    selector {
      match_labels = {
        # It must match the labels in the Pod template (spec.template.metadata.labels).
        pod_selector_lbl = local.pod_selector_label
      }
    }
    # The Pod template.
    template {
      metadata {
        /***
        Labels attach to the Pod.
        The pod-template-hash label is added by the Deployment controller to every ReplicaSet
        that a Deployment creates or adopts.
        ***/
        labels = local.pod_labels
      }
      # The Pod template's specification.
      spec {
        dynamic "affinity" {
          for_each = var.pod.affinity == {} ? [] : [1]
          content {
            dynamic "pod_anti_affinity" {
              for_each = var.pod.affinity.pod_anti_affinity == {} ? [] : [1]
              content {
                dynamic "required_during_scheduling_ignored_during_execution" {
                  for_each = var.pod.affinity.pod_anti_affinity.required_during_scheduling_ignored_during_execution
                  iterator = it
                  content {
                    label_selector {
                      match_labels = it.value.label_selector.match_labels
                      dynamic "match_expressions" {
                        for_each = it.value.label_selector.match_expressions
                        iterator = it1
                        content {
                          key = it1.value.key
                          operator = it1.value.operator
                          values = it1.value.values
                        }
                      }
                    }
                    topology_key = "kubernetes.io/hostname"
                  }
                }
              }
            }
            dynamic "pod_affinity" {
              for_each = var.pod.affinity.pod_affinity == {} ? [] : [1]
              content {
                dynamic "required_during_scheduling_ignored_during_execution" {
                  for_each = var.pod.affinity.pod_affinity.required_during_scheduling_ignored_during_execution
                  iterator = it
                  content {
                    label_selector {
                      match_labels = it.value.label_selector.match_labels
                      dynamic "match_expressions" {
                        for_each = it.value.label_selector.match_expressions
                        iterator = it1
                        content {
                          key = it1.value.key
                          operator = it1.value.operator
                          values = it1.value.values
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
        termination_grace_period_seconds = var.pod.termination_grace_period_seconds
        image_pull_secrets {
          name = kubernetes_secret.secrets[0].metadata[0].name  # registry-credentials
        }
        service_account_name = var.service_account == null ? "default" : var.service_account.name
        /***
        In version 1.6+, you can opt out of automounting API credentials for a service account by
        setting 'automountServiceAccountToken: false' on the service account.
        In version 1.6+, you can also opt out of automounting API credentials for a particular
        pod.
        The pod spec takes precedence over the service account if both specify an
        'automountServiceAccountToken' value.
        ***/
        automount_service_account_token = var.pod.automount_service_account_token
        /***
        Security context options at the pod level serve as a default for all the pod's containers
        but can be overridden at the container level.
        ***/
        dynamic "security_context" {
          for_each = var.pod.security_context == {} ? [] : [1]
          content {
            /***
            Set the group that owns the pod volumes. This group will be used by K8s to change the
            permissions of all files/directories in the volumes, when the volumes are mounted by
            a pod.
            ***/
            fs_group = var.pod.security_context.fs_group
            /***
            By default, Kubernetes recursively changes ownership and permissions for the contents
            of each volume to match the fsGroup specified in a Pod's securityContext when that
            volume is mounted. For large volumes, checking and changing ownership and permissions
            can take a lot of time, slowing Pod startup. You can use the fsGroupChangePolicy
            field inside a securityContext to control the way that Kubernetes checks and manages
            ownership and permissions for a volume.
            ***/
            fs_group_change_policy = var.pod.security_context.fs_group_change_policy
            run_as_group = var.pod.security_context.run_as_group
            run_as_non_root = var.pod.security_context.run_as_non_root
            run_as_user = var.pod.security_context.run_as_user
            dynamic "se_linux_options" {
              for_each = var.pod.security_context.se_linux_options == null ? [] : [1]
              content {
                user = var.pod.security_context.se_linux_options.user
                role = var.pod.security_context.se_linux_options.role
                type = var.pod.security_context.se_linux_options.type
                level = var.pod.security_context.se_linux_options.level
              }
            }
            dynamic "seccomp_profile" {
              for_each = var.pod.security_context.seccomp_profile == null ? [] : [1]
              content {
                type = var.pod.security_context.seccomp_profile.type
                localhost_profile = var.pod.security_context.seccomp_profile.localhost_profile
              }
            }
            supplemental_groups = var.pod.security_context.supplemental_groups
            dynamic "sysctl" {
              for_each = var.pod.security_context.sysctl
              iterator = it
              content {
                name = it.name
                value = it.value
              }
            }
            dynamic "windows_options" {
              for_each = var.pod.security_context.windows_options == null ? [] : [1]
              content {
                gmsa_credential_spec = var.pod.security_context.windows_options.gmsa_credential_spec
                gmsa_credential_spec_name = var.pod.security_context.windows_options.gmsa_credential_spec_name
                host_process = var.pod.security_context.windows_options.host_process
                run_as_username = var.pod.security_context.windows_options.run_as_username
              }
            }
          }
        }
        # These containers are run during pod initialization.
        dynamic "init_container" {
          for_each = var.pod.init_container
          iterator = it
          content {
            name = it.value.name
            args = it.value.args
            command = it.value.command
            dynamic "env" {
              for_each = it.value.env
              content {
                name = env.key
                value = env.value
              }
            }
            dynamic "env_from" {
              for_each = it.value.env_from_secrets
              iterator = it1
              content {
                secret_ref {
                  name = it1.value
                }
              }
            }
            image = it.value.image
            image_pull_policy = it.value.image_pull_policy
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
        dynamic "container" {
          for_each = var.pod.container
          iterator = it
          content {
            name = it.value.name
            args = it.value.args
            command = it.value.command
            /***
            To list all of the environment variables:
            Linux: $ printenv
            ***/
            dynamic "env" {
              for_each = it.value.env
              content {
                name = env.key
                value = env.value
              }
            }
            dynamic "env" {
              for_each = it.value.env_field
              content {
                name = env.value.env_name
                value_from {
                  field_ref {
                    field_path = env.value.field_path
                  }
                }
              }
            }
            dynamic "env_from" {
              for_each = var.config_map
              content {
                config_map_ref {
                  name = env_from.value.name
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
              for_each = it.value.env_from_secrets
              iterator = it1
              content {
                secret_ref {
                  name = it1.value
                }
              }
            }
            image = local.image_tag  #jct
            image_pull_policy = it.value.image_pull_policy
            /***
            Security settings that you specify for a container apply only to the individual
            container, and they override settings made at the Pod level when there is overlap.
            Container settings do not affect the Pod's Volumes.
            ***/
            dynamic "security_context" {
              for_each = it.value.security_context == {} ? [] : [1]
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
              iterator = it1
              content {
                name = it1.value.name
                container_port = it1.value.target_port  # The port the app is listening.
                protocol = it1.value.protocol
              }
            }
            dynamic "liveness_probe" {
              for_each = it.value.liveness_probe
              iterator = it1
              content {
                initial_delay_seconds = it1.value.initial_delay_seconds
                period_seconds = it1.value.period_seconds
                timeout_seconds = it1.value.timeout_seconds
                failure_threshold = it1.value.failure_threshold
                success_threshold = it1.value.success_threshold
                /***
                K8s can probe a container using one of the three probes:
                The HTTP GET probe performs an HTTP GET request on the container. If the probe
                receives a response that doesn't represent an error (HTTP response code is 2xx or
                3xx), the probe is considered successful. If the server returns an error response
                code or it doesn't respond at all, the probe is considered a failure and the
                container will be restarted as a result.
                ***/
                dynamic "http_get" {
                  for_each = it1.value.http_get
                  iterator = it2
                  content {
                    host = it2.value.host
                    path = it2.value.path
                    port = it2.value.port
                    scheme = it2.value.scheme
                    dynamic "http_header" {
                      for_each = it2.value.http_header
                      iterator = it3
                      content {
                        name = it3.value.name
                        value = it3.value.value
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
                  for_each = it1.value.exec != null ? [it1.value.exec] : []
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
                  for_each = it1.value.tcp_socket != null ? [it1.value.tcp_socket] : []
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
              for_each = it.value.readiness_probe
              iterator = it1
              content {
                initial_delay_seconds = it1.value.initial_delay_seconds
                period_seconds = it1.value.period_seconds
                timeout_seconds = it1.value.timeout_seconds
                failure_threshold = it1.value.failure_threshold
                success_threshold = it1.value.success_threshold
                /***
                K8s can probe a container using one of the three probes:
                The HTTP GET probe sends an HTTP GET request to the container, and the HTTP status
                code of the response determines whether the container is ready or not.
                ***/
                dynamic "http_get" {
                  for_each = it1.value.http_get
                  iterator = it2
                  content {
                    host = it2.value.host
                    path = it2.value.path
                    port = it2.value.port
                    scheme = it2.value.scheme
                    dynamic "http_header" {
                      for_each = it2.value.http_header
                      iterator = it3
                      content {
                        name = it3.value.name
                        value = it3.value.value
                      }
                    }
                  }
                }
                /***
                The Exec probe executes a process. The container's status is determined by the
                process' exit status code.
                ***/
                dynamic "exec" {
                  for_each = it1.value.exec != null ? [it1.value.exec] : []
                  content {
                    command = exec.value.command
                  }
                }
                /***
                The TCP Socket probe opens a TCP connection to a specified port of the container.
                If the connection is established, the container is considered ready.
                ***/
                dynamic "tcp_socket" {
                  for_each = it1.value.tcp_socket != null ? [it1.value.tcp_socket] : []
                  content {
                    port = tcp_socket.value.port
                  }
                }
              }
            }
            dynamic "resources" {
              for_each = it.value.resources == {} ? [] : [1]
              content {
                requests = {
                  cpu = it.value.resources.requests_cpu
                  memory = it.value.resources.requests_memory
                }
                limits = {
                  cpu = it.value.resources.limits_cpu
                  memory = it.value.resources.limits_memory
                }
              }
            }
            dynamic "volume_mount" {
              for_each = it.value.volume_mounts
              iterator = it1
              content {
                name = it1.value.name
                mount_path = it1.value.mount_path
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
                sub_path = it1.value.sub_path
                read_only = it1.value.read_only
              }
            }
          }
        }
        /***
        Set volumes at the Pod level, then mount them into containers inside that Pod.

        By default, K8s emptyDir volumes are created with root:root ownership and 750
        permissions. This means that the directory created by K8s for the emptyDir volume is
        owned by the root user and group, which translates to read-write-execute permissions for
        the owner (root), read-execute permissions for the group, and no permissions for others.
        (For directories, execute permission is required to access the contents of the directory.)
        In many cases, especially when running containers as non-root users, this default
        ownership can lead to permission issues when containers try to write to the emptyDir
        volume. To address this, you might need to adjust the ownership and permissions of the
        emptyDir volume or consider using other volume types or approaches.
        ***/
        dynamic "volume" {
          for_each = var.pod.volume_empty_dir
          iterator = it
          content {
            name = it.value.name
            empty_dir {
              medium = it.value.medium
              size_limit = it.value.size_limit
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

        While a ConfigMap is typically mounted directly as a read-only volume, it is possible to
        make its contents writable within a Pod by copying them into an emptyDir volume using an
        initContainer. This is often necessary for applications that require modifying their
        configuration files during runtime.
        ***/
        dynamic "volume" {
          for_each = var.pod.volume_config_map
          iterator = it
          content {
            name = it.value.name
            config_map {
              name = it.value.config_map_name
              default_mode = it.value.default_mode
              dynamic "items" {
                for_each = it.value.items
                iterator = it1
                content {
                  key = it1.value.key
                  path = it1.value.path
                }
              }
            }
          }
        }
        dynamic "volume" {
          for_each = var.pod.volume_secrets
          iterator = it
          content {
            name = it.value.name
            secret {
              secret_name = it.value.secret_name
              default_mode = it.value.default_mode
              dynamic "items" {
                for_each = it.value.items
                iterator = it1
                content {
                  key = it1.value.key
                  path = it1.value.path
                }
              }
            }
          }
        }
        /***
        Pods access storage by using the claim as a volume. Claims must exist in the same
        namespace as the Pod using the claim. The cluster finds the claim in the Pod's namespace
        and uses it to get the PersistentVolume backing the claim. The volume is then mounted to
        the host and into the Pod.
        ***/
        dynamic "volume" {
          for_each = var.pod.volume_pv
          iterator = it
          content {
            name = it.value.name
            persistent_volume_claim {
              claim_name = it.value.claim_name
            }
          }
        }
      }
    }
  }
}

/***
Declare a K8s service to create a DNS record to make the microservice accessible within the
cluster.
***/
resource "kubernetes_service" "service" {
  metadata {
    name = var.service.name
    namespace = var.service.namespace
    labels = var.service.labels
  }
  #
  spec {
    # The label selector determines which pods belong to the service.
    selector = local.service_selector
    session_affinity = var.service.session_affinity
    dynamic "port" {
      for_each = var.service.ports
      iterator = it
      content {
        name = it.value.name
        port = it.value.service_port
        target_port = it.value.target_port
        node_port = it.value.node_port
        protocol = it.value.protocol
      }
    }
    type = var.service.type
  }
}
