package pureport

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testAccDataSourceNetworksConfig_empty = `
data "pureport_accounts" "main" {
	name_regex = "Terraform .*"
}

data "pureport_networks" "empty" {
	account_id = "${data.pureport_accounts.main.accounts.0.id}"
}
`

const testAccDataSourceNetworksConfig_name_regex = `
data "pureport_accounts" "main" {
	name_regex = "Terraform .*"
}

data "pureport_networks" "name_regex" {
	account_id = "${data.pureport_accounts.main.accounts.0.id}"
	name_regex = "Clash.*"
}
`

func TestNetworks_empty(t *testing.T) {

	resourceName := "data.pureport_networks.empty"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworksConfig_empty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceNetworks(resourceName),
					resource.TestCheckResourceAttr(resourceName, "networks.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "networks.0.id"),
					resource.TestCheckResourceAttr(resourceName, "networks.0.href", "/networks/network-EhlpJLhAcHMOmY75J91H3g"),
					resource.TestCheckResourceAttr(resourceName, "networks.0.name", "Siouxsie And The Banshees"),
					resource.TestCheckResourceAttr(resourceName, "networks.0.description", "Test Network #2"),
					resource.TestCheckResourceAttr(resourceName, "networks.0.account_id", "ac-8QVPmcPb_EhapbGHBMAo6Q"),
					resource.TestCheckResourceAttr(resourceName, "networks.1.id", "network-fN6NX6utBCoE5L_H261P4A"),
					resource.TestCheckResourceAttr(resourceName, "networks.1.href", "/networks/network-fN6NX6utBCoE5L_H261P4A"),
					resource.TestCheckResourceAttr(resourceName, "networks.1.name", "The Clash"),
					resource.TestCheckResourceAttr(resourceName, "networks.1.description", "Test Network #1"),
					resource.TestCheckResourceAttr(resourceName, "networks.1.account_id", "ac-8QVPmcPb_EhapbGHBMAo6Q"),
				),
			},
		},
	})
}

func TestNetworks_name_regex(t *testing.T) {

	resourceName := "data.pureport_networks.name_regex"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworksConfig_name_regex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceNetworks(resourceName),
					resource.TestCheckResourceAttr(resourceName, "networks.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networks.0.id", "network-fN6NX6utBCoE5L_H261P4A"),
					resource.TestCheckResourceAttr(resourceName, "networks.0.href", "/networks/network-fN6NX6utBCoE5L_H261P4A"),
					resource.TestCheckResourceAttr(resourceName, "networks.0.name", "The Clash"),
					resource.TestCheckResourceAttr(resourceName, "networks.0.description", "Test Network #1"),
					resource.TestCheckResourceAttr(resourceName, "networks.0.account_id", "ac-8QVPmcPb_EhapbGHBMAo6Q"),
				),
			},
		},
	})
}

func testAccCheckDataSourceNetworks(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		// Find the state object
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find Network data source: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}
