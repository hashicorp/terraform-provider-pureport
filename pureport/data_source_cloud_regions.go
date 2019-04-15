// Package pureport provides ...
package pureport

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCloudRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudRegionsRead,

		Schema: map[string]*schema.Schema{
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type: schema.TypeString,
						},
						"name": {
							Type: schema.TypeString,
						},
						"provider": {
							Type: schema.TypeString,
						},
						"identifier": {
							Type: schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudRegionsRead(d *schema.ResourceData, m interface{}) error {
	return nil
}
