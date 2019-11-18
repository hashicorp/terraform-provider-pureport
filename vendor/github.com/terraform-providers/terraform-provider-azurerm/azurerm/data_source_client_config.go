package azurerm

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceArmClientConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmClientConfigRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subscription_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"object_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_principal_application_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "This has been deprecated in favour of the `client_id` property",
			},
			"service_principal_object_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "This has been deprecated in favour of the unified `authenticated_object_id` property",
			},
		},
	}
}

func dataSourceArmClientConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	var servicePrincipal *graphrbac.ServicePrincipal
	if client.usingServicePrincipal {
		spClient := client.Graph.ServicePrincipalsClient
		// Application & Service Principal is 1:1 per tenant. Since we know the appId (client_id)
		// here, we can query for the Service Principal whose appId matches.
		filter := fmt.Sprintf("appId eq '%s'", client.clientId)
		listResult, listErr := spClient.List(ctx, filter)

		if listErr != nil {
			return fmt.Errorf("Error listing Service Principals: %#v", listErr)
		}

		if listResult.Values() == nil || len(listResult.Values()) != 1 {
			return fmt.Errorf("Unexpected Service Principal query result: %#v", listResult.Values())
		}

		servicePrincipal = &(listResult.Values())[0]
	}

	d.SetId(time.Now().UTC().String())
	d.Set("client_id", client.clientId)
	d.Set("tenant_id", client.tenantId)
	d.Set("subscription_id", client.subscriptionId)

	if principal := servicePrincipal; principal != nil {
		d.Set("service_principal_application_id", client.clientId)
		d.Set("service_principal_object_id", principal.ObjectID)
	} else {
		d.Set("service_principal_application_id", "")
		d.Set("service_principal_object_id", "")
	}

	d.Set("object_id", "")

	// TODO remove this when we confirm that MSI no longer returns nil with getAuthenticatedObjectID
	objectId := ""
	if client.getAuthenticatedObjectID != nil {
		v, err := client.getAuthenticatedObjectID(ctx)
		if err != nil {
			return fmt.Errorf("Error getting authenticated object ID: %v", err)
		}
		objectId = v
	}
	d.Set("object_id", objectId)

	return nil
}
