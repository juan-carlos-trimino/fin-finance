
terraform {
  backend "local" {
    path = "../../../tf-states/IaC-app/terraform.tfstate"
  }
}
