package pureport

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCloudServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudServicesRead,

		Schema: map[string]*schema.Schema{
			"services": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv4_prefix_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ipv6_prefix_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cloud_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudServicesRead(d *schema.ResourceData, m interface{}) error {
	return nil
}
