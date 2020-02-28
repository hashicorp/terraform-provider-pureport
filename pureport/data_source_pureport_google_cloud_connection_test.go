package pureport

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccDataSourceGoogleConnectionConfig_common = `
data "pureport_accounts" "main" {
  filter {
    name = "Name"
    values = ["Terraform .*"]
  }
}

data "pureport_networks" "main" {
  account_href = data.pureport_accounts.main.accounts.0.href
  filter {
    name = "Name"
    values = ["A Flock of Seagulls"]
  }
}

data "pureport_connections" "main" {
  network_href = data.pureport_networks.main.networks.0.href
  filter {
    name = "Name"
    values = ["Google"]
  }
}
`

const testAccDataSourceGoogleConnectionConfig_basic = testAccDataSourceGoogleConnectionConfig_common + `
data "pureport_google_cloud_connection" "basic" {
  connection_id = data.pureport_connections.main.connections.0.id
}
`

func TestDataSourceGoogleConnection_basic(t *testing.T) {

	resourceName := "data.pureport_google_cloud_connection.basic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleConnectionConfig_basic,
				Check: resource.ComposeTestCheckFunc(

					resource.ComposeAggregateTestCheckFunc(

						testAccCheckDataSourceGoogleConnection(resourceName),

						resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile("conn-.{16}")),
						resource.TestCheckResourceAttr(resourceName, "name", "GoogleCloud_DataSource"),
						resource.TestCheckResourceAttr(resourceName, "description", ""),
						resource.TestCheckResourceAttr(resourceName, "speed", "50"),
						resource.TestMatchResourceAttr(resourceName, "href", regexp.MustCompile("/connections/conn-.{16}")),
						resource.TestCheckResourceAttr(resourceName, "high_availability", "false"),
						resource.TestCheckResourceAttr(resourceName, "state", "ACTIVE"),
						resource.TestMatchResourceAttr(resourceName, "network_href", regexp.MustCompile("/networks/network-.{16}")),
						resource.TestCheckResourceAttr(resourceName, "cloud_service_hrefs.#", "0"),

						resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),

						resource.TestCheckResourceAttr(resourceName, "gateways.#", "1"),

						resource.TestCheckResourceAttr(resourceName, "gateways.0.availability_domain", "PRIMARY"),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.name", "GoogleCloud_DataSource - Primary"),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.description", ""),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.customer_asn", "16550"),
						resource.TestMatchResourceAttr(resourceName, "gateways.0.customer_ip", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.pureport_asn", "394351"),
						resource.TestMatchResourceAttr(resourceName, "gateways.0.pureport_ip", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.bgp_password", ""),
						resource.TestMatchResourceAttr(resourceName, "gateways.0.peering_subnet", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.public_nat_ip", ""),
						resource.TestCheckResourceAttrSet(resourceName, "gateways.0.vlan"),
						resource.TestCheckResourceAttrSet(resourceName, "gateways.0.remote_id"),
					),
				),
			},
		},
	})
}

func testAccCheckDataSourceGoogleConnection(resourceName string) resource.TestCheckFunc {
	if testEnvironmentName == "Production" {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceName, "location_href", "/locations/us-chi"),
		)
	}

	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "location_href", "/locations/us-sea"),
	)
}
