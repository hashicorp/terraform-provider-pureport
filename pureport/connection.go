package pureport

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pureport/pureport-sdk-go/pureport/client"
)

var (
	StandardGatewaySchema = map[string]*schema.Schema{
		"availability_domain": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"link_state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"customer_asn": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"customer_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"pureport_asn": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"pureport_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"bgp_password": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"peering_subnet": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"public_nat_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"remote_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vlan": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}

	VpnGatewaySchema = map[string]*schema.Schema{
		"availability_domain": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"link_state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"customer_asn": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"customer_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"pureport_asn": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"pureport_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"bgp_password": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"peering_subnet": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"public_nat_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"customer_gateway_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"customer_vti_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"pureport_gateway_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"pureport_vti_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vpn_auth_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vpn_auth_key": {
			Type:     schema.TypeString,
			Computed: true,
		},
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
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
					"mappings": {
						Type:     schema.TypeSet,
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
		"network_href": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

// FlattenGateway flattens the provide gateway to a map for use with terraform
func FlattenStandardGateway(gateway *client.StandardGateway) map[string]interface{} {
	return map[string]interface{}{
		"availability_domain": gateway.AvailabilityDomain,
		"name":                gateway.Name,
		"description":         gateway.Description,
		"link_state":          gateway.LinkState,
		"remote_id":           gateway.RemoteId,
		"vlan":                gateway.Vlan,
		"customer_asn":        gateway.BgpConfig.CustomerASN,
		"customer_ip":         gateway.BgpConfig.CustomerIP,
		"pureport_asn":        gateway.BgpConfig.PureportASN,
		"pureport_ip":         gateway.BgpConfig.PureportIP,
		"bgp_password":        gateway.BgpConfig.Password,
		"peering_subnet":      gateway.BgpConfig.PeeringSubnet,
		"public_nat_ip":       gateway.BgpConfig.PublicNatIp,
	}
}

// FlattenGateway flattens the provide gateway to a map for use with terraform
func FlattenVpnGateway(gateway *client.VpnGateway) map[string]interface{} {
	return map[string]interface{}{
		"availability_domain": gateway.AvailabilityDomain,
		"name":                gateway.Name,
		"description":         gateway.Description,
		"link_state":          gateway.LinkState,
		"customer_gateway_ip": gateway.CustomerGatewayIP,
		"customer_vti_ip":     gateway.CustomerVtiIP,
		"pureport_gateway_ip": gateway.PureportGatewayIP,
		"pureport_vti_ip":     gateway.PureportVtiIP,
		"vpn_auth_type":       gateway.Auth.Type_,
		"vpn_auth_key":        gateway.Auth.Key,
		"customer_asn":        gateway.BgpConfig.CustomerASN,
		"customer_ip":         gateway.BgpConfig.CustomerIP,
		"pureport_asn":        gateway.BgpConfig.PureportASN,
		"pureport_ip":         gateway.BgpConfig.PureportIP,
		"bgp_password":        gateway.BgpConfig.Password,
		"peering_subnet":      gateway.BgpConfig.PeeringSubnet,
		"public_nat_ip":       gateway.BgpConfig.PublicNatIp,
	}
}

func flattenCustomerNetworks(customerNetworks []client.CustomerNetwork) []map[string]string {

	var out []map[string]string
	for _, cn := range customerNetworks {
		out = append(out, map[string]string{
			"name":    cn.Name,
			"address": cn.Address,
		})
	}

	return out
}

func FlattenNatConfig(config *client.NatConfig) []map[string]interface{} {

	output := make([]map[string]interface{}, 1)
	output[0] = map[string]interface{}{
		"blocks":    config.Blocks,
		"enabled":   config.Enabled,
		"pnat_cidr": config.PnatCidr,
		"mappings":  flattenMappings(config.Mappings),
	}

	return output
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

func WaitForConnection(name string, d *schema.ResourceData, m interface{}) error {

	config := m.(*Config)
	ctx := config.Session.GetSessionContext()
	connectionId := d.Id()

	log.Printf("[Info] Waiting for connection to come up.")

	createStateConf := &resource.StateChangeConf{
		Pending: []string{
			"INITIALIZING",
			"PROVISIONING",
			"UPDATING",
		},
		Target: []string{
			"ACTIVE",
		},
		Refresh: func() (interface{}, string, error) {

			c, resp, err := config.Session.Client.ConnectionsApi.GetConnection(ctx, connectionId)
			if err != nil {
				return 0, "", fmt.Errorf("Error reading data for %s: %s", name, err)
			}

			if resp.StatusCode >= 300 {
				return 0, "", fmt.Errorf("Error received while waiting for creation of %s: code=%v", name, resp.StatusCode)
			}

			conn := reflect.ValueOf(c)
			state := conn.FieldByName("State").String()

			return c, state, nil

		},
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     5 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	_, err := createStateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for connection (%s) to be created: %s", connectionId, err)
	}

	return nil
}

func DeleteConnection(name string, d *schema.ResourceData, m interface{}) error {

	config := m.(*Config)
	ctx := config.Session.GetSessionContext()
	connectionId := d.Id()

	// Wait until we are in a state that we can trigger a delete from
	log.Printf("[Info] Waiting to trigger a delete.")

	waitingStateConf := &resource.StateChangeConf{
		Pending: []string{
			"INITIALIZING",
			"PROVISIONING",
			"UPDATING",
			"DELETING",
		},
		Target: []string{
			"FAILED_TO_PROVISION",
			"ACTIVE",
			"DOWN",
			"FAILED_TO_UPDATE",
			"FAILED_TO_DELETE",
			"DELETED",
		},
		Refresh: func() (interface{}, string, error) {

			c, resp, err := config.Session.Client.ConnectionsApi.GetConnection(ctx, connectionId)
			if err != nil {
				return 0, "", fmt.Errorf("Error deleting data for %s: %s", name, err)
			}

			if resp.StatusCode >= 300 {
				return 0, "", fmt.Errorf("Error Response while attempting to delete %s: code=%v", name, resp.StatusCode)
			}

			conn := reflect.ValueOf(c)
			state := conn.FieldByName("State").String()

			return c, state, nil

		},
		Timeout:                   d.Timeout(schema.TimeoutDelete),
		Delay:                     5 * time.Second,
		MinTimeout:                1 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	_, err := waitingStateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for connection (%s) to be deletable: %s", connectionId, err)
	}

	// Delete
	_, resp, err := config.Session.Client.ConnectionsApi.DeleteConnection(ctx, connectionId)
	if err != nil {
		return fmt.Errorf("Error deleting data for %s: %s", name, err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while deleting %s: code=%v", name, resp.StatusCode)
	}

	log.Printf("[Info] Waiting for connection to be deleted")

	deleteStateConf := &resource.StateChangeConf{
		Pending: []string{
			"INITIALIZING",
			"PROVISIONING",
			"UPDATING",
			"DELETING",
		},
		Target: []string{
			"DELETED",
		},
		Refresh: func() (interface{}, string, error) {

			c, resp, err := config.Session.Client.ConnectionsApi.GetConnection(ctx, connectionId)

			if resp.StatusCode == 404 {
				return 0, "DELETED", nil
			}

			if err != nil {
				return 0, "", fmt.Errorf("Error Response while deleting %s: error=%s", name, err)
			}

			conn := reflect.ValueOf(c)
			state := conn.FieldByName("State").String()

			return c, state, nil

		},
		Timeout:                   d.Timeout(schema.TimeoutDelete),
		Delay:                     5 * time.Second,
		MinTimeout:                1 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	_, err = deleteStateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for connection (%s) to be created: %s", connectionId, err)
	}

	d.SetId("")

	return nil
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

		tmp_config := data.([]interface{})
		config := tmp_config[0].(map[string]interface{})
		natConfig.Enabled = config["enabled"].(bool)

		for _, m := range config["mappings"].(*schema.Set).List() {

			mapping := m.(map[string]interface{})

			new := client.NatMapping{
				NativeCidr: mapping["native_cidr"].(string),
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
