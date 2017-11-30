// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/oracle/bmcs-go-sdk"

	"github.com/oracle/terraform-provider-oci/crud"
)

var dhcpOptions = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"type": {
			Type:     schema.TypeString,
			Required: true,
		},
		"custom_dns_servers": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"server_type": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"search_domain_names": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	},
}

func DHCPOptionsResource() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: crud.DefaultTimeout,
		Create:   createDHCPOptions,
		Read:     readDHCPOptions,
		Update:   updateDHCPOptions,
		Delete:   deleteDHCPOptions,
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
			"default_options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dhcpOptions,
			},
			"options": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     dhcpOptions,
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

func createDHCPOptions(d *schema.ResourceData, m interface{}) (e error) {
	client := m.(*OracleClients)
	crd := &DHCPOptionsResourceCrud{}
	crd.D = d
	crd.Client = client.client
	return crud.CreateResource(d, crd)
}

func readDHCPOptions(d *schema.ResourceData, m interface{}) (e error) {
	crd := &DHCPOptionsResourceCrud{}
	client := m.(*OracleClients)
	crd.D = d
	crd.Client = client.client
	return crud.ReadResource(crd)
}

func updateDHCPOptions(d *schema.ResourceData, m interface{}) (e error) {
	client := m.(*OracleClients)
	crd := &DHCPOptionsResourceCrud{}
	crd.D = d
	crd.Client = client.client
	return crud.UpdateResource(d, crd)
}

func deleteDHCPOptions(d *schema.ResourceData, m interface{}) (e error) {
	client := m.(*OracleClients)
	crd := &DHCPOptionsResourceCrud{}
	crd.D = d
	crd.Client = client.clientWithoutNotFoundRetries
	return crud.DeleteResource(d, crd)
}

type DHCPOptionsResourceCrud struct {
	crud.BaseCrud
	Res *baremetal.DHCPOptions
}

func (s *DHCPOptionsResourceCrud) ID() string {
	return s.Res.ID
}

func (s *DHCPOptionsResourceCrud) CreatedPending() []string {
	return []string{baremetal.ResourceProvisioning}
}

func (s *DHCPOptionsResourceCrud) CreatedTarget() []string {
	return []string{baremetal.ResourceAvailable}
}

func (s *DHCPOptionsResourceCrud) DeletedPending() []string {
	return []string{baremetal.ResourceTerminating}
}

func (s *DHCPOptionsResourceCrud) DeletedTarget() []string {
	return []string{baremetal.ResourceTerminated}
}

func (s *DHCPOptionsResourceCrud) State() string {
	return s.Res.State
}

func (s *DHCPOptionsResourceCrud) Create() (e error) {
	// If we are creating a default resource, then don't have to
	// actually create it. Just set the ID and update it.
	if defaultId := s.D.Get("default_id").(string); defaultId != "" {
		// We query the default DHCP options; this can be restored later
		// when we try to delete a default resource.
		var res *baremetal.DHCPOptions
		if res, e = s.Client.GetDHCPOptions(defaultId); e != nil {
			return
		}
		s.D.Set("default_options", optionsToMapArray(res.Options))

		s.D.SetId(defaultId)
		e = s.Update()
		return
	}

	compartmentID := s.D.Get("compartment_id").(string)
	vcnID := s.D.Get("vcn_id").(string)

	opts := &baremetal.CreateOptions{}
	opts.DisplayName = s.D.Get("display_name").(string)

	s.Res, e = s.Client.CreateDHCPOptions(compartmentID, vcnID, s.buildEntities(), opts)
	s.D.Set("default_options", nil)

	return
}

func (s *DHCPOptionsResourceCrud) Get() (e error) {
	res, e := s.Client.GetDHCPOptions(s.D.Id())
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

func (s *DHCPOptionsResourceCrud) Update() (e error) {
	opts := &baremetal.UpdateDHCPDNSOptions{}
	opts.Options = s.buildEntities()

	s.Res, e = s.Client.UpdateDHCPOptions(s.D.Id(), opts)
	return
}

func optionsToMapArray(options []baremetal.DHCPDNSOption) (res []map[string]interface{}) {
	for _, val := range options {
		entity := map[string]interface{}{
			"type":                val.Type,
			"custom_dns_servers":  val.CustomDNSServers,
			"server_type":         val.ServerType,
			"search_domain_names": val.SearchDomainNames,
		}
		res = append(res, entity)
	}

	return
}

func (s *DHCPOptionsResourceCrud) SetData() {
	s.D.Set("compartment_id", s.Res.CompartmentID)
	s.D.Set("display_name", s.Res.DisplayName)

	s.D.Set("options", optionsToMapArray(s.Res.Options))

	s.D.Set("state", s.Res.State)
	s.D.Set("time_created", s.Res.TimeCreated.String())
}

func (s *DHCPOptionsResourceCrud) Delete() (e error) {
	if s.D.Get("default_id") != "" {
		// We can't actually delete a default resource.
		// Instead, revert it to default settings and mark it as terminated
		s.D.Set("options", s.D.Get("default_options"))
		e = s.Update()

		s.D.Set("state", baremetal.ResourceTerminated)
		return
	}

	return s.Client.DeleteDHCPOptions(s.D.Id(), nil)
}

func (s *DHCPOptionsResourceCrud) buildEntities() (entities []baremetal.DHCPDNSOption) {
	entities = []baremetal.DHCPDNSOption{}
	for _, val := range s.D.Get("options").([]interface{}) {
		data := val.(map[string]interface{})

		servers := []string{}
		for _, val := range data["custom_dns_servers"].([]interface{}) {
			servers = append(servers, val.(string))
		}
		if len(servers) == 0 {
			servers = nil
		}
		searchDomains := []string{}
		for _, val := range data["search_domain_names"].([]interface{}) {
			searchDomains = append(searchDomains, val.(string))
		}
		if len(searchDomains) == 0 {
			searchDomains = nil
		}
		entity := baremetal.DHCPDNSOption{
			Type:              data["type"].(string),
			CustomDNSServers:  servers,
			ServerType:        data["server_type"].(string),
			SearchDomainNames: searchDomains,
		}
		entities = append(entities, entity)
	}
	return
}
