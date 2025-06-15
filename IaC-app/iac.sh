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
  echo -e "Usage: $0 deploy/destroy [OPTIONS]\n"
  echo -e "  -h, --help               Display help\n"
  echo "  deploy                   Apply changes to the infrastructure."
  echo "    -av, --app_version       Application version; e.g., 1.0.0 (default)"
  echo "    -dt, --deployment_type   Deployment type: empty-dir (default) or persistent-disk"
  echo "    -rp, --reverse_proxy     Use a reverse proxy; e.g., true or false (default)"
  echo -e "    -p, --pprof              Enable/disable pprof; e.g., true or false (default)\n"
  echo -e "  destroy                  Destroy the infrastructure.\n"
  exit 0
}

checkOptions() {
  # echo "Elements in \$@: $@"
  arr=()
  for arg in "$@"
  do
    arr+=("$arg")
  done
  # echo "Elements in arr: ${arr[@]}"
  local -i ndx=0
  local -i size="${#arr[@]}"
  local app_version="1.0.0"
  local reverse_proxy="false"
  local k8s_crds=$reverse_proxy
  local pprof="false"
  local deployment_type="empty-dir"
  for (( ndx = 1; ndx < size; ))  # Move pass the first argument.
  do
    flag=${arr[ndx]}
    # echo "Element at index $ndx: $flag"
    ndx=$(( ndx + 1 ))
    value=${arr[ndx]}
    # echo "Element at index $ndx: $value"
    ndx=$(( ndx + 1 ))
    case "$flag" in
      "-h" | "--help")
        display_help
        ;;
      "-av" | "--app_version")
        app_version=$value
        ;;
      "-p" | "--pprof")
        pprof=$value
        ;;
      "-dt" | "--deployment_type")
        deployment_type=$value
        ;;
      "-rp" | "--reverse_proxy")
        reverse_proxy=$value
        k8s_crds=$reverse_proxy
        ;;
      *)  #Default.
        echo -e "Unknown flag: $flag\n"
        display_help
        ;;
    esac
  done
  #
  if [ "$reverse_proxy" != "true" ] && [ "$reverse_proxy" != "false" ]
  then
    echo -e "Valid values for the flag -rp/--reverse_proxy: 'true' or 'false'.\n"
    exit 1
  elif [ "$pprof" != "true" ] && [ "$pprof" != "false" ]
  then
    echo -e "Valid values for the flag -p/--pprof: 'true' or 'false'.\n"
    exit 1
  elif [ "$deployment_type" != "empty-dir" ] && [ "$deployment_type" != "persistent-disk" ]
  then
    echo -e "Valid values for the flag -dt/--deployment_type: 'empty-dir' or 'persistent-disk'.\n"
    exit 1
  else
    export APP_VERSION=$app_version
    export REVERSE_PROXY=$reverse_proxy
    export K8S_CRDS=$k8s_crds
    export DEPLOYMENT_TYPE=$deployment_type
    export PPROF=$pprof
  fi
  return
}

# Main program
echo -e "\n"
if [ "$1" == "destroy" ]
then
  declare -i arguments="$#"
  if (( arguments > 1 ))  # Number of arguments > 1?
  then
    display_help
  else
    terraform destroy -auto-approve
    echo -e "\n"
  fi
elif [ "$1" == "deploy" ]
then
  checkOptions $@
  echo "****************************"
  echo "Starting deployment..."
  start_time=$(date)
  echo "$start_time"
  echo -e "****************************\n"
  # Treats unset or undefined variables as an error when substituting (during parameter expansion).
  # Does not apply to special parameters such as wildcard * or @.
  set -u
  # Checks that expected input environment variables are provided.
  : "$APP_VERSION"
  : "$K8S_CRDS"
  : "REVERSE_PROXY"
  : "$DEPLOYMENT_TYPE"
  : "$PPROF"
  echo "*********************"
  echo "Environment variables"
  echo "*********************"
  echo "APP_VERSION = $(printenv APP_VERSION)"
  echo "REVERSE_PROXY = $(printenv REVERSE_PROXY)"
  echo "K8S_CRDS = $(printenv K8S_CRDS)"
  echo "DEPLOYMENT_TYPE = $(printenv DEPLOYMENT_TYPE)"
  echo -e "PPROF = $(printenv PPROF)\n"
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
    echo "*****************************************"
    echo "Creating CustomResourceDefinitions (CRDs)"
    echo "*****************************************"
    terraform apply -auto-approve \
      -var "reverse_proxy=$REVERSE_PROXY" \
      -var "k8s_crds=$K8S_CRDS"
    export K8S_CRDS="false"
  fi
  echo "*************************"
  echo "Deploying the application"
  echo "*************************"
  terraform apply -auto-approve \
    -var "app_version=$APP_VERSION" \
    -var "reverse_proxy=$REVERSE_PROXY" \
    -var "k8s_crds=$K8S_CRDS" \
    -var "deployment_type=$DEPLOYMENT_TYPE" \
    -var "pprof=$PPROF"
  echo -e "\n*********************"
  echo "Copying the lock file"
  echo "*********************"
  # Copy lock file so that it can be saved in the repo.
  if [ -f ../../../tf-states/IaC-app/.terraform.lock.hcl ]
  then
    echo "$(cp -v -a ../../../tf-states/IaC-app/.terraform.lock.hcl .terraform.lock.hcl)"
  else
    echo "The lock file does not exist."
  fi
  echo -e "\n****************************"
  echo "Done..."
  end_time=$(date)
  echo "$end_time"
  start_time_seconds=$(date -d "$start_time" +"%s")
  end_time_seconds=$(date -d "$end_time" +"%s")
  duration=$(( $end_time_seconds - $start_time_seconds ))
  echo "Time Elapsed: $(($duration / 3600)) hours $(($duration % 3600 / 60)) minutes and $(($duration % 60)) seconds."
  echo -e "****************************\n"
else
  display_help
fi
