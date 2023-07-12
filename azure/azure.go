package azure

import (
	"context"
	"log"
	"os"

	// "github.com/atoscerebro/jarvis/utils" //???

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/vision-cli/vision-plugin-infra-v1/placeholders"
	// "github.com/vision-cli/vision/config"
)

// Create the infra prequisites for Azure. We create these using az commands so that
// we can use the same random string suffix
// Infra prequisites are:
// - resource group - will be projct name + "-rg" + random string
// - storage account - will be projct name + "-sa" + random string
//   we're storing the terraform state in azure blob storage in the project resource group
//   there is no central resource group for all project states to reduce the blast
//   radius of a compromised account or deleted / corrupted state file
// - keyvault - will be projct name + "-kv" + random string
// - container registry - will be projct name + "-cr" + random string

var (
	resourceGroupName     string = placeholders.GetDefaultAzureResourceGroupName()
	// resourceGroupLocation string = placeholders.GetDefaultAzureLocation()
	accountName						string = "infra-plugin-stg-acc"
	// templateFile          string = "template.json"
	ctx                          = context.Background()
)

var accCreateParams armstorage.AccountCreateParameters = armstorage.AccountCreateParameters{
	Kind:     to.Ptr(armstorage.KindStorage),
	Location: to.Ptr("uksouth"),
	Properties: &armstorage.AccountPropertiesCreateParameters{
		AllowBlobPublicAccess:        to.Ptr(false),
		AllowSharedKeyAccess:         to.Ptr(true),
		DefaultToOAuthAuthentication: to.Ptr(false),
		Encryption: &armstorage.Encryption{
			KeySource:                       to.Ptr(armstorage.KeySourceMicrosoftStorage),
			RequireInfrastructureEncryption: to.Ptr(false),
			Services: &armstorage.EncryptionServices{
				Blob: &armstorage.EncryptionService{
					Enabled: to.Ptr(true),
					KeyType: to.Ptr(armstorage.KeyTypeAccount),
				},
				File: &armstorage.EncryptionService{
					Enabled: to.Ptr(true),
					KeyType: to.Ptr(armstorage.KeyTypeAccount),
				},
			},
		},
		IsHnsEnabled:  to.Ptr(true),
		IsSftpEnabled: to.Ptr(true),
		KeyPolicy: &armstorage.KeyPolicy{
			KeyExpirationPeriodInDays: to.Ptr[int32](20),
		},
		MinimumTLSVersion: to.Ptr(armstorage.MinimumTLSVersionTLS12),
		RoutingPreference: &armstorage.RoutingPreference{
			PublishInternetEndpoints:  to.Ptr(true),
			PublishMicrosoftEndpoints: to.Ptr(true),
			RoutingChoice:             to.Ptr(armstorage.RoutingChoiceMicrosoftRouting),
		},
		SasPolicy: &armstorage.SasPolicy{
			ExpirationAction:    to.Ptr(armstorage.ExpirationActionLog),
			SasExpirationPeriod: to.Ptr("1.15:59:59"),
		},
	},
	SKU: &armstorage.SKU{
		Name: to.Ptr(armstorage.SKUNameStandardGRS),
	},
	Tags: map[string]*string{
		"key1": to.Ptr("value1"),
		"key2": to.Ptr("value2"),
	},
} 

func EngageAzure() error {
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")

	err := createStorageAccount(subscriptionId)
	if err != nil {
		log.Printf("failed to create storage account: %v", err)
		return err
	}

	return nil
}

func createStorageAccount(subscriptionId string) error {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}

	clientFactory, err := armstorage.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return err
	}

	client := clientFactory.NewAccountsClient()

	poller, err := client.BeginCreate(ctx, resourceGroupName, accountName, accCreateParams, nil)
	if err != nil {
		return err
	}

	// res not used currently, so is a blank identifier 
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	
	return nil
}





