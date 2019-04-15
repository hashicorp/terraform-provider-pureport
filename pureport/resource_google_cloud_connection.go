// Package pureport provides ...
package pureport

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGoogleCloudConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"primary_pairing_key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"secondary_pairing_key": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	// Add the base items
	for k, v := range getBaseConnectionSchema() {
		connection_schema[k] = v
	}

	return &schema.Resource{
		Create: resourceGoogleCloudConnectionCreate,
		Read:   resourceGoogleCloudConnectionRead,
		Update: resourceGoogleCloudConnectionUpdate,
		Delete: resourceGoogleCloudConnectionDelete,

		Schema: connection_schema,
	}
}

func resourceGoogleCloudConnectionCreate(d *schema.ResourceData, m interface{}) error {
	return resourceGoogleCloudConnectionRead(d, m)
}

func resourceGoogleCloudConnectionRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceGoogleCloudConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceGoogleCloudConnectionRead(d, m)
}

func resourceGoogleCloudConnectionDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
