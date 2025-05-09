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
# If the directory already exists, mkdir will not create it again and will not produce an error.
# To preserve linebreaks, quote the command or variable.
echo "$(mkdir --verbose --parents ../../../tf-states/IaC-app/)"
echo "$(chmod -v -R 700 ../../../tf-states/)"
echo "*************************"
echo "Deploying the application"
echo "*************************"
terraform init
terraform apply -auto-approve \
  -var "app_version=$APP_VERSION" \
  -var "k8s_manifest_crd=$K8S_MANIFEST_CRD"
echo "*********************************"
echo "Saving the vars file with secrets"
echo "*********************************"
echo "$(mkdir --verbose --parents ../../../tf-secret-vars/IaC-app/)"
echo "$(chmod -v -R 700 ../../../tf-secret-vars/)"
echo "$(cp -v --update --archive --backup=numbered tf_secrets.auto.tfvars ../../../tf-secret-vars/IaC-app/tf_secrets.auto.tfvars)"
echo "Done"
