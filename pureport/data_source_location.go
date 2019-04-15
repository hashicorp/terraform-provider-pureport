// Package pureport provides ...
package pureport

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceLocation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocationRead,

		Schema: map[string]*schema.Schema{
			"href": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
			},
		},
	}
}

func dataSourceLocationRead(d *schema.ResourceData, m interface{}) error {
	return nil
}
