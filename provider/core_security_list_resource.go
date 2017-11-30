// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/oracle/bmcs-go-sdk"

	"github.com/oracle/terraform-provider-oci/crud"
)

var transportSchema = &schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"max": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"min": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	},
}

var icmpSchema = &schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"code": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	},
}

var egressRules = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"destination": {
			Type:     schema.TypeString,
			Required: true,
		},
		"icmp_options": icmpSchema,
		"protocol": {
			Type:     schema.TypeString,
			Required: true,
		},
		"tcp_options": transportSchema,
		"udp_options": transportSchema,
		"stateless": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	},
}

var ingressRules = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"icmp_options": icmpSchema,
		"protocol": {
			Type:     schema.TypeString,
			Required: true,
		},
		"source": {
			Type:     schema.TypeString,
			Required: true,
		},
		"tcp_options": transportSchema,
		"udp_options": transportSchema,
		"stateless": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	},
}

func SecurityListResource() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: crud.DefaultTimeout,
		Create:   createSecurityList,
		Read:     readSecurityList,
		Update:   updateSecurityList,
		Delete:   deleteSecurityList,
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
			"default_egress_security_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     egressRules,
			},
			"egress_security_rules": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     egressRules,
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
			"default_ingress_security_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     ingressRules,
			},
			"ingress_security_rules": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     ingressRules,
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

func createSecurityList(d *schema.ResourceData, m interface{}) (e error) {
	client := m.(*OracleClients)
	crd := &SecurityListResourceCrud{}
	crd.D = d
	crd.Client = client.client
	return crud.CreateResource(d, crd)
}

func readSecurityList(d *schema.ResourceData, m interface{}) (e error) {
	client := m.(*OracleClients)
	crd := &SecurityListResourceCrud{}
	crd.D = d
	crd.Client = client.client
	return crud.ReadResource(crd)
}

func updateSecurityList(d *schema.ResourceData, m interface{}) (e error) {
	client := m.(*OracleClients)
	crd := &SecurityListResourceCrud{}
	crd.D = d
	crd.Client = client.client
	return crud.UpdateResource(d, crd)
}

func deleteSecurityList(d *schema.ResourceData, m interface{}) (e error) {
	client := m.(*OracleClients)
	crd := &SecurityListResourceCrud{}
	crd.D = d
	crd.Client = client.clientWithoutNotFoundRetries
	return crud.DeleteResource(d, crd)
}

type SecurityListResourceCrud struct {
	crud.BaseCrud
	Res *baremetal.SecurityList
}

func (s *SecurityListResourceCrud) ID() string {
	return s.Res.ID
}

func (s *SecurityListResourceCrud) CreatedPending() []string {
	return []string{baremetal.ResourceProvisioning}
}

func (s *SecurityListResourceCrud) CreatedTarget() []string {
	return []string{baremetal.ResourceAvailable}
}

func (s *SecurityListResourceCrud) DeletedPending() []string {
	return []string{baremetal.ResourceTerminating}
}

func (s *SecurityListResourceCrud) DeletedTarget() []string {
	return []string{baremetal.ResourceTerminated}
}

func (s *SecurityListResourceCrud) State() string {
	return s.Res.State
}

func (s *SecurityListResourceCrud) Create() (e error) {
	// If we are creating a default resource, then don't have to
	// actually create it. Just set the ID and update it.
	if defaultId := s.D.Get("default_id").(string); defaultId != "" {
		var res *baremetal.SecurityList
		if res, e = s.Client.GetSecurityList(defaultId); e != nil {
			return
		}
		s.D.Set("default_egress_security_rules", egressRulesToMapArray(res.EgressSecurityRules))
		s.D.Set("default_ingress_security_rules", ingressRulesToMapArray(res.IngressSecurityRules))

		s.D.SetId(defaultId)
		e = s.Update()
		return
	}

	compartmentID := s.D.Get("compartment_id").(string)
	egress := s.buildEgressRules()
	ingress := s.buildIngressRules()
	vcnID := s.D.Get("vcn_id").(string)

	opts := &baremetal.CreateOptions{}
	opts.DisplayName = s.D.Get("display_name").(string)

	s.Res, e = s.Client.CreateSecurityList(compartmentID, vcnID, egress, ingress, opts)
	s.D.Set("default_ingress_security_rules", nil)
	s.D.Set("default_egress_security_rules", nil)

	return
}

func (s *SecurityListResourceCrud) Get() (e error) {
	res, e := s.Client.GetSecurityList(s.D.Id())
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

func (s *SecurityListResourceCrud) Update() (e error) {
	opts := &baremetal.UpdateSecurityListOptions{}

	if displayName, ok := s.D.GetOk("display_name"); ok {
		opts.DisplayName = displayName.(string)
	}

	if egress := s.buildEgressRules(); egress != nil {
		opts.EgressRules = egress
	}
	if ingress := s.buildIngressRules(); ingress != nil {
		opts.IngressRules = ingress
	}

	s.Res, e = s.Client.UpdateSecurityList(s.D.Id(), opts)
	return
}

func egressRulesToMapArray(rules []baremetal.EgressSecurityRule) (res []map[string]interface{}) {
	for _, egressRule := range rules {
		confEgressRule := map[string]interface{}{}
		confEgressRule["destination"] = egressRule.Destination
		confEgressRule = buildConfRule(
			confEgressRule,
			egressRule.Protocol,
			egressRule.ICMPOptions,
			egressRule.TCPOptions,
			egressRule.UDPOptions,
			&egressRule.IsStateless,
		)
		res = append(res, confEgressRule)
	}

	return
}

func ingressRulesToMapArray(rules []baremetal.IngressSecurityRule) (res []map[string]interface{}) {
	for _, ingressRule := range rules {
		confIngressRule := map[string]interface{}{}
		confIngressRule["source"] = ingressRule.Source
		confIngressRule = buildConfRule(
			confIngressRule,
			ingressRule.Protocol,
			ingressRule.ICMPOptions,
			ingressRule.TCPOptions,
			ingressRule.UDPOptions,
			&ingressRule.IsStateless,
		)
		res = append(res, confIngressRule)
	}

	return
}

func (s *SecurityListResourceCrud) SetData() {
	s.D.Set("compartment_id", s.Res.CompartmentID)
	s.D.Set("display_name", s.Res.DisplayName)

	s.D.Set("egress_security_rules", egressRulesToMapArray(s.Res.EgressSecurityRules))
	s.D.Set("ingress_security_rules", ingressRulesToMapArray(s.Res.IngressSecurityRules))

	s.D.Set("state", s.Res.State)
	s.D.Set("time_created", s.Res.TimeCreated.String())
	s.D.Set("vcn_id", s.Res.VcnID)
}

func (s *SecurityListResourceCrud) Delete() (e error) {
	if s.D.Get("default_id") != "" {
		// We can't actually delete a default resource.
		// Instead, revert it to default settings and mark it as terminated
		s.D.Set("egress_security_rules", s.D.Get("default_egress_security_rules"))
		s.D.Set("ingress_security_rules", s.D.Get("default_ingress_security_rules"))
		e = s.Update()

		s.D.Set("state", baremetal.ResourceTerminated)
		return
	}

	return s.Client.DeleteSecurityList(s.D.Id(), nil)
}

func (s *SecurityListResourceCrud) buildEgressRules() (sdkRules []baremetal.EgressSecurityRule) {
	sdkRules = []baremetal.EgressSecurityRule{}
	for _, val := range s.D.Get("egress_security_rules").([]interface{}) {
		confRule := val.(map[string]interface{})

		sdkRule := baremetal.EgressSecurityRule{
			Destination: confRule["destination"].(string),
			ICMPOptions: s.buildICMPOptions(confRule),
			Protocol:    confRule["protocol"].(string),
			TCPOptions:  s.buildTCPOptions(confRule),
			UDPOptions:  s.buildUDPOptions(confRule),
			IsStateless: confRule["stateless"].(bool),
		}

		sdkRules = append(sdkRules, sdkRule)
	}
	return
}

func (s *SecurityListResourceCrud) buildIngressRules() (sdkRules []baremetal.IngressSecurityRule) {
	sdkRules = []baremetal.IngressSecurityRule{}
	for _, val := range s.D.Get("ingress_security_rules").([]interface{}) {
		confRule := val.(map[string]interface{})

		sdkRule := baremetal.IngressSecurityRule{
			ICMPOptions: s.buildICMPOptions(confRule),
			Protocol:    confRule["protocol"].(string),
			Source:      confRule["source"].(string),
			TCPOptions:  s.buildTCPOptions(confRule),
			UDPOptions:  s.buildUDPOptions(confRule),
			IsStateless: confRule["stateless"].(bool),
		}

		sdkRules = append(sdkRules, sdkRule)
	}
	return
}

func (s *SecurityListResourceCrud) buildICMPOptions(conf map[string]interface{}) (opts *baremetal.ICMPOptions) {
	l := conf["icmp_options"].([]interface{})
	if len(l) > 0 {
		confOpts := l[0].(map[string]interface{})
		opts = &baremetal.ICMPOptions{
			Code: uint64(confOpts["code"].(int)),
			Type: uint64(confOpts["type"].(int)),
		}
	}
	return
}

func (s *SecurityListResourceCrud) buildTCPOptions(conf map[string]interface{}) (opts *baremetal.TCPOptions) {
	l := conf["tcp_options"].([]interface{})
	if len(l) > 0 {
		confOpts := l[0].(map[string]interface{})
		opts = &baremetal.TCPOptions{
			baremetal.PortRange{
				Max: uint64(confOpts["max"].(int)),
				Min: uint64(confOpts["min"].(int)),
			},
		}
	}
	return
}

func (s *SecurityListResourceCrud) buildUDPOptions(conf map[string]interface{}) (opts *baremetal.UDPOptions) {
	l := conf["udp_options"].([]interface{})
	if len(l) > 0 {
		confOpts := l[0].(map[string]interface{})
		opts = &baremetal.UDPOptions{
			baremetal.PortRange{
				Max: uint64(confOpts["max"].(int)),
				Min: uint64(confOpts["min"].(int)),
			},
		}
	}
	return
}

func buildConfICMPOptions(opts *baremetal.ICMPOptions) (list []interface{}) {
	confOpts := map[string]interface{}{
		"code": int(opts.Code),
		"type": int(opts.Type),
	}
	return []interface{}{confOpts}
}

func buildConfTransportOptions(portRange baremetal.PortRange) (list []interface{}) {
	confOpts := map[string]interface{}{
		"max": int(portRange.Max),
		"min": int(portRange.Min),
	}
	return []interface{}{confOpts}
}

func buildConfRule(
	confRule map[string]interface{},
	protocol string,
	icmpOpts *baremetal.ICMPOptions,
	tcpOpts *baremetal.TCPOptions,
	udpOpts *baremetal.UDPOptions,
	stateless *bool,
) map[string]interface{} {
	confRule["protocol"] = protocol
	if icmpOpts != nil {
		confRule["icmp_options"] = buildConfICMPOptions(icmpOpts)
	}
	if tcpOpts != nil {
		confRule["tcp_options"] = buildConfTransportOptions(tcpOpts.DestinationPortRange)
	}
	if udpOpts != nil {
		confRule["udp_options"] = buildConfTransportOptions(udpOpts.DestinationPortRange)
	}
	if stateless != nil {
		confRule["stateless"] = *stateless
	}
	return confRule
}
