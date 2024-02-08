#############################################################################################################
# In Terraform to fetch data, you use a data source. Fetching data from a data source is similar to the GET #
# method in REST APIs.                                                                                      #
#############################################################################################################

# Tenancy is the root or parent to all compartments.
# Use the value of <tenancy-ocid> for the compartment OCID.
# The data source gets a list of availability domains in your entire tenancy. The tenancy is the compartment
# OCID for the root compartment. Providing a specific "<compartment-ocid>" or the "<tenancy-ocid>" outputs
# the same list.
# https://registry.terraform.io/providers/oracle/oci/latest/docs/data-sources/identity_availability_domains
data "oci_identity_availability_domains" "availability_domains" {
  compartment_id = var.tenancy_ocid
}
