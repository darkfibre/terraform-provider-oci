# oci\_core\_default\_route\_table

[RouteTable Reference][e98ebc48]

  [e98ebc48]: https://docs.us-phoenix-1.oraclecloud.com/api/#/en/iaas/20160918/RouteTable/ "RouteTableReference"

Configure a VCN's default route table resource.

For more information on default resources, see [Managing Default VCN Resources](https://github.com/oracle/terraform-provider-oci/blob/master/docs/Managing%20Default%20Resources.md)
## Example Usage

```
resource "oci_core_default_route_table" "t" {
    manage_default_resource_id = "id of VCN's default security list"
    display_name = "display_name"
    route_rules {
        cidr_block = "cidr_block"
        network_entity_id = "network_entity_id"
    }
    route_rules {
        cidr_block = "cidr_block"
        network_entity_id = "network_entity_id"
    }
    vcn_id = "vcn_id"
}
```

## Argument Reference

The following arguments are supported:

* `manage_default_resource_id` - (Required) The OCID of a [default route table resource](https://github.com/oracle/terraform-provider-oci/blob/master/docs/Managing%20Default%20Resources.md) to manage.
* `display_name` - (Optional) A user-friendly name. Does not have to be unique, and it's changeable. Avoid entering confidential information.
* `route_rules` - (Required) The collection of rules for routing destination IPs to network devices.

## Attributes reference

* `display_name` - A user-friendly name. Does not have to be unique, and it's changeable. Avoid entering confidential information.
* `id` - The route table's Oracle Cloud ID (OCID).
* `state` - The route table's current state. Allowed values are: [PROVISIONING, AVAILABLE, TERMINATING, TERMINATED]
* `route_rules` - The collection of rules for routing destination IPs to network devices.
* `time_created` - The date and time the route table was created, in the format defined by RFC3339. Example: `2016-08-25T21:10:29.600Z`.
