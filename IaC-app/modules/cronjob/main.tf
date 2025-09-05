/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
-------------------------------------------------------
Define input variables to the module.
***/
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
variable cron_job {
  type = object({
    name = optional(string)
    labels = optional(map(string))
    namespace = optional(string, "default")
    concurrency_policy = optional(string, "Allow")
    failed_jobs_history_limit = optional(number, 1)
    starting_deadline_seconds = optional(number)
    #
    # https://crontab.guru/
    # minute hour day_of_month month day_of_week command_to_execute.
    #                    ┌───────────── minute (0 - 59)
    #                    │ ┌───────────── hour (0 - 23)
    #                    │ │ ┌───────────── day_of_month (1 - 31)
    #                    │ │ │ ┌───────────── month (1 - 12)
    #                    │ │ │ │ ┌───────────── day_of_week (0 - 6) (Sun to Sat)
    #                    │ │ │ │ │              OR sun, mon, tue, wed, thu, fri, sat
    #                    │ │ │ │ │
    #                    │ │ │ │ │
    #                    │ │ │ │ │
    schedule = string  # * * * * *
    successful_jobs_history_limit = optional(number, 3)
    suspend = optional(bool, false)
  })
}
variable env_from_secrets {
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
# https://kubernetes.io/docs/concepts/workloads/controllers/job/
variable job_template {
  type = object({
    # metadata = optional(object({
      name = optional(string)
    #   labels = optional(map(string))
      namespace = optional(string, "default")
    # }), {})
    # pod_metadata = optional(object({
    #   name = optional(string)
    #   labels = optional(map(string))
    # }), {})
    active_deadline_seconds = optional(number)
    backoff_limit = optional(number, 6)
    backoff_limit_per_index = optional(number)
    max_failed_indexes = optional(number)
    completion_mode = optional(string, "NonIndexed")
    completions = optional(number)
    manual_selector = optional(bool)
    parallelism = optional(number)
    pod_failure_policy = optional(list(object({
      rule = list(object({
        action = string
        on_exit_codes = list(object({
          values = list(number)
          container_name = optional(string)
          operator = optional(string)
        }))
        on_pod_condition = list(object({
          status = string
          type = string
        }))
      }))
    })))
    selector = optional(list(object({
      match_expressions = optional(list(object({
        key = string
        operator = string
        values = set(string)
      })))
      match_labels = map(string)
    })))
    ttl_seconds_after_finished = optional(string)
    container = list(object({
      name = string
      args = optional(list(string))
      command = optional(list(string))
      env = optional(map(any))
      env_field = optional(list(object({
        name = string
        field_path = string
      })), [])
      # Passing all entries of a ConfigMap as environment variables at once (envFrom).
      env_from_config_map = optional(list(string), [])
      env_from_secrets = optional(list(string), [])
      image = optional(string)
      image_pull_policy = optional(string, "Always")
      port = optional(object({
        requests_cpu = optional(string)
        requests_memory = optional(string)
        limits_cpu = optional(string)
        limits_memory = optional(string)
      }))
      security_context = optional(object({
        allow_privilege_escalation = optional(bool)
        capabilities = optional(object({
          add = optional(list(string))
          drop = optional(list(string))
        }))
        privileged = optional(bool, false)
        read_only_root_filesystem = optional(bool, false)
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
      }), {})
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
    }))
    restart_policy = optional(string, "Always")
    security_context = optional(object({
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
  })
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
  count = length(var.env_from_secrets)
  metadata {
    name = var.env_from_secrets[count.index].name
    namespace = var.env_from_secrets[count.index].namespace
    labels = var.env_from_secrets[count.index].labels
    annotations = var.env_from_secrets[count.index].annotations
  }
  # Plain-text data.
  data = var.env_from_secrets[count.index].data
  /***
  ***/
  binary_data = var.env_from_secrets[count.index].binary_data
  # https://kubernetes.io/docs/concepts/configuration/secret/#secret-types
  # https://kubernetes.io/docs/concepts/configuration/secret/#serviceaccount-token-secrets
  type = var.env_from_secrets[count.index].type
  immutable = var.env_from_secrets[count.index].immutable
}

# https://registry.terraform.io/providers/hashicorp/kubernetes/1.10.0/docs/resources/cron_job
# https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/cron_job#spec-2
resource "kubernetes_cron_job_v1" "cronjob" {
  metadata {
    name = var.cron_job.name
    labels = var.cron_job.labels
    namespace = var.cron_job.namespace
  }
  spec {
    concurrency_policy = var.cron_job.concurrency_policy
    failed_jobs_history_limit = var.cron_job.failed_jobs_history_limit
    schedule = var.cron_job.schedule
    starting_deadline_seconds = var.cron_job.starting_deadline_seconds
    successful_jobs_history_limit = var.cron_job.successful_jobs_history_limit
    suspend = var.cron_job.suspend
    job_template {  # The pod.
      metadata {
        name = var.job_template.name
        # labels = var.job_template.labels
        namespace = var.job_template.namespace
      }
      spec {
        template {  # Describe the pod.
          metadata {
            # name = var.job_template.pod_metadata.name
            # labels = var.job_template.pod_metadata.labels
          }
          spec {
            dynamic "container" {
              for_each = var.job_template.container
              iterator = it
              content {
                name = it.value.name
                args = it.value.args
                command = it.value.command

                dynamic "env_from" {
                  for_each = it.value.env_from_config_map
                  iterator = it1
                  content {
                    config_map_ref {
                      name = it1.value
                    }
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
                  for_each = it.value.security_context == {} ? [] : [1]
                  content {
                    allow_privilege_escalation = it.value.security_context.allow_privilege_escalation
                    dynamic "capabilities" {
                      for_each = it.value.security_context.capabilities == null ? [] : [1]
                      content {
                        add = it.value.security_context.capabilities.value.add
                        drop = it.value.security_context.capabilities.value.drop
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
                  for_each = it.value.volume_mounts
                  iterator = it1
                  content {
                    name = it1.value.name
                    mount_path = it1.value.mount_path
                    sub_path = it1.value.sub_path
                    read_only = it1.value.read_only
                  }
                }

              }
            }
            restart_policy = var.job_template.restart_policy
            termination_grace_period_seconds = var.job_template.termination_grace_period_seconds
            dynamic "security_context" {
              for_each = var.job_template.security_context == {} ? [] : [1]
              content {
                /***
                Set the group that owns the pod volumes. This group will be used by K8s to change the
                permissions of all files/directories in the volumes, when the volumes are mounted by
                a pod.
                ***/
                fs_group = var.job_template.security_context.fs_group
                /***
                By default, Kubernetes recursively changes ownership and permissions for the contents
                of each volume to match the fsGroup specified in a Pod's securityContext when that
                volume is mounted. For large volumes, checking and changing ownership and permissions
                can take a lot of time, slowing Pod startup. You can use the fsGroupChangePolicy
                field inside a securityContext to control the way that Kubernetes checks and manages
                ownership and permissions for a volume.
                ***/
                fs_group_change_policy = var.job_template.security_context.fs_group_change_policy
                run_as_group = var.job_template.security_context.run_as_group
                run_as_non_root = var.job_template.security_context.run_as_non_root
                run_as_user = var.job_template.security_context.run_as_user
                dynamic "se_linux_options" {
                  for_each = var.job_template.security_context.se_linux_options == null ? [] : [1]
                  content {
                    user = var.job_template.security_context.se_linux_options.user
                    role = var.job_template.security_context.se_linux_options.role
                    type = var.job_template.security_context.se_linux_options.type
                    level = var.job_template.security_context.se_linux_options.level
                  }
                }
                dynamic "seccomp_profile" {
                  for_each = var.job_template.security_context.seccomp_profile == null ? [] : [1]
                  content {
                    type = var.job_template.security_context.seccomp_profile.type
                    localhost_profile = var.job_template.security_context.seccomp_profile.localhost_profile
                  }
                }
                supplemental_groups = var.job_template.security_context.supplemental_groups
                dynamic "sysctl" {
                  for_each = var.job_template.security_context.sysctl
                  iterator = it
                  content {
                    name = it.name
                    value = it.value
                  }
                }
                dynamic "windows_options" {
                  for_each = var.job_template.security_context.windows_options == null ? [] : [1]
                  content {
                    gmsa_credential_spec = var.job_template.security_context.gmsa_credential_spec
                    gmsa_credential_spec_name = var.job_template.security_context.gmsa_credential_spec_name
                    host_process = var.job_template.security_context.host_process
                    run_as_username = var.job_template.security_context.run_as_username
                  }
                }
              }
            }

            dynamic "volume" {
              for_each = var.job_template.volume_config_map
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

          }
        }
      }
    }
  }
}
