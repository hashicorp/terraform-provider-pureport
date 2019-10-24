package pureport

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/configuration"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/filter"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/tags"
)

func dataSourceAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAccountsRead,

		Schema: map[string]*schema.Schema{
			"filter": filter.DataSourceFiltersSchema(),
			"accounts": {
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

						"tags": tags.TagsSchemaComputed(),
					},
				},
			},
		},
	}
}

func dataSourceAccountsRead(d *schema.ResourceData, m interface{}) error {

	config := m.(*configuration.Config)
	filters, filtersOk := d.GetOk("filter")

	ctx := config.Session.GetSessionContext()

	accounts, resp, err := config.Session.Client.AccountsApi.FindAllAccounts(ctx, nil)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error when Reading Pureport Account data: %v", err)
	}

	if resp.StatusCode >= 300 {

		d.SetId("")

		if resp.StatusCode == 404 {
			// Need to gracefully handle 404, for refresh
			return nil
		}

		return fmt.Errorf("Error Response while Reading Pureport Account data")
	}

	// Filter the results
	var filteredAccounts []client.Account
	if filtersOk {

		input := make([]interface{}, len(accounts))
		for i, x := range accounts {
			input[i] = x
		}

		output := filter.FilterType(input, filter.BuildDataSourceFilters(filters.(*schema.Set)))
		for _, x := range output {
			filteredAccounts = append(filteredAccounts, x.(client.Account))
		}

	} else {
		filteredAccounts = accounts
	}

	// Sort the list
	sort.Slice(filteredAccounts, func(i int, j int) bool {
		return filteredAccounts[i].Name < filteredAccounts[j].Name
	})

	// Convert to Map
	out := flattenAccounts(filteredAccounts)
	if err := d.Set("accounts", out); err != nil {
		return fmt.Errorf("Error reading accounts: %s", err)
	}

	data, err := json.Marshal(accounts)
	if err != nil {
		return fmt.Errorf("Error generating Id: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", hashcode.String(string(data))))

	return nil
}

func flattenAccounts(accounts []client.Account) (out []map[string]interface{}) {

	for _, account := range accounts {

		l := map[string]interface{}{
			"id":          account.Id,
			"href":        account.Href,
			"name":        account.Name,
			"description": account.Description,
			"tags":        account.Tags,
		}

		out = append(out, l)
	}

	return
}
