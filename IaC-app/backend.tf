
terraform {
  backend "local" {
    path = "$HOME/repos/IaC-app/terraform.tfstate"
  }
}
