package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	// "github.com/atoscerebro/jarvis/utils" //???

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/vision-cli/vision/config"
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
	resourceGroupName string 			= config.DefaultAzureResourceGroupName()
	resourceGroupLcoation string  = config.DefaultAzureLocation()
	deploymentName string 				= "vision-cli-deployment" // what is this in vision?
	templateFile string 					= "template.json"
	ctx = context.Background()
)


func readJSON(path string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("failed to read file: %v", err)
		return nil, err
	}
	contents := make(map[string]interface{})
	_ = json.Unmarshal(data, &contents)
	return contents, nil
}

func main() {
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")
	
}