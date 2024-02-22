# $ kubectl get ns
# $ kubectl describe ns
resource "kubernetes_namespace" "ns" {
  metadata {
    name = var.app_name
    labels = {
      app = var.app_name
    }
  }
}
