package pureport

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/configuration"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/filter"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/tags"
)

func dataSourceConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConnectionsRead,

		Schema: map[string]*schema.Schema{
			"filter": filter.DataSourceFiltersSchema(),
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

						"tags": tags.TagsSchemaComputed(),
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
		return fmt.Errorf("Error when Reading Connections data: %v", err)
	}

	if resp.StatusCode >= 300 {
		d.SetId("")
		return fmt.Errorf("Error Response while Reading Connections data")
	}

	// Filter the results
	var filteredConnections []client.Connection

	filters, filtersOk := d.GetOk("filter")
	if filtersOk {

		input := make([]interface{}, len(connections))
		for i, x := range connections {
			input[i] = x
		}

		output := filter.FilterType(input, filter.BuildDataSourceFilters(filters.(*schema.Set)))
		for _, x := range output {
			filteredConnections = append(filteredConnections, x.(client.Connection))
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
			"tags":          c.Tags,
		})
	}

	return
}
