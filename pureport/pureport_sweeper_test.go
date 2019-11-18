package pureport

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-pureport/pureport/configuration"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// sharedClientForRegion returns a common provider client configured for the specified region
func sharedClientForRegion(region string) (interface{}, error) {

	config := configuration.Config{}

	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}
