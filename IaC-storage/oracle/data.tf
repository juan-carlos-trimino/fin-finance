###################################################################################################
# In Terraform to fetch data, you use a data source. Fetching data from a data source is similar  #
# to the GET method in REST APIs.                                                                 #
###################################################################################################
# Each tenancy is assigned one unique system-generated and immutable namespace name.
# To view the Object Storage namespace string, do the following:
# Open the Profile menu and click Tenancy: <tenancy_name>. The namespace string is listed under
# Object Storage Settings.
# https://docs.oracle.com/en-us/iaas/Content/Object/Tasks/understandingnamespaces.htm
data "oci_objectstorage_namespace" "ns" {
  # compartment_id = var.tenancy_ocid
  compartment_id = var.compartment_id
}
