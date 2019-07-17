package pureport

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-aws/aws"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm"
	"github.com/terraform-providers/terraform-provider-google/google"
	"github.com/terraform-providers/terraform-provider-template/template"
)

var (
	testAccProviders        map[string]terraform.ResourceProvider
	testAccProvider         *schema.Provider
	testAccGoogleProvider   *schema.Provider
	testAccAWSProvider      *schema.Provider
	testAccTemplateProvider *schema.Provider
	testAccAzureProvider    *schema.Provider
)

var testEnvironmentName string = "Production"

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccTemplateProvider = template.Provider().(*schema.Provider)
	testAccGoogleProvider = google.Provider().(*schema.Provider)
	testAccAWSProvider = aws.Provider().(*schema.Provider)
	testAccAzureProvider = azurerm.Provider().(*schema.Provider)

	testAccProviders = map[string]terraform.ResourceProvider{
		"template": testAccTemplateProvider,
		"pureport": testAccProvider,
		"google":   testAccGoogleProvider,
		"aws":      testAccAWSProvider,
		"azurerm":  testAccAzureProvider,
	}

	// Environment Variables for the Test Environment
	if env := os.Getenv("PUREPORT_ACC_TEST_ENVIRONMENT"); env != "" {
		testEnvironmentName = env
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
		"ARM_CLIENT_ID",
		"ARM_CLIENT_SECRET",
		"ARM_SUBSCRIPTION_ID",
		"ARM_TENANT_ID",
		"ARG_USE_MSI",
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
