
resource "null_resource" "kubeconfig" {
  triggers = {
    always_run = timestamp()
  }
  #
  provisioner "local-exec" {
    command = <<-EOT
      oci ce cluster create-kubeconfig --cluster-id ${module.cluster.cluster-id} \
       --file ${var.kubeconfig_path} --region ${var.region} --token-version 2.0.0 \
       --kube-endpoint PUBLIC_ENDPOINT
    EOT
  }
}
