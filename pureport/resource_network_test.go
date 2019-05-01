package pureport

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pureport/pureport-sdk-go/pureport/session"
	"github.com/pureport/pureport-sdk-go/pureport/swagger"
)

const testAccResourceNetworkConfig_basic = `
data "pureport_accounts" "main" {
	name_regex = "Terraform"
}

resource "pureport_network" "main" {
	name = "NetworkTest"
	description = "Network Terraform Test"
	account_id = "${data.pureport_accounts.main.accounts.0.id}"
}
`

func TestNetwork_basic(t *testing.T) {

	resourceName := "pureport_network.main"
	var instance swagger.Network

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
					resource.TestCheckResourceAttr(resourceName, "name", "NetworkTest"),
					resource.TestCheckResourceAttr(resourceName, "description", "Network Terraform Test"),
					resource.TestCheckResourceAttrPtr(resourceName, "account.0.id", &instance.Account.Id),
					resource.TestCheckResourceAttrSet(resourceName, "account.0.href"),
				),
			},
		},
	})
}

func testAccCheckResourceNetwork(name string, instance *swagger.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		sess, ok := testAccProvider.Meta().(*session.Session)
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

		ctx := sess.GetSessionContext()
		found, resp, err := sess.Client.NetworksApi.GetNetwork(ctx, id)

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
		_, resp, err := sess.Client.NetworksApi.GetNetwork(ctx, id)

		if err != nil && resp.StatusCode != 404 {
			return fmt.Errorf("should not get error for Network with ID %s after delete: %s", id, err)
		}
	}

	return nil
}
