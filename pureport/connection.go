package pureport

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pureport/pureport-sdk-go/pureport/swagger"
)

func getBaseConnectionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"customer_networks": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"address": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.CIDRNetwork(16, 32),
					},
				},
			},
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"high_availability": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"location": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeString,
						Required: true,
					},
					"href": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"nat_config": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     schema.TypeBool,
						Required: true,
					},
					"mappings": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"native_cidr": {
									Type:     schema.TypeString,
									Required: true,
								},
								"nat_cidr": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"blocks": {
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"pnat_cidr": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"network": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeString,
						Required: true,
					},
					"href": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"speed": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntInSlice([]int{50, 100, 200, 300, 400, 500, 1000, 10000}),
		},
	}
}

func flattenConnection(connection swagger.Connection) map[string]interface{} {
	return map[string]interface{}{
		"customer_asn":      connection.CustomerASN,
		"customer_networks": flattenCustomerNetworks(connection.CustomerNetworks),
		"description":       connection.Description,
		"high_availability": connection.HighAvailability,
		"location":          flattenLink(connection.Location),
		"name":              connection.Name,
		"network":           flattenLink(connection.Network),
		"speed":             connection.Speed,
	}
}

func flattenLink(link *swagger.Link) map[string]interface{} {
	return map[string]interface{}{
		"id":   link.Id,
		"href": link.Href,
	}
}

func flattenCustomerNetworks(networks []swagger.CustomerNetwork) (out []map[string]interface{}) {

	for _, network := range networks {

		n := map[string]interface{}{
			"name":    network.Name,
			"address": network.Address,
		}

		out = append(out, n)
	}

	return
}

func flattenNatConfig(config swagger.NatConfig) map[string]interface{} {
	return map[string]interface{}{
		"blocks":    config.Blocks,
		"enabled":   config.Enabled,
		"pnat_cidr": config.PnatCidr,
		"mappings":  flattenMappings(config.Mappings),
	}
}

func flattenMappings(mappings []swagger.NatMapping) (out []map[string]interface{}) {

	for _, mapping := range mappings {

		m := map[string]interface{}{
			"nat_cidr":    mapping.NatCidr,
			"native_cidr": mapping.NativeCidr,
		}

		out = append(out, m)
	}

	return
}
