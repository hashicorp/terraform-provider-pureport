package pureport

import (
	"fmt"
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

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccountsConfig_empty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAccounts(resourceName),
					resource.TestCheckResourceAttr(resourceName, "accounts.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0.id", "account-EhlpJLhAcHMOmY75J91H3g"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0.href", "/accounts/account-"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0.name", ""),
					resource.TestCheckResourceAttr(resourceName, "accounts.0.description", "Test Account #1"),
					resource.TestCheckResourceAttr(resourceName, "accounts.1.id", "account-fN6NX6utBCoE5L_H261P4A"),
					resource.TestCheckResourceAttr(resourceName, "accounts.1.href", "/accounts/account-"),
					resource.TestCheckResourceAttr(resourceName, "accounts.1.name", ""),
					resource.TestCheckResourceAttr(resourceName, "accounts.1.description", "Test Account #2"),
				),
			},
		},
	})
}

func TestAccounts_name_regex(t *testing.T) {

	resourceName := "data.pureport_accounts.name_regex"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccountsConfig_name_regex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAccounts(resourceName),
					resource.TestCheckResourceAttr(resourceName, "accounts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0.id", "account-fN6NX6utBCoE5L_H261P4A"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0.href", "/accounts/account-fN6NX6utBCoE5L_H261P4A"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0.name", "The Clash"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0.description", "Test Account #1"),
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
			return fmt.Errorf("Can't find Cloud Region data source: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}
