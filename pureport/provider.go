package pureport

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Global MutexKV
var mutexKV = mutexkv.NewMutexKV()

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{}
}

//func ResourceMap() map[string]*schema.Resource {
//	resourceMap, _ := ResourceMapWithErrors()
//	return resourceMap
//}
//
//func ResourceMapWithErrors() (map[string]*schema.Resource, error) {
//	return mergeResourceMaps()
//}
//
//func providerConfigure(d *schema.ResourceData) (interface{}, error) {
//	config := Config{}
//}
