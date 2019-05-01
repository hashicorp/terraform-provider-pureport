package pureport

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/pureport/pureport-sdk-go/pureport/session"
)

func resourceSiteVPNConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"auth_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "PSK",
			ValidateFunc: validation.StringInSlice([]string{"psk"}, true),
		},
		"enable_bgp_password": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"ike_version": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"V1", "V2"}, true),
		},
		"ikev1_config": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"esp": {
						Type:     schema.TypeList,
						MaxItems: 1,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"dh_group": {
									Type:     schema.TypeString,
									Required: true,
								},
								"encryption": {
									Type:     schema.TypeString,
									Required: true,
								},
								"integrity": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
					"ike": {
						Type:     schema.TypeList,
						MaxItems: 1,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"dh_group": {
									Type:     schema.TypeString,
									Required: true,
								},
								"encryption": {
									Type:     schema.TypeString,
									Required: true,
								},
								"integrity": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
		"ikev2_config": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"esp": {
						Type:     schema.TypeList,
						MaxItems: 1,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"dh_group": {
									Type:     schema.TypeString,
									Required: true,
								},
								"encryption": {
									Type:     schema.TypeString,
									Required: true,
								},
								"integrity": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
					"ike": {
						Type:     schema.TypeList,
						MaxItems: 1,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"dh_group": {
									Type:     schema.TypeString,
									Required: true,
								},
								"encryption": {
									Type:     schema.TypeString,
									Required: true,
								},
								"integrity": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"prf": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
		"primary_customer_router_ip": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"primary_key": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"routing_type": {
			Type:     schema.TypeString,
			Required: true,
		},
		"secondary_customer_router_ip": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"secondary_key": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"traffic_selectors": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"customer_side": {
						Type:     schema.TypeString,
						Required: true,
					},
					"pureport_side": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
	}

	// Add the base items
	for k, v := range getBaseConnectionSchema() {
		connection_schema[k] = v
	}

	return &schema.Resource{
		Create: resourceSiteVPNConnectionCreate,
		Read:   resourceSiteVPNConnectionRead,
		Update: resourceSiteVPNConnectionUpdate,
		Delete: resourceSiteVPNConnectionDelete,

		Schema: connection_schema,
	}
}

func addTrafficSelectorMappings(d *schema.ResourceData) []client.TrafficSelectorMapping {

	if data, ok := d.GetOk("traffic_selectors"); ok {

		mappings := []client.TrafficSelectorMapping{}

		for _, m := range data.([]map[string]string) {

			new := client.TrafficSelectorMapping{
				CustomerSide: m["customer_side"],
				PureportSide: m["pureport_side"],
			}

			mappings = append(mappings, new)
		}

		return mappings
	}

	return nil
}

func addIkeVersion1(d *schema.ResourceData) *client.Ikev1Config {

	config := &client.Ikev1Config{}

	if data, ok := d.GetOk("ikev1_config"); ok {

		raw_config := data.(map[string]interface{})
		esp := raw_config["esp"].(map[string]interface{})
		ike := raw_config["ike"].(map[string]interface{})

		config.Esp = &client.Ikev1EspConfig{
			DhGroup:    esp["dh_group"].(string),
			Encryption: esp["encryption"].(string),
			Integrity:  esp["integrity"].(string),
		}

		config.Ike = &client.Ikev1IkeConfig{
			DhGroup:    ike["dh_group"].(string),
			Encryption: ike["encryption"].(string),
			Integrity:  ike["integrity"].(string),
		}

	} else {

		config.Esp = &client.Ikev1EspConfig{
			DhGroup:    "MODP_2048",
			Encryption: "AES_128",
			Integrity:  "SHA256_HMAC",
		}

		config.Ike = &client.Ikev1IkeConfig{
			DhGroup:    "MODP_2048",
			Encryption: "AES_128",
			Integrity:  "SHA256_HMAC",
		}
	}

	return config

}

func addIkeVersion2(d *schema.ResourceData) *client.Ikev2Config {

	config := &client.Ikev2Config{}

	if data, ok := d.GetOk("ikev2_config"); ok {

		raw_config := data.(map[string]interface{})
		esp := raw_config["esp"].(map[string]interface{})
		ike := raw_config["ike"].(map[string]interface{})

		config.Esp = &client.Ikev2EspConfig{
			DhGroup:    esp["dh_group"].(string),
			Encryption: esp["encryption"].(string),
			Integrity:  esp["integrity"].(string),
		}

		config.Ike = &client.Ikev2IkeConfig{
			DhGroup:    ike["dh_group"].(string),
			Encryption: ike["encryption"].(string),
			Integrity:  ike["integrity"].(string),
			Prf:        ike["prf"].(string),
		}

	} else {

		config.Esp = &client.Ikev2EspConfig{
			DhGroup:    "MODP_2048",
			Encryption: "AES_128",
			Integrity:  "SHA256_HMAC",
		}

		config.Ike = &client.Ikev2IkeConfig{
			DhGroup:    "MODP_2048",
			Encryption: "AES_128",
			Integrity:  "SHA256_HMAC",
		}
	}

	return config
}

func resourceSiteVPNConnectionCreate(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)

	// Generic Connection values
	network := d.Get("network").([]interface{})
	speed := d.Get("speed").(int)
	name := d.Get("name").(string)
	location := d.Get("location").([]interface{})
	billingTerm := d.Get("billing_term").(string)

	// Create the body of the request
	connection := client.SiteIpSecVpnConnection{
		Type_:       "SITE_IPSEC_VPN",
		Name:        name,
		Speed:       int32(speed),
		AuthType:    d.Get("auth_type").(string),
		IkeVersion:  d.Get("ike_version").(string),
		RoutingType: d.Get("routing_type").(string),
		PrimaryKey:  d.Get("primary_key").(string),

		Location: &client.Link{
			Id:   location[0].(map[string]interface{})["id"].(string),
			Href: location[0].(map[string]interface{})["href"].(string),
		},
		Network: &client.Link{
			Id:   network[0].(map[string]interface{})["id"].(string),
			Href: network[0].(map[string]interface{})["href"].(string),
		},
		BillingTerm: billingTerm,
	}

	// Generic Optionals
	connection.CustomerNetworks = AddCustomerNetworks(d)
	connection.Nat = AddNATConfiguration(d)

	if description, ok := d.GetOk("description"); ok {
		connection.Description = description.(string)
	}

	if highAvailability, ok := d.GetOk("high_availability"); ok {
		connection.HighAvailability = highAvailability.(bool)
	}

	if customerASN, ok := d.GetOk("customer_asn"); ok {
		connection.CustomerASN = int64(customerASN.(int))
	}

	// SiteVPN Optionals
	connection.TrafficSelectors = addTrafficSelectorMappings(d)

	if connection.IkeVersion == "1" {
		connection.IkeV1 = addIkeVersion1(d)
	} else {
		connection.IkeV2 = addIkeVersion2(d)
	}

	if authType, ok := d.GetOk("auth_type"); ok {
		connection.AuthType = authType.(string)
	}

	if enableBGPPassword, ok := d.GetOk("enable_bgp_password"); ok {
		connection.EnableBGPPassword = enableBGPPassword.(bool)
	}

	if primaryCustomerRouterIP, ok := d.GetOk("primary_customer_router_ip"); ok {
		connection.PrimaryCustomerRouterIP = primaryCustomerRouterIP.(string)
	}

	if secondaryCustomerRouterIP, ok := d.GetOk("secondary_customer_router_ip"); ok {
		connection.SecondaryCustomerRouterIP = secondaryCustomerRouterIP.(string)
	}

	if secondaryKey, ok := d.GetOk("secondary_key"); ok {
		connection.SecondaryKey = secondaryKey.(string)
	}

	ctx := sess.GetSessionContext()

	opts := client.AddConnectionOpts{
		Body: optional.NewInterface(connection),
	}

	resp, err := sess.Client.ConnectionsApi.AddConnection(
		ctx,
		network[0].(map[string]interface{})["id"].(string),
		&opts,
	)

	if err != nil {
		log.Printf("Error Creating new SiteVPN Connection: %v", err)
		d.SetId("")
		return nil
	}

	if resp.StatusCode >= 300 {
		log.Printf("Error Response while creating new SiteVPN Connection: code=%v", resp.StatusCode)
		d.SetId("")
		return nil
	}

	loc := resp.Header.Get("location")
	u, err := url.Parse(loc)
	if err != nil {
		log.Printf("Error when decoding Connection ID")
		return nil
	}

	id := filepath.Base(u.Path)
	d.SetId(id)

	if id == "" {
		log.Printf("Error when decoding location header")
		return nil
	}

	return resourceSiteVPNConnectionRead(d, m)
}

func resourceSiteVPNConnectionRead(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)
	connectionId := d.Id()
	ctx := sess.GetSessionContext()

	c, resp, err := sess.Client.ConnectionsApi.GetConnection(ctx, connectionId)
	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf("Error Response while reading SiteVPN Connection: code=%v", resp.StatusCode)
			d.SetId("")
		}
		return fmt.Errorf("Error reading data for SiteVPN Connection: %s", err)
	}

	if resp.StatusCode >= 300 {
		fmt.Errorf("Error Response while reading SiteVPN Connection: code=%v", resp.StatusCode)
	}

	conn := c.(client.SiteIpSecVpnConnection)
	d.Set("speed", conn.Speed)
	d.Set("description", conn.Description)
	d.Set("high_availability", conn.HighAvailability)

	var customerNetworks []map[string]string
	for _, cn := range conn.CustomerNetworks {
		customerNetworks = append(customerNetworks, map[string]string{
			"name":    cn.Name,
			"address": cn.Address,
		})
	}
	if err := d.Set("customer_networks", customerNetworks); err != nil {
		return fmt.Errorf("Error setting customer networks for VPN Connection %s: %s", d.Id(), err)
	}

	if err := d.Set("location", []map[string]string{
		{
			"id":   conn.Location.Id,
			"href": conn.Location.Href,
		},
	}); err != nil {
		return fmt.Errorf("Error setting location for VPN Connection %s: %s", d.Id(), err)
	}

	if err := d.Set("network", []map[string]string{
		{
			"id":   conn.Network.Id,
			"href": conn.Network.Href,
		},
	}); err != nil {
		return fmt.Errorf("Error setting network for VPN Connection %s: %s", d.Id(), err)
	}

	d.Set("auth_type", conn.AuthType)
	d.Set("enable_bgp_password", conn.EnableBGPPassword)
	d.Set("ike_version", conn.IkeVersion)

	if conn.IkeVersion == "V1" {
		if err := d.Set("ikev1_config", []map[string]interface{}{
			{
				"esp": []map[string]string{
					{
						"dh_group":   conn.IkeV1.Esp.DhGroup,
						"encryption": conn.IkeV1.Esp.Encryption,
						"integrity":  conn.IkeV1.Esp.Integrity,
					},
				},
				"ike": []map[string]string{
					{
						"dh_group":   conn.IkeV1.Ike.DhGroup,
						"encryption": conn.IkeV1.Ike.Encryption,
						"integrity":  conn.IkeV1.Ike.Integrity,
					},
				},
			},
		}); err != nil {
			return fmt.Errorf("Error setting IKE V1 Configuration for VPN Connection %s: %s", d.Id(), err)
		}
	}

	if conn.IkeVersion == "V2" {
		if err := d.Set("ikev2_config", []map[string]interface{}{
			{
				"esp": []map[string]string{
					{
						"dh_group":   conn.IkeV2.Esp.DhGroup,
						"encryption": conn.IkeV2.Esp.Encryption,
						"integrity":  conn.IkeV2.Esp.Integrity,
					},
				},
				"ike": []map[string]string{
					{
						"dh_group":   conn.IkeV2.Ike.DhGroup,
						"encryption": conn.IkeV2.Ike.Encryption,
						"integrity":  conn.IkeV2.Ike.Integrity,
						"prf":        conn.IkeV2.Ike.Prf,
					},
				},
			},
		}); err != nil {
			return fmt.Errorf("Error setting IKE V2 Configuration for VPN Connection %s: %s", d.Id(), err)
		}
	}

	d.Set("routing_type", conn.RoutingType)
	d.Set("primary_customer_router_ip", conn.PrimaryCustomerRouterIP)
	d.Set("primary_key", conn.PrimaryKey)
	d.Set("secondary_customer_router_ip", conn.SecondaryCustomerRouterIP)
	d.Set("secondary_key", conn.SecondaryKey)

	trafficSelectors := []map[string]string{}
	for _, v := range conn.TrafficSelectors {
		trafficSelectors = append(trafficSelectors, map[string]string{
			"customer_side": v.CustomerSide,
			"pureport_side": v.PureportSide,
		})
	}

	if err := d.Set("traffic_selectors", trafficSelectors); err != nil {
		return fmt.Errorf("Error setting traffics selectors for VPN Connection %s: %s", d.Id(), err)
	}

	return nil
}

func resourceSiteVPNConnectionUpdate(d *schema.ResourceData, m interface{}) error {

	return resourceSiteVPNConnectionRead(d, m)
}

func resourceSiteVPNConnectionDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteConnection(d, m)
}
