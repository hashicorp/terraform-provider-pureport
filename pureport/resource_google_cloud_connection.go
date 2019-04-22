// Package pureport provides ...
package pureport

import (
	"fmt"
	"log"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pureport/pureport-sdk-go/pureport/session"
	"github.com/pureport/pureport-sdk-go/pureport/swagger"
)

func resourceGoogleCloudConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"primary_pairing_key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"secondary_pairing_key": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	// Add the base items
	for k, v := range getBaseConnectionSchema() {
		connection_schema[k] = v
	}

	return &schema.Resource{
		Create: resourceGoogleCloudConnectionCreate,
		Read:   resourceGoogleCloudConnectionRead,
		Update: resourceGoogleCloudConnectionUpdate,
		Delete: resourceGoogleCloudConnectionDelete,

		Schema: connection_schema,
	}
}

func resourceGoogleCloudConnectionCreate(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)

	// Generic Connection values
	network := d.Get("network").([]interface{})
	speed := d.Get("speed").(int)
	name := d.Get("name").(string)
	location := d.Get("location").([]interface{})
	billingTerm := d.Get("billing_term").(string)

	// Google specific values
	primaryPairingKey := d.Get("primary_pairing_key").(string)

	// Create the body of the request
	connection := swagger.GoogleCloudInterconnectConnection{
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
		BillingTerm:       billingTerm,
		PrimaryPairingKey: primaryPairingKey,
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

	// Google Optionals
	if secondaryPairingKey, ok := d.GetOk("secondary_pairing_key"); ok {
		connection.SecondaryPairingKey = secondaryPairingKey.(string)
	}

	connection.Type_ = "GOOGLE_CLOUD_INTERCONNECT"

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
		log.Printf("[Error] Error Creating new Google Cloud Connection: %v", err)
		d.SetId("")
		return nil
	}

	if resp.StatusCode >= 300 {
		log.Printf("[Error] Error Response while creating new Google Cloud Connection: code=%v", resp.StatusCode)
		d.SetId("")
		return nil
	}

	d.SetId(id)
	return resourceGoogleCloudConnectionRead(d, m)
}

func resourceGoogleCloudConnectionRead(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)
	connectionId := d.Id()
	ctx := sess.GetSessionContext()

	c, resp, err := sess.Client.ConnectionsApi.Get11(ctx, connectionId)
	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf("[Error] Error Response while reading Google Cloud Connection: code=%v", resp.StatusCode)
			d.SetId("")
		}
		return fmt.Errorf("[Error] Error reading data for Google Cloud Connection: %s", err)
	}

	if resp.StatusCode >= 300 {
		fmt.Errorf("[Error] Error Response while reading Google Cloud Connection: code=%v", resp.StatusCode)
	}

	conn := c.(swagger.GoogleCloudInterconnectConnection)

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

	d.Set("primary_pairing_key", conn.PrimaryPairingKey)
	d.Set("secondary_pairing_key", conn.SecondaryPairingKey)

	return nil
}

func resourceGoogleCloudConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceGoogleCloudConnectionRead(d, m)
}

func resourceGoogleCloudConnectionDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteConnection(d, m)
}
