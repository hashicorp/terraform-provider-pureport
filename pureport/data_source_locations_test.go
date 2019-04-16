package pureport

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testAccDataSourceLocationsConfig_foo = `
data "pureport_locations" "foo" {
}
`

func TestLocations_basic(t *testing.T) {

	resourceName := "data.pureport_cloud_regions.any"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLocationsConfig_foo,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceLocations(resourceName),
					resource.TestCheckResourceAttr(resourceName, "locations.0.id", ""),
					resource.TestCheckResourceAttr(resourceName, "locations.0.name", ""),
					resource.TestCheckResourceAttr(resourceName, "locations.0.links.#", ""),
					resource.TestCheckResourceAttr(resourceName, "locations.0.links.0.location_id", ""),
					resource.TestCheckResourceAttr(resourceName, "locations.0.links.0.speed", ""),
					resource.TestCheckResourceAttr(resourceName, "locations.#", "20"),
				),
			},
		},
	})
}

func testAccCheckDataSourceLocations(name string) resource.TestCheckFunc {
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
