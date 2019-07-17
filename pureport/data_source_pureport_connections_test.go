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
  filter {
    name = "Name"
    values = ["Terraform .*"]
  }
}

data "pureport_networks" "main" {
  account_href = "${data.pureport_accounts.main.accounts.0.href}"
  filter {
    name = "Name"
    values = ["Connections"]
  }
}
`

const testAccDataSourceConnectionsConfig_empty = testAccDataSourceConnectionsConfig_common + `
data "pureport_connections" "empty" {
  network_href = "${data.pureport_networks.main.networks.0.href}"
}
`

const testAccDataSourceConnectionsConfig_name_filter = testAccDataSourceConnectionsConfig_common + `
data "pureport_connections" "name_filter" {
  network_href = "${data.pureport_networks.main.networks.0.href}"
  filter {
    name = "Name"
    values = [".*Test-2"]
  }
}
`

func TestDataSourceConnections_empty(t *testing.T) {

	resourceName := "data.pureport_connections.empty"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectionsConfig_empty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceConnections(resourceName),
					resource.TestCheckResourceAttr(resourceName, "connections.#", "3"),

					testAccCheckDataSourceConnectionsTest1(resourceName, "connections.0"),
					testAccCheckDataSourceConnectionsTest2(resourceName, "connections.1"),
					testAccCheckDataSourceConnectionsTest3(resourceName, "connections.2"),
				),
			},
		},
	})
}

func TestDataSourceConnections_name_filter(t *testing.T) {

	resourceName := "data.pureport_connections.name_filter"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConnectionsConfig_name_filter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceConnections(resourceName),
					resource.TestCheckResourceAttr(resourceName, "connections.#", "1"),
					testAccCheckDataSourceConnectionsTest2(resourceName, "connections.0"),
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

func testAccCheckDataSourceConnectionsTest1(resourceName, connection string) resource.TestCheckFunc {
	if testEnvironmentName == "Production" {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestMatchResourceAttr(resourceName, connection+".id", regexp.MustCompile("conn-.{16}")),
			resource.TestMatchResourceAttr(resourceName, connection+".href", regexp.MustCompile("/connections/conn-.{16}")),
			resource.TestCheckResourceAttr(resourceName, connection+".name", "ConnectionsTest-1"),
			resource.TestCheckResourceAttr(resourceName, connection+".description", "ACC Test - 1"),
			resource.TestCheckResourceAttr(resourceName, connection+".type", "AWS_DIRECT_CONNECT"),
			resource.TestCheckResourceAttr(resourceName, connection+".speed", "50"),
			resource.TestCheckResourceAttr(resourceName, connection+".location_href", "/locations/us-wdc"),
			resource.TestCheckResourceAttr(resourceName, connection+".state", "ACTIVE"),

			resource.TestCheckResourceAttr(resourceName, connection+".tags.#", "0"),
		)
	}

	return resource.ComposeAggregateTestCheckFunc(
		resource.TestMatchResourceAttr(resourceName, connection+".id", regexp.MustCompile("conn-.{16}")),
		resource.TestMatchResourceAttr(resourceName, connection+".href", regexp.MustCompile("/connections/conn-.{16}")),
		resource.TestCheckResourceAttr(resourceName, connection+".name", "ConnectionsTest-1"),
		resource.TestCheckResourceAttr(resourceName, connection+".description", "ACC Test - 1"),
		resource.TestCheckResourceAttr(resourceName, connection+".type", "AWS_DIRECT_CONNECT"),
		resource.TestCheckResourceAttr(resourceName, connection+".speed", "50"),
		resource.TestCheckResourceAttr(resourceName, connection+".location_href", "/locations/us-sea"),
		resource.TestCheckResourceAttr(resourceName, connection+".state", "ACTIVE"),

		resource.TestCheckResourceAttr(resourceName, connection+".tags.#", "0"),
	)
}

func testAccCheckDataSourceConnectionsTest2(resourceName, connection string) resource.TestCheckFunc {
	if testEnvironmentName == "Production" {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestMatchResourceAttr(resourceName, connection+".id", regexp.MustCompile("conn-.{16}")),
			resource.TestMatchResourceAttr(resourceName, connection+".href", regexp.MustCompile("/connections/conn-.{16}")),
			resource.TestCheckResourceAttr(resourceName, connection+".name", "ConnectionsTest-2"),
			resource.TestCheckResourceAttr(resourceName, connection+".description", "ACC Test - 2"),
			resource.TestCheckResourceAttr(resourceName, connection+".type", "AWS_DIRECT_CONNECT"),
			resource.TestCheckResourceAttr(resourceName, connection+".speed", "50"),
			resource.TestCheckResourceAttr(resourceName, connection+".location_href", "/locations/us-sjc"),
			resource.TestCheckResourceAttr(resourceName, connection+".state", "ACTIVE"),

			resource.TestCheckResourceAttr(resourceName, connection+".tags.#", "0"),
		)
	}
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestMatchResourceAttr(resourceName, connection+".id", regexp.MustCompile("conn-.{16}")),
		resource.TestMatchResourceAttr(resourceName, connection+".href", regexp.MustCompile("/connections/conn-.{16}")),
		resource.TestCheckResourceAttr(resourceName, connection+".name", "ConnectionsTest-2"),
		resource.TestCheckResourceAttr(resourceName, connection+".description", "ACC Test - 2"),
		resource.TestCheckResourceAttr(resourceName, connection+".type", "AWS_DIRECT_CONNECT"),
		resource.TestCheckResourceAttr(resourceName, connection+".speed", "50"),
		resource.TestCheckResourceAttr(resourceName, connection+".location_href", "/locations/us-sea"),
		resource.TestCheckResourceAttr(resourceName, connection+".state", "ACTIVE"),

		resource.TestCheckResourceAttr(resourceName, connection+".tags.#", "0"),
	)
}

func testAccCheckDataSourceConnectionsTest3(resourceName, connection string) resource.TestCheckFunc {
	if testEnvironmentName == "Production" {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestMatchResourceAttr(resourceName, connection+".id", regexp.MustCompile("conn-.{16}")),
			resource.TestMatchResourceAttr(resourceName, connection+".href", regexp.MustCompile("/connections/conn-.{16}")),
			resource.TestCheckResourceAttr(resourceName, connection+".name", "ConnectionsTest-3"),
			resource.TestCheckResourceAttr(resourceName, connection+".description", "ACC Test - 3"),
			resource.TestCheckResourceAttr(resourceName, connection+".type", "AWS_DIRECT_CONNECT"),
			resource.TestCheckResourceAttr(resourceName, connection+".speed", "50"),
			resource.TestCheckResourceAttr(resourceName, connection+".location_href", "/locations/us-chi"),
			resource.TestCheckResourceAttr(resourceName, connection+".state", "ACTIVE"),

			resource.TestCheckResourceAttr(resourceName, connection+".tags.#", "0"),
		)
	}
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestMatchResourceAttr(resourceName, connection+".id", regexp.MustCompile("conn-.{16}")),
		resource.TestMatchResourceAttr(resourceName, connection+".href", regexp.MustCompile("/connections/conn-.{16}")),
		resource.TestCheckResourceAttr(resourceName, connection+".name", "ConnectionsTest-3"),
		resource.TestCheckResourceAttr(resourceName, connection+".description", "ACC Test - 3"),
		resource.TestCheckResourceAttr(resourceName, connection+".type", "AWS_DIRECT_CONNECT"),
		resource.TestCheckResourceAttr(resourceName, connection+".speed", "50"),
		resource.TestCheckResourceAttr(resourceName, connection+".location_href", "/locations/us-sea"),
		resource.TestCheckResourceAttr(resourceName, connection+".state", "ACTIVE"),

		resource.TestCheckResourceAttr(resourceName, connection+".tags.#", "0"),
	)
}
