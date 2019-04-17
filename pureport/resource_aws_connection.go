// Package main provides AWSConnection resource
package pureport

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pureport/pureport-sdk-go/pureport/session"
	"github.com/pureport/pureport-sdk-go/pureport/swagger"
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

	// AWS specific values
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

	// AWS Optionals
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

	ctx := sess.GetSessionContext()

	opts := swagger.AddConnectionOpts{
		Body: optional.NewInterface(connection),
	}

	resp, err := sess.Client.ConnectionsApi.AddConnection(
		ctx,
		network[0].(map[string]interface{})["id"].(string),
		&opts,
	)

	if err != nil {
		log.Printf("[Error] Error Creating new AWS Connection: %s", err)
		d.SetId("")
		return nil
	}

	if resp.StatusCode != 200 {
		log.Printf("[Error] Error Response while creating new AWS Connection: code=%v", resp.StatusCode)
		d.SetId("")
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Error] Error reading create AWS Connection response")
		return nil
	}

	var data map[string]string
	if err = json.Unmarshal(body, &data); err != nil {
		log.Printf("[Error] Error reading response from creating new AWS Connection")
		return nil
	}

	d.SetId(data["id"])

	return resourceAWSConnectionRead(d, m)
}

func resourceAWSConnectionRead(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)
	connectionId := d.Id()
	ctx := sess.GetSessionContext()

	c, resp, err := sess.Client.ConnectionsApi.Get11(ctx, connectionId)
	conn := c.(swagger.AwsDirectConnectConnection)

	if err != nil {
		log.Printf("[Error] Error reading data for AWS Connection: %s", err)

		if resp.StatusCode == 404 {
			log.Printf("[Error] Error Response while reading AWS Connection: code=%v", resp.StatusCode)
			d.SetId("")
		}
		return nil
	}

	if resp.StatusCode >= 300 {
		log.Printf("[Error] Error Response while reading AWS Connection: code=%v", resp.StatusCode)
	}

	fmt.Printf("Connection: %+v", conn)
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

	fmt.Printf("Peering: %+v", conn.Peering.Type_)
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

	return nil
}

func resourceAWSConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceAWSConnectionRead(d, m)
}

func resourceAWSConnectionDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
