variable "tenancy_ocid" {}
variable "user_ocid" {}
variable "fingerprint" {}
variable "private_key_path" {}
variable "compartment_ocid" {}
variable "region" {}

provider "oci" {
  tenancy_ocid = "${var.tenancy_ocid}"
  user_ocid = "${var.user_ocid}"
  fingerprint = "${var.fingerprint}"
  private_key_path = "${var.private_key_path}"
  region = "${var.region}"
}

resource "oci_core_virtual_network" "vcn1" {
  cidr_block = "10.0.0.0/16"
  dns_label = "vcn1"
  compartment_id = "${var.compartment_ocid}"
  display_name = "vcn1"
}

resource "oci_core_internet_gateway" "ig1" {
  compartment_id = "${var.compartment_ocid}"
  display_name = "ig1"
  vcn_id = "${oci_core_virtual_network.vcn1.id}"
}

resource "oci_core_route_table" "default-route-table" {
  default_id = "${oci_core_virtual_network.vcn1.default_route_table_id}"
  display_name = "default-route-table"
  route_rules {
    cidr_block = "0.0.0.0/0"
    network_entity_id = "${oci_core_internet_gateway.ig1.id}"
  }
}

resource "oci_core_route_table" "route-table1" {
  compartment_id = "${var.compartment_ocid}"
  vcn_id = "${oci_core_virtual_network.vcn1.id}"
  display_name = "route-table1"
  route_rules {
    cidr_block = "0.0.0.0/0"
    network_entity_id = "${oci_core_internet_gateway.ig1.id}"
  }
}

resource "oci_core_dhcp_options" "default-dhcp-options" {
  default_id = "${oci_core_virtual_network.vcn1.default_dhcp_options_id}"
  display_name = "default-dhcp-options"

  // required
  options {
    type = "DomainNameServer"
    server_type = "VcnLocalPlusInternet"
  }

  // optional
  options {
    type = "SearchDomain"
    search_domain_names = [ "abc.com" ]
  }
}

resource "oci_core_dhcp_options" "dhcp-options1" {
  compartment_id = "${var.compartment_ocid}"
  vcn_id = "${oci_core_virtual_network.vcn1.id}"
  display_name = "dhcp-options1"

  // required
  options {
    type = "DomainNameServer"
    server_type = "VcnLocalPlusInternet"
  }

  // optional
  options {
    type = "SearchDomain"
    search_domain_names = [ "test123.com" ]
  }
}

resource "oci_core_security_list" "default-security-list" {
  default_id = "${oci_core_virtual_network.vcn1.default_security_list_id}"
  display_name = "default-security-list"

  // allow outbound tcp traffic on all ports
  egress_security_rules {
    destination = "0.0.0.0/0"
    protocol = "6"
  }

  // allow outbound udp traffic on a port range
  egress_security_rules {
    destination = "0.0.0.0/0"
    protocol = "17" // udp
    stateless = true

    udp_options {
      "min" = 319
      "max" = 320
    }
  }

  // allow inbound ssh traffic
  ingress_security_rules {
    protocol = "6" // tcp
    source = "0.0.0.0/0"
    stateless = false

    tcp_options {
      "min" = 22
      "max" = 22
    }
  }

  // allow inbound icmp traffic of a specific type
  ingress_security_rules {
    protocol  = 1
    source    = "0.0.0.0/0"
    stateless = true

    icmp_options {
      "type" = 3
      "code" = 4
    }
  }
}
