// Package pureport provides ...
package pureport

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceAzureConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"service_key": {
			Type:     schema.TypeString,
			Required: true,
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
		Create: resourceAzureConnectionCreate,
		Read:   resourceAzureConnectionRead,
		Update: resourceAzureConnectionUpdate,
		Delete: resourceAzureConnectionDelete,

		Schema: connection_schema,
	}
}

func resourceAzureConnectionCreate(d *schema.ResourceData, m interface{}) error {
	return resourceAzureConnectionRead(d, m)
}

func resourceAzureConnectionRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceAzureConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceAzureConnectionRead(d, m)
}

func resourceAzureConnectionDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
