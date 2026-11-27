#!/bin/bash
# The '-e' option tells the shell to exit immediately if any invoked command exits with a non-zero status, which is the convention
# for indicating errors in Unix-like operating systems. It's a handy option for scripting, as it helps to catch errors and bugs early.
#
# The '-u' option treats unset or undefined variables as an error when substituting (during parameter expansion).
# Does not apply to special parameters such as wildcard * or @.
set -euo pipefail

#
# The "Dot Space Script" Trick for Running Bash Scripts in the "Current Shell"
# In Bash, a dot (".") is a command, which is synonymous with the "source" command. It executes the
# content of a specified script in the current shell environment. This differs from running a
# script directly, which executes in a subshell. Therefore, when using ".", any changes made by the
# script, such as variables assignments or directory changes, will persist in the current shell. To
# persist the environment variables, run the script like so
# $ . ./ev_app.sh
# or
# $ source ./ev_app.sh
# Please note that "source" is a synonym for dot (".") in the Bash shell, but not in POSIX sh, so
# for maximum compatibility use the period.
#

display_help() {
  printf 'Usage: %s <deploy|destroy> [OPTIONS]\n\n' "$0"
  printf "Commands:\n"
  printf "  deploy                  Plan or apply changes to the infrastructure\n"
  printf "  destroy                 Destroy the managed infrastructure\n\n"
  printf "Options:\n"
  printf "  -p, --plan              Run terraform plan (deploy only)\n"
  printf "  -a, --apply             Run terraform apply (default for deploy)\n"
  printf "                          terraform plan and terraform apply require the exact same input variables\n"
  printf "  -b, --build             Build the image (deploy only)\n"
  printf "  -n, --no-build          Use an image from the repo (deploy only)\n"
  printf "                          (The image has to exist in the repo)\n"
  printf "  -e, --empty-dir         Use empty-dir (deploy only)\n"
  printf "  -i, --persistent-disk   Use persistent disk (deploy only)\n"
  printf "  -d, --postgres          Use the postgres database (deploy only)\n"
  printf "  -f, --pprof             If present, pprof is enabled (deploy only)\n"
  printf "                          (If not present, pprof is disabled)\n"
  printf "  -g, --gateway           If present, use a reverse proxy (deploy only)\n"
  printf "                          (If not present, reverse proxy is not used)\n"
  printf "  -v, --version <ver>     Application version; e.g., 1.0.0  (deploy only)\n"
  printf "  -y, --auto-approve      Skip interactive approval (apply or destroy)\n"
  printf "  -h, --help              Display help\n\n"
  printf "Examples:\n"
  printf "Deploy app by building the image, version 1.0.0, and empty-dir as storage.\n"
  printf "$ ./iac.sh deploy -b -v 1.0.0 -e\n\n"
  printf "Deploy app by not building an image, version 1.0.1, persistent storage, and enable the endpoints of pprof.\n"
  printf "$ ./iac.sh deploy -n -v 1.0.1 -i -f\n\n"
  printf "Deploy app by not building an image, version 1.0.0, persistent storage, and reverse proxy.\n"
  printf "$ ./iac.sh deploy -n -v 1.0.0 -i -g\n\n"
  exit 0
}

check_options() {
  local arr=("$@")
  local -i size="${#arr[@]}"
  if [[ "$size" -eq 0 ]];
  then
    display_help
    exit 1
  fi
  local command="${arr[0]}"
  local deploy_mode=""
  local storage_type=""
  local auto_approve=false
  # Initialize all global variables with safe defaults (so destroy & set -u work cleanly)
  action=""
  build_image=""
  app_version=""
  use_postgres=false
  enable_pprof=false
  use_gateway=false
  persistent_disk=false
  empty_dir=false
  case "$command" in
    -h|--help)
      display_help
      ;;
    deploy)
      action="apply"
      ;;
    destroy)
      action="destroy"
      ;;
    *)
      printf "Unknown command: %s\n\n" "$command" >&2
      display_help
      ;;
  esac
  local -i ndx=1
  # Parse remaining flags starting at index 1
  while [[ $ndx -lt $size ]];
  do
    local flag="${arr[ndx]}"
    case "$flag" in
      -h|--help)
        display_help
        ;;
      -y|--auto-approve)
        auto_approve=true
        ;;
      -p|--plan)
        if [[ "$command" != "deploy" ]];
        then
          printf "Error: %s is only valid with 'deploy'\n" "$flag" >&2
          exit 1
        fi
        action="plan"
        ;;
      -a|--apply)
        [[ "$command" == "deploy" ]] || { printf "Error: %s is only valid with 'deploy'\n" "$flag" >&2; exit 1; }
        [[ -z "$deploy_mode" ]] || { printf "Error: Cannot combine -p/--plan and -a/--apply\n" >&2; exit 1; }
        deploy_mode="apply"
        ;;
      -b|--build)
        [[ "$command" == "deploy" ]] || { printf "Error: %s is only valid with 'deploy'\n" "$flag" >&2; exit 1; }
        [[ -z "$build_image" ]] || { printf "Error: Cannot combine -b/--build and -n/--no-build\n" >&2; exit 1; }
        build_image=true
        ;;
      -n|--no-build)
        [[ "$command" == "deploy" ]] || { printf "Error: %s is only valid with 'deploy'\n" "$flag" >&2; exit 1; }
        [[ -z "$build_image" ]] || { printf "Error: Cannot combine -b/--build and -n/--no-build\n" >&2; exit 1; }
        build_image=false
        ;;
      -e|--empty-dir)
        [[ "$command" == "deploy" ]] || { printf "Error: %s is only valid with 'deploy'\n" "$flag" >&2; exit 1; }
        [[ -z "$storage_type" ]] || { printf "Error: Cannot combine -e/--empty-dir and -i/--persistent-disk\n" >&2; exit 1; }
        storage_type="empty-dir"
        ;;
      -i|--persistent-disk)
        [[ "$command" == "deploy" ]] || { printf "Error: %s is only valid with 'deploy'\n" "$flag" >&2; exit 1; }
        [[ -z "$storage_type" ]] || { printf "Error: Cannot combine -e/--empty-dir and -i/--persistent-disk\n" >&2; exit 1; }
        storage_type="persistent-disk"
        ;;
      -d|--postgres)
        [[ "$command" == "deploy" ]] || { printf "Error: %s is only valid with 'deploy'\n" "$flag" >&2; exit 1; }
        use_postgres=true
        ;;
      -f|--pprof)
        [[ "$command" == "deploy" ]] || { printf "Error: %s is only valid with 'deploy'\n" "$flag" >&2; exit 1; }
        enable_pprof=true
        ;;
      -g|--gateway)
        [[ "$command" == "deploy" ]] || { printf "Error: %s is only valid with 'deploy'\n" "$flag" >&2; exit 1; }
        use_gateway=true
        ;;
      -v|--version)
        [[ "$command" == "deploy" ]] || { printf "Error: %s is only valid with 'deploy'\n" "$flag" >&2; exit 1; }
        ndx=$(( ndx + 1 ))
        if [[ $ndx -ge $size || "${arr[ndx]}" == -* ]];
        then
          printf "Error: %s requires a version argument (e.g., -v 1.0.0)\n" "$flag" >&2
          exit 1
        fi
        app_version="${arr[ndx]}"
        ;;
      *)
        printf "Error: Unknown flag '%s'\n\n" "$flag" >&2
        exit 1
        ;;
    esac
    ndx=$(( ndx + 1 ))
  done
  # Validate required deploy options
  if [[ "$command" == "deploy" ]];
  then
    [[ -n "$app_version" ]] || { printf "Error: -v/--version is required for 'deploy'\n" >&2; exit 1; }
    [[ -n "$build_image" ]] || { printf "Error: Either -b/--build or -n/--no-build is required for 'deploy'\n" >&2; exit 1; }
    [[ -n "$storage_type" ]] || { printf "Error: Either -e/--empty-dir or -i/--persistent-disk is required for 'deploy'\n" >&2; exit 1; }
  fi
  #
  if [[ "$storage_type" == "persistent-disk" ]];
  then
    persistent_disk=true
  elif [[ "$storage_type" == "empty-dir" ]];
  then
    empty_dir=true
  fi
  # Build global Terraform arguments array
  tf_args=()
  if [[ "$auto_approve" == true && "$action" != "plan" ]];
  then
    tf_args+=("-auto-approve")
  fi
}

print_time_elapsed() {
  printf "\n*************************************************\n"
  printf "Done...\n"
  local end_time
  end_time=$(date)
  printf '%s\n' "$end_time"
  start_time_seconds=$(date -d "$1" +"%s")
  end_time_seconds=$(date -d "$end_time" +"%s")
  duration=$(( $end_time_seconds - $start_time_seconds ))
  printf "Time Elapsed: %02d hours %02d minutes and %02d seconds." "$(($duration / 3600))" \
         "$(($duration % 3600 / 60))" "$(($duration % 60))"
  printf "\n*************************************************\n\n"
  return
}

# Main program
printf "\n"
check_options "$@"
printf "****************************\n"
printf "Starting deployment...\n"
start_time=$(date)
printf "$start_time\n"
printf '****************************\n\n'
printf "*************************************************************\n"
printf "Creating the directory structure for the Terraform state file\n"
printf "*************************************************************\n"
# Because ../../../tf-states/IaC-app/ is a relative path, it resolves relative to whatever directory the shell is in when
# the script runs -- not necessarily the directory where the .sh file lives. If you want it anchored to your script's location
# regardless of where you invoke it from, resolve the script directory first.
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ "$1" == "destroy" ]
then
  printf "******************************\n"
  printf "Removing the infrastructure...\n"
  printf "******************************\n"
  # Terraform automatically loads ANY file ending in .auto.tfvars during EVERY command (plan, apply, and destroy); hence, you
  # explicitly pass -var-file="last-deploy.tfvars" to terraform destroy. Only destroy ever reads the saved file -- zero risk
  # of plan or apply accidentally picking up stale flags from a previous deployment.
  terraform -chdir="${script_dir}" "$action" "${tf_args[@]}" -var-file="last-deploy.tfvars"
  # See https://cert-manager.io/docs/installation/helm/
  # Uncomment and delete them (with --ignore-not-found=true) only when you are decommissioning a dedicated, single-purpose cluster
  # environment and want a 100% clean slate.
  # printf "\n*****************************\n"
  # printf "Deleting cert-manager's CRDs.\n"
  # printf "*****************************\n"
  # kubectl delete crd \
  #   issuers.cert-manager.io \
  #   clusterissuers.cert-manager.io \
  #   certificates.cert-manager.io \
  #   certificaterequests.cert-manager.io \
  #   orders.acme.cert-manager.io \
  #   challenges.acme.cert-manager.io
  printf "\n*****************************\n"
  printf "Deleting traefik's CRDs.\n"
  printf "*****************************\n"
  kubectl delete crd --ignore-not-found=true \
    ingressroutes.traefik.io \
    ingressroutetcps.traefik.io \
    ingressrouteudps.traefik.io \
    middlewares.traefik.io \
    middlewaretcps.traefik.io \
    serverstransports.traefik.io \
    serverstransporttcps.traefik.io \
    tlsoptions.traefik.io \
    tlsstores.traefik.io \
    traefikservices.traefik.io \
    accesscontrolpolicies.hub.traefik.io \
    aiservices.hub.traefik.io \
    apibundles.hub.traefik.io \
    apicatalogitems.hub.traefik.io \
    apiplans.hub.traefik.io \
    apiportals.hub.traefik.io \
    apiratelimits.hub.traefik.io \
    apis.hub.traefik.io \
    apiversions.hub.traefik.io \
    managedsubscriptions.hub.traefik.io
    #   gatewayclasses.gateway.networking.k8s.io \
    #   gateways.gateway.networking.k8s.io \
    #   grpcroutes.gateway.networking.k8s.io \
    #   httproutes.gateway.networking.k8s.io \
    #   referencegrants.gateway.networking.k8s.io
  rm -f "${script_dir}/last-deploy.tfvars"
  print_time_elapsed "$start_time"
elif [ "$1" == "deploy" ]
then
  # Verify all variables are bound
  : "$action"
  : "$build_image"
  : "$app_version"
  : "$use_postgres"
  : "$enable_pprof"
  : "$use_gateway"
  : "$persistent_disk"
  : "$empty_dir"
  #
  printf "*********************\n"
  printf "Debug - variables    \n"
  printf "*********************\n"
  printf "action = $action\n"
  printf "build_image = $build_image\n"
  printf "app_version = $app_version\n"
  printf "use_postgres = $use_postgres\n"
  printf "enable_pprof = $enable_pprof\n"
  printf "use_gatewate = $use_gateway\n"
  printf "persistent_disk = $persistent_disk\n"
  printf "empty_dir = $empty_dir\n\n"
  #
  printf "*****************\n"
  printf "Current directory\n"
  printf "*****************\n"
  printf '%s\n\n' "$(pwd)"
  printf "*********************************\n"
  printf "Saving the vars file with secrets\n"
  printf "*********************************\n"
  # If the directory already exists, mkdir will not create it again and will not produce an error.
  # -v, --verbose => Print a message for each created directory.
  # -p, --parents => No error if existing, make parent directories as needed.
  # Adding -- before the path in mkdir ensures that even if a parent directory name somehow starts with a hyphen (-), mkdir won't
  # mistake it for a command-line flag.
  mkdir -v -p -- "${script_dir}/../../../tf-secret-vars/IaC-app/" || exit 1
  # -v, --verbose => Explain what is being done.
  # -a, --archive	=> Preserve the source's metadata, such as creation date, permissions, and
  #                  extended attributes.
  # --backup=simple => Make a backup of each existing destination file.
  #   simple, never => Always make simple backups
  # -S "SUFFIX", --suffix=SUFFIX => Override the usual backup suffix.
  cp -v -a --backup="simple" -S ".bak" -- \
     "${script_dir}/tf_secrets.auto.tfvars" \
     "${script_dir}/../../../tf-secret-vars/IaC-app/tf_secrets.auto.tfvars"
  # Lock down permissions AFTER copying the secret file and any .bak backup
  # -v, --verbose => Output a diagnostic for every file processed.
  # -R, --recursive => Change files and directories recursively.
  chmod -v -R 700 -- "${script_dir}/../../../tf-secret-vars/" || exit 1
  printf "*************************************************************\n"
  printf "Creating the directory structure for the Terraform state file\n"
  printf "*************************************************************\n"
  tfstate_dir="$(realpath -m "${script_dir}/../../../tf-states/IaC-app")"
  mkdir -v -p -- "${tfstate_dir}" || exit 1
  chmod -v -R 700 -- "${tfstate_dir}" || exit 1
  printf "**********************\n"
  printf "Initializing Terraform\n"
  printf "**********************\n"
  terraform -chdir="${script_dir}" init \
            -reconfigure \
            -input=false \
            -backend-config="path=${tfstate_dir}/terraform.tfstate"
  common_vars=(
    -var "build_image=$build_image"
    -var "app_version=$app_version"
    -var "db_postgres=$use_postgres"
    -var "pprof=$enable_pprof"
    -var "reverse_proxy=$use_gateway"
    -var "persistent_disk=$persistent_disk"
    -var "empty_dir=$empty_dir"
  )
  # Only run this step when the Traefik CRD is not already present in the cluster.
  if [[ "$use_gateway" == "true" ]] && ! kubectl get crd ingressroutes.traefik.io &>/dev/null;
  then
    printf "\n*****************************************\n"
    printf "Creating CustomResourceDefinitions (CRDs)\n"
    printf "*****************************************\n"
    # When action=apply (with auto_approve=true)
    # terraform -chdir="${script_dir}" apply -auto-approve "${tf_args[@]}" ... passes -auto-approve twice, which Terraform
    # accepts without any issue.
    terraform -chdir="${script_dir}" apply -auto-approve \
      "${tf_args[@]}" \
      "${common_vars[@]}" \
      -var "k8s_crds=true"
  fi
  printf "\n*************************"
  printf "\nDeploying the application"
  printf "\n*************************\n"
  terraform -chdir="${script_dir}" "$action" \
    "${tf_args[@]}" \
    "${common_vars[@]}" \
    -var "k8s_crds=false"
  # To ensure no other local user can even peek inside ${tfstate_dir} while Terraform is running, lock the folder immediately
  # after mkdir, and run chmod again at the very end after terraform apply finishes writing the state files.
  chmod -v -R 700 -- "${tfstate_dir}" || exit 1
  # If action is apply, save current variables for a future destroy.
  if [ "$action" == "apply" ];
  then
    # Add last-deploy.tfvars to your .gitignore file so Git ignores it and never commits it to your repository.
    # You need quotes around "${app_version}" because it is a string (e.g., "1.0.0" or "latest"), whereas booleans (true/false)
    # remain unquoted.
    cat <<EOF > "${script_dir}/last-deploy.tfvars"
build_image     = ${build_image}
app_version     = "${app_version}"
db_postgres     = ${use_postgres}
pprof           = ${enable_pprof}
reverse_proxy   = ${use_gateway}
persistent_disk = ${persistent_disk}
empty_dir       = ${empty_dir}
k8s_crds        = false
EOF
  fi
  print_time_elapsed "$start_time"
else
  printf 'Command "%s" -> Not implemented\n\n' "$1"
  display_help
fi
