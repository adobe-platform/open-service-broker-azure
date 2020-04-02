package storage

import (
	"context"
	"fmt"

	storageSDK "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (b *blobServicesManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("unsetBlobServices", b.unsetBlobServices),
		service.NewDeprovisioningStep("deleteARMDeployment", b.deleteARMDeployment),
	)
}

func (b *blobServicesManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*instanceDetails)

	if err := b.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return instance.Details, nil
}

// You don't actually delete Blob Services, but you can revert any settings
// back to the defaults
func (b *blobServicesManager) unsetBlobServices(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt := instance.Details.(*instanceDetails)
	blobProperties, err := b.blobServicesClient.GetServiceProperties(
		ctx,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		dt.StorageAccountName,
	)

	if err != nil {
		return nil, fmt.Errorf("Error getting Blob Services properties: %s", err)
	}

	newCorsRules := new([]storageSDK.CorsRule)

	fmt.Printf("Blob Services Properties Properties Before:\n%+v\n", *blobProperties.BlobServicePropertiesProperties)

	// Reset all the settings to nil
	if blobProperties.BlobServicePropertiesProperties.Cors.CorsRules != nil {
		blobProperties.BlobServicePropertiesProperties.Cors.CorsRules = newCorsRules
	}
	if blobProperties.BlobServicePropertiesProperties.DeleteRetentionPolicy != nil {
		*blobProperties.BlobServicePropertiesProperties.DeleteRetentionPolicy.Enabled = false
	}
	if blobProperties.BlobServicePropertiesProperties.DefaultServiceVersion != nil {
		*blobProperties.BlobServicePropertiesProperties.DefaultServiceVersion = ""
	}
	if blobProperties.BlobServicePropertiesProperties.AutomaticSnapshotPolicyEnabled != nil {
		*blobProperties.BlobServicePropertiesProperties.AutomaticSnapshotPolicyEnabled = false
	}

	fmt.Printf("Blob Services Properties Properties After:\n%+v\n", *blobProperties.BlobServicePropertiesProperties)

	_, err = b.blobServicesClient.SetServiceProperties(
		ctx,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		dt.StorageAccountName,
		blobProperties,
	)
	if err != nil {
		return nil, fmt.Errorf("Error unsetting Blob Services Properties: %s", err)
	}

	return dt, nil
}
