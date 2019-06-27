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

const testAccDataSourceAccountsConfig_name_regex = `
data "pureport_accounts" "name_regex" {
	name_regex = "Terraform .*"
}
`

func TestAccounts_empty(t *testing.T) {

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

					resource.TestMatchResourceAttr(resourceName, "accounts.0.id", regexp.MustCompile("ac-.{16}")),
					resource.TestMatchResourceAttr(resourceName, "accounts.0.href", regexp.MustCompile("/accounts/ac-.{16}")),
					resource.TestCheckResourceAttr(resourceName, "accounts.0.name", "HashiCorp"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0.description", "Developer Account for Testing Terraform Provider"),

					resource.TestMatchResourceAttr(resourceName, "accounts.1.id", regexp.MustCompile("ac-.{16}")),
					resource.TestMatchResourceAttr(resourceName, "accounts.1.href", regexp.MustCompile("/accounts/ac-.{16}")),
					resource.TestCheckResourceAttr(resourceName, "accounts.1.name", "Terraform Acceptance Tests"),
					resource.TestCheckResourceAttr(resourceName, "accounts.1.description", "Terraform Provider Acceptance Tests"),
				),
			},
		},
	})
}

func TestAccounts_name_regex(t *testing.T) {

	resourceName := "data.pureport_accounts.name_regex"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccountsConfig_name_regex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAccounts(resourceName),
					resource.TestCheckResourceAttr(resourceName, "accounts.#", "1"),

					resource.TestMatchResourceAttr(resourceName, "accounts.0.id", regexp.MustCompile("ac-.{16}")),
					resource.TestMatchResourceAttr(resourceName, "accounts.0.href", regexp.MustCompile("/accounts/ac-.{16}")),
					resource.TestCheckResourceAttr(resourceName, "accounts.0.name", "Terraform Acceptance Tests"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0.description", "Terraform Provider Acceptance Tests"),
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
