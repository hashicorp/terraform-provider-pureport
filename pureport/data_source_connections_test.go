package pureport

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testAccDataSourceConnectionsConfig_common = `
data "pureport_accounts" "main" {
  name_regex = "Terraform .*"
}

data "pureport_networks" "main" {
  account_href = "${data.pureport_accounts.main.accounts.0.href}"
  name_regex = "Connections"
}
`

const testAccDataSourceConnectionsConfig_empty = testAccDataSourceConnectionsConfig_common + `
data "pureport_connections" "empty" {
  network_href = "${data.pureport_networks.main.networks.0.href}"
}
`

const testAccDataSourceConnectionsConfig_name_regex = testAccDataSourceConnectionsConfig_common + `
data "pureport_connections" "name_filter" {
  network_href = "${data.pureport_networks.main.networks.0.href}"
  name_regex = ".*Test-2"
}
`

func TestConnections_empty(t *testing.T) {

	resourceName := "data.pureport_connections.empty"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectionsConfig_empty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceConnections(resourceName),
					resource.TestCheckResourceAttr(resourceName, "connections.#", "3"),

					resource.TestMatchResourceAttr(resourceName, "connections.0.id", regexp.MustCompile("conn-.{16}")),
					resource.TestMatchResourceAttr(resourceName, "connections.0.href", regexp.MustCompile("/connections/conn-.{16}")),
					resource.TestCheckResourceAttr(resourceName, "connections.0.name", "ConnectionsTest-2"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.description", "ACC Test - 2"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.type", "AWS_DIRECT_CONNECT"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.speed", "50"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.location_href", "/locations/us-sea"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.state", "ACTIVE"),

					resource.TestMatchResourceAttr(resourceName, "connections.1.id", regexp.MustCompile("conn-.{16}")),
					resource.TestMatchResourceAttr(resourceName, "connections.1.href", regexp.MustCompile("/connections/conn-.{16}")),
					resource.TestCheckResourceAttr(resourceName, "connections.1.name", "ConnectionsTest-1"),
					resource.TestCheckResourceAttr(resourceName, "connections.1.description", "ACC Test - 1"),
					resource.TestCheckResourceAttr(resourceName, "connections.1.type", "AWS_DIRECT_CONNECT"),
					resource.TestCheckResourceAttr(resourceName, "connections.1.speed", "50"),
					resource.TestCheckResourceAttr(resourceName, "connections.1.location_href", "/locations/us-sea"),
					resource.TestCheckResourceAttr(resourceName, "connections.1.state", "ACTIVE"),

					resource.TestMatchResourceAttr(resourceName, "connections.2.id", regexp.MustCompile("conn-.{16}")),
					resource.TestMatchResourceAttr(resourceName, "connections.2.href", regexp.MustCompile("/connections/conn-.{16}")),
					resource.TestCheckResourceAttr(resourceName, "connections.2.name", "ConnectionsTest-3"),
					resource.TestCheckResourceAttr(resourceName, "connections.2.description", "ACC Test - 3"),
					resource.TestCheckResourceAttr(resourceName, "connections.2.type", "AWS_DIRECT_CONNECT"),
					resource.TestCheckResourceAttr(resourceName, "connections.2.speed", "50"),
					resource.TestCheckResourceAttr(resourceName, "connections.2.location_href", "/locations/us-sea"),
					resource.TestCheckResourceAttr(resourceName, "connections.2.state", "ACTIVE"),
				),
			},
		},
	})
}

func TestConnections_name_regex(t *testing.T) {

	resourceName := "data.pureport_connections.name_filter"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectionsConfig_name_regex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceConnections(resourceName),
					resource.TestCheckResourceAttr(resourceName, "connections.#", "1"),

					resource.TestMatchResourceAttr(resourceName, "connections.0.id", regexp.MustCompile("conn-.{16}")),
					resource.TestMatchResourceAttr(resourceName, "connections.0.href", regexp.MustCompile("/connections/conn-.{16}")),
					resource.TestCheckResourceAttr(resourceName, "connections.0.name", "ConnectionsTest-2"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.description", "ACC Test - 2"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.type", "AWS_DIRECT_CONNECT"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.speed", "50"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.location_href", "/locations/us-sea"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.state", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckDataSourceConnections(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		// Find the state object
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find Connections data source: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}
