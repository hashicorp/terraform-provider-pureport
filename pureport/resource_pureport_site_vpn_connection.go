package pureport

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/configuration"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/connection"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/tags"
)

func resourceSiteVPNConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"speed": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntInSlice([]int{50, 100, 200, 300, 400, 500, 1000, 10000}),
		},
		"ike_version": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"V1", "V2"}, true),
			StateFunc: func(val interface{}) string {
				return strings.ToUpper(val.(string))
			},
		},
		"primary_customer_router_ip": {
			Type:     schema.TypeString,
			Required: true,
		},
		"routing_type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"ROUTE_BASED_BGP", "ROUTE_BASED_STATIC", "POLICY_BASED"}, true),
			StateFunc: func(val interface{}) string {
				return strings.ToUpper(val.(string))
			},
		},
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
		},
		"ike_config": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
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
		"primary_key": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"secondary_customer_router_ip": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"secondary_key": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"traffic_selectors": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"customer_side": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.CIDRNetwork(8, 32),
					},
					"pureport_side": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.CIDRNetwork(8, 32),
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
	for k, v := range connection.GetBaseResourceConnectionSchema() {
		connection_schema[k] = v
	}

	return &schema.Resource{
		Create: resourceSiteVPNConnectionCreate,
		Read:   resourceSiteVPNConnectionRead,
		Update: resourceSiteVPNConnectionUpdate,
		Delete: resourceSiteVPNConnectionDelete,

		Schema: connection_schema,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},
	}
}

func expandTrafficSelectorMappings(d *schema.ResourceData) []client.TrafficSelectorMapping {

	if data, ok := d.GetOk("traffic_selectors"); ok {

		mappings := []client.TrafficSelectorMapping{}

		for _, i := range data.(*schema.Set).List() {

			m := i.(map[string]interface{})

			ts := client.TrafficSelectorMapping{
				CustomerSide: m["customer_side"].(string),
				PureportSide: m["pureport_side"].(string),
			}

			mappings = append(mappings, ts)
		}

		return mappings
	}

	return nil
}

func expandIkeVersion1(d *schema.ResourceData) *client.Ikev1Config {

	config := &client.Ikev1Config{}

	if data, ok := d.GetOk("ike_config"); ok {

		tmp_config := data.([]interface{})
		raw_config := tmp_config[0].(map[string]interface{})

		tmp_esp := raw_config["esp"].([]interface{})
		esp := tmp_esp[0].(map[string]interface{})

		tmp_ike := raw_config["ike"].([]interface{})
		ike := tmp_ike[0].(map[string]interface{})

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

		tmp_config := data.([]interface{})
		raw_config := tmp_config[0].(map[string]interface{})

		tmp_esp := raw_config["esp"].([]interface{})
		esp := tmp_esp[0].(map[string]interface{})

		tmp_ike := raw_config["ike"].([]interface{})
		ike := tmp_ike[0].(map[string]interface{})

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
	c.CustomerNetworks = connection.ExpandCustomerNetworks(d)
	c.Nat = connection.ExpandNATConfiguration(d)

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

	if t, ok := d.GetOk("tags"); ok {
		c.Tags = tags.FilterTags(t.(map[string]interface{}))
	}

	return c
}

func resourceSiteVPNConnectionCreate(d *schema.ResourceData, m interface{}) error {

	c := expandSiteVPNConnection(d)

	config := m.(*configuration.Config)
	ctx := config.Session.GetSessionContext()

	opts := client.AddConnectionOpts{
		Body: optional.NewInterface(c),
	}

	_, resp, err := config.Session.Client.ConnectionsApi.AddConnection(
		ctx,
		filepath.Base(c.Network.Href),
		&opts,
	)

	if err != nil {

		http_err := err
		json_response := string(err.(client.GenericSwaggerError).Body()[:])
		response, err := structure.ExpandJsonFromString(json_response)
		if err != nil {
			log.Printf("Error Creating new %s: %v", connection.SiteVPNConnectionName, err)
		} else {
			statusCode := int(response["status"].(float64))
			log.Printf("Error Creating new %s: %d\n", connection.SiteVPNConnectionName, statusCode)
			log.Printf("  %s\n", response["code"])
			log.Printf("  %s\n", response["message"])
		}

		d.SetId("")
		return fmt.Errorf("Error while creating %s: err=%s", connection.SiteVPNConnectionName, http_err)
	}

	if resp.StatusCode >= 300 {
		d.SetId("")
		return fmt.Errorf("Error while creating %s: code=%v", connection.SiteVPNConnectionName, resp.StatusCode)
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

	if err := connection.WaitForConnection(connection.SiteVPNConnectionName, d, m); err != nil {
		return fmt.Errorf("Error waiting for %s: err=%s", connection.SiteVPNConnectionName, err)
	}

	return resourceSiteVPNConnectionRead(d, m)
}

func resourceSiteVPNConnectionRead(d *schema.ResourceData, m interface{}) error {

	config := m.(*configuration.Config)
	connectionId := d.Id()
	ctx := config.Session.GetSessionContext()

	c, resp, err := config.Session.Client.ConnectionsApi.GetConnection(ctx, connectionId)
	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf("Error Response while reading %s: code=%v", connection.SiteVPNConnectionName, resp.StatusCode)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading data for %s: %s", connection.SiteVPNConnectionName, err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while reading %s: code=%v", connection.SiteVPNConnectionName, resp.StatusCode)
	}

	conn := c.(client.SiteIpSecVpnConnection)
	d.Set("auth_type", conn.AuthType)
	d.Set("description", conn.Description)
	d.Set("enable_bgp_password", conn.EnableBGPPassword)
	d.Set("high_availability", conn.HighAvailability)
	d.Set("href", conn.Href)
	d.Set("ike_version", conn.IkeVersion)
	d.Set("name", conn.Name)
	d.Set("primary_customer_router_ip", conn.PrimaryCustomerRouterIP)
	d.Set("primary_key", conn.PrimaryKey)
	d.Set("routing_type", conn.RoutingType)
	d.Set("secondary_customer_router_ip", conn.SecondaryCustomerRouterIP)
	d.Set("secondary_key", conn.SecondaryKey)
	d.Set("speed", conn.Speed)
	d.Set("state", conn.State)

	// Add Gateway information
	var gateways []map[string]interface{}
	if g := conn.PrimaryGateway; g != nil {
		gateways = append(gateways, connection.FlattenVpnGateway(g))
	}
	if g := conn.SecondaryGateway; g != nil {
		gateways = append(gateways, connection.FlattenVpnGateway(g))
	}
	if err := d.Set("gateways", gateways); err != nil {
		return fmt.Errorf("Error setting gateway information for %s %s: %s", connection.SiteVPNConnectionName, d.Id(), err)
	}

	if err := d.Set("customer_networks", connection.FlattenCustomerNetworks(conn.CustomerNetworks)); err != nil {
		return fmt.Errorf("Error setting customer networks for %s %s: %s", connection.SiteVPNConnectionName, d.Id(), err)
	}

	if err := d.Set("nat_config", connection.FlattenNatConfig(conn.Nat)); err != nil {
		return fmt.Errorf("Error setting NAT Configuration for %s %s: %s", connection.SiteVPNConnectionName, d.Id(), err)
	}

	if err := d.Set("location_href", conn.Location.Href); err != nil {
		return fmt.Errorf("Error setting location for %s %s: %s", connection.SiteVPNConnectionName, d.Id(), err)
	}
	if err := d.Set("network_href", conn.Network.Href); err != nil {
		return fmt.Errorf("Error setting network for %s %s: %s", connection.SiteVPNConnectionName, d.Id(), err)
	}

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
			return fmt.Errorf("Error setting IKE V1 Configuration for %s %s: %s", connection.SiteVPNConnectionName, d.Id(), err)
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
			return fmt.Errorf("Error setting IKE V2 Configuration for %s %s: %s", connection.SiteVPNConnectionName, d.Id(), err)
		}
	}

	trafficSelectors := []map[string]string{}
	for _, v := range conn.TrafficSelectors {
		trafficSelectors = append(trafficSelectors, map[string]string{
			"customer_side": v.CustomerSide,
			"pureport_side": v.PureportSide,
		})
	}

	if err := d.Set("traffic_selectors", trafficSelectors); err != nil {
		return fmt.Errorf("Error setting traffics selectors for %s %s: %s", connection.SiteVPNConnectionName, d.Id(), err)
	}

	if err := d.Set("tags", conn.Tags); err != nil {
		return fmt.Errorf("Error setting tags for %s %s: %s", connection.SiteVPNConnectionName, d.Id(), err)
	}

	return nil
}

func resourceSiteVPNConnectionUpdate(d *schema.ResourceData, m interface{}) error {

	c := expandSiteVPNConnection(d)

	d.Partial(true)

	config := m.(*configuration.Config)
	ctx := config.Session.GetSessionContext()

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
		c.CustomerNetworks = connection.ExpandCustomerNetworks(d)
	}

	if d.HasChange("nat_config") {
		c.Nat = connection.ExpandNATConfiguration(d)
	}

	if d.HasChange("billing_term") {
		c.BillingTerm = d.Get("billing_term").(string)
	}

	if d.HasChange("enable_bgp_password") {
		c.EnableBGPPassword = d.Get("enable_bgp_password").(bool)
	}

	if d.HasChange("ike_version") {
		c.IkeVersion = d.Get("ike_version").(string)

		if c.IkeVersion == "V1" {
			c.IkeV1 = expandIkeVersion1(d)
			c.IkeV2 = nil
		} else {
			c.IkeV2 = expandIkeVersion2(d)
			c.IkeV1 = nil
		}
	}

	if d.HasChange("ike_config") {
		if c.IkeVersion == "V1" {
			c.IkeV1 = expandIkeVersion1(d)
			c.IkeV2 = nil
		} else {
			c.IkeV2 = expandIkeVersion2(d)
			c.IkeV1 = nil
		}
	}

	if d.HasChange("primary_customer_router_ip") {
		c.PrimaryCustomerRouterIP = d.Get("primary_customer_router_ip").(string)
	}

	if d.HasChange("primary_key") {
		c.PrimaryKey = d.Get("primary_key").(string)
	}

	if d.HasChange("routing_type") {
		c.RoutingType = d.Get("routing_type").(string)
	}

	if d.HasChange("secondary_customer_router_ip") {
		c.SecondaryCustomerRouterIP = d.Get("secondary_customer_router_ip").(string)
	}

	if d.HasChange("secondary_key") {
		c.SecondaryKey = d.Get("secondary_key").(string)
	}

	if d.HasChange("traffic_selectors") {
		c.TrafficSelectors = expandTrafficSelectorMappings(d)
	}

	if d.HasChange("tags") {
		_, nraw := d.GetChange("tags")
		c.Tags = tags.FilterTags(nraw.(map[string]interface{}))
	}

	opts := client.UpdateConnectionOpts{
		Body: optional.NewInterface(c),
	}

	_, resp, err := config.Session.Client.ConnectionsApi.UpdateConnection(
		ctx,
		d.Id(),
		&opts,
	)

	if err != nil {

		if swerr, ok := err.(client.GenericSwaggerError); ok {

			json_response := string(swerr.Body()[:])
			response, jerr := structure.ExpandJsonFromString(json_response)

			if jerr == nil {
				statusCode := int(response["status"].(float64))
				log.Printf("Error updating %s: %d\n", connection.SiteVPNConnectionName, statusCode)
				log.Printf("  %s\n", response["code"])
				log.Printf("  %s\n", response["message"])
			}
		}

		return fmt.Errorf("Error while updating %s: err=%s", connection.SiteVPNConnectionName, err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while updating %s: code=%v", connection.SiteVPNConnectionName, resp.StatusCode)
	}

	if err := connection.WaitForConnection(connection.SiteVPNConnectionName, d, m); err != nil {
		return fmt.Errorf("Error waiting for %s: err=%s", connection.SiteVPNConnectionName, err)
	}

	d.Partial(false)

	return resourceSiteVPNConnectionRead(d, m)
}

func resourceSiteVPNConnectionDelete(d *schema.ResourceData, m interface{}) error {
	return connection.DeleteConnection(connection.SiteVPNConnectionName, d, m)
}
