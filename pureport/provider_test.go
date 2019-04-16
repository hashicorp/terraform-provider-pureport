package pureport

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"pureport": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {

	if err := testAccProvider.Configure(terraform.NewResourceConfig(nil)); err != nil {
		t.Fatal(err)
	}
}
