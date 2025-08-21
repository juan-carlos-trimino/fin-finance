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
  deployment_finances = "fin-finances"
  service_name_finances = "fin-finances"
  svc_gateway = "fin-gateway"
  svc_error_page = "fin-error-page"
  svc_traefik = "fin-traefik"
  svc_mysql = "fin-mysql"
  svc_mysql_router = "fin-mysql-router"
  statefulset_postgres_master = "fin-postgres-master"
  service_name_postgres_master = "fin-postgres-master-headless"
  statefulset_postgres_replica = "fin-postgres-replica"
  service_name_postgres_replica = "fin-postgres-replica-headless"
  ############
  # Services #
  ############
  # DNS translates hostnames to IP addresses; the container name is the hostname. When using Docker
  # and Docker Compose, DNS works automatically.
  # In K8s, a service makes the deployment accessible by other containers via DNS.
  # FQDN: service-name.namespace.svc.cluster.local
  svc_dns_error_page = "${local.svc_error_page}.${local.namespace}${var.cluster_domain_suffix}"
  svc_dns_finances = "${local.deployment_finances}.${local.namespace}${var.cluster_domain_suffix}"
  svc_dns_mysql_server = "${local.svc_mysql}-headless.${local.namespace}${var.cluster_domain_suffix}"
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
  svc_finances = local.deployment_finances
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
  #
  app_name = var.app_name
  app_version = var.app_version
  build_image = var.build_image
  cr_login_server = local.cr_login_server
  cr_username = var.cr_username
  cr_password = var.cr_password
  # Configure environment variables specific to the app.
  env = {
    PPROF = var.pprof
    K8S = true
    HTTP_PORT = "8080"
    SVC_NAME = local.service_name_finances
    APP_NAME_VER = "${var.app_name} ${var.app_version}"
    MAX_RETRIES = 3
  }
  # *** env_field ***
  # env_field = [{
  #   env_name = "POD_ID"
  #   field_path = "status.podIP"
  # }]
  # *** env_field ***
  # *** s3 storage ***
  # env_secret = [{
  #   env_name = "AWS_SECRET_ACCESS_KEY"
  #   secret_name = "${local.deployment_finances}-s3-storage"
  #   secret_key = "aws_secret_access_key"
  # },
  # {
  #   env_name = "OBJ_STORAGE_NS"
  #   secret_name = "${local.deployment_finances}-s3-storage"
  #   secret_key = "obj_storage_ns"
  # },
  # {
  #   env_name = "AWS_REGION"
  #   secret_name = "${local.deployment_finances}-s3-storage"
  #   secret_key = "region"
  # },
  # {
  #   env_name = "AWS_ACCESS_KEY_ID"
  #   secret_name = "${local.deployment_finances}-s3-storage"
  #   secret_key = "aws_access_key_id"
  # }]
  # *** s3 storage ***
  image_tag = var.build_image ? "" : "${var.cr_username}/${local.deployment_finances}:${var.app_version}"
  labels = {
    "app" = var.app_name
  }
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
  namespace = local.namespace
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
  replicas = 1
  # Limits and requests for CPU resources are measured in millicores. If the container needs one
  # full core to run, use the value '1000m.' If the container only needs 1/4 of a core, use the
  # value of '250m.'
  resources = {  # QoS - Guaranteed
    limits_cpu = "300m"
    limits_memory = "300Mi"
  }
  role = {
    name = "${local.deployment_finances}-role"
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
    name = "${local.deployment_finances}-role-binding"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    annotations = {}
    role_ref = {
      kind = "Role"
      name = "${local.deployment_finances}-role"
      api_group = "rbac.authorization.k8s.io"
    }
    subjects = [{
      kind = "ServiceAccount"
      name = "${local.deployment_finances}-service-account"
      namespace = local.namespace
    }]
  }
  # https://kubernetes.io/docs/concepts/configuration/secret/
  # If the order of Secrets changes, the Deployment must be changed accordingly. See
  # spec.image_pull_secrets.
  secrets = [{
    name = "${local.deployment_finances}-registry-credentials"
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
  #   name = "${local.deployment_finances}-s3-storage"
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
  security_context = [{
    run_as_non_root = true
    run_as_user = 1100
    run_as_group = 1100
    read_only_root_filesystem = true
  }]
  service = {
    name = local.service_name_finances
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
      "svc_selector_label" = "svc-${local.service_name_finances}"
    }
    type = "ClusterIP"
  }
  service_account = {
    name = "${local.deployment_finances}-service-account"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    annotations = {}
    automount_service_account_token = false
    secrets = [{
      name = "${local.deployment_finances}-registry-credentials"
    },
    # {
    #   name = "${local.deployment_finances}-s3-storage"
    # }
    ]
  }
  deployment_name = local.deployment_finances
  volume_mount = [{
    name = "wsf"
    mount_path = "/wsf_data_dir"
    read_only = false
  }]
  volume_pv = [{
    name = "wsf"
    claim_name = "${var.app_name}-pvc"
  }]
}

module "fin-finances-empty" {  # Using emptyDir.
  count = var.empty_dir && !var.k8s_crds ? 1 : 0
  # Specify the location of the module, which contains the file main.tf.
  source = "./modules/deployment"
  dir_path = ".."
  #
  affinity = {
    pod_anti_affinity = {
      required_during_scheduling_ignored_during_execution = [{
        topology_key = "kubernetes.io/hostname"
        label_selector = {
          # Tell K8s to avoid scheduling a replica in a node where there is already a replica with
          # the label "aff-finances: running".
          match_expressions = [{
            "key" = "aff-finances"
            "operator" = "In"
            "values" = ["running"]
          }]
        }
      }]
    }
  }
  app_name = var.app_name
  app_version = var.app_version
  build_image = var.build_image
  cr_login_server = local.cr_login_server
  cr_password = var.cr_password
  cr_username = var.cr_username
  # Configure environment variables specific to the app.
  env = {
    PPROF = var.pprof
    K8S = true
    HTTP_PORT = "8080"
    SVC_NAME = local.service_name_finances
    APP_NAME_VER = "${var.app_name} ${var.app_version}"
    MAX_RETRIES = 3
  }
  image_tag = var.build_image == true ? "" : "${var.cr_username}/${local.deployment_finances}:${var.app_version}"
  # See empty_dir.
  init_container = [{
    name = "file-permission"
    #
    command = ["/bin/sh",
      "-c",
      "chown -v -R 1100:1100 /wsf_data_dir && chmod -R 750 /wsf_data_dir"
    ]
    image = "busybox:1.34.1"
    image_pull_policy = "IfNotPresent"
    security_context = [{
      run_as_non_root = false
      run_as_user = 0
      run_as_group = 0
      read_only_root_filesystem = true
      privileged = true
    }]
    volume_mounts = [{
      name = "wsf"
      mount_path = "/wsf_data_dir"
      read_only = false
    }]
  }]
  labels = {
    "aff-finances" = "running"
    "app" = var.app_name
  }
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
  namespace = local.namespace
  pod_security_context = [{
    fs_group = 2200
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
  replicas = 3
  # Limits and requests for CPU resources are measured in millicores. If the container needs one
  # full core to run, use the value '1000m.' If the container only needs 1/4 of a core, use the
  # value of '250m.'
  resources = {  # QoS - Guaranteed
    limits_cpu = "300m"
    limits_memory = "300Mi"
  }
  role = {
    name = "${local.deployment_finances}-role"
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
    name = "${local.deployment_finances}-role-binding"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    annotations = {}
    role_ref = {
      kind = "Role"
      name = "${local.deployment_finances}-role"
      api_group = "rbac.authorization.k8s.io"
    }
    subjects = [{
      kind = "ServiceAccount"
      name = "${local.deployment_finances}-service-account"
      namespace = local.namespace
    }]
  }
  # https://kubernetes.io/docs/concepts/configuration/secret/
  # If the order of Secrets changes, the Deployment must be changed accordingly. See
  # spec.image_pull_secrets.
  secrets = [{
    name = "${local.deployment_finances}-registry-credentials"
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
  security_context = [{
    run_as_non_root = true
    run_as_user = 1100
    run_as_group = 1100
    read_only_root_filesystem = true
  }]
  service = {
    name = local.service_name_finances
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
      "svc_selector_label" = "svc-${local.service_name_finances}"
    }
    type = "ClusterIP"
  }
  service_account = {
    name = "${local.deployment_finances}-service-account"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    annotations = {}
    automount_service_account_token = false
    secrets = [{
      name = "${local.deployment_finances}-registry-credentials"
    },
    # {
    #   name = "${local.deployment_finances}-s3-storage"
    # }
    ]
  }
  deployment_name = local.deployment_finances
  strategy = {
    type = "RollingUpdate"
    max_surge = 1
    max_unavailable = 0
  }
  # When using the emptyDir{}, the init_container is required.
  volume_empty_dir = [{
    name = "wsf"
  }]
  volume_mount = [{
    name = "wsf"
    mount_path = "/wsf_data_dir"
    read_only = false
  }]
}

module "fin-PostgresMaster" {
  count = var.db_postgres && !var.k8s_crds ? 1 : 0
  # Specify the location of the module, which contains the file main.tf.
  source = "./modules/statefulset"
  #
  app_name = var.app_name
  app_version = var.app_version
  env = {
    PGDATA = var.pgdata
  }
  env_field = [{
    name = "POP_IP"
    field_path = "status.podIP"
  }]
  config_map = [{
    # Same as volume_config_map.config_map_name.
    name = "${local.statefulset_postgres_master}-postgres-conf-files"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
      "db" = var.postgres_db_label
    }
    data = {
      # https://www.postgresql.org/docs/current/auth-pg-hba-conf.html
      "pg_hba.conf" = "${file("${var.path_postgres_configs}/pg_hba.conf")}"
      # https://www.postgresql.org/docs/current/auth-username-maps.html
      # "pg_ident.conf" = "${file("${var.path_postgres_configs}/pg_ident.conf")}"
      # Share by master and replica.
      "postgresql.conf" = "${file("${var.path_postgres_configs}/postgresql.conf")}"
      # Only for master.
      "master.conf" = "${file("${var.path_postgres_configs}/master.conf")}"
    }
  }, {
    # Same as volume.volume_config_map.config_map_name.name.
    name = "${local.statefulset_postgres_master}-postgres-script-files"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
      "db" = var.postgres_db_label
    }
    data = {
      "create-replication-user.sh" = "${file("${var.path_postgres_scripts}/create-replication-user.sh")}"
    }
  }]
  /***
  "-c": This is the first argument. It's typically used in conjunction with a shell command (like
  /bin/sh or /bin/bash) to indicate that the following string should be interpreted as a command
  string to be executed by the shell. The -c flag tells the shell to read commands from the string
  argument that follows.
  ***/
    # chown -R $(POSTGRES_USER):$(POSTGRES_USER) /var/lib/postgresql/data &&
    # chown -R 1999:1999 /wsf_data_dir/config &&
    # chown -R 1999:1999 /wsf_data_dir
    # chown -R 1999:1999 /var/lib/postgresql/data
  args = ["-c",
    <<-EOT
    /usr/local/bin/docker-entrypoint.sh postgres
    config_file=/wsf_data_dir/config/postgres/postgresql.conf
    EOT
  ]
  command = ["/bin/bash"]
  image_pull_policy = "IfNotPresent"
  image_tag = var.postgres_image_tag
  init_container = [{
    name = "file-permission"
    command = ["/bin/sh",
      "-c",
      "mkdir -p /wsf_data_dir/data/archive"
      # && chown -R 1999:1999 /wsf_data_dir/data/archive"
      # "mkdir -p /wsf_data_dir/data/archive && chown -R 1999:1999 /wsf_data_dir/data/archive && chown 1999:1999 /var/run/postgresql"
    ]
    image = "busybox:1.34.1"
    image_pull_policy = "IfNotPresent"
    security_context = [{
      run_as_non_root = true
      run_as_user = 1999
      run_as_group = 1999
      read_only_root_filesystem = true
      privileged = true
    }]
    volume_mounts = [{
      name = "wsf"
      mount_path = "/wsf_data_dir"
      read_only = false
    }]
  }]
  labels = {
    # "aff-mysql-server" = "running"
    "app" = var.app_name
    "db" = var.postgres_db_label
  }
  liveness_probe = [{
    initial_delay_seconds = 60  # Delay before the first probe.
    period_seconds = 10  # How often to perform the probe.
    timeout_seconds = 3  # Timeout for the probe command.
    failure_threshold = 3  # Number of consecutive failures before marking unready.
    success_threshold = 1
    exec = {
      command = ["pg_isready",
        # "--host", "$POD_IP",
        # "--port", "5432",
        "--username", "${var.postgres_user}",
        "--dbname", "${var.postgres_db}"
      ]
    }
  }]
  namespace = local.namespace
  # Ensure that the non-root user running the container has the necessary group permissions to
  # access files in mounted volumes.
  pod_security_context = [{
    fs_group = 1999
  }]
  readiness_probe = [{
    initial_delay_seconds = 30
    period_seconds = 5
    timeout_seconds = 3
    failure_threshold = 3
    # success_threshold = 1
    exec = {
      # https://www.postgresql.org/docs/current/app-pg-isready.html
      command = ["pg_isready",
        "-U", "${var.postgres_user}",
        "-d", "${var.postgres_db}"
      ]
    }
  }]
  # Always use 3 or more nodes for fault tolerance.
  replicas = 1
  resources = {  # QoS - Guaranteed
    limits_cpu = "250m"
    limits_memory = "1Gi"
  }
  secrets = [{
    name = "${local.statefulset_postgres_master}-secret"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
      "db" = var.postgres_db_label
    }
    # Plain-text data.
    data = {
      POSTGRES_DB = var.postgres_db
      POSTGRES_USER = var.postgres_user
      POSTGRES_PASSWORD = var.postgres_password
      REPLICATION_PASSWORD = var.replication_password
    }
    type = "Opaque"
    immutable = true
  }]
  service = {
    name = local.service_name_postgres_master
    namespace = local.namespace
    labels = {
      "app" = var.app_name
      "db" = var.postgres_db_label
    }
    ports = [{
      name = "postgres"
      service_port = 5432
      target_port = 5432
      protocol = "TCP"
    }]
    selector = {
      "svc_selector_label" = "svc-${local.service_name_postgres_master}"
    }
    publish_not_ready_addresses = true
    type = "ClusterIP"
  }
  security_context = [{
    run_as_non_root = true
    run_as_user = 1999
    run_as_group = 1999
    read_only_root_filesystem = false
  }]
  statefulset_name = local.statefulset_postgres_master
  update_strategy = {
    type = "RollingUpdate"
  }
  volume_claim_templates = [{
    name = "wsf"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
      "db" = var.postgres_db_label
    }
    volume_mode = "Filesystem"
    # https://kubernetes.io/docs/concepts/storage/persistent-volumes/?ref=kodekloud.com#access-modes
    access_modes = ["ReadWriteOnce"]
    # The minimum amount of persistent storage that a PVC can request is 50GB. If the request is
    # for less than 50GB, the request is rounded up to 50GB.
    storage = "50Gi"
    storage_class_name = "oci-bv"
  }]
  volume_config_map = [{
    name = "config"
    config_map_name = "${local.statefulset_postgres_master}-postgres-conf-files"
    default_mode = "0660"
    items = [{
      key = "pg_hba.conf"
      path = "pg_hba.conf"
    }, /*{
      key = "pg_ident.conf"
      path = "pg_ident.conf"
    },*/ {
      key = "postgresql.conf"
      path = "postgresql.conf"
    }, {
      key = "master.conf"
      path = "master.conf"
    }]
  }, {
    name = "script"
    config_map_name = "${local.statefulset_postgres_master}-postgres-script-files"
    default_mode = "0760"
  }]
  volume_mount = [{
    name = "wsf"
    mount_path = "/wsf_data_dir"
    read_only = false
  }, {
    name = "config"
    mount_path = "/wsf_data_dir/config"
    read_only = false
  }, {
    name = "script"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
      "db" = var.postgres_db_label
    }
    # https://hub.docker.com/_/postgres#initialization-scripts
    mount_path = "/docker-entrypoint-initdb.d/create-replication-user.sh"
    sub_path = "create-replication-user.sh"
  }]
}

module "fin-PostgresReplica" {
  count = var.db_postgres && !var.k8s_crds ? 1 : 0
  depends_on = [
    module.fin-PostgresMaster
  ]
  # Specify the location of the module, which contains the file main.tf.
  source = "./modules/statefulset"
  #
  app_name = var.app_name
  app_version = var.app_version
  env = {
    PGDATA = var.pgdata
  }
  env_field = [{
    name = "POP_IP"
    field_path = "status.podIP"
  }]
  config_map = [{
    # Same as volume_config_map.config_map_name.
    name = "${local.statefulset_postgres_replica}-postgres-conf-files"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    data = {
      # https://www.postgresql.org/docs/current/auth-pg-hba-conf.html
      # "pg_hba.conf" = "${file("${var.path_postgres_configs}/pg_hba.conf")}"
      # https://www.postgresql.org/docs/current/auth-username-maps.html
      # "pg_ident.conf" = "${file("${var.path_postgres_configs}/pg_ident.conf")}"
      # Share by master and slave.
      "postgresql.conf" = "${file("${var.path_postgres_configs}/postgresql.conf")}"
      # Only for master.
      "replica.conf" = "${file("${var.path_postgres_configs}/replica.conf")}"
    }
  }/*, {
    # Same as volume_config_map.config_map_name.
    name = "${var.app_name}-postgres-script-files"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    data = {
      "create-replication-user.sh" = "${file("${var.path_postgres_scripts}/create-replication-user.sh")}"
    }
  }*/]
  /***
  "-c": This is the first argument. It's typically used in conjunction with a shell command (like
  /bin/sh or /bin/bash) to indicate that the following string should be interpreted as a command
  string to be executed by the shell. The -c flag tells the shell to read commands from the string
  argument that follows.
  ***/
    # chown -R $(POSTGRES_USER):$(POSTGRES_USER) /var/lib/postgresql/data &&
    # chown -R 1999:1999 /wsf_data_dir/config &&
    # chown -R 1999:1999 /wsf_data_dir
    # chown -R 1999:1999 /var/lib/postgresql/data
  args = ["-c",
    <<-EOT
    /usr/local/bin/docker-entrypoint.sh postgres
    config_file=/wsf_data_dir/config/postgres/postgresql.conf
    EOT
  ]
  command = ["/bin/bash"]
  image_pull_policy = "IfNotPresent"
  image_tag = var.postgres_image_tag
  init_container = [{
    name = "setup-replica-data-directory"
    env = {
      PGDATA = var.pgdata
      PGHOST = local.service_name_postgres_master
      PGPASSWORD = "rpassword"
    }


 /***
        The PGPASSWORD environment variable in PostgreSQL allows the specification of a password for database connections without requiring interactive input. This variable can be set in the shell before executing PostgreSQL client applications like psql or pg_dump.

For improved security, consider using the following alternatives:
.pgpass file:
.
A password file (~/.pgpass on Linux/macOS, %APPDATA%\postgresql\pgpass.conf on Windows) can store connection details, including passwords, in a more secure manner. The file permissions must be set correctly (e.g., chmod 600 ~/.pgpass) to prevent unauthorized access.

        ***/
    env_from_secrets = [
      "replica-data-directory-secret"
    ]


    command = ["/bin/bash",
      "-c",
      # https://www.postgresql.org/docs/current/app-pgbasebackup.html
      <<-EOT
      mkdir -p $PGDATA
      if [ -z "$(ls -A $PGDATA)" ];
      then
        echo "Running pg_basebackup to catch up replication server...";
        pg_basebackup -h $PGHOST -R -D $PGDATA -P -U replication;
        chown -R 1999:1999 $PGDATA;
      else
        echo "Skipping pg_basebackup because directory is not empty";
      fi
      EOT
    ]
    image = var.postgres_image_tag
    image_pull_policy = "IfNotPresent"
    security_context = [{
      run_as_non_root = true
      run_as_user = 1999
      run_as_group = 1999
      read_only_root_filesystem = false
      privileged = true
    }]
    volume_mounts = [{
      name = "wsf"
      mount_path = "/wsf_data_dir"
      read_only = false
    }]
  }]
  labels = {
    # "aff-mysql-server" = "running"
    "app" = var.app_name
    "db" = var.postgres_db_label
  }
  liveness_probe = [{
    initial_delay_seconds = 60  # Delay before the first probe.
    period_seconds = 10  # How often to perform the probe.
    timeout_seconds = 3  # Timeout for the probe command.
    failure_threshold = 3  # Number of consecutive failures before marking unready.
    success_threshold = 1
    exec = {
      command = ["pg_isready",
        "--host", "$POD_IP",
        "--port", "5432",
        "--username", "${var.postgres_user}",
        "--dbname", "${var.postgres_db}"
      ]
    }
  }]
  namespace = local.namespace
  # Ensure that the non-root user running the container has the necessary group permissions to
  # access files in mounted volumes.
  pod_security_context = [{
    fs_group = 1999
  }]
  readiness_probe = [{
    initial_delay_seconds = 30
    period_seconds = 5
    timeout_seconds = 3
    failure_threshold = 3
    # success_threshold = 1
    exec = {
      # https://www.postgresql.org/docs/current/app-pg-isready.html
      command = ["pg_isready",
        "--host", "$POD_IP",
        "--port", "5432",
        "-U", "${var.postgres_user}",
        "-d", "${var.postgres_db}"
      ]
    }
  }]
  # Always use 3 or more nodes for fault tolerance.
  replicas = 1
  resources = {  # QoS - Guaranteed
    limits_cpu = "100m"
    limits_memory = "256Mi"
  }
  secrets = [{
    name = "${local.statefulset_postgres_replica}-secret"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    # Plain-text data.
    data = {
      POSTGRES_DB = var.postgres_db
      POSTGRES_USER = var.postgres_user
      POSTGRES_PASSWORD = var.postgres_password
      REPLICATION_PASSWORD = var.replication_password
    }
    type = "Opaque"
    immutable = true
  }, {
    name = "replica-data-directory-secret"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    # Plain-text data.
    data = {
      PGPASSWORD = var.replication_password
    }
    type = "Opaque"
    immutable = true
  }]
  service = {
    name = local.service_name_postgres_replica
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    ports = [{
      name = "postgres"
      service_port = 5432
      target_port = 5432
      protocol = "TCP"
    }]
    selector = {
      "svc_selector_label" = "svc-${local.service_name_postgres_replica}"
    }
    publish_not_ready_addresses = true
    type = "ClusterIP"
  }
  security_context = [{
    run_as_non_root = true
    run_as_user = 1999
    run_as_group = 1999
    read_only_root_filesystem = false
  }]
  statefulset_name = local.statefulset_postgres_replica
  update_strategy = {
    type = "RollingUpdate"
  }
  volume_claim_templates = [{
    name = "wsf"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
    volume_mode = "Filesystem"
    # https://kubernetes.io/docs/concepts/storage/persistent-volumes/?ref=kodekloud.com#access-modes
    access_modes = ["ReadWriteOnce"]
    # The minimum amount of persistent storage that a PVC can request is 50GB. If the request is
    # for less than 50GB, the request is rounded up to 50GB.
    storage = "50Gi"
    storage_class_name = "oci-bv"
  }]
  volume_config_map = [{
    name = "config"
    config_map_name = "${local.statefulset_postgres_replica}-postgres-conf-files"
    default_mode = "0660"
    items = [/*{
      key = "pg_hba.conf"
      path = "pg_hba.conf"
    }, {
      key = "pg_ident.conf"
      path = "pg_ident.conf"
    },*/ {
      key = "postgresql.conf"
      path = "postgresql.conf"
    }, {
      key = "replica.conf"
      path = "replica.conf"
    }]
  }]
  volume_mount = [{
    name = "wsf"
    mount_path = "/wsf_data_dir"
    read_only = false
  }, {
    name = "config"
    mount_path = "/wsf_data_dir/config"
    read_only = false
  }/*, {
    name = "init-scripts"
    # https://hub.docker.com/_/postgres#initialization-scripts
    mount_path = "/docker-entrypoint-initdb.d/create-replication-user.sh"
    sub_path = "create-replication-user.sh"
    read_only = false
  }*/]
}



module "fin-MySqlServer" {
  count = var.db_postgres && !var.k8s_crds ? /*1*/0 : 0
  # Specify the location of the module, which contains the file main.tf.
  source = "./modules/statefulset"
  #
  affinity = {
    pod_anti_affinity = {
      required_during_scheduling_ignored_during_execution = [{
        topology_key = "kubernetes.io/hostname"
        label_selector = {
          match_expressions = [{
            "key" = "aff-mysql-server"
            "operator" = "In"
            "values" = ["running"]
          }]
        }
      }]
    }
  }
  app_name = var.app_name
  app_version = var.app_version
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
  /***
  "-c": This is the first argument. It's typically used in conjunction with a shell command (like
  /bin/sh or /bin/bash) to indicate that the following string should be interpreted as a command
  string to be executed by the shell. The -c flag tells the shell to read commands from the string
  argument that follows.
  ***/
  args = ["-c",
    /***
    --server-id: https://dev.mysql.com/doc/refman/8.0/en/replication-options.html#sysvar_server_id
    --log-bin: https://dev.mysql.com/doc/refman/8.0/en/replication-options-binary-log.html#sysvar_log_bin
    --binlog-checksum: https://dev.mysql.com/doc/refman/8.0/en/replication-options-binary-log.html#option_mysqld_binlog-checksum
    --max-relay-log-size: https://dev.mysql.com/doc/refman/8.0/en/replication-options-replica.html#option_mysqld_max-relay-log-size
    --max-binlog-size: https://dev.mysql.com/doc/refman/8.0/en/replication-options-binary-log.html#sysvar_max_binlog_size
    --gtid-mode: https://dev.mysql.com/doc/refman/8.0/en/replication-options-gtids.html#sysvar_gtid_mode
    --enforce-gtid-consistency: https://dev.mysql.com/doc/refman/8.0/en/replication-options-gtids.html#sysvar_enforce_gtid_consistency
    --datadir: https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html#sysvar_datadir
    --verbose: https://dev.mysql.com/doc/refman/8.0/en/server-options.html#option_mysqld_verbose
    --default-authentication-plugin: https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html#sysvar_default_authentication_plugin
    --report-host: https://dev.mysql.com/doc/refman/8.0/en/replication-options-replica.html#sysvar_report_host
    ***/
    /***
    In Terraform, both << and <<- are used to define heredoc strings, which are multi-line string
    literals. The key difference lies in how they handle indentation:
    <<DELIMITER (Standard Heredoc): This marker defines a standard heredoc. The content within the
    heredoc, including any leading whitespace or indentation, is preserved exactly as written. This
    means if you indent the content within your heredoc for readability in your HCL code, that
    indentation will be part of the resulting string value.
    <<-DELIMITER (Indented Heredoc): This marker defines an indented heredoc. Terraform
    automatically removes the common leading whitespace from all lines within the heredoc. This is
    particularly useful when you want to indent your heredoc content to align with the surrounding
    HCL code for readability, but you do not want that indentation to be part of the actual string
    value. Terraform calculates the smallest amount of leading whitespace present on any line
    within the heredoc and removes that amount from all lines.
    Terraform's heredoc syntax can be combined with YAML to define multi-line YAML content directly
    within a Terraform configuration. The >- syntax in YAML indicates a "folded block" scalar,
    which removes newlines within the content, while preserving a single space between lines, and
    also trims a single trailing newline if present.
    >-: This is a YAML multiline string indicator. The ">" indicates a folded style, where newlines
    are folded into spaces, and the "-" removes a trailing newline. This syntax is used to define a
    multi-line string that will be passed as a single argument to the command. This is commonly
    used when you want to execute a complex script or multiple commands within a single argument to
    the shell. See https://yaml.org/spec/1.2.2/#8112-block-chomping-indicator.
    ***/
    # >-
    <<-EOT
    /usr/local/bin/docker-entrypoint.sh mysqld
    --server-id=$((40 +  $(echo $HOSTNAME | grep -o '[^-]*$') + 1))
    --log-bin=OFF
    --binlog-checksum=NONE
    --max-relay-log-size=0
    --max-binlog-size=524288
    --enforce-gtid-consistency=ON
    --gtid-mode=ON
    --datadir=$HOMEwsf_data_dir/mysql
    --verbose
    --default-authentication-plugin=caching_sha2_password
    --report-host=$HOSTNAME.${local.svc_dns_mysql_server}
    EOT
  ]
  command = ["/bin/bash"]
  # Configure environment variables specific to the app.
  env = {
    MYSQL_DATABASE = var.mysql_database
    MYSQL_ROOT_HOST = var.mysql_root_host
  }
  image_tag = var.mysql_image_tag
  labels = {
    "aff-mysql-server" = "running"
    "app" = var.app_name
  }
  /*
  liveness_probe = [{
    initial_delay_seconds = 150
    period_seconds = 30
    timeout_seconds = 30
    # failure_threshold = 60
    exec = {
      command = [
        "bash",
        "-c",
        <<-EOT
        |
        mysqladmin -uroot -p$MYSQL_ROOT_PASSWORD ping
        EOT
    ]}
  }]
  */
  namespace = local.namespace
  pod_security_context = [{
    fs_group = 2200
  }]
  /*
  readiness_probe = [{
    initial_delay_seconds = 150
    period_seconds = 30
    timeout_seconds = 30
    # failure_threshold = 60
    # success_threshold = 1
    exec = {
      command = [
        "bash",
        "-c",
        <<-EOT
        |
        mysql -h127.0.0.1 -uroot -p$MYSQL_ROOT_PASSWORD -e'SELECT 1'
        EOT
    ]}
  }]
  */
  # Always use 3 or more nodes for fault tolerance.
  replicas = 3
  resources = {  # QoS - Guaranteed
    limits_cpu = "500m"
    limits_memory = "1Gi"
  }
  # If the order of Secrets changes, the Deployment must be changed accordingly. See
  # spec.image_pull_secrets.
  secrets = [{
    name = "${local.svc_mysql}-secret"
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
  security_context = [{
    run_as_non_root = true
    run_as_user = 1100
    run_as_group = 1100
    read_only_root_filesystem = false
  }]
  service = {
    name = "${local.svc_mysql}-headless"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
    }
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
    selector = {
      "svc_selector_label" = "svc-${local.svc_mysql}-headless"
    }
    publish_not_ready_addresses = true
    type = "ClusterIP"
  }
  statefulset_name = local.svc_mysql
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
  volume_mount = [{
    name = "wsf"
    mount_path = "/wsf_data_dir/mysql"
    read_only = false
  }, {
    name = "mysql"
    mount_path = "/var/lib/mysql"
    sub_path = "mysql"
    read_only = false
  }]
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
  # count = var.db_mysql && !var.k8s_crds ? 1 : 0
  # Specify the location of the module, which contains the file main.tf.
  source = "./modules/deployment"
  dir_path = ""
  #
  affinity = {
    pod_affinity = {
      required_during_scheduling_ignored_during_execution = [{
        topology_key = "kubernetes.io/hostname"
        label_selector = {
          match_expressions = [{
            "key" = "aff-finances"
            "operator" = "In"
            "values" = ["running"]
          }]
        }
      }]
    }
    pod_anti_affinity = {
      required_during_scheduling_ignored_during_execution = [{
        topology_key = "kubernetes.io/hostname"
        label_selector = {
          match_expressions = [{
            "key" = "aff-mysql-router"
            "operator" = "In"
            "values" = ["running"]
          }]
        }
      }]
    }
  }
  app_name = var.app_name
  app_version = var.app_version
  build_image = false
  env = {
    MYSQL_HOST = local.svc_dns_mysql_server
    MYSQL_PORT = var.mysql_router_port
    MYSQL_INNODB_CLUSTER_MEMBERS = var.mysql_router_cluster_members
    MYSQL_ROUTER_BOOTSTRAP_EXTRA_OPTIONS = var.mysql_router_bootstrap
  }
  image_tag = var.mysql_router_image_tag
  labels = {
    "aff-mysql-router" = "running"
    "app" = var.app_name
  }
  liveness_probe = [{
    initial_delay_seconds = 150
    period_seconds = 30
    timeout_seconds = 30
    failure_threshold = 60
    tcp_socket = {
      port = 6446
    }
  }]
  namespace = local.namespace
  readiness_probe = [{
    initial_delay_seconds = 150
    period_seconds = 30
    timeout_seconds = 30
    failure_threshold = 60
    # success_threshold = 1
    tcp_socket = {
      port = 6446
    }
  }]
  replicas = 3
  resources = {
    limits_cpu = "500m"
    limits_memory = "0.5Gi"
  }
  secrets = [{
    name = "${local.svc_mysql_router}-secret"
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
  security_context = [{
    run_as_non_root = true
    run_as_user = 1100
    run_as_group = 1100
    read_only_root_filesystem = false
  }]
  service = {
    name = "${local.svc_mysql_router}"
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
    },
    {
      name = "x-read-write"  # X Protocol.
      service_port = 6448
      target_port = 6448
      protocol = "TCP"
    },
    {
      name = "x-read-only"  # X Protocol.
      service_port = 6449
      target_port = 6449
      protocol = "TCP"
    }]
    selector = {
      "svc_selector_label" = "svc-${local.svc_mysql_router}"
    }
    type = "ClusterIP"
  }
  deployment_name = local.svc_mysql_router
}
