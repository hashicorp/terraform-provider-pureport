package pureport

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"sort"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pureport/pureport-sdk-go/pureport/client"
)

const (
	awsConnectionName = "AWS Cloud Connection"
)

func resourceAWSConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"cloud_service_hrefs": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"peering_type": {
			Type:         schema.TypeString,
			Description:  "The peering type to use for this connection: [PUBLIC, PRIVATE]",
			Default:      "PRIVATE",
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"private", "public"}, true),
		},
		"aws_account_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"aws_region": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"gateways": {
			Computed: true,
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 2,
			Elem: &schema.Resource{
				Schema: StandardGatewaySchema,
			},
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

func expandAWSConnection(d *schema.ResourceData) client.AwsDirectConnectConnection {

	// Generic Connection values
	speed := d.Get("speed").(int)

	// Create the body of the request
	c := client.AwsDirectConnectConnection{
		Type_: "AWS_DIRECT_CONNECT",
		Name:  d.Get("name").(string),
		Speed: int32(speed),
		Location: &client.Link{
			Href: d.Get("location_href").(string),
		},
		Network: &client.Link{
			Href: d.Get("network_href").(string),
		},
		AwsAccountId: d.Get("aws_account_id").(string),
		AwsRegion:    d.Get("aws_region").(string),
		BillingTerm:  d.Get("billing_term").(string),
	}

	c.CustomerNetworks = ExpandCustomerNetworks(d)
	c.Nat = ExpandNATConfiguration(d)
	c.CloudServices = ExpandCloudServices(d)
	c.Peering = ExpandPeeringType(d)

	if description, ok := d.GetOk("description"); ok {
		c.Description = description.(string)
	}

	if highAvailability, ok := d.GetOk("high_availability"); ok {
		c.HighAvailability = highAvailability.(bool)
	}

	return c
}

func resourceAWSConnectionCreate(d *schema.ResourceData, m interface{}) error {

	connection := expandAWSConnection(d)

	config := m.(*Config)
	ctx := config.Session.GetSessionContext()

	opts := client.AddConnectionOpts{
		Body: optional.NewInterface(connection),
	}

	resp, err := config.Session.Client.ConnectionsApi.AddConnection(
		ctx,
		filepath.Base(connection.Network.Href),
		&opts,
	)

	if err != nil {

		http_err := err
		json_response := string(err.(client.GenericSwaggerError).Body()[:])
		response, err := structure.ExpandJsonFromString(json_response)

		if err != nil {
			log.Printf("Error Creating new %s: %v", awsConnectionName, err)

		} else {
			statusCode := int(response["status"].(float64))

			log.Printf("Error Creating new %s: %d\n", awsConnectionName, statusCode)
			log.Printf("  %s\n", response["code"])
			log.Printf("  %s\n", response["message"])
		}

		d.SetId("")
		return fmt.Errorf("Error while creating %s: err=%s", awsConnectionName, http_err)
	}

	if resp.StatusCode >= 300 {
		d.SetId("")
		return fmt.Errorf("Error while creating %s: code=%v", awsConnectionName, resp.StatusCode)
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

	WaitForConnection(awsConnectionName, d, m)

	return resourceAWSConnectionRead(d, m)
}

func resourceAWSConnectionRead(d *schema.ResourceData, m interface{}) error {

	config := m.(*Config)
	connectionId := d.Id()
	ctx := config.Session.GetSessionContext()

	c, resp, err := config.Session.Client.ConnectionsApi.GetConnection(ctx, connectionId)
	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf("Error Response while reading %s: code=%v", awsConnectionName, resp.StatusCode)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading data for %s: %s", awsConnectionName, err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while reading %s: code=%v", awsConnectionName, resp.StatusCode)
	}

	conn := c.(client.AwsDirectConnectConnection)
	d.Set("aws_account_id", conn.AwsAccountId)
	d.Set("aws_region", conn.AwsRegion)

	var cloudServiceHrefs []string
	for _, cs := range conn.CloudServices {
		cloudServiceHrefs = append(cloudServiceHrefs, cs.Href)
	}

	sort.Slice(cloudServiceHrefs, func(i int, j int) bool {
		return cloudServiceHrefs[i] < cloudServiceHrefs[j]
	})

	if err := d.Set("cloud_service_hrefs", cloudServiceHrefs); err != nil {
		return fmt.Errorf("Error setting cloud services for %s %s: %s", awsConnectionName, d.Id(), err)
	}

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
		return fmt.Errorf("Error setting customer networks for %s %s: %s", awsConnectionName, d.Id(), err)
	}

	// Add Gateway information
	var gateways []map[string]interface{}
	if g := conn.PrimaryGateway; g != nil {
		gateways = append(gateways, FlattenStandardGateway(g))
	}
	if g := conn.SecondaryGateway; g != nil {
		gateways = append(gateways, FlattenStandardGateway(g))
	}
	if err := d.Set("gateways", gateways); err != nil {
		return fmt.Errorf("Error setting gateway information for %s %s: %s", awsConnectionName, d.Id(), err)
	}

	d.Set("description", conn.Description)
	d.Set("high_availability", conn.HighAvailability)

	if err := d.Set("location_href", conn.Location.Href); err != nil {
		return fmt.Errorf("Error setting location for %s %s: %s", awsConnectionName, d.Id(), err)
	}

	if err := d.Set("network_href", conn.Network.Href); err != nil {
		return fmt.Errorf("Error setting network for %s %s: %s", awsConnectionName, d.Id(), err)
	}

	return nil
}

func resourceAWSConnectionUpdate(d *schema.ResourceData, m interface{}) error {

	c := expandAWSConnection(d)

	d.Partial(true)

	config := m.(*Config)
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
		c.CustomerNetworks = ExpandCustomerNetworks(d)
	}

	if d.HasChange("nat_config") {
		c.Nat = ExpandNATConfiguration(d)
	}

	if d.HasChange("billing_term") {
		c.BillingTerm = d.Get("billing_term").(string)
	}

	if d.HasChange("cloud_services") {
		c.CloudServices = ExpandCloudServices(d)
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

		json_response := string(err.(client.GenericSwaggerError).Body()[:])
		response, err := structure.ExpandJsonFromString(json_response)
		if err != nil {
			log.Printf("Error updating %s: %v", awsConnectionName, err)
		} else {
			statusCode := int(response["status"].(float64))
			log.Printf("Error updating %s: %d\n", awsConnectionName, statusCode)
			log.Printf("  %s\n", response["code"])
			log.Printf("  %s\n", response["message"])
		}

		return fmt.Errorf("Error while updating %s: err=%s", awsConnectionName, err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while updating %s: code=%v", awsConnectionName, resp.StatusCode)
	}

	d.Partial(false)

	return resourceAWSConnectionRead(d, m)
}

func resourceAWSConnectionDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteConnection(awsConnectionName, d, m)
}
