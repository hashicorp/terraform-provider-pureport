package pureport

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/pureport/pureport-sdk-go/pureport/session"
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
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"speed": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntInSlice([]int{50, 100, 200, 300, 400, 500, 1000, 10000}),
		},
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
		"nat_config": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
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
		"billing_term": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "HOURLY",
		},
		"customer_asn": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(0, 4294967295),
		},
		"high_availability": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"location_href": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"network": {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: true,
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
	}
}

func flattenConnection(connection client.Connection) map[string]interface{} {
	return map[string]interface{}{
		"customer_asn":      connection.CustomerASN,
		"customer_networks": flattenCustomerNetworks(connection.CustomerNetworks),
		"high_availability": connection.HighAvailability,
		"location":          flattenLink(connection.Location),
		"network":           flattenLink(connection.Network),
		"name":              connection.Name,
		"description":       connection.Description,
		"speed":             connection.Speed,
	}
}

func flattenLink(link *client.Link) map[string]interface{} {
	return map[string]interface{}{
		"id":   link.Id,
		"href": link.Href,
	}
}

func flattenCustomerNetworks(networks []client.CustomerNetwork) (out []map[string]interface{}) {

	for _, network := range networks {

		n := map[string]interface{}{
			"name":    network.Name,
			"address": network.Address,
		}

		out = append(out, n)
	}

	return
}

func flattenNatConfig(config client.NatConfig) map[string]interface{} {
	return map[string]interface{}{
		"blocks":    config.Blocks,
		"enabled":   config.Enabled,
		"pnat_cidr": config.PnatCidr,
		"mappings":  flattenMappings(config.Mappings),
	}
}

func flattenMappings(mappings []client.NatMapping) (out []map[string]interface{}) {

	for _, mapping := range mappings {

		m := map[string]interface{}{
			"nat_cidr":    mapping.NatCidr,
			"native_cidr": mapping.NativeCidr,
		}

		out = append(out, m)
	}

	return
}

func DeleteConnection(name string, d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)
	ctx := sess.GetSessionContext()
	connectionId := d.Id()

	// Wait until we are in a state that we can trigger a delete from
	log.Printf("[Info] Waiting to trigger a delete.")

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 3 * time.Minute

	wait_to_delete := func() error {

		c, resp, err := sess.Client.ConnectionsApi.GetConnection(ctx, connectionId)
		if err != nil {
			return backoff.Permanent(
				fmt.Errorf("Error deleting data for %s: %s", name, err),
			)
		}

		if resp.StatusCode >= 300 {
			return backoff.Permanent(
				fmt.Errorf("Error Response while attempting to delete %s: code=%v", name, resp.StatusCode),
			)
		}

		conn := reflect.ValueOf(c)
		if DeletableState[conn.FieldByName("State").String()] {
			return nil
		} else {
			return fmt.Errorf("Not Completed ...")
		}
	}

	if err := backoff.Retry(wait_to_delete, b); err != nil {
		return err
	}

	// Delete
	_, resp, err := sess.Client.ConnectionsApi.DeleteConnection(ctx, connectionId)
	if err != nil {
		return fmt.Errorf("Error deleting data for %s: %s", name, err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while deleting %s: code=%v", name, resp.StatusCode)
	}

	log.Printf("[Info] Waiting for connection to be deleted")
	wait_for_delete := func() error {

		log.Printf("Retrying ...%+v", b.GetElapsedTime())
		_, resp, _ := sess.Client.ConnectionsApi.GetConnection(ctx, connectionId)

		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		} else {
			return fmt.Errorf("Not Completed ...")
		}
	}

	return backoff.Retry(wait_for_delete, b)
}

// ExpandCustomerNetworks to decode the customer network information
func ExpandCustomerNetworks(d *schema.ResourceData) []client.CustomerNetwork {

	if data, ok := d.GetOk("customer_networks"); ok {

		customerNetworks := []client.CustomerNetwork{}

		for _, cn := range data.([]map[string]string) {

			new := client.CustomerNetwork{
				Name:    cn["name"],
				Address: cn["address"],
			}

			customerNetworks = append(customerNetworks, new)
		}

		return customerNetworks
	}

	return nil
}

func ExpandNATConfiguration(d *schema.ResourceData) *client.NatConfig {

	if data, ok := d.GetOk("nat_config"); ok {

		natConfig := &client.NatConfig{}

		config := data.(map[string]interface{})
		natConfig.Enabled = config["enabled"].(bool)

		for _, m := range config["mappings"].([]map[string]string) {

			new := client.NatMapping{
				NativeCidr: m["native_cidr"],
			}

			natConfig.Mappings = append(natConfig.Mappings, new)
		}
		return natConfig
	}

	return nil
}

func ExpandCloudServices(d *schema.ResourceData) []client.Link {

	if data, ok := d.GetOk("cloud_service_hrefs"); ok {

		cloudServices := []client.Link{}
		for _, cs := range data.([]interface{}) {
			cloudServices = append(cloudServices, client.Link{Href: cs.(string)})
		}

		sort.Slice(cloudServices, func(i int, j int) bool {
			return cloudServices[i].Href < cloudServices[j].Href
		})

		return cloudServices
	}

	return nil
}

func ExpandPeeringType(d *schema.ResourceData) *client.PeeringConfiguration {

	peeringConfig := &client.PeeringConfiguration{}

	if data, ok := d.GetOk("peering_type"); ok {
		peeringConfig.Type_ = data.(string)
	} else {
		peeringConfig.Type_ = "Private"
	}

	return peeringConfig
}
