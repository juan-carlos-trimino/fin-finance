/***
Define input variables to the module.
***/
variable compartment_id {
  type = string
}
variable shape {
  type = string
  default = "flexible"
}
variable public_subnet_id {
  type = string
}
variable is_private {
  type = bool
  default = false
}
variable vcn_id {
  type = string
}
variable target_ids {
  type = list
}
variable shape_details_minimum_bandwidth_in_mbps {
  type = number
  default = 10
}
variable shape_details_maximum_bandwidth_in_mbps {
  type = number
  default = 100
}

######################
# Load Balancer (LB) #
###################################################################################################
# The load balancer service provides a reverse proxy solution that hides the IP of the client     #
# from backend application server and vice versa. It is capable of performing advanced layer 7    #
# (HTTP/HTTPS), layer 4 (TCP) load balancing and SSL offloading.                                  #
#                                                                                                 #
# Best for: Websites, mobile apps, SSL termination, and advanced HTTP handling.                   #
###################################################################################################

# https://registry.terraform.io/providers/oracle/oci/latest/docs/resources/load_balancer_load_balancer
resource "oci_load_balancer_load_balancer" "lb" {
  display_name = "LB-Layer-7"
  compartment_id = var.compartment_id
  ip_mode = "IPV4"
  # If "true", the service assigns a private IP address to the load balancer.
  # If "false", the service assigns a public IP address to the load balancer.
  is_private = var.is_private
  shape = var.shape
  # The OCID of the subnet where you want to create this load balancer.
  subnet_ids = [
    var.public_subnet_id
  ]
  # The configuration details to create load balancer using Flexible shape. This is required only
  # if shapeName is Flexible.
  dynamic "shape_details" {
    for_each = var.shape == "flexible" ? [1] : []
    content {
      minimum_bandwidth_in_mbps = var.shape_details_minimum_bandwidth_in_mbps
      maximum_bandwidth_in_mbps = var.shape_details_maximum_bandwidth_in_mbps
    }
  }
  network_security_group_ids = [
    oci_core_network_security_group.public_lb_nsg.id,
    oci_core_network_security_group.private_lb_nsg.id
  ]
}

# HTTP
# https://registry.terraform.io/providers/oracle/oci/latest/docs/resources/load_balancer_listener
resource "oci_load_balancer_listener" "lb-listener-http" {
  name = "LB-HTTP-Listener"
  port = 80
  protocol = "HTTP"
  default_backend_set_name = oci_load_balancer_backend_set.lb-backend-set.name
  load_balancer_id = oci_load_balancer_load_balancer.lb.id
  connection_configuration {
  #   backend_tcp_proxy_protocol_version = "0"
    idle_timeout_in_seconds = "60"
  }
}

# HTTPS
resource "oci_load_balancer_listener" "lb-listener-https" {
  name = "LB-HTTPS-Listener"
  port = 443
  protocol = "HTTP"
  default_backend_set_name = oci_load_balancer_backend_set.lb-backend-set.name
  load_balancer_id = oci_load_balancer_load_balancer.lb.id
  connection_configuration {
  #   backend_tcp_proxy_protocol_version = "0"
    idle_timeout_in_seconds = "60"
  }
}

# https://registry.terraform.io/providers/oracle/oci/latest/docs/resources/load_balancer_backend_set
resource "oci_load_balancer_backend_set" "lb-backend-set" {
  name = "lb-backend-set"
  health_checker {
    port = 443
    protocol = "HTTP"
    return_code = 200
    # url_path = "/api/health"
    url_path = "/"
    # interval_ms = 10000
    # response_body_regex = ".*"
    # retries = 3
    # timeout_in_millis = 3000
  }
  load_balancer_id = oci_load_balancer_load_balancer.lb.id
  policy = "ROUND_ROBIN"
  # policy = "IP_HASH"
}

# https://registry.terraform.io/providers/oracle/oci/latest/docs/resources/load_balancer_backend
resource "oci_load_balancer_backend" "lb-backend" {
  count = length(var.target_ids)
  port = 80
  backup = false
  drain = false
  offline = false
  weight = 1
  load_balancer_id = oci_load_balancer_load_balancer.lb.id
  backendset_name = oci_load_balancer_backend_set.lb-backend-set.name
  ip_address = var.target_ids[count.index].private_ip
}

resource "oci_core_network_security_group" "public_lb_nsg" {
  display_name = "Public LB NSG"
  compartment_id = var.compartment_id
  vcn_id = var.vcn_id
}

resource "oci_core_network_security_group_security_rule" "allow_http_from_all" {
  description = "Allow HTTP from all"
  stateless = false
  network_security_group_id = oci_core_network_security_group.public_lb_nsg.id
  direction = "INGRESS"
  protocol = 6  # TCP
  source = "0.0.0.0/0"
  source_type = "CIDR_BLOCK"
  tcp_options {
    destination_port_range {
      max = 80
      min = 80
    }
  }
}

resource "oci_core_network_security_group" "private_lb_nsg" {
  display_name = "Private LB NSG"
  compartment_id = var.compartment_id
  vcn_id = var.vcn_id
}

resource "oci_core_network_security_group_security_rule" "allow_https_from_all" {
  description = "Allow HTTPS from all"
  stateless = false
  network_security_group_id = oci_core_network_security_group.private_lb_nsg.id
  direction = "INGRESS"
  protocol = 6  # TCP
  source = "0.0.0.0/0"
  source_type = "CIDR_BLOCK"
  tcp_options {
    destination_port_range {
      max = 443
      min = 443
    }
  }
}

output "lb_public_ip" {
  value = ([
    for ip in oci_load_balancer_load_balancer.lb.ip_address_details :
      ip.ip_address
  ])
}
