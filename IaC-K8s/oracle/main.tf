###################################################################################################
# Login                                                                                           #
# (1) https://www.oracle.com/cloud/sign-in.html                                                   #
# (2) Sign In using a Cloud Account Name                                                          #
# (3) Cloud Account Name                                                                          #
###################################################################################################
# OCI CLI                                                                                         #
# https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/climanualinst.htm#Manual_Installation    #
# Manual Installation: Ubuntu                                                                     #
# Step 1: Installing Python                                                                       #
#  Before you install the CLI, run the following commands on a new Ubuntu image.                  #
#  ~$ sudo apt update                                                                             #
#  ~$ sudo apt install build-essential zlib1g-dev libncurses5-dev libgdbm-dev libnss3-dev \       #
#     libssl-dev libreadline-dev libffi-dev libsqlite3-dev wget libbz2-dev                        #
#  ~$ sudo apt update && sudo apt install python3.12.0 python3.12.0-pip python3.12.0-venv         #
# Step 2: Creating and Configuring a Virtual Environment                                          #
#  The venv Python module is a virtual environment builder that lets you create isolated Python   #
#  environments.                                                                                  #
#  Installing and Activating your Virtual Environment                                             #
#  After Python is installed, set up a virtual environment for your operating system using the    #
#  following steps.                                                                               #
#  1. Navigate to the directory in which you would like to create the virtual environment.        #
#     $ mkdir -p ~/oci/python && cd ~/oci/python                                                  #
#  2. Create the virtual environment using the version of Python installed.                       #
#     ~/oci/python$ python3.12 -m venv oracle-cli                                                 #
#  3. Activate the virtual environment.                                                           #
#     ~/oci/python$ source oracle-cli/bin/activate                                                #
# Step 3: Installing the Command Line Interface                                                   #
#  To install using PyPI, run the following command:                                              #
#  (oracle-cli) ~/oci/python$ pip install oci-cli                                                 #
# Step 4: Setting up the Configuration File                                                       #
#  Before using the CLI, you must create a configuration file that contains the required          #
#  credentials for working with Oracle Cloud Infrastructure. The default location for the         #
#  configuration file is ~/.oci.                                                                  #
#  Use the Setup Dialog                                                                           #
#  To have the CLI guide you through the first-time setup process, use the setup config command:  #
#  (oracle-cli) ~/oci/python$ oci setup config                                                    #
#  This command prompts you for the information required to create the configuration file and the #
#  API public and private keys. The setup dialog uses this information to generate an API key pair#
#  and creates the configuration file. After API keys are created, upload the public key using the#
#  Console.                                                                                       #
#  You will need the following:                                                                   #
#  (1) user's OCID (Profile->My profile)                                                          #
#  (2) tenancy's OCID (Profile->Tenancy: <tenancy-name>)                                          #
#  (3) the region                                                                                 #
#  When creating the keys, decline creating a passphrase. Once the keys are generated, you'll need#
#  to associate the public key to the user. From the Oracle Cloud web console, click on "Profile->#
#  My profile->API keys" on the left and click on "Add API Key." Upload the public key's pem file.#
# Step 5: Verify that everything is configured properly                                           #
#  You can verify that everything is configured properly by running the following command:        #
#  (oracle-cli) ~/oci/python$ oci iam compartment list -c <tenancy-ocid>                          #
#  where <tenancy-ocid> is your tenancy's OCID.                                                   #
#  If there are no errors in the JSON reply, the config file was create (by default in ~/.oci). At#
#  this point, you need to run Terraform to allocate your resources.                              #
# Step 6: Deactivate the virtual environment.                                                     #
#  (oracle-cli) ~/oci/python$ deactivate                                                          #
# Step 7: Activate the virtual environment.                                                       #
#  ~$ source ~/oci/python/oracle-cli/bin/activate                                                 #
###################################################################################################
# Terraform                                                                                       #
# https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli                 #
#                                                                                                 #
# To install/upgrade Terraform on WSL                                                             #
# $ sudo apt update && sudo apt upgrade -y                                                        #
# $ sudo apt install wget unzip                                                                   #
# $ wget https://releases.hashicorp.com/terraform/1.10.4/terraform_1.10.4_linux_amd64.zip \       #
#   -O terraform.zip                                                                              #
# $ unzip terraform.zip                                                                           #
# $ sudo mv terraform /usr/local/bin                                                              #
# $ rm terraform.zip                                                                              #
###################################################################################################
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
# $ terraform destroy -auto-approve
# $ terraform destroy -var-file="../tf_secrets.auto.tfvars"
# $ terraform destroy -var="app_version=1.0.0" -auto-approve
#
# To troubleshoot the OCI Terraform Provider:
# https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/terraformtroubleshooting.htm
#
# Once Terraform finish setting up your resources, you need to set up kubectl to access the cluster.
# See null-resources.tf.
#
# Finally, let's try to list the available nodes in the cluster.
# $ kubectl get nodes
# If the nodes are displayed, you are done.

# Virtual Cloud Network (VCN) or Virtual Private Cloud (VPC).
module "vcn" {
  source = "oracle-terraform-modules/vcn/oci"
  version = "3.1.0"
  vcn_name = "vcn-fiv"
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
  },
  {
    stateless = false
    source = "10.0.0.0/24"
    source_type = "CIDR_BLOCK"
    protocol = "6"
    tcp_options = [{  # Port to bind the health check server (kube-proxy).
      min = 10256
      max = 10256
    }]
  },
  {
    stateless = false
    source = "10.0.0.0/24"
    source_type = "CIDR_BLOCK"
    protocol = "6"
    tcp_options = [{  # NodePort
      min = var.nlb_node_port
      max = var.nlb_node_port
    }]
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
  route_table_id = module.vcn.ig_route_id
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
  },
  {
    stateless = false
    destination = "10.0.1.0/24"
    destination_type = "CIDR_BLOCK"
    protocol = "6"
    tcp_options = [{  # Port to bind the health check server (kube-proxy).
      min = 10256
      max = 10256
    }]
  },
  {
    stateless = false
    destination = "10.0.1.0/24"
    destination_type = "CIDR_BLOCK"
    protocol = "6"
    tcp_options = [{
      min = var.nlb_node_port
      max = var.nlb_node_port
    }]
  }]
  # Allow traffic for all ports within the range of the VCN (10.0.0.0/16).
  sl_ingress_security_rules = [{
    stateless = false
    # Allow VCN traffic to come in as well as traffic from anywhere on port 6443 TCP (for kubectl
    # to communicate with the K8S cluster).
    source = "0.0.0.0/0"
    source_type = "CIDR_BLOCK"
    protocol = "6"
    tcp_options = [{  # kubectl.
      min = 6443
      max = 6443
    }]
  },
  # Allow the load balancer to communicate with the public subnet.
  {
    stateless = false
    source = "0.0.0.0/0"
    source_type = "CIDR_BLOCK"
    protocol = "6"
    tcp_options = [{
      max = 80
      min = 80
    }]
  },
  {
    stateless = false
    source = "0.0.0.0/0"
    source_type = "CIDR_BLOCK"
    protocol = "6"
    tcp_options = [{
      max = 443
      min = 443
    }]
  },
  {
    stateless = false
    source = "10.0.0.0/16"
    source_type = "CIDR_BLOCK"
    protocol = "all"
  }]
}

module "cluster" {
  depends_on = [
    module.public-subnet
  ]
  source = "./modules/cluster"
  name = "k8s-cluster"
  type = "BASIC_CLUSTER"
  compartment_id = oci_identity_compartment.fin-compartment.id
  vcn_id = module.vcn.vcn_id
  # k8s_version = "v1.29.1"
  k8s_version = "v1.30.1"
  subnet_ids = [
    module.public-subnet.subnet-id
  ]
}

module "arm64-node-pool" {
  depends_on = [
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
  k8s_version = var.k8s_version
  public_key = local.public_key
  memory_per_node = var.memory_per_node
  ocpus_per_node = var.ocpus_per_node
}

module "node-port-nlb" {
  depends_on = [
    module.arm64-node-pool
  ]
  count = var.network_load_balancer ? 1 : 0
  source = "./modules/network-load-balancer"
  nlb_node_port = var.nlb_node_port
  compartment_id = oci_identity_compartment.fin-compartment.id
  subnet_id = module.public-subnet.subnet-id
  target_ids = local.active_nodes
  nodes = var.nodes
}

module "load-balancer" {
  depends_on = [
    module.arm64-node-pool
  ]
  count = var.load_balancer ? 1 : 0
  source = "./modules/load-balancer"
  compartment_id = oci_identity_compartment.fin-compartment.id
  shape = "flexible"
  is_private = false
  public_subnet_id = module.public-subnet.subnet-id
  vcn_id = module.vcn.vcn_id
  target_ids = local.active_nodes
  shape_details_maximum_bandwidth_in_mbps = 10

  # node_pool_id = module.arm64-node-pool.node-pool-id
  # instance_pool_id   = module.instance-pool.instance_pool_id

  # region             = var.region
  # node_pool_size = module.arm64-node-pool.instance_pool_size
  # private_subnet_id  = module.private-vcn.private_subnet_id
}

module "igw" {
  depends_on = [
    module.arm64-node-pool
  ]
  count = var.igw ? 1 : 0
  source = "./modules/igw"
  compartment_id = oci_identity_compartment.fin-compartment.id
  name = "igw"
  enabled = true
  vcn_id = module.vcn.vcn_id
}



/***
module "nat" {
  count = var.nat ? 1 : 0
  source = "../modules/nat"
  region = var.region
  # compartment_ocid = xxxxxxxxxxxxxxxxvar.compartment_ocid
  # availability_domain = xxxxxxxxxxxxxxvar.availability_domain
  # vcn_id = module.private-vcn.vcn_id
  # private_subnet_id = module.private-vcn.private_subnet_id
  # public_subnet_id = module.private-vcn.public_subnet_id
  # environment = var.environment
}
***/
