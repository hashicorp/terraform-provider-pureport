package pureport

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/connection"
)

func dataSourceGoogleCloudConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"connection_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile("conn-.{16}"), "Connection ID must start with 'conn-' with 16 trailing characters."),
		},
		"primary_pairing_key": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"speed": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"secondary_pairing_key": {
			Type:     schema.TypeString,
			Computed: true,
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
		Read:   dataSourceGoogleCloudConnectionRead,
		Schema: connection_schema,
	}
}

func dataSourceGoogleCloudConnectionRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("connection_id").(string))

	return resourceGoogleCloudConnectionRead(d, m)
}
