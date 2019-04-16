// Package pureport provides ...
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

func dataSourceLocations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocationsRead,

		Schema: map[string]*schema.Schema{
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
	ctx := sess.GetSessionContext()

	locations, resp, err := sess.Client.LocationsApi.FindLocations(ctx)
	if err != nil {
		log.Printf("[Error] Error when Reading Pureport Location data")
		d.SetId("")
		return nil
	}

	if resp.StatusCode != 200 {
		log.Printf("[Error] Error Response while Reading Pureport Location data")
		d.SetId("")
		return nil
	}

	// Sort the list
	sort.Slice(locations, func(i int, j int) bool {
		return locations[i].Id < locations[j].Id
	})

	// Convert to Map
	out := flattenLocations(locations)
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

func flattenLocations(locations []swagger.Location) (out []map[string]interface{}) {

	for _, loc := range locations {

		l := map[string]interface{}{
			"id":    loc.Id,
			"name":  loc.Name,
			"links": flattenLinks(loc.LocationLinks),
		}

		out = append(out, l)
	}

	return
}

func flattenLinks(links []swagger.LocationLinkConnection) (out []map[string]interface{}) {

	for _, link := range links {

		l := map[string]interface{}{
			"location_id": link.Location.Id,
			"speed":       link.Speed,
		}

		out = append(out, l)
	}

	return
}
