package pureport

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccDataSourceAzureConnectionConfig_common = `
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
    values = ["Azure"]
  }
}
`

const testAccDataSourceAzureConnectionConfig_basic = testAccDataSourceAzureConnectionConfig_common + `
data "pureport_azure_connection" "basic" {
  connection_id = data.pureport_connections.main.connections.0.id
}
`

func TestDataSourceAzureConnection_basic(t *testing.T) {

	resourceName := "data.pureport_azure_connection.basic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureConnectionConfig_basic,
				Check: resource.ComposeTestCheckFunc(

					resource.ComposeAggregateTestCheckFunc(

						resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile("conn-.{16}")),

						resource.TestCheckResourceAttr(resourceName, "name", "AzureExpressRoute_DataSource"),
						resource.TestCheckResourceAttr(resourceName, "description", "Some random description"),
						resource.TestCheckResourceAttr(resourceName, "speed", "100"),
						resource.TestCheckResourceAttr(resourceName, "high_availability", "true"),
						resource.TestMatchResourceAttr(resourceName, "service_key", regexp.MustCompile("[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}")),

						resource.TestCheckResourceAttr(resourceName, "peering_type", "PRIVATE"),
						resource.TestMatchResourceAttr(resourceName, "href", regexp.MustCompile("/connections/conn-.{16}")),
						resource.TestCheckResourceAttr(resourceName, "state", "ACTIVE"),
						resource.TestCheckResourceAttr(resourceName, "location_href", "/locations/us-sea"),
						resource.TestMatchResourceAttr(resourceName, "network_href", regexp.MustCompile("/networks/network-.{16}")),
						resource.TestCheckResourceAttr(resourceName, "cloud_service_hrefs.#", "0"),

						resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),

						resource.TestCheckResourceAttr(resourceName, "gateways.#", "2"),

						resource.TestCheckResourceAttr(resourceName, "gateways.0.availability_domain", "PRIMARY"),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.name", "AzureExpressRoute_DataSource - Primary"),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.description", ""),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.customer_asn", "12076"),
						resource.TestMatchResourceAttr(resourceName, "gateways.0.customer_ip", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}/30")),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.pureport_asn", "394351"),
						resource.TestMatchResourceAttr(resourceName, "gateways.0.pureport_ip", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}/30")),
						resource.TestCheckResourceAttrSet(resourceName, "gateways.0.bgp_password"),
						resource.TestMatchResourceAttr(resourceName, "gateways.0.peering_subnet", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.public_nat_ip", ""),
						resource.TestCheckResourceAttrSet(resourceName, "gateways.0.vlan"),
						resource.TestCheckResourceAttrSet(resourceName, "gateways.0.remote_id"),

						resource.TestCheckResourceAttr(resourceName, "gateways.1.availability_domain", "SECONDARY"),
						resource.TestCheckResourceAttr(resourceName, "gateways.1.name", "AzureExpressRoute_DataSource - Secondary"),
						resource.TestCheckResourceAttr(resourceName, "gateways.1.description", ""),
						resource.TestCheckResourceAttr(resourceName, "gateways.1.customer_asn", "12076"),
						resource.TestMatchResourceAttr(resourceName, "gateways.1.customer_ip", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}/30")),
						resource.TestCheckResourceAttr(resourceName, "gateways.1.pureport_asn", "394351"),
						resource.TestMatchResourceAttr(resourceName, "gateways.1.pureport_ip", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}/30")),
						resource.TestCheckResourceAttrSet(resourceName, "gateways.1.bgp_password"),
						resource.TestMatchResourceAttr(resourceName, "gateways.1.peering_subnet", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
						resource.TestCheckResourceAttr(resourceName, "gateways.1.public_nat_ip", ""),
						resource.TestCheckResourceAttrSet(resourceName, "gateways.1.vlan"),
						resource.TestCheckResourceAttrSet(resourceName, "gateways.1.remote_id"),
					),
				),
			},
		},
	})
}
