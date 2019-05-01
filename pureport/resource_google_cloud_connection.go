package pureport

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/pureport/pureport-sdk-go/pureport/session"
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
	connection := client.GoogleCloudInterconnectConnection{
		Type_: "GOOGLE_CLOUD_INTERCONNECT",
		Name:  name,
		Speed: int32(speed),
		Location: &client.Link{
			Id:   location[0].(map[string]interface{})["id"].(string),
			Href: location[0].(map[string]interface{})["href"].(string),
		},
		Network: &client.Link{
			Id:   network[0].(map[string]interface{})["id"].(string),
			Href: network[0].(map[string]interface{})["href"].(string),
		},
		BillingTerm:       billingTerm,
		PrimaryPairingKey: primaryPairingKey,
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

	// Google Optionals
	if secondaryPairingKey, ok := d.GetOk("secondary_pairing_key"); ok {
		connection.SecondaryPairingKey = secondaryPairingKey.(string)
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
		log.Printf("Error Creating new Google Cloud Connection: %v", err)
		d.SetId("")
		return nil
	}

	if resp.StatusCode >= 300 {
		log.Printf("Error Response while creating new Google Cloud Connection: code=%v", resp.StatusCode)
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

	return resourceGoogleCloudConnectionRead(d, m)
}

func resourceGoogleCloudConnectionRead(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)
	connectionId := d.Id()
	ctx := sess.GetSessionContext()

	c, resp, err := sess.Client.ConnectionsApi.GetConnection(ctx, connectionId)
	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf("Error Response while reading Google Cloud Connection: code=%v", resp.StatusCode)
			d.SetId("")
		}
		return fmt.Errorf("Error reading data for Google Cloud Connection: %s", err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while reading Google Cloud Connection: code=%v", resp.StatusCode)
	}

	conn := c.(client.GoogleCloudInterconnectConnection)
	d.Set("speed", conn.Speed)

	var customerNetworks []map[string]string
	for _, cn := range conn.CustomerNetworks {
		customerNetworks = append(customerNetworks, map[string]string{
			"name":    cn.Name,
			"address": cn.Address,
		})
	}
	if err := d.Set("customer_networks", customerNetworks); err != nil {
		return fmt.Errorf("Error setting customer networks for Google Cloud Connection %s: %s", d.Id(), err)
	}

	d.Set("description", conn.Description)
	d.Set("high_availability", conn.HighAvailability)
	if err := d.Set("location", []map[string]string{
		{
			"id":   conn.Location.Id,
			"href": conn.Location.Href,
		},
	}); err != nil {
		return fmt.Errorf("Error setting location for Google Cloud Connection %s: %s", d.Id(), err)
	}

	if err := d.Set("network", []map[string]string{
		{
			"id":   conn.Network.Id,
			"href": conn.Network.Href,
		},
	}); err != nil {
		return fmt.Errorf("Error setting network for Google Cloud Connection %s: %s", d.Id(), err)
	}

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
