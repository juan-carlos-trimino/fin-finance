#!/bin/bash
#
# Treats unset or undefined variables as an error when substituting (during parameter expansion).
# Does not apply to special parameters such as wildcard * or @.
set -u
# Checks that expected input environment variables are provided.
: "$APP_VERSION"
: "$K8S_MANIFEST_CRD"
echo "*********************"
echo "Environment variables"
echo "*********************"
echo "APP_VERSION = $(printenv APP_VERSION)"
echo "K8S_MANIFEST_CRD = $(printenv K8S_MANIFEST_CRD)"
echo "*****************"
echo "Current directory"
echo "*****************"
echo $(pwd)
echo "*************************************************************"
echo "Creating the directory structure for the Terraform state file"
echo "*************************************************************"
# To preserve linebreaks, quote the command or variable.
# If the directory already exists, mkdir will not create it again and will not produce an error.
# -v, --verbose => Print a message for each created directory.
# -p, --parents => No error if existing, make parent directories as needed.
echo "$(mkdir -v -p ../../../tf-states/IaC-app/)"
# -v, --verbose => Output a diagnostic for every file processed.
# -R, --recursive => Change files and directories recursively.
echo "$(chmod -v -R 700 ../../../tf-states/)"
echo "*********************************"
echo "Saving the vars file with secrets"
echo "*********************************"
echo "$(mkdir -v -p ../../../tf-secret-vars/IaC-app/)"
echo "$(chmod -v -R 700 ../../../tf-secret-vars/)"
# -v, --verbose => Explain what is being done.
# -a, --archive	=> Preserve the source's metadata, such as creation date, permissions, and
#                  extended attributes.
# --backup=simple => Make a backup of each existing destination file.
#   simple, never => Always make simple backups
# -S, --suffix=SUFFIX => Override the usual backup suffix.
echo "$(cp -v -a --backup="simple" -S=".bak" tf_secrets.auto.tfvars ../../../tf-secret-vars/IaC-app/tf_secrets.auto.tfvars)"
echo "*************************"
echo "Deploying the application"
echo "*************************"
terraform init
terraform apply -auto-approve \
  -var "app_version=$APP_VERSION" \
  -var "k8s_manifest_crd=$K8S_MANIFEST_CRD"
echo "*********************"
echo "Copying the lock file"
echo "*********************"
# Copy lock file so that it can be saved in the repo.
if [ -f ../../../tf-states/IaC-app/.terraform.lock.hcl ]
then
  echo "$(cp -v -a ../../../tf-states/IaC-app/.terraform.lock.hcl .terraform.lock.hcl)"
else
  echo "The lock file does not exists."
fi
echo "Done..."
