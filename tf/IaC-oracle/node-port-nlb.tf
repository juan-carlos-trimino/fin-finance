########################################
# Network Load Balancer (NLB)          #
# (It routes traffic into the cluster) #
########################################

data "oci_containerengine_node_pool" "k8s-node-port" {
  node_pool_id = module.arm64-node-pool.node-pool-id
}

locals {
  # Let's load the active nodes from the Node Pool; the NLB has to point to active nodes.
  active_nodes = (
    [for node in data.oci_containerengine_node_pool.k8s-node-port.nodes :
     node if node.state == "ACTIVE"]
  )
}

resource "oci_network_load_balancer_network_load_balancer" "node-port-nlb" {
  compartment_id = oci_identity_compartment.fin-compartment.id
  display_name = "node-port-nlb"
  subnet_id = module.public-subnet.subnet-id
  is_private = false
  is_preserve_source_destination = false
}

resource "oci_network_load_balancer_backend_set" "node-port-nlb-backend-set" {
  name = "node-port-nlb-backend-set"
  health_checker {
    port = 10256
    protocol = "TCP"
  }
  network_load_balancer_id = oci_network_load_balancer_network_load_balancer.node-port-nlb.id
  policy = "FIVE_TUPLE"
  is_preserve_source = false
}

resource "oci_network_load_balancer_backend" "node-port-nlb-backend" {
  count = length(local.active_nodes)
  backend_set_name = oci_network_load_balancer_backend_set.node-port-nlb-backend-set.name
  network_load_balancer_id = oci_network_load_balancer_network_load_balancer.node-port-nlb.id
  port = 31600
  target_id = local.active_nodes[count.index].id
}

resource "oci_network_load_balancer_listener" "node-port-nlb-listener" {
  name = "node-port-nlb-listener"
  default_backend_set_name = oci_network_load_balancer_backend_set.node-port-nlb-backend-set.name
  network_load_balancer_id = oci_network_load_balancer_network_load_balancer.node-port-nlb.id
  port = "80"
  protocol = "TCP"
}
