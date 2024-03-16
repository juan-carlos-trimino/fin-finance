




resource "oci_core_instance" "nat" {
  display_name = "NAT"
  agent_config {
    is_management_disabled = "false"
    is_monitoring_disabled = "false"
    plugins_config {
      name = "Vulnerability Scanning"
      desired_state = "DISABLED"
    }
    plugins_config {
      name = "Compute Instance Monitoring"
      desired_state = "ENABLED"
    }
    plugins_config {
      name = "Bastion"
      desired_state = "DISABLED"
    }
  }
  #
  availability_config {
    recovery_action = "RESTORE_INSTANCE"
  }
  availability_domain = var.availability_domain
  compartment_id = var.compartment_id
  fault_domain = var.default_fault_domain
  create_vnic_details {
    assign_private_dns_record = true
    assign_public_ip = true
    subnet_id = var.public_subnet_id
    skip_source_dest_check = true
  }
  #
  instance_options {
    are_legacy_imds_endpoints_disabled = false
  }
  is_pv_encryption_in_transit_enabled = true
  metadata = {
    "ssh_authorized_keys" = file(var.PATH_TO_PUBLIC_KEY)
    "user_data" = data.cloudinit_config.nat_instance_init.rendered
  }
  shape = "VM.Standard.A1.Flex"
  shape_config {
    memory_in_gbs = "6"
    ocpus = "1"
  }
  source_details {
    source_id = var.os_image_id
    source_type = "image"
  }
  freeform_tags = local.tags
}

output "nat_id" {
  value = oci_core_instance.nat_instance.id
}

output "nat_public_ip" {
  value = oci_core_instance.nat_instance.public_ip
}
