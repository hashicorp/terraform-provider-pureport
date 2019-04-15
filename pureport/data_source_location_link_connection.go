// Package pureport provides ...
package pureport

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceLocationLinkConnection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocationLinkConnectionRead,

		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"speed": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLocationLinkConnectionRead(d *schema.ResourceData, m interface{}) error {
	return nil
}
