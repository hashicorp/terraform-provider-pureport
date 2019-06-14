package pureport

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pureport/pureport-sdk-go/pureport/client"
)

const testAccResourceAWSConnectionConfig_common = `
data "pureport_accounts" "main" {
  name_regex = "Terraform"
}

data "pureport_cloud_regions" "main" {
  name_regex = "Oregon"
}

data "pureport_locations" "main" {
  name_regex = "^Sea*"
}

data "pureport_networks" "main" {
  account_href = "${data.pureport_accounts.main.accounts.0.href}"
  name_regex = "Bansh.*"
}

data "aws_caller_identity" "current" {}
`

const testAccResourceAWSConnectionConfig_basic = testAccResourceAWSConnectionConfig_common + `
resource "pureport_aws_connection" "main" {
  name = "AwsDirectConnectTest"
  speed = "100"
  high_availability = true

  location_href = "${data.pureport_locations.main.locations.0.href}"
  network_href = "${data.pureport_networks.main.networks.0.href}"

  aws_region = "${data.pureport_cloud_regions.main.regions.0.identifier}"
  aws_account_id = "${data.aws_caller_identity.current.account_id}"
}
`

const testAccResourceAWSConnectionConfig_basic_update_speed = testAccResourceAWSConnectionConfig_common + `
resource "pureport_aws_connection" "main" {
  name = "AwsDirectConnectTest"
  speed = "200"
  high_availability = true

  location_href = "${data.pureport_locations.main.locations.0.href}"
  network_href = "${data.pureport_networks.main.networks.0.href}"

  aws_region = "${data.pureport_cloud_regions.main.regions.0.identifier}"
  aws_account_id = "${data.aws_caller_identity.current.account_id}"
}
`

const testAccResourceAWSConnectionConfig_basic_update_no_respawn = testAccResourceAWSConnectionConfig_common + `
resource "pureport_aws_connection" "main" {
  name = "Aws DirectConnect Test"
  description = "AWS Basic Test"
  speed = "100"
  high_availability = true

  location_href = "${data.pureport_locations.main.locations.0.href}"
  network_href = "${data.pureport_networks.main.networks.0.href}"

  aws_region = "${data.pureport_cloud_regions.main.regions.0.identifier}"
  aws_account_id = "${data.aws_caller_identity.current.account_id}"
}
`

const testAccResourceAWSConnectionConfig_basic_update_respawn = testAccResourceAWSConnectionConfig_common + `
resource "pureport_aws_connection" "main" {
  name = "AwsDirectConnectTest"
  speed = "100"
  high_availability = true

  location_href = "${data.pureport_locations.main.locations.0.href}"
  network_href = "${data.pureport_networks.main.networks.0.href}"

  aws_region = "${data.pureport_cloud_regions.main.regions.0.identifier}"
  aws_account_id = "${data.aws_caller_identity.current.account_id}"
}
`

const testAccResourceAWSConnectionConfig_cloudServices = testAccResourceAWSConnectionConfig_common + `
data "pureport_cloud_services" "s3" {
  name_regex = ".*S3"
}

resource "pureport_aws_connection" "main" {
  name = "AwsDirectConnectCloudServicesTest"
  speed = "100"
  high_availability = true

  location_href = "${data.pureport_locations.main.locations.0.href}"
  network_href = "${data.pureport_networks.main.networks.0.href}"

  cloud_service_hrefs = data.pureport_cloud_services.s3.services.*.href
  peering_type = "PUBLIC"

  aws_region = "${data.pureport_cloud_regions.main.regions.0.identifier}"
  aws_account_id = "${data.aws_caller_identity.current.account_id}"
}
`

const testAccResourceAWSConnectionConfig_nat_mapping = testAccResourceAWSConnectionConfig_common + `
resource "pureport_aws_connection" "main" {
  name = "AwsDirectConnectNatMappingTest"
  speed = "100"
  high_availability = true

  location_href = "${data.pureport_locations.main.locations.0.href}"
  network_href = "${data.pureport_networks.main.networks.0.href}"

  aws_region = "${data.pureport_cloud_regions.main.regions.0.identifier}"
  aws_account_id = "${data.aws_caller_identity.current.account_id}"

  nat_config {
    enabled = true

    mappings {
      native_cidr = "192.168.0.0/24"
    }

    mappings {
      native_cidr = "192.200.0.0/16"
    }
  }
}
`

func TestAWSConnection_basic(t *testing.T) {

	resourceName := "pureport_aws_connection.main"
	var instance client.AwsDirectConnectConnection
	var respawn_instance client.AwsDirectConnectConnection

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAWSConnectionConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceAWSConnection(resourceName, &instance),

					resource.TestCheckResourceAttrPtr(resourceName, "id", &instance.Id),
					resource.TestCheckResourceAttr(resourceName, "name", "AwsDirectConnectTest"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "speed", "100"),
					resource.TestCheckResourceAttr(resourceName, "high_availability", "true"),
					resource.TestCheckResourceAttr(resourceName, "location_href", "/locations/us-sea"),
					resource.TestMatchResourceAttr(resourceName, "network_href", regexp.MustCompile("/networks/network-.{16}")),

					resource.TestCheckResourceAttr(resourceName, "gateways.#", "2"),

					resource.TestCheckResourceAttr(resourceName, "gateways.0.availability_domain", "PRIMARY"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.name", "AWS_DIRECT_CONNECT"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.link_state", "PENDING"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.customer_asn", "64512"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.customer_ip", "169.254.1.2/30"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.pureport_asn", "394351"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.pureport_ip", "169.254.1.1/30"),
					resource.TestCheckResourceAttrSet(resourceName, "gateways.0.bgp_password"),
					resource.TestMatchResourceAttr(resourceName, "gateways.0.peering_subnet", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.public_nat_ip", ""),
					resource.TestCheckResourceAttrSet(resourceName, "gateways.0.vlan"),
					resource.TestCheckResourceAttrSet(resourceName, "gateways.0.remote_id"),

					resource.TestCheckResourceAttr(resourceName, "gateways.1.availability_domain", "SECONDARY"),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.name", "AWS_DIRECT_CONNECT 2"),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.description", ""),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.link_state", "PENDING"),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.customer_asn", "64512"),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.customer_ip", "169.254.2.2/30"),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.pureport_asn", "394351"),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.pureport_ip", "169.254.2.1/30"),
					resource.TestCheckResourceAttrSet(resourceName, "gateways.1.bgp_password"),
					resource.TestMatchResourceAttr(resourceName, "gateways.1.peering_subnet", regexp.MustCompile("169.254.[0-9]{1,3}.[0-9]{1,3}")),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.public_nat_ip", ""),
					resource.TestCheckResourceAttrSet(resourceName, "gateways.1.vlan"),
					resource.TestCheckResourceAttrSet(resourceName, "gateways.1.remote_id"),
				),
			},
			{
				Config: testAccResourceAWSConnectionConfig_basic_update_no_respawn,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceName, "id", &instance.Id),
					resource.TestCheckResourceAttr(resourceName, "name", "Aws DirectConnect Test"),
					resource.TestCheckResourceAttr(resourceName, "description", "AWS Basic Test"),
				),
			},
			{
				Config: testAccResourceAWSConnectionConfig_basic_update_respawn,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceAWSConnection(resourceName, &respawn_instance),
					resource.TestCheckResourceAttrPtr(resourceName, "id", &respawn_instance.Id),
					TestCheckResourceConnectionIdChanged(&instance.Id, &respawn_instance.Id),
					resource.TestCheckResourceAttr(resourceName, "name", "AwsDirectConnectTest"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestMatchResourceAttr(resourceName, "aws_account_id", regexp.MustCompile("[0-9]{12}")),
				),
			},
		},
	})
}

func TestAWSConnection_updateSpeed(t *testing.T) {

	resourceName := "pureport_aws_connection.main"
	var instance client.AwsDirectConnectConnection
	var respawn_instance client.AwsDirectConnectConnection

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAWSConnectionConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceAWSConnection(resourceName, &instance),
					resource.TestCheckResourceAttrPtr(resourceName, "id", &instance.Id),
					resource.TestCheckResourceAttr(resourceName, "name", "AwsDirectConnectTest"),
					resource.TestCheckResourceAttr(resourceName, "speed", "100"),
				),
			},
			{
				Config: testAccResourceAWSConnectionConfig_basic_update_speed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceAWSConnection(resourceName, &respawn_instance),
					resource.TestCheckResourceAttrPtr(resourceName, "id", &respawn_instance.Id),
					TestCheckResourceConnectionIdChanged(&instance.Id, &respawn_instance.Id),
					resource.TestCheckResourceAttr(resourceName, "name", "AwsDirectConnectTest"),
					resource.TestCheckResourceAttr(resourceName, "speed", "200"),
				),
			},
		},
	})
}

func TestAWSConnection_cloudServices(t *testing.T) {

	resourceName := "pureport_aws_connection.main"
	var instance client.AwsDirectConnectConnection

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAWSConnectionConfig_cloudServices,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceAWSConnection(resourceName, &instance),
					resource.TestCheckResourceAttrPtr(resourceName, "id", &instance.Id),
					resource.TestCheckResourceAttr(resourceName, "name", "AwsDirectConnectCloudServicesTest"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "speed", "100"),
					resource.TestCheckResourceAttr(resourceName, "high_availability", "true"),
					resource.TestCheckResourceAttr(resourceName, "location_href", "/locations/us-sea"),
					resource.TestMatchResourceAttr(resourceName, "network_href", regexp.MustCompile("/networks/network-.{16}")),
					resource.TestCheckResourceAttr(resourceName, "nat_config.0.enabled", "false"),

					resource.TestCheckResourceAttr(resourceName, "gateways.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.availability_domain", "PRIMARY"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.name", "AWS_DIRECT_CONNECT"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.link_state", "PENDING"),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.customer_asn", "7224"),
					resource.TestMatchResourceAttr(resourceName, "gateways.0.customer_ip", regexp.MustCompile("45.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}")),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.pureport_asn", "394351"),
					resource.TestMatchResourceAttr(resourceName, "gateways.0.pureport_ip", regexp.MustCompile("45.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}")),
					resource.TestCheckResourceAttrSet(resourceName, "gateways.0.bgp_password"),
					resource.TestMatchResourceAttr(resourceName, "gateways.1.peering_subnet", regexp.MustCompile("45.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}")),
					resource.TestCheckResourceAttr(resourceName, "gateways.0.public_nat_ip", ""),
					resource.TestCheckResourceAttrSet(resourceName, "gateways.0.vlan"),
					resource.TestCheckResourceAttrSet(resourceName, "gateways.0.remote_id"),

					resource.TestCheckResourceAttr(resourceName, "gateways.1.availability_domain", "SECONDARY"),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.name", "AWS_DIRECT_CONNECT 2"),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.description", ""),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.link_state", "PENDING"),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.customer_asn", "7224"),
					resource.TestMatchResourceAttr(resourceName, "gateways.1.customer_ip", regexp.MustCompile("45.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}")),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.pureport_asn", "394351"),
					resource.TestMatchResourceAttr(resourceName, "gateways.1.pureport_ip", regexp.MustCompile("45.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}")),
					resource.TestCheckResourceAttrSet(resourceName, "gateways.1.bgp_password"),
					resource.TestMatchResourceAttr(resourceName, "gateways.1.peering_subnet", regexp.MustCompile("45.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}")),
					resource.TestCheckResourceAttr(resourceName, "gateways.1.public_nat_ip", ""),
					resource.TestCheckResourceAttrSet(resourceName, "gateways.1.vlan"),
					resource.TestCheckResourceAttrSet(resourceName, "gateways.1.remote_id"),
				),
			},
		},
	})
}

func TestAWSConnection_nat_mappings(t *testing.T) {

	resourceName := "pureport_aws_connection.main"
	var instance client.AwsDirectConnectConnection

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAWSConnectionConfig_nat_mapping,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceAWSConnection(resourceName, &instance),
					resource.TestCheckResourceAttrPtr(resourceName, "id", &instance.Id),
					resource.TestCheckResourceAttr(resourceName, "name", "AwsDirectConnectNatMappingTest"),
					resource.TestCheckResourceAttr(resourceName, "nat_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "nat_config.0.mappings.#", "2"),
					//					resource.TestCheckResourceAttr(resourceName, "nat_config.0.mappings.0.native_cidr", "192.168.0.0/24"),
					//					resource.TestCheckResourceAttrSet(resourceName, "nat_config.0.mappings.0.nat_cidr"),
					//					resource.TestCheckResourceAttr(resourceName, "nat_config.0.mappings.1.native_cidr", "192.200.0.0/16"),
					//					resource.TestCheckResourceAttrSet(resourceName, "nat_config.0.mappings.1.nat_cidr"),
					resource.TestCheckResourceAttr(resourceName, "nat_config.0.blocks.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "nat_config.0.pnat_cidr"),
				),
			},
		},
	})
}

func testAccCheckResourceAWSConnection(name string, instance *client.AwsDirectConnectConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		config, ok := testAccProvider.Meta().(*Config)
		if !ok {
			return fmt.Errorf("Error getting Pureport client")
		}

		// Find the state object
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find AWS Connnection resource: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id := rs.Primary.ID

		ctx := config.Session.GetSessionContext()
		found, resp, err := config.Session.Client.ConnectionsApi.GetConnection(ctx, id)

		if err != nil {
			return fmt.Errorf("receive error when requesting AWS Connection %s", id)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("Error getting AWS Connection ID %s: %s", id, err)
		}

		*instance = found.(client.AwsDirectConnectConnection)

		return nil
	}
}

func testAccCheckAWSConnectionDestroy(s *terraform.State) error {

	config, ok := testAccProvider.Meta().(*Config)
	if !ok {
		return fmt.Errorf("Error getting Pureport client")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pureport_aws_connection" {
			continue
		}

		id := rs.Primary.ID

		ctx := config.Session.GetSessionContext()
		_, resp, err := config.Session.Client.ConnectionsApi.GetConnection(ctx, id)

		if err != nil && resp.StatusCode != 404 {
			return fmt.Errorf("should not get error for AWS Connection with ID %s after delete: %s", id, err)
		}
	}

	return nil
}
