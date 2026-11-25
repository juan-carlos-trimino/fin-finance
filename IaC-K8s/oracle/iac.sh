#!/bin/bash
# The '-e' option tells the shell to exit immediately if any invoked command exits with a non-zero status, which is the convention
# for indicating errors in Unix-like operating systems. It's a handy option for scripting, as it helps to catch errors and bugs early.
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

# Global variable (do NOT use 'local' outside a function)
# action=""

display_help() {
  printf "Usage: $0 <deploy|destroy> [OPTIONS]\n\n"
  printf "Commands:\n"
  printf "  deploy                    Plan or apply changes to the infrastructure\n"
  printf "  destroy                   Destroy the managed infrastructure\n\n"
  printf "Options:\n"
  printf "  -p, --plan                Run terraform plan (deploy only)\n"
  printf "  -a, --apply               Run terraform apply (default for deploy)\n"
  printf "  -y, --auto-approve        Skip interactive approval (apply or destroy)\n"
  printf "  -h, --help                Display this help message\n\n"
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
  # Validate the first argument (arr[0])
  local command="${arr[0]}"
  local action=""
  local auto_approve=false
  case "$command" in
    -h|--help)
      display_help
      ;;
    deploy)
      action="apply"  # Default action if no -p/-a flag is given
      ;;
    destroy)
      action="destroy"
      ;;
    *)  # Default.
      printf "Unknown command: %s\n\n" "$command" >&2
      display_help
      ;;
  esac
  local -i ndx=1
  # Parse remaining flags starting at index 1
  for (( ndx = 1; ndx < size; ndx++ ));
  do
    local flag="${arr[ndx]}"
    case "$flag" in
      -h|--help)
        display_help
        ;;
      -p|--plan)
        if [[ "$command" != "deploy" ]];
        then
          printf "Error: %s can only be used with 'deploy'\n" "$flag" >&2
          exit 1
        fi
        action="plan"
        ;;
      -a|--apply)
        if [[ "$command" != "deploy" ]]; then
          printf "Error: %s can only be used with 'deploy'\n" "$flag" >&2
          exit 1
        fi
        action="apply"
        ;;
      -y|--auto-approve)
        auto_approve=true
        ;;
      *)
        printf "Unknown flag: %s\n\n" "$flag" >&2
        display_help
        ;;
    esac
  done
  # Build the global array for Terraform.
  # Variables assigned inside a function without the local keyword are global by default.
  tf_args=("$action")
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
echo
check_options "$@"
printf "****************************\n"
printf "Starting deployment...\n"
start_time=$(date)
printf "$start_time\n"
printf '****************************\n\n'
if [ "$1" == "deploy" ]
then
  printf "*****************\n"
  printf "Current directory\n"
  printf "*****************\n"
  printf '%s\n\n' "$(pwd)"
  # To preserve linebreaks, quote the command or variable.
  #
  # Because ../../../../tf-states/IaC-K8s/ is a relative path, it resolves relative to whatever directory the shell is in when
  # the script runs -- not necessarily the directory where the .sh file lives. If you want it anchored to your script's location
  # regardless of where you invoke it from, resolve the script directory first.
  script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  printf "*********************************\n"
  printf "Saving the vars file with secrets\n"
  printf "*********************************\n"
  mkdir -v -p -- "${script_dir}/../../../../tf-secret-vars/IaC-K8s/" || exit 1
  # -v, --verbose => Explain what is being done.
  # -a, --archive	=> Preserve the source's metadata, such as creation date, permissions, and
  #                  extended attributes.
  # --backup=simple => Make a backup of each existing destination file.
  #   simple, never => Always make simple backups
  # -S "SUFFIX", --suffix=SUFFIX => Override the usual backup suffix.
  cp -v -a --backup="simple" -S ".bak" -- \
     "${script_dir}/tf_secrets.auto.tfvars" \
     "${script_dir}/../../../../tf-secret-vars/IaC-K8s/tf_secrets.auto.tfvars"
  # Lock down permissions AFTER copying the secret file and any .bak backup
  chmod -v -R 700 -- "${script_dir}/../../../../tf-secret-vars/" || exit 1
  printf "*************************************************************\n"
  printf "Creating the directory structure for the Terraform state file\n"
  printf "*************************************************************\n"
  tfstate_dir="$(realpath -m "${script_dir}/../../../../tf-states/IaC-K8s")"
  # If the directory already exists, mkdir will not create it again and will not produce an error.
  # -v, --verbose => Print a message for each created directory.
  # -p, --parents => No error if existing, make parent directories as needed.
  # Adding -- before the path in mkdir ensures that even if a parent directory name somehow starts with a hyphen (-), mkdir won't
  # mistake it for a command-line flag.
  mkdir -v -p -- "${tfstate_dir}" || exit 1
  # -v, --verbose => Output a diagnostic for every file processed.
  # -R, --recursive => Change files and directories recursively.
  chmod -v -R 700 -- "${tfstate_dir}" || exit 1
  printf "**********************\n"
  printf "Initializing Terraform\n"
  printf "**********************\n"
  terraform -chdir="${script_dir}" init \
            -reconfigure \
            -input=false \
            -backend-config="path=${tfstate_dir}/terraform.tfstate"
  printf "\n*************************"
  printf "\nDeploying the application"
  printf "\n*************************\n"
  terraform -chdir="${script_dir}" "${tf_args[@]}"
  # To ensure no other local user can even peek inside ${tfstate_dir} while Terraform is running, lock the folder immediately after
  # mkdir, and run chmod again at the very end after terraform apply finishes writing the state files.
  chmod -v -R 700 -- "${tfstate_dir}" || exit 1
  print_time_elapsed "$start_time"
elif [ "$1" == "destroy" ]
then
  printf 'Command "%s" -> Not implemented\n\n' "$1"
  display_help
fi
