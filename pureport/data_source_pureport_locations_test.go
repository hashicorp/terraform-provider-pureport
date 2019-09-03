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

const testAccDataSourceLocationsConfig_name_filter = `
data "pureport_locations" "name_filter" {
  filter {
    name = "Name"
    values = ["^Sea*"]
  }
}
`

func TestDataSourceLocations_empty(t *testing.T) {

	resourceName := "data.pureport_locations.empty"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLocationsConfig_empty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceLocations(resourceName),
					testAccCheckDataSourceLocationAll(resourceName),
				),
			},
		},
	})
}

func TestDataSourceLocations_name_filter(t *testing.T) {

	resourceName := "data.pureport_locations.name_filter"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLocationsConfig_name_filter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceLocations(resourceName),

					resource.TestCheckResourceAttr(resourceName, "locations.#", "1"),
					testAccCheckDataSourceLocationSeattle(resourceName, "locations.0"),
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

func testAccCheckDataSourceLocationAll(resourceName string) resource.TestCheckFunc {
	if testEnvironmentName == "Production" {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceName, "locations.#", "5"),
			testAccCheckDataSourceLocationChicago(resourceName, "locations.0"),
			testAccCheckDataSourceLocationDallas(resourceName, "locations.1"),
			testAccCheckDataSourceLocationSeattle(resourceName, "locations.2"),
			testAccCheckDataSourceLocationSanJose(resourceName, "locations.3"),
			testAccCheckDataSourceLocationWashington(resourceName, "locations.4"),
		)
	}

	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "locations.#", "3"),
		testAccCheckDataSourceLocationRaleigh(resourceName, "locations.0"),
		testAccCheckDataSourceLocationSeattle(resourceName, "locations.1"),
		testAccCheckDataSourceLocationVirtualPod(resourceName, "locations.2"),
	)
}

func testAccCheckDataSourceLocationRaleigh(resourceName, location string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, location+".id", "us-ral"),
		resource.TestCheckResourceAttr(resourceName, location+".href", "/locations/us-ral"),
		resource.TestCheckResourceAttr(resourceName, location+".name", "Raleigh, NC"),
		resource.TestCheckResourceAttr(resourceName, location+".links.#", "0"),
	)
}

func testAccCheckDataSourceLocationChicago(resourceName, location string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, location+".id", "us-chi"),
		resource.TestCheckResourceAttr(resourceName, location+".href", "/locations/us-chi"),
		resource.TestCheckResourceAttr(resourceName, location+".name", "Chicago, IL"),
		resource.TestCheckResourceAttr(resourceName, location+".links.#", "4"),

		resource.TestCheckResourceAttr(resourceName, location+".links.0.location_href", "/locations/us-wdc"),
		resource.TestCheckResourceAttr(resourceName, location+".links.0.speed", "1000"),

		resource.TestCheckResourceAttr(resourceName, location+".links.1.location_href", "/locations/us-sjc"),
		resource.TestCheckResourceAttr(resourceName, location+".links.1.speed", "1000"),

		resource.TestCheckResourceAttr(resourceName, location+".links.2.location_href", "/locations/us-sea"),
		resource.TestCheckResourceAttr(resourceName, location+".links.2.speed", "1000"),

		resource.TestCheckResourceAttr(resourceName, location+".links.3.location_href", "/locations/us-dal"),
		resource.TestCheckResourceAttr(resourceName, location+".links.3.speed", "1000"),
	)
}

func testAccCheckDataSourceLocationDallas(resourceName, location string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, location+".id", "us-dal"),
		resource.TestCheckResourceAttr(resourceName, location+".href", "/locations/us-dal"),
		resource.TestCheckResourceAttr(resourceName, location+".name", "Dallas, TX"),
		resource.TestCheckResourceAttr(resourceName, location+".links.#", "4"),

		resource.TestCheckResourceAttr(resourceName, location+".links.0.location_href", "/locations/us-wdc"),
		resource.TestCheckResourceAttr(resourceName, location+".links.0.speed", "1000"),

		resource.TestCheckResourceAttr(resourceName, location+".links.1.location_href", "/locations/us-sjc"),
		resource.TestCheckResourceAttr(resourceName, location+".links.1.speed", "1000"),

		resource.TestCheckResourceAttr(resourceName, location+".links.2.location_href", "/locations/us-sea"),
		resource.TestCheckResourceAttr(resourceName, location+".links.2.speed", "1000"),

		resource.TestCheckResourceAttr(resourceName, location+".links.3.location_href", "/locations/us-chi"),
		resource.TestCheckResourceAttr(resourceName, location+".links.3.speed", "1000"),
	)
}

func testAccCheckDataSourceLocationSeattle(resourceName, location string) resource.TestCheckFunc {
	if testEnvironmentName == "Production" {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resourceName, location+".id", "us-sea"),
			resource.TestCheckResourceAttr(resourceName, location+".href", "/locations/us-sea"),
			resource.TestCheckResourceAttr(resourceName, location+".name", "Seattle, WA"),
			resource.TestCheckResourceAttr(resourceName, location+".links.#", "4"),

			resource.TestCheckResourceAttr(resourceName, location+".links.0.location_href", "/locations/us-wdc"),
			resource.TestCheckResourceAttr(resourceName, location+".links.0.speed", "1000"),

			resource.TestCheckResourceAttr(resourceName, location+".links.1.location_href", "/locations/us-sjc"),
			resource.TestCheckResourceAttr(resourceName, location+".links.1.speed", "1000"),

			resource.TestCheckResourceAttr(resourceName, location+".links.2.location_href", "/locations/us-chi"),
			resource.TestCheckResourceAttr(resourceName, location+".links.2.speed", "1000"),

			resource.TestCheckResourceAttr(resourceName, location+".links.3.location_href", "/locations/us-dal"),
			resource.TestCheckResourceAttr(resourceName, location+".links.3.speed", "1000"),
		)
	}

	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, location+".id", "us-sea"),
		resource.TestCheckResourceAttr(resourceName, location+".href", "/locations/us-sea"),
		resource.TestCheckResourceAttr(resourceName, location+".name", "Seattle, WA"),
		resource.TestCheckResourceAttr(resourceName, location+".links.#", "0"),
	)
}

func testAccCheckDataSourceLocationSanJose(resourceName, location string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, location+".id", "us-sjc"),
		resource.TestCheckResourceAttr(resourceName, location+".href", "/locations/us-sjc"),
		resource.TestCheckResourceAttr(resourceName, location+".name", "Silicon Valley, CA"),
		resource.TestCheckResourceAttr(resourceName, location+".links.#", "4"),

		resource.TestCheckResourceAttr(resourceName, location+".links.0.location_href", "/locations/us-wdc"),
		resource.TestCheckResourceAttr(resourceName, location+".links.0.speed", "1000"),

		resource.TestCheckResourceAttr(resourceName, location+".links.1.location_href", "/locations/us-sea"),
		resource.TestCheckResourceAttr(resourceName, location+".links.1.speed", "1000"),

		resource.TestCheckResourceAttr(resourceName, location+".links.2.location_href", "/locations/us-chi"),
		resource.TestCheckResourceAttr(resourceName, location+".links.2.speed", "1000"),

		resource.TestCheckResourceAttr(resourceName, location+".links.3.location_href", "/locations/us-dal"),
		resource.TestCheckResourceAttr(resourceName, location+".links.3.speed", "1000"),
	)
}

func testAccCheckDataSourceLocationWashington(resourceName, location string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, location+".id", "us-wdc"),
		resource.TestCheckResourceAttr(resourceName, location+".href", "/locations/us-wdc"),
		resource.TestCheckResourceAttr(resourceName, location+".name", "Washington, DC"),
		resource.TestCheckResourceAttr(resourceName, location+".links.#", "4"),

		resource.TestCheckResourceAttr(resourceName, location+".links.0.location_href", "/locations/us-sjc"),
		resource.TestCheckResourceAttr(resourceName, location+".links.0.speed", "1000"),

		resource.TestCheckResourceAttr(resourceName, location+".links.1.location_href", "/locations/us-sea"),
		resource.TestCheckResourceAttr(resourceName, location+".links.1.speed", "1000"),

		resource.TestCheckResourceAttr(resourceName, location+".links.2.location_href", "/locations/us-chi"),
		resource.TestCheckResourceAttr(resourceName, location+".links.2.speed", "1000"),

		resource.TestCheckResourceAttr(resourceName, location+".links.3.location_href", "/locations/us-dal"),
		resource.TestCheckResourceAttr(resourceName, location+".links.3.speed", "1000"),
	)
}

func testAccCheckDataSourceLocationVirtualPod(resourceName, location string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, location+".id", "virt-sea"),
		resource.TestCheckResourceAttr(resourceName, location+".href", "/locations/virt-sea"),
		resource.TestCheckResourceAttr(resourceName, location+".name", "Tacoma, WA"),
		resource.TestCheckResourceAttr(resourceName, location+".links.#", "0"),
	)
}
