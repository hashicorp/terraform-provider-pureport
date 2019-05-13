package pureport

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/pureport/pureport-sdk-go/pureport/session"
)

const (
	sitevpnConnectionName = "SiteVPN Connection"
)

func resourceSiteVPNConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"auth_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "PSK",
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"psk"}, true),
		},
		"enable_bgp_password": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},
		"ike_version": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"V1", "V2"}, true),
		},
		"ike_config": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			ForceNew: true,
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
			Required: true,
			ForceNew: true,
		},
		"primary_key": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"routing_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"secondary_customer_router_ip": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"secondary_key": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"traffic_selectors": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			ForceNew: true,
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

func expandTrafficSelectorMappings(d *schema.ResourceData) []client.TrafficSelectorMapping {

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

func expandIkeVersion1(d *schema.ResourceData) *client.Ikev1Config {

	config := &client.Ikev1Config{}

	if data, ok := d.GetOk("ike_config"); ok {

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

func expandIkeVersion2(d *schema.ResourceData) *client.Ikev2Config {

	config := &client.Ikev2Config{}

	if data, ok := d.GetOk("ike_config"); ok {

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

func expandSiteVPNConnection(d *schema.ResourceData) client.SiteIpSecVpnConnection {

	// Generic Connection values
	speed := d.Get("speed").(int)

	// Create the body of the request
	c := client.SiteIpSecVpnConnection{
		Type_:       "SITE_IPSEC_VPN",
		Name:        d.Get("name").(string),
		Speed:       int32(speed),
		AuthType:    d.Get("auth_type").(string),
		IkeVersion:  d.Get("ike_version").(string),
		RoutingType: d.Get("routing_type").(string),
		PrimaryKey:  d.Get("primary_key").(string),

		Location: &client.Link{
			Href: d.Get("location_href").(string),
		},
		Network: &client.Link{
			Href: d.Get("network_href").(string),
		},
		BillingTerm: d.Get("billing_term").(string),
	}

	// Generic Optionals
	c.CustomerNetworks = ExpandCustomerNetworks(d)
	c.Nat = ExpandNATConfiguration(d)

	if description, ok := d.GetOk("description"); ok {
		c.Description = description.(string)
	}

	if highAvailability, ok := d.GetOk("high_availability"); ok {
		c.HighAvailability = highAvailability.(bool)
	}

	if customerASN, ok := d.GetOk("customer_asn"); ok {
		c.CustomerASN = int64(customerASN.(int))
	}

	// SiteVPN Optionals
	c.TrafficSelectors = expandTrafficSelectorMappings(d)

	if c.IkeVersion == "V1" {
		c.IkeV1 = expandIkeVersion1(d)
	} else {
		c.IkeV2 = expandIkeVersion2(d)
	}

	if authType, ok := d.GetOk("auth_type"); ok {
		c.AuthType = authType.(string)
	}

	if enableBGPPassword, ok := d.GetOk("enable_bgp_password"); ok {
		c.EnableBGPPassword = enableBGPPassword.(bool)
	}

	if primaryCustomerRouterIP, ok := d.GetOk("primary_customer_router_ip"); ok {
		c.PrimaryCustomerRouterIP = primaryCustomerRouterIP.(string)
	}

	if secondaryCustomerRouterIP, ok := d.GetOk("secondary_customer_router_ip"); ok {
		c.SecondaryCustomerRouterIP = secondaryCustomerRouterIP.(string)
	}

	if secondaryKey, ok := d.GetOk("secondary_key"); ok {
		c.SecondaryKey = secondaryKey.(string)
	}

	return c
}

func resourceSiteVPNConnectionCreate(d *schema.ResourceData, m interface{}) error {

	connection := expandSiteVPNConnection(d)

	sess := m.(*session.Session)

	ctx := sess.GetSessionContext()

	opts := client.AddConnectionOpts{
		Body: optional.NewInterface(connection),
	}

	resp, err := sess.Client.ConnectionsApi.AddConnection(
		ctx,
		filepath.Base(connection.Network.Href),
		&opts,
	)

	if err != nil {

		http_err := err
		json_response := string(err.(client.GenericSwaggerError).Body()[:])
		response, err := structure.ExpandJsonFromString(json_response)
		if err != nil {
			log.Printf("Error Creating new %s: %v", sitevpnConnectionName, err)
		} else {
			statusCode := int(response["status"].(float64))
			log.Printf("Error Creating new %s: %d\n", sitevpnConnectionName, statusCode)
			log.Printf("  %s\n", response["code"])
			log.Printf("  %s\n", response["message"])
		}

		d.SetId("")
		return fmt.Errorf("Error while creating %s: err=%s", sitevpnConnectionName, http_err)
	}

	if resp.StatusCode >= 300 {
		d.SetId("")
		return fmt.Errorf("Error while creating %s: code=%v", sitevpnConnectionName, resp.StatusCode)
	}

	loc := resp.Header.Get("location")
	u, err := url.Parse(loc)
	if err != nil {
		return fmt.Errorf("Error when decoding Connection ID")
	}

	id := filepath.Base(u.Path)
	d.SetId(id)

	if id == "" {
		log.Printf("Error when decoding location header")
		return fmt.Errorf("Error when decoding Connection ID")
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
			log.Printf("Error Response while reading %s: code=%v", sitevpnConnectionName, resp.StatusCode)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading data for %s: %s", sitevpnConnectionName, err)
	}

	if resp.StatusCode >= 300 {
		fmt.Errorf("Error Response while reading %s: code=%v", sitevpnConnectionName, resp.StatusCode)
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
		return fmt.Errorf("Error setting customer networks for %s %s: %s", sitevpnConnectionName, d.Id(), err)
	}

	if err := d.Set("location_href", conn.Location.Href); err != nil {
		return fmt.Errorf("Error setting location for %s %s: %s", sitevpnConnectionName, d.Id(), err)
	}
	if err := d.Set("network_href", conn.Network.Href); err != nil {
		return fmt.Errorf("Error setting network for %s %s: %s", sitevpnConnectionName, d.Id(), err)
	}

	d.Set("auth_type", conn.AuthType)
	d.Set("enable_bgp_password", conn.EnableBGPPassword)
	d.Set("ike_version", conn.IkeVersion)

	if conn.IkeVersion == "V1" {
		if err := d.Set("ike_config", []map[string]interface{}{
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
			return fmt.Errorf("Error setting IKE V1 Configuration for %s %s: %s", sitevpnConnectionName, d.Id(), err)
		}
	}

	if conn.IkeVersion == "V2" {
		if err := d.Set("ike_config", []map[string]interface{}{
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
			return fmt.Errorf("Error setting IKE V2 Configuration for %s %s: %s", sitevpnConnectionName, d.Id(), err)
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
		return fmt.Errorf("Error setting traffics selectors for %s %s: %s", sitevpnConnectionName, d.Id(), err)
	}

	return nil
}

func resourceSiteVPNConnectionUpdate(d *schema.ResourceData, m interface{}) error {

	c := expandSiteVPNConnection(d)

	d.Partial(true)

	sess := m.(*session.Session)
	ctx := sess.GetSessionContext()

	if d.HasChange("name") {
		c.Name = d.Get("name").(string)
		d.SetPartial("name")
	}

	if d.HasChange("description") {
		c.Description = d.Get("description").(string)
		d.SetPartial("description")
	}

	if d.HasChange("speed") {
		c.Speed = int32(d.Get("speed").(int))
		d.SetPartial("speed")
	}

	if d.HasChange("customer_networks") {
		c.CustomerNetworks = ExpandCustomerNetworks(d)
	}

	if d.HasChange("nat_config") {
		c.Nat = ExpandNATConfiguration(d)
	}

	if d.HasChange("billing_term") {
		c.BillingTerm = d.Get("billing_term").(string)
	}

	opts := client.UpdateConnectionOpts{
		Body: optional.NewInterface(c),
	}

	_, resp, err := sess.Client.ConnectionsApi.UpdateConnection(
		ctx,
		d.Id(),
		&opts,
	)

	if err != nil {

		json_response := string(err.(client.GenericSwaggerError).Body()[:])
		response, err := structure.ExpandJsonFromString(json_response)
		if err != nil {
			log.Printf("Error updating %s: %v", sitevpnConnectionName, err)
		} else {
			statusCode := int(response["status"].(float64))
			log.Printf("Error updating %s: %d\n", sitevpnConnectionName, statusCode)
			log.Printf("  %s\n", response["code"])
			log.Printf("  %s\n", response["message"])
		}

		return fmt.Errorf("Error while updating %s: err=%s", sitevpnConnectionName, err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while updating %s: code=%v", sitevpnConnectionName, resp.StatusCode)
	}

	d.Partial(false)
	return resourceSiteVPNConnectionRead(d, m)
}

func resourceSiteVPNConnectionDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteConnection(sitevpnConnectionName, d, m)
}
