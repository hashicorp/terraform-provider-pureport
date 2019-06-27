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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLocationsConfig_empty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceLocations(resourceName),

					resource.TestCheckResourceAttr(resourceName, "locations.#", "5"),

					resource.TestCheckResourceAttr(resourceName, "locations.0.id", "us-chi"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.href", "/locations/us-chi"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.name", "Chicago, IL"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.links.#", "4"),

					resource.TestCheckResourceAttr(resourceName, "locations.0.links.0.location_href", "/locations/us-wdc"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.links.0.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.0.links.1.location_href", "/locations/us-sjc"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.links.1.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.0.links.2.location_href", "/locations/us-sea"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.links.2.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.0.links.3.location_href", "/locations/us-dal"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.links.3.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.1.id", "us-dal"),
					resource.TestCheckResourceAttr(resourceName, "locations.1.href", "/locations/us-dal"),
					resource.TestCheckResourceAttr(resourceName, "locations.1.name", "Dallas, TX"),
					resource.TestCheckResourceAttr(resourceName, "locations.1.links.#", "4"),

					resource.TestCheckResourceAttr(resourceName, "locations.1.links.0.location_href", "/locations/us-wdc"),
					resource.TestCheckResourceAttr(resourceName, "locations.1.links.0.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.1.links.1.location_href", "/locations/us-sjc"),
					resource.TestCheckResourceAttr(resourceName, "locations.1.links.1.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.1.links.2.location_href", "/locations/us-sea"),
					resource.TestCheckResourceAttr(resourceName, "locations.1.links.2.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.1.links.3.location_href", "/locations/us-chi"),
					resource.TestCheckResourceAttr(resourceName, "locations.1.links.3.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.2.id", "us-sea"),
					resource.TestCheckResourceAttr(resourceName, "locations.2.href", "/locations/us-sea"),
					resource.TestCheckResourceAttr(resourceName, "locations.2.name", "Seattle, WA"),
					resource.TestCheckResourceAttr(resourceName, "locations.2.links.#", "4"),

					resource.TestCheckResourceAttr(resourceName, "locations.2.links.0.location_href", "/locations/us-wdc"),
					resource.TestCheckResourceAttr(resourceName, "locations.2.links.0.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.2.links.1.location_href", "/locations/us-sjc"),
					resource.TestCheckResourceAttr(resourceName, "locations.2.links.1.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.2.links.2.location_href", "/locations/us-chi"),
					resource.TestCheckResourceAttr(resourceName, "locations.2.links.2.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.2.links.3.location_href", "/locations/us-dal"),
					resource.TestCheckResourceAttr(resourceName, "locations.2.links.3.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.3.id", "us-sjc"),
					resource.TestCheckResourceAttr(resourceName, "locations.3.href", "/locations/us-sjc"),
					resource.TestCheckResourceAttr(resourceName, "locations.3.name", "Silicon Valley, CA"),
					resource.TestCheckResourceAttr(resourceName, "locations.3.links.#", "4"),

					resource.TestCheckResourceAttr(resourceName, "locations.3.links.0.location_href", "/locations/us-wdc"),
					resource.TestCheckResourceAttr(resourceName, "locations.3.links.0.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.3.links.1.location_href", "/locations/us-sea"),
					resource.TestCheckResourceAttr(resourceName, "locations.3.links.1.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.3.links.2.location_href", "/locations/us-chi"),
					resource.TestCheckResourceAttr(resourceName, "locations.3.links.2.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.3.links.3.location_href", "/locations/us-dal"),
					resource.TestCheckResourceAttr(resourceName, "locations.3.links.3.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.4.id", "us-wdc"),
					resource.TestCheckResourceAttr(resourceName, "locations.4.href", "/locations/us-wdc"),
					resource.TestCheckResourceAttr(resourceName, "locations.4.name", "Washington, DC"),
					resource.TestCheckResourceAttr(resourceName, "locations.4.links.#", "4"),

					resource.TestCheckResourceAttr(resourceName, "locations.4.links.0.location_href", "/locations/us-sjc"),
					resource.TestCheckResourceAttr(resourceName, "locations.4.links.0.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.4.links.1.location_href", "/locations/us-sea"),
					resource.TestCheckResourceAttr(resourceName, "locations.4.links.1.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.4.links.2.location_href", "/locations/us-chi"),
					resource.TestCheckResourceAttr(resourceName, "locations.4.links.2.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.4.links.3.location_href", "/locations/us-dal"),
					resource.TestCheckResourceAttr(resourceName, "locations.4.links.3.speed", "1000"),
				),
			},
		},
	})
}

func TestLocations_name_regex(t *testing.T) {

	resourceName := "data.pureport_locations.name_regex"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLocationsConfig_name_regex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceLocations(resourceName),

					resource.TestCheckResourceAttr(resourceName, "locations.#", "1"),

					resource.TestCheckResourceAttr(resourceName, "locations.0.id", "us-sea"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.href", "/locations/us-sea"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.name", "Seattle, WA"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.links.#", "4"),

					resource.TestCheckResourceAttr(resourceName, "locations.0.links.0.location_href", "/locations/us-wdc"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.links.0.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.0.links.1.location_href", "/locations/us-sjc"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.links.1.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.0.links.2.location_href", "/locations/us-chi"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.links.2.speed", "1000"),

					resource.TestCheckResourceAttr(resourceName, "locations.0.links.3.location_href", "/locations/us-dal"),
					resource.TestCheckResourceAttr(resourceName, "locations.0.links.3.speed", "1000"),
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
			return fmt.Errorf("Can't find Locations data source: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}
