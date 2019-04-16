package pureport

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testAccDataSourceCloudServicesConfig_foo = `
data "pureport_cloud_services" "foo" {
}
`

func TestCloudServices_basic(t *testing.T) {

	resourceName := "data.pureport_cloud_services.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudServicesConfig_foo,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceCloudServices(resourceName),
					resource.TestCheckResourceAttr(resourceName, "services.0.id", ""),
					resource.TestCheckResourceAttr(resourceName, "services.0.name", ""),
					resource.TestCheckResourceAttr(resourceName, "services.0.provider", ""),
					resource.TestCheckResourceAttr(resourceName, "services.0.service", ""),
					resource.TestCheckResourceAttr(resourceName, "services.0.ipv4_prefix_count", ""),
					resource.TestCheckResourceAttr(resourceName, "services.0.ipv6_prefix_count", ""),
					resource.TestCheckResourceAttr(resourceName, "services.0.cloud_region_id", ""),
					resource.TestCheckResourceAttr(resourceName, "services.#", "20"),
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
