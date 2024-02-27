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
    HTTP_PORT="8080"
    # Set USER to any string to avoid the error:
    # user: Current requires cgo or $USER set in environment
    USER="wsf-user"
    SVC_NAME: local.svc_finances
    APP_NAME_VER: "${var.app_name} ${var.app_version}"
    MAX_RETRIES: 20
    SERVER: "http://${local.svc_dns_finances}"
  }
  # For NodePort, it's required to allow communication on ALL protocols in the worker node subnet.
  ports = [{
    name = "ports"
    service_port = 80
    target_port = 8080
    node_port = 31600
    protocol = "TCP"
  }]
  service_type = "NodePort"
  # service_type = "LoadBalancer"
  security_context = [{
    run_as_non_root = true
    run_as_user = 1100
    run_as_group = 1100
    read_only_root_filesystem = true
  }]
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
}
