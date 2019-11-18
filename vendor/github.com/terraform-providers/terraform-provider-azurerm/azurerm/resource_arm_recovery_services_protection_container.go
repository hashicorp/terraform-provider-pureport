package azurerm

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-01-10/siterecovery"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRecoveryServicesProtectionContainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRecoveryServicesProtectionContainerCreate,
		Read:   resourceArmRecoveryServicesProtectionContainerRead,
		Update: nil,
		Delete: resourceArmSiteRecoveryProtectionContainerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),

			"recovery_vault_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateRecoveryServicesVaultName,
			},
			"recovery_fabric_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
		},
	}
}

func resourceArmRecoveryServicesProtectionContainerCreate(d *schema.ResourceData, meta interface{}) error {
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	fabricName := d.Get("recovery_fabric_name").(string)
	name := d.Get("name").(string)

	client := meta.(*ArmClient).RecoveryServices.ProtectionContainerClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, fabricName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing recovery services protection container %s (fabric %s): %+v", name, fabricName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_recovery_services_protection_container", azure.HandleAzureSdkForGoBug2824(*existing.ID))
		}
	}

	parameters := siterecovery.CreateProtectionContainerInput{
		Properties: &siterecovery.CreateProtectionContainerInputProperties{},
	}

	future, err := client.Create(ctx, fabricName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating recovery services protection container %s (fabric %s): %+v", name, fabricName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error creating recovery services protection container %s (fabric %s): %+v", name, fabricName, err)
	}

	resp, err := client.Get(ctx, fabricName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving site recovery protection container %s (fabric %s): %+v", name, fabricName, err)
	}

	d.SetId(azure.HandleAzureSdkForGoBug2824(*resp.ID))

	return resourceArmRecoveryServicesProtectionContainerRead(d, meta)
}

func resourceArmRecoveryServicesProtectionContainerRead(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	fabricName := id.Path["replicationFabrics"]
	name := id.Path["replicationProtectionContainers"]

	client := meta.(*ArmClient).RecoveryServices.ProtectionContainerClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, fabricName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on recovery services protection container %s (fabric %s): %+v", name, fabricName, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("recovery_vault_name", vaultName)
	d.Set("recovery_fabric_name", fabricName)
	return nil
}

func resourceArmSiteRecoveryProtectionContainerDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	fabricName := id.Path["replicationFabrics"]
	name := id.Path["replicationProtectionContainers"]

	client := meta.(*ArmClient).RecoveryServices.ProtectionContainerClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	future, err := client.Delete(ctx, fabricName, name)
	if err != nil {
		return fmt.Errorf("Error deleting recovery services protection container %s (fabric %s): %+v", name, fabricName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of recovery services protection container %s (fabric %s): %+v", name, fabricName, err)
	}

	return nil
}
