package pureport

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/configuration"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/tags"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkCreate,
		Read:   resourceNetworkRead,
		Update: resourceNetworkUpdate,
		Delete: resourceNetworkDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_href": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tags.TagsSchema(),
			"href": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func expandNetwork(d *schema.ResourceData) client.Network {

	n := client.Network{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if t, ok := d.GetOk("tags"); ok {
		n.Tags = tags.FilterTags(t.(map[string]interface{}))
	}

	return n
}

func resourceNetworkCreate(d *schema.ResourceData, m interface{}) error {

	network := expandNetwork(d)
	accountHref := d.Get("account_href").(string)
	accountId := filepath.Base(accountHref)

	config := m.(*configuration.Config)
	ctx := config.Session.GetSessionContext()

	opts := client.AddNetworkOpts{
		Body: optional.NewInterface(network),
	}

	_, resp, err := config.Session.Client.NetworksApi.AddNetwork(
		ctx,
		accountId,
		&opts,
	)

	if err != nil {

		http_err := err
		json_response := string(err.(client.GenericSwaggerError).Body()[:])
		response, err := structure.ExpandJsonFromString(json_response)
		if err != nil {
			log.Printf("Error Creating new Network: %v", err)
		} else {
			statusCode := int(response["status"].(float64))
			log.Printf("Error Creating new Network: %d\n", statusCode)
			log.Printf("  %s\n", response["code"])
			log.Printf("  %s\n", response["message"])
		}

		d.SetId("")
		return fmt.Errorf("Error while creating Network: err=%s", http_err)
	}

	if resp.StatusCode >= 300 {
		d.SetId("")
		return fmt.Errorf("Error while creating network: code=%v", resp.StatusCode)
	}

	loc := resp.Header.Get("location")
	u, err := url.Parse(loc)
	if err != nil {
		return fmt.Errorf("Error when decoding Network ID")
	}

	id := filepath.Base(u.Path)
	d.SetId(id)

	if id == "" {
		log.Printf("Error when decoding location header")
		return fmt.Errorf("Error when decoding Network ID")
	}

	return resourceNetworkRead(d, m)
}

func resourceNetworkRead(d *schema.ResourceData, m interface{}) error {

	config := m.(*configuration.Config)
	networkId := d.Id()
	ctx := config.Session.GetSessionContext()

	n, resp, err := config.Session.Client.NetworksApi.GetNetwork(ctx, networkId)
	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf("Error Response while reading Network: code=%v", resp.StatusCode)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading data for Network: %s", err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while reading Network: code=%v", resp.StatusCode)
	}

	d.Set("name", n.Name)
	d.Set("description", n.Description)
	d.Set("href", n.Href)
	d.Set("account_href", n.Account.Href)

	if err := d.Set("tags", n.Tags); err != nil {
		return fmt.Errorf("Error setting tags for Network %s: %s", d.Id(), err)
	}

	return nil
}

func resourceNetworkUpdate(d *schema.ResourceData, m interface{}) error {

	n := expandNetwork(d)

	d.Partial(true)

	if d.HasChange("name") {
		n.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		n.Description = d.Get("description").(string)
	}

	if d.HasChange("tags") {
		_, nraw := d.GetChange("tags")
		n.Tags = tags.FilterTags(nraw.(map[string]interface{}))
	}

	config := m.(*configuration.Config)
	ctx := config.Session.GetSessionContext()

	opts := client.UpdateNetworkOpts{
		Body: optional.NewInterface(n),
	}

	_, resp, err := config.Session.Client.NetworksApi.UpdateNetwork(
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
				log.Printf("Error updating Network: %d\n", statusCode)
				log.Printf("  %s\n", response["code"])
				log.Printf("  %s\n", response["message"])
			}
		}

		return fmt.Errorf("Error while updating Network: err=%s", err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while updating Network : code=%v", resp.StatusCode)
	}

	d.Partial(false)
	return resourceNetworkRead(d, m)
}

func resourceNetworkDelete(d *schema.ResourceData, m interface{}) error {

	config := m.(*configuration.Config)
	ctx := config.Session.GetSessionContext()
	networkId := d.Id()

	// Delete
	resp, err := config.Session.Client.NetworksApi.DeleteNetwork(ctx, networkId)

	if err != nil {
		return fmt.Errorf("Error deleting Network: %s", err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while Network: code=%v", resp.StatusCode)
	}

	d.SetId("")

	return nil
}
