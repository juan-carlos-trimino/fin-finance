
# Execute the following command to clean up all the pending resources and remove all the K8s
# objects: (Caution: The command in this step removes all the resources.)
# $ kubectl delete all --all -n <namespace>
# Execute the following command to delete the namespace:
# $ kubectl delete namespace <namespace>

locals {
  namespace = kubernetes_namespace.ns.metadata[0].name
  cr_login_server = "docker.io"
  ####################
  # Name of Services #
  ####################
  svc_finances = "fin-finances"
  ############
  # Services #
  ############
  svc_dns_finances = "${local.svc_finances}.${local.namespace}.svc.cluster.local"
}

###################################################################################################
# Application                                                                                     #
###################################################################################################
module "fin-finances" {
#  count = var.k8s_manifest_crd ? 0 : 1
  # Specify the location of the module, which contains the file main.tf.
  source = "./modules/deployment"
  dir_path = ".."
  app_name = var.app_name
  app_version = var.app_version
  namespace = local.namespace
  replicas = 1
  qos_limits_cpu = "400m"
  qos_limits_memory = "400Mi"
  cr_login_server = local.cr_login_server
  cr_username = var.cr_username
  cr_password = var.cr_password
  # Configure environment variables specific to the app.
  env = {
    SVC_NAME: local.svc_finances
    APP_NAME_VER: "${var.app_name} ${var.app_version}"
    PORT: "80"
    MAX_RETRIES: 20
    SERVER: "http://${local.svc_dns_finances}"
  }
  # readiness_probe = [{
  #   http_get = [{
  #     path = "/readiness"
  #     port = 0
  #     scheme = "HTTP"
  #   }]
  #   initial_delay_seconds = 30
  #   period_seconds = 20
  #   timeout_seconds = 2
  #   failure_threshold = 4
  #   success_threshold = 1
  # }]
  service_name = local.svc_finances
  service_type = "LoadBalancer"
}
