package pureport

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pureport/pureport-sdk-go/pureport/client"
)

const testAccResourceSiteVPNConnectionConfig_basic = `
data "pureport_accounts" "main" {
	name_regex = "Terraform"
}

data "pureport_cloud_regions" "main" {
	name_regex = "Oregon"
}

data "pureport_locations" "main" {
	name_regex = "^Sea*"
}

data "pureport_networks" "main" {
	account_href = "${data.pureport_accounts.main.accounts.0.href}"
	name_regex = "Bansh.*"
}

resource "pureport_site_vpn_connection" "main" {
	name = "SiteVPNTest"
	speed = "100"
	high_availability = true

	location_href = "${data.pureport_locations.main.locations.0.href}"
	network_href = "${data.pureport_networks.main.networks.0.href}"

	ike_version = "V2"
	
	routing_type = "ROUTE_BASED_BGP"
	customer_asn = 30000

	primary_customer_router_ip = "123.123.123.123"
	secondary_customer_router_ip = "124.124.124.124"
}
`

func TestSiteVPNConnection_basic(t *testing.T) {

	resourceName := "pureport_site_vpn_connection.main"
	var instance client.SiteIpSecVpnConnection

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSiteVPNConnectionConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceSiteVPNConnection(resourceName, &instance),
					resource.TestCheckResourceAttrPtr(resourceName, "id", &instance.Id),
					resource.TestCheckResourceAttr(resourceName, "name", "SiteVPNTest"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "speed", "100"),
					resource.TestCheckResourceAttr(resourceName, "high_availability", "true"),

					resource.TestCheckResourceAttr(resourceName, "gateways.#", "2"),

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
					resource.TestMatchResourceAttr(resourceName, "gateways.0.pureport_gateway_ip", regexp.MustCompile("45.56.[0-9]{1,3}.[0-9]{1,3}")),
					resource.TestMatchResourceAttr(resourceName, "gateways.0.pureport_vti_ip", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.vpn_auth_type", "PSK"),

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
					resource.TestMatchResourceAttr(resourceName, "gateways.1.pureport_gateway_ip", regexp.MustCompile("45.56.[0-9]{1,3}.[0-9]{1,3}")),
					resource.TestMatchResourceAttr(resourceName, "gateways.1.pureport_vti_ip", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.vpn_auth_type", "PSK"),
				),
			},
		},
	})
}

func testAccCheckResourceSiteVPNConnection(name string, instance *client.SiteIpSecVpnConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		config, ok := testAccProvider.Meta().(*Config)
		if !ok {
			return fmt.Errorf("Error getting Pureport client")
		}

		// Find the state object
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find SiteVPN Connection resource: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id := rs.Primary.ID

		ctx := config.Session.GetSessionContext()
		found, resp, err := config.Session.Client.ConnectionsApi.GetConnection(ctx, id)

		if err != nil {
			return fmt.Errorf("receive error when requesting SiteVPN Connection %s", id)
		}

		if resp.StatusCode != 200 {
			fmt.Errorf("Error getting SiteVPN Connection ID %s: %s", id, err)
		}

		*instance = found.(client.SiteIpSecVpnConnection)

		return nil
	}
}

func testAccCheckSiteVPNConnectionDestroy(s *terraform.State) error {

	config, ok := testAccProvider.Meta().(*Config)
	if !ok {
		return fmt.Errorf("Error getting Pureport client")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pureport_aws_connection" {
			continue
		}

		id := rs.Primary.ID

		ctx := config.Session.GetSessionContext()
		_, resp, err := config.Session.Client.ConnectionsApi.GetConnection(ctx, id)

		if err != nil && resp.StatusCode != 404 {
			return fmt.Errorf("should not get error for SiteVPN Connection with ID %s after delete: %s", id, err)
		}
	}

	return nil
}
