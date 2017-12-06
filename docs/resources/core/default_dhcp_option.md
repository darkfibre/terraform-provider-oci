# oci\_core\_default\_dhcp\_options

[DhcpOptions Reference][82b94f0a]

  [82b94f0a]: https://docs.us-phoenix-1.oraclecloud.com/api/#/en/iaas/20160918/DhcpOptions/ "DhcpOptionsReference"

Configures a VCN's default DHCP options resource.

For more information, see
[DNS in Your Virtual Cloud Network](https://docs.us-phoenix-1.oraclecloud.com/Content/Network/Concepts/dns.htm)

For more information on default resources, see [Managing Default VCN Resources](https://github.com/oracle/terraform-provider-oci/blob/master/docs/Managing%20Default%20Resources.md)
## Example Usage

#### VCN Local with Internet
```
resource "oci_core_default_dhcp_options" "dhcp-options1" {
  manage_default_resource_id = "id of VCN's default DHCP options"
  display_name = "dhcp-options1"

  // required
  options {
    type = "DomainNameServer"
    server_type = "VcnLocalPlusInternet"
  }

  // optional
  options {
    type = "SearchDomain"
    search_domain_names = [ "test.com" ]
  }
}
```

#### Custom DNS Server

```
resource "oci_core_default_dhcp_options" "dhcp-options2" {
  manage_default_resource_id = "id of VCN's default security list"
  display_name = "dhcp-options3"

  // required
  options {
    type = "DomainNameServer"
    server_type = "CustomDnsServer"
    custom_dns_servers = [ "192.168.0.2", "192.168.0.11", "192.168.0.19" ]
  }

  // optional
  options {
    type = "SearchDomain"
    search_domain_names = [ "test.com" ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `manage_default_resource_id` - (Required) The OCID of a [default DHCP option resource](https://github.com/oracle/terraform-provider-oci/blob/master/docs/Managing%20Default%20Resources.md) to manage.
* `display_name` - (Optional) A user-friendly name. Does not have to be unique, and it's changeable. Avoid entering confidential information.
* `options` - (Required) A set of [DHCP Options](https://docs.us-phoenix-1.oraclecloud.com/api/#/en/iaas/20160918/DhcpDnsOption/).

## Attributes Reference
* `display_name` - A user-friendly name. Does not have to be unique, and it's changeable. Avoid entering confidential information.
* `id` - Oracle ID (OCID) for the set of DHCP options.
* `state` - The DRG's current state. Allowed values are: [PROVISIONING, AVAILABLE, TERMINATING, TERMINATED]
* `options` - The collection of individual DHCP options.
* `time_created` - The date and time the set of DHCP options was created, in the format defined by RFC3339. Example: `2016-08-25T21:10:29.600Z`.
