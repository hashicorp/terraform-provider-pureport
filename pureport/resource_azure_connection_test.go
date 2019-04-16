package pureport

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pureport/pureport-sdk-go/pureport/session"
	"github.com/pureport/pureport-sdk-go/pureport/swagger"
)

const testAccDataSourceAzureConnectionConfig_basic = `
resource "pureport_azure_connection" "main" {
	name = "AzureExpressRouteTest"
	speed = "50"
	location_id = "us-ral"
	network_id = "network-RgwELBcU0ATnC5JezEAsSg"
	service_key = "8d892e3a-caae-48ac-9b71-4760de0b1d2c"
}
`

func TestAzureConnection_basic(t *testing.T) {

	resourceName := "resource.pureport_aws_connection.main"
	var instance swagger.AzureExpressRouteConnection

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAzureConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureConnectionConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAzureConnection(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "id", instance.Id),
					resource.TestCheckResourceAttr(resourceName, "name", ""),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
		},
	})
}

func testAccCheckDataSourceAzureConnection(name string, instance *swagger.AzureExpressRouteConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		sess, ok := testAccProvider.Meta().(*session.Session)
		if !ok {
			return fmt.Errorf("Error getting Pureport client")
		}

		// Find the state object
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find AWS Connnection resource: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id := rs.Primary.ID

		ctx := sess.GetSessionContext()
		found, resp, err := sess.Client.ConnectionsApi.Get11(ctx, id)

		if err != nil {
			return fmt.Errorf("receive error when requesting Azure Connection %s", id)
		}

		if resp.StatusCode != 200 {
			fmt.Errorf("Error getting Azure Connection ID %s: %s", id, err)
		}

		*instance = *found.(*swagger.AzureExpressRouteConnection)

		return nil
	}
}

func testAccCheckAzureConnectionDestroy(s *terraform.State) error {

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
		_, resp, err := sess.Client.ConnectionsApi.Get11(ctx, id)

		if err != nil {
			return fmt.Errorf("should not get error for Azure Connection with ID %s after delete: %s", id, err)
		}

		if resp.StatusCode != 404 {
			return fmt.Errorf("should not find Azure Connection with ID %s existing after delete", id)
		}
	}

	return nil
}
