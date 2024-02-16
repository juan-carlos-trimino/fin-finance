# $ terraform init
# $ terraform apply -auto-approve
# $ terraform plan -var-file="../tf_secrets.auto.tfvars"
# $ terraform apply -var-file="../tf_secrets.auto.tfvars"
# $ terraform destroy -var-file="../tf_secrets.auto.tfvars"

#############################################################################################################
# Accessing the K8s cluster                                                                                 #
# 1. Let's create a kubeconfig file for kubectl. The command below will create a kubeconfig file in the     #
#    ~/.kube directory; you will need the cluster OCID, the name of the config file, and the region. After  #
#    executing the command, the kubeconfig file will contain the keys and all the configuration for kubectl #
#    to access the cluster.                                                                                 #
#    $ oci ce cluster create-kubeconfig --cluster-id <cluster OCID> --file ~/.kube/<name-of-config-file>    #
#      --region <region> --token-version 2.0.0 --kube-endpoint PUBLIC_ENDPOINT                              #
# 2. Set the KUBECONFIG environment variable with the kubeconfig file path to connect to the cluster.       #
#    $ export KUBECONFIG=~/.kube/<name-of-config-file>                                                      #
# 3. Check if the environment variable was set.                                                             #
#    $ printenv KUBECONFIG                                                                                  #
# 4. Finally, let's try to list the available nodes in the cluster.                                         #
#    $ kubectl get nodes                                                                                    #
#############################################################################################################

#############################################################################################################
# In Terraform, resources are objects such as virtual cloud networks or compute instances. You can create,  #
# update, and delete them with Terraform.                                                                   #
#############################################################################################################

# Declare an Oracle Cloud Infrastructure compartment resource and then define the specifics for the
# compartment.
# https://registry.terraform.io/providers/oracle/oci/latest/docs/resources/identity_compartment
resource "oci_identity_compartment" "fin-compartment" {
  description = "Compartment for Terraform resources."
  name = "fin-compartment"
  compartment_id = var.tenancy_ocid
  # By default, the Terraform provider does not delete a compartment when using the destroy command.
  # Note: To destroy a compartment, the compartment must also be empty. Use the depends_on argument to ensure
  #       that any hidden dependencies are defined.
  enable_delete = true
}



###########################
# Outputs for compartment #
###########################
# output "compartment-name" {
#   value = oci_identity_compartment.fin-compartment.name
# }

output "compartment-id" {
  value = oci_identity_compartment.fin-compartment.id
}
