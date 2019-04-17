package pureport

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pureport/pureport-sdk-go/pureport/session"
	"github.com/pureport/pureport-sdk-go/pureport/swagger"
)

const testAccDataSourceAWSConnectionConfig_basic = `
data "pureport_cloud_regions" "main" {
	name_regex = "Oregon"
}

data "pureport_locations" "main" {
	name_regex = "^Sea*"
}

data "pureport_networks" "main" {
	account_id = "ac-8QVPmcPb_EhapbGHBMAo6Q"
	name_regex = "Bansh.*"
}

resource "pureport_aws_connection" "main" {
	name = "AwsDirectConnectTest"
	speed = "100"
	location_id = "${data.pureport_locations.main.locations.0.id}"
	network_id = "${data.pureport_networks.main.networks.0.id}"
	aws_region = "${data.pureport_cloud_regions.main.regions.0.identifier}"
	aws_account_id = "123456789012"
}
`

func TestAWSConnection_basic(t *testing.T) {

	resourceName := "resource.pureport_aws_connection.main"
	var instance swagger.AwsDirectConnectConnection

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAWSConnectionConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAWSConnection(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "id", instance.Id),
					resource.TestCheckResourceAttr(resourceName, "name", ""),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
		},
	})
}

func testAccCheckDataSourceAWSConnection(name string, instance *swagger.AwsDirectConnectConnection) resource.TestCheckFunc {
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
			return fmt.Errorf("receive error when requesting AWS Connection %s", id)
		}

		if resp.StatusCode != 200 {
			fmt.Errorf("Error getting AWS Connection ID %s: %s", id, err)
		}

		*instance = *found.(*swagger.AwsDirectConnectConnection)

		return nil
	}
}

func testAccCheckAWSConnectionDestroy(s *terraform.State) error {

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
			return fmt.Errorf("should not get error for AWS Connection with ID %s after delete: %s", id, err)
		}

		if resp.StatusCode != 404 {
			return fmt.Errorf("should not find AWS Connection with ID %s existing after delete", id)
		}
	}

	return nil
}
