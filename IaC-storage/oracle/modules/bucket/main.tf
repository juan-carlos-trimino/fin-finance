/***
-------------------------------------------------------
A Terraform reusable module for deploying microservices
-------------------------------------------------------
Define input variables to the module.
***/
variable namespace {
  type = string
}
variable compartment_id {
  type = string
  sensitive = true
}
variable bucket_name {
  type = string
  sensitive = true
}
variable access_type {
  type = string
  sensitive = true
}
variable auto_tiering {
  type = string
  sensitive = true
}
variable storage_tier {
  type = string
  sensitive = true
}

resource "oci_objectstorage_bucket" "bucket" {
  compartment_id = var.compartment_id
  namespace = var.namespace
  name = var.bucket_name
  access_type = var.access_type
  auto_tiering = var.auto_tiering
  storage_tier = var.storage_tier
}
