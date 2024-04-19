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
variable dockerfile_name {
  default = "Dockerfile-prod"
  type = string
}
variable dir_path {
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
variable obj_storage {
  default = []
  type = list(object({
    obj_storage_ns = string
    aws_access_key_id = string
    aws_secret_access_key = string
  }))
  sensitive = true
}
variable region {
  type = string
  sensitive = true
}
variable readiness_probe {
  default = []
  type = list(object({
    http_get = list(object({
      # Host name to connect to, defaults to the pod IP.
      #host = string
      # Path to access on the HTTP server. Defaults to /.
      path = string
      # Name or number of the port to access on the container. Number must be in the range 1 to
      # 65535.
      port = number
      # Scheme to use for connecting to the host (HTTP or HTTPS). Defaults to HTTP.
      scheme = string
    }))
    # Number of seconds after the container has started before liveness or readiness probes are
    # initiated. Defaults to 0 seconds. Minimum value is 0.
    initial_delay_seconds = number
    # How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
    period_seconds = number
    # Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1.
    timeout_seconds = number
    # When a probe fails, Kubernetes will try failureThreshold times before giving up. Giving up in
    # case of liveness probe means restarting the container. In case of readiness probe the Pod
    # will be marked Unready. Defaults to 3. Minimum value is 1.
    failure_threshold = number
    # Minimum consecutive successes for the probe to be considered successful after having failed.
    # Defaults to 1. Must be 1 for liveness and startup Probes. Minimum value is 1.
    success_threshold = number
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
  type = map
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
variable pvc_access_modes {
  default = []
  type = list(any)
}
variable pvc_storage_class_name {
  default = ""
  type = string
}
variable pvc_storage_size {
  default = "20Gi"
  type = string
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
    # An array of keys from the ConfigMap to create as files.
    items = optional(list(object({
      # Include the entry under this key.
      key = string
      # The entry's value should be stored in this file.
      path = string
      # Although ConfigMaps should be used for non-sensitive configuration data, you may want to
      # make the file readable and writeble only to the user and group that owned the file; e.g.,
      # default_mode = "6600" (-rw-rw------)
      # The default permission is "6440" (-rw-r--r----)
      default_mode = optional(string)
    })), [])
  }))
}
variable volume_pvc {  # PersistentVolumeClaim
  default = []
  type = list(object({
    volume_name = string
    claim_name = string
  }))
}
variable persistent_volume_claims {
  default = []
  type = list(object({
    name = string
    # ReadWriteOnce (RWO) - Only a single NODE can mount the volume for reading and writing.
    # ReadOnlyMany (ROX) - Multiple NODES can mount the volume for reading.
    # ReadWriteMany (RWX) - Multiple NODES can mount the volume for both reading and writing.
    access_modes = list(string)
    # Filesystem (default) or Block.
    volume_mode = optional(string)
    storage = string
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
    command = "docker build -t ${local.image_tag} --file ${var.dir_path}/${var.dockerfile_name} ${var.dir_path}"
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
    # command = "docker login ${var.cr_login_server} -T -u ${var.cr_username} --password-stdin"
    command = "docker login ${var.cr_login_server} -u ${var.cr_username} -p ${var.cr_password}"
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

resource "kubernetes_secret" "registry_credentials" {
  metadata {
    name = "${var.service_name}-registry-credentials"
    namespace = var.namespace
    labels = {
      app = var.app_name
    }
  }
  data = {
    ".dockerconfigjson" = jsonencode({
      auths = {
        "${var.cr_login_server}" = {
          auth = base64encode("${var.cr_username}:${var.cr_password}")
        }
      }
    })
  }
  type = "kubernetes.io/dockerconfigjson"
}

resource "kubernetes_secret" "obj_storage" {  # For object storage.
  count = length(var.obj_storage)
  metadata {
    name = "${var.service_name}-s3-storage"
    namespace = var.namespace
    labels = {
      app = var.app_name
    }
  }
  # Plain-text data.
  data = {
    obj_storage_ns = var.obj_storage[count.index].obj_storage_ns
    region = var.region
    aws_access_key_id = var.obj_storage[count.index].aws_access_key_id
    aws_secret_access_key = var.obj_storage[count.index].aws_secret_access_key
  }
  type = "Opaque"
}

# PersistentVolumeClaims can only be created in a specific namespace; they can then only be used by
# pods in the same namespace.
resource "kubernetes_persistent_volume_claim" "pvc" {
  count = length(var.persistent_volume_claims)
  metadata {
    name = var.persistent_volume_claims[count.index].name
    namespace = var.namespace
    labels = {
      app = var.app_name
    }
  }
  spec {
    access_modes = var.persistent_volume_claims[count.index].access_modes
    volume_mode = var.persistent_volume_claims[count.index].volume_mode
    resources {
      requests = {
        storage = var.persistent_volume_claims[count.index].storage
      }
    }
    storage_class_name = ""
  }
}

# Declare a K8s deployment to deploy a microservice; it instantiates the container for the
# microservice into the K8s cluster.
resource "kubernetes_deployment" "deployment" {
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
          name = kubernetes_secret.registry_credentials.metadata[0].name
        }
        container {
          name = var.service_name
          image_pull_policy = var.image_pull_policy
          image = local.image_tag
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
          dynamic "readiness_probe" {
            for_each = var.readiness_probe
            content {
              initial_delay_seconds = readiness_probe.value["initial_delay_seconds"]
              period_seconds = readiness_probe.value["period_seconds"]
              timeout_seconds = readiness_probe.value["timeout_seconds"]
              failure_threshold = readiness_probe.value["failure_threshold"]
              success_threshold = readiness_probe.value["success_threshold"]
              dynamic "http_get" {
                for_each = readiness_probe.value.http_get
                content {
                  #host = http_get.value["host"]
                  path = http_get.value["path"]
                  port = http_get.value["port"] != 0 ? http_get.value["port"] : 8080
                  scheme = http_get.value["scheme"]
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
                  # name = kubernetes_secret.obj_storage.metadata[0].name
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
              dynamic "items" {
                for_each = it.value["items"]
                iterator = itn
                content {
                  key = itn.value["key"]
                  path = itn.value["path"]
                  default_mode = itn.value["default_mode"]
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
          for_each = var.volume_pvc
          content {
            name = volume.value["volume_name"]
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
