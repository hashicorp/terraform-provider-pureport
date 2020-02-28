package pureport

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/configuration"
)

func init() {
	resource.AddTestSweepers("pureport_google_cloud_connection", &resource.Sweeper{
		Name: "pureport_google_cloud_connection",
		F: func(region string) error {
			c, err := sharedClientForRegion(region)
			if err != nil {
				return fmt.Errorf("Error getting client: %s", err)
			}

			config := c.(*configuration.Config)
			connections, err := config.GetAccConnections()
			if err != nil {
				return fmt.Errorf("Error getting connections %s", err)
			}

			if err = config.SweepConnections(connections); err != nil {
				return fmt.Errorf("Error occurred sweeping connections")
			}

			return nil
		},
	})
}

const testAccResourceGoogleCloudConnectionConfig_common = `
data "pureport_accounts" "main" {
  filter {
    name = "Name"
    values = ["Terraform"]
  }
}

data "pureport_locations" "main" {
  filter {
    name = "Name"
    values = ["^Sea.*"]
  }
}

data "pureport_networks" "main" {
  account_href = data.pureport_accounts.main.accounts.0.href
  filter {
    name = "Name"
    values = ["Bansh.*"]
  }
}
`

func testAccResourceGoogleCloudConnectionConfig_basic() string {
	format := testAccResourceGoogleCloudConnectionConfig_common + `
data "google_compute_network" "default" {
  name = "terraform-acc-network-%s"
}

resource "google_compute_router" "main" {
  name    = "terraform-acc-%s-${count.index + 1}"
  network = data.google_compute_network.default.name

  bgp {
    asn = "16550"
  }

  count = 2
}

resource "google_compute_interconnect_attachment" "main" {
  name   = "terraform-acc-%s-${count.index + 1}"
  router = element(google_compute_router.main.*.self_link, count.index)
  type   = "PARTNER"
  edge_availability_domain = "AVAILABILITY_DOMAIN_${count.index + 1}"

  lifecycle {
    ignore_changes = ["vlan_tag8021q"]
  }

  count = 2
}

resource "pureport_google_cloud_connection" "main" {
  name = "%s"
  speed = "50"

  location_href = data.pureport_locations.main.locations.0.href
  network_href = data.pureport_networks.main.networks.0.href

  primary_pairing_key = google_compute_interconnect_attachment.main.0.pairing_key

  tags = {
    Environment = "tf-test"
    Owner       = "ksk-google"
    sweep       = "TRUE"
  }
}
`

	router_name := acctest.RandomWithPrefix("router")
	interconnect_name := acctest.RandomWithPrefix("interconnect")
	connection_name := acctest.RandomWithPrefix("GoogleCloudTest")
	environment_name := "dev1"

	if testEnvironmentName == "Production" {
		environment_name = "prod"
	}

	return fmt.Sprintf(format, environment_name, router_name, interconnect_name, connection_name)
}

func TestResourceGoogleCloudConnection_basic(t *testing.T) {

	resourceName := "pureport_google_cloud_connection.main"
	var instance client.GoogleCloudInterconnectConnection

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleCloudConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGoogleCloudConnectionConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGoogleCloudConnection(resourceName, &instance),
					resource.TestCheckResourceAttrPtr(resourceName, "id", &instance.Id),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile("^GoogleCloudTest-.*")),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "speed", "50"),
					resource.TestCheckResourceAttr(resourceName, "high_availability", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_pairing_key"),
					resource.TestCheckResourceAttr(resourceName, "secondary_pairing_key", ""),

					resource.TestCheckResourceAttr(resourceName, "gateways.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.availability_domain", "PRIMARY"),
					resource.TestMatchResourceAttr(resourceName, "gateways.0.name", regexp.MustCompile("^GoogleCloudTest-.* - Primary")),
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

					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "tf-test"),
					resource.TestCheckResourceAttr(resourceName, "tags.Owner", "ksk-google"),
				),
			},
		},
	})
}

func testAccCheckResourceGoogleCloudConnection(name string, instance *client.GoogleCloudInterconnectConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		config, ok := testAccProvider.Meta().(*configuration.Config)
		if !ok {
			return fmt.Errorf("Error getting Pureport client")
		}

		// Find the state object
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find Google Cloud Connection resource: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id := rs.Primary.ID

		ctx := config.Session.GetSessionContext()
		found, resp, err := config.Session.Client.ConnectionsApi.GetConnection(ctx, id)

		if err != nil {
			return fmt.Errorf("receive error when requesting Google Cloud Connection %s", id)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("Error getting Google Cloud Connection ID %s: %s", id, err)
		}

		*instance = found.(client.GoogleCloudInterconnectConnection)

		return nil
	}
}

func testAccCheckGoogleCloudConnectionDestroy(s *terraform.State) error {

	config, ok := testAccProvider.Meta().(*configuration.Config)
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
			return fmt.Errorf("should not get error for Google Cloud Connection with ID %s after delete: %s", id, err)
		}

		if resp.StatusCode != 404 {
			return fmt.Errorf("should not find Google Cloud Connection with ID %s existing after delete", id)
		}
	}

	return nil
}
