package pureport

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pureport/pureport-sdk-go/pureport/client"
)

const testAccResourceDummyConnectionConfig_basic = `
data "pureport_accounts" "main" {
	name_regex = "Terraform"
}

data "pureport_cloud_regions" "main" {
	name_regex = "Oregon"
}

data "pureport_locations" "main" {
	name_regex = ".*ttle.*"
}

data "pureport_networks" "main" {
	account_href = "${data.pureport_accounts.main.accounts.0.href}"
	name_regex = "Bansh.*"
}

resource "pureport_dummy_connection" "main" {
	name = "DummyTest"
	speed = "100"
	high_availability = true

	location_href = "${data.pureport_locations.main.locations.0.href}"
	network_href = "${data.pureport_networks.main.networks.0.href}"
}
`

func TestDummyConnection_basic(t *testing.T) {

	resourceName := "pureport_dummy_connection.main"
	var instance client.DummyConnection

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDummyConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDummyConnectionConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceDummyConnection(resourceName, &instance),
					resource.TestCheckResourceAttrPtr(resourceName, "id", &instance.Id),
					resource.TestCheckResourceAttr(resourceName, "name", "DummyTest"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "speed", "100"),
					resource.TestCheckResourceAttr(resourceName, "high_availability", "true"),

					resource.TestCheckResourceAttr(resourceName, "gateways.#", "2"),

					resource.TestCheckResourceAttr(resourceName, "gateways.0.availability_domain", "PRIMARY"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.name", "DUMMY"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.link_state", "PENDING"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.customer_asn", "65000"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.customer_ip", "169.254.100.2/30"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.pureport_asn", "394351"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.pureport_ip", "169.254.100.1/30"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.bgp_password", ""),
					resource.TestMatchResourceAttr(resourceName, "gateways.0.peering_subnet", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.public_nat_ip", ""),
					resource.TestCheckResourceAttrSet(resourceName, "gateways.0.vlan"),

					resource.TestCheckResourceAttr(resourceName, "gateways.1.availability_domain", "SECONDARY"),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.name", "DUMMY 2"),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.description", ""),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.link_state", "PENDING"),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.customer_asn", "65000"),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.customer_ip", "169.254.200.2/30"),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.pureport_asn", "394351"),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.pureport_ip", "169.254.200.1/30"),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.bgp_password", ""),
					resource.TestMatchResourceAttr(resourceName, "gateways.1.peering_subnet", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.public_nat_ip", ""),
					resource.TestCheckResourceAttrSet(resourceName, "gateways.1.vlan"),
				),
			},
		},
	})
}

func testAccCheckResourceDummyConnection(name string, instance *client.DummyConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		config, ok := testAccProvider.Meta().(*Config)
		if !ok {
			return fmt.Errorf("Error getting Pureport client")
		}

		// Find the state object
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find Dummy Connection resource: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id := rs.Primary.ID

		ctx := config.Session.GetSessionContext()
		found, resp, err := config.Session.Client.ConnectionsApi.GetConnection(ctx, id)

		if err != nil {
			return fmt.Errorf("receive error when requesting Dummy Connection %s", id)
		}

		if resp.StatusCode != 200 {
			fmt.Errorf("Error getting Dummy Connection ID %s: %s", id, err)
		}

		*instance = found.(client.DummyConnection)

		return nil
	}
}

func testAccCheckDummyConnectionDestroy(s *terraform.State) error {

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

		if err != nil {
			return fmt.Errorf("should not get error for Dummy Connection with ID %s after delete: %s", id, err)
		}

		if resp.StatusCode != 404 {
			return fmt.Errorf("should not find Dummy Connection with ID %s existing after delete", id)
		}
	}

	return nil
}
