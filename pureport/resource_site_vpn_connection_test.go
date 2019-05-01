package pureport

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/pureport/pureport-sdk-go/pureport/session"
)

const testAccResourceSiteVPNConnectionConfig_basic = `
data "pureport_accounts" "main" {
	name_regex = "Terraform"
}

data "pureport_cloud_regions" "main" {
	name_regex = "Oregon"
}

data "pureport_locations" "main" {
	name_regex = "^Ral*"
}

data "pureport_networks" "main" {
	account_id = "${data.pureport_accounts.main.accounts.0.id}"
	name_regex = "Bansh.*"
}

resource "pureport_site_vpn_connection" "main" {
	name = "SiteVPNTest"
	speed = "100"
	high_availability = true

	location {
		id = "${data.pureport_locations.main.locations.0.id}"
		href = "${data.pureport_locations.main.locations.0.href}"
	}
	network {
		id = "${data.pureport_networks.main.networks.0.id}"
		href = "${data.pureport_networks.main.networks.0.href}"
	}

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
				),
			},
		},
	})
}

func testAccCheckResourceSiteVPNConnection(name string, instance *client.SiteIpSecVpnConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		sess, ok := testAccProvider.Meta().(*session.Session)
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

		ctx := sess.GetSessionContext()
		found, resp, err := sess.Client.ConnectionsApi.GetConnection(ctx, id)

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

	sess, ok := testAccProvider.Meta().(*session.Session)
	if !ok {
		return fmt.Errorf("Error getting Pureport client")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pureport_aws_connection" {
			continue
		}

		id := rs.Primary.ID

		ctx := sess.GetSessionContext()
		_, resp, err := sess.Client.ConnectionsApi.GetConnection(ctx, id)

		if err != nil && resp.StatusCode != 404 {
			return fmt.Errorf("should not get error for SiteVPN Connection with ID %s after delete: %s", id, err)
		}
	}

	return nil
}
