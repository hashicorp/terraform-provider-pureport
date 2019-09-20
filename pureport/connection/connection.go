package connection

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
	"github.com/terraform-providers/terraform-provider-pureport/pureport/configuration"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/filter"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/tags"
)

const (
	AwsConnectionName     = "AWS Cloud Connection"
	AzureConnectionName   = "Azure Cloud Connection"
	GoogleConnectionName  = "Google Cloud Connection"
	SiteVPNConnectionName = "SiteVPN Connection"
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

func GetBaseResourceConnectionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"href": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
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
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"customer_networks": {
			Type:     schema.TypeSet,
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
			Type:     schema.TypeInt,
			Optional: true,
			// This should be 4,294,967,295(64bit) but for 32bit Arch `int` is only 32bits so
			// 2,147,483,647.
			ValidateFunc: validation.IntBetween(0, 2147483647),
		},
		"high_availability": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"tags": tags.TagsSchema(),
	}
}

func GetBaseDataSourceConnectionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": filter.DataSourceFiltersSchema(),
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"href": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"location_href": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"network_href": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"customer_networks": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"address": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"nat_config": {
			Type:     schema.TypeList,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"mappings": {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"native_cidr": {
									Type:     schema.TypeString,
									Computed: true,
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
			Computed: true,
		},
		"customer_asn": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"high_availability": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"tags": tags.TagsSchema(),
	}
}

// FlattenGateway flattens the provide gateway to a map for use with terraform
func FlattenStandardGateway(gateway *client.StandardGateway) (out map[string]interface{}) {

	out = map[string]interface{}{
		"availability_domain": gateway.AvailabilityDomain,
		"name":                gateway.Name,
		"description":         gateway.Description,
		"remote_id":           gateway.RemoteId,
		"vlan":                gateway.Vlan,
		"customer_asn":        0,
		"customer_ip":         "",
		"pureport_asn":        0,
		"pureport_ip":         "",
		"bgp_password":        "",
		"peering_subnet":      "",
		"public_nat_ip":       "",
	}

	// If we are using BGP, include the confiuration
	if gateway.BgpConfig != nil {
		out["customer_asn"] = gateway.BgpConfig.CustomerASN
		out["customer_ip"] = gateway.BgpConfig.CustomerIP
		out["pureport_asn"] = gateway.BgpConfig.PureportASN
		out["pureport_ip"] = gateway.BgpConfig.PureportIP
		out["bgp_password"] = gateway.BgpConfig.Password
		out["peering_subnet"] = gateway.BgpConfig.PeeringSubnet
		out["public_nat_ip"] = gateway.BgpConfig.PublicNatIp
	}

	return
}

// FlattenGateway flattens the provide gateway to a map for use with terraform
func FlattenVpnGateway(gateway *client.VpnGateway) (out map[string]interface{}) {

	out = map[string]interface{}{
		"availability_domain": gateway.AvailabilityDomain,
		"name":                gateway.Name,
		"description":         gateway.Description,
		"customer_gateway_ip": gateway.CustomerGatewayIP,
		"customer_vti_ip":     gateway.CustomerVtiIP,
		"pureport_gateway_ip": gateway.PureportGatewayIP,
		"pureport_vti_ip":     gateway.PureportVtiIP,
		"vpn_auth_type":       gateway.Auth.Type_,
		"vpn_auth_key":        gateway.Auth.Key,
		"customer_asn":        0,
		"customer_ip":         "",
		"pureport_asn":        0,
		"pureport_ip":         "",
		"bgp_password":        "",
		"peering_subnet":      "",
		"public_nat_ip":       "",
	}

	// If we are using BGP, include the confiuration
	if gateway.BgpConfig != nil {
		out["customer_asn"] = gateway.BgpConfig.CustomerASN
		out["customer_ip"] = gateway.BgpConfig.CustomerIP
		out["pureport_asn"] = gateway.BgpConfig.PureportASN
		out["pureport_ip"] = gateway.BgpConfig.PureportIP
		out["bgp_password"] = gateway.BgpConfig.Password
		out["peering_subnet"] = gateway.BgpConfig.PeeringSubnet
		out["public_nat_ip"] = gateway.BgpConfig.PublicNatIp
	}

	return
}

func FlattenCustomerNetworks(customerNetworks []client.CustomerNetwork) (out []map[string]string) {

	for _, cn := range customerNetworks {
		out = append(out, map[string]string{
			"name":    cn.Name,
			"address": cn.Address,
		})
	}

	return
}

func FlattenNatConfig(config *client.NatConfig) (out []map[string]interface{}) {

	return append(out, map[string]interface{}{
		"blocks":    config.Blocks,
		"enabled":   config.Enabled,
		"pnat_cidr": config.PnatCidr,
		"mappings":  flattenMappings(config.Mappings),
	})
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

	config := m.(*configuration.Config)
	ctx := config.Session.GetSessionContext()
	connectionId := d.Id()

	log.Printf("[Info] Waiting for connection to come up.")

	createStateConf := &resource.StateChangeConf{
		Pending: []string{
			"INITIALIZING",
			"PROVISIONING",
			"UPDATING",
			"WAITING_TO_PROVISION",
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

	config := m.(*configuration.Config)
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
			"WAITING_TO_PROVISION",
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
			"WAITING_TO_PROVISION",
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

		for _, cn := range data.(*schema.Set).List() {

			network := cn.(map[string]interface{})

			new := client.CustomerNetwork{
				Name:    network["name"].(string),
				Address: network["address"].(string),
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
