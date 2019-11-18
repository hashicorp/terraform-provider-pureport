package pureport

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/configuration"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/filter"
)

func dataSourceLocations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocationsRead,

		Schema: map[string]*schema.Schema{
			"filter": filter.DataSourceFiltersSchema(),
			"locations": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"href": {
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
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"location_href": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"speed": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceLocationsRead(d *schema.ResourceData, m interface{}) error {

	config := m.(*configuration.Config)
	filters, filtersOk := d.GetOk("filter")

	ctx := config.Session.GetSessionContext()

	locations, resp, err := config.Session.Client.LocationsApi.FindLocations(ctx)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error when Reading Pureport Location data: %v", err)
	}

	if resp.StatusCode >= 300 {
		d.SetId("")

		if resp.StatusCode == 404 {
			// Need to gracefully handle 404, for refresh
			return nil
		}
		return fmt.Errorf("Error Response while Reading Pureport Location data")
	}

	// Filter the results
	var filteredLocations []client.Location
	if filtersOk {

		input := make([]interface{}, len(locations))
		for i, x := range locations {
			input[i] = x
		}

		output := filter.FilterType(input, filter.BuildDataSourceFilters(filters.(*schema.Set)))
		for _, x := range output {
			filteredLocations = append(filteredLocations, x.(client.Location))
		}

	} else {
		filteredLocations = locations
	}

	// Sort the list
	sort.Slice(filteredLocations, func(i int, j int) bool {
		return filteredLocations[i].Id < filteredLocations[j].Id
	})

	// Convert to Map
	out := flattenLocations(filteredLocations)
	if err := d.Set("locations", out); err != nil {
		return fmt.Errorf("Error reading locations: %s", err)
	}

	data, err := json.Marshal(locations)
	if err != nil {
		return fmt.Errorf("Error generating Id: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", hashcode.String(string(data))))

	return nil
}

func flattenLocations(locations []client.Location) (out []map[string]interface{}) {

	for _, loc := range locations {

		l := map[string]interface{}{
			"id":    loc.Id,
			"href":  loc.Href,
			"name":  loc.Name,
			"links": flattenLinks(loc.LocationLinks),
		}

		out = append(out, l)
	}

	return
}

func flattenLinks(links []client.LocationLinkConnection) (out []map[string]interface{}) {

	for _, link := range links {

		l := map[string]interface{}{
			"location_href": link.Location.Href,
			"speed":         link.Speed,
		}

		out = append(out, l)
	}

	return
}
