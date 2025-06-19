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
    #
    helm = {
      source = "hashicorp/helm"
      version = ">= 2.17.0"
    }
    #
    digitalocean = {
      # Using an environment variable to set the DIGITALOCEAN_TOKEN.
      source = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

# Load and connect to Helm.
provider "helm" {
  kubernetes {
    config_path = "~/.kube/k8s-config"
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

provider "digitalocean" {
}
