package pureport

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pureport/pureport-sdk-go/pureport/client"
)

const testAccResourceNetworkConfig_common = `
data "pureport_accounts" "main" {
  name_regex = "Terraform"
}
`

const testAccResourceNetworkConfig_basic = testAccResourceNetworkConfig_common + `
resource "pureport_network" "main" {
  name = "NetworkTest"
  description = "Network Terraform Test"
  account_href = "${data.pureport_accounts.main.accounts.0.href}"
}
`

func TestNetwork_basic(t *testing.T) {

	resourceName := "pureport_network.main"
	var instance client.Network

	resource.Test(t, resource.TestCase{
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
					resource.TestCheckResourceAttr(resourceName, "account_href", "/accounts/ac-8QVPmcPb_EhapbGHBMAo6Q"),
				),
			},
		},
	})
}

func testAccCheckResourceNetwork(name string, instance *client.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		config, ok := testAccProvider.Meta().(*Config)
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
			fmt.Errorf("Error getting Network ID %s: %s", id, err)
		}

		*instance = found

		return nil
	}
}

func testAccCheckNetworkDestroy(s *terraform.State) error {

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
		_, resp, err := config.Session.Client.NetworksApi.GetNetwork(ctx, id)

		if err != nil && resp.StatusCode != 404 {
			return fmt.Errorf("should not get error for Network with ID %s after delete: %s", id, err)
		}
	}

	return nil
}
