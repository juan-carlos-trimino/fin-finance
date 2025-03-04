resource "null_resource" "kubeconfig" {
  # The trigger value forces the provisioner to run during each plan/apply.
  triggers = {
    always_run = timestamp()
  }
  provisioner "local-exec" {
    command = <<-EOT
      oci ce cluster create-kubeconfig --cluster-id ${module.cluster.cluster-id} \
       --file ${var.kubeconfig_path} --region ${var.region} --token-version 2.0.0 \
       --kube-endpoint PUBLIC_ENDPOINT
    EOT
  }
}

# Next, set the KUBECONFIG environment variable with the kubeconfig file path.
# $ export KUBECONFIG=~/.kube/<name-of-config-file>
#
# Check if the environment variable was set.
# $ printenv KUBECONFIG
# resource "null_resource" "kubeconfig-ev" {
#   triggers = {
#     always_run = timestamp()
#   }
#   #
#   provisioner "local-exec" {
#     interpreter = ["/bin/bash", "-c"]
#     command = <<-EOT
#       echo 'KUBECONFIG="${var.kubeconfig_path}}"' >> ~/.bashrc
#       source ~/.bashrc
#     EOT
#   }
# }
