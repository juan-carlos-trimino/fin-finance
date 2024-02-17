

data "oci_core_images" "use_image" {
  compartment_id = oci_identity_compartment.fin-compartment.id
  operating_system = "Oracle Linux"
  operating_system_version = "8"
  shape = "VM.Standard.A1.Flex"
  sort_by = "TIMECREATED"
  sort_order = "DESC"
}


