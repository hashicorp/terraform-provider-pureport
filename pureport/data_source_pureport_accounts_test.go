package pureport

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testAccDataSourceAccountsConfig_empty = `
data "pureport_accounts" "empty" {
}
`

const testAccDataSourceAccountsConfig_name_filter = `
data "pureport_accounts" "name_filter" {
  filter {
    name = "Name"
    values = ["Terraform .*"]
  }
}
`

func TestDataSourceAccounts_empty(t *testing.T) {

	resourceName := "data.pureport_accounts.empty"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccountsConfig_empty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAccounts(resourceName),
					resource.TestCheckResourceAttr(resourceName, "accounts.#", "2"),
					testAccCheckDataSourceAccountMain(resourceName, "accounts.0"),
					testAccCheckDataSourceAccountChildAccount(resourceName, "accounts.1"),
				),
			},
		},
	})
}

func TestDataSourceAccounts_name_filter(t *testing.T) {

	resourceName := "data.pureport_accounts.name_filter"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccountsConfig_name_filter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAccounts(resourceName),
					resource.TestCheckResourceAttr(resourceName, "accounts.#", "1"),
					testAccCheckDataSourceAccountChildAccount(resourceName, "accounts.0"),
				),
			},
		},
	})
}

func testAccCheckDataSourceAccounts(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		// Find the state object
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find Accounts data source: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}

func testAccCheckDataSourceAccountMain(resourceName, account string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestMatchResourceAttr(resourceName, account+".id", regexp.MustCompile("ac-.{16}")),
		resource.TestMatchResourceAttr(resourceName, account+".href", regexp.MustCompile("/accounts/ac-.{16}")),
		resource.TestCheckResourceAttr(resourceName, account+".name", "HashiCorp"),
		resource.TestCheckResourceAttr(resourceName, account+".description", "Developer Account for Testing Terraform Provider"),
		resource.TestCheckResourceAttr(resourceName, account+".tags.#", "0"),
	)
}

func testAccCheckDataSourceAccountChildAccount(resourceName, account string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestMatchResourceAttr(resourceName, account+".id", regexp.MustCompile("ac-.{16}")),
		resource.TestMatchResourceAttr(resourceName, account+".href", regexp.MustCompile("/accounts/ac-.{16}")),
		resource.TestCheckResourceAttr(resourceName, account+".name", "Terraform Acceptance Tests"),
		resource.TestCheckResourceAttr(resourceName, account+".description", "Terraform Provider Acceptance Tests"),
		resource.TestCheckResourceAttr(resourceName, account+".tags.#", "0"),
	)
}
