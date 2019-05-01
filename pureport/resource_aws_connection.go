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

func resourceAWSConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"aws_account_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"aws_region": {
			Type:     schema.TypeString,
			Required: true,
		},
		"cloud_services": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeString,
						Required: true,
					},
					"href": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
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
		Create: resourceAWSConnectionCreate,
		Read:   resourceAWSConnectionRead,
		Update: resourceAWSConnectionUpdate,
		Delete: resourceAWSConnectionDelete,

		Schema: connection_schema,
	}
}

func resourceAWSConnectionCreate(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)

	// Generic Connection values
	network := d.Get("network").([]interface{})
	speed := d.Get("speed").(int)
	name := d.Get("name").(string)
	location := d.Get("location").([]interface{})
	billingTerm := d.Get("billing_term").(string)

	// Create the body of the request
	connection := client.AwsDirectConnectConnection{
		Type_: "AWS_DIRECT_CONNECT",
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
		AwsAccountId: d.Get("aws_account_id").(string),
		AwsRegion:    d.Get("aws_region").(string),
		BillingTerm:  billingTerm,
	}

	// Generic Optionals
	connection.CustomerNetworks = AddCustomerNetworks(d)
	connection.Nat = AddNATConfiguration(d)
	connection.CloudServices = AddCloudServices(d)
	connection.Peering = AddPeeringType(d)

	if description, ok := d.GetOk("description"); ok {
		connection.Description = description.(string)
	}

	if highAvailability, ok := d.GetOk("high_availability"); ok {
		connection.HighAvailability = highAvailability.(bool)
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
		log.Printf("Error Creating new AWS Connection: %v", err)
		d.SetId("")
		return nil
	}

	if resp.StatusCode >= 300 {
		log.Printf("Error Response while creating new AWS Connection: code=%v", resp.StatusCode)
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

	return resourceAWSConnectionRead(d, m)
}

func resourceAWSConnectionRead(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)
	connectionId := d.Id()
	ctx := sess.GetSessionContext()

	c, resp, err := sess.Client.ConnectionsApi.GetConnection(ctx, connectionId)
	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf("Error Response while reading AWS Connection: code=%v", resp.StatusCode)
			d.SetId("")
		}
		return fmt.Errorf("Error reading data for AWS Connection: %s", err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while reading AWS Connection: code=%v", resp.StatusCode)
	}

	conn := c.(client.AwsDirectConnectConnection)
	d.Set("aws_account_id", conn.AwsAccountId)
	d.Set("aws_region", conn.AwsRegion)

	var cloudServices []map[string]string
	for _, cs := range conn.CloudServices {
		cloudServices = append(cloudServices, map[string]string{
			"id":   cs.Id,
			"href": cs.Href,
		})
	}
	if err := d.Set("cloud_services", cloudServices); err != nil {
		return fmt.Errorf("Error setting cloud services for AWS Cloud Connection %s: %s", d.Id(), err)
	}

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
		return fmt.Errorf("Error setting customer networks for AWS Cloud Connection %s: %s", d.Id(), err)
	}

	d.Set("description", conn.Description)
	d.Set("high_availability", conn.HighAvailability)

	if err := d.Set("location", []map[string]string{
		{
			"id":   conn.Location.Id,
			"href": conn.Location.Href,
		},
	}); err != nil {
		return fmt.Errorf("Error setting location for AWS Cloud Connection %s: %s", d.Id(), err)
	}

	if err := d.Set("network", []map[string]string{
		{
			"id":   conn.Network.Id,
			"href": conn.Network.Href,
		},
	}); err != nil {
		return fmt.Errorf("Error setting network for AWS Cloud Connection %s: %s", d.Id(), err)
	}

	return nil
}

func resourceAWSConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceAWSConnectionRead(d, m)
}

func resourceAWSConnectionDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteConnection(d, m)
}
