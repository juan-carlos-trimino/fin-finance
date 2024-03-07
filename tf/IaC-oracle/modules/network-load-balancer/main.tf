/***
Define input variables to the module.
***/
variable compartment_id {
  type = string
}
variable subnet_id {
  type = string
}
# variable node_pool_id {
#   type = string
# }
variable nodes {
  type = number
}
variable nlb_node_port {
  type = number
}

variable target_ids {
  type = list
}

# variable active_nodes {
#   type = list
# }


########################################
# Network Load Balancer (NLB)          #
# (It routes traffic into the cluster) #
###################################################################################################
# The network load balancer service provides a pass-through (non-proxy solution) that is capable  #
# of preserving the client header (source and destination IP). It is built for speed, optimized   #
# for long running connections, high throughput and low latency.                                  #
#                                                                                                 #
# Best for: Scaling network virtual appliances such as firewalls, real-time streaming, long       #
# running connections, Voice over IP (VoIP), Internet of Things (IoT), and trading platforms.     #
###################################################################################################

# data "oci_containerengine_node_pool" "k8s-node-port" {
#   node_pool_id = var.node_pool_id
# }

# locals {
#   # Let's load the active nodes from the Node Pool; the NLB has to point to active nodes.
#   active_nodes = (
#     [for node in data.oci_containerengine_node_pool.k8s-node-port.nodes :
#      node if node.state == "ACTIVE"]
#   )
# }


# data "oci_core_instance_pool_instances" "instance_pool_instances" {
#   compartment_id = var.compartment_id
#   instance_pool_id = var.node_pool_id
# }

# data "oci_core_instance" "core_instances_ips" {
#   count = var.nodes
#   instance_id = data.oci_core_instance_pool_instances.instance_pool_instances.instances[count.index].id
# }




resource "oci_network_load_balancer_network_load_balancer" "node-port-nlb" {
  display_name = "node-port-nlb"
  compartment_id = var.compartment_id
  subnet_id = var.subnet_id
  is_private = false
  is_preserve_source_destination = false
}

# https://registry.terraform.io/providers/oracle/oci/latest/docs/resources/network_load_balancer_backend_set
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
  count = length(var.target_ids) # var.nodes #length(local.active_nodes)
  # for_each = toset([var.target_ids])
  backend_set_name = oci_network_load_balancer_backend_set.node-port-nlb-backend-set.name
  network_load_balancer_id = oci_network_load_balancer_network_load_balancer.node-port-nlb.id
  port = var.nlb_node_port
  # target_id = local.active_nodes[count.index].id
//  target_id = data.oci_core_instance_pool_instances.instance_pool_instances.instances[count.index].id
  target_id = var.target_ids[count.index].id
  # target_id = each.key
}

resource "oci_network_load_balancer_listener" "node-port-nlb-listener" {
  name = "node-port-nlb-listener"
  default_backend_set_name = oci_network_load_balancer_backend_set.node-port-nlb-backend-set.name
  network_load_balancer_id = oci_network_load_balancer_network_load_balancer.node-port-nlb.id
  port = "80"
  protocol = "TCP"
}

output "node_port_nlb_public_ip" {
  value = ([
    for ip in oci_network_load_balancer_network_load_balancer.node-port-nlb.ip_addresses :
      ip if ip.is_public == true
  ])
}
