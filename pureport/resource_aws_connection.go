// Package main provides AWSConnection resource
package pureport

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
	return resourceAWSConnectionRead(d, m)
}

func resourceAWSConnectionRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceAWSConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceAWSConnectionRead(d, m)
}

func resourceAWSConnectionDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
