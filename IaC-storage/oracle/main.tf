#############
# Terraform #
#############
# $ terraform init
# $ terraform apply -auto-approve
# $ terraform destroy -auto-approve

module "create_buckets" {
  source = "./modules/bucket"
  count = length(var.bucket_names)
  compartment_id = var.compartment_id
  namespace = data.oci_objectstorage_namespace.ns.namespace
  bucket_name = var.bucket_names[count.index]
  access_type = var.access_types[count.index]
  auto_tiering = var.auto_tiering[count.index]
  storage_tier = var.storage_tiers[count.index]
}
