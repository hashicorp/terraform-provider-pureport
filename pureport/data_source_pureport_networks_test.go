package pureport

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testAccDataSourceNetworksConfig_empty = `
data "pureport_accounts" "main" {
	name_regex = "Terraform .*"
}

data "pureport_networks" "empty" {
	account_href = "${data.pureport_accounts.main.accounts.0.href}"
}
`

const testAccDataSourceNetworksConfig_name_regex = `
data "pureport_accounts" "main" {
	name_regex = "Terraform .*"
}

data "pureport_networks" "name_regex" {
	account_href = "${data.pureport_accounts.main.accounts.0.href}"
	name_regex = "Clash.*"
}
`

func TestNetworksDataSource_empty(t *testing.T) {

	resourceName := "data.pureport_networks.empty"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworksConfig_empty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceNetworks(resourceName),
					resource.TestCheckResourceAttr(resourceName, "networks.#", "4"),

					resource.TestMatchResourceAttr(resourceName, "networks.0.id", regexp.MustCompile("network-.{16}")),
					resource.TestMatchResourceAttr(resourceName, "networks.0.href", regexp.MustCompile("/networks/network-.{16}")),
					resource.TestCheckResourceAttr(resourceName, "networks.0.name", "A Flock of Seagulls"),
					resource.TestCheckResourceAttr(resourceName, "networks.0.description", "Test Network DataSource"),
					resource.TestMatchResourceAttr(resourceName, "networks.0.account_href", regexp.MustCompile("/accounts/ac-.{16}")),

					resource.TestMatchResourceAttr(resourceName, "networks.1.id", regexp.MustCompile("network-.{16}")),
					resource.TestMatchResourceAttr(resourceName, "networks.1.href", regexp.MustCompile("/networks/network-.{16}")),
					resource.TestCheckResourceAttr(resourceName, "networks.1.name", "Connections"),
					resource.TestCheckResourceAttr(resourceName, "networks.1.description", "Data Source Testing"),
					resource.TestMatchResourceAttr(resourceName, "networks.1.account_href", regexp.MustCompile("/accounts/ac-.{16}")),

					resource.TestMatchResourceAttr(resourceName, "networks.2.id", regexp.MustCompile("network-.{16}")),
					resource.TestMatchResourceAttr(resourceName, "networks.2.href", regexp.MustCompile("/networks/network-.{16}")),
					resource.TestCheckResourceAttr(resourceName, "networks.2.name", "Siouxsie & The Banshees"),
					resource.TestCheckResourceAttr(resourceName, "networks.2.description", "Test Network #2"),
					resource.TestMatchResourceAttr(resourceName, "networks.2.account_href", regexp.MustCompile("/accounts/ac-.{16}")),

					resource.TestMatchResourceAttr(resourceName, "networks.3.id", regexp.MustCompile("network-.{16}")),
					resource.TestMatchResourceAttr(resourceName, "networks.3.href", regexp.MustCompile("/networks/network-.{16}")),
					resource.TestCheckResourceAttr(resourceName, "networks.3.name", "The Clash"),
					resource.TestCheckResourceAttr(resourceName, "networks.3.description", "Test Network #1"),
					resource.TestMatchResourceAttr(resourceName, "networks.3.account_href", regexp.MustCompile("/accounts/ac-.{16}")),
				),
			},
		},
	})
}

func TestNetworksDataSource_name_regex(t *testing.T) {

	resourceName := "data.pureport_networks.name_regex"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworksConfig_name_regex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceNetworks(resourceName),

					resource.TestCheckResourceAttr(resourceName, "networks.#", "1"),

					resource.TestMatchResourceAttr(resourceName, "networks.0.id", regexp.MustCompile("network-.{16}")),
					resource.TestMatchResourceAttr(resourceName, "networks.0.href", regexp.MustCompile("/networks/network-.{16}")),
					resource.TestCheckResourceAttr(resourceName, "networks.0.name", "The Clash"),
					resource.TestCheckResourceAttr(resourceName, "networks.0.description", "Test Network #1"),
					resource.TestMatchResourceAttr(resourceName, "networks.0.account_href", regexp.MustCompile("/accounts/ac-.{16}")),
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
