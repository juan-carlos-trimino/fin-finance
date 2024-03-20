############################
# DECLARE GLOBAL VARIABLES #
############################
variable tenancy_ocid {
  description = "Organize and control access to your resources."
  type = string
  sensitive = true
}

variable compartment_id {
  description = "Organize and control access to your resources."
  type = string
  sensitive = true
}

# From the Profile menu, go to "My profile" and copy OCID.
variable user_ocid {
  type = string
  sensitive = true
}

# From the Profile menu (user avatar), go to "My profile" and click "API Keys" (on the left side).
# Copy the fingerprint associated with the RSA public key. The format is:
# xx:xx:xx...xx
variable fingerprint {
  type = string
  sensitive = true
}

variable region {
  description = "An OCI region."
  type = string
  sensitive = true
}

variable bucket_names {
  type = list(string)
  sensitive = true
}

variable access_types {
  type = list(string)
  sensitive = true
}

variable auto_tiering {
  type = list(string)
  sensitive = true
}

# variable retention_rules {
#   type = list(list(object({
#     display_name = string
#     duration = list(object({
#       string
#     ))}
#     prefix = string
#   })))
# }

variable storage_tiers {
  type = list(string)
  sensitive = true
}

# Collect the following information from your environment.
# Path to the RSA private key file to use for authentication against OCI. An API key can be created
# in the UI under Profile -> My Profile -> API keys.
# Example for Oracle Linux: /home/opc/.oci/<your-rsa-key-name>.pem
# See /IaC-K8s/oracle/variables.tf
variable private_key_path {
  description = "Path to the RSA private key."
  type = string
  sensitive = true
}


