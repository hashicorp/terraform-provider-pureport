package pureport

import (
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pureport/terraform-provider-pureport/pureport/connection"
)

func dataSourceAWSConnection() *schema.Resource {

	connection_schema := map[string]*schema.Schema{
		"connection_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile("conn-.{16}"), "Connection ID must start with 'conn-' with 16 trailing characters."),
		},
		"aws_account_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"aws_region": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"speed": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"cloud_service_hrefs": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
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
		Read:   dataSourceAWSConnectionRead,
		Schema: connection_schema,
	}
}

func dataSourceAWSConnectionRead(d *schema.ResourceData, m interface{}) error {

	d.SetId(d.Get("connection_id").(string))

	return resourceAWSConnectionRead(d, m)
}
