package storage

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (b *gpv2BlobContainerManager) Bind(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return nil, nil
}

// nolint: lll
func (b *gpv2BlobContainerManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt := instance.Details.(*instanceDetails)
	credential := credentials{
		StorageAccountName:         dt.StorageAccountName,
		AccessKey:                  dt.AccessKey,
		ContainerName:              dt.ContainerName,
		PrimaryBlobServiceEndPoint: fmt.Sprintf("https://%s.blob.core.windows.net/", dt.StorageAccountName),
	}
	return credential, nil
}
