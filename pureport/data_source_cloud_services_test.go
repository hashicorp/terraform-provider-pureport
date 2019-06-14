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

const testAccDataSourceCloudServicesConfig_name_regex = `
data "pureport_cloud_services" "name_regex" {
	name_regex = ".*S3 us-west-2"
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

					resource.TestCheckResourceAttr(resourceName, "services.#", "4"),

					resource.TestCheckResourceAttr(resourceName, "services.0.id", "aws-s3-us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "services.0.href", "/cloudServices/aws-s3-us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "services.0.name", "AWS S3 us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "services.0.provider", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "services.0.service", "S3"),
					resource.TestCheckResourceAttr(resourceName, "services.0.ipv4_prefix_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "services.0.ipv6_prefix_count", "4"),
					resource.TestCheckResourceAttr(resourceName, "services.0.cloud_region_id", "aws-us-east-1"),
				),
			},
		},
	})
}

func TestCloudServices_name_regex(t *testing.T) {

	resourceName := "data.pureport_cloud_services.name_regex"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudServicesConfig_name_regex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceCloudServices(resourceName),

					resource.TestCheckResourceAttr(resourceName, "services.#", "1"),

					resource.TestCheckResourceAttr(resourceName, "services.0.id", "aws-s3-us-west-2"),
					resource.TestCheckResourceAttr(resourceName, "services.0.href", "/cloudServices/aws-s3-us-west-2"),
					resource.TestCheckResourceAttr(resourceName, "services.0.name", "AWS S3 us-west-2"),
					resource.TestCheckResourceAttr(resourceName, "services.0.provider", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "services.0.service", "S3"),
					resource.TestCheckResourceAttr(resourceName, "services.0.ipv4_prefix_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "services.0.ipv6_prefix_count", "4"),
					resource.TestCheckResourceAttr(resourceName, "services.0.cloud_region_id", "aws-us-west-2"),
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
