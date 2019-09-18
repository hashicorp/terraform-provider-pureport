package pureport

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"time"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/configuration"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/connection"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/tags"
)

func resourceAzureConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"service_key": {
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
		"peering_type": {
			Type:         schema.TypeString,
			Description:  "The peering type to use for this connection: [PUBLIC, PRIVATE]",
			Default:      "PRIVATE",
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"private", "public"}, true),
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
		Create: resourceAzureConnectionCreate,
		Read:   resourceAzureConnectionRead,
		Update: resourceAzureConnectionUpdate,
		Delete: resourceAzureConnectionDelete,

		Schema: connection_schema,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},
	}
}

func expandAzureConnection(d *schema.ResourceData) client.AzureExpressRouteConnection {

	// Generic Connection values
	speed := d.Get("speed").(int)

	// Azure specific values
	serviceKey := d.Get("service_key").(string)

	// Create the body of the request
	c := client.AzureExpressRouteConnection{
		Type_: "AZURE_EXPRESS_ROUTE",
		Name:  d.Get("name").(string),
		Speed: int32(speed),
		Location: &client.Link{
			Href: d.Get("location_href").(string),
		},
		Network: &client.Link{
			Href: d.Get("network_href").(string),
		},
		BillingTerm: d.Get("billing_term").(string),
		ServiceKey:  serviceKey,
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

	if t, ok := d.GetOk("tags"); ok {
		c.Tags = tags.FilterTags(t.(map[string]interface{}))
	}

	// Azure Optionals
	c.Peering = connection.ExpandPeeringType(d)

	return c
}

func resourceAzureConnectionCreate(d *schema.ResourceData, m interface{}) error {

	c := expandAzureConnection(d)

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
			log.Printf("Error Creating new %s: %v", connection.AzureConnectionName, err)

		} else {
			statusCode := int(response["status"].(float64))

			log.Printf("Error Creating new %s: %d\n", connection.AzureConnectionName, statusCode)
			log.Printf("  %s\n", response["code"])
			log.Printf("  %s\n", response["message"])
		}

		d.SetId("")
		return fmt.Errorf("Error while creating %s: err=%s", connection.AzureConnectionName, http_err)
	}

	if resp.StatusCode >= 300 {
		d.SetId("")
		return fmt.Errorf("Error while creating %s: code=%v", connection.AzureConnectionName, resp.StatusCode)
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

	if err := connection.WaitForConnection(connection.AzureConnectionName, d, m); err != nil {
		return fmt.Errorf("Error waiting for %s: err=%s", connection.AzureConnectionName, err)
	}

	return resourceAzureConnectionRead(d, m)
}

func resourceAzureConnectionRead(d *schema.ResourceData, m interface{}) error {

	config := m.(*configuration.Config)
	connectionId := d.Id()
	ctx := config.Session.GetSessionContext()

	c, resp, err := config.Session.Client.ConnectionsApi.GetConnection(ctx, connectionId)
	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf("Error Response while reading %s: code=%v", connection.AzureConnectionName, resp.StatusCode)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading data for %s: %s", connection.AzureConnectionName, err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while reading %s: code=%v", connection.AzureConnectionName, resp.StatusCode)
	}

	conn := c.(client.AzureExpressRouteConnection)
	d.Set("description", conn.Description)
	d.Set("high_availability", conn.HighAvailability)
	d.Set("href", conn.Href)
	d.Set("name", conn.Name)
	d.Set("peering_type", conn.Peering.Type_)
	d.Set("service_key", conn.ServiceKey)
	d.Set("speed", conn.Speed)
	d.Set("state", conn.State)

	if err := d.Set("customer_networks", connection.FlattenCustomerNetworks(conn.CustomerNetworks)); err != nil {
		return fmt.Errorf("Error setting customer networks for %s %s: %s", connection.AzureConnectionName, d.Id(), err)
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
		return fmt.Errorf("Error setting gateway information for %s %s: %s", connection.AzureConnectionName, d.Id(), err)
	}

	// NAT Configuration
	if err := d.Set("nat_config", connection.FlattenNatConfig(conn.Nat)); err != nil {
		return fmt.Errorf("Error setting NAT Configuration for %s %s: %s", connection.AzureConnectionName, d.Id(), err)
	}

	if err := d.Set("location_href", conn.Location.Href); err != nil {
		return fmt.Errorf("Error setting location for %s %s: %s", connection.AzureConnectionName, d.Id(), err)
	}
	if err := d.Set("network_href", conn.Network.Href); err != nil {
		return fmt.Errorf("Error setting network for %s %s: %s", connection.AzureConnectionName, d.Id(), err)
	}

	if err := d.Set("tags", conn.Tags); err != nil {
		return fmt.Errorf("Error setting tags for %s %s: %s", connection.AzureConnectionName, d.Id(), err)
	}

	return nil
}

func resourceAzureConnectionUpdate(d *schema.ResourceData, m interface{}) error {

	c := expandAzureConnection(d)

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
				log.Printf("Error updating %s: %d\n", connection.AzureConnectionName, statusCode)
				log.Printf("  %s\n", response["code"])
				log.Printf("  %s\n", response["message"])
			}
		}

		return fmt.Errorf("Error while updating %s: err=%s", connection.AzureConnectionName, err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while updating %s: code=%v", connection.AzureConnectionName, resp.StatusCode)
	}

	if err := connection.WaitForConnection(connection.AzureConnectionName, d, m); err != nil {
		return fmt.Errorf("Error waiting for %s: err=%s", connection.AzureConnectionName, err)
	}

	d.Partial(false)

	return resourceAzureConnectionRead(d, m)
}

func resourceAzureConnectionDelete(d *schema.ResourceData, m interface{}) error {
	return connection.DeleteConnection(connection.AzureConnectionName, d, m)
}
