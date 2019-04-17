package pureport

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testAccDataSourceLocationsConfig_empty = `
data "pureport_locations" "empty" {
}
`

const testAccDataSourceLocationsConfig_name_regex = `
data "pureport_locations" "name_regex" {
	name_regex = "^Sea*"
}
`

func TestLocations_empty(t *testing.T) {

	resourceName := "data.pureport_locations.empty"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLocationsConfig_empty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceLocations(resourceName),
					resource.TestCheckResourceAttr(resourceName, "locations.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.id", "us-ral"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.name", "Raleigh"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.links.#", "0"),
				),
			},
		},
	})
}

func TestLocations_name_regex(t *testing.T) {

	resourceName := "data.pureport_locations.name_regex"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLocationsConfig_name_regex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceLocations(resourceName),
					resource.TestCheckResourceAttr(resourceName, "locations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.id", "us-sea"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.name", "Seattle"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.links.#", "0"),
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
