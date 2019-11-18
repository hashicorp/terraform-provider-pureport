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

func dataSourceNetworks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworksRead,

		Schema: map[string]*schema.Schema{
			"filter": filter.DataSourceFiltersSchema(),
			"account_href": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"networks": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_href": {
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

func dataSourceNetworksRead(d *schema.ResourceData, m interface{}) error {

	config := m.(*configuration.Config)
	accountHref := d.Get("account_href").(string)
	accountId := filepath.Base(accountHref)

	ctx := config.Session.GetSessionContext()

	networks, resp, err := config.Session.Client.NetworksApi.FindNetworks(ctx, accountId)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error when Reading Pureport Network data: %v", err)
	}

	if resp.StatusCode >= 300 {
		d.SetId("")

		if resp.StatusCode == 404 {
			// Need to gracefully handle 404, for refresh
			return nil
		}
		return fmt.Errorf("Error Response while Reading Pureport Network data")
	}

	// Filter the results
	var filteredNetworks []client.Network

	filters, filtersOk := d.GetOk("filter")
	if filtersOk {

		input := make([]interface{}, len(networks))
		for i, x := range networks {
			input[i] = x
		}

		output := filter.FilterType(input, filter.BuildDataSourceFilters(filters.(*schema.Set)))
		for _, x := range output {
			filteredNetworks = append(filteredNetworks, x.(client.Network))
		}

	} else {
		filteredNetworks = networks
	}

	// Sort the list
	sort.Slice(filteredNetworks, func(i int, j int) bool {
		return filteredNetworks[i].Name < filteredNetworks[j].Name
	})

	// Convert to Map
	if err := d.Set("networks", flattenNetworks(filteredNetworks)); err != nil {
		return fmt.Errorf("Error reading networks: %s", err)
	}

	data, err := json.Marshal(networks)
	if err != nil {
		return fmt.Errorf("Error generating Id: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", hashcode.String(string(data))))

	return nil
}

func flattenNetworks(networks []client.Network) (out []map[string]interface{}) {

	for _, n := range networks {

		l := map[string]interface{}{
			"id":           n.Id,
			"href":         n.Href,
			"name":         n.Name,
			"description":  n.Description,
			"account_href": n.Account.Href,
			"tags":         n.Tags,
		}

		out = append(out, l)
	}

	return
}
