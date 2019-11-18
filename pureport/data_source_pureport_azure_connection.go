package pureport

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/connection"
)

func dataSourceAzureConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"connection_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile("conn-.{16}"), "Connection ID must start with 'conn-' with 16 trailing characters."),
		},
		"speed": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"service_key": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"peering_type": {
			Type:        schema.TypeString,
			Description: "The peering type to use for this connection: [PUBLIC, PRIVATE]",
			Computed:    true,
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
	for k, v := range connection.GetBaseDataSourceConnectionSchema() {
		connection_schema[k] = v
	}

	return &schema.Resource{
		Read:   dataSourceAzureConnectionRead,
		Schema: connection_schema,
	}
}

func dataSourceAzureConnectionRead(d *schema.ResourceData, m interface{}) error {

	d.SetId(d.Get("connection_id").(string))

	return resourceAzureConnectionRead(d, m)
}
