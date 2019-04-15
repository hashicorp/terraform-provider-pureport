// Package main provides Connection resource
package pureport

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectionCreate,
		Read:   resourceConnectionRead,
		Update: resourceConnectionUpdate,
		Delete: resourceConnectionDelete,

		Schema: map[string]*schema.Schema{},
	}
}

func resourceConnectionCreate(d *schema.ResourceData, m interface{}) error {
	return resourceConnectionRead(d, m)
}

func resourceConnectionRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceConnectionRead(d, m)
}

func resourceConnectionDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
