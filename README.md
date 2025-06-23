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
 | | |   | ├ traefik-scc.yaml
 | | |   | └ values.yaml
 | | |   └ main.tf
 | | └ .gitkeep
 | ├ .terraform.lock.hcl
 | ├ backend.tf
 | ├ iac.sh
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
### Generated Files
`Terraform` generates several types of files during the provisioning of infrastructure.

#### terraform.tfstate
It is the file that `Terraform` uses to track the state of the infrastructure it manages. The state file contains information about the resources that `Terraform` is managing to determine which resources need to be created, updated, or deleted to match the current configuration. After invoking the command `terraform apply` for the first time, `Terraform` will generate the state file `terraform.tfstate`. Subsequent invocations of `terraform apply` use the state file as input. `Terraform` loads the state file and then *refreshes* it from the current infrastructure.

`Terraform` stores its state files in a [backend](https://developer.hashicorp.com/terraform/cli/commands/init#backend-initialization) ([backend.tf](https://github.com/juan-carlos-trimino/fin-finance/blob/main/IaC-app/backend.tf)). But if a backend is not explicitly specified, `Terraform` uses the local backend. By default, the local backend stores the file (`terraform.tfstate`) in the same directory where the command `terraform apply` was run. But since `Terraform` state files contain all data in plain text, it is a security problem for certain `Terraform` resources that need to store sensitive data. Hence, **it is not recommended to store state files in a source control**; sensitive data must always be stored in a secure location and never in a source control. Furthermore, when working with a team of developers, each member of the team requires access to the same state file; i.e., the file needs to be stored in a shared location. Unfortunately, most source control systems do not allow the locking of files, which may cause issues when multiple users attempt to access the file at the same time.

The question then is how to persist the state file? `Terraform` provides a solution that uses external storage to store the state file. To implement the solution, a three-step process is required:<br>
**(1)**	Write `Terraform` code to create the external storage (let’s say *Simple Storage Service* or *S3*). Then deploy the code with a local backend.<br>
**(2)**	In the `Terraform` code written to deploy the S3 storage, change the backend to use the S3 storage and run the command `terraform init` to copy the local state to the S3 storage.<br>
**(3)**	In the `Terraform` code written to deploy the application, change the backend to use the S3 storage and run the command `terraform init` to copy the local state to the S3 storage. ENSURE THE PATH/FILE NAMES ARE DIFFERENT TO AVOID OVERRIDING THE FILE STATES.

If there is a need to delete the S3 storage, repeat the three-step process in reverse:<br>
**(1)**	In the `Terraform` code written to deploy the application, remove the backend configuration. Then run the command `terraform init` to copy the `Terraform` state back to the local disk.<br>
**(2)**	In the `Terraform` code written to deploy the S3 storage, remove the backend configuration. Then run the command `terraform init` to copy the `Terraform` state back to the local disk.<br>
**(3)**	Execute the command `terraform destroy` to delete the S3 storage.

Although this solution is a bit awkward, keep in mind that after the S3 storage exists, the rest of the `Terraform` code can simply specify the backend configuration right from the start without any extra steps. But perhaps the most [awkward limitation](https://developer.hashicorp.com/terraform/language/backend#define-a-backend-block) is that the backend block does not allow the use of variables or references.

#### terraform.tfstate.backup
It is the backup file of the `terraform.tfstate` file. `Terraform` automatically creates a backup of the state file before making any changes to the state file. This ensures recovering from a corrupted or lost state file possible. This file is stored in the same directory as the `terraform.tfstate` file.

The `terraform.tfstate.backup` file can be used to restore the `Terraform` state to the previous version. To do so, just rename the `terraform.tfstate.backup` file to `terraform.tfstate` and run the command `terraform init`.

#### .terraform.lock.hcl
`Terraform` automatically creates or updates the [dependency lock](https://developer.hashicorp.com/terraform/language/files/dependency-lock) file each time the command `terraform init` is run. This file tracks the versions of providers and modules used in a configuration thereby ensuring all subsequent runs of `terraform apply` or `terraform plan` use the same provider versions, preventing unexpected changes due to updates or different environments. The file is [typically located](https://developer.hashicorp.com/terraform/language/files/dependency-lock#lock-file-location) in the same directory as the root module and **is recommended to be included in [version control](https://developer.hashicorp.com/terraform/language/files/dependency-lock#lock-file-location)**.

When the backend configuration changes to a different location, the state and lock files will be moved to the new location specified in the backend configuration. These files are managed by the backend and are typically located alongside each other. To prevent the corruption of the state file, the lock file ensures that only one process can modify the state file at a time.

#### .terraform/
`Terraform` creates a hidden `.terraform/` directory, which serves as a working directory to cache provider plugins and modules, records which workspace is currently active, and records the last known backend configuration. `Terraform` automatically manages this directory and creates it during initialization. Since this directory may contain sensitive credentials for the remote backend, **it should not be included in a [version control](https://developer.hashicorp.com/terraform/language/backend#initialize-the-backend)**, nor should the contents of the directory be modified directly.

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

<a id="shell"></a>
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
#### cluster
Retrieve cluster details.
```
$ kubectl cluster-info
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

#### exec
---
**Note**

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

#### namespace
List all the namespaces within a cluster.
```
$ kubectl get namespaces --all-namespaces

or

$  kubectl get namespaces
```

---
**Note**

If the **STATUS** of a namespace is displaying **Terminating**, you will need to follow the steps below to delete the namespace.

1. Get the name of the namespace that is stuck in the **Terminating** state.
```
$ kubectl get namespaces
```

2. Select the namespace that is stuck in the **Terminating** state and save the contents of the namespace in a JSON file.
```
$ kubectl get ns <terminating-namespace> -o json > ns.json
```

3. Edit the JSON file by removing the **kubernetes** value from the **finalizers** field. Save the file.
```
$ vi ns.json
...
"spec": {
  "finalizers": [
    "kubernetes"
  ]
},
...
```

4. After removing the **kubernetes** value from the **finalizers** field, the contents of the JSON file should resemble the following listing.
```
$ cat ns.json
...
"spec": {
  "finalizers": [
  ]
},
...
```

5. To apply the change, run the command below.
```
$ kubectl replace --raw "/api/v1/namespaces/<terminating-namespace>/finalize" -f ./ns.json
```

6. Verify that the terminating namespace has been removed.
```
$ kubectl get namespaces
```

---

#### node
Confirm what platform is running on the cluster.
```
$ kubectl describe node | grep "kubernetes.io/arch"
```

To retrieve nodes information.
```
$ kubectl get nodes
```

#### pods
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
$ kubectl delete pods --all -n <name-space>
```

#### port-forward
Forward a local port to a port on a K8s resource. This command establishes a *secure tunnel* between your local machine and the resource allowing you to access the resource as if it were running locally. This is particularly useful for troubleshooting and debugging resources that are not exposed externally.
```
$ kubectl port-forward -n <namespace> <resource-type>/<resource-name> <local-port>:<resource-port>
```
Where:<br>
*resource-type* specifies the type of K8s resource. It defaults to pod, if omitted.<br>
*resource-name* specifies the name of the K8s resource.<br>
*local-port* is the port number on your local machine; i.e, this is the port that you'll use to access the resource from your local machine.<br>
*resource-port* is the port number for the K8s resource; i.e., the traffic sent to the local port on your machine will be forwarded to this port of the resource.

---
**Note**

`kubectl port-forward` does not return. To continue, you will need to open another terminal. The port-forwarding session ends when you manually close it by pressing ***Ctrl+C***, which sends an interrupt signal to the process, in the terminal where the command was initiated, when the resource being forwarded terminates, or when a timeout is reached due to inactivity.

To run this command in the background, you can add the ***&*** operator at the end of the command.
```
$ kubectl port-forward -n <namespace> <resource-type>/<resource-name> <local-port>:<resource-port> &
```

To stop the background process, you will need to find its process ID by executing this command.
```
$ ps -aux | grep -i "kubectl port-forward"
```
Where:
*ps aux* displays information about all active processes from all users.
*grep -i "kubectl port-forward"* filters the results to only show lines containing the `kubectl port-forward` expression; the search is case-insensitive (-i).

Locate the relevant line, then kill the process by running the command below.
```
$ kill <pid>
```
Where:
*pid* is the process ID of the `kubectl port-forward` command.

---

If you don't need a specific local port, you can omit the local port; `kubectl` automatically choose and allocate a local port.
```
$ kubectl port-forward -n <namespace> <resource-type>/<resource-name> :<resource-port>
```

To get a list of all forwarded ports active, you can use the following command.
```
$ ps -ef | grep -i "port-forward"
```
Where:<br>
*ps -ef* lists all processes running on the system with detailed information, including the process ID, command line arguments, and user.
*grep -i "port-forward"* filters the results to only show lines containing the port-forward expression; the search is case-insensitive (-i).

#### resources
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

#### version
Display the `Kubernetes` version running on the client and server.
```
$ kubectl version
```

#### volumes
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

## Docker
`Docker` is a platform for packaging and deploying containers. You use `Docker` to:
1. Package your microservice into a `Docker` image.
2. Publish your image to your private image registry.
3. Run your microservice in a container.

To install `Docker`, go to the `Docker` website at https://docs.docker.com. Once there, find the `Get Docker` link and follow the instructions to install `Docker` for your platform.

With `Docker` installed, you package your microservice with the following steps:
1. Create a `Dockerfile` for your microservice.
2. Package your microservice as a `Docker` image.
3. Publish your image to the image registry.
4. Test the published image by running it as a container.

Some popular `Docker` registries for publishing images are:
1. https://hub.docker.com/
2. https://quay.io/
3. https://cloud.google.com/container-registry
4. https://aws.amazon.com/ecr/

---
**Permission Issues with WSL**

When you run `Docker` commands from within `WSL`, you may encounter `access denied` errors. If so, the issue often originates from how `Docker` manages credentials across different environments, particularly related to the `config.json` file and `credsStore` property. To fix the issue, locate the `config.json` file, which is usually found in the `~/.docker/` directory of the `WSL` distribution. Next, open the file with a text editor and look for a line containing the property `"credsStore":`. If it's set to `"desktop.exe"` or `"wincred"`, it indicates that `Docker` is trying to use `Windows-based` credentials, which can cause issues in `WSL`. The simplest solution is to delete the entire line containing `"credsStore": "desktop.exe"` or rename `credsStore` to `credStore` (note the missing `**s**`). This forces `Docker` to rely on the default credentials within WSL.

---

### Useful Commands
---
**Note**

The `$()` syntax in `Bash` is used for command substitution. It allows you to execute a command and substitute its output in place of the `$()` expression. This is useful when you want to use the output of a command as an argument for another command or assign it to a variable.

---

#### Version
Display the current version of Docker.
```
$ docker --version
```

#### Private registry
Before you can push to your registry, you must first login. If the password or username contains special characters, you need to use quotes (""); otherwise, you can omit them.
```
$ docker login docker.io --username <user-name> --password "<password>"
```

The basic syntax of the command follows:
1. The image will be built on a `linux/amd64` platform (my machine), but the image will be run on a `linux/arm64` platform (Oracle). Create an image for each platform.
2. Tagging (`--tag or -t`) the image requires the following: image registry (`docker.io` is the default), the user's or organization's name, image repository, and the tag (`latest` is the default).
3. By default, `Docker` assumes that the `Dockerfile` is named `Dockerfile` and is located in the build context's root. If the `Dockerfile` has a different name or is located in a different directory, the `--file or -f` option can specify its path.
4. The final `..` in the command provides the path to the [build context](https://docs.docker.com/build/concepts/context/#what-is-a-build-context). At this location, the builder will find the `Dockerfile` and other referenced files.
```
$ docker build --platform linux/amd64,linux/arm64 --tag <registry>/<user's name>/<repo>:<tag> --file ../Dockerfile-prod ..
```

Publish the image to the registry.
```
$ docker push <registry>/<user's name>/<repo>:<tag>
```

Pull a single image from a repository.
```
$ docker pull <registry>/<user's name>/<repo>:<tag>
```

Pull all images from a repository; use the `--all-tags or -a` option.
```
docker pull --all-tags <registry>/<user's name>/<repo>
```

#### Listing images
List the images.
```
$ docker images
$ docker image ls
```

#### Removing images
Delete one or more images.
```
$ docker rmi <image-id> <image-id>
```

Remove all dangling images. If `-a` is specified, also remove all images not referenced by any container. To avoid confirmation prompts, specify the `-f` option.
```
$ docker image prune --all --force
```

Remove all images **older** than 6 months => 4320h = 24 hour/day * 30 days/month * 6 months. To avoid confirmation prompts, specify the `--force or -f` option.
```
$ docker image prune --all --filter "until=4320h"
$ docker image prune --all --force --filter "until=4320h"
```

Remove all images. The `docker images -q` command lists only the image ids.
```
$ docker rmi $(docker images -q)
```

Remove all dangling images. Docker images consist of multiple layers. Dangling images are layers that have no relationship to any tagged images. They no longer serve a purpose and consume disk space.
```
$ docker rmi $(docker images -qf "dangling=true")
```

#### Running containers
Run a command in a new container; pull the image, if needed, and start the container. The `--rm` option automatically removes the container and its associated anonymous volumes when it exits.
```
$ docker run --rm <user's name>/<repo>:<tag> ls -al
```

#### Starting containers
Run a command in a new container; pull the image, if needed, and start the container. The options are:<br>
`--detach or -d` runs container in the background and displays the container id.<br>
`--name` assigns a name to the container.
```
$ docker run -d --name finances <user's name>/<repo>:<tag>
```

#### Executing commands
Execute a command in a running container.
(Alpine images provide the Almquist shell (ash or sh) from BusyBox.)
```
$ docker exec -it <container-id-or-name> ash
$ docker exec -it <container-id-or-name> ls -al
$ docker exec -it <container-id-or-name> sh -c "echo a && echo b"
```

#### List containers
List all running containers.
```
$ docker ps
```

List all containers.
```
$ docker ps -a
```

#### Stopping containers
Stop one or more running containers.
```
$ docker stop <container-id> <container-id>
```

#### Killing containers
Kill one or more running containers.
```
$ docker kill <container-id-or-name> <container-id-or-name>
```

Kill a container and send a custom signal.
```
$ docker kill --signal=SIGKILL <container-id-or-name>
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

##### CPU Profiling
When it is activated, the application asks the OS to interrupt it every 10ms (default). When the application is interrupted, it suspends the current activity and transfers the execution to the profiler. The profiler collects execution statistics, and then it transfers execution back to the application.

To active the CPU profiling, you access the `debug/pprof/profile` endpoint. Accessing this endpoint will execute CPU profiling for 30 seconds by default. For 30 seconds, the application is interrupted every 10ms. If you need to customize the duration for which a CPU profile is collected, you can append the parameter `?seconds=<number>` to the `pprof` command; e.g., to collect a 5-minute CPU profile, you would use `debug/pprof/profile?seconds=300`.

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
This will set the `BROWSER` environment variable to the path of your desired browser. To make this [change permanent](#shell), you need to add this command to your `.bashrc` file or another relevant startup script.

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
$ go tool pprof -http={url1}:{port1} http://{url2}:{port2}/debug/pprof/{endpoint}?seconds={x}
```
In the preceding command, **url1:port1** will be used by the server displaying the profiling statistics whereas **url2:port2** will be used by the application being profiled.

### Compile and Run the App
---
**Note**

To see all environment variables supported by the app, see [config.go](./src/config/config.go). To run the app in a `K8s` environment, set the environment variable `K8S` to true.

---

A short description of the compiler switches being use follows:

`-a`: Force rebuilding of packages that are already up-to-date.

`-o`: The `go build` command specifies the output file name for the compiled binary. By default, without the `-o` flag, the output executable will have the same name as the source file (or the directory name if building a package) and no extension (or source_file.exe on `Windows`). The `-o` flag allows the customization of the name and optionally the location of the output file. The `-o` flag only affects the name of the final executable; it does not change the package name or any other aspects of the Go code.

`-ldflags="-s -w"`: In Go, the `go build` command allows the use of `ldflags` to modify the behavior of the linker. The `-s` and `-w` flags are commonly used with `ldflags` to reduce the size of the final executable.<br>
`-s`: This flag strips the symbol table from the executable. The symbol table is used for debugging and contains information about function and variable names. Removing it reduces the binary size but makes debugging more difficult.<br>
`-w`: This flag strips the DWARF debugging information from the executable. DWARF is a standardized debugging data format. Removing it further reduces the binary size but also hinders debugging capabilities.

`-installsuffix cgo`: The `go build` command uses this flag when the environment variable `CGO_ENABLED=1` thereby building and installing `Go` packages that utilize `cgo`, which enables `Go` code to interface with `C` code. The flag modifies the install path of the compiled package by adding *cgo* as a suffix to the installation directory. This ensures that `cgo-enabled packages` are installed separately from `regular Go packages` preventing potential conflicts between `cgo` and `non-cgo` builds.<br>
`CGO_ENABLED`: It is an environment variable that controls whether the `Go` tool enables the `cgo` tool. `Cgo` enables `Go` programs to call `C` code. `CGO_ENABLED` can have one of two values:<br>
`0`: Disable `cgo`. `Go` code **cannot call** `C` functions.<br>
`1`: Enable `cgo`. `Go` code **can call** `C` functions.<br>
If `cgo` is disabled, `Go` programs cannot use `C` libraries thereby limiting the functionality of some programs, especially those that need to interact with system-level APIs or existing `C` codebases.

Compile and run the app as a standalone HTTP server (default) on port 8080 (default).
```
$ CGO_ENABLED=0 go build -o finance && ./finance

or (the application does not require cgo to be enabled)

$ CGO_ENABLED=0 go build -o finance && HTTP=true HTTP_PORT=8080 ./finance

or (once the application does not require debugging information)

$ CGO_ENABLED=0 go build -a -ldflags="-s -w" -o finance && HTTP=true HTTP_PORT=8080 ./finance

or (same as above, but with cgo enabled)

$ CGO_ENABLED=1 go build -a -ldflags="-s -w" -installsuffix cgo -o finance && HTTP=true HTTP_PORT=8080 ./finance

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
To provision or deprovision resources required by the application, you use the script `iac.sh`. For a brief description of its usage, invoke the script with the option `-h` or `--help`.
```
$ ./iac.sh -h
```

The command below will deploy the application the following options:
reverse proxy: `Traefik`
deployment type: persistent disk
application version: 1.31.2
pprof: disable
```
$ ./iac.sh deploy -rp "true" -dt "persistent-disk" -av "1.31.2"
```

To deprovision the resources, use the command below.
```
$ ./iac.sh destroy
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
