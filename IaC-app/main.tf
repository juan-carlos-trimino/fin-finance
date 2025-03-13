#############
# Terraform #
#############
# $ terraform init
# $ terraform apply -var="app_version=1.0.0" -auto-approve
# $ terraform apply -auto-approve
# $ terraform apply -var="app_version=1.0.0" -var="k8s_manifest_crd=false" -auto-approve
# $ terraform apply -var="k8s_manifest_crd=false" -auto-approve
# $ terraform destroy -var="app_version=1.0.0" -auto-approve
# $ terraform destroy -auto-approve
####################
# Kubectl Commands #
####################
###########################
# Troubleshooting Traefik #
###########################
# Execute commands in a running Traefik container.
# $ kubectl exec -it -n finances $(kubectl get pods -n finances --selector "app.kubernetes.io/name=traefik" --output=name) -- /bin/sh
#
# kubectl get pod,middleware,ingressroute,svc -n finances
# kubectl get all -l "app.kubernetes.io/name=traefik" -n finances
# kubectl get all -l "app=finances" -n finances
################################
# Troubleshooting Certificates #
################################
# $ kubectl get svc,pods -n finances
# $ kubectl get Issuers,ClusterIssuers,Certificates,CertificateRequests,Orders,Challenges -n finances
# $ kubectl get Issuers,ClusterIssuers,Certificates,CertificateRequests,Orders,Challenges --all-namespaces
# $ kubectl describe Issuers,ClusterIssuers,Certificates,CertificateRequests,Orders,Challenges -A
# $ kubectl describe Issuers,ClusterIssuers,Certificates,CertificateRequests,Orders,Challenges -n finances
#
# To check the certificate:
# $ kubectl -n finances describe certificate <certificate-name>
# $ kubectl -n finances delete certificate <certificate-name>
#
# To describe a specific resource (the resource name can be obtained from the kubectl get command):
# $ kubectl -n finances describe Issuer <issuer-name>
# $ kubectl get ingressroute -A
# $ kubectl get ingressroute -n finances
#
# To delete a pending Challenge, see here and here. As per documentation, the order is important!!!
# $ kubectl delete Issuer <issuer-name> -n finances
# $ kubectl delete Certificate <certificate-name> -n finances
#
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
  svc_traefik = "fin-traefik"
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
module "traefik" {
  count = var.reverse_proxy ? 1 : 0
  source = "./modules/traefik/traefik"
  app_name = var.app_name
  namespace = local.namespace
  chart_version = "26.1.0"
  api_auth_token = var.traefik_dns_api_token
  timeout = var.helm_traefik_timeout_seconds
  service_name = local.svc_traefik
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
  chart_version = "1.14.4"
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
  region = var.region
  cr_login_server = local.cr_login_server
  cr_username = var.cr_username
  cr_password = var.cr_password
  replicas = 1
  /*** init_container
  init_container = [{
    name = "file-permission"
    image = "busybox:1.34.1"
    image_pull_policy = "IfNotPresent"
    command = [
      "/bin/sh",
      "-c",
      "chown -v -R 1100:1100 /wsf_data_dir"
      # "chmod -R 660 /wsf_data_dir"
    ]
    volume_mounts = [{
      name = "wsf"
      mount_path = "/wsf_data_dir"
      read_only = false
    }]
    security_context = [{
      run_as_non_root = false
      run_as_user = 0
      run_as_group = 0
      read_only_root_filesystem = true
      privileged = true
    }]
  }]
  init_container ***/
  # Limits and requests for CPU resources are measured in millicores. If the container needs one
  # full core to run, use the value '1000m.' If the container only needs 1/4 of a core, use the
  # value of '250m.'
  resources = {  # QoS - Guaranteed
    limits_cpu = "300m"
    limits_memory = "300Mi"
  }
  #### https://kubernetes.io/docs/concepts/configuration/secret/
  # If the order of Secrets changes, the Deployment must be changed accordingly.
  secrets = [{
    name = "${local.svc_finances}-registry-credentials"
    # Plain-text data.
    data = {
      ".dockerconfigjson" = jsonencode({
        auths = {
          "${local.cr_login_server}" = {
            auth = base64encode("${var.cr_username}:${var.cr_password}")
          }
        }
      })
    }
    type = "kubernetes.io/dockerconfigjson"
  },
  /*** s3 storage
  {
    name = "${local.svc_finances}-s3-storage"
    data = {
      obj_storage_ns = var.obj_storage_ns
      region = var.region
      aws_access_key_id = var.aws_access_key_id
      aws_secret_access_key = var.aws_secret_access_key
    }
    type = "Opaque"
  }
  s3 storage ***/
  ]
  # service_account = {
  #   name = "${local.svc_finances}-service-account"
  #   # Note: The keys and the values in the map must be strings. In other words, you cannot use
  #   #       numeric, boolean, list or other types for either the keys or the values.
  #   # https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
  #   # annotations = {
  #   #   "kubernetes.io/enforce-mountable-secrets" = "true"
  #   #   "kubernetes.io/service-account.name" = "${local.svc_finances}-service-account"
  #   # }
  #   automount_service_account_token = true
  #   secret = [{
  #     name = "${local.svc_finances}-secrets"
  #   },
  #   {
  #     name = "${local.svc_finances}-s3-storage"
  #   }]
  # }
  # Configure environment variables specific to the app.
  env = {
    # Set USER to any string to avoid the error:
    # user: Current requires cgo or $USER set in environment
    USER = "wsf-user"
    #
    K8S = true
    HTTP_PORT = "8080"
    SVC_NAME = local.svc_finances
    APP_NAME_VER = "${var.app_name} ${var.app_version}"
    MAX_RETRIES = 20
    SERVER = "http://${local.svc_dns_finances}"
  }
  /*** s3 storage
  env_secret = [{
    env_name = "AWS_SECRET_ACCESS_KEY"
    secret_name = "${local.svc_finances}-s3-storage"
    secret_key = "aws_secret_access_key"
  },
  {
    env_name = "OBJ_STORAGE_NS"
    secret_name = "${local.svc_finances}-s3-storage"
    secret_key = "obj_storage_ns"
  },
  {
    env_name = "AWS_REGION"
    secret_name = "${local.svc_finances}-s3-storage"
    secret_key = "region"
  },
  {
    env_name = "AWS_ACCESS_KEY_ID"
    secret_name = "${local.svc_finances}-s3-storage"
    secret_key = "aws_access_key_id"
  }]
  s3 storage ***/
  /*** env_field
  env_field = [{
    env_name = "POD_ID"
    field_path = "status.podIP"
  }]
  env_field ***/
  /*** NodePort
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
  NodePort ***/
  ports = [{
    name = "ports"
    service_port = 80
    target_port = 8080
    protocol = "TCP"
  }]
  volume_mount = [{
    name = "wsf"
    mount_path = "/wsf_data_dir"
    read_only = false
  }]
  /*** empty_dir
  volume_empty_dir = [{
    name = "wsf"
  }]
  empty_dir ***/
  volume_pv = [{
    pv_name = "wsf"
    claim_name = "finances-pvc"
  }]
  persistent_volume_claims = [{
    pvc_name = "finances-pvc"
    volume_mode = "Filesystem"
    access_modes = ["ReadWriteOnce"]
    storage_size = "2Gi"
    storage_class_name = "oci-bv"
  }]
  pod_security_context = [{
    fs_group = 2200
  }]
  security_context = [{
    run_as_non_root = true
    run_as_user = 1100
    run_as_group = 1100
    read_only_root_filesystem = true
  }]
  # /*** readiness_probe
  readiness_probe = [{
    initial_delay_seconds = 5
    period_seconds = 25
    timeout_seconds = 1
    failure_threshold = 3
    success_threshold = 1
    # http_get = [{
    #   path = "/readiness"
    #   port = 8080
    #   scheme = "HTTP"
    # }]
    tcp_socket = {
      port = 8080
    }
    # exec = {
    #   command = [
    #     "/bin/sh",
    #     "-c",
    #     "ls -al /wsf_data_dir"
    # ]}
  }]
  # readiness_probe **/
  # /*** liveness_probe
  liveness_probe = [{
    initial_delay_seconds = 5
    period_seconds = 20
    timeout_seconds = 1
    failure_threshold = 3
    success_threshold = 1
    # http_get = [{
    #   path = "/liveness"
    #   port = 8080
    #   scheme = "HTTP"
    # }]
    tcp_socket = {
      port = 8080
    }
    # exec = {
    #   command = [
    #     "/bin/sh",
    #     "-c",
    #     "ls -al /wsf_data_dir"
    # ]}
  }]
  # liveness_probe ***/
  service_type = "ClusterIP"
  service_name = local.svc_finances
}
