terraform {
  # Terraform version.
  required_version = ">= 1.7.2"
  required_providers {
    null = {
      source = "hashicorp/null"
      version = ">= 3.1.0"
    }
  }
}

provider "null" {
}
