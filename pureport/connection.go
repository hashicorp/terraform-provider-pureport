package pureport

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pureport/pureport-sdk-go/pureport/session"
	"github.com/pureport/pureport-sdk-go/pureport/swagger"
)

var (
	DeletableState = map[string]bool{
		"FAILED_TO_PROVISION": true,
		"ACTIVE":              true,
		"DOWN":                true,
		"FAILED_TO_UPDATE":    true,
		"FAILED_TO_DELETE":    true,
		"DELETED":             true,
	}
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
		"billing_term": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "HOURLY",
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

func DeleteConnection(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)
	ctx := sess.GetSessionContext()
	connectionId := d.Id()

	// Wait until we are in a state that we can trigger a delete from
	log.Printf("[Info] Waiting to trigger a delete.")
	for i := 0; i < 100; i++ {

		c, resp, err := sess.Client.ConnectionsApi.GetConnection(ctx, connectionId)
		if err != nil {
			return fmt.Errorf("[Error] Error deleting data for AWS Connection: %s", err)
		}

		if resp.StatusCode >= 300 {
			return fmt.Errorf("Error Response while attempting to delete AWS Connection: code=%v", resp.StatusCode)
		}

		conn := reflect.ValueOf(c)
		if DeletableState[conn.FieldByName("State").String()] {
			break
		}

		time.Sleep(time.Second)
	}

	// Delete
	_, resp, err := sess.Client.ConnectionsApi.DeleteConnection(ctx, connectionId)

	if err != nil {
		return fmt.Errorf("[Error] Error deleting data for AWS Connection: %s", err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while deleting AWS Connection: code=%v", resp.StatusCode)
	}

	for i := 0; i < 100; i++ {

		log.Printf("[Info] Waiting for channel to be deleted: attempt %d", i)
		_, resp, _ := sess.Client.ConnectionsApi.GetConnection(ctx, connectionId)

		if resp.StatusCode == 404 {
			d.SetId("")
			break
		}

		time.Sleep(time.Second)
	}

	return nil

}

// AddCustomerNetworks to decode the customer network information
func AddCustomerNetworks(d *schema.ResourceData) []swagger.CustomerNetwork {
	customerNetworks := []swagger.CustomerNetwork{}

	if data, ok := d.GetOk("customer_networks"); ok {
		for _, cn := range data.([]map[string]string) {

			new := swagger.CustomerNetwork{
				Name:    cn["name"],
				Address: cn["Address"],
			}

			customerNetworks = append(customerNetworks, new)
		}
	}

	return customerNetworks
}

func AddNATConfiguration(d *schema.ResourceData) *swagger.NatConfig {

	natConfig := &swagger.NatConfig{
		Enabled: false,
	}

	if data, ok := d.GetOk("nat_config"); ok {

		config := data.(map[string]interface{})
		natConfig.Enabled = config["enabled"].(bool)

		for _, m := range config["mappings"].([]map[string]string) {

			new := swagger.NatMapping{
				NativeCidr: m["native_cidr"],
			}

			natConfig.Mappings = append(natConfig.Mappings, new)
		}
	}

	return natConfig
}

func AddCloudServices(d *schema.ResourceData) []swagger.Link {

	cloudServices := []swagger.Link{}

	if data, ok := d.GetOk("cloud_services"); ok {
		for _, cs := range data.([]map[string]string) {

			new := swagger.Link{
				Id:   cs["id"],
				Href: cs["href"],
			}

			cloudServices = append(cloudServices, new)
		}
	}

	return cloudServices
}

func AddPeeringType(d *schema.ResourceData) *swagger.PeeringConfiguration {

	peeringConfig := &swagger.PeeringConfiguration{}

	if data, ok := d.GetOk("peering"); ok {
		peeringConfig.Type_ = data.(string)
	} else {
		peeringConfig.Type_ = "Private"
	}

	return peeringConfig
}
