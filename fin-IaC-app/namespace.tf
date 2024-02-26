# $ kubectl get ns
# $ kubectl describe ns
# Execute the following command to clean up all the pending resources and remove all the K8s
# objects: (Caution: The command in this step removes all the resources.)
# $ kubectl delete all --all -n <namespace>
# Execute the following command to delete the namespace:
# $ kubectl delete namespace <namespace>
resource "kubernetes_namespace" "ns" {
  metadata {
    name = var.app_name
    labels = {
      app = var.app_name
    }
  }
}
