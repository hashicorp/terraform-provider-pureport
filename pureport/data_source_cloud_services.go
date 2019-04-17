package pureport

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pureport/pureport-sdk-go/pureport/session"
	"github.com/pureport/pureport-sdk-go/pureport/swagger"
)

func dataSourceCloudServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudServicesRead,

		Schema: map[string]*schema.Schema{
			"services": {
				Type:     schema.TypeList,
				Computed: true,
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

	sess := m.(*session.Session)
	ctx := sess.GetSessionContext()

	services, resp, err := sess.Client.CloudServicesApi.GetCloudServices(ctx)
	if err != nil {
		log.Printf("[Error] Error when Reading Cloud Services data")
		d.SetId("")
		return nil
	}

	if resp.StatusCode != 200 {
		log.Printf("[Error] Error Response while Reading Cloud Services data")
		d.SetId("")
		return nil
	}

	// Sort the list
	sort.Slice(services, func(i int, j int) bool {
		return services[i].Id < services[j].Id
	})

	// Convert to Map
	out := flattenServices(services)
	if err := d.Set("services", out); err != nil {
		return fmt.Errorf("Error reading cloud services: %s", err)
	}

	data, err := json.Marshal(services)
	if err != nil {
		return fmt.Errorf("Error generating Id: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", hashcode.String(string(data))))

	return nil
}

func flattenServices(services []swagger.CloudService) (out []map[string]interface{}) {

	for _, cs := range services {

		s := map[string]interface{}{
			"id":                cs.Id,
			"name":              cs.Name,
			"provider":          cs.Provider,
			"service":           cs.Service,
			"ipv4_prefix_count": cs.Ipv4PrefixCount,
			"ipv6_prefix_count": cs.Ipv6PrefixCount,
			"cloud_region_id":   cs.CloudRegion.Id,
		}

		out = append(out, s)
	}

	return
}
