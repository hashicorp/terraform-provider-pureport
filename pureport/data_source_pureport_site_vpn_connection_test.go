package pureport

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const testAccDataSourceSiteVPNConnectionConfig_common = `
data "pureport_accounts" "main" {
  name_regex = "Terraform .*"
}

data "pureport_networks" "main" {
  account_href = "${data.pureport_accounts.main.accounts.0.href}"
  name_regex = "A Flock of Seagulls"
}

data "pureport_connections" "main" {
  network_href = "${data.pureport_networks.main.networks.0.href}"
  name_regex = "SiteVPN"
}
`

const testAccDataSourceSiteVPNConnectionConfig_basic = testAccDataSourceSiteVPNConnectionConfig_common + `
data "pureport_site_vpn_connection" "basic" {
  connection_id = "${data.pureport_connections.main.connections.0.id}"
}
`

func TestSiteVPNConnectionDataSource_basic(t *testing.T) {

	resourceName := "data.pureport_site_vpn_connection.basic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSiteVPNConnectionConfig_basic,
				Check: resource.ComposeTestCheckFunc(

					resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile("conn-.{16}")),
						resource.TestCheckResourceAttr(resourceName, "name", "SiteVPN_RouteBasedBGP_DataSource"),
						resource.TestCheckResourceAttr(resourceName, "description", ""),
						resource.TestCheckResourceAttr(resourceName, "speed", "100"),
						resource.TestMatchResourceAttr(resourceName, "href", regexp.MustCompile("/connections/conn-.{16}")),
						resource.TestCheckResourceAttr(resourceName, "high_availability", "true"),
						resource.TestCheckResourceAttr(resourceName, "state", "ACTIVE"),
						resource.TestCheckResourceAttr(resourceName, "location_href", "/locations/us-chi"),
						resource.TestMatchResourceAttr(resourceName, "network_href", regexp.MustCompile("/networks/network-.{16}")),
						resource.TestCheckResourceAttr(resourceName, "cloud_service_hrefs.#", "0"),

						resource.TestCheckResourceAttr(resourceName, "gateways.#", "2"),

						resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "gateways.0.availability_domain", "PRIMARY"),
							resource.TestCheckResourceAttr(resourceName, "gateways.0.name", "SITE_IPSEC_VPN"),
							resource.TestCheckResourceAttr(resourceName, "gateways.0.description", ""),
							resource.TestCheckResourceAttr(resourceName, "gateways.0.link_state", "PENDING"),
							resource.TestCheckResourceAttr(resourceName, "gateways.0.customer_asn", "30000"),
							resource.TestMatchResourceAttr(resourceName, "gateways.0.customer_ip", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
							resource.TestCheckResourceAttr(resourceName, "gateways.0.pureport_asn", "394351"),
							resource.TestMatchResourceAttr(resourceName, "gateways.0.pureport_ip", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
							resource.TestCheckResourceAttr(resourceName, "gateways.0.bgp_password", ""),
							resource.TestMatchResourceAttr(resourceName, "gateways.0.peering_subnet", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
							resource.TestCheckResourceAttr(resourceName, "gateways.0.public_nat_ip", ""),
							resource.TestCheckResourceAttr(resourceName, "gateways.0.customer_gateway_ip", "123.123.123.123"),
							resource.TestMatchResourceAttr(resourceName, "gateways.0.customer_vti_ip", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
							resource.TestMatchResourceAttr(resourceName, "gateways.0.pureport_gateway_ip", regexp.MustCompile("45.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}")),
							resource.TestMatchResourceAttr(resourceName, "gateways.0.pureport_vti_ip", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
							resource.TestCheckResourceAttr(resourceName, "gateways.0.vpn_auth_type", "PSK"),
							resource.TestCheckResourceAttrSet(resourceName, "gateways.0.vpn_auth_key"),
						),

						resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "gateways.1.availability_domain", "SECONDARY"),
							resource.TestCheckResourceAttr(resourceName, "gateways.1.name", "SITE_IPSEC_VPN 2"),
							resource.TestCheckResourceAttr(resourceName, "gateways.1.description", ""),
							resource.TestCheckResourceAttr(resourceName, "gateways.1.link_state", "PENDING"),
							resource.TestCheckResourceAttr(resourceName, "gateways.1.customer_asn", "30000"),
							resource.TestMatchResourceAttr(resourceName, "gateways.1.customer_ip", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
							resource.TestCheckResourceAttr(resourceName, "gateways.1.pureport_asn", "394351"),
							resource.TestMatchResourceAttr(resourceName, "gateways.1.pureport_ip", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
							resource.TestCheckResourceAttr(resourceName, "gateways.1.bgp_password", ""),
							resource.TestMatchResourceAttr(resourceName, "gateways.1.peering_subnet", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
							resource.TestCheckResourceAttr(resourceName, "gateways.1.public_nat_ip", ""),
							resource.TestCheckResourceAttr(resourceName, "gateways.1.customer_gateway_ip", "124.124.124.124"),
							resource.TestMatchResourceAttr(resourceName, "gateways.1.customer_vti_ip", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
							resource.TestMatchResourceAttr(resourceName, "gateways.1.pureport_gateway_ip", regexp.MustCompile("45.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}")),
							resource.TestMatchResourceAttr(resourceName, "gateways.1.pureport_vti_ip", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
							resource.TestCheckResourceAttr(resourceName, "gateways.1.vpn_auth_type", "PSK"),
							resource.TestCheckResourceAttrSet(resourceName, "gateways.1.vpn_auth_key"),
						),
					),
				),
			},
		},
	})
}
