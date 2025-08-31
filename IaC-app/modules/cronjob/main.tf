/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
-------------------------------------------------------
Define input variables to the module.
***/
variable cron_job {
  type = object({
    metadata = object({
      name = optional(string)
      labels = optional(map(string))
      namespace = optional(string, "default")
    })
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
    # https://kubernetes.io/docs/concepts/workloads/controllers/job/
    job_template = object({
      metadata = object({
        name = optional(string)
        labels = optional(map(string))
        namespace = optional(string, "default")
      })
      pod_metadata = optional(object({
        name = optional(string)
        labels = optional(map(string))
      }))
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
        env_from = optional(list(object({
          name = string
          field_path = string
        })))
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
    })
  })
}

# https://registry.terraform.io/providers/hashicorp/kubernetes/1.10.0/docs/resources/cron_job
# https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/cron_job#spec-2
resource "kubernetes_cron_job_v1" "cronjob" {
  metadata {
    name = var.cron_job.metadata.name
    labels = var.cron_job.metadata.labels
    namespace = var.cron_job.metadata.namespace
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
        name = var.cron_job.job_template.metadata.name
        labels = var.cron_job.job_template.metadata.labels
        namespace = var.cron_job.job_template.metadata.namespace
      }
      spec {
        template {  # Describe the pod.
          metadata {
            name = var.cron_job.job_template.pod_metadata.name
            labels = var.cron_job.job_template.pod_metadata.labels
          }
          spec {
            dynamic "container" {
              for_each = var.cron_job.job_template.container
              iterator = it
              content {
                name = it.value.name
                args = it.value.args
                command = it.value.command
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
              }
            }
            restart_policy = var.cron_job.job_template.restart_policy
            dynamic "security_context" {
              for_each = var.cron_job.job_template.security_context == {} ? [] : [1]
              content {
                /***
                Set the group that owns the pod volumes. This group will be used by K8s to change the
                permissions of all files/directories in the volumes, when the volumes are mounted by
                a pod.
                ***/
                fs_group = var.cron_job.job_template.security_context.fs_group
                /***
                By default, Kubernetes recursively changes ownership and permissions for the contents
                of each volume to match the fsGroup specified in a Pod's securityContext when that
                volume is mounted. For large volumes, checking and changing ownership and permissions
                can take a lot of time, slowing Pod startup. You can use the fsGroupChangePolicy
                field inside a securityContext to control the way that Kubernetes checks and manages
                ownership and permissions for a volume.
                ***/
                fs_group_change_policy = var.cron_job.job_template.security_context.fs_group_change_policy
                run_as_group = var.cron_job.job_template.security_context.run_as_group
                run_as_non_root = var.cron_job.job_template.security_context.run_as_non_root
                run_as_user = var.cron_job.job_template.security_context.run_as_user
                dynamic "se_linux_options" {
                  for_each = var.cron_job.job_template.security_context.se_linux_options == null ? [] : [1]
                  content {
                    user = var.cron_job.job_template.security_context.se_linux_options.user
                    role = var.cron_job.job_template.security_context.se_linux_options.role
                    type = var.cron_job.job_template.security_context.se_linux_options.type
                    level = var.cron_job.job_template.security_context.se_linux_options.level
                  }
                }
                dynamic "seccomp_profile" {
                  for_each = var.cron_job.job_template.security_context.seccomp_profile == null ? [] : [1]
                  content {
                    type = var.cron_job.job_template.security_context.seccomp_profile.type
                    localhost_profile = var.cron_job.job_template.security_context.seccomp_profile.localhost_profile
                  }
                }
                supplemental_groups = var.cron_job.job_template.security_context.supplemental_groups
                dynamic "sysctl" {
                  for_each = var.cron_job.job_template.security_context.sysctl
                  iterator = it
                  content {
                    name = it.name
                    value = it.value
                  }
                }
                dynamic "windows_options" {
                  for_each = var.cron_job.job_template.security_context.windows_options == null ? [] : [1]
                  content {
                    gmsa_credential_spec = var.cron_job.job_template.security_context.gmsa_credential_spec
                    gmsa_credential_spec_name = var.cron_job.job_template.security_context.gmsa_credential_spec_name
                    host_process = var.cron_job.job_template.security_context.host_process
                    run_as_username = var.cron_job.job_template.security_context.run_as_username
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
