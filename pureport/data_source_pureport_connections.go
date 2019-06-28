package pureport

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/pureport/terraform-provider-pureport/pureport/configuration"
)

func dataSourceConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConnectionsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"network_href": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"connections": {
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
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"speed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"location_href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceConnectionsRead(d *schema.ResourceData, m interface{}) error {

	config := m.(*configuration.Config)
	networkHref := d.Get("network_href").(string)
	networkId := filepath.Base(networkHref)

	ctx := config.Session.GetSessionContext()

	connections, resp, err := config.Session.Client.ConnectionsApi.GetConnections(ctx, networkId)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error when Reading Cloud Services data: %v", err)
	}

	if resp.StatusCode >= 300 {
		d.SetId("")
		return fmt.Errorf("Error Response while Reading Connections data")
	}

	// Filter the results
	var filteredConnections []client.Connection

	nameRegex, nameRegexOk := d.GetOk("name_regex")
	if nameRegexOk {
		r := regexp.MustCompile(nameRegex.(string))
		for _, connection := range connections {
			if r.MatchString(connection.Name) {
				filteredConnections = append(filteredConnections, connection)
			}
		}
	} else {
		filteredConnections = connections
	}

	// Sort the list
	sort.Slice(filteredConnections, func(i int, j int) bool {
		return filteredConnections[i].Name < filteredConnections[j].Name
	})

	// Convert to Map
	if err := d.Set("connections", flattenConnections(filteredConnections)); err != nil {
		return fmt.Errorf("Error reading cloud connections: %s", err)
	}

	data, err := json.Marshal(connections)
	if err != nil {
		return fmt.Errorf("Error generating Id: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", hashcode.String(string(data))))

	return nil
}

func flattenConnections(connections []client.Connection) (out []map[string]interface{}) {

	for _, c := range connections {

		out = append(out, map[string]interface{}{
			"id":            c.Id,
			"href":          c.Href,
			"name":          c.Name,
			"description":   c.Description,
			"type":          c.Type_,
			"speed":         c.Speed,
			"location_href": c.Location.Href,
			"state":         c.State,
		})
	}

	return
}
