package pureport

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testAccDataSourceCloudServicesConfig_empty = `
data "pureport_cloud_services" "empty" {
}
`

func TestCloudServices_empty(t *testing.T) {

	resourceName := "data.pureport_cloud_services.empty"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudServicesConfig_empty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceCloudServices(resourceName),
					resource.TestCheckResourceAttr(resourceName, "services.0.id", "aws-cloud9-us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "services.0.name", "AWS Cloud9 us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "services.0.provider", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "services.0.service", "CLOUD9"),
					resource.TestCheckResourceAttr(resourceName, "services.0.ipv4_prefix_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "services.0.ipv6_prefix_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "services.0.cloud_region_id", "aws-us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "services.#", "16"),
				),
			},
		},
	})
}

func testAccCheckDataSourceCloudServices(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		// Find the state object
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find Cloud Services data source: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}
