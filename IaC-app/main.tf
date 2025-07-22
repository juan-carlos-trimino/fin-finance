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
  svc_gateway = "fin-gateway"
  svc_error_page = "fin-error-page"
  svc_traefik = "fin-traefik"
  svc_my_sql = "fin-mysql"
  svc_my_sql_router = "fin-my-sql-router"
  ############
  # Services #
  ############
  # DNS translates hostnames to IP addresses; the container name is the hostname. When using Docker
  # and Docker Compose, DNS works automatically.
  # In K8s, a service makes the deployment accessible by other containers via DNS.
  # FQDN: service-name.namespace.svc.cluster.local
  svc_dns_error_page = "${local.svc_error_page}.${local.namespace}${var.cluster_domain_suffix}"
  svc_dns_finances = "${local.svc_finances}.${local.namespace}${var.cluster_domain_suffix}"
}

###################################################################################################
# traefik                                                                                         #
###################################################################################################
module "traefik" {
  count = var.reverse_proxy ? 1 : 0
  source = "./modules/traefik/traefik"
  app_name = var.app_name
  namespace = local.namespace
  chart_version = "v36.1.0"  # Released: 2025-06-11.
  api_auth_token = var.traefik_dns_api_token
  timeout = var.helm_traefik_timeout_seconds
  service_name = local.svc_traefik
}

module "middleware-gateway-basic-auth" {
  count = var.reverse_proxy && !var.k8s_crds ? 1 : 0
  source = "./modules/traefik/middlewares/middleware-gateway-basic-auth"
  app_name = var.app_name
  namespace = local.namespace
  traefik_gateway_username = var.traefik_gateway_username
  traefik_gateway_password = var.traefik_gateway_password
  service_name = local.middleware_gateway_basic_auth
}

module "middleware-dashboard-basic-auth" {
  count = var.reverse_proxy && !var.k8s_crds ? 1 : 0
  source = "./modules/traefik/middlewares/middleware-dashboard-basic-auth"
  app_name = var.app_name
  namespace = local.namespace
  # While the dashboard in itself is read-only, it is good practice to secure access to it.
  traefik_dashboard_username = var.traefik_dashboard_username
  traefik_dashboard_password = var.traefik_dashboard_password
  service_name = local.middleware_dashboard_basic_auth
}

module "middleware-compress" {
  count = var.reverse_proxy && !var.k8s_crds ? 1 : 0
  source = "./modules/traefik/middlewares/middleware-compress"
  app_name = var.app_name
  namespace = local.namespace
  service_name = local.middleware_compress
}

module "middleware-rate-limit" {
  count = var.reverse_proxy && !var.k8s_crds ? 1 : 0
  source = "./modules/traefik/middlewares/middleware-rate-limit"
  app_name = var.app_name
  namespace = local.namespace
  average = 6
  period = "1m"
  burst = 12
  service_name = local.middleware_rate_limit
}

module "middleware-security-headers" {
  count = var.reverse_proxy && !var.k8s_crds ? 1 : 0
  source = "./modules/traefik/middlewares/middleware-security-headers"
  app_name = var.app_name
  namespace = local.namespace
  service_name = local.middleware_security_headers
}

module "middleware-redirect-https" {
  count = var.reverse_proxy && !var.k8s_crds ? 1 : 0
  source = "./modules/traefik/middlewares/middleware-redirect-https"
  app_name = var.app_name
  namespace = local.namespace
  service_name = local.middleware_redirect_https
}

module "tlsstore" {
  count = var.reverse_proxy && !var.k8s_crds ? 1 : 0
  source = "./modules/traefik/tlsstore"
  app_name = var.app_name
  namespace = "default"
  secret_name = local.traefik_secret_cert_name
  service_name = local.tls_store
}

module "tlsoptions" {
  count = var.reverse_proxy && !var.k8s_crds ? 1 : 0
  source = "./modules/traefik/tlsoptions"
  app_name = var.app_name
  namespace = local.namespace
  service_name = local.tls_options
}

module "ingress-route" {
  count = var.reverse_proxy && !var.k8s_crds ? 1 : 0
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
  secret_name = local.traefik_secret_cert_name
  issuer_name = local.issuer_name
  # host_name = "169.46.98.220.nip.io"
  # host_name = "memories.mooo.com"
  host_name = ["trimino.xyz", "www.trimino.xyz"]
  service_name = local.ingress_route
}

###################################################################################################
# cert manager                                                                                    #
###################################################################################################
module "cert-manager" {
  count = var.reverse_proxy ? 1 : 0
  source = "./modules/traefik/cert-manager/cert-manager"
  namespace = local.namespace
  chart_version = "1.18.0"
  service_name = "fin-cert-manager"
}

module "acme-issuer" {
  count = var.reverse_proxy && !var.k8s_crds ? 1 : 0
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
  count = var.reverse_proxy && !var.k8s_crds ? 1 : 0
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
module "fin-finances-persistent" {
  count = var.persistent_disk && !var.k8s_crds ? 1 : 0
  # Specify the location of the module, which contains the file main.tf.
  source = "./modules/deployment"
  dir_path = ".."
  app_name = var.app_name
  app_version = var.app_version
  namespace = local.namespace
  image_tag = var.build_image ? "" : "${var.cr_username}/${local.svc_finances}:${var.app_version}"
  build_image = var.build_image
  cr_login_server = local.cr_login_server
  cr_username = var.cr_username
  cr_password = var.cr_password
  labels = {
    "app" = var.app_name
  }
  replicas = 1
  # Limits and requests for CPU resources are measured in millicores. If the container needs one
  # full core to run, use the value '1000m.' If the container only needs 1/4 of a core, use the
  # value of '250m.'
  resources = {  # QoS - Guaranteed
    limits_cpu = "300m"
    limits_memory = "300Mi"
  }
  # https://kubernetes.io/docs/concepts/configuration/secret/
  # If the order of Secrets changes, the Deployment must be changed accordingly. See
  # spec.image_pull_secrets.
  secrets = [{
    name = "${local.svc_finances}-registry-credentials"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
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
  # *** s3 storage ***
  # {
  #   name = "${local.svc_finances}-s3-storage"
  #   data = {
  #     obj_storage_ns = var.obj_storage_ns
  #     region = var.region
  #     aws_access_key_id = var.aws_access_key_id
  #     aws_secret_access_key = var.aws_secret_access_key
  #   }
  #   type = "Opaque"
  # }
  # *** s3 storage ***
  ]
  service_account = {
    name = "${local.svc_finances}-service-account"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    annotations = {}
    automount_service_account_token = false
    secrets = [{
      name = "${local.svc_finances}-registry-credentials"
    },
    # {
    #   name = "${local.svc_finances}-s3-storage"
    # }
    ]
  }
  role = {
    name = "${local.svc_finances}-role"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    annotations = {}
    rules = [{
      # It provides read-only access to information without allowing modification.
      api_groups = [""]
      resources = ["pods", "configmaps"]
      verbs = ["get", "watch", "list"]
    }]
  }
  role_binding = {
    name = "${local.svc_finances}-role-binding"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    annotations = {}
    role_ref = {
      kind = "Role"
      name = "${local.svc_finances}-role"
      api_group = "rbac.authorization.k8s.io"
    }
    subjects = [{
      kind = "ServiceAccount"
      name = "${local.svc_finances}-service-account"
      namespace = local.namespace
    }]
  }
  # Configure environment variables specific to the app.
  env = {
    PPROF = var.pprof
    K8S = true
    HTTP_PORT = "8080"
    SVC_NAME = local.svc_finances
    APP_NAME_VER = "${var.app_name} ${var.app_version}"
    MAX_RETRIES = 3
  }
  # *** s3 storage ***
  # env_secret = [{
  #   env_name = "AWS_SECRET_ACCESS_KEY"
  #   secret_name = "${local.svc_finances}-s3-storage"
  #   secret_key = "aws_secret_access_key"
  # },
  # {
  #   env_name = "OBJ_STORAGE_NS"
  #   secret_name = "${local.svc_finances}-s3-storage"
  #   secret_key = "obj_storage_ns"
  # },
  # {
  #   env_name = "AWS_REGION"
  #   secret_name = "${local.svc_finances}-s3-storage"
  #   secret_key = "region"
  # },
  # {
  #   env_name = "AWS_ACCESS_KEY_ID"
  #   secret_name = "${local.svc_finances}-s3-storage"
  #   secret_key = "aws_access_key_id"
  # }]
  # *** s3 storage ***
  # *** env_field ***
  # env_field = [{
  #   env_name = "POD_ID"
  #   field_path = "status.podIP"
  # }]
  # *** env_field ***
  volume_mount = [{
    name = "wsf"
    mount_path = "/wsf_data_dir"
    read_only = false
  }]
  volume_pv = [{
    name = "wsf"
    claim_name = "${var.app_name}-pvc"
  }]
  persistent_volume_claims = [{
    name = "${var.app_name}-pvc"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    volume_mode = "Filesystem"
    # The volume can be mounted as read-write by many nodes.
    access_modes = ["ReadWriteOnce"]
    # The minimum amount of persistent storage that a PVC can request is 50GB. If the request is
    # for less than 50GB, the request is rounded up to 50GB.
    storage_size = "50Gi"
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
  # You should always define a readiness probe, even if it's as simple as sending an HTTP request
  # to the base URL.
  readiness_probe = [{
    # Always remember to set an initial delay to account for your app's startup time.
    initial_delay_seconds = 2
    period_seconds = 25
    timeout_seconds = 1
    failure_threshold = 3
    success_threshold = 1
    http_get = [{
      path = "/readiness"
      port = 8080
      scheme = "HTTP"
    }]
    # tcp_socket = {
    #   port = 8088
    # }
    # exec = {
    #   command = [
    #     "/bin/sh",
    #     "-c",
    #     "ls -al /wsf_data_dir"
    # ]}
  }]
  # You should always define a liveness probe. Keep probes light.
  liveness_probe = [{
    initial_delay_seconds = 5
    period_seconds = 20
    timeout_seconds = 1
    # Don't bother implementing retry loops; K8s will retry the probe.
    failure_threshold = 1
    success_threshold = 1
    http_get = [{
      path = "/liveness"
      port = 8080
      scheme = "HTTP"
    }]
    # tcp_socket = {
    #   port = 8080
    # }
    # exec = {
    #   command = [
    #     "/bin/sh",
    #     "-c",
    #     "ls -al /wsf_data_dir"
    # ]}
  }]
  service = {
    name = local.svc_finances
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    ports = [{
      name = "ports"
      service_port = 80
      target_port = 8080
      protocol = "TCP"
    }]
    selector = {
      "svc_selector_label" = "svc-${local.svc_finances}"
    }
    type = "ClusterIP"
  }
  service_name = local.svc_finances
}

module "fin-finances-empty" {  # Using emptyDir.
  count = var.empty_dir && !var.k8s_crds ? 1 : 0
  # Specify the location of the module, which contains the file main.tf.
  source = "./modules/deployment"
  dir_path = ".."
  app_name = var.app_name
  app_version = var.app_version
  namespace = local.namespace
  image_tag = var.build_image == true ? "" : "${var.cr_username}/${local.svc_finances}:${var.app_version}"
  build_image = var.build_image
  cr_login_server = local.cr_login_server
  cr_username = var.cr_username
  cr_password = var.cr_password
  labels = {
    "app" = var.app_name
  }
  affinity = [{
    pod_anti_affinity = [{
      required_during_scheduling_ignored_during_execution = [{
        topology_key = "kubernetes.io/hostname"
        # match_labels = {
        #   "finances" = "running"
        # }
        match_expressions = [{
          "key" = "finances"
          "operator" = "In"
          "values" = ["running"]
        }]
      }]
    }]
  }]
  replicas = 3
  # See empty_dir.
  init_container = [{
    name = "file-permission"
    image = "busybox:1.34.1"
    image_pull_policy = "IfNotPresent"
    command = [
      "/bin/sh",
      "-c",
      "chown -v -R 1100:1100 /wsf_data_dir && chmod -R 750 /wsf_data_dir"
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
  # Limits and requests for CPU resources are measured in millicores. If the container needs one
  # full core to run, use the value '1000m.' If the container only needs 1/4 of a core, use the
  # value of '250m.'
  resources = {  # QoS - Guaranteed
    limits_cpu = "300m"
    limits_memory = "300Mi"
  }
  # https://kubernetes.io/docs/concepts/configuration/secret/
  # If the order of Secrets changes, the Deployment must be changed accordingly. See
  # spec.image_pull_secrets.
  secrets = [{
    name = "${local.svc_finances}-registry-credentials"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
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
  }]
  service_account = {
    name = "${local.svc_finances}-service-account"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    annotations = {}
    automount_service_account_token = false
    secrets = [{
      name = "${local.svc_finances}-registry-credentials"
    },
    # {
    #   name = "${local.svc_finances}-s3-storage"
    # }
    ]
  }
  role = {
    name = "${local.svc_finances}-role"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    annotations = {}
    rules = [{
      # It provides read-only access to information without allowing modification.
      api_groups = [""]
      resources = ["pods", "configmaps"]
      verbs = ["get", "watch", "list"]
    }]
  }
  role_binding = {
    name = "${local.svc_finances}-role-binding"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    annotations = {}
    role_ref = {
      kind = "Role"
      name = "${local.svc_finances}-role"
      api_group = "rbac.authorization.k8s.io"
    }
    subjects = [{
      kind = "ServiceAccount"
      name = "${local.svc_finances}-service-account"
      namespace = local.namespace
    }]
  }
  # Configure environment variables specific to the app.
  env = {
    PPROF = var.pprof
    K8S = true
    HTTP_PORT = "8080"
    SVC_NAME = local.svc_finances
    APP_NAME_VER = "${var.app_name} ${var.app_version}"
    MAX_RETRIES = 3
  }
  # *** env_field ***
  # env_field = [{
  #  env_name = "POD_ID"
  #  field_path = "status.podIP"
  # }]
  # *** env_field ***
  # When using the emptyDir{}, the init_container is required.
  volume_empty_dir = [{
    name = "wsf"
  }]
  volume_mount = [{
    name = "wsf"
    mount_path = "/wsf_data_dir"
    read_only = false
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
  # You should always define a readiness probe, even if it's as simple as sending an HTTP request
  # to the base URL.
  readiness_probe = [{
    # Always remember to set an initial delay to account for your app's startup time.
    initial_delay_seconds = 2
    period_seconds = 25
    timeout_seconds = 1
    failure_threshold = 3
    success_threshold = 1
    http_get = [{
      path = "/readiness"
      port = 8080
      scheme = "HTTP"
    }]
    # tcp_socket = {
    #   port = 8088
    # }
    # exec = {
    #   command = [
    #     "/bin/sh",
    #     "-c",
    #     "ls -al /wsf_data_dir"
    # ]}
  }]
  # You should always define a liveness probe. Keep probes light.
  liveness_probe = [{
    initial_delay_seconds = 5
    period_seconds = 20
    timeout_seconds = 1
    # Don't bother implementing retry loops; K8s will retry the probe.
    failure_threshold = 1
    success_threshold = 1
    http_get = [{
      path = "/liveness"
      port = 8080
      scheme = "HTTP"
    }]
    # tcp_socket = {
    #   port = 8080
    # }
    # exec = {
    #   command = [
    #     "/bin/sh",
    #     "-c",
    #     "ls -al /wsf_data_dir"
    # ]}
  }]
  service = {
    name = local.svc_finances
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    ports = [{
      name = "ports"
      service_port = 80
      target_port = 8080
      protocol = "TCP"
    }]
    selector = {
      "svc_selector_label" = "svc-${local.svc_finances}"
    }
    type = "ClusterIP"
  }
  service_name = local.svc_finances
}




# /*
module "fin-MySqlServer" {
  count = var.db_mysql && !var.k8s_crds ? 1 : 0
  # Specify the location of the module, which contains the file main.tf.
  source = "./modules/statefulset"
  app_name = var.app_name
  app_version = var.app_version
  namespace = local.namespace
  image_tag = var.mysql_image_tag
  publish_not_ready_addresses = true
  labels = {
    "app" = var.app_name
  }
  # config_map = [{
  #   name = "${var.app_name}-mysql-config-map"
  #   namespace = local.namespace
  #   labels = {
  #     "app" = var.app_name
  #   }
  #   data = {
  #     "MYSQL_DATABASE" = var.mysql_database
  #     "MYSQL_USER" = var.mysql_user
  #     "MYSQL_PASSWORD" = var.mysql_password
  #     "MYSQL_ROOT_PASSWORD" = var.mysql_root_password
  #   }
  # }]
  # Configure environment variables specific to the app.
  env = {
    MYSQL_DATABASE = var.mysql_database
  }
  # If the order of Secrets changes, the Deployment must be changed accordingly. See
  # spec.image_pull_secrets.
  secrets = [{
    name = "${local.svc_my_sql}-secret"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    # Plain-text data.
    data = {
      MYSQL_USER = var.mysql_user
      MYSQL_PASSWORD = var.mysql_password
      MYSQL_ROOT_PASSWORD = var.mysql_root_password
    }
    type = "Opaque"
    immutable = true
  }]
  # Always use 3 or more nodes for fault tolerance.
  replicas = 3
  resources = {  # QoS - Guaranteed
    limits_cpu = "500m"
    limits_memory = "1Gi"
  }
  volume_mount = [{
    name = "wsf"
    mount_path = "/wsf_data_dir/mysql"
    read_only = false
  }]
  volume_claim_templates = [{
    name = "wsf"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    volume_mode = "Filesystem"
    # The volume can be mounted as read-write by many nodes.
    access_modes = ["ReadWriteOnce"]
    # The minimum amount of persistent storage that a PVC can request is 50GB. If the request is
    # for less than 50GB, the request is rounded up to 50GB.
    storage = "50Gi"
    storage_class_name = "oci-bv"
  }]
  # readiness_probe = [{
  #   initial_delay_seconds = 5
  #   period_seconds = 2
  #   timeout_seconds = 1
  #   # failure_threshold = 3
  #   # success_threshold = 1
  #   # http_get = [{
  #   #   path = "/readiness"
  #   #   port = 8080
  #   #   scheme = "HTTP"
  #   # }]
  #   # tcp_socket = {
  #   #   port = 8088
  #   # }
  #   exec = {
  #     # command: ["mysql", "-h", "127.0.0.1", "-e", "SELECT 1"]
  #     command = [
  #       "mysqladmin",
  #       "ping",
  #       "-u", "root",
  #       "-h", "127.0.0.1",
  #       "-p${MYSQL_ROOT_PASSWORD}"
  #   ]}
  # }]
  # liveness_probe = [{
  #   initial_delay_seconds = 30
  #   period_seconds = 10
  #   timeout_seconds = 5
  #   # failure_threshold = 3
  #   # success_threshold = 1
  #   # http_get = [{
  #   #   path = "/readiness"
  #   #   port = 8080
  #   #   scheme = "HTTP"
  #   # }]
  #   # tcp_socket = {
  #   #   port = 8088
  #   # }
  #   exec = {
  #     # command: ["mysqladmin", "ping"]
  #     command = [
  #       "mysqladmin",
  #       "ping",
  #       "-u", "root",
  #       "-h", "127.0.0.1",
  #       "-p${MYSQL_ROOT_PASSWORD}"
  #   ]}
  # }]
  # pod_security_context = [{
  #   fs_group = 2200
  # }]
  security_context = [{
    run_as_non_root = true
    run_as_user = 1100
    run_as_group = 1100
    read_only_root_filesystem = false
  }]
  ports = [{
    name = "mysql"
    # Standard MySQL port – used by clients and MySQL Router to connect to MySQL server.
    service_port = 3306
    target_port = 3306
    protocol = "TCP"
  },
  {
    name = "x"
    # X Protocol – used by MySQL Shell, MySQL Router, and Group Replication for
    # administrative/configuration tasks
    service_port = 33060
    target_port = 33060
    protocol = "TCP"
  },
  {
    name = "group"
    # Group Replication port – used by MySQL nodes for internal communication and replication
    # coordination.
    service_port = 33061
    target_port = 33061
    protocol = "TCP"
  },
  {
    name = "split"
    # MySQL Router read/write split ports (configurable).
    service_port = 6606
    target_port = 6606
    protocol = "TCP"
  }]
  service_name = local.svc_my_sql
}

/***
https://dev.mysql.com/doc/mysql-router/8.0/en/
Notes:
(1) When used with a MySQL InnoDB Cluster, MySQL Router acts as a proxy to hide the multiple MySQL
    instances on your network and map the data requests to one of the cluster instances.
(2) The recommended deployment model for MySQL Router is with InnoDB Cluster, with Router sitting
    on the same host as the application.
(3) Running in a container requires a working InnoDB cluster. If supplied, the run script waits for
    the given mysql host to start, the InnoDB cluster to have the
    MYSQL_INNODB_CLUSTER_MEMBERS-defined number of members, and then uses the supplied host for
    bootstrapping.
***/
module "fin-MySqlRouter" {
  count = var.k8s_crds ? 0 : 0
  # Specify the location of the module, which contains the file main.tf.
  source = "./modules/deployment"
  dir_path = "rrrr"
  app_name = var.app_name
  app_version = var.app_version
  namespace = local.namespace
  image_tag = var.mysql_router_image_tag
  build_image = false
  labels = {
    "app" = var.app_name
  }
  env = {
    MYSQL_HOST = var.mysql_router_host
    MYSQL_PORT = var.mysql_router_port
    MYSQL_INNODB_CLUSTER_MEMBERS = var.mysql_router_cluster_members
    MYSQL_ROUTER_BOOTSTRAP_EXTRA_OPTIONS = var.mysql_router_bootstrap
  }
  secrets = [{
    name = "${local.svc_my_sql_router}-secret"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    # Plain-text data.
    data = {
      MYSQL_USER = var.mysql_router_user
      MYSQL_PASSWORD = var.mysql_router_password
    }
    type = "Opaque"
  }]



  replicas = 3
  resources = {
    limits_cpu = "500m"
    limits_memory = "0.5Gi"
  }
  # command = ["/bin/sh"]
  # args = ["-c",  # https://www.man7.org/linux/man-pages/man1/bash.1.html
  #   "while true; do sleep 3600; done"]
  service = {
    name = local.svc_my_sql_router
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    ports = [{  # https://hub.docker.com/r/mysql/mysql-router#exposed-ports
      name = "read-write"  # Primary.
      service_port = 6446
      target_port = 6446
      protocol = "TCP"
    },
    {
      name = "read-only"  # Secondary.
      service_port = 6447
      target_port = 6447
      protocol = "TCP"
    }]
    selector = {
      "svc_selector_label" = "svc-${local.svc_finances}"
    }
    type = "ClusterIP"
  }


  service_name = local.svc_my_sql_router
}
