terraform {
  # Terraform version.
  required_version = ">= 1.7.2"
  required_providers {
    oci = {
      source = "oracle/oci"
      version = ">= 5.0.0"
    }
    #
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = ">= 2.8.0"
    }
    #
    null = {
      source = "hashicorp/null"
      version = ">= 3.1.1"
    }
  }
}

# https://docs.oracle.com/en-us/iaas/developer-tutorials/tutorials/tf-provider/01-summary.htm
provider "oci" {  # Oracle Cloud Infrastructure (OCI)
  region = var.region
}

# Configure the K8s Provider.
provider "kubernetes" {
  config_path = "~/.kube/k8s-config"
}

provider "null" {
}
