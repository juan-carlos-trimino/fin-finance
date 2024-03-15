###########################################
# kubectl                                 #
# https://kubernetes.io/docs/tasks/tools/ #
###########################################
# To create a kubeconfig file for kubectl to access the cluster, execute the following command:
# $ oci ce cluster create-kubeconfig --cluster-id <cluster OCID> --file ~/.kube/<name-of-config-file>
#   --region <region> --token-version 2.0.0 --kube-endpoint PUBLIC_ENDPOINT
# You will need the following:
# (1) cluster's OCID (Navigation menu->Developer Services->Kubernetes Clusters (OKE) [Under
#     Containers & Artifacts]->Select the compartment that contains the cluster[Compartment]->
#     On the Clusters page, click the name of the cluster)
# (2) name for the config file
# (3) the region
#
# The command will create a kubeconfig file in the ~/.kube directory; the kubeconfig file will
# contain the keys and all of the configuration for kubectl to access the cluster.
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
