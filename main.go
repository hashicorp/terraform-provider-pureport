package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/pureport/terraform-provider-pureport/pureport"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: pureport.Provider})
}
