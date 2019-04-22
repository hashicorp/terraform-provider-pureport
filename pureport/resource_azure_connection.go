// Package pureport provides ...
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

func resourceAzureConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"service_key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"peering": {
			Type:         schema.TypeString,
			Description:  "The peering configuration to use for this connection Public/Private",
			Default:      "Private",
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
	connection := swagger.AzureExpressRouteConnection{
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
		BillingTerm: billingTerm,
		ServiceKey:  serviceKey,
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

	// Azure Optionals
	if peeringType, ok := d.GetOk("peering"); ok {
		connection.Peering = &swagger.PeeringConfiguration{
			Type_: peeringType.(string),
		}
	} else {
		connection.Peering = &swagger.PeeringConfiguration{
			Type_: "",
		}
	}

	connection.Type_ = "AZURE_EXPRESS_ROUTE"

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
		log.Printf("[Error] Error Creating new Azure Connection: %v", err)
		d.SetId("")
		return nil
	}

	if resp.StatusCode >= 300 {
		log.Printf("[Error] Error Response while creating new Azure Connection: code=%v", resp.StatusCode)
		d.SetId("")
		return nil
	}

	d.SetId(id)
	return resourceAzureConnectionRead(d, m)
}

func resourceAzureConnectionRead(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)
	connectionId := d.Id()
	ctx := sess.GetSessionContext()

	c, resp, err := sess.Client.ConnectionsApi.Get11(ctx, connectionId)
	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf("[Error] Error Response while reading Azure Connection: code=%v", resp.StatusCode)
			d.SetId("")
		}
		return fmt.Errorf("[Error] Error reading data for Azure Connection: %s", err)
	}

	if resp.StatusCode >= 300 {
		fmt.Errorf("[Error] Error Response while reading Azure Connection: code=%v", resp.StatusCode)
	}

	conn := c.(swagger.AzureExpressRouteConnection)
	d.Set("service_key", conn.ServiceKey)

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

func resourceAzureConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceAzureConnectionRead(d, m)
}

func resourceAzureConnectionDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteConnection(d, m)
}
