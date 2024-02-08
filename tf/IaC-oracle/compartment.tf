# $ terraform init
# $ terraform apply -auto-approve
# $ terraform plan -var-file="../tf_secrets.auto.tfvars"
# $ terraform apply -var-file="../tf_secrets.auto.tfvars"
# $ terraform destroy -var-file="../tf_secrets.auto.tfvars"

#############################################################################################################
# In Terraform, resources are objects such as virtual cloud networks or compute instances. You can create,  #
# update, and delete them with Terraform.                                                                   #
#############################################################################################################

# Declare an Oracle Cloud Infrastructure compartment resource and then define the specifics for the
# compartment.
resource "oci_identity_compartment" "tf-compartment" {
  description = "Compartment for Terraform resources."
  name = "fin-compartment"
  compartment_id = var.tenancy_ocid
  # By default, the Terraform provider does not delete a compartment when using the destroy command.
  # Note: To destroy a compartment, the compartment must also be empty. Use the depends_on argument to ensure
  #       that any hidden dependencies are defined.
  enable_delete = true
}
