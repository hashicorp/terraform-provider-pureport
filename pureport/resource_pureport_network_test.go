package pureport

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/pureport/terraform-provider-pureport/pureport/configuration"
)

func init() {
	resource.AddTestSweepers("pureport_network", &resource.Sweeper{
		Name: "pureport_network",
		F: func(region string) error {
			c, err := sharedClientForRegion(region)
			if err != nil {
				return fmt.Errorf("Error getting client: %s", err)
			}

			config := c.(*configuration.Config)
			networks, err := config.GetAccNetworks()
			if err != nil {
				return fmt.Errorf("Error getting networks %s", err)
			}

			if err = config.SweepNetworks(networks); err != nil {
				return fmt.Errorf("Error occurred sweeping networks")
			}

			return nil
		},
	})
}

const testAccResourceNetworkConfig_common = `
data "pureport_accounts" "main" {
  filter {
    name = "Name"
    values = ["Terraform"]
  }
}
`

const testAccResourceNetworkConfig_basic = testAccResourceNetworkConfig_common + `
resource "pureport_network" "main" {
  name = "NetworkTest"
  description = "Network Terraform Test"
  account_href = "${data.pureport_accounts.main.accounts.0.href}"

  tags = {
    Environment = "tf-test"
    Owner       = "the-rockit"
    sweep       = "TRUE"
  }
}
`

func TestResourceNetwork_basic(t *testing.T) {

	resourceName := "pureport_network.main"
	var instance client.Network

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNetworkConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceNetwork(resourceName, &instance),
					resource.TestCheckResourceAttrPtr(resourceName, "id", &instance.Id),
					resource.TestCheckResourceAttrPtr(resourceName, "href", &instance.Href),
					resource.TestCheckResourceAttr(resourceName, "name", "NetworkTest"),
					resource.TestCheckResourceAttr(resourceName, "description", "Network Terraform Test"),
					resource.TestMatchResourceAttr(resourceName, "account_href", regexp.MustCompile("/accounts/ac-.{16}")),

					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "tf-test"),
					resource.TestCheckResourceAttr(resourceName, "tags.Owner", "the-rockit"),
				),
			},
		},
	})
}

func testAccCheckResourceNetwork(name string, instance *client.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		config, ok := testAccProvider.Meta().(*configuration.Config)
		if !ok {
			return fmt.Errorf("Error getting Pureport client")
		}

		// Find the state object
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find Network resource: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id := rs.Primary.ID

		ctx := config.Session.GetSessionContext()
		found, resp, err := config.Session.Client.NetworksApi.GetNetwork(ctx, id)

		if err != nil {
			return fmt.Errorf("receive error when requesting Network ID %s", id)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("Error getting Network ID %s: %s", id, err)
		}

		*instance = found

		return nil
	}
}

func testAccCheckNetworkDestroy(s *terraform.State) error {

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
		_, resp, err := config.Session.Client.NetworksApi.GetNetwork(ctx, id)

		if err != nil && resp.StatusCode != 404 {
			return fmt.Errorf("should not get error for Network with ID %s after delete: %s", id, err)
		}
	}

	return nil
}
