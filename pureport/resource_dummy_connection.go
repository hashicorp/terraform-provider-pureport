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
	dummyConnectionName = "Dummy Cloud Connection"
)

func resourceDummyConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"peering_type": {
			Type:         schema.TypeString,
			Description:  "The peering type to use for this connection Public/Private",
			Default:      "PRIVATE",
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"private", "public"}, true),
			Optional:     true,
		},
	}

	// Add the base items
	for k, v := range getBaseConnectionSchema() {
		connection_schema[k] = v
	}

	return &schema.Resource{
		Create: resourceDummyConnectionCreate,
		Read:   resourceDummyConnectionRead,
		Update: resourceDummyConnectionUpdate,
		Delete: resourceDummyConnectionDelete,

		Schema: connection_schema,
	}
}

func expandDummyConnection(d *schema.ResourceData) client.DummyConnection {

	// Generic Connection values
	network := d.Get("network").([]interface{})
	speed := d.Get("speed").(int)
	name := d.Get("name").(string)
	billingTerm := d.Get("billing_term").(string)
	location_href := d.Get("location_href").(string)

	// Create the body of the request
	c := client.DummyConnection{
		Type_: "DUMMY",
		Name:  name,
		Speed: int32(speed),
		Location: &client.Link{
			Href: location_href,
		},
		Network: &client.Link{
			Id:   network[0].(map[string]interface{})["id"].(string),
			Href: network[0].(map[string]interface{})["href"].(string),
		},
		BillingTerm: billingTerm,
	}

	// Generic Optionals
	c.CustomerNetworks = ExpandCustomerNetworks(d)
	c.Nat = ExpandNATConfiguration(d)
	c.Peering = ExpandPeeringType(d)

	if description, ok := d.GetOk("description"); ok {
		c.Description = description.(string)
	}

	if highAvailability, ok := d.GetOk("high_availability"); ok {
		c.HighAvailability = highAvailability.(bool)
	}

	return c
}

func resourceDummyConnectionCreate(d *schema.ResourceData, m interface{}) error {

	connection := expandDummyConnection(d)

	sess := m.(*session.Session)
	ctx := sess.GetSessionContext()

	opts := client.AddConnectionOpts{
		Body: optional.NewInterface(connection),
	}

	resp, err := sess.Client.ConnectionsApi.AddConnection(
		ctx,
		connection.Network.Id,
		&opts,
	)

	if err != nil {

		json_response := string(err.(client.GenericSwaggerError).Body()[:])
		response, err := structure.ExpandJsonFromString(json_response)
		if err != nil {
			log.Printf("Error Creating new %s: %v", dummyConnectionName, err)
		} else {
			statusCode := int(response["status"].(float64))
			log.Printf("Error Creating new %s: %d\n", dummyConnectionName, statusCode)
			log.Printf("  %s\n", response["code"])
			log.Printf("  %s\n", response["message"])
		}

		d.SetId("")
		return fmt.Errorf("Error while creating %s: err=%s", dummyConnectionName, err)
	}

	if resp.StatusCode >= 300 {
		d.SetId("")
		return fmt.Errorf("Error while creating %s: code=%v", dummyConnectionName, resp.StatusCode)
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
		return fmt.Errorf("Error decoding Connection ID")
	}

	return resourceDummyConnectionRead(d, m)
}

func resourceDummyConnectionRead(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)
	connectionId := d.Id()
	ctx := sess.GetSessionContext()

	c, resp, err := sess.Client.ConnectionsApi.GetConnection(ctx, connectionId)
	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf("Error Response while reading %s: code=%v", dummyConnectionName, resp.StatusCode)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading data for %s: %s", dummyConnectionName, err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while reading %s: code=%v", dummyConnectionName, resp.StatusCode)
	}

	conn := c.(client.DummyConnection)

	d.Set("peering_type", conn.Peering.Type_)
	d.Set("speed", conn.Speed)

	var customerNetworks []map[string]string
	for _, cn := range conn.CustomerNetworks {
		customerNetworks = append(customerNetworks, map[string]string{
			"name":    cn.Name,
			"address": cn.Address,
		})
	}
	if err := d.Set("customer_networks", customerNetworks); err != nil {
		return fmt.Errorf("Error setting customer networks for %s %s: %s", dummyConnectionName, d.Id(), err)
	}

	d.Set("description", conn.Description)
	d.Set("high_availability", conn.HighAvailability)

	if err := d.Set("location_href", conn.Location.Href); err != nil {
		return fmt.Errorf("Error setting location for %s %s: %s", dummyConnectionName, d.Id(), err)
	}

	if err := d.Set("network", []map[string]string{
		{
			"id":   conn.Network.Id,
			"href": conn.Network.Href,
		},
	}); err != nil {
		return fmt.Errorf("Error setting network for %s %s: %s", dummyConnectionName, d.Id(), err)
	}

	return nil
}

func resourceDummyConnectionUpdate(d *schema.ResourceData, m interface{}) error {

	c := expandDummyConnection(d)

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
			log.Printf("Error Creating new %s: %v", dummyConnectionName, err)
		} else {
			statusCode := int(response["status"].(float64))
			log.Printf("Error updating %s: %d\n", dummyConnectionName, statusCode)
			log.Printf("  %s\n", response["code"])
			log.Printf("  %s\n", response["message"])
		}

		return fmt.Errorf("Error while updating %s: err=%s", dummyConnectionName, err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while updating %s: code=%v", dummyConnectionName, resp.StatusCode)
	}

	d.Partial(false)
	return resourceDummyConnectionRead(d, m)
}

func resourceDummyConnectionDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteConnection(dummyConnectionName, d, m)
}
