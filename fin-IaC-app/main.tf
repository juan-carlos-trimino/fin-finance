# $ terraform init
# $ terraform apply -var="app_version=1.0.0" -auto-approve
# $ terraform apply -auto-approve
# $ terraform apply -var="app_version=1.0.0" -var="k8s_manifest_crd=false" -auto-approve
# $ terraform apply -var="k8s_manifest_crd=false" -auto-approve
# $ terraform destroy -var="app_version=1.0.0" -auto-approve
# $ terraform destroy -auto-approve
# $ kubectl get all -n finances
locals {
  namespace = kubernetes_namespace.ns.metadata[0].name
  cr_login_server = "docker.io"
  ###########
  # Traefik #
  ###########
  traefik_secret_cert_name = "le-secret-cert"
  issuer_name = "le-acme-issuer"
  tls_store = "default"
  tls_options = "tlsoptions"
  ##################
  # Ingress Routes #
  ##################
  ingress_route = "fin-ingress-route"
  # ingress_route_tcp_rabbitmq = "mem-ingress-route-tcp-rabbitmq"
  ###############
  # Middlewares #
  ###############
  middleware_compress = "fin-mw-compress"
  middleware_rate_limit = "fin-mw-rate-limit"
  middleware_error_page = "fin-mw-error-page"
  middleware_gateway_basic_auth = "fin-mw-gateway-basic-auth"
  middleware_dashboard_basic_auth = "fin-mw-dashboard-basic-auth"
  middleware_security_headers = "fin-mw-security-headers"
  middleware_redirect_https = "fin-mw-redirect-https"
  ####################
  # Name of Services #
  ####################
  svc_finances = "fin-finances"
  svc_error_page = "fin-error-page"
  ############
  # Services #
  ############
  # DNS translates hostnames to IP addresses; the container name is the hostname. When using Docker
  # and Docker Compose, DNS works automatically.
  # In K8s, a service makes the deployment accessible by other containers via DNS.
  svc_dns_error_page = "${local.svc_error_page}.${local.namespace}.svc.cluster.local"
  svc_dns_finances = "${local.svc_finances}.${local.namespace}.svc.cluster.local"
}

###################################################################################################
# traefik                                                                                         #
###################################################################################################
# kubectl get pod,middleware,ingressroute,svc -n finances
# kubectl get all -l "app.kubernetes.io/instance=traefik" -n finances
# kubectl get all -l "app=finances" -n finances
module "traefik" {
  count = var.reverse_proxy ? 1 : 0
  source = "./modules/traefik/traefik"
  app_name = var.app_name
  namespace = local.namespace
  chart_version = "26.1.0"
  api_auth_token = var.traefik_dns_api_token
  timeout = var.helm_traefik_timeout_seconds
  service_name = "fin-traefik"
}

module "middleware-gateway-basic-auth" {
  count = var.reverse_proxy && !var.k8s_manifest_crd ? 1 : 0
  source = "./modules/traefik/middlewares/middleware-gateway-basic-auth"
  app_name = var.app_name
  namespace = local.namespace
  traefik_gateway_username = var.traefik_gateway_username
  traefik_gateway_password = var.traefik_gateway_password
  service_name = local.middleware_gateway_basic_auth
}

module "middleware-dashboard-basic-auth" {
  count = var.reverse_proxy && !var.k8s_manifest_crd ? 1 : 0
  source = "./modules/traefik/middlewares/middleware-dashboard-basic-auth"
  app_name = var.app_name
  namespace = local.namespace
  # While the dashboard in itself is read-only, it is good practice to secure access to it.
  traefik_dashboard_username = var.traefik_dashboard_username
  traefik_dashboard_password = var.traefik_dashboard_password
  service_name = local.middleware_dashboard_basic_auth
}

module "middleware-compress" {
  count = var.reverse_proxy && !var.k8s_manifest_crd ? 1 : 0
  source = "./modules/traefik/middlewares/middleware-compress"
  app_name = var.app_name
  namespace = local.namespace
  service_name = local.middleware_compress
}

module "middleware-rate-limit" {
  count = var.reverse_proxy && !var.k8s_manifest_crd ? 1 : 0
  source = "./modules/traefik/middlewares/middleware-rate-limit"
  app_name = var.app_name
  namespace = local.namespace
  average = 6
  period = "1m"
  burst = 12
  service_name = local.middleware_rate_limit
}

module "middleware-security-headers" {
  count = var.reverse_proxy && !var.k8s_manifest_crd ? 1 : 0
  source = "./modules/traefik/middlewares/middleware-security-headers"
  app_name = var.app_name
  namespace = local.namespace
  service_name = local.middleware_security_headers
}

module "middleware-redirect-https" {
  count = var.reverse_proxy && !var.k8s_manifest_crd ? 1 : 0
  source = "./modules/traefik/middlewares/middleware-redirect-https"
  app_name = var.app_name
  namespace = local.namespace
  service_name = local.middleware_redirect_https
}

module "tlsstore" {
  count = var.reverse_proxy && !var.k8s_manifest_crd ? 1 : 0
  source = "./modules/traefik/tlsstore"
  app_name = var.app_name
  namespace = "default"
  secret_name = local.traefik_secret_cert_name
  service_name = local.tls_store
}

module "tlsoptions" {
  count = var.reverse_proxy && !var.k8s_manifest_crd ? 1 : 0
  source = "./modules/traefik/tlsoptions"
  app_name = var.app_name
  namespace = local.namespace
  service_name = local.tls_options
}

module "ingress-route" {
  count = var.reverse_proxy && !var.k8s_manifest_crd ? 1 : 0
  source = "./modules/traefik/ingress-route"
  app_name = var.app_name
  namespace = local.namespace
  tls_store = local.tls_store
  tls_options = local.tls_options
  middleware_rate_limit = local.middleware_rate_limit
  middleware_compress = local.middleware_compress
  middleware_gateway_basic_auth = local.middleware_gateway_basic_auth
  middleware_dashboard_basic_auth = local.middleware_dashboard_basic_auth
  middleware_security_headers = local.middleware_security_headers
  svc_finances = local.svc_finances
  # svc_gateway = local.svc_gateway
  secret_name = local.traefik_secret_cert_name
  issuer_name = local.issuer_name
  # host_name = "169.46.98.220.nip.io"
  # host_name = "memories.mooo.com"
  host_name = "trimino.xyz"
  service_name = local.ingress_route
}

###################################################################################################
# cert manager                                                                                    #
###################################################################################################
module "cert-manager" {
  count = var.reverse_proxy ? 1 : 0
  source = "./modules/traefik/cert-manager/cert-manager"
  namespace = local.namespace
  chart_version = "1.14.3"
  service_name = "fin-cert-manager"
}

module "acme-issuer" {
  count = var.reverse_proxy && !var.k8s_manifest_crd ? 1 : 0
  source = "./modules/traefik/cert-manager/acme-issuer"
  app_name = var.app_name
  namespace = local.namespace
  issuer_name = local.issuer_name
  acme_email = var.traefik_le_email
  # Let's Encrypt has two different services, one for staging (letsencrypt-staging) and one for
  # production (letsencrypt-prod).
  acme_server = "https://acme-staging-v02.api.letsencrypt.org/directory"
  # acme_server = "https://acme-v02.api.letsencrypt.org/directory"
  dns_names = ["trimino.xyz", "www.trimino.xyz"]
  # Digital Ocean token requires base64 encoding.
  traefik_dns_api_token = var.traefik_dns_api_token
}

module "certificate" {
  count = var.reverse_proxy && !var.k8s_manifest_crd ? 1 : 0
  source = "./modules/traefik/cert-manager/certificates"
  app_name = var.app_name
  namespace = local.namespace
  issuer_name = local.issuer_name
  certificate_name = "le-cert"
  # The A record maps a name to one or more IP addresses when the IP are known and stable.
  # The CNAME record maps a name to another name. It should only be used when there are no other
  # records on that name.
  # common_name = "trimino.xyz"
  dns_names = ["trimino.xyz", "www.trimino.xyz"]
  secret_name = local.traefik_secret_cert_name
}

###################################################################################################
# Application                                                                                     #
###################################################################################################
module "fin-finances" {
  count = var.k8s_manifest_crd ? 0 : 1
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
  #########################################
  # Exposing services to external clients #
  #########################################
  # Use a NodePort service #
  ##########################
  # Setting the service type to NodePort – For a NodePort service, each node in the cluster opens
  # a port on the node itself (the same port number is used across all nodes) and redirects
  # traffic received on that port to the underlying service. The service isn't accessible only at
  # the internal cluster IP and port, but also through a dedicated port on all nodes. Specifying
  # the port isn't mandatory; K8s will choose a random port if it is omitted.
  # Note: By default, the range of the service NodePorts is 30000-32768. This range contains 2768
  # ports, which means that you can create up to 2768 services with NodePorts.
  #
  # For NodePort, it's required to allow communication on ALL protocols in the worker node subnet.
  ports = [{
    name = "ports"
    service_port = 80
    target_port = 8080
    node_port = var.nlb_node_port
    protocol = "TCP"
  }]
  service_type = "NodePort"
  # ports = [{
  #   name = "ports"
  #   service_port = 80
  #   target_port = 8080
  #   protocol = "TCP"
  # }]
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
