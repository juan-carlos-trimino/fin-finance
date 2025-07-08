#!/bin/bash
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
  # -e   Enable interpretation of backslash escapes.
  # -n   Do not output the trailing newline.
  printf "Usage: $0 deploy/destroy [OPTIONS]\n"
  printf "  -h, --help                Display help\n"
  printf "  deploy                  Apply changes to the infrastructure.\n"
  printf "    -b, --build             Build the image\n"
  printf "    -nb, --no-build         Use an image from the repo\n"
  printf "                            (The image has to exist in the repo)\n"
  printf "    -v, --version           Application version; e.g., 1.0.0\n"
  printf "    -ed, --empty-dir        Use empty-dir\n"
  printf "    -pd, --persistent-disk  Use persistent disk\n"
  printf "    -gw, --gateway          If present, use a reverse proxy\n"
  printf "                            (If not present, reverse proxy is not used)\n"
  printf "    -p, --pprof             If present, pprof is enabled\n"
  printf "                            (If not present, pprof is disabled)\n"
  printf "  destroy                 Destroy the infrastructure.\n\n"
  printf "Examples:\n"
  printf "Deploy app using empty-dir as storage.\n"
  printf "Build the image using v1.0.0.\n"
  printf "$ ./iac.sh deploy -b -v 1.0.0 -ed\n\n"
  printf "Deploy app using persistent storage, enable the endpoints of pprof.\n"
  printf "Use image in repo v1.0.1.\n"
  printf "$ ./iac.sh deploy -nb -v 1.0.1 -pd -p\n\n"
  printf "Deploy app using persistent storage and the reverse proxy.\n"
  printf "Use image in repo v1.0.0.\n"
  printf "$ ./iac.sh deploy -nb -v 1.0.0 -pd -gw\n\n"
  exit 0
}

check_options() {
  # echo "Elements in \$@: $@"
  arr=()
  for arg in "$@"
  do
    arr+=("$arg")
  done
  # echo "Elements in arr: ${arr[@]}"
  local -i ndx=0
  local -i size="${#arr[@]}"
  local app_version=""
  local build_image=""
  local reverse_proxy="false"
  local k8s_crds=$reverse_proxy
  local pprof="false"
  local deployment_type=""
  for (( ndx = 1; ndx < size; ))  # Move pass the first argument.
  do
    flag=${arr[ndx]}
    # echo "Element at index $ndx: $flag"
    case "$flag" in
      "-h" | "--help")
        display_help
        ;;
      "-b" | "--build")
        build_image="true"
        ;;
      "-nb" | "--no-build")
        build_image="false"
        ;;
      "-v" | "--version")
        ndx=$(( ndx + 1 ))
        # echo "Element at index $ndx: ${arr[ndx]}"
        app_version=${arr[ndx]}
        ;;
      "-p" | "--pprof")
        pprof="true"
        ;;
      "-ed" | "--empty-dir")
        deployment_type="empty-dir"
        ;;
      "-pd" | "--persistent-disk")
        deployment_type="persistent-disk"
        ;;
      "-gw" | "--gateway")
        reverse_proxy="true"
        k8s_crds=$reverse_proxy
        ;;
      *)  #Default.
        echo -e "Unknown flag: $flag\n"
        exit 1
        ;;
    esac
    ndx=$(( ndx + 1 ))
  done
  #
  if [ ${arr[0]} == "deploy" ]
  then
    if [ "$build_image" == "" ]
    then
      printf "One of these two options is required:\n"
      printf "  -b or --build.\n"
      printf "  -nb or --no-build.\n\n"
      exit 1
    elif [ "$app_version" == "" ]
    then
      echo -e "The option -v/--version is required.\n"
      exit 1
    elif [ "$deployment_type" == "" ]
    then
      printf "One of these two options is required:\n"
      printf "  -ed or --empty-dir.\n"
      printf "  -pd or --persistent-disk.\n\n"
      exit 1
    fi
    export APP_VERSION=$app_version
    export REVERSE_PROXY=$reverse_proxy
    export K8S_CRDS=$k8s_crds
    export DEPLOYMENT_TYPE=$deployment_type
    export PPROF=$pprof
    export BUILD_IMAGE=$build_image
  elif [ ${arr[0]} == "destroy" ]
  then
    echo ""  # Placeholder...
  else
    display_help
    exit 1
  fi
  return
}

print_time_elapsed() {
  printf "\n*************************************************\n"
  echo "Done..."
  local end_time=$(date)
  echo "$end_time"
  start_time_seconds=$(date -d "$1" +"%s")
  end_time_seconds=$(date -d "$end_time" +"%s")
  duration=$(( $end_time_seconds - $start_time_seconds ))
  printf "Time Elapsed: %02d hours %02d minutes and %02d seconds." "$(($duration / 3600))" \
         "$(($duration % 3600 / 60))" "$(($duration % 60))"
  printf "\n*************************************************\n\n"
  return
}

# Main program
echo
check_options $@
echo "****************************"
echo "Starting deployment..."
start_time=$(date)
echo "$start_time"
echo -e "****************************\n"
if [ "$1" == "destroy" ]
then
  declare -i arguments="$#"
  if (( arguments > 1 ))  # Number of arguments > 1?
  then
    display_help
  else
    printf "******************************\n"
    printf "Removing the infrastructure...\n"
    printf "******************************\n"
    terraform destroy -auto-approve
    # See https://cert-manager.io/docs/installation/helm/
    printf "\n*****************************"
    printf "\nDeleting cert-manager's CRDs."
    printf "\n*****************************\n"
    # kubectl delete crd \
    #   issuers.cert-manager.io \
    #   clusterissuers.cert-manager.io \
    #   certificates.cert-manager.io \
    #   certificaterequests.cert-manager.io \
    #   orders.acme.cert-manager.io \
    #   challenges.acme.cert-manager.io
    printf "\n*****************************"
    printf "\nDeleting traefik's CRDs."
    printf "\n*****************************\n"
    kubectl delete crd \
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
    print_time_elapsed "$start_time"
  fi
elif [ "$1" == "deploy" ]
then
  # Treats unset or undefined variables as an error when substituting (during parameter expansion).
  # Does not apply to special parameters such as wildcard * or @.
  set -u
  # Checks that expected input environment variables are provided.
  : "$APP_VERSION"
  : "$K8S_CRDS"
  : "REVERSE_PROXY"
  : "$DEPLOYMENT_TYPE"
  : "$PPROF"
  : "$BUILD_IMAGE"
  echo "*********************"
  echo "Environment variables"
  echo "*********************"
  echo "APP_VERSION = $(printenv APP_VERSION)"
  echo "REVERSE_PROXY = $(printenv REVERSE_PROXY)"
  echo "K8S_CRDS = $(printenv K8S_CRDS)"
  echo "DEPLOYMENT_TYPE = $(printenv DEPLOYMENT_TYPE)"
  echo "PPROF = $(printenv PPROF)"
  printf "BUILD_IMAGE = $(printenv BUILD_IMAGE)\n\n"
  echo "*****************"
  echo "Current directory"
  echo "*****************"
  echo -e "$(pwd)\n"
  echo "*************************************************************"
  echo "Creating the directory structure for the Terraform state file"
  echo "*************************************************************"
  # To preserve linebreaks, quote the command or variable.
  # If the directory already exists, mkdir will not create it again and will not produce an error.
  # -v, --verbose => Print a message for each created directory.
  # -p, --parents => No error if existing, make parent directories as needed.
  echo -n "$(mkdir -v -p ../../../tf-states/IaC-app/)"
  # -v, --verbose => Output a diagnostic for every file processed.
  # -R, --recursive => Change files and directories recursively.
  echo -e "$(chmod -v -R 700 ../../../tf-states/)\n"
  echo "*********************************"
  echo "Saving the vars file with secrets"
  echo "*********************************"
  echo -n "$(mkdir -v -p ../../../tf-secret-vars/IaC-app/)"
  echo "$(chmod -v -R 700 ../../../tf-secret-vars/)"
  # -v, --verbose => Explain what is being done.
  # -a, --archive	=> Preserve the source's metadata, such as creation date, permissions, and
  #                  extended attributes.
  # --backup=simple => Make a backup of each existing destination file.
  #   simple, never => Always make simple backups
  # -S, --suffix=SUFFIX => Override the usual backup suffix.
  echo -e "$(cp -v -a --backup="simple" -S=".bak" tf_secrets.auto.tfvars ../../../tf-secret-vars/IaC-app/tf_secrets.auto.tfvars)\n"
  echo "**********************"
  echo "Initializing Terraform"
  echo "**********************"
  terraform init
  if [ "$REVERSE_PROXY" == "true" ]
  then
    echo -e "\n*****************************************"
    echo "Creating CustomResourceDefinitions (CRDs)"
    printf "*****************************************\n"
    terraform apply -auto-approve \
      -var "reverse_proxy=$REVERSE_PROXY" \
      -var "k8s_crds=$K8S_CRDS"
    export K8S_CRDS="false"
  fi
  printf "\n*************************"
  printf "\nDeploying the application"
  printf "\n*************************\n"
  terraform apply -auto-approve \
    -var "app_version=$APP_VERSION" \
    -var "reverse_proxy=$REVERSE_PROXY" \
    -var "k8s_crds=$K8S_CRDS" \
    -var "deployment_type=$DEPLOYMENT_TYPE" \
    -var "pprof=$PPROF" \
    -var "build_image=$BUILD_IMAGE"
  printf "\n*********************\n"
  echo "Copying the lock file"
  echo "*********************"
  # Copy lock file so that it can be saved in the repo.
  if [ -f ../../../tf-states/IaC-app/.terraform.lock.hcl ]
  then
    echo "$(cp -v -a ../../../tf-states/IaC-app/.terraform.lock.hcl .terraform.lock.hcl)"
  else
    echo "The lock file does not exist."
  fi
  print_time_elapsed "$start_time"
else
  display_help
fi
