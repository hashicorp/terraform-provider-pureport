package pureport

import (
	//	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	//	"github.com/hashicorp/terraform/terraform"
	//	"github.com/pureport/pureport-sdk-go"
)

// Global MutexKV
//var mutexKV = mutexkv.NewMutexKV()

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key":  "",
		"secret_key":  "",
		"profile":     "",
		"token":       "",
		"max_retries": "",
	}
}

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["access_key"],
			},

			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["secret_key"],
			},

			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["profile"],
			},

			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["token"],
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     25,
				Description: descriptions["max_retries"],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"pureport_aws_connection":          resourceAWSConnection(),
			"pureport_azure_connection":        resourceAzureConnection(),
			"pureport_google_cloud_connection": resourceGoogleCloudConnection(),
			"pureport_dummy_connection":        resourceDummyConnection(),
			"pureport_network":                 resourceNetwork(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
	}
}

//func providerConfigure(d *schema.ResourceData) (interface{}, error) {
//	config := Config{}
//}
