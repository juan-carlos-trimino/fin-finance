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
variable dir_path {
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
variable region {
  type = string
  sensitive = true
}
variable dockerfile_name {
  default = "Dockerfile-prod"
  type = string
}
variable cr_login_server {
  type = string
}
variable cr_username {
  type = string
  sensitive = true
}
variable cr_password {
  type = string
  sensitive = true
}
variable readiness_probe {
  default = []
  type = list(object({
    # Number of seconds after the container has started before liveness or readiness probes are
    # initiated. Defaults to 0 seconds. Minimum value is 0.
    initial_delay_seconds = optional(number)
    # How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
    period_seconds = optional(number)
    # Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1.
    timeout_seconds = optional(number)
    # When a probe fails, Kubernetes will try failureThreshold times before giving up. Giving up in
    # case of liveness probe means restarting the container. In case of readiness probe the Pod
    # will be marked Unready. Defaults to 3. Minimum value is 1.
    failure_threshold = optional(number)
    # Minimum consecutive successes for the probe to be considered successful after having failed.
    # Defaults to 1. Must be 1 for liveness and startup Probes. Minimum value is 1.
    success_threshold = optional(number)
    http_get = optional(list(object({
      # Host name to connect to, defaults to the pod IP.
      host = optional(string)
      # Path to access on the HTTP server. Defaults to /.
      path = optional(string)
      # Name or number of the port to access on the container. Number must be in the range 1 to
      # 65535.
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
variable init_container {
  default = []
  type = list(object({
    name = string
    image = string
    image_pull_policy = optional(string)
    command = optional(list(string))
    security_context = optional(list(object({
      run_as_non_root = bool
      run_as_user = number
      run_as_group = number
      read_only_root_filesystem = bool
      privileged = bool
    })), [])
    volume_mounts = optional(list(object({
      name = string
      mount_path = string
      sub_path = optional(string)
      read_only = optional(bool)
    })), [])
  }))
}
# Be aware that the default imagePullPolicy depends on the image tag. If a container refers to the
# latest tag (either explicitly or by not specifying the tag at all), imagePullPolicy defaults to
# Always, but if the container refers to any other tag, the policy defaults to IfNotPresent.
#
# When using a tag other than latest, the imagePullPolicy property must be set if changes are made
# to an image without changing the tag. Better yet, always push changes to an image under a new
# tag.
variable image_pull_policy {
  default = "Always"
  type = string
}
variable pod_security_context {
  default = []
  type = list(object({
    run_as_user = optional(string)
    run_as_group = optional(string)
    fs_group = optional(string)
    fs_group_change_policy = optional(string)
    supplemental_groups = optional(set(number))
  }))
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
    annotations = optional(map(string), {})
    automount_service_account_token = optional(bool, true)
    secret = optional(list(object({
      name = string
    })), [])
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
# In Linux when a filesystem is mounted into a non-empty directory, the directory will only contain
# the files from the newly mounted filesystem. The files in the original directory are inaccessible
# for as long as the filesystem is mounted. In cases when the original directory contains crucial
# files, mounting a volume could break the container. To overcome this limitation, K8s provides an
# additional subPath property on the volumeMount; this property mounts a single file or a single
# directory from the volume instead of mounting the whole volume, and it does not hide the existing
# files in the original directory.
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
variable persistent_volume_claims {
  default = []
  type = list(object({
    pvc_name = string
    # ReadWriteOnce (RWO) - Only a single NODE can mount the volume for reading and writing.
    # ReadOnlyMany (ROX) - Multiple NODES can mount the volume for reading.
    # ReadWriteMany (RWX) - Multiple NODES can mount the volume for both reading and writing.
    access_modes = list(string)
    # Filesystem (default) or Block.
    volume_mode = optional(string)
    storage_size = string
    # By specifying an empty string ("") as the storage class name, the PVC binds to a
    # pre-provisioned PV instead of dynamically provisioning a new one.
    storage_class_name = optional(string)
  }))
}

/***
Define local variables.
***/
locals {
  pod_selector_label = "rs-${var.service_name}"
  svc_selector_label = "svc-${var.service_name}"
  image_tag = (
    var.image_tag == "" ?
    "${var.cr_login_server}/${var.cr_username}/${var.service_name}:${var.app_version}" :
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
  triggers = {
    always_run = timestamp()
  }
  #
  provisioner "local-exec" {
    command = "docker build --platform linux/amd64,linux/arm64 --tag ${local.image_tag} --file ${var.dir_path}/${var.dockerfile_name} ${var.dir_path}"
  }
}

/***
Login to the Container Registry.
***/
resource "null_resource" "docker_login" {
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
    namespace = var.namespace
    labels = {
      app = var.app_name
    }
    annotations = var.service_account.annotations
  }
  # To enable automatic mounting of the service account token; it defaults to true.
  automount_service_account_token = var.service_account.automount_service_account_token
  dynamic "secret" {
    for_each = var.service_account.secret
    content {
      name = secret.value["name"]
    }
  }
}

# PersistentVolumeClaims can only be created in a specific namespace; they can then only be used by
# pods in the same namespace.
resource "kubernetes_persistent_volume_claim" "pvc" {
  count = length(var.persistent_volume_claims)
  metadata {
    name = var.persistent_volume_claims[count.index].pvc_name
    namespace = var.namespace
    labels = {
      app = var.app_name
    }
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
    storage_class_name = var.persistent_volume_claims[count.index].storage_class_name
  }
}

# Deployment -> Stateless.
resource "kubernetes_deployment" "stateless" {
  depends_on = [
    null_resource.docker_push
  ]
  metadata {
    name = var.service_name
    namespace = var.namespace
    # Labels attach to the Deployment.
    labels = {
      app = var.app_name
    }
  }
  # The Deployment's specification.
  spec {
    # The desired number of pods that should be running.
    replicas = var.replicas
    # revision_history_limit = var.revision_history_limit
    # The label selector determines the pods the ReplicaSet manages.
    selector {
      match_labels = {
        # It must match the labels in the Pod template.
        pod_selector_lbl = local.pod_selector_label
      }
    }
    # The Pod template.
    template {
      metadata {
        # Labels attach to the Pod.
        # The pod-template-hash label is added by the Deployment controller to every ReplicaSet
        # that a Deployment creates or adopts.
        labels = {
          app = var.app_name
          # It must match the label selector of the ReplicaSet.
          pod_selector_lbl = local.pod_selector_label
          # It must match the label selector of the Service.
          svc_selector_lbl = local.svc_selector_label
        }
      }
      # The Pod template's specification.
      spec {
        termination_grace_period_seconds = var.termination_grace_period_seconds
        image_pull_secrets {
          name = kubernetes_secret.secrets[0].metadata[0].name  # registry-credentials
        }
        service_account_name = var.service_account == null ? "default" : var.service_account.name
        # Security context options at the pod level serve as a default for all the pod's containers
        # but can be overridden at the container level.
        dynamic "security_context" {
          for_each = var.pod_security_context
          iterator = it
          content {
            run_as_user = it.value["run_as_user"]
            run_as_group = it.value["run_as_group"]
            # Set the group that owns the pod volumes. This group will be used by K8s to change the
            # permissions of all files/directories in the volumes, when the volumes are mounted by
            # a pod.
            fs_group = it.value["fs_group"]
            supplemental_groups = it.value["supplemental_groups"]
            # By default, Kubernetes recursively changes ownership and permissions for the contents
            # of each volume to match the fsGroup specified in a Pod's securityContext when that
            # volume is mounted. For large volumes, checking and changing ownership and permissions
            # can take a lot of time, slowing Pod startup. You can use the fsGroupChangePolicy
            # field inside a securityContext to control the way that Kubernetes checks and manages
            # ownership and permissions for a volume.
            fs_group_change_policy = it.value["fs_group_change_policy"]
          }
        }
        # These containers are run during pod initialization.
        dynamic "init_container" {
          for_each = var.init_container
          iterator = it
          content {
            name = it.value["name"]
            image = it.value["image"]
            image_pull_policy = it.value["image_pull_policy"]
            command = it.value["command"]
            dynamic "security_context" {
              for_each = it.value["security_context"]
              iterator = it1
              content {
                run_as_non_root = it1.value["run_as_non_root"]
                run_as_user = it1.value["run_as_user"]
                run_as_group = it1.value["run_as_group"]
                read_only_root_filesystem = it1.value["read_only_root_filesystem"]
                privileged = it1.value["privileged"]
              }
            }
            dynamic "volume_mount" {
              for_each = it.value["volume_mounts"]
              iterator = it2
              content {
                name = it2.value["name"]
                mount_path = it2.value["mount_path"]
                sub_path = it2.value["sub_path"]
                read_only = it2.value["read_only"]
              }
            }
          }
        }
        container {
          name = var.service_name
          image_pull_policy = var.image_pull_policy
          image = local.image_tag
          # Security settings that you specify for a container apply only to the individual
          # container, and they override settings made at the Pod level when there is overlap.
          # Container settings do not affect the Pod's Volumes.
          dynamic "security_context" {
            for_each = var.security_context
            content {
              run_as_non_root = security_context.value["run_as_non_root"]
              run_as_user = security_context.value["run_as_user"]
              run_as_group = security_context.value["run_as_group"]
              read_only_root_filesystem = security_context.value["read_only_root_filesystem"]
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
              name = port.value["name"]
              container_port = port.value["target_port"]  # The port the app is listening.
              protocol = port.value["protocol"]
            }
          }
          # Liveness probes keep pods healthy by killing unhealthy containers and replacing them
          # with new healthy containers; readiness probes ensure that only pods with containers
          # that are ready to serve requests receive them. Unlike liveness probes, if a container
          # fails the readiness check, it won't be killed or restarted.
          dynamic "readiness_probe" {
            for_each = var.readiness_probe
            content {
              initial_delay_seconds = readiness_probe.value["initial_delay_seconds"]
              period_seconds = readiness_probe.value["period_seconds"]
              timeout_seconds = readiness_probe.value["timeout_seconds"]
              failure_threshold = readiness_probe.value["failure_threshold"]
              success_threshold = readiness_probe.value["success_threshold"]
              # K8s can probe a container using one of the three probes:
              # The HTTP GET probe sends an HTTP GET request to the container, and the HTTP status
              # code of the response determines whether the container is ready or not.
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
              # The Exec probe executes a process. The container's status is determined by the
              # process' exit status code.
              dynamic "exec" {
                # for_each = it.value["exec"] != null ? [it.value["exec"]] : []
                for_each = readiness_probe.value["exec"] != null ? [readiness_probe.value["exec"]] : []
                content {
                  command = exec.value.command
                }
              }
              # The TCP Socket probe opens a TCP connection to a specified port of the container.
              # If the connection is established, the container is considered ready.
              dynamic "tcp_socket" {
                # for_each = it.value["tcp_socket"] != null ? [it.value["tcp_socket"]] : []
                for_each = readiness_probe.value["tcp_socket"] != null ? [readiness_probe.value["tcp_socket"]] : []
                content {
                  port = tcp_socket.value.port
                }
              }
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
              # K8s can probe a container using one of the three probes:
              # The HTTP GET probe performs an HTTP GET request on the container. If the probe
              # receives a response that doesn't represent an error (HTTP response code is 2xx or
              # 3xx), the probe is considered successful. If the server returns an error response
              # code or it doesn't respond at all, the probe is considered a failure and the
              # container will be restarted as a result.
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
              # The Exec probe executes an arbitrary command inside the container and checks the
              # command's exit status code. If the status code is 0, the probe is successful. All
              # other codes are considered failures.
              dynamic "exec" {
                for_each = it.value["exec"] != null ? [it.value["exec"]] : []
                content {
                  command = exec.value.command
                }
              }
              # The TCP Socket probe tries to open a TCP connection to the specified port of the
              # container. If the connection is established successfully, the probe is successful.
              # Otherwise, the container is restarted.
              dynamic "tcp_socket" {
                for_each = it.value["tcp_socket"] != null ? [it.value["tcp_socket"]] : []
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
          # To list all of the environment variables:
          # Linux: $ printenv
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
        # Set volumes at the Pod level, then mount them into containers inside that Pod.
        #
        # By default, K8s emptyDir volumes are created with root:root ownership and 750
        # permissions. This means that the directory created by K8s for the emptyDir volume is
        # owned by the root user and group, which translates to read-write-execute permissions for
        # the owner (root), read-execute permissions for the group, and no permissions for others.
        # (For directories, execute permission is required to access the contents of the
        # directory.)
        # In many cases, especially when running containers as non-root users, this default
        # ownership can lead to permission issues when containers try to write to the emptyDir
        # volume. To address this, you might need to adjust the ownership and permissions of the
        # emptyDir volume or consider using other volume types or approaches.
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
        dynamic "volume" {
          for_each = var.volume_pv
          content {
            name = volume.value["pv_name"]
            persistent_volume_claim {
              claim_name = volume.value["claim_name"]
            }
          }
        }
      }
    }
  }
}

# Declare a K8s service to create a DNS record to make the microservice accessible within the
# cluster.
resource "kubernetes_service" "service" {
  metadata {
    name = var.service_name
    namespace = var.namespace
    labels = {
      app = var.app_name
    }
  }
  #
  spec {
    # The label selector determines which pods belong to the service.
    selector = {
      svc_selector_lbl = local.svc_selector_label
    }
    session_affinity = var.service_session_affinity
    dynamic "port" {
      for_each = var.ports
      iterator = it
      content {
        name = it.value["name"]
        port = it.value["service_port"]
        target_port = it.value["target_port"]
        node_port = it.value["node_port"]
        protocol = it.value["protocol"]
      }
    }
    type = var.service_type
  }
}
