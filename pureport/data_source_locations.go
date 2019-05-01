// Package pureport provides ...
package pureport

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/pureport/pureport-sdk-go/pureport/session"
)

func dataSourceLocations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocationsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
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
									"location_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"speed": {
										Type:     schema.TypeString,
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

	sess := m.(*session.Session)
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	ctx := sess.GetSessionContext()

	locations, resp, err := sess.Client.LocationsApi.FindLocations(ctx)
	if err != nil {
		log.Printf("[Error] Error when Reading Pureport Location data")
		d.SetId("")
		return nil
	}

	if resp.StatusCode >= 300 {
		log.Printf("[Error] Error Response while Reading Pureport Location data")
		d.SetId("")
		return nil
	}

	// Filter the results
	var filteredLocations []client.Location
	if nameRegexOk {
		r := regexp.MustCompile(nameRegex.(string))
		for _, location := range locations {
			if r.MatchString(location.Name) {
				filteredLocations = append(filteredLocations, location)
			}
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
			"location_id": link.Location.Id,
			"speed":       link.Speed,
		}

		out = append(out, l)
	}

	return
}
