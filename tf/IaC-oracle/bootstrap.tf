# $ terraform init
# $ terraform -chdir=../tf init
# where -chdir=../tf allows you to declare where the root of your terraform project is located.
# $ terraform apply -var="app_version=1.0.0" -auto-approve
# $ terraform destroy -var="app_version=1.0.0" -auto-approve
# $ terraform version


# To troubleshoot the OCI Terraform Provider:
# https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/terraformtroubleshooting.htm
#############################################################################################################
# IMPORTANT: Resources you provision accrue costs while they are running. It's a good idea, as you learn,   #
#            to always run "terraform destroy" on your project.                                             #
#############################################################################################################
locals {
}
