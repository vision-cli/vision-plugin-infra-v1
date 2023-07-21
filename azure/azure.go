package azure

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/vision-cli/common/execute"

	// "github.com/vision-cli/common/tmpl"
	"github.com/vision-cli/vision-plugin-infra-v1/placeholders"
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
	resourceGroupName string = placeholders.GetDefaultAzureResourceGroupName()
	accountName       string = "infrapluginstgacc"
	ctx                      = context.Background()
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
		KeyPolicy: &armstorage.KeyPolicy{
			KeyExpirationPeriodInDays: to.Ptr[int32](20),
		},
		MinimumTLSVersion: to.Ptr(armstorage.MinimumTLSVersionTLS12),
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

func EngageAzure(executor execute.Executor) error {
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")

	fmt.Println("creating storage account (azure)")
	if err := createStorageAccount(subscriptionId); err != nil {
		return fmt.Errorf("failed to create storage account: %v", err)
	}

	//create tfstate container
	if err := createTfStateContainer(subscriptionId); err != nil {
		return fmt.Errorf("creating terraform state container: %v", err)
	}

	fmt.Println("executing make init (terraform)")
	c := exec.Command("make", "init")

	if err := executor.Errors(c, "./azure/_templates/az/tf/", "inititalise Terraform"); err != nil {
		return fmt.Errorf("executing Terraform make init: %v", err)
	}

	// make plan

	// make apply

	return nil
}

func createStorageAccount(subscriptionId string) error {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("creating new default Azure credential: %v", err)
	}


	clientFactory, err := armstorage.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return fmt.Errorf("creating client factory: %v", err)
	}

	client := clientFactory.NewAccountsClient()

	poller, err := client.BeginCreate(ctx, resourceGroupName, accountName, accCreateParams, nil)
	if err != nil {
		return fmt.Errorf("creating poller: %v", err)
	}

	// res not used currently, so is a blank identifier
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("polling until done: %v", err)
	}

	return nil
}

func createTfStateContainer(subscriptionId string) error {
	storageAccountNameUrl := "https://infrapluginstgacc.blob.core.windows.net"
	containerName := "tfstate"

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("creating new default Azure credential: %v", err)
	}

	opts := azblob.ClientOptions {}

	client, err := azblob.NewClient(storageAccountNameUrl, cred, &opts)
	if err != nil {
		return fmt.Errorf("create container client: %v", err)
	}

	_, err = client.CreateContainer(ctx, containerName, nil)
	if err != nil {
		return fmt.Errorf("create tfstate container: %v", err)
	}

	return nil
}

// func handleError(str string, err error) error {
// 	if err != nil {
// 		return fmt.Errorf(str + ": %v", err)
// 	}
// 	return nil
// }