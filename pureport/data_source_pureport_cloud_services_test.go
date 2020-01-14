package pureport

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const testAccDataSourceCloudServicesConfig_empty = `
data "pureport_cloud_services" "empty" {
}
`

const testAccDataSourceCloudServicesConfig_name_filter = `
data "pureport_cloud_services" "name_filter" {
  filter {
    name = "Name"
    values = [".*S3 us-west-2"]
  }
}
`

func TestDataSourceCloudServicesDataSource_empty(t *testing.T) {

	resourceName := "data.pureport_cloud_services.empty"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudServicesConfig_empty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceCloudServices(resourceName),

					resource.TestCheckResourceAttr(resourceName, "services.#", "8"),

					resource.TestCheckResourceAttr(resourceName, "services.0.id", "aws-dynamodb-us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "services.0.href", "/cloudServices/aws-dynamodb-us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "services.0.name", "AWS Dynamodb us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "services.0.provider", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "services.0.service", "DYNAMODB"),
					resource.TestMatchResourceAttr(resourceName, "services.0.ipv4_prefix_count", regexp.MustCompile("[0-9]{1,2}")),
					resource.TestMatchResourceAttr(resourceName, "services.0.ipv6_prefix_count", regexp.MustCompile("[0-9]{1,2}")),
					resource.TestCheckResourceAttr(resourceName, "services.0.cloud_region_id", "aws-us-east-1"),

					resource.TestCheckResourceAttr(resourceName, "services.4.id", "aws-s3-us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "services.4.href", "/cloudServices/aws-s3-us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "services.4.name", "AWS S3 us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "services.4.provider", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "services.4.service", "S3"),
					resource.TestMatchResourceAttr(resourceName, "services.0.ipv4_prefix_count", regexp.MustCompile("[0-9]{1,2}")),
					resource.TestMatchResourceAttr(resourceName, "services.0.ipv6_prefix_count", regexp.MustCompile("[0-9]{1,2}")),
					resource.TestCheckResourceAttr(resourceName, "services.4.cloud_region_id", "aws-us-east-1"),
				),
			},
		},
	})
}

func TestDataSourceCloudServicesDataSource_name_filter(t *testing.T) {

	resourceName := "data.pureport_cloud_services.name_filter"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudServicesConfig_name_filter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceCloudServices(resourceName),

					resource.TestCheckResourceAttr(resourceName, "services.#", "1"),

					resource.TestCheckResourceAttr(resourceName, "services.0.id", "aws-s3-us-west-2"),
					resource.TestCheckResourceAttr(resourceName, "services.0.href", "/cloudServices/aws-s3-us-west-2"),
					resource.TestCheckResourceAttr(resourceName, "services.0.name", "AWS S3 us-west-2"),
					resource.TestCheckResourceAttr(resourceName, "services.0.provider", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "services.0.service", "S3"),
					resource.TestMatchResourceAttr(resourceName, "services.0.ipv4_prefix_count", regexp.MustCompile("[0-9]{1,2}")),
					resource.TestMatchResourceAttr(resourceName, "services.0.ipv6_prefix_count", regexp.MustCompile("[0-9]{1,2}")),
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
