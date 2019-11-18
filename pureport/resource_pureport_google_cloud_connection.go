package pureport

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"
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

func resourceGoogleCloudConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"primary_pairing_key": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"speed": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntInSlice([]int{50, 100, 200, 300, 400, 500, 1000, 10000}),
		},
		"secondary_pairing_key": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"gateways": {
			Computed: true,
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 2,
			Elem: &schema.Resource{
				Schema: connection.StandardGatewaySchema,
			},
		},
	}

	// Add the base items
	for k, v := range connection.GetBaseResourceConnectionSchema() {
		connection_schema[k] = v
	}

	return &schema.Resource{
		Create: resourceGoogleCloudConnectionCreate,
		Read:   resourceGoogleCloudConnectionRead,
		Update: resourceGoogleCloudConnectionUpdate,
		Delete: resourceGoogleCloudConnectionDelete,

		Schema: connection_schema,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},
	}
}

func expandGoogleCloudConnection(d *schema.ResourceData) client.GoogleCloudInterconnectConnection {

	// Generic Connection values
	speed := d.Get("speed").(int)

	// Google specific values
	primaryPairingKey := d.Get("primary_pairing_key").(string)

	// Create the body of the request
	c := client.GoogleCloudInterconnectConnection{
		Type_: "GOOGLE_CLOUD_INTERCONNECT",
		Name:  d.Get("name").(string),
		Speed: int32(speed),
		Location: &client.Link{
			Href: d.Get("location_href").(string),
		},
		Network: &client.Link{
			Href: d.Get("network_href").(string),
		},
		BillingTerm:       d.Get("billing_term").(string),
		PrimaryPairingKey: primaryPairingKey,
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

	// Google Optionals
	if secondaryPairingKey, ok := d.GetOk("secondary_pairing_key"); ok {
		c.SecondaryPairingKey = secondaryPairingKey.(string)
	}

	if t, ok := d.GetOk("tags"); ok {
		c.Tags = tags.FilterTags(t.(map[string]interface{}))
	}

	return c
}

func resourceGoogleCloudConnectionCreate(d *schema.ResourceData, m interface{}) error {

	c := expandGoogleCloudConnection(d)

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

		switch e := err.(type) {
		case client.GenericSwaggerError:
			json_response := string(e.Body()[:])
			response, err := structure.ExpandJsonFromString(json_response)

			if err != nil {
				log.Printf("Error creating new %s: %v", connection.GoogleConnectionName, err)
			} else {
				statusCode := int(response["status"].(float64))
				log.Printf("Error creating new %s: %d\n", connection.GoogleConnectionName, statusCode)
				log.Printf("  %s\n", response["code"])
				log.Printf("  %s\n", response["message"])
			}
		case *url.Error:
			log.Printf("Error creating new %s: %s", connection.GoogleConnectionName, e.Error())
		default:
			log.Printf("Error creating new %s: %v", connection.GoogleConnectionName, e)
		}

		d.SetId("")
		return fmt.Errorf("Error while creating %s: err=%s", connection.GoogleConnectionName, http_err)
	}

	if resp.StatusCode >= 300 {
		d.SetId("")
		return fmt.Errorf("Error while creating %s: code=%v", connection.GoogleConnectionName, resp.StatusCode)
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

	if err := connection.WaitForConnection(connection.GoogleConnectionName, d, m); err != nil {
		return fmt.Errorf("Error waiting for %s: err=%s", connection.GoogleConnectionName, err)
	}

	return resourceGoogleCloudConnectionRead(d, m)
}

func resourceGoogleCloudConnectionRead(d *schema.ResourceData, m interface{}) error {

	config := m.(*configuration.Config)
	connectionId := d.Id()
	ctx := config.Session.GetSessionContext()

	c, resp, err := config.Session.Client.ConnectionsApi.GetConnection(ctx, connectionId)
	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf("Error Response while reading %s: code=%v", connection.GoogleConnectionName, resp.StatusCode)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading data for %s: %s", connection.GoogleConnectionName, err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while reading %s: code=%v", connection.GoogleConnectionName, resp.StatusCode)
	}

	conn := c.(client.GoogleCloudInterconnectConnection)
	d.Set("description", conn.Description)
	d.Set("high_availability", conn.HighAvailability)
	d.Set("href", conn.Href)
	d.Set("name", conn.Name)
	d.Set("primary_pairing_key", conn.PrimaryPairingKey)
	d.Set("secondary_pairing_key", conn.SecondaryPairingKey)
	d.Set("speed", conn.Speed)
	d.Set("state", conn.State)

	if err := d.Set("customer_networks", connection.FlattenCustomerNetworks(conn.CustomerNetworks)); err != nil {
		return fmt.Errorf("Error setting customer networks for %s %s: %s", connection.GoogleConnectionName, d.Id(), err)
	}

	// Add Gateway information
	var gateways []map[string]interface{}
	if g := conn.PrimaryGateway; g != nil {
		gateways = append(gateways, connection.FlattenStandardGateway(g))
	}
	if g := conn.SecondaryGateway; g != nil {
		gateways = append(gateways, connection.FlattenStandardGateway(g))
	}
	if err := d.Set("gateways", gateways); err != nil {
		return fmt.Errorf("Error setting gateway information for %s %s: %s", connection.GoogleConnectionName, d.Id(), err)
	}

	// NAT Configuration
	if err := d.Set("nat_config", connection.FlattenNatConfig(conn.Nat)); err != nil {
		return fmt.Errorf("Error setting NAT Configuration for %s %s: %s", connection.GoogleConnectionName, d.Id(), err)
	}

	if err := d.Set("location_href", conn.Location.Href); err != nil {
		return fmt.Errorf("Error setting location for %s %s: %s", connection.GoogleConnectionName, d.Id(), err)
	}
	if err := d.Set("network_href", conn.Network.Href); err != nil {
		return fmt.Errorf("Error setting network for %s %s: %s", connection.GoogleConnectionName, d.Id(), err)
	}

	if err := d.Set("tags", conn.Tags); err != nil {
		return fmt.Errorf("Error setting tags for %s %s: %s", connection.GoogleConnectionName, d.Id(), err)
	}

	return nil
}

func resourceGoogleCloudConnectionUpdate(d *schema.ResourceData, m interface{}) error {

	c := expandGoogleCloudConnection(d)

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
				log.Printf("Error updating %s: %d\n", connection.GoogleConnectionName, statusCode)
				log.Printf("  %s\n", response["code"])
				log.Printf("  %s\n", response["message"])
			}
		}

		return fmt.Errorf("Error while updating %s: err=%s", connection.GoogleConnectionName, err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while updating %s: code=%v", connection.GoogleConnectionName, resp.StatusCode)
	}

	if err := connection.WaitForConnection(connection.GoogleConnectionName, d, m); err != nil {
		return fmt.Errorf("Error waiting for %s: err=%s", connection.GoogleConnectionName, err)
	}

	d.Partial(false)

	return resourceGoogleCloudConnectionRead(d, m)
}

func resourceGoogleCloudConnectionDelete(d *schema.ResourceData, m interface{}) error {
	return connection.DeleteConnection(connection.GoogleConnectionName, d, m)
}
