# Declare an Oracle Cloud Infrastructure compartment resource and then define the specifics for the
# compartment.
# https://registry.terraform.io/providers/oracle/oci/latest/docs/resources/identity_compartment
resource "oci_identity_compartment" "fin-compartment" {
  description = "Compartment for Terraform resources."
  name = var.compartment_name
  compartment_id = var.tenancy_ocid
  # By default, the Terraform provider does not delete a compartment when using the destroy
  # command.
  # Note: To destroy a compartment, the compartment must also be empty. Use the depends_on argument
  #       to ensure that any hidden dependencies are defined.
  enable_delete = true
}
