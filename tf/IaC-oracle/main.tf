#########
# Login #
#########
# (1) https://www.oracle.com/cloud/sign-in.html
# (2) Sign In using a Cloud Account Name
# (3) Cloud Account Name

####################################################################################
# OCI CLI                                                                          #
# https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm#Quickstart #
####################################################################################
# Once you have OCI CLI installed on your machine, execute the following command:
# $ oci setup config
# It will prompt you for all of the information require to generate the config file. You will need
# the following:
# (1) user's OCID
# (2) tenancy's OCID
# (3) the region
#
# When creating the keys, decline creating a passphrase. Once the keys are generated, you'll need
# to associate the public key to the user. From the Oracle Cloud web console, click on "API keys"
# on the left and click on "Add API Key." Upload the public key's pem file.
#
# You can verify that everything is configured properly by running the following command:
# $ oci iam compartment list -c <tenancy-ocid>
#   where <tenancy-ocid> is your tenancy's OCID.
# If there are no errors, you are done.

###########################################
# kubectl                                 #
# https://kubernetes.io/docs/tasks/tools/ #
###########################################
# To create a kubeconfig file for kubectl to access the cluster, execute the following command:
# $ oci ce cluster create-kubeconfig --cluster-id <cluster OCID> --file ~/.kube/<name-of-config-file>
#   --region <region> --token-version 2.0.0 --kube-endpoint PUBLIC_ENDPOINT
# You will need the following:
# (1) cluster's OCID
# (2) name for the config file
# (3) the region
#
# The command will create a kubeconfig file in the ~/.kube directory; the kubeconfig file will
# contain the keys and all of the configuration for kubectl to access the cluster.
#
# Next, set the KUBECONFIG environment variable with the kubeconfig file path.
# $ export KUBECONFIG=~/.kube/<name-of-config-file>
#
# Check if the environment variable was set.
# $ printenv KUBECONFIG
#
# Finally, let's try to list the available nodes in the cluster.
# $ kubectl get nodes
# If the nodes are displayed, you are done.

###################################################################################
# Terraform                                                                       #
# https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli #
###################################################################################
# $ terraform version
#
# $ terraform init
# $ terraform -chdir=../tf init
#   where -chdir=../tf allows you to declare where the root of your terraform project is located.
#
# $ terraform plan
# $ terraform plan -var-file="../tf_secrets.auto.tfvars"
#
# $ terraform apply
# $ terraform apply -auto-approve
# $ terraform apply -var-file="../tf_secrets.auto.tfvars"
# $ terraform apply -var="app_version=1.0.0" -auto-approve
#
###################################################################################################
# IMPORTANT: Resources you provision accrue costs while they are running. It's a good idea, as you#
#            learn, to always run "terraform destroy" on your project.                            #
###################################################################################################
# $ terraform destroy
# $ terraform destroy -var-file="../tf_secrets.auto.tfvars"
# $ terraform destroy -var="app_version=1.0.0" -auto-approve
#
# To troubleshoot the OCI Terraform Provider:
# https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/terraformtroubleshooting.htm

# Virtual Cloud Network (VCN) or Virtual Private Cloud (VPC).
module "vcn" {
  source = "oracle-terraform-modules/vcn/oci"
  version = "3.1.0"
  vcn_name = "fin-vcn"
  vcn_dns_label = "findnslbl"
  compartment_id = oci_identity_compartment.fin-compartment.id
  region = var.region
  # The DNS Domain Name for your virtual cloud network is: <your-dns-label>.oraclevcn.com
  # Alphanumeric string that begins with a letter.
  vcn_cidrs = ["10.0.0.0/16"]
  local_peering_gateways = null
  internet_gateway_route_rules = null
  nat_gateway_route_rules = null
  create_internet_gateway = true
  create_nat_gateway = true
  create_service_gateway = true
}

module "arm64-node-pool" {
  depends_on = [
    module.private-subnet,
    module.public-subnet,
    module.cluster
  ]
  source = "./modules/node"
  name = "arm64-worker-pool"
  tenancy_ocid = var.tenancy_ocid
  compartment_id = oci_identity_compartment.fin-compartment.id
  subnet_id = module.private-subnet.subnet-id
  cluster_id = module.cluster.cluster-id
  cluster_cni_type = module.cluster.cluster-cni-type
  nodes = var.nodes
  memory_per_node = var.memory_per_node
  ocpus_per_node = var.ocpus_per_node
}

module "private-subnet" {
  depends_on = [
    module.vcn
  ]
  source = "./modules/subnet"
  vcn_id = module.vcn.vcn_id
  subnet_display_name = "private-subnet"
  compartment_id = oci_identity_compartment.fin-compartment.id
  cidr_block = "10.0.1.0/24"  # Private subnet's CIDR block.
  # Caution: For the route table id, use module.vcn.nat_route_id.
  # Do not use module.vcn.nat_gateway_id, because it is the OCID for the gateway and not the route
  # table.
  route_table_id = module.vcn.nat_route_id
  # VNICs created in this subnet cannot have public IP addresses.
  prohibit_public_ip_on_vnic = true
  #
  sl_display_name = "private-subnet-security-list"
  sl_egress_security_rules = [{
    stateless = false  # No
    destination = "0.0.0.0/0"  # Allow all traffic to go out anywhere.
    destination_type = "CIDR_BLOCK"
    protocol = "all"  # All protocols
  }]
  # Allow traffic for all ports within the range of the VCN (10.0.0.0/16).
  sl_ingress_security_rules = [{
    stateless = false
    source = "10.0.0.0/16"  # VCN
    source_type = "CIDR_BLOCK"
    protocol = "all"
  }]
}

module "public-subnet" {
  depends_on = [
    module.vcn
  ]
  source = "./modules/subnet"
  vcn_id = module.vcn.vcn_id
  subnet_display_name = "public-subnet"
  compartment_id = oci_identity_compartment.fin-compartment.id
  cidr_block = "10.0.0.0/24"    # Public subnet's CIDR block.
  # Caution: For the route table id, use module.vcn.nat_route_id.
  # Do not use module.vcn.nat_gateway_id, because it is the OCID for the gateway and not the route
  # table.
  route_table_id = module.vcn.nat_route_id
  # VNICs created in this subnet will automatically be assigned public IP addresses unless
  # specified otherwise during instance launch or VNIC creation.
  prohibit_public_ip_on_vnic = false
  #
  sl_display_name = "public-subnet-security-list"
  sl_egress_security_rules = [{
    stateless = false  # No
    destination = "0.0.0.0/0"  # Allow all traffic to go out anywhere.
    destination_type = "CIDR_BLOCK"
    protocol = "all"  # All protocols
  }]
  # Allow traffic for all ports within the range of the VCN (10.0.0.0/16).
  sl_ingress_security_rules = [#{
  #   stateless = false
  #   source = "10.0.0.0/16"
  #   source_type = "CIDR_BLOCK"
  #   protocol = "all"
  # },
  {
    stateless = false
    # Allow VCN traffic to come in as well as traffic from anywhere on port 6443 TCP; we'll use
    # kubectl to communicate with the K8S cluster.
    source = "0.0.0.0/0"
    source_type = "CIDR_BLOCK"
    protocol = "6"
    tcp_options = [{
      min = 6443
      max = 6443
    }]
  },
  {
    stateless = false
    source = "10.0.0.0/16"
    source_type = "CIDR_BLOCK"
    protocol = "6"
    tcp_options = [{
      min = 22
      max = 22
    }]
  }]
}

  #
  #
  # ingress_security_rules {
  #   stateless = false
  #   source = "0.0.0.0/0"
  #   source_type = "CIDR_BLOCK"
  #   # Get protocol numbers from https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml ICMP is 1
  #   protocol = "1"
  #   # For ICMP type and code see: https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml
  #   icmp_options {
  #     type = 3
  #     code = 4
  #   }
  # }
  #
  # ingress_security_rules {
  #   stateless = false
  #   source = "10.0.0.0/16"
  #   source_type = "CIDR_BLOCK"
  #   # Get protocol numbers from https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml ICMP is 1
  #   protocol = "1"
  #   # For ICMP type and code see: https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml
  #   icmp_options {
  #     type = 3
  #   }
  # }

module "cluster" {
  depends_on = [
    module.public-subnet
  ]
  source = "./modules/cluster"
  name = "k8s-cluster"
  type = "BASIC_CLUSTER"
  compartment_id = oci_identity_compartment.fin-compartment.id
  vcn_id = module.vcn.vcn_id
  subnet_ids = [module.public-subnet.subnet-id]
}
