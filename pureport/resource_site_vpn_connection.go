// Package main provides SiteVPNConnection resource
package pureport

import (
	"fmt"
	"log"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pureport/pureport-sdk-go/pureport/session"
	"github.com/pureport/pureport-sdk-go/pureport/swagger"
)

func resourceSiteVPNConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"auth_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "psk",
			ValidateFunc: validation.StringInSlice([]string{"psk"}, true),
		},
		"enable_bgp_password": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"ike_version": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"1", "2"}, true),
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
			Required: true,
		},
		"primary_key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"routing_type": {
			Type:     schema.TypeString,
			Required: true,
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

func resourceSiteVPNConnectionCreate(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)

	// Generic Connection values
	network := d.Get("network").([]interface{})
	speed := d.Get("speed").(int)
	name := d.Get("name").(string)
	location := d.Get("location").([]interface{})
	billingTerm := d.Get("billing_term").(string)

	// SiteVPN specific values
	awsAccountId := d.Get("aws_account_id").(string)
	awsRegion := d.Get("aws_region").(string)

	// Create the body of the request
	connection := swagger.AwsDirectConnectConnection{
		Name:  name,
		Speed: int32(speed),
		Location: &swagger.Link{
			Id:   location[0].(map[string]interface{})["id"].(string),
			Href: location[0].(map[string]interface{})["href"].(string),
		},
		Network: &swagger.Link{
			Id:   network[0].(map[string]interface{})["id"].(string),
			Href: network[0].(map[string]interface{})["href"].(string),
		},
		AwsAccountId: awsAccountId,
		AwsRegion:    awsRegion,
		BillingTerm:  billingTerm,
	}

	// Generic Optionals
	if customerNetworks, ok := d.GetOk("customer_networks"); ok {
		for _, cn := range customerNetworks.([]map[string]string) {

			new := swagger.CustomerNetwork{
				Name:    cn["name"],
				Address: cn["Address"],
			}

			connection.CustomerNetworks = append(connection.CustomerNetworks, new)
		}
	}

	if description, ok := d.GetOk("description"); ok {
		connection.Description = description.(string)
	}

	if highAvailability, ok := d.GetOk("high_availability"); ok {
		connection.HighAvailability = highAvailability.(bool)
	}

	if natConfig, ok := d.GetOk("nat_config"); ok {

		config := natConfig.(map[string]interface{})
		connection.Nat = &swagger.NatConfig{
			Enabled: config["enabled"].(bool),
		}

		for _, m := range config["mappings"].([]map[string]string) {

			new := swagger.NatMapping{
				NativeCidr: m["native_cidr"],
			}

			connection.Nat.Mappings = append(connection.Nat.Mappings, new)
		}
	}

	// SiteVPN Optionals
	if cloudServices, ok := d.GetOk("cloud_services"); ok {
		for _, cs := range cloudServices.([]map[string]string) {

			new := swagger.Link{
				Id:   cs["id"],
				Href: cs["href"],
			}

			connection.CloudServices = append(connection.CloudServices, new)
		}
	}

	if peeringType, ok := d.GetOk("peering"); ok {
		connection.Peering = &swagger.PeeringConfiguration{
			Type_: peeringType.(string),
		}
	} else {
		connection.Peering = &swagger.PeeringConfiguration{
			Type_: "",
		}
	}

	connection.Type_ = "SiteVPN_DIRECT_CONNECT"

	ctx := sess.GetSessionContext()

	opts := swagger.AddConnectionOpts{
		Body: optional.NewInterface(connection),
	}

	id, resp, err := sess.Client.ConnectionsApi.AddConnection(
		ctx,
		network[0].(map[string]interface{})["id"].(string),
		&opts,
	)

	if err != nil {
		log.Printf("[Error] Error Creating new SiteVPN Connection: %v", err)
		d.SetId("")
		return nil
	}

	if resp.StatusCode >= 300 {
		log.Printf("[Error] Error Response while creating new SiteVPN Connection: code=%v", resp.StatusCode)
		d.SetId("")
		return nil
	}

	d.SetId(id)

	return resourceSiteVPNConnectionRead(d, m)
}

func resourceSiteVPNConnectionRead(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)
	connectionId := d.Id()
	ctx := sess.GetSessionContext()

	c, resp, err := sess.Client.ConnectionsApi.Get11(ctx, connectionId)
	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf("[Error] Error Response while reading SiteVPN Connection: code=%v", resp.StatusCode)
			d.SetId("")
		}
		return fmt.Errorf("[Error] Error reading data for SiteVPN Connection: %s", err)
	}

	if resp.StatusCode >= 300 {
		fmt.Errorf("[Error] Error Response while reading SiteVPN Connection: code=%v", resp.StatusCode)
	}

	conn := c.(swagger.AwsDirectConnectConnection)
	d.Set("aws_account_id", conn.AwsAccountId)
	d.Set("aws_region", conn.AwsRegion)

	var cloudServices []map[string]string
	for _, cs := range conn.CloudServices {
		cloudServices = append(cloudServices, map[string]string{
			"id":   cs.Id,
			"href": cs.Href,
		})
	}
	d.Set("cloud_services", cloudServices)
	d.Set("peering", conn.Peering.Type_)

	var customerNetworks []map[string]string
	for _, cn := range conn.CustomerNetworks {
		customerNetworks = append(customerNetworks, map[string]string{
			"name":    cn.Name,
			"address": cn.Address,
		})
	}
	d.Set("customer_networks", customerNetworks)

	d.Set("description", conn.Description)
	d.Set("high_availability", conn.HighAvailability)
	d.Set("location", map[string]string{
		"id":   conn.Location.Id,
		"href": conn.Location.Href,
	})
	d.Set("network", map[string]string{
		"id":   conn.Network.Id,
		"href": conn.Network.Href,
	})

	return nil
}

func resourceSiteVPNConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceSiteVPNConnectionRead(d, m)
}

func resourceSiteVPNConnectionDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteConnection(d, m)
}
