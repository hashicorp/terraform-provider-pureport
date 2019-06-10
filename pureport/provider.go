package pureport

import (
	//	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Global MutexKV
//var mutexKV = mutexkv.NewMutexKV()

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"api_key":      "Pureport API Key",
		"api_secret":   "Pureport API Secret",
		"api_url":      "Pureport API URL to execute against",
		"auth_profile": "The authentication profile in your local Pureport configuration file.",
	}
}

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["api_key"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"PUREPORT_API_KEY",
				}, nil),
			},

			"api_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["secret_key"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"PUREPORT_API_SECRET",
				}, nil),
			},

			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["api_url"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"PUREPORT_ENDPOINT",
				}, nil),
			},

			"auth_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["auth_profile"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"PUREPORT_PROFILE",
				}, nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"pureport_aws_connection":          resourceAWSConnection(),
			"pureport_azure_connection":        resourceAzureConnection(),
			"pureport_google_cloud_connection": resourceGoogleCloudConnection(),
			"pureport_site_vpn_connection":     resourceSiteVPNConnection(),
			"pureport_network":                 resourceNetwork(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"pureport_cloud_regions":  dataSourceCloudRegions(),
			"pureport_cloud_services": dataSourceCloudServices(),
			"pureport_locations":      dataSourceLocations(),
			"pureport_networks":       dataSourceNetworks(),
			"pureport_accounts":       dataSourceAccounts(),
			"pureport_connections":    dataSourceConnections(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	config := Config{}

	if v, ok := d.GetOk("auth_profile"); ok {
		config.AuthenticationProfile = v.(string)
	}

	if v, ok := d.GetOk("api_key"); ok {
		config.APIKey = v.(string)
	}

	if v, ok := d.GetOk("api_secret"); ok {
		config.APISecret = v.(string)
	}

	if v, ok := d.GetOk("api_url"); ok {
		config.EndPoint = v.(string)
	}

	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}
