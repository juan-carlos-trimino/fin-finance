
terraform {
  backend "local" {
    path = "../../../tf-states/IaC-app/terraform.tfstate"
    # To prevent S3-native locking (or moving the lock file), set this to false.
    # use_lockfile = false
  }
}
