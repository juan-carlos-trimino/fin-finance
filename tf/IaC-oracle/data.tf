###################################################################################################
# In Terraform to fetch data, you use a data source. Fetching data from a data source is similar  #
# to the GET method in REST APIs.                                                                 #
###################################################################################################
data "oci_containerengine_cluster_kube_config" "kubeconfig" {
  cluster_id = module.cluster.cluster-id
  endpoint = "PUBLIC_ENDPOINT"  # LEGACY_KUBERNETES,PUBLIC_ENDPOINT,PRIVATE_ENDPOINT,VCN_HOSTNAME.
  token_version = "2.0.0"
}
