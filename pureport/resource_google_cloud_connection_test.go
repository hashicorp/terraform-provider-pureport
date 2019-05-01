package pureport

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pureport/pureport-sdk-go/pureport/session"
	"github.com/pureport/pureport-sdk-go/pureport/swagger"
)

const testAccResourceGoogleCloudConnectionConfig_basic = `
data "pureport_accounts" "main" {
	name_regex = "Terraform"
}

data "pureport_cloud_regions" "main" {
	name_regex = "Oregon"
}

data "pureport_locations" "main" {
	name_regex = "^Ral.*"
}

data "pureport_networks" "main" {
	account_id = "${data.pureport_accounts.main.accounts.0.id}"
	name_regex = "Bansh.*"
}

resource "pureport_google_cloud_connection" "main" {
	name = "GoogleCloudTest"
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

	primary_pairing_key = "ff040d5e-476a-43eb-a5c0-ec069f0d1ca1/us-west2/1"
	secondary_pairing_key = "9826da45-2c3b-44e7-8d2f-ae856ce8a90d/us-west2/2"
}
`

func TestGoogleCloudConnection_basic(t *testing.T) {

	resourceName := "pureport_google_cloud_connection.main"
	var instance swagger.GoogleCloudInterconnectConnection

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleCloudConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGoogleCloudConnectionConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGoogleCloudConnection(resourceName, &instance),
					resource.TestCheckResourceAttrPtr(resourceName, "id", &instance.Id),
					resource.TestCheckResourceAttr(resourceName, "name", "GoogleCloudTest"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "speed", "100"),
					resource.TestCheckResourceAttr(resourceName, "high_availability", "true"),
					resource.TestCheckResourceAttr(resourceName, "high_availability", "true"),
					resource.TestCheckResourceAttr(resourceName, "primary_pairing_key", "ff040d5e-476a-43eb-a5c0-ec069f0d1ca1/us-west2/1"),
					resource.TestCheckResourceAttr(resourceName, "secondary_pairing_key", "9826da45-2c3b-44e7-8d2f-ae856ce8a90d/us-west2/2"),
				),
			},
		},
	})
}

func testAccCheckResourceGoogleCloudConnection(name string, instance *swagger.GoogleCloudInterconnectConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		sess, ok := testAccProvider.Meta().(*session.Session)
		if !ok {
			return fmt.Errorf("Error getting Pureport client")
		}

		// Find the state object
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find Dummy Connnection resource: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id := rs.Primary.ID

		ctx := sess.GetSessionContext()
		found, resp, err := sess.Client.ConnectionsApi.GetConnection(ctx, id)

		if err != nil {
			return fmt.Errorf("receive error when requesting Google Cloud Connection %s", id)
		}

		if resp.StatusCode != 200 {
			fmt.Errorf("Error getting Google Cloud Connection ID %s: %s", id, err)
		}

		*instance = found.(swagger.GoogleCloudInterconnectConnection)

		return nil
	}
}

func testAccCheckGoogleCloudConnectionDestroy(s *terraform.State) error {

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

		if err != nil {
			return fmt.Errorf("should not get error for Google Cloud Connection with ID %s after delete: %s", id, err)
		}

		if resp.StatusCode != 404 {
			return fmt.Errorf("should not find Google Cloud Connection with ID %s existing after delete", id)
		}
	}

	return nil
}
