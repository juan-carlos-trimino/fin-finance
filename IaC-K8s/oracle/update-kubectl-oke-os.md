# Updating `kubectl` to Match Oracle Kubernetes Engine (OKE)

## Overview

When managing an Oracle Kubernetes Engine (OKE) cluster from any local machine or CI/CD environment (Linux, WSL2, macOS, or Windows), you may encounter a minor version skew warning when running `kubectl` commands:

```text
Client Version: v1.31.0
Server Version: v1.36.1
WARNING: version difference between client (1.31) and server (1.36) exceeds the supported minor version skew of +/-1
```

*(Note: While the installation steps below are tailored for WSL2/Linux, this version skew behavior applies universally across all operating systems whenever the `kubectl` binary falls behind or gets ahead of the cluster control plane.)*

---

## Why This Update Is Required

### 1. Kubernetes Version Skew Policy
Kubernetes enforces a strict support policy between the client (`kubectl`) and the cluster control plane (`kube-apiserver`):
* **`kubectl` is officially supported within one minor version (`+/- 1`)** of the cluster's control plane version.
* For an OKE server running **`v1.36.1`**, supported `kubectl` versions are **`v1.35.x`**, **`v1.36.x`**, and **`v1.37.x`**.

### 2. Risks of Running an Outdated Client
* **API & Feature Mismatches:** Newer server-side API fields, flags, or resource schemas introduced between `v1.31` and `v1.36` may not be recognized or properly validated by an older `kubectl` binary.
* **Subcommand & Output Discrepancies:** Built-in helpers (such as `kubectl debug`, `kubectl create`, or `kubectl explain`) rely on client-side logic that can fall out of sync with the cluster.
* **Operational Noise:** Repeated version skew warnings clutter terminal output and CI/CD or automation logs.

---

## How to Solve It (WSL2 / Linux)

### Step 1: Identify the Current `kubectl` Path
Verify where `kubectl` is currently installed so you overwrite the active binary rather than creating a duplicate in your `PATH`:

```bash
which kubectl
```

*(Typically `/usr/local/bin/kubectl`.)*

### Step 2: Clean Up Any Partial Downloads
If an earlier download attempt saved an HTML error page or incomplete file named `kubectl` in your working directory, remove it first:

```bash
rm -f kubectl kubectl.sha256
```

### Step 3: Download, Verify, and Install `kubectl` `v1.36.1`
Construct the download URL using shell variables to prevent browser extensions or clipboard tools from truncating or rewriting the URL during copy-paste:

```bash
VER="v1.36.1"
HOST="dl.k8s.io"
PATH_URL="release/${VER}/bin/linux/amd64"

curl -L "https://${HOST}/${PATH_URL}/kubectl" -o kubectl && \
curl -L "https://${HOST}/${PATH_URL}/kubectl.sha256" -o kubectl.sha256 && \
echo "$(cat kubectl.sha256)  kubectl" | sha256sum --check && \
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && \
rm kubectl kubectl.sha256
```

> **Note:** If you are running WSL2 or Linux on an ARM64 processor, change `amd64` to `arm64` in `PATH_URL`.

### Step 4: Verify the Installation
Confirm that both the client and server versions match and the skew warning is resolved:

```bash
kubectl version
```

Expected output:
```text
Client Version: v1.36.1
Server Version: v1.36.1
```

### Step 5: Check Worker Node Versions (Recommended)
In OKE, control plane upgrades and worker node pool upgrades happen independently. Verify that your worker nodes (`kubelet`) are also aligned with `v1.36.1`:

```bash
kubectl get nodes
```
