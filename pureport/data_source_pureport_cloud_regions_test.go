package pureport

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const testAccDataSourceCloudRegionsConfig_empty = `
data "pureport_cloud_regions" "empty" {
}
`

const testAccDataSourceCloudRegionsConfig_name_filter = `
data "pureport_cloud_regions" "name_filter" {
  filter {
    name = "DisplayName"
    values = ["US East.*"]
  }
}
`

func TestDataSourceCloudRegionsDataSource_empty(t *testing.T) {

	resourceName := "data.pureport_cloud_regions.empty"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudRegionsConfig_empty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceCloudRegions(resourceName),

					resource.TestMatchResourceAttr(resourceName, "regions.#", regexp.MustCompile("[0-9]{1,2}")),

					resource.TestCheckResourceAttr(resourceName, "regions.0.id", "aws-ap-northeast-1"),
					resource.TestCheckResourceAttr(resourceName, "regions.0.name", "Asia Pacific (Tokyo)"),
					resource.TestCheckResourceAttr(resourceName, "regions.0.provider", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "regions.0.identifier", "ap-northeast-1"),
				),
			},
		},
	})
}

func TestDataSourceCloudRegionsDataSource_name_filter(t *testing.T) {

	resourceName := "data.pureport_cloud_regions.name_filter"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudRegionsConfig_name_filter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceCloudRegions(resourceName),

					resource.TestCheckResourceAttr(resourceName, "regions.#", "2"),

					resource.TestCheckResourceAttr(resourceName, "regions.0.id", "aws-us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "regions.0.name", "US East (N. Virginia)"),
					resource.TestCheckResourceAttr(resourceName, "regions.0.provider", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "regions.0.identifier", "us-east-1"),

					resource.TestCheckResourceAttr(resourceName, "regions.1.id", "aws-us-east-2"),
					resource.TestCheckResourceAttr(resourceName, "regions.1.name", "US East (Ohio)"),
					resource.TestCheckResourceAttr(resourceName, "regions.1.provider", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "regions.1.identifier", "us-east-2"),
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
