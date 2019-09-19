package dns

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	RecordSetsClient *dns.RecordSetsClient
	ZonesClient      *dns.ZonesClient
}

func BuildClient(o *common.ClientOptions) *Client {

	RecordSetsClient := dns.NewRecordSetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&RecordSetsClient.Client, o.ResourceManagerAuthorizer)

	ZonesClient := dns.NewZonesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ZonesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		RecordSetsClient: &RecordSetsClient,
		ZonesClient:      &ZonesClient,
	}
}
