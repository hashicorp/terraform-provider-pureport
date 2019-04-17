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
	"github.com/pureport/pureport-sdk-go/pureport/session"
	"github.com/pureport/pureport-sdk-go/pureport/swagger"
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
	nameRegex, nameRegexOk := d.GetOk("name_regex")
	accountId, accountIdOk := d.GetOk("account_id")

	if !accountIdOk {
		log.Printf("[Error] Invalid Pureport Account ID")
	}

	ctx := sess.GetSessionContext()

	networks, resp, err := sess.Client.NetworksApi.FindNetworks(ctx, accountId.(string))
	if err != nil {
		log.Printf("[Error] Error when Reading Pureport Network data")
		d.SetId("")
		return nil
	}

	if resp.StatusCode != 200 {
		log.Printf("[Error] Error Response while Reading Pureport Network data")
		d.SetId("")
		return nil
	}

	// Filter the results
	var filteredNetworks []swagger.Network
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

func flattenNetworks(networks []swagger.Network) (out []map[string]interface{}) {

	for _, network := range networks {

		l := map[string]interface{}{
			"id":          network.Id,
			"name":        network.Name,
			"description": network.Description,
			"account_id":  network.Account.Id,
		}

		out = append(out, l)
	}

	return
}
