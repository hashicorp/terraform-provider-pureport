package pureport

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/configuration"
)

func init() {
	resource.AddTestSweepers("pureport_azure_connection", &resource.Sweeper{
		Name: "pureport_azure_connection",
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

func testAccResourceAzureConnectionConfig_common() string {
	format := `
data "pureport_accounts" "main" {
  filter {
    name = "Name"
    values = ["Terraform .*"]
  }
}

data "pureport_locations" "main" {
  filter {
    name = "Name"
    values = ["Sea.*"]
  }
}

data "pureport_networks" "main" {
  account_href = "${data.pureport_accounts.main.accounts.0.href}"
  filter {
    name = "Name"
    values = ["Bansh.*"]
  }
}

data "azurerm_express_route_circuit" "main" {
  name                = "terraform-acc-express-route-%s"
  resource_group_name = "terraform-acceptance-tests"
}
`

	if testEnvironmentName == "Production" {
		return fmt.Sprintf(format, "prod")
	}

	return fmt.Sprintf(format, "dev1")
}

func testAccResourceAzureConnectionConfig() string {

	format := testAccResourceAzureConnectionConfig_common() + `
resource "pureport_azure_connection" "main" {
  name = "%s"
  description = "Some random description"
  speed = "100"
  high_availability = true

  location_href = "${data.pureport_locations.main.locations.0.href}"
  network_href = "${data.pureport_networks.main.networks.0.href}"

  service_key = "${data.azurerm_express_route_circuit.main.service_key}"

  tags = {
    Environment = "tf-test"
    Owner       = "ksk-azure"
    sweep       = "TRUE"
  }
}
`

	connection_name := acctest.RandomWithPrefix("AzureExpressRouteTest-")

	return fmt.Sprintf(format, connection_name)
}

func TestResourceAzureConnection_basic(t *testing.T) {

	resourceName := "pureport_azure_connection.main"
	var instance client.AzureExpressRouteConnection

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAzureConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAzureConnectionConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceAzureConnection(resourceName, &instance),
					resource.TestCheckResourceAttrPtr(resourceName, "id", &instance.Id),

					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "AzureExpressRouteTest"),
						resource.TestCheckResourceAttr(resourceName, "description", "Some random description"),
						resource.TestCheckResourceAttr(resourceName, "speed", "100"),
						resource.TestCheckResourceAttr(resourceName, "high_availability", "true"),
						resource.TestMatchResourceAttr(resourceName, "service_key", regexp.MustCompile("[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}")),

						resource.TestCheckResourceAttr(resourceName, "gateways.#", "2"),

						resource.TestCheckResourceAttr(resourceName, "gateways.0.availability_domain", "PRIMARY"),
						resource.TestCheckResourceAttr(resourceName, "gateways.0.name", "AZURE_EXPRESS_ROUTE"),
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
						resource.TestCheckResourceAttr(resourceName, "gateways.1.name", "AZURE_EXPRESS_ROUTE 2"),
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

						resource.TestCheckResourceAttr(resourceName, "tags.Environment", "tf-test"),
						resource.TestCheckResourceAttr(resourceName, "tags.Owner", "ksk-azure"),
					),
				),
			},
		},
	})
}

func testAccCheckResourceAzureConnection(name string, instance *client.AzureExpressRouteConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		config, ok := testAccProvider.Meta().(*configuration.Config)
		if !ok {
			return fmt.Errorf("Error getting Pureport client")
		}

		// Find the state object
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find Azure Connnection resource: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id := rs.Primary.ID

		ctx := config.Session.GetSessionContext()
		found, resp, err := config.Session.Client.ConnectionsApi.GetConnection(ctx, id)

		if err != nil {
			return fmt.Errorf("receive error when requesting Azure Connection %s", id)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("Error getting Azure Connection ID %s: %s", id, err)
		}

		*instance = found.(client.AzureExpressRouteConnection)

		return nil
	}
}

func testAccCheckAzureConnectionDestroy(s *terraform.State) error {

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
			return fmt.Errorf("should not get error for Azure Connection with ID %s after delete: %s", id, err)
		}

		if resp.StatusCode != 404 {
			return fmt.Errorf("should not find Azure Connection with ID %s existing after delete", id)
		}
	}

	return nil
}
