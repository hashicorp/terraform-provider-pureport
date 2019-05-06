package pureport

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-google/google"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testAccGoogleProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccGoogleProvider = google.Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"pureport": testAccProvider,
		"google":   testAccGoogleProvider,
	}
}

func testAccPreCheck(t *testing.T) {

	pureportEnvVars := []string{
		"PUREPORT_API_KEY",
		"PUREPORT_API_SECRET",
		"PUREPORT_ENDPOINT",
	}

	googleEnvVars := []string{
		"GOOGLE_CREDENTIALS",
		"GOOGLE_PROJECT",
		"GOOGLE_REGION",
	}

	// Pureport Provider Configuration
	for _, e := range pureportEnvVars {
		if v := os.Getenv(e); v == "" {
			t.Fatalf("%s must be specified for acceptance tests", e)
		}
	}

	// Google Cloud Provider Configuration
	for _, e := range googleEnvVars {
		if v := os.Getenv(e); v == "" {
			t.Fatalf("%s must be specified for acceptance tests", e)
		}
	}

	if err := testAccProvider.Configure(terraform.NewResourceConfig(nil)); err != nil {
		t.Fatal(err)
	}
}
