terraform {
  # Terraform version.
  required_version = ">= 1.7.2"
  required_providers {
    oci = {
      source = "oracle/oci"
      version = ">= 5.0.0"
    }
  }
}

# https://docs.oracle.com/en-us/iaas/developer-tutorials/tutorials/tf-provider/01-summary.htm
provider "oci" {  # Oracle Cloud Infrastructure (OCI)
  tenancy_ocid = var.tenancy_ocid
  user_ocid = var.user_ocid
  private_key_path = var.private_key_path
  fingerprint = var.fingerprint
  # config_file_profile = "DEFAULT"
  region = var.region
}
