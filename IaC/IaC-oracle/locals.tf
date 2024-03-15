locals {
  # Let's load the active nodes from the Node Pool; the NLB has to point to active nodes.
  active_nodes = (
    [for node in data.oci_containerengine_node_pool.k8s-node-port.nodes :
     node if node.state == "ACTIVE"]
  )
}
