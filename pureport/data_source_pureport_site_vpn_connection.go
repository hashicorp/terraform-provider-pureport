package pureport

import (
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pureport/terraform-provider-pureport/pureport/connection"
)

func dataSourceSiteVPNConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"connection_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile("conn-.{16}"), "Connection ID must start with 'conn-' with 16 trailing characters."),
		},
		"auth_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"speed": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"primary_customer_router_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"routing_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ike_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ike_config": {
			Type:     schema.TypeList,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"esp": {
						Type:     schema.TypeList,
						MaxItems: 1,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"dh_group": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"encryption": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"integrity": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"ike": {
						Type:     schema.TypeList,
						MaxItems: 1,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"dh_group": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"encryption": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"integrity": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"prf": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
		"enable_bgp_password": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"primary_key": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"secondary_customer_router_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"secondary_key": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"traffic_selectors": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"customer_side": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"pureport_side": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"gateways": {
			Computed: true,
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 2,
			Elem: &schema.Resource{
				Schema: connection.VpnGatewaySchema,
			},
		},
	}

	// Add the base items
	for k, v := range connection.GetBaseDataSourceConnectionSchema() {
		connection_schema[k] = v
	}

	return &schema.Resource{
		Read:   dataSourceSiteVPNConnectionRead,
		Schema: connection_schema,
	}
}

func dataSourceSiteVPNConnectionRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("connection_id").(string))

	return resourceSiteVPNConnectionRead(d, m)
}
