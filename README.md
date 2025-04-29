# The fin-finance repo
This repo contains the following four (4) directories:
1. **IaC-app**: The `Infrastructure-as-Code (IaC)` application directory contains the `Terraform` code to deploy and manage the application in the cloud.
2. **IaC-K8s**: This directory contains the `Terraform` code to set up the `Oracle Cloud Infrastructure`.
3. **IaC-storage**: The storage directory contains the `Terraform` code to set up access to the `Oracle Cloud Infrastructure`, which enables access to the `Simple Storage Service (S3)`.
4. **src**: This directory contains the code for the application.

```
fin-finance
 ├ .github
 | ├ workflows
 | | └ github-actions-demo.yml
 | └ .gitkeep
 ├ .vscode
 | ├ .gitignore
 | ├ LICENSE
 | ├ README.md
 | └ settings.json
 ├ IaC-K8s
 | ├ oracle
 | | ├ modules
 | | | ├ cluster
 | | | | └ main.tf
 | | | ├ igw
 | | | | └ main.tf
 | | | ├ load-balancer
 | | | | └ main.tf
 | | | ├ nat
 | | | | └ main.tf
 | | | ├ network-load-balancer
 | | | | └ main.tf
 | | | ├ node
 | | | | └ main.tf
 | | | └ subnet
 | | |   └ main.tf
 | | ├ compartment.tf
 | | ├ data.tf
 | | ├ locals.tf
 | | ├ main.tf
 | | ├ null-resource.tf
 | | ├ outputs.tf
 | | ├ providers.tf
 | | ├ ssh-key.tf
 | | └ variables.tf
 | └ .gitkeep
 ├ IaC-app
 | ├ modules
 | | ├ deployment
 | | | └ main.tf
 | | ├ traefik
 | | | ├ cert-manager
 | | | | ├ acme-issuer
 | | | | | └ main.tf
 | | | | ├ cert-manager
 | | | | | └ main.tf
 | | | | └ certificates
 | | | |   └ main.tf
 | | | ├ error-page
 | | | | └ main.tf
 | | | ├ ingress-route
 | | | | └ main.tf
 | | | ├ middlewares
 | | | | ├ middleware-compress
 | | | | | └ main.tf
 | | | | ├ middleware-dashboard-basic-auth
 | | | | | └ main.tf
 | | | | ├ middleware-error-page
 | | | | | └ main.tf
 | | | | ├ middleware-gateway-basic-auth
 | | | | | └ main.tf
 | | | | ├ middleware-kibana-basic-auth
 | | | | | └ main.tf
 | | | | ├ middleware-rabbitmq-basic-auth
 | | | | | └ main.tf
 | | | | ├ middleware-rate-limit
 | | | | | └ main.tf
 | | | | ├ middleware-redirect-https
 | | | | | └ main.tf
 | | | | └ middleware-security-headers
 | | | |   └ main.tf
 | | | ├ tlsoption
 | | | | └ main.tf
 | | | ├ tlsstore
 | | | | └ main.tf
 | | | └ traefik
 | | |   ├ util
 | | |   | └ values.yaml
 | | |   └ main.tf
 | | └ .gitkeep
 | ├ main.tf
 | ├ namespace.tf
 | ├ providers.tf
 | └ variables.tf
 ├ IaC-storage
 | ├ oracle
 | | ├ modules
 | | | ├ bucket
 | | | | └ main.tf
 | | | └ .gitkeep
 | | ├ data.tf
 | | ├ main.tf
 | | ├ providers.tf
 | | └ variables.tf
 | └ .gitkeep
 ├ src
 | ├ concurrency
 | | ├ barrier.go
 | | ├ channel.go
 | | ├ rpReadWriteLock.go
 | | ├ semaphore.go
 | | ├ waitGroup.go
 | | ├ waitGroupL.go
 | | └ wpReadWriteLock.go
 | ├ config
 | | └ config.go
 | ├ finances
 | | ├ Annuities.go
 | | ├ Annuities_test.go
 | | ├ Bonds.go
 | | ├ Bonds_test.go
 | | ├ Mortgage.go
 | | ├ Periods.go
 | | ├ SimpleInterest.go
 | | └ SimpleInterest_test.go
 | ├ mathutil
 | | ├ mathutil.go
 | | └ mathutil_test.go
 | ├ security
 | | └ security.go
 | ├ webfinances
 | | ├ public
 | | | ├ css
 | | | | └ home.css
 | | | ├ js
 | | | | ├ AdCompoundingPeriods.js
 | | | | ├ AdEqualPeriodicPayments.js
 | | | | ├ AdFutureValue.js
 | | | | ├ AdPresentValue.js
 | | | | ├ OaCompoundingPeriods.js
 | | | | ├ OaEqualPeriodicPayments.js
 | | | | ├ OaFutureValue.js
 | | | | ├ OaGrowingAnnuity.js
 | | | | ├ OaInterestRate.js
 | | | | ├ OaPerpetuity.js
 | | | | ├ OaPresentValue.js
 | | | | ├ SimpleInterestAccurate.js
 | | | | ├ SimpleInterestBankers.js
 | | | | ├ SimpleInterestOrdinary.js
 | | | | ├ bonds.js
 | | | | ├ bondsYTM.js
 | | | | ├ getParams.js
 | | | | ├ miscellaneous.js
 | | | | └ mortgage.js
 | | | └ .gitkeep
 | | ├ templates
 | | | ├ annuitydue
 | | | | ├ cp
 | | | | | ├ cp.html
 | | | | | ├ i-PMT-FV.html
 | | | | | ├ i-PMT-PV.html
 | | | | | └ i-PV-FV.html
 | | | | ├ epp
 | | | | | ├ epp.html
 | | | | | ├ n-i-FV.html
 | | | | | └ n-i-PV.html
 | | | | ├ fv
 | | | | | ├ fv.html
 | | | | | ├ n-i-PMT.html
 | | | | | └ n-i-PV.html
 | | | | └ pv
 | | | |   ├ n-i-FV.html
 | | | |   ├ n-i-PMT.html
 | | | |   └ pv.html
 | | | ├ bonds
 | | | | ├ bonds.html
 | | | | ├ convexity.html
 | | | | ├ currentprice.html
 | | | | ├ duration.html
 | | | | ├ macaulayduration.html
 | | | | ├ modifiedduration.html
 | | | | ├ taxfree.html
 | | | | ├ yieldtocall.html
 | | | | └ yieldtomaturity.html
 | | | ├ miscellaneous
 | | | | ├ averagerate.html
 | | | | ├ compfrequencyconv.html
 | | | | ├ depreciation.html
 | | | | ├ effectiveannualrate.html
 | | | | ├ growthdecay.html
 | | | | ├ miscellaneous.html
 | | | | ├ nominalrate.html
 | | | | └ nominalratevs.html
 | | | ├ mortgage
 | | | | ├ amortizationtable.html
 | | | | ├ costofmortgage.html
 | | | | ├ heloc.html
 | | | | └ mortgage.html
 | | | ├ ordinaryannuity
 | | | | ├ cp
 | | | | | ├ cp.html
 | | | | | ├ i-PMT-FV.html
 | | | | | ├ i-PMT-PV.html
 | | | | | └ i-PV-FV.html
 | | | | ├ epp
 | | | | | ├ epp.html
 | | | | | ├ n-i-FV.html
 | | | | | └ n-i-PV.html
 | | | | ├ fv
 | | | | | ├ fv.html
 | | | | | ├ n-i-PMT.html
 | | | | | └ n-i-PV.html
 | | | | ├ ga
 | | | | | ├ FV.html
 | | | | | ├ PV.html
 | | | | | └ ga.html
 | | | | ├ interestrate
 | | | | | ├ interestrate.html
 | | | | | └ n-PV-FV.html
 | | | | ├ perpetuity
 | | | | | ├ gp.html
 | | | | | ├ p.html
 | | | | | └ perpetuity.html
 | | | | └ pv
 | | | |   ├ n-i-FV.html
 | | | |   ├ n-i-PMT.html
 | | | |   └ pv.html
 | | | ├ simpleinterestaccurate
 | | | | ├ accurate.html
 | | | | ├ amountofinterest.html
 | | | | ├ interestrate.html
 | | | | ├ principal.html
 | | | | └ time.html
 | | | ├ simpleinterestbankers
 | | | | ├ amountofinterest.html
 | | | | ├ bankers.html
 | | | | ├ interestrate.html
 | | | | ├ principal.html
 | | | | └ time.html
 | | | ├ simpleinterestordinary
 | | | | ├ amountofinterest.html
 | | | | ├ interestrate.html
 | | | | ├ ordinary.html
 | | | | ├ principal.html
 | | | | └ time.html
 | | | ├ about.html
 | | | ├ annuitydue.html
 | | | ├ contact.html
 | | | ├ finances.html
 | | | ├ footer.html
 | | | ├ header.html
 | | | ├ index.html
 | | | ├ ordinaryannuity.html
 | | | ├ simpleinterest.html
 | | | └ welcome.html
 | | ├ WfAdCpPages.go
 | | ├ WfAdEppPages.go
 | | ├ WfAdFvPages.go
 | | ├ WfAdPvPages.go
 | | ├ WfBondsPages.go
 | | ├ WfMiscellaneousPages.go
 | | ├ WfMortgagePages.go
 | | ├ WfOaCpPages.go
 | | ├ WfOaEppPages.go
 | | ├ WfOaFvPages.go
 | | ├ WfOaGaPages.go
 | | ├ WfOaInterestRatePages.go
 | | ├ WfOaPerpetuityPages.go
 | | ├ WfOaPvPages.go
 | | ├ WfPages.go
 | | ├ WfSiAccuratePages.go
 | | ├ WfSiBankersPages.go
 | | ├ WfSiOrdinaryPages.go
 | | └ fields.go
 | ├ go.mod
 | ├ go.sum
 | └ main.go
 ├ .dockerignore
 ├ .gitignore
 ├ Dockerfile-prod
 ├ LICENSE
 ├ README.md
 └ fin-finance.code-workspace
```

## Terraform
[To install `Terraform`](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli), download the binary executable for the Operating System (OS) being used to a directory in the system's PATH environment variable.

To troubleshoot the `OCI Terraform Provider`:<br>
https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/terraformtroubleshooting.htm

To install/upgrade `Terraform on Windows Subsystem for Linux (WSL)`.
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
<a id="terraform_init"></a>
This command initializes a working directory containing Terraform configuration files. ***This is the first command you should run after writing a new Terraform configuration or cloning an existing configuration from version control. It is safe to run this command multiple times.*** See [init](https://developer.hashicorp.com/terraform/cli/commands/init) for more information.
```
$ terraform init
$ terraform -chdir=../tf init
```
Note the usage of the `-chdir` option. This option allows you to declare where the root of your terraform project is located.

This command creates an execution plan, which lets you preview the changes that Terraform plans to make to your infrastructure.
```
$ terraform plan
$ terraform plan -var-file="../tf_secrets.auto.tfvars"
```
Note the use of the `-var-file` option. It sets values for potentially many input variables declared in the root module of the configuration, using definitions from a `tfvars` file. Use this option multiple times to include values from more than one file.

<a id="terraform_apply"></a>
This command executes the actions proposed in a Terraform plan. See [apply](https://developer.hashicorp.com/terraform/cli/commands/apply) for more information.
```
$ terraform apply
$ terraform apply -auto-approve
$ terraform apply -var-file="../tf_secrets.auto.tfvars"
$ terraform apply -var="app_version=1.0.0" -auto-approve
```
Notice the usage of the `-auto-approve` and `-var` options. The former skips interactive approval of the plan before applying. `Terraform` ignores this option when you pass a previously-saved plan file. This is because `Terraform` interprets the act of passing the plan file as the approval. The latter sets a value for a single input variable declared in the root module of the configuration. Use this option multiple times to set more than one variable.

---
**Important**

Resources you provision accrue costs while they are running. It's a good idea, as you learn, to always run `terraform destroy` on your project.

---
<a id="terraform_destroy"></a>
To deprovision all objects managed by a `Terraform` configuration. See  [destroy](https://developer.hashicorp.com/terraform/cli/commands/destroy)  for more information.
```
$ terraform destroy
$ terraform destroy -auto-approve
$ terraform destroy -var-file="../tf_secrets.auto.tfvars"
$ terraform destroy -var="app_version=1.0.0" -auto-approve
```

Once `Terraform` finish setting up your resources, you need to set up `kubectl` to access the cluster. See [kubectl](#kubectl).

Finally, let's try to list the available nodes in the cluster.
```
$ kubectl get nodes
```
If the nodes are displayed, you are done.

## Installing kubectl
`kubectl` is a [command line tool](https://kubernetes.io/docs/tasks/tools/) for communicating with a `Kubernetes` cluster's control plane, using the `Kubernetes API`.

---
**Note**

A file that is used to configure access to a cluster is usually referred to as a `kubeconfig file`. This is a conventional way of referring to a configuration file, often shortened to config file. It does not imply that a file named kubeconfig exists.

---
You will need to create a kubeconfig file with authentication and configuration details, which will allow kubectl to communicate with your cluster. To create the kubeconfig file, you execute the command below, which requires the following information:<br>
**(1)** Cluster's OCID (Navigation menu->Developer Services->Kubernetes Clusters (OKE) [Under Containers & Artifacts]->Select the compartment that contains the cluster[Compartment]-> On the Clusters page, click the name of the cluster)<br>
**(2)** Name for the config file<br>
**(3)** The region
```
$ oci ce cluster create-kubeconfig --cluster-id <cluster OCID> --file ~/.kube/<name-of-config-file> --region <region> --token-version 2.0.0 --kube-endpoint PUBLIC_ENDPOINT
```
The command will create a kubeconfig file in the `~/.kube` directory; the kubeconfig file will contain the keys and all of the configuration for `kubectl` to access the cluster. See `IaC-K8s/oracle/data.tf` for appropriate values to the parameters `--token-version` and `--kube-endpoint`.

---
**Note**

Setting the permissions of your `~/.kube/<name-of-config-file>` file to `600` ensures that only the owner (you) can read and write to it, enhancing security by limiting access to the Kubernetes configuration file.

```
$ chmod 600  ~/.kube/<name-of-config-file>
```
---
By default, `kubectl` looks for a file named `config` in the `$HOME/.kube (~/.kube)` directory; hence, if the `KUBECONFIG` environment variable is not set, `kubectl` uses the default values `~/.kube/config`. You can specify other kubeconfig files by setting the `KUBECONFIG` environment variable or by setting the `--kubeconfig` flag.

To export the `KUBECONFIG` environment variable ***only*** for the current shell and its children processes, you use the `export` command.
```
export KUBECONFIG=<name-of-config-file>
```
To reiterate, when an environment variable is set from the shell using the export command, its existence ends when the current session ends.

To set the `KUBECONFIG` environment variable as a `user-specific environment variable`, add the `export` command to `~/.bashrc (bash), ~/.kshrc (ksh), or ~/.zshrc (zsh)`, depending on which shell you are using. By modifying the shell-specific configuration file, the environment variable will persist across sessions and system restarts. Note the use of the `bash shell` below.
```
$ echo 'export KUBECONFIG=<name-of-config-file>' >> ~/.bashrc
```
Next, reload the file to apply the changes.
```
$ source ~/.bashrc
```

To view all environment variables, use the `printenv` command. Since there are many variables on the list, use the `less` command to control the view.
```
$ printenv | less
```
The output shows the first page of the list and allows you to move forward by pressing `Space` to see the next page or `Enter` to display the next line. Exit the view with `q`.

To view a specific environment variable, use the `set` command.
```
$ set | grep KUBECONFIG

or

$ echo $KUBECONFIG

or

$ printenv KUBECONFIG
```

## Kubenertes (K8s)
`K8s` is an open source container orchestration platform.

### Useful Commands
#### version
Display the `Kubernetes` version running on the client and server.
```
$ kubectl version
```

#### config
Display the configuration of the cluster.
```
$ kubectl config view
```

Display all users.
```
$ kubectl config view -o jsonpath='{range .users[*]}{.name}{"\n"}{end}'
```

Retrieve one user details.
```
$ kubectl config view -o jsonpath='{.users[?(@.name == "<user-name>")].user}{"\n"}'
```

#### cluster
Retrieve cluster details.
```
$ kubectl cluster-info
```

#### node
Confirm what platform is running on the cluster.
```
$ kubectl describe node | grep "kubernetes.io/arch"
```

To retrieve nodes information.
```
$ kubectl get nodes
```

#### exec
---
***Note***

The double dash (`--`) in the command signals the end of command options for `kubectl`. Everything after the double dash is the command that should be executed inside the pod; the double dash is required.

The command takes the following options:<br>
`-i` or `--stdin`: Keep stdin open even if not attached.<br>
`-t` or `--tty`: Allocate a pseudo-TTY.<br>
`-c` or `--container`: Specify the container name (useful for pods hosting multiple containers).<br>
`-n` or `--namespace`: Specify the namespace of the pod.

---
To open an interactive shell (e.g.; `bash`) in a pod hosting one container, execute the command below.
```
$ kubectl exec -it <pod-name> -n <name-space> -- /bin/bash
```

Since pods are capable of hosting multiple containers, you can specify a specific container by using the `-c` flag.
```
$ kubectl exec -it <pod-name> -c <container-name> -n <name-space> -- /bin/bash
```

To execute a single command without entering an interactive shell, use.
```
$ kubectl exec <pod-name> -n <name-space> -- env
```

#### Pods
Display what node a pod is scheduled.
```
$ kubectl get pods -o wide -n <name-space>
```

Retrieve a list of host IP addresses with the additional `phase` field indicating if the pod is running or not.
```
$ kubectl get pods -o jsonpath='{range .items[*]}{.status.hostIP}{"\t"}{.status.phase}{"\n"}{end}' -n <name-space>
```

Retrieve pods across all namespaces.
```
$ kubectl get pods --all-namespaces
```

Retrieve pods under a specific namespace.
```
$ kubectl get pods -n <name-space>
```

Display details of a specific pod under a specific namespace.
```
$ kubectl describe pod <pod-name> -n <name-space>
```

Retrieve all containers in a pod with all of their ports.
```
$ kubectl get pod <pod-name> -o jsonpath='{range .spec.containers[*]}{.name}{"\t\t"}{range .ports[*]}{.name}{"="}{.containerPort}{"\t\t"}{end
}{"\n"}{end}' -n <name-space>
```

Delete a pod under a specific namespace.
```
$ kubectl delete pod <pod-name> -n <name-space>
```

Delete all pods without specifying their names.
```
$ kubectl delete pods --all
```

#### logs
Retrieve the logs for a pod under a specific namespace.
```
$ kubectl logs <pod-name> -n <name-space>
```

Retrieve the logs for a pod that was previously running under a specific namespace.
```
$ kubectl logs <pod-name> -n <name-space> --previous
```

Retrieve the logs for a specific container running in a pod under a specific namespace. If the pod has only one container, the container name is optional.
```
$ kubectl logs <pod-name> -c <container-name> -n <name-space>
```

#### Volumes
---
**Note**

PersistentVolume resources are cluster-scoped and thus cannot be created in a specific namespace. On the other hand, PersistentVolumeClaims can only be created in a specific namespace, and they can then only be used by pods in the same namespace.

For more in-depth information, see [Persistent Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)

---
List all PersistentVolumes.
```
$ kubectl get pv
```

List all PersistentVolumeClaims.
```
$ kubectl get pvc -n <name-space>
```

List storage classes.
```
$ kubectl get storageclass

or

$ kubectl get sc
```

Display more information about the given storage class.
```
$ kubectl get sc <storage-class-name> -o yaml
```

#### Resources
Retrieve built-in resource types (pods, services, daemon sets, deployments, replica sets, jobs, cronjobs, and stateful sets) under a specific namespace.
```
$ kubectl get all -n <name-space>
```
---
**Note**

The `kubectl delete` command might not be successful initially if you use `finalizers` to prevent accidental deletion. Finalizers are keys on resources that signal pre-delete operations. Finalizers control the garbage collection on resources, and they're designed to alert controllers about what cleanup operations to do before they remove a resource.

If you try to delete a resource that has a finalizer on it, the resource remains in finalization until the controller removes the finalizer keys, or the finalizers are removed by using `kubectl`. After the finalizer list is emptied, `Kubernetes` can reclaim the resource and put it into a queue to be deleted from the registry.

See [Using Finalizers to Control Deletion](https://kubernetes.io/blog/2021/05/14/using-finalizers-to-control-deletion/) for more information.

---
Remove a resource in the `Terminating` state.<br>
To remove a `finalizer` from a resource, you typically update the resource's metadata to remove the finalizer entry. This action signals `Kubernetes` that the cleanup tasks are complete, allowing the resource to be fully deleted.

To ensure the resource has one or more finalizers attach, you can use `kubectl get` or `kubectl describe`. If finalizers are attached, you remove them by executing the command below.
```
$ kubectl patch <resource> <resource-name> -p '{"metadata":{"finalizers":null}}'
```

## IaC-K8s
IaC-K8s contains the `Terraform` code for provisioning (i.e., creating, preparing, and activating the underlying infrastructure of a cloud environment) the `Oracle Cloud Infrastructure (OCI)`, which is an `Infrastructure as a Service (IaaS)` and `Platform as a Service (PaaS)` offering. The `OCI` is a set of complementary cloud services that enable you to build and run a range of applications and services in a highly available hosted environment. `OCI` provides high-performance compute capabilities (as physical hardware instances) and storage capacity in a flexible overlay virtual network that is securely accessible from your on-premises network.

For more `OCI` and `Terraform` documentation, please see [Using Terraform and Oracle Cloud Infrastructure](https://docs.oracle.com/en-us/iaas/Content/dev/terraform/tutorials.htm).

### Login
1. https://www.oracle.com/cloud/sign-in.html
2. Sign In using a Cloud Account Name
3. Cloud Account Name

### Installing Oracle Cloud Infrastructure (OCI) Command Line Interface (CLI)
[OCI CLI](https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/climanualinst.htm#Manual_Installation) is a tool that allows users to interact with OCI services directly from the command line.

**Manual Installation: Ubuntu**<br>
**Step 1: Installing Python**<br>
Before you install the CLI, run the following commands on a new Ubuntu image.<br>
```
~$ sudo apt update
~$ sudo apt install build-essential zlib1g-dev libncurses5-dev libgdbm-dev libnss3-dev libssl-dev libreadline-dev libffi-dev libsqlite3-dev wget libbz2-dev
~$ sudo apt update && sudo apt install python3.12.0 python3.12.0-pip python3.12.0-venv
```

**Step 2: Creating and Configuring a Virtual Environment**<br>
The `venv` Python module is a virtual environment builder that lets you create isolated Python environments.<br>
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
To install using `PyPI`, run the following command.
```
(oracle-cli) ~/oci/python$ pip install oci-cli
```

**Step 4: Setting up the Configuration File**<br>
Before using the CLI, you must create a configuration file that contains the required credentials for working with `Oracle Cloud Infrastructure`. The default location for the configuration file is `~/.oci`.<br>
Use the Setup Dialog<br>
To have the CLI guide you through the first-time setup process, use the setup config command.
```
(oracle-cli) ~/oci/python$ oci setup config
```
This command prompts you for the information required to create the configuration file and the API public and private keys. The setup dialog uses this information to generate an API key pair and creates the configuration file. After API keys are created, upload the public key using the Console. You will need the following:<br>
1. User's OCID (Profile->User settings or My profile)
2. Tenancy's OCID (Profile->Tenancy: \<tenancy-name\>)
3. The region

When creating the keys, decline creating a passphrase. Once the keys are generated, you'll need to associate the public key to the user. From the Oracle Cloud web console, click on `Profile-> My profile->API keys`* on the left and click on `Add API Key`. Upload the public key's pem file.

**Step 5: Verify that everything is configured properly**<br>
You can verify that everything is configured properly by running the following command.
```
(oracle-cli) ~/oci/python$ oci iam compartment list -c <tenancy-ocid>
```
where \<tenancy-ocid\> is your tenancy's OCID.

If there are no errors in the `JSON` reply, the config file was create (by default in `~/.oci`). At this point, you need to run `Terraform` to allocate your resources.

**Step 6: Deactivate the virtual environment.**
```
(oracle-cli) ~/oci/python$ deactivate
```

**Step 7: Activate the virtual environment.**
```
~$ source ~/oci/python/oracle-cli/bin/activate
```

## Traefik (Gateway/Reverse Proxy and Load Balancer)
In a typical microservices deployment, microservices are **not** exposed directly to client applications; i.e., microservices are behind a set of APIs that is exposed to the outside world via a gateway. **The gateway is the entry point to the microservices deployment, which screens all incoming messages for security and other quality of service (QoS) features.** Since the gateway deals with **north/south traffic**, it is mostly about **edge security**. To reiterate, **the gateway is the only component publicly accessible for requests originating from the external internet**. Besides `Traefik`, there are many options for reverse proxies available such as `Nginx`, `Pomerium` (free), `Apache`, `Caddy`, `Envoy`, `Zuul`, and `HAProxy`.

### Troubleshooting Traefik
For a more detailed explanation, please see [Traefik, Let’s Encrypt, Cert-Manager, and OpenShift using Terraform (Part 4)](https://trimino.com/simple-app/traefik-lets-encrypt-cert-manager-and-openshift-using-terraform-part-4/)

Execute commands in a running `Traefik` container.
```
$ kubectl exec -it -n finances $(kubectl get pods -n finances --selector "app.kubernetes.io/name=traefik" --output=name) -- /bin/sh
```

```
$ kubectl get pod,middleware,ingressroute,svc -n finances
$ kubectl get all -l "app.kubernetes.io/name=traefik" -n finances
$ kubectl get all -l "app=finances" -n finances
```

### Troubleshooting Certificates
```
$ kubectl get svc,pods -n finances
$ kubectl get Issuers,ClusterIssuers,Certificates,CertificateRequests,Orders,Challenges -n finances
$ kubectl get Issuers,ClusterIssuers,Certificates,CertificateRequests,Orders,Challenges --all-namespaces
$ kubectl describe Issuers,ClusterIssuers,Certificates,CertificateRequests,Orders,Challenges -A
$ kubectl describe Issuers,ClusterIssuers,Certificates,CertificateRequests,Orders,Challenges -n finances
```

Check the certificate.
```
$ kubectl -n finances describe certificate <certificate-name>
```

Delete a certificate.
```
$ kubectl -n finances delete certificate <certificate-name>
```

To describe a specific resource (the resource name can be obtained from the kubectl get command).
```
$ kubectl -n finances describe Issuer <issuer-name>
$ kubectl get ingressroute -A
$ kubectl get ingressroute -n finances
```

To delete a pending `Challenge`, see [here](https://cert-manager.io/docs/installation/helm/#uninstalling) and [here](https://cert-manager.io/docs/installation/uninstall/). As per documentation, the order is important!!!
```
$ kubectl delete Issuer <issuer-name> -n finances
$ kubectl delete Certificate <certificate-name> -n finances
```

## src
### Initializing a Go Project
In version 1.13, `Go` added a new way of managing the libraries a `Go project` depends on, called [Go modules](https://go.dev/ref/mod). A `Go module` has a number of `Go` code files implementing the functionality of a `package`, but it also has two additional and important files in the root: the `go.mod` and `go.sum` files. These files contain information the go tool uses to keep track of the module's configuration and are commonly maintained by the tool.

With `Go modules`, it is possible for `Go projects` to be located anywhere on the filesystem instead of a specific directory defined by `Go`. Having said that, you'll create the project directory `fin-finance` with the module directory `src`.
```
$ mkdir -p fin-finance/src
$ cd fin-finance/src
```
Notice the usage of the `-p` option to create parent directories that do not exist. Once you've executed the `mkdir` command, the directory structure will look like
```
fin-finance
 └ src
```
Next, you'll create a `go.mod` file within the `src` directory to define the `Go module` itself. To do this, you'll use the go tool's `mod init` command and provide it with the module's name, which in this case is `finance`. Usually, the module's directory name is the same as the module name, but in your case they will not be the same.
```
~/fin-finance/src$ go mod init finance
```
Once created, the `go.mod` file contains the name of the module and versions of other modules your own module depends on. It can also contain other directives, such as replace, which can be helpful for doing development on multiple modules at once.

### Upgrading Go
Find the system architecture type.
```
$ dpkg --print-architecture
```
Update the package index files on the system, which contain information about available packages and their versions. It downloads the most recent package information from the sources listed in the `/etc/apt/sources.list` file that contains your sources list.
```
$ sudo apt-get update
```
Install the newest versions of all packages currently installed on the system from the sources enumerated in `/etc/apt/sources.list`.
```
$ sudo apt-get -y upgrade
```
Remove the existing `Go` installation, if any. Find the location of your current `Go` installation.
```
$ which go
```
Remove the directory found in the previous step, if any.
```
$ sudo rm -rf /usr/local/go
```
Download the [latest Go version](https://go.dev/dl/). Notice the version and system architecture type; e.g., go**1.22.0.linux-amd64**.tar.gz
```
$ wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
```
Extract the archive (replace with the actual filename).
```
$ sudo tar -xvf go1.22.0.linux-amd64.tar.gz -C /usr/local
```
Update the `PATH environment variable`. Open a shell configuration file; e.g., `~/.bashrc`, `~/.zshrc`, `~/.profile`, etc. Add the following line to the `~/.bashrc` config file; note the use of the bash shell below.
```
$ echo "export PATH=/usr/local/go/bin:$PATH" >> ~/.bashrc
```
Reload the file to apply the changes.
```
$ source ~/.bashrc
```
Verify the installation.
```
$ go version
```
Finally, you may need to update the `go.mod` file with the new compiler version.

### Debugging Go
#### Delve
To install the debugger in VS Code:<br>
**(1)** Open the Command Palette (***Ctrl + Shift + P***).<br>
**(2)** Find ***Go: Install/Update Tools*** and select ***dlv***.

The settings for the debugger can be stored in the `.code-workspace` file or the `.vscode/launch.json` directory. For this project, the settings are stored in the `.code-workspace` file under the `launch` section.

### Profiling Go with pprof
Because performance is an important aspect of software development, profiling is a valuable tool for understanding and improving the performance of Go applications. The `Go profiler (pprof)` is a tool for profiling Go applications. It is part of the Go standard library and can be used to generate detailed profiles of Go applications, including CPU, memory, and concurrency profiles. It reads profiling samples in the `profile.proto` format and generates both text and graphical reports.

One way to enable `pprof` is to use the [net/http/pprof](https://pkg.go.dev/net/http/pprof) package to serve the profiling data via `HTTP`. (This assumes that your application has an HTTP server running; otherwise, you will need to start one.) If you use the `blank import`, the profile package will **only** register its handlers with the default multiplexer (`http.DefaultServeMux`). If you are not using the default multiplexer, you will need to register the handlers with the multiplexer you're using. Once the handlers are registered, you can reach the `pprof URL` via `http://{url}:{port}/debug/pprof`.

Note that [enabling pprof is safe](https://go.dev/doc/diagnostics#profiling) even in production. The profiles that impact performance, such as CPU profiling, aren't enabled by default, nor do they run continuously; they are activated only for a specific period. Nonetheless, exposing `pprof` endpoints can lead to potential security risks and unintended performance degradation; it's crucial to secure access to the endpoints.

---
**Note**

The `graphviz` package is part of the `Universe` repository, which is enabled by default on most Ubuntu installations.

The `Universe` repository is a standard repository for `Ubuntu`. The repository is community-maintained and provides free and open-source software. By default, the repository is enabled in the latest versions of Ubuntu, but if for some reason is not enabled, use the following command to enable it.
```
$ sudo add-apt-repository universe
$ sudo apt update
```
The second command performs an update to the package list cache.

To remove the repository, use the command below. The command will **not** remove packages that were installed from the repository, if any.
```
$ sudo add-apt-repository -r universe
$ sudo apt update
```
Again, the second command performs an update to the package list cache thereby ensuring the repository is no longer usable.

---
Please note that you will need to have [graphviz](https://graphviz.org/) installed for web visualizations. To install it, run the commands below.
```
$ sudo apt install graphviz
$ sudo apt update
```
To confirm the installation, display the version of `graphviz`.
```
$ dot -V
```
To remove `graphviz`, use the commands below.
```
$ sudo apt purge graphviz
$ sudo apt update
```

#### Pprof Endpoints (Profiles)
To view all available profiles, open your browser and type the following address into the browser's address bar: `http://{url}:{port}/debug/pprof/`.

#### Analyzing the Results of pprof
To analyze the results of `pprof`, you can use the command `go tool pprof`. You can run the command on interactive mode or via a web interface; furthermore, the command can read a profile from a file or directly from a server via `HTTP`.

##### Interactive Mode
**To read the profiling statistics from a file**, export the statistics using [`curl`](https://curl.se/docs/manpage.html).
```
$ curl http://{url}:{port}/debug/pprof/{endpoint}?seconds={x} --output ./{filename}
```
In the preceding command, `pprof` is being requested to profile the application for **x** seconds.<br>
Next, use the interactive mode to read the profile. This opens an interactive shell that allows running interactive commands for analyzing the profile.
```
$ go tool pprof ./{filename}
```
The `pprof` tool provides various commands for interacting with it. To display a list of the available commands and their usage, use one of the two commands shown below.
```
$ go tool pprof -h

or

$ go tool pprof -help
```
**To run `pprof` in interactive mode directly from a server via `HTTP`**, you will use the command shown below.
```
$ go tool pprof http://{url}:{port}/debug/pprof/{endpoint}?seconds={x}
```

##### Web Interface
---
**Note**

In order to use `pprof` with a web interface, your system will require a default browser.

***Linux***<br>
To set a default browser in `Linux` using environment variables, you can use the command
```
export BROWSER=/path/to/desired/browser
```
This will set the `BROWSER` environment variable to the path of your desired browser. To make this change permanent, you need to add this command to your `.bashrc` file or another relevant startup script.

***Windows Subsystem for Linux (WSL)***<br>
To set the `BROWSER` environment variable, you can use the command
```
export BROWSER='/mnt/c/Windows/explorer.exe'
```
This command sets the default web browser to the `Windows File Explorer`; i.e., it will use the default browser of the host `Windows`. To make this change permanent, as with Linux, you need to add this command to your `.bashrc` file or another relevant startup script.

---
**To analyze the profiling statistics from a file**, export the statistics using [`curl`](https://curl.se/docs/manpage.html).
```
$ curl http://{url}:{port}/debug/pprof/{endpoint}?seconds={x} -o ./{filename}
```
Next, analyze the result using the web interface. Note that when the `-http` flag is specified, `pprof` starts a web server at the specified **url:port** that provides an interactive web-based interface. Both **url** and **port** are optional. The default value for **url** is **localhost**. For **port**, its default value is **a random available port**; otherwise, any port not being used by another application. Hence, **-http=":"** starts a server locally at a random port. This command should automatically open the web browser at the right page; if not, manually navigate to the specified **url:port** in the web browser.
```
$ go tool pprof -http="{url}:{port}" ./{filename}
```
**To analyze the profiling statistics directly from a server via `HTTP`**, use the command shown below.
```
go tool pprof -http={url1}:{port1} http://{url2}:{port2}/debug/pprof/{endpoint}?seconds{x}
```
In the preceding command, **url1:port1** will be used by the server displaying the profiling statistics whereas **url2:port2** refers the application being profiled.


[text](https://pkg.go.dev/net/http/pprof#pkg-overview)
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx


**CPU Profiling**
When it is activated, the application asks the OS to interrupt it every 10ms (default). When the application is interrupted, it suspends the current activity and transfers the execution to the profiler. The profiler collects execution statistics, and then it transfers execution back to the application.

To active the CPU profiling, you access the `debug/pprof/profile` endpoint. Accessing this endpoint will execute CPU profiling for 30 seconds by default. For 30 seconds, the application is interrupted every 10ms.

To write the output to a file, use the command below:
```
$ curl http://{url}:{port}/debug/pprof/{prof1}?seconds={x} --output {filename}
```
    where {prof1} is trace or profile.

```
$ curl http://{url}:{port}/debug/pprof/{prof2} --output {filename}
```
where {prof2} is heap.

To inspect a file.
```
$ go tool pprof {filename}
```

To inspect the result using the graphical user interface.
```
$ go tool pprof -http=:{port1} {filename}
```

To directly connect to the debug point.
```
$ go tool pprof http://{url}:{port}/debug/pprof/{prof1}?seconds={x}
$ go tool pprof http://{url}:{port}/debug/pprof/{prof2}
```

To inspect the result using the graphical user interface, use the command below:
```
$ go tool pprof -http=:{port1} http://{url}:{port2}/debug/pprof/{prof1}?seconds={x}
$ go tool pprof -http=:{port1} http://{url}:{port2}/debug/pprof/{prof2}
```




### Compile and Run the App
---
**Note**

To see all environment variables supported by the app, see `//Environment variables.` in [main.go](./src/main.go). To run the app in a `K8s` environment, set the environment variable `K8S` to true.

---
Compile and run the app as a standalone HTTP server (default) on port 8080 (default).
```
$ go build -o finance && ./finance

or

$ go build -o finance && HTTP=true HTTP_PORT=8080 ./finance
```
Compile and run the app as a standalone HTTPS server on port 8443 (default).
```
$ go build -o finance && HTTP=false HTTPS=true ./finance

or

$ go build -o finance && HTTP=false HTTPS=true HTTP_PORT=8443 ./finance
```
Compile and run the app as two standalone servers (HTTP and HTTPS ) using the default ports.
```
$ go build -o finance && HTTPS=true ./finance
```
Force rebuilding of packages that are already up-to-date and run the app.
```
$ go build -o finance -a && ./finance
```
To change an environment variable's value, set the environment variable to its new value; e.g., to change the default value of the environment variable HTTP_PORT, execute the command below.
```
$ HTTP_PORT=18080 ./finance

for multiple environment variables

$ HTTP_PORT=18080 HTTPS=true ./finance
```
The ampersand symbol (`&`) instructs the shell to execute the command as a separate background process. To compile and run the app in the background.
```
$ go build -o finance && ./finance &
```
Compile and run the named `main` Go package in the background.
```
$ go run main.go &
```

### Application Deployment with Terraform
See [terraform init](#terraform_init).
```
$ terraform init
```
From the same directory where you invoked the `init` command, run the `apply` command; this command gathers together and executes all of our `Terraform` code files. The option `-auto-approve` runs Terraform in `non-interactive` mode. See [terraform apply](#terraform_apply).
```
$ terraform apply -var="k8s_manifest_crd=false" -auto-approve
```
This command sets the variable `app_version`, enables non-interactive mode, and invokes `apply`.
```
$ terraform apply -var="k8s_manifest_crd=false" -var="app_version=1.0.1" -auto-approve
```
To deploy the reverse proxy `Traefik` after initializing `Terraform`, you’ll execute any one of the two commands below (they are equivalent since the default value for the variable `k8s_manifest_crd` is true; see [variables.tf](./IaC-app/variables.tf)). For more information see [Deploying Traefik in Our OpenShift Cluster (Part 3)](https://trimino.com/simple-app/deploy-traefik-openshift/), section `Building and Deploying Traefik`.
```
$ terraform apply -var="app_version=1.0.1" -auto-approve

or

$ terraform apply -var="k8s_manifest_crd=true" -var="app_version=1.0.1" -auto-approve
```
This command destroys your current infrastructure that was created by `Terraform`. See [terraform destroy](#terraform_destroy).
```
$ terraform destroy -auto-approve

or

$ terraform destroy -var="app_version=1.0.1" -auto-approve
```

#### How to kill a process.
##### Windows
```
C:\> netstat -ano | findstr :<port>
C:\> taskkill /PID <PID> /F

or

C:\> npx kill-port <port>
```

##### Linux
```
$ ps -a
$ kill <PID>
```

#### Display Headers
##### Windows (PowerShell)
```
PS> curl.exe -IL "http://localhost:8080"
```

##### Linux
```
$ curl.exe -IL "http://localhost:8080"
```
