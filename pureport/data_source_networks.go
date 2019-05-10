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

func dataSourceNetworks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworksRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"account_id": {
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
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNetworksRead(d *schema.ResourceData, m interface{}) error {

	sess := m.(*session.Session)
	accountId := d.Get("account_id")

	ctx := sess.GetSessionContext()

	networks, resp, err := sess.Client.NetworksApi.FindNetworks(ctx, accountId.(string))
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

	nameRegex, nameRegexOk := d.GetOk("name_regex")
	if nameRegexOk {
		r := regexp.MustCompile(nameRegex.(string))
		for _, network := range networks {
			if r.MatchString(network.Name) {
				filteredNetworks = append(filteredNetworks, network)
			}
		}
	} else {
		filteredNetworks = networks
	}

	// Sort the list
	sort.Slice(filteredNetworks, func(i int, j int) bool {
		return filteredNetworks[i].Id < filteredNetworks[j].Id
	})

	// Convert to Map
	out := flattenNetworks(filteredNetworks)
	if err := d.Set("networks", out); err != nil {
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

	for _, network := range networks {

		l := map[string]interface{}{
			"id":          network.Id,
			"href":        network.Href,
			"name":        network.Name,
			"description": network.Description,
			"account_id":  network.Account.Id,
		}

		out = append(out, l)
	}

	return
}
