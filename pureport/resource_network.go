package pureport

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/pureport/pureport-sdk-go/pureport/session"
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
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
		},
	}
}

func resourceNetworkCreate(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	accountId := d.Get("account_id").(string)

	network := client.Network{
		Name:        name,
		Description: description,
	}

	ctx := sess.GetSessionContext()

	opts := client.AddNetworkOpts{
		Body: optional.NewInterface(network),
	}

	resp, err := sess.Client.NetworksApi.AddNetwork(
		ctx,
		accountId,
		&opts,
	)

	if err != nil {

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
		return nil
	}

	if resp.StatusCode >= 300 {
		log.Printf("Error Response while creating new Network: code=%v", resp.StatusCode)
		d.SetId("")
		return nil
	}

	loc := resp.Header.Get("location")
	u, err := url.Parse(loc)
	if err != nil {
		log.Printf("Error when decoding Network ID")
		return nil
	}

	id := filepath.Base(u.Path)
	d.SetId(id)

	if id == "" {
		log.Printf("Error when decoding location header")
		return nil
	}

	return resourceNetworkRead(d, m)
}

func resourceNetworkRead(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)
	networkId := d.Id()
	ctx := sess.GetSessionContext()

	n, resp, err := sess.Client.NetworksApi.GetNetwork(ctx, networkId)
	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf("Error Response while reading Network: code=%v", resp.StatusCode)
			d.SetId("")
		}
		return fmt.Errorf("Error reading data for Network: %s", err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while reading Network: code=%v", resp.StatusCode)
	}

	d.Set("account_id", n.Account.Id)
	d.Set("name", n.Name)
	d.Set("description", n.Description)

	if err := d.Set("account", []map[string]string{
		{
			"id":   n.Account.Id,
			"href": n.Account.Href,
		},
	}); err != nil {
		return fmt.Errorf("Error while setting Network: code=%v", resp.StatusCode)
	}

	return nil
}

func resourceNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceNetworkRead(d, m)
}

func resourceNetworkDelete(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)
	ctx := sess.GetSessionContext()
	networkId := d.Id()

	// Delete
	resp, err := sess.Client.NetworksApi.DeleteNetwork(ctx, networkId)

	if err != nil {
		return fmt.Errorf("Error deleting Network: %s", err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while Network: code=%v", resp.StatusCode)
	}

	d.SetId("")

	return nil
}
