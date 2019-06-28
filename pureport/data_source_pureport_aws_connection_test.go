package pureport

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const testAccDataSourceAwsConnectionConfig_common = `
data "pureport_accounts" "main" {
  name_regex = "Terraform .*"
}

data "pureport_networks" "main" {
  account_href = "${data.pureport_accounts.main.accounts.0.href}"
  name_regex = "A Flock of Seagulls"
}

data "pureport_connections" "main" {
  network_href = "${data.pureport_networks.main.networks.0.href}"
  name_regex = "AWS"
}
`

const testAccDataSourceAwsConnectionConfig_basic = testAccDataSourceAwsConnectionConfig_common + `
data "pureport_aws_connection" "basic" {
  connection_id = "${data.pureport_connections.main.connections.0.id}"
}
`

func TestAwsConnectionDataSource_basic(t *testing.T) {

	resourceName := "data.pureport_aws_connection.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsConnectionConfig_basic,
				Check: resource.ComposeTestCheckFunc(

					resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile("conn-.{16}")),
						resource.TestMatchResourceAttr(resourceName, "aws_account_id", regexp.MustCompile("[0-9]{12}")),
						resource.TestCheckResourceAttr(resourceName, "aws_region", "us-west-1"),
						resource.TestCheckResourceAttr(resourceName, "speed", "50"),
						resource.TestCheckResourceAttr(resourceName, "peering_type", "PRIVATE"),
						resource.TestMatchResourceAttr(resourceName, "href", regexp.MustCompile("/connections/conn-.{16}")),
						resource.TestCheckResourceAttr(resourceName, "name", "AWS Connection DataSource"),
						resource.TestCheckResourceAttr(resourceName, "description", ""),
						resource.TestCheckResourceAttr(resourceName, "state", "ACTIVE"),
						resource.TestCheckResourceAttr(resourceName, "high_availability", "false"),
						resource.TestCheckResourceAttr(resourceName, "location_href", "/locations/us-sjc"),
						resource.TestMatchResourceAttr(resourceName, "network_href", regexp.MustCompile("/networks/network-.{16}")),
						resource.TestCheckResourceAttr(resourceName, "cloud_service_hrefs.#", "0"),

						resource.TestCheckResourceAttr(resourceName, "gateways.#", "1"),

						resource.TestCheckResourceAttr(resourceName, "gateways.0.availability_domain", "PRIMARY"),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.name", "AWS_DIRECT_CONNECT"),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.description", ""),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.link_state", "PENDING"),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.customer_asn", "64512"),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.customer_ip", "169.254.1.2/30"),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.pureport_asn", "394351"),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.pureport_ip", "169.254.1.1/30"),
						resource.TestCheckResourceAttrSet(resourceName, "gateways.0.bgp_password"),
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
