// Package main provides Connection resource
package pureport

import (
	"fmt"
	"log"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pureport/pureport-sdk-go/pureport/session"
	"github.com/pureport/pureport-sdk-go/pureport/swagger"
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
			"account": {
				Type:     schema.TypeList,
				Required: true,
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
	account := d.Get("account").(map[string]string)

	network := swagger.Network{
		Account: &swagger.Link{
			Id:   account["id"],
			Href: account["href"],
		},
		Name:        name,
		Description: description,
	}

	ctx := sess.GetSessionContext()

	opts := swagger.AddNetworkOpts{
		Body: optional.NewInterface(network),
	}

	resp, err := sess.Client.NetworksApi.AddNetwork(
		ctx,
		account["id"],
		&opts,
	)

	if err != nil {
		log.Printf("[Error] Error Creating new Network: %v", err)
		d.SetId("")
		return nil
	}

	if resp.StatusCode >= 300 {
		log.Printf("[Error] Error Response while creating new Network: code=%v", resp.StatusCode)
		d.SetId("")
		return nil
	}

	id := resp.Header.Get("location")
	d.SetId(id)

	if id == "" {
		log.Printf("[Error] Error when decoding location header")
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
			log.Printf("[Error] Error Response while reading Network: code=%v", resp.StatusCode)
			d.SetId("")
		}
		return fmt.Errorf("[Error] Error reading data for Network: %s", err)
	}

	if resp.StatusCode >= 300 {
		fmt.Errorf("[Error] Error Response while reading AWS Connection: code=%v", resp.StatusCode)
	}

	account := map[string]string{
		"id":   n.Account.Id,
		"href": n.Account.Href,
	}
	d.Set("account", account)
	d.Set("name", n.Name)
	d.Set("description", n.Description)

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
		return fmt.Errorf("[Error] Error deleting Network: %s", err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error Response while Network: code=%v", resp.StatusCode)
	}

	d.SetId("")

	return nil
}
