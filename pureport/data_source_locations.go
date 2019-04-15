// Package pureport provides ...
package pureport

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceLocations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocationsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceLocationsRead(d *schema.ResourceData, m interface{}) error {
	return nil
}
