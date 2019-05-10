package pureport

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/pureport/pureport-sdk-go/pureport/session"
)

func dataSourceCloudServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudServicesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
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
						"href": {
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
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	ctx := sess.GetSessionContext()

	services, resp, err := sess.Client.CloudServicesApi.GetCloudServices(ctx)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error when Reading Cloud Services data: %v", err)
	}

	if resp.StatusCode >= 300 {
		d.SetId("")
		return fmt.Errorf("Error Response while Reading Cloud Services data")
	}

	// Filter the results
	var filteredServices []client.CloudService
	if nameRegexOk {
		r := regexp.MustCompile(nameRegex.(string))
		for _, service := range services {
			if r.MatchString(service.Name) {
				filteredServices = append(filteredServices, service)
			}
		}
	} else {
		filteredServices = services
	}

	// Sort the list
	sort.Slice(filteredServices, func(i int, j int) bool {
		return filteredServices[i].Id < filteredServices[j].Id
	})

	// Convert to Map
	out := flattenServices(filteredServices)
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

func flattenServices(services []client.CloudService) (out []map[string]interface{}) {

	for _, cs := range services {

		s := map[string]interface{}{
			"id":                cs.Id,
			"href":              cs.Href,
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
