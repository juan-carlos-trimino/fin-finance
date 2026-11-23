# OKE Basic Cluster Upgrade Guide (Worker Nodes & Control Plane via Terraform and OCI CLI)

This guide covers upgrading an Oracle Kubernetes Engine (OKE) **Basic** cluster when a control plane upgrade is blocked by Kubernetes version skew rules (`409-IncorrectState`), including dynamically resolving ARM (`aarch64`) OKE worker images in Terraform and manually cycling worker nodes without incurring Enhanced cluster fees.

---

## 1. Understanding the Version Skew Constraint

Kubernetes requires worker nodes (`kubelet`) to be within **3 minor versions** of the control plane (`kube-apiserver`).

- **Blocked Upgrade Example:** Control plane on `v1.36.1` upgrading to `v1.37.1` while worker nodes are on `v1.33.0` (`1.37 - 1.33 = 4` minor versions).
- **Resolution Order:**
  1. Update the worker node pool template (`kubernetes_version` and `image_id`) to `v1.36.1` in Terraform.
  2. Replace existing `v1.33.x` worker nodes one by one so they recreate on `v1.36.1`.
  3. Upgrade the control plane (`oci_containerengine_cluster`) to `v1.37.1` in Terraform.

---

## 2. Terraform Configuration for Dynamic ARM (`aarch64`) OKE Images

Standard compute image data sources (`oci_core_images`) return generic OS images without Kubernetes binaries and will fail with `400-InvalidParameter: Invalid nodeSourceDetails.imageId`. Always use `oci_containerengine_node_pool_option` to fetch official OKE worker images.

### Step 2.1: Root `data.tf`
Fetch all OKE node pool options and filter for ARM (`aarch64`) Oracle Linux 8 images matching `var.k8s_version`:

```hcl
data "oci_containerengine_node_pool_option" "oke-images" {
  node_pool_option_id = "all"
}

locals {
  # Filter for ARM (aarch64) Oracle Linux 8 OKE images matching k8s_version
  arm_oke_images = [
    for source in data.oci_containerengine_node_pool_option.oke-images.sources : source.image_id
    if length(regexall("Oracle-Linux-8\..*-aarch64-.*-OKE-${replace(var.k8s_version, "v", "")}", source.source_name)) > 0
    && length(regexall("GPU", source.source_name)) == 0
  ]
}
```

### Step 2.2: Root `main.tf`
Pass the resolved image OCID into the node pool module:

```hcl
module "arm64-node-pool" {
  source = "./modules/node"

  # ... existing module arguments ...
  node_image_id = local.arm_oke_images[0]
}
```

### Step 2.3: `modules/node/variables.tf`
Declare the input variable inside the node pool module:

```hcl
variable "node_image_id" {
  type        = string
  description = "The OCID of the OKE worker node image."
}
```

### Step 2.4: `modules/node/main.tf`
Use `var.node_image_id` in `node_source_details`, and keep `node_pool_cycling_details` commented out (as automated node cycling requires an Enhanced cluster):

```hcl
resource "oci_containerengine_node_pool" "node-pool" {
  name               = var.name
  cluster_id         = var.cluster_id
  compartment_id     = var.compartment_id
  kubernetes_version = var.k8s_version

  node_source_details {
    boot_volume_size_in_gbs = 50
    image_id                = var.node_image_id
    source_type             = "image"
  }

  # Enhanced cluster feature (keep commented out on Basic clusters)
  # node_pool_cycling_details {
  #   is_node_cycling_enabled = false
  #   maximum_surge           = 1
  #   maximum_unavailable     = 0
  # }
}
```

### Step 2.5: Apply the Node Pool Template Update
Run Terraform from your root directory to update the node pool template in OCI:

```bash
terraform plan
terraform apply -auto-approve
```

---

## 3. Replacing Worker Nodes One by One (Basic Cluster)

Because Basic clusters do not automatically cycle existing nodes when a node pool template is updated, replace the 4 old nodes sequentially so OKE recreates each one using the `v1.36.1` template.

### Step 3.1: List Current Worker Nodes and Their OCIDs
Run the following OCI CLI command using your node pool OCID (note the outer single quotes `'...'` around `--query` so Bash preserves the inner double quotes around hyphenated keys):

```bash
oci ce node-pool get   --node-pool-id ocid1.nodepool.oc1.iad.aaaaaaaaww2vjshfm2rvjpnizz7kaqppnhicgolkws3thb4ajnsscnytdmfq   --query 'data.nodes[].{name:name, ip:"private-ip", id:id, state:"lifecycle-state"}'   --output table
```

### Step 3.2: Delete One Old Node at a Time
Replace `<NODE_OCID>` with the OCID (`ocid1.instance...`) of a single `v1.33.x` node. Keep `--is-decrement-size false` so OKE immediately provisions a replacement node to maintain your pool size of 4:

```bash
oci ce node-pool delete-node   --node-pool-id ocid1.nodepool.oc1.iad.aaaaaaaaww2vjshfm2rvjpnizz7kaqppnhicgolkws3thb4ajnsscnytdmfq   --node-id <NODE_OCID>   --is-decrement-size false   --force
```

### Step 3.3: Verify the Replacement Node Is Ready
Monitor Kubernetes until the new node joins as `Ready` on `v1.36.1` (typically 3–5 minutes):

```bash
kubectl get nodes -w
```

### Step 3.4: Repeat for Remaining Old Nodes
Repeat **Step 3.2** and **Step 3.3** one node at a time until all 4 worker nodes show `Ready` and `v1.36.1`.

---

## 4. Upgrading the Control Plane to `v1.37.1`

Once all 4 worker nodes are confirmed on `v1.36.1`:

1. Update `k8s_version` (or the `kubernetes_version` attribute on `oci_containerengine_cluster`) to `"v1.37.1"` in Terraform.
   - *Note:* If `var.k8s_version` is shared between `oci_containerengine_cluster` and `oci_containerengine_node_pool`, verify via `terraform plan` whether you want the node pool template updated to `v1.37.1` at the same time or kept pinned to `v1.36.1`.
2. Run Terraform to upgrade the control plane:
   ```bash
   terraform plan
   terraform apply -auto-approve
   ```
3. Verify the cluster and nodes after completion:
   ```bash
   kubectl version
   kubectl get nodes
   ```
