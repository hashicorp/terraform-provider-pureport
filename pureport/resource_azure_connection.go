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
	azureConnectionName = "Azure Cloud Connection"
)

func resourceAzureConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"service_key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"peering": {
			Type:         schema.TypeString,
			Description:  "The peering configuration to use for this connection Public/Private",
			Default:      "PRIVATE",
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"private", "public"}, true),
		},
	}

	// Add the base items
	for k, v := range getBaseConnectionSchema() {
		connection_schema[k] = v
	}

	return &schema.Resource{
		Create: resourceAzureConnectionCreate,
		Read:   resourceAzureConnectionRead,
		Update: resourceAzureConnectionUpdate,
		Delete: resourceAzureConnectionDelete,

		Schema: connection_schema,
	}
}

func resourceAzureConnectionCreate(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)

	// Generic Connection values
	network := d.Get("network").([]interface{})
	speed := d.Get("speed").(int)
	name := d.Get("name").(string)
	location := d.Get("location").([]interface{})
	billingTerm := d.Get("billing_term").(string)

	// Azure specific values
	serviceKey := d.Get("service_key").(string)

	// Create the body of the request
	connection := client.AzureExpressRouteConnection{
		Type_: "AZURE_EXPRESS_ROUTE",
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
		BillingTerm: billingTerm,
		ServiceKey:  serviceKey,
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

	// Azure Optionals
	connection.Peering = AddPeeringType(d)

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

		json_response := string(err.(client.GenericSwaggerError).Body()[:])
		response, err := structure.ExpandJsonFromString(json_response)
		if err != nil {
			log.Printf("Error Creating new %s: %v", azureConnectionName, err)
		} else {
			log.Printf("Error Creating new %s: %f\n", azureConnectionName, response["status"])
			log.Printf("  %s\n", response["code"])
			log.Printf("  %s\n", response["message"])
		}
		d.SetId("")
		return nil
	}

	if resp.StatusCode >= 300 {
		log.Printf("Error Response while creating new %s: code=%v", azureConnectionName, resp.StatusCode)
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

	return resourceAzureConnectionRead(d, m)
}

func resourceAzureConnectionRead(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)
	connectionId := d.Id()
	ctx := sess.GetSessionContext()

	c, resp, err := sess.Client.ConnectionsApi.GetConnection(ctx, connectionId)
	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf("Error Response while reading %s: code=%v", azureConnectionName, resp.StatusCode)
			d.SetId("")
		}
		return fmt.Errorf("Error reading data for %s: %s", azureConnectionName, err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while reading %s: code=%v", azureConnectionName, resp.StatusCode)
	}

	conn := c.(client.AzureExpressRouteConnection)
	d.Set("service_key", conn.ServiceKey)
	d.Set("peering", conn.Peering.Type_)
	d.Set("speed", conn.Speed)

	var customerNetworks []map[string]string
	for _, cn := range conn.CustomerNetworks {
		customerNetworks = append(customerNetworks, map[string]string{
			"name":    cn.Name,
			"address": cn.Address,
		})
	}
	if err := d.Set("customer_networks", customerNetworks); err != nil {
		return fmt.Errorf("Error setting customer networks for %s %s: %s", azureConnectionName, d.Id(), err)
	}

	d.Set("description", conn.Description)
	d.Set("high_availability", conn.HighAvailability)

	if err := d.Set("location", []map[string]string{
		{
			"id":   conn.Location.Id,
			"href": conn.Location.Href,
		},
	}); err != nil {
		return fmt.Errorf("Error setting location for %s %s: %s", azureConnectionName, d.Id(), err)
	}

	if err := d.Set("network", []map[string]string{
		{
			"id":   conn.Network.Id,
			"href": conn.Network.Href,
		},
	}); err != nil {
		return fmt.Errorf("Error setting location for %s %s: %s", azureConnectionName, d.Id(), err)
	}

	return nil
}

func resourceAzureConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceAzureConnectionRead(d, m)
}

func resourceAzureConnectionDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteConnection(azureConnectionName, d, m)
}
