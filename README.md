# The fin-finance repo
This repo contains the following four (4) directories:
1. **IaC-app**: The *Infrastructure-as-Code (IaC)* application directory contains the Terraform code to deploy and manage the application in the cloud.
2. **IaC-K8s**: This directory contains the Terraform code to set up the Oracle Cloud Infrastructure.
3. **IaC-storage**: The storage directory contains the Terraform code to set up access to the Oracle Cloud Infrastructure, which enables access to the Simple Storage Service (S3).
4. **src**: This directory contains the code for the application.

## Terraform
https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli

To install/upgrade Terraform on Windows Subsystem for Linux (WSL)
```
$ sudo apt update && sudo apt upgrade -y
$ sudo apt install wget unzip
$ wget https://releases.hashicorp.com/terraform/1.10.4/terraform_1.10.4_linux_amd64.zip -O terraform.zip
$ unzip terraform.zip
$ sudo mv terraform /usr/local/bin
$ rm terraform.zip
```
### Useful Commands
To obtain the current version of Terraform and all installed plugins.
```
$ terraform version
```
To initialize a working directory containing Terraform configuration files. ***This is the first command you should run after writing a new Terraform configuration or cloning an existing configuration from version control. It is safe to run this command multiple times.***
```
$ terraform init
$ terraform -chdir=../tf init
  where -chdir=../tf allows you to declare where the root of your terraform project is located.
```

```
$ terraform plan
$ terraform plan -var-file="../tf_secrets.auto.tfvars"

$ terraform apply
$ terraform apply -auto-approve
$ terraform apply -var-file="../tf_secrets.auto.tfvars"
$ terraform apply -var="app_version=1.0.0" -auto-approve

IMPORTANT: Resources you provision accrue costs while they are running. It's a good idea, as you learn, to always run "terraform destroy" on your project.
$ terraform destroy
$ terraform destroy -auto-approve
$ terraform destroy -var-file="../tf_secrets.auto.tfvars"
$ terraform destroy -var="app_version=1.0.0" -auto-approve
```

## IaC-K8s
TBD Enter description

### Login
1. https://www.oracle.com/cloud/sign-in.html
2. Sign In using a Cloud Account Name
3. Cloud Account Name

### OCI CLI
https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/climanualinst.htm#Manual_Installation<br><br>
Manual Installation: Ubuntu<br>
**Step 1: Installing Python**<br>
Before you install the CLI, run the following commands on a new Ubuntu image.<br>
```
~$ sudo apt update
~$ sudo apt install build-essential zlib1g-dev libncurses5-dev libgdbm-dev libnss3-dev libssl-dev libreadline-dev libffi-dev libsqlite3-dev wget libbz2-dev
~$ sudo apt update && sudo apt install python3.12.0 python3.12.0-pip python3.12.0-venv
```

**Step 2: Creating and Configuring a Virtual Environment**<br>
The ***venv*** Python module is a virtual environment builder that lets you create isolated Python environments.<br>
Installing and Activating your Virtual Environment<br>
After Python is installed, set up a virtual environment for your operating system using the following steps.
1. Navigate to the directory in which you would like to create the virtual environment.
   ```
   $ mkdir -p ~/oci/python && cd ~/oci/python
   ```
2. Create the virtual environment using the version of Python installed.
   ```
   ~/oci/python$ python3.12 -m venv oracle-cli
   ```
3. Activate the virtual environment.
   ```
   ~/oci/python$ source oracle-cli/bin/activate
   ```

**Step 3: Installing the Command Line Interface**<br>
To install using PyPI, run the following command:
```
(oracle-cli) ~/oci/python$ pip install oci-cli
```

**Step 4: Setting up the Configuration File**<br>
Before using the CLI, you must create a configuration file that contains the required credentials for working with Oracle Cloud Infrastructure. The default location for the configuration file is ***~/.oci***.<br>
**Use the Setup Dialog**<br>
To have the CLI guide you through the first-time setup process, use the setup config command:
```
(oracle-cli) ~/oci/python$ oci setup config
```
This command prompts you for the information required to create the configuration file and the API public and private keys. The setup dialog uses this information to generate an API key pair and creates the configuration file. After API keys are created, upload the public key using the Console. You will need the following:<br>
1. User's OCID (Profile->My profile)
2. Tenancy's OCID (Profile->Tenancy: \<tenancy-name\>)
3. The region

When creating the keys, decline creating a passphrase. Once the keys are generated, you'll need to associate the public key to the user. From the Oracle Cloud web console, click on ***Profile-> My profile->API keys*** on the left and click on ***Add API Key***. Upload the public key's pem file.

**Step 5: Verify that everything is configured properly**<br>
You can verify that everything is configured properly by running the following command:
```
(oracle-cli) ~/oci/python$ oci iam compartment list -c <tenancy-ocid>
```
where \<tenancy-ocid\> is your tenancy's OCID.

If there are no errors in the JSON reply, the config file was create (by default in ***~/.oci***). At this point, you need to run Terraform to allocate your resources.

**Step 6: Deactivate the virtual environment.**
```
(oracle-cli) ~/oci/python$ deactivate
```

**Step 7: Activate the virtual environment.**
```
~$ source ~/oci/python/oracle-cli/bin/activate
```
