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
    spec = object({
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
      #                    │ │ │ │ ┌───────────── day_of_week (0 - 6) (Sun to Sat;
      #                    │ │ │ │ │              7 is also Sunday on some systems)
      #                    │ │ │ │ │              OR sun, mon, tue, wed, thu, fri, sat
      #                    │ │ │ │ │
      #                    │ │ │ │ │
      schedule = string  # * * * * *
      successful_jobs_history_limit = optional(number, 3)
      suspend = optional(bool, false)
      job_template = object({
        metadata = object({
          name = optional(string)
          labels = optional(map(string))
          namespace = optional(string, "default")
        })
        spec = object({
          template = object({
            metadata = optional(object({
              name = optional(string)
              labels = optional(map(string))
              namespace = optional(string, "default")
            }))
            spec = object({
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
              }))
            })
          })
        })
      })
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
    concurrency_policy = var.cron_job.spec.concurrency_policy
    failed_jobs_history_limit = var.cron_job.spec.failed_jobs_history_limit
    schedule = var.cron_job.spec.schedule
    starting_deadline_seconds = var.cron_job.spec.starting_deadline_seconds
    successful_jobs_history_limit = var.cron_job.spec.successful_jobs_history_limit
    suspend = var.cron_job.spec.suspend
    job_template {  # The pod.
      metadata {
        name = var.cron_job.spec.job_template.metadata.name
        labels = var.cron_job.spec.job_template.metadata.labels
        namespace = var.cron_job.spec.job_template.metadata.namespace
      }
      spec {
        template {  # Describe the pod.
          metadata {
          }
          spec {
            # restart_policy = var.restart_policy
            dynamic "container" {
              for_each = var.cron_job.spec.job_template.spec.template.spec.container
              iterator = it
              content {
                name = it.value.name
                args = it.value.args
                command = it.value.command
                image = it.value.image
                image_pull_policy = it.value.image_pull_policy
              }
            }
          }
        }
      }
    }
  }
}
