# Login:
# (1) https://www.oracle.com/cloud/sign-in.html
# (2) Sign In using a Cloud Account Name
# (3) astros69

#First of all, when you have OCI CLI installed on your machine, execute the following command:
#$ oci setup config
#This will prompt you all the data it’ll need to generate the proper configuration for you.


# After you’re done with the key generation, there’s one more thing we need to do. Associating the public key that was generated during setup with the user. Go back to your user in the Oracle Cloud web console, click on API keys on the left and click on Add API Key. Upload your public key’s pem file and you’re done.

# You can verify that everything is configured properly by running the following command:
# $ oci iam compartment list -c <tenancy-ocid>
# Where <tenancy-ocid> is your tenancy’s OCID. If you don’t get some authorization error but your compartments’ data, you’re good to go.

# And then we’ll create a kubeconfig for kubectl to access the cluster. Let's execute the following command:
# 


# $ terraform init
# $ terraform apply -auto-approve
# $ terraform plan
# $ terraform apply
# $ terraform destroy

# $ terraform plan -var-file="../tf_secrets.auto.tfvars"
# $ terraform apply -var-file="../tf_secrets.auto.tfvars"
# $ terraform destroy -var-file="../tf_secrets.auto.tfvars"


# $ terraform init
# $ terraform -chdir=../tf init
# where -chdir=../tf allows you to declare where the root of your terraform project is located.
# $ terraform apply -var="app_version=1.0.0" -auto-approve
# $ terraform destroy -var="app_version=1.0.0" -auto-approve
# $ terraform version


# To troubleshoot the OCI Terraform Provider:
# https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/terraformtroubleshooting.htm
#############################################################################################################
# IMPORTANT: Resources you provision accrue costs while they are running. It's a good idea, as you learn,   #
#            to always run "terraform destroy" on your project.                                             #
#############################################################################################################


#############################################################################################################
# Accessing the K8s cluster                                                                                 #
# 1. Let's create a kubeconfig file for kubectl. The command below will create a kubeconfig file in the     #
#    ~/.kube directory; you will need the cluster OCID, the name of the config file, and the region. After  #
#    executing the command, the kubeconfig file will contain the keys and all the configuration for kubectl #
#    to access the cluster.                                                                                 #
#    $ oci ce cluster create-kubeconfig --cluster-id <cluster OCID> --file ~/.kube/<name-of-config-file>    #
#      --region <region> --token-version 2.0.0 --kube-endpoint PUBLIC_ENDPOINT                              #
# 2. Set the KUBECONFIG environment variable with the kubeconfig file path to connect to the cluster.       #
#    $ export KUBECONFIG=~/.kube/<name-of-config-file>                                                      #
# 3. Check if the environment variable was set.                                                             #
#    $ printenv KUBECONFIG                                                                                  #
# 4. Finally, let's try to list the available nodes in the cluster.                                         #
#    $ kubectl get nodes                                                                                    #
#############################################################################################################


# Virtual Cloud Network (VCN) or Virtual Private Cloud (VPC).
module "vcn" {
  source = "oracle-terraform-modules/vcn/oci"
  version = "3.1.0"
  vcn_name = "fin-vcn"
  # The DNS Domain Name for your virtual cloud network is: <your-dns-label>.oraclevcn.com
  # Alphanumeric string that begins with a letter.
  vcn_dns_label = "findnslbl"
  vcn_cidrs = ["10.0.0.0/16"]
  compartment_id = oci_identity_compartment.fin-compartment.id
  region = var.region
  internet_gateway_route_rules = null
  local_peering_gateways = null
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
  subnet_display_name = "private-subnet"
  cidr_block = "10.0.1.0/24"
  compartment_id = oci_identity_compartment.fin-compartment.id
  vcn_id = module.vcn.vcn_id
  # Caution: For the route table id, use module.vcn.nat_route_id.
  # Do not use module.vcn.nat_gateway_id, because it is the OCID for the gateway and not the route
  # table.
  route_table_id = module.vcn.nat_route_id
  #
  sl_display_name = "private-subnet-security-list"
  sl_egress_security_rules = [{
    stateless = false  # No
    destination = "0.0.0.0/0"
    destination_type = "CIDR_BLOCK"
    protocol = "all"  # All protocols
  }]
  sl_ingress_security_rules = [{
    stateless = false
    source = "10.0.0.0/16"
    source_type = "CIDR_BLOCK"
    protocol = "all"
  }]
}

module "public-subnet" {
  depends_on = [
    module.vcn
  ]
  source = "./modules/subnet"
  subnet_display_name = "public-subnet"
  cidr_block = "10.0.0.0/24"
  compartment_id = oci_identity_compartment.fin-compartment.id
  vcn_id = module.vcn.vcn_id
  # Caution: For the route table id, use module.vcn.nat_route_id.
  # Do not use module.vcn.nat_gateway_id, because it is the OCID for the gateway and not the route
  # table.
  route_table_id = module.vcn.nat_route_id
  #
  sl_display_name = "public-subnet-security-list"
  sl_egress_security_rules = [{
    stateless = false  # No
    destination = "0.0.0.0/0"
    destination_type = "CIDR_BLOCK"
    protocol = "all"  # All protocols
  }]
  #
  sl_ingress_security_rules = [{
    stateless = false
    source = "10.0.0.0/16"
    source_type = "CIDR_BLOCK"
    protocol = "all"
  },
  {
    stateless = false
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
    module.vcn
  ]
  source = "./modules/cluster"
  name = "k8s-cluster"
  type = "BASIC_CLUSTER"
  compartment_id = oci_identity_compartment.fin-compartment.id
  vcn_id = module.vcn.vcn_id
  subnet_id = module.public-subnet.subnet-id
}
