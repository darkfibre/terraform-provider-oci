// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/oracle/bmcs-go-sdk"

	"fmt"
	"time"

	"github.com/oracle/terraform-provider-oci/crud"
)

var routeRules = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"cidr_block": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"network_entity_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
	},
}

func RouteTableResource() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: crud.DefaultTimeout,
		Create:   createRouteTable,
		Read:     readRouteTable,
		Update:   updateRouteTable,
		Delete:   deleteRouteTable,
		Schema: map[string]*schema.Schema{
			"compartment_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: crud.DefaultResourceSuppressDiff,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"compartment_id", "vcn_id"},
				ForceNew:      true,
			},
			"default_route_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     routeRules,
			},
			"route_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     routeRules,
			},
			"time_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcn_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: crud.DefaultResourceSuppressDiff,
			},
		},
	}
}

func createRouteTable(d *schema.ResourceData, m interface{}) (e error) {
	client := m.(*OracleClients)
	crd := &RouteTableResourceCrud{}
	crd.D = d
	crd.Client = client.client
	return crud.CreateResource(d, crd)
}

func readRouteTable(d *schema.ResourceData, m interface{}) (e error) {
	client := m.(*OracleClients)
	crd := &RouteTableResourceCrud{}
	crd.D = d
	crd.Client = client.client
	return crud.ReadResource(crd)
}

func updateRouteTable(d *schema.ResourceData, m interface{}) (e error) {
	client := m.(*OracleClients)
	crd := &RouteTableResourceCrud{}
	crd.D = d
	crd.Client = client.client
	return crud.UpdateResource(d, crd)
}

func deleteRouteTable(d *schema.ResourceData, m interface{}) (e error) {
	client := m.(*OracleClients)
	crd := &RouteTableResourceCrud{}
	crd.D = d
	crd.Client = client.clientWithoutNotFoundRetries
	return crud.DeleteResource(d, crd)
}

type RouteTableResourceCrud struct {
	crud.BaseCrud
	Res *baremetal.RouteTable
}

func (s *RouteTableResourceCrud) ID() string {
	return s.Res.ID
}

func (s *RouteTableResourceCrud) CreatedPending() []string {
	return []string{baremetal.ResourceProvisioning}
}

func (s *RouteTableResourceCrud) CreatedTarget() []string {
	return []string{baremetal.ResourceAvailable}
}

func (s *RouteTableResourceCrud) DeletedPending() []string {
	return []string{baremetal.ResourceTerminating}
}

func (s *RouteTableResourceCrud) DeletedTarget() []string {
	return []string{baremetal.ResourceTerminated}
}

func (s *RouteTableResourceCrud) State() string {
	return s.Res.State
}

func (s *RouteTableResourceCrud) Create() (e error) {
	// If we are creating a default resource, then don't have to
	// actually create it. Just set the ID and update it.
	if defaultId := s.D.Get("default_id").(string); defaultId != "" {
		var res *baremetal.RouteTable
		if res, e = s.Client.GetRouteTable(defaultId); e != nil {
			return
		}
		s.D.Set("default_route_rules", routeRulesToMapArray(res.RouteRules))

		s.D.SetId(defaultId)
		e = s.Update()
		return
	}

	compartmentID := s.D.Get("compartment_id").(string)
	vcnID := s.D.Get("vcn_id").(string)

	opts := &baremetal.CreateOptions{}
	opts.DisplayName = s.D.Get("display_name").(string)

	rr, e := s.buildRouteRules()

	if e != nil {
		return e
	}
	s.Res, e = s.Client.CreateRouteTable(compartmentID, vcnID, rr, opts)
	s.D.Set("default_route_rules", nil)

	return
}

func (s *RouteTableResourceCrud) Get() (e error) {
	res, e := s.Client.GetRouteTable(s.D.Id())
	if e == nil {
		s.Res = res

		// If this is a default resource that we removed earlier, then
		// we need to assume that the parent resource will remove it
		// and notify terraform of it. Otherwise, terraform will
		// see that the resource is still available and error out
		if s.D.Get("default_id") != "" &&
			s.D.Get("state") == baremetal.ResourceTerminated {
			s.Res.State = baremetal.ResourceTerminated
		}
	}
	return
}

func (s *RouteTableResourceCrud) Update() (e error) {
	opts := &baremetal.UpdateRouteTableOptions{}

	if displayName, ok := s.D.GetOk("display_name"); ok {
		opts.DisplayName = displayName.(string)
	}

	opts.RouteRules, e = s.buildRouteRules()

	if e != nil {
		return e
	}

	s.Res, e = s.Client.UpdateRouteTable(s.D.Id(), opts)
	return
}

func routeRulesToMapArray(rules []baremetal.RouteRule) (res []map[string]interface{}) {
	for _, val := range rules {
		rule := map[string]interface{}{
			"cidr_block":        val.CidrBlock,
			"network_entity_id": val.NetworkEntityID,
		}
		res = append(res, rule)
	}

	return
}

func (s *RouteTableResourceCrud) SetData() {
	s.D.Set("compartment_id", s.Res.CompartmentID)
	s.D.Set("display_name", s.Res.DisplayName)

	s.D.Set("route_rules", routeRulesToMapArray(s.Res.RouteRules))

	s.D.Set("time_modified", s.Res.TimeModified.String())
	s.D.Set("state", s.Res.State)
	s.D.Set("time_created", s.Res.TimeCreated.String())
}

func (s *RouteTableResourceCrud) Delete() (e error) {
	if s.D.Get("default_id") != "" {
		// We can't actually delete a default resource.
		// Instead, revert it to default settings and mark it as terminated
		s.D.Set("route_rules", s.D.Get("default_route_rules"))
		e = s.Update()

		s.D.Set("state", baremetal.ResourceTerminated)
		return
	}

	return s.Client.DeleteRouteTable(s.D.Id(), nil)
}

func (s *RouteTableResourceCrud) ExtraWaitPostCreateDelete() time.Duration {
	return time.Duration(15 * time.Second)
}

func (s *RouteTableResourceCrud) buildRouteRules() (routeRules []baremetal.RouteRule, e error) {
	routeRules = []baremetal.RouteRule{}
	for _, val := range s.D.Get("route_rules").([]interface{}) {

		if val == nil {
			return nil, fmt.Errorf("Empty route_rules are not permitted. Instead, the route_rules block may be omitted entirely.")
		}

		data := val.(map[string]interface{})
		routeRule := baremetal.RouteRule{
			CidrBlock:       data["cidr_block"].(string),
			NetworkEntityID: data["network_entity_id"].(string),
		}
		routeRules = append(routeRules, routeRule)
	}
	return
}
