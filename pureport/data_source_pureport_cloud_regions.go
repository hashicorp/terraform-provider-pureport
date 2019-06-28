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
	"github.com/pureport/terraform-provider-pureport/pureport/configuration"
)

func dataSourceCloudRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudRegionsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"regions": {
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
						"identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudRegionsRead(d *schema.ResourceData, m interface{}) error {

	config := m.(*configuration.Config)
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	ctx := config.Session.GetSessionContext()

	regions, resp, err := config.Session.Client.CloudRegionsApi.GetCloudRegions(ctx)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error when Reading Cloud Region data: %v", err)
	}

	if resp.StatusCode >= 300 {
		d.SetId("")

		if resp.StatusCode == 404 {
			// Need to gracefully handle 404, for refresh
			return nil
		}
		return fmt.Errorf("Error Response while Reading Cloud Region data")
	}

	// Filter the results
	var filteredRegions []client.CloudRegion
	if nameRegexOk {
		r := regexp.MustCompile(nameRegex.(string))
		for _, region := range regions {
			if r.MatchString(region.DisplayName) {
				filteredRegions = append(filteredRegions, region)
			}
		}
	} else {
		filteredRegions = regions
	}

	// Sort the list
	sort.Slice(filteredRegions, func(i int, j int) bool {
		return filteredRegions[i].Id < filteredRegions[j].Id
	})

	// Convert to Map
	out := flattenRegions(filteredRegions)
	if err := d.Set("regions", out); err != nil {
		return fmt.Errorf("Error reading Cloud Regions: %s", err)
	}

	data, err := json.Marshal(regions)
	if err != nil {
		return fmt.Errorf("Error generating Id: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", hashcode.String(string(data))))

	return nil
}

func flattenRegions(regions []client.CloudRegion) (out []map[string]interface{}) {

	for _, cr := range regions {

		r := map[string]interface{}{
			"id":         cr.Id,
			"name":       cr.DisplayName,
			"provider":   cr.Provider,
			"identifier": cr.ProviderAssignedId,
		}

		out = append(out, r)
	}

	return
}
