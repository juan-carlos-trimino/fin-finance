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
  containers_security_context = {
    allow_privilege_escalation = false
    privileged = false
    read_only_root_filesystem = true
    run_as_non_root = true
  }
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
  # Ensure that the non-root user running the container has the necessary group permissions to
  # access files in mounted volumes.
  pod_security_context = {
    fs_group = 1100
    run_as_non_root = true
    run_as_user = 1100
    run_as_group = 1100
  }
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
  containers_security_context = {
    allow_privilege_escalation = false
    privileged = false
    read_only_root_filesystem = true
  }
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
    security_context = {
      run_as_non_root = false
      run_as_user = 0
      run_as_group = 0
      read_only_root_filesystem = true
      privileged = true
    }
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
  # Ensure that the non-root user running the container has the necessary group permissions to
  # access files in mounted volumes.
  pod_security_context = {
    fs_group = 1100
    run_as_non_root = true
    run_as_user = 1100
    run_as_group = 1100
  }
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
    }]
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
  /***
  "-c": This is the first argument. It's typically used in conjunction with a shell command (like
  /bin/sh or /bin/bash) to indicate that the following string should be interpreted as a command
  string to be executed by the shell. The -c flag tells the shell to read commands from the string
  argument that follows.
  ***/
  args = ["-c",
    "config_file=/postgres/config/postgresql.conf"
  ]
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
      "pg_hba.conf" = "${file("${var.postgres_config_path}/pg_hba.conf")}"
      # https://www.postgresql.org/docs/current/auth-username-maps.html
      "pg_ident.conf" = "${file("${var.postgres_config_path}/pg_ident.conf")}"
      "postgresql.conf" = "${file("${var.postgres_config_path}/postgresql.conf")}"
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
      "create-replication-user.sh" = "${file("${var.postgres_script_path}/create-replication-user.sh")}"
    }
  }]
  containers_security_context = {
    allow_privilege_escalation = false
    privileged = false
    read_only_root_filesystem = true
  }
  env = {
    PGDATA = var.postgres_data
  }
  env_field = [{
    name = "POD_IP"
    field_path = "status.podIP"
  }]
  image_pull_policy = "IfNotPresent"
  image_tag = var.postgres_image_tag
  init_container = [{
    name = "init-master"
    args = ["-c",
      <<-EOT
      # Check if variable is set and not empty.
      if [[ ! -n "$STANDBY_MODE" ]];
      then
        printf "The environment variable STANDBY_MODE is unset or empty.\n"
        printf "Set to \"on\" for primary server.\n"
        printf "Set to \"off\" for secondary (backup) server.\n"
        exit 1
      fi
      # The variable holding the path to the directory is enclosed in double quotes to handle
      # potential spaces or special characters in the path.
      if [ ! -d "$PGDATA" ];
      then
        printf "Creating data directory...\n"
        mkdir -p "$PGDATA"
        printf "Creating archive directory...\n"
        mkdir -p /wsf_data_dir/data/archive
      fi
      #
      if [ "$STANDBY_MODE" == "on" ];
      then
        # Initialize from backup if data directory is empty.
        if [ -z "$(ls -A "$PGDATA")" ];
        then
          printf "Initializing from backup...\n"
          # The PGPASSWORD environment variable in PostgreSQL allows the specification of a password
          # for database connections without requiring interactive input. This variable can be set in
          # the shell before executing PostgreSQL client applications like psql or pg_dump.
          export PGPASSWORD=$REPLICATION_PASSWORD
          # https://www.postgresql.org/docs/current/app-pgbasebackup.html
          pg_basebackup -v -h $PGHOST -p 5432 -D $PGDATA -U replication -R -Xs -Fp
        fi
      else
        # cat /postgres/initdb/create-replication-user.sh
        # https://hub.docker.com/_/postgres#initialization-scripts
        if [ ! -d "/docker-entrypoint-initdb.d" ];
        then
          mkdir /docker-entrypoint-initdb.d
        fi
        cp /postgres/initdb/* /docker-entrypoint-initdb.d/
        printf "Adding credential...\n"
        sed -i 's/#POSTGRES_USER/$(POSTGRES_USER)/g' /docker-entrypoint-initdb.d/create-replication-user.sh
        sed -i 's/#POSTGRES_DB/$(POSTGRES_DB)/g' /docker-entrypoint-initdb.d/create-replication-user.sh
        sed -i 's/#REPLICATION_PASSWORD/$(REPLICATION_PASSWORD)/g' /docker-entrypoint-initdb.d/create-replication-user.sh
        printf "Changing permissions for emptyDir...\n"
        chown -v -R 1999:1999 /docker-entrypoint-initdb.d && chmod -R 750 /docker-entrypoint-initdb.d
        # cat /docker-entrypoint-initdb.d/create-replication-user.sh
      fi
      printf "Changing permissions for /wsf_data_dir...\n"
      chown -v -R 1999:1999 /wsf_data_dir && chmod -R 760 /wsf_data_dir
      EOT
    ]
    command = ["/bin/sh"]
    env = {
      PGDATA = var.postgres_data
      PGHOST = ""
      # on - Secondary.
      # off - Primary.
      STANDBY_MODE = "off"
    }
    env_from_secrets = [
      "${local.statefulset_postgres_master}-secret"
    ]
    image = var.postgres_image_tag
    image_pull_policy = "IfNotPresent"
    security_context = {
      run_as_user = 0
      run_as_group = 0
      run_as_non_root = false
      read_only_root_filesystem = false
      privileged = false
    }
    volume_mounts = [{
      name = "wsf"
      mount_path = "/wsf_data_dir"
      read_only = false
    }, {
      name = "readonly-initdb-volume"
      # https://hub.docker.com/_/postgres#initialization-scripts
      mount_path = "/postgres/initdb"
      read_only = false
    }, {
      name = "writable-initdb-volume"
      mount_path = "/docker-entrypoint-initdb.d"
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
      command = ["/bin/sh",
        "-c",
        "pg_isready --host=$(POD_IP) --port=5432 --username=${var.postgres_user} --dbname=${var.postgres_db}"
      ]
    }
  }]
  namespace = local.namespace
  # Ensure that the non-root user running the container has the necessary group permissions to
  # access files in mounted volumes.
  pod_security_context = {
    fs_group = 1999
    run_as_non_root = true
    run_as_user = 1999
    run_as_group = 1999
  }
  readiness_probe = [{
    initial_delay_seconds = 30
    period_seconds = 5
    timeout_seconds = 3
    failure_threshold = 3
    # success_threshold = 1
    exec = {
      # https://www.postgresql.org/docs/current/app-pg-isready.html
      command = ["/bin/sh",
        "-c",
        "pg_isready --host=$(POD_IP) --port=5432 --username=${var.postgres_user} --dbname=${var.postgres_db}"
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
      REPLICATION_PASSWORD = var.postgres_replication_password
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
    name = "config-volume"
    config_map_name = "${local.statefulset_postgres_master}-postgres-conf-files"
    default_mode = "0760"
    items = [{
      key = "pg_hba.conf"
      path = "pg_hba.conf"
    }, {
      key = "pg_ident.conf"
      path = "pg_ident.conf"
    }, {
      key = "postgresql.conf"
      path = "postgresql.conf"
    }]
  }, {
    name = "readonly-initdb-volume"
    config_map_name = "${local.statefulset_postgres_master}-postgres-script-files"
    default_mode = "0750"
  }]
  volume_mount = [{
    name = "wsf"
    mount_path = "/wsf_data_dir"
    read_only = false
  }, {
    name = "wsf"
    /***
    For security reasons, you want to prevent processes running in a container from writing to the
    container's filesystem. If you make the container's filesystem read-only, you will need to
    mount a volume in every directory the app writes information; e.g., logs.
    ***/
    mount_path = "/var/run/postgresql"  # The path where Postgres stores its data.
    read_only = false
  }, {
    name = "config-volume"
    mount_path = "/postgres/config"
    read_only = false
  }, {
    name = "writable-initdb-volume"
    mount_path = "docker-entrypoint-initdb.d"
    read_only = false
  }]
  volume_empty_dir = [{
    name = "writable-initdb-volume"
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
  args = ["-c",
    "config_file=/postgres/config/postgresql.conf"
  ]
  config_map = [{
    # Same as volume_config_map.config_map_name.
    name = "${local.statefulset_postgres_replica}-postgres-conf-files"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
      "db" = var.postgres_db_label
    }
    data = {
      # https://www.postgresql.org/docs/current/auth-pg-hba-conf.html
      "pg_hba.conf" = "${file("${var.postgres_config_path}/pg_hba.conf")}"
      # https://www.postgresql.org/docs/current/auth-username-maps.html
      "pg_ident.conf" = "${file("${var.postgres_config_path}/pg_ident.conf")}"
      "postgresql.conf" = "${file("${var.postgres_config_path}/postgresql.conf")}"
    }
  }, {
    # Same as volume.volume_config_map.config_map_name.name.
    name = "${local.statefulset_postgres_replica}-postgres-script-files"
    namespace = local.namespace
    labels = {
      "app" = var.app_name
      "db" = var.postgres_db_label
    }
    data = {
      "create-replication-user.sh" = "${file("${var.postgres_script_path}/create-replication-user.sh")}"
    }
  }]
  containers_security_context = {
    allow_privilege_escalation = false
    privileged = false
    read_only_root_filesystem = true
  }
  env = {
    PGDATA = var.postgres_data
  }
  env_field = [{
    name = "POD_IP"
    field_path = "status.podIP"
  }]
  image_pull_policy = "IfNotPresent"
  image_tag = var.postgres_image_tag
  init_container = [{
    name = "init-replica"
    args = ["-c",
      <<-EOT
      # Check if variable is set and not empty.
      if [[ ! -n "$STANDBY_MODE" ]];
      then
        printf "The environment variable STANDBY_MODE is unset or empty.\n"
        printf "Set to \"on\" for primary server.\n"
        printf "Set to \"off\" for secondary (backup) server.\n"
        exit 1
      fi
      # The variable holding the path to the directory is enclosed in double quotes to handle
      # potential spaces or special characters in the path.
      if [ ! -d "$PGDATA" ];
      then
        printf "Creating data directory...\n"
        mkdir -p "$PGDATA"
        printf "Creating archive directory...\n"
        mkdir -p /wsf_data_dir/data/archive
      fi
      #
      if [ "$STANDBY_MODE" == "on" ];
      then
        # Initialize from backup if data directory is empty.
        if [ -z "$(ls -A "$PGDATA")" ];
        then
          printf "Initializing from backup...\n"
          # The PGPASSWORD environment variable in PostgreSQL allows the specification of a password
          # for database connections without requiring interactive input. This variable can be set in
          # the shell before executing PostgreSQL client applications like psql or pg_dump.
          export PGPASSWORD=$REPLICATION_PASSWORD
          # https://www.postgresql.org/docs/current/app-pgbasebackup.html
          pg_basebackup -v -h $PGHOST -p 5432 -D $PGDATA -U replication -R -Xs -Fp
        fi
      else
        # cat /postgres/initdb/create-replication-user.sh
        # https://hub.docker.com/_/postgres#initialization-scripts
        if [ ! -d "/docker-entrypoint-initdb.d" ];
        then
          mkdir /docker-entrypoint-initdb.d
        fi
        cp /postgres/initdb/* /docker-entrypoint-initdb.d/
        printf "Adding credential...\n"
        sed -i 's/#POSTGRES_USER/$(POSTGRES_USER)/g' /docker-entrypoint-initdb.d/create-replication-user.sh
        sed -i 's/#POSTGRES_DB/$(POSTGRES_DB)/g' /docker-entrypoint-initdb.d/create-replication-user.sh
        sed -i 's/#REPLICATION_PASSWORD/$(REPLICATION_PASSWORD)/g' /docker-entrypoint-initdb.d/create-replication-user.sh
        printf "Changing permissions for emptyDir...\n"
        chown -v -R 1999:1999 /docker-entrypoint-initdb.d && chmod -R 750 /docker-entrypoint-initdb.d
        # cat /docker-entrypoint-initdb.d/create-replication-user.sh
      fi
      printf "Changing permissions for /wsf_data_dir...\n"
      chown -v -R 1999:1999 /wsf_data_dir && chmod -R 760 /wsf_data_dir
      EOT
    ]
    command = ["/bin/sh"]
    env = {
      PGDATA = var.postgres_data
      PGHOST = local.service_name_postgres_master
      # on - Secondary.
      # off - Primary.
      STANDBY_MODE = "on"
    }
    env_from_secrets = [
      "${local.statefulset_postgres_replica}-secret"
    ]
    image = var.postgres_image_tag
    image_pull_policy = "IfNotPresent"
    security_context = {
      run_as_user = 0
      run_as_group = 0
      run_as_non_root = false
      read_only_root_filesystem = false
      privileged = false
    }
    volume_mounts = [{
      name = "wsf"
      mount_path = "/wsf_data_dir"
      read_only = false
    }, {
      name = "readonly-initdb-volume"
      # https://hub.docker.com/_/postgres#initialization-scripts
      mount_path = "/postgres/initdb"
      read_only = false
    }, {
      name = "writable-initdb-volume"
      mount_path = "/docker-entrypoint-initdb.d"
      read_only = false
    }/*, {
      name = "writable-lib"
      mount_path = "/var/lib/pgsql"
      read_only = false
    }*/]
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
      command = ["/bin/sh",
        "-c",
        "pg_isready --host=$(POD_IP) --port=5432 --username=${var.postgres_user} --dbname=${var.postgres_db}"
      ]
    }
  }]
  namespace = local.namespace
  # Ensure that the non-root user running the container has the necessary group permissions to
  # access files in mounted volumes.
  pod_security_context = {
    fs_group = 1999
    run_as_non_root = true
    run_as_user = 1999
    run_as_group = 1999
  }
  readiness_probe = [{
    initial_delay_seconds = 30
    period_seconds = 5
    timeout_seconds = 3
    failure_threshold = 3
    # success_threshold = 1
    exec = {
      # https://www.postgresql.org/docs/current/app-pg-isready.html
      command = ["/bin/sh",
        "-c",
        "pg_isready --host=$(POD_IP) --port=5432 --username=${var.postgres_user} --dbname=${var.postgres_db}"
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
      "db" = var.postgres_db_label
    }
    # Plain-text data.
    data = {
      POSTGRES_DB = var.postgres_db
      POSTGRES_USER = var.postgres_user
      POSTGRES_PASSWORD = var.postgres_password
      REPLICATION_PASSWORD = var.postgres_replication_password
    }
    type = "Opaque"
    immutable = true
  }]
  service = {
    name = local.service_name_postgres_replica
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
      "svc_selector_label" = "svc-${local.service_name_postgres_replica}"
    }
    publish_not_ready_addresses = true
    type = "ClusterIP"
  }
  statefulset_name = local.statefulset_postgres_replica
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
    name = "config-volume"
    config_map_name = "${local.statefulset_postgres_replica}-postgres-conf-files"
    default_mode = "0760"
    items = [{
      key = "pg_hba.conf"
      path = "pg_hba.conf"
    }, {
      key = "pg_ident.conf"
      path = "pg_ident.conf"
    }, {
      key = "postgresql.conf"
      path = "postgresql.conf"
    }]
  }, {
    name = "readonly-initdb-volume"
    config_map_name = "${local.statefulset_postgres_replica}-postgres-script-files"
    default_mode = "0750"
  }]
  volume_mount = [{
    name = "wsf"
    mount_path = "/wsf_data_dir"
    read_only = false
  }, {
    name = "wsf"
    mount_path = "/var/run/postgresql"  # The path where Postgres stores its data.
    read_only = false
  }, {
    name = "config-volume"
    mount_path = "/postgres/config"
    read_only = false
  }, {
    name = "writable-initdb-volume"
    mount_path = "docker-entrypoint-initdb.d"
    read_only = false
  }]
  volume_empty_dir = [{
    name = "writable-initdb-volume"
  }]
}

module "fin-PostgresBackup" {
  count = var.db_postgres && !var.k8s_crds ? 1 : 0
  depends_on = [
    module.fin-PostgresMaster
  ]
  # Specify the location of the module, which contains the file main.tf.
  source = "./modules/cronjob"
  #
  cron_job = {
    metadata = {
      name = "postgres-backup-cronjon"
      labels = {
        "app" = var.app_name
        "db" = var.postgres_db_label
      }
      namespace = local.namespace
    }
    spec = {
      schedule = "* * * * *"
      job_template = {
        metadata = {
          name = "postgres-backup-job"
          labels = {
            "app" = var.app_name
            "db" = var.postgres_db_label
          }
          namespace = local.namespace
        }
        spec = {
          template = {
            spec = {
              container = [{
                name = "postgres-backup-container"
                command = ["/bin/sh",
                  "-c",
                  "date; echo Hello from K8s"
                ]
                image = var.busybox
                image_pull_policy = "IfNotPresent"
              }]
            }
          }
        }
      }
    }
  }
}

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
