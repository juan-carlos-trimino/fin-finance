# The fin-finance repo
This repo contains the following four (4) directories:
1. **IaC-app**: The *Infrastructure-as-Code (IaC)* application directory contains the Terraform code to deploy and manage the application in the cloud.
2. **IaC-K8s**: This directory contains the Terraform code to set up the Oracle Cloud Infrastructure.
3. **IaC-storage**: The storage directory contains the Terraform code to set up access to the Oracle Cloud Infrastructure, which enables access to the Simple Storage Service (S3).
4. **src**: This directory contains the code for the application.

## Terraform
https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli

To troubleshoot the ***OCI Terraform Provider***:<br>
https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/terraformtroubleshooting.htm

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
```
where -chdir=../tf allows you to declare where the root of your terraform project is located.

This command creates an execution plan, which lets you preview the changes that Terraform plans to make to your infrastructure.
```
$ terraform plan
$ terraform plan -var-file="../tf_secrets.auto.tfvars"
```
where -var-file="../tf_secrets.auto.tfvars" sets values for potentially many input variables declared in the root module of the configuration, using definitions from a ***tfvars*** file. Use this option multiple times to include values from more than one file.

This command executes the actions proposed in a Terraform plan.
```
$ terraform apply
$ terraform apply -auto-approve
$ terraform apply -var-file="../tf_secrets.auto.tfvars"
$ terraform apply -var="app_version=1.0.0" -auto-approve
```
where -auto-approve skips interactive approval of the plan before applying. Terraform ignores this option when you pass a previously-saved plan file. This is because Terraform interprets the act of passing the plan file as the approval.<br>
and<br>
-var sets a value for a single input variable declared in the root module of the configuration. Use this option multiple times to set more than one variable.

***IMPORTANT***: Resources you provision accrue costs while they are running. It's a good idea, as you learn, to always run **terraform destroy** on your project.<br>
To deprovision all objects managed by a Terraform configuration.
```
$ terraform destroy
$ terraform destroy -auto-approve
$ terraform destroy -var-file="../tf_secrets.auto.tfvars"
$ terraform destroy -var="app_version=1.0.0" -auto-approve
```

Once Terraform finish setting up your resources, you need to set up ***kubectl*** to access the cluster. See ***kubectl***.

Finally, let's try to list the available nodes in the cluster.
```
$ kubectl get nodes
```
If the nodes are displayed, you are done.

## kubectl
https://kubernetes.io/docs/tasks/tools/

---
**Note**

A file that is used to configure access to a cluster is usually referred to as a ***kubeconfig file***. This is a conventional way of referring to a configuration file, often shortened to config file. It does not imply that a file named kubeconfig exists.

---

You will need to create a kubeconfig file with authentication and configuration details, which will allow kubectl to communicate with your cluster. To create the kubeconfig file, you execute the command below, which requires the following information:<br>
**(1)** Cluster's OCID (Navigation menu->Developer Services->Kubernetes Clusters (OKE) [Under Containers & Artifacts]->Select the compartment that contains the cluster[Compartment]-> On the Clusters page, click the name of the cluster)<br>
**(2)** Name for the config file<br>
**(3)** The region
```
$ oci ce cluster create-kubeconfig --cluster-id <cluster OCID> --file ~/.kube/<name-of-config-file> --region <region> --token-version 2.0.0 --kube-endpoint PUBLIC_ENDPOINT
```
The command will create a kubeconfig file in the ***~/.kube*** directory; the kubeconfig file will contain the keys and all of the configuration for kubectl to access the cluster. See ***IaC-K8s/oracle/data.tf*** for appropriate values to the parameters ***--token-version*** and ***--kube-endpoint***.

---
**Note**

Setting the permissions of your ***~/.kube/\<name-of-config-file\>*** file to ***600*** ensures that only the owner (you) can read and write to it, enhancing security by limiting access to the Kubernetes configuration file.

```
$ chmod 600  ~/.kube/<name-of-config-file>
```
---

By default, kubectl looks for a file named ***config*** in the ***$HOME/.kube (~/.kube)*** directory; hence, if the ***KUBECONFIG*** environment variable is not set, kubectl uses the default values ***~/.kube/config***. You can specify other kubeconfig files by setting the ***KUBECONFIG*** environment variable or by setting the ***--kubeconfig*** flag.

To export the KUBECONFIG environment variable ***only*** for the current shell and its children processes, you use the ***export*** command:
```
export KUBECONFIG=<name-of-config-file>
```
To reiterate, when an environment variable is set from the shell using the export command, its existence ends when the current session ends.

To set the KUBECONFIG environment variable as a ***user-specific environment variable***, add the ***export*** command to ***~/.bashrc (bash), ~/.kshrc (ksh), or ~/.zshrc (zsh)***, depending on which shell you are using. By modifying the shell-specific configuration file, the environment variable will persist across sessions and system restarts. Below the bash shell is used:
```
$ echo 'export KUBECONFIG=<name-of-config-file>' >> ~/.bashrc
```
Next, reload the file to apply the changes:
```
$ source ~/.bashrc
```

To view all environment variables, use the ***printenv*** command. Since there are many variables on the list, use the ***less*** command to control the view:
```
$ printenv | less
```
The output shows the first page of the list and allows you to move forward by pressing ***Space*** to see the next page or ***Enter*** to display the next line. Exit the view with ***q***.

To view a specific environment variable, use the ***set*** command:
```
$ set | grep KUBECONFIG
```
or
```
$ echo $KUBECONFIG
```
or
```
$ printenv KUBECONFIG
```

## Kubenertes (K8s)
### Useful Commands
Display the Kubernetes version running on the client and server.
```
$ kubectl version
```

Display the configuration of the cluster.
```
$ kubectl config view
```

Display a list of users.
```
$ kubectl config view -o jsonpath='{.users[*].name}'
```

Retrieve cluster details.
```
$ kubectl cluster-info
```

Confirm what platform is running on the cluster.
```
$ kubectl describe node | grep "kubernetes.io/arch"
```

To retrieve nodes information.
```
$ kubectl get nodes
```

Retrieve built-in resource types (pods, services, daemon sets, deployments, replica sets, jobs, cronjobs, and stateful sets) under a specific namespace.
```
$ kubectl get all -n <name-space>
```

Retrieve pods under a specific namespace.
```
$ kubectl get pods -n <name-space>
```

Delete a pod under a specific namespace.
```
$ kubectl delete pod <pod-name> -n <name-space>
```



[text](https://spacelift.io/blog/kubectl-delete-pod)
# $ kubectl describe -n <name-space> pod <pod-name>
# $ kubectl get -n <name-space> -o jsonpath='{.spec.containers[*].ports[*].containerPort}' pod <pod-name>

Display what node a pod is scheduled.
```
$ kubectl get po -o wide -n <name-space>
```

Retrieve detail pod information.
```
$ kubectl describe pod <pod-name> -n <name-space>
```

## IaC-K8s
IaC-K8s contains the Terraform code for provisioning (i.e., creating, preparing, and activating the underlying infrastructure of a cloud environment) the Oracle Cloud Infrastructure (OCI), which is an Infrastructure as a Service (IaaS) and Platform as a Service (PaaS) offering. The OCI is a set of complementary cloud services that enable you to build and run a range of applications and services in a highly available hosted environment. OCI provides high-performance compute capabilities (as physical hardware instances) and storage capacity in a flexible overlay virtual network that is securely accessible from your on-premises network.

For more OCI and Terraform documentation, please see [Using Terraform and Oracle Cloud Infrastructure](https://docs.oracle.com/en-us/iaas/Content/dev/terraform/tutorials.htm).

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
Use the Setup Dialog<br>
To have the CLI guide you through the first-time setup process, use the setup config command:
```
(oracle-cli) ~/oci/python$ oci setup config
```
This command prompts you for the information required to create the configuration file and the API public and private keys. The setup dialog uses this information to generate an API key pair and creates the configuration file. After API keys are created, upload the public key using the Console. You will need the following:<br>
1. User's OCID (Profile->User settings or My profile)
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

## src
### Debugging Go
#### Delve
To install the debugger in VS Code:<br>
**(1)** Open the Command Palette (***Ctrl + Shift + P***).<br>
**(2)** Find ***Go: Install/Update Tools*** and select ***dlv***.

The settings for the debugger can be stored in the ***.code-workspace*** file or the ***.vscode/launch.json*** directory. For this project, the settings are stored in the ***.code-workspace*** file under the ***launch*** section.
