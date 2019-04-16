package pureport

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testAccDataSourceCloudRegionsConfig_foo = `
data "pureport_cloud_regions" "foo" {
}
`

func TestCloudRegions_basic(t *testing.T) {

	resourceName := "data.pureport_cloud_regions.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudRegionsConfig_foo,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceCloudRegions(resourceName),
					resource.TestCheckResourceAttr(resourceName, "regions.0.id", ""),
					resource.TestCheckResourceAttr(resourceName, "regions.0.name", ""),
					resource.TestCheckResourceAttr(resourceName, "regions.0.provider", ""),
					resource.TestCheckResourceAttr(resourceName, "regions.0.service", ""),
					resource.TestCheckResourceAttr(resourceName, "regions.0.ipv4_prefix_count", ""),
					resource.TestCheckResourceAttr(resourceName, "regions.0.ipv6_prefix_count", ""),
					resource.TestCheckResourceAttr(resourceName, "regions.0.cloud_region_id", ""),
					resource.TestCheckResourceAttr(resourceName, "regions.#", "20"),
				),
			},
		},
	})
}

func testAccCheckDataSourceCloudRegions(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		// Find the state object
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find Cloud Region data source: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}
