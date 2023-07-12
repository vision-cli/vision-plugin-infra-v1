package placeholders

import (
	"fmt"
	"math/rand"
	"regexp"

	api_v1 "github.com/vision-cli/api/v1"
)

const (
	ArgsCommandIndex = 0
	ArgsNameIndex    = 1
	// include any other arg indexes here
)

var (
	viperAzurePrefix            = "azure"
	viperAzureResourceGroupKey  = viperAzurePrefix + ".resource-group"
	viperAzureLocationKey       = viperAzurePrefix + ".location"
	viperAzureStorageAccountKey = viperAzurePrefix + ".storage-account"
	viperAzureTennantKey        = viperAzurePrefix + ".tenant"
	viperAzureKeyvaultKey       = viperAzurePrefix + ".keyvault"
	viperAzureAcrNameKey        = viperAzurePrefix + ".acr.name"
	viperAzureAcrLoginServerKey = viperAzurePrefix + ".acr.login-server"
	maxAzureKeyLen              = 20
)

var nonAlphaRegex = regexp.MustCompile(`[^a-zA-Z]+`)

type Placeholders struct {
	Name string
}

func SetupPlaceholders(req api_v1.PluginRequest) (*Placeholders, error) {
	// setup your placeholders here
	// you can also deepcopy the Placeholders in the plugin request and use it
	// this is just an example:
	name := clearString(req.Args[ArgsNameIndex])
	return &Placeholders{
		Name: name,
	}, nil
}

func clearString(str string) string {
	return nonAlphaRegex.ReplaceAllString(str, "")
}

type InfraConfig struct {
	
}

func SetAzureResourceGroup(s string) {
	viperAzureResourceGroupKey = s
}

func GetAzureResourceGroup() string {
	return viperAzureResourceGroupKey
}

func GetDefaultAzureResourceGroupName() string {
	return "vision-infra-plugin-rg"
}

func SetAzureLocation(s string) {
	viperAzureLocationKey = s
}

func GetAzureLocation() string {
	return viperAzureLocationKey
}

func GetDefaultAzureLocation() string {
	return "uksouth"
}

func GetProjectName() string {
	return "project.name"
}

func GetAzureStorageAccount() string {
	return viperAzureStorageAccountKey
}

func GetAzureKeyvault() string {
	return viperAzureKeyvaultKey
}

func GetAzureAcr() string {
	return viperAzureAcrNameKey
}

func GetDefaultAzureApp() string {
	return GetProjectName() + "-" + fmt.Sprintf("%d", rand.Intn(100000))

}

func SaveConfig() error {
	return nil
}

func InfraDirectory() string {
	return "."
}
