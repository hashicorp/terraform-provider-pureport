package pureport

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-aws/aws"
	"github.com/terraform-providers/terraform-provider-google/google"
	"github.com/terraform-providers/terraform-provider-template/template"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testAccGoogleProvider *schema.Provider
var testAccAWSProvider *schema.Provider
var testAccTemplateProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccGoogleProvider = google.Provider().(*schema.Provider)
	testAccTemplateProvider = template.Provider().(*schema.Provider)
	testAccAWSProvider = aws.Provider().(*schema.Provider)

	testAccProviders = map[string]terraform.ResourceProvider{
		"pureport": testAccProvider,
		"google":   testAccGoogleProvider,
		"template": testAccTemplateProvider,
		"aws":      testAccAWSProvider,
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

	amazonEnvVars := []string{
		"AWS_DEFAULT_REGION",
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
	}

	azureEnvVars := []string{
		"TF_VAR_azurerm_express_route_circuit_service_key",
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

	// AWS Cloud Provider Configuration
	for _, e := range amazonEnvVars {
		if v := os.Getenv(e); v == "" {
			t.Fatalf("%s must be specified for acceptance tests", e)
		}
	}

	// Azure Cloud Provider Configuration
	for _, e := range azureEnvVars {
		if v := os.Getenv(e); v == "" {
			t.Fatalf("%s must be specified for acceptance tests", e)
		}
	}

	if err := testAccProvider.Configure(terraform.NewResourceConfig(nil)); err != nil {
		t.Fatal(err)
	}
}
