package plugin

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/vision-cli/vision/cli"
	"github.com/atoscerebro/jarvis/config"
	"github.com/atoscerebro/jarvis/execute"
	"github.com/atoscerebro/jarvis/utils"
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

func AzureChecks(silent bool, executor execute.Executor) error {
	if !executor.CommandExists("az") {
		return fmt.Errorf("azure cli is not installed")
	}

	if err := checkContainerAppExtension(executor); err != nil {
		return err
	}

	config.SetProvider(ProviderAzure)

	if err := setAzureTenant(executor); err != nil {
		return err
	}

	// Create the resource group if it doesnt exist
	if err := checkResourceGrp(silent, executor); err != nil {
		return err
	}

	// Create the main components in the resource group
	if err := checkAzureStorageAccount(silent, executor); err != nil {
		return err
	}

	if err := checkKeyvault(silent, executor); err != nil {
		return err
	}

	if err := checkAcr(silent, executor); err != nil {
		return err
	}

	if err := loginAcr(executor); err != nil {
		return err
	}

	// Create key items needed in the components
	// if err := checkCert(); err != nil {
	// 	return err
	// }

	if err := buildGraphqlContainer(executor); err != nil {
		return err
	}

	return nil
}

func checkContainerAppExtension(executor execute.Executor) error {
	c := exec.Command("az", "extension", "list")
	o, err := executor.Output(c, "", "Checking az cli container app extension")
	if err != nil {
		return err
	}

	type Extension struct {
		Name string `json:"name"`
	}

	var extensions []Extension
	err = json.Unmarshal([]byte(o), &extensions)
	if err != nil {
		return err
	}

	found := false
	for _, e := range extensions {
		if e.Name == "containerapp" {
			found = true
			break
		}
	}

	if !found {
		cli.Warningf("The Azure CLI containerapp extension was not found.")
		if !cli.Confirmed("Do you want to insatll it") {
			return fmt.Errorf("The containerapp extension is required")
		}

		c := exec.Command("az", "extension", "add", "--name", "containerapp")
		_, err := executor.Output(c, "", "Installing az cli container app extension")
		if err != nil {
			return err
		}
	}
	return nil
}

func setAzureTenant(executor execute.Executor) error {
	c := exec.Command("az", "account", "show")
	a, err := executor.Output(c, "", "Checking resource group")
	if err != nil {
		log.Fatal(err)
	}

	type Account struct {
		Name string `json:"name"`
	}

	var account Account
	err = json.Unmarshal([]byte(a), &account)
	if err != nil {
		log.Fatal(err)
	}

	config.SetAzureTenant(account.Name)

	return nil
}

func execCommandCheckResultExists[T comparable](c *exec.Cmd, message string, match T, executor execute.Executor) (bool, error) {
	output, err := executor.Output(c, "", message)
	if err != nil {
		return false, err
	}

	var result []T
	err = json.Unmarshal([]byte(output), &result)
	if err != nil {
		return false, err
	}

	for _, r := range result {
		if r == match {
			return true, nil
		}
	}

	return false, nil
}

func checkResourceGrp(silent bool, executor execute.Executor) error {
	type ResourceGroup struct {
		Name     string `json:"name"`
		Location string `json:"location"`
	}

	var match ResourceGroup
	if config.AzureResourceGroup() == "" {
		match = ResourceGroup{
			Name:     config.DefaultAzureResourceGroupName(),
			Location: config.DefaultAzureLocation(),
		}
	} else {
		match = ResourceGroup{
			Name:     config.AzureResourceGroup(),
			Location: config.AzureLocation(),
		}
	}

	c := exec.Command("az", "group", "list")
	exists, err := execCommandCheckResultExists(c, "Checking resource group", match, executor)
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		cli.Warningf(fmt.Sprintf("Resource group %s exists", match.Name))
		return nil
	}

	if !silent {
		cli.Warningf(fmt.Sprintf("Resource group %s not found", match.Name))
		if !cli.Confirmed("Do you want to create it ") {
			return fmt.Errorf(fmt.Sprintf("Resource group %s not found", match.Name))
		}
		c = exec.Command("az", "account", "list-locations", "--output", "table")
		var locations string
		locations, err = executor.Output(c, "", "Display available locations")
		if err != nil {
			panic(err)
		}
		fmt.Println(locations)
		match.Location = cli.Input("Enter location: ", "uksouth", true)
		match.Location = utils.CleanString(match.Location)
	}
	if silent {
		match.Location = config.DefaultAzureLocation()
	}
	config.SetAzureLocation(match.Location)

	c = exec.Command("az", "group", "create", "--name", match.Name, "--location", match.Location) //nolint:gosec // we are cleaning the location string
	return executor.Errors(c, "", "Creating resource group")
}

func checkAzureStorageAccount(silent bool, executor execute.Executor) error {
	type StorageAccount struct {
		Name string `json:"name"`
	}

	var match StorageAccount
	if config.AzureStorageAccount() == "" {
		match = StorageAccount{
			Name: config.DefaultAzureStorageAccount(),
		}
	} else {
		match = StorageAccount{
			Name: config.AzureStorageAccount(),
		}
	}

	c := exec.Command("az", "storage", "account", "list")
	exists, err := execCommandCheckResultExists(c, "Checking storage account", match, executor)
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		cli.Warningf(fmt.Sprintf("Storage account %s exists", match.Name))
		return nil
	}

	if !silent {
		sa := cli.InputWithValidation("Enter storage account name", config.DefaultAzureStorageAccount(), true, storageAccValidation, executor)
		config.SetAzureStorageAccount(sa)
	} else {
		config.SetAzureStorageAccount(config.DefaultAzureStorageAccount())
	}

	{
		c := exec.Command("az", //nolint:gosec // we are cleaning the storage account string
			"storage",
			"account",
			"create",
			"--resource-group", config.AzureResourceGroup(),
			"--name", config.AzureStorageAccount(),
			"--sku", "Standard_LRS",
			"--encryption-service", "blob",
			"--allow-blob-public-access",
			"--min-tls-version", "TLS1_2",
		)
		if err := executor.Errors(c, "", "Creating Azure storage account"); err != nil {
			log.Fatal(err)
		}
	}

	{
		c := exec.Command("az", //nolint:gosec // we are cleaning the storage account string
			"storage",
			"container",
			"create",
			"--name", "tfstate",
			"--account-name", config.AzureStorageAccount(),
		)
		if err := executor.Errors(c, "", "Creating Azure storage container"); err != nil {
			log.Fatal(err)
		}
	}

	{
		c := exec.Command("az", //nolint:gosec // we are cleaning the storage account string
			"storage",
			"account",
			"blob-service-properties",
			"update",
			"-g", config.AzureResourceGroup(),
			"--account-name", config.AzureStorageAccount(),
			"--enable-versioning",
			"--enable-delete-retention", "true",
			"--delete-retention-days", "7",
			"--enable-container-delete-retention", "true",
			"--container-delete-retention-days", "7",
		)
		if err := executor.Errors(c, "", "Creating Azure storage container"); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func checkKeyvault(silent bool, executor execute.Executor) error {
	type KeyVault struct {
		Name string `json:"name"`
	}

	var match KeyVault
	if config.AzureStorageAccount() == "" {
		match = KeyVault{
			Name: config.DefaultAzureKeyvault(),
		}
	} else {
		match = KeyVault{
			Name: config.AzureKeyvault(),
		}
	}

	c := exec.Command("az", "keyvault", "list")
	exists, err := execCommandCheckResultExists(c, "Checking keyvault", match, executor)
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		cli.Warningf(fmt.Sprintf("Keyvault %s exists", match.Name))
		return nil
	}

	if !silent {
		kv := cli.InputWithValidation("Enter keyvault name", config.DefaultAzureKeyvault(), true, keyvaultValidation, executor)
		config.SetAzureKeyvault(kv)
	} else {
		config.SetAzureKeyvault(config.DefaultAzureKeyvault())
	}

	c = exec.Command("az", //nolint:gosec // we are cleaning the storage account string
		"keyvault",
		"create",
		"--location", config.AzureLocation(),
		"--resource-group", config.AzureResourceGroup(),
		"--name", config.AzureKeyvault(),
	)
	if err := executor.Errors(c, "", "Creating Azure keyvault"); err != nil {
		return err
	}

	return nil
}

// func checkCert() error {
// 	c := exec.Command("az",
// 		"keyvault",
// 		"certificate",
// 		"get-default-policy",
// 	)
// 	defaultPolicy, err := execute.Output(c, "", "Getting Keyvault default policy")
// 	if err != nil {
// 		return err
// 	}

// 	c = exec.Command("az", //nolint:gosec // we are cleaning the storage account string
// 		"keyvault",
// 		"certificate",
// 		"create",
// 		"--vault-name", config.AzureKeyvault(),
// 		"-n", "project-cert",
// 		"-p", defaultPolicy,
// 	)
// 	output, err := execute.Output(c, "", "Creating Project certificate")
// 	if err != nil {
// 		return err
// 	}
// 	cli.Warningf(output)

// 	c = exec.Command("az", //nolint:gosec // we are cleaning the storage account string
// 		"keyvault",
// 		"certificate",
// 		"download",
// 		"--vault-name", config.AzureKeyvault(),
// 		"-n", "project-cert",
// 		"-f", "cert.pem",
// 	)
// 	output, err = execute.Output(c, "", "Downloading Project certificate")
// 	if err != nil {
// 		return err
// 	}
// 	cli.Warningf(output)

// 	return nil
// }

func checkAcr(silent bool, executor execute.Executor) error {
	type Acr struct {
		Name string `json:"name"`
	}

	var match Acr
	if config.AzureStorageAccount() == "" {
		match = Acr{
			Name: config.DefaultAzureAcr(),
		}
	} else {
		match = Acr{
			Name: config.AzureAcr(),
		}
	}

	c := exec.Command("az", "acr", "list")
	exists, err := execCommandCheckResultExists(c, "Checking acr", match, executor)
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		cli.Warningf(fmt.Sprintf("Acr %s exists", match.Name))
		return nil
	}

	if !silent {
		acr := cli.InputWithValidation("Enter acr name", config.DefaultAzureAcr(), true, acrValidation, executor)
		config.SetAzureAcr(acr)
	} else {
		config.SetAzureAcr(config.DefaultAzureAcr())
	}

	c = exec.Command("az", //nolint:gosec // we are cleaning the storage account string
		"acr",
		"create",
		"--resource-group", config.AzureResourceGroup(),
		"--name", config.AzureAcr(),
		"--sku", "Standard",
		"--admin-enabled",
	)
	output, err := executor.Output(c, "", "Creating Azure acr")
	if err != nil {
		return err
	}

	type ContainerRegistry struct {
		LoginServer string `json:"loginServer"`
	}

	var containerRegistry ContainerRegistry
	err = json.Unmarshal([]byte(output), &containerRegistry)
	if err != nil {
		return err
	}
	config.SetAzureAcrLoginServer(containerRegistry.LoginServer)

	return nil
}

func loginAcr(executor execute.Executor) error {
	c := exec.Command("az", //nolint:gosec // we are cleaning the storage account string
		"acr",
		"login",
		"--name", config.AzureAcr(),
	)

	return executor.Errors(c, "", "Logging in to Azure acr")
}

func buildGraphqlContainer(executor execute.Executor) error {
	output, err := executor.Output(exec.Command("docker", "images"), "", "Checking if graphql container exists")
	if err != nil {
		return err
	}

	if strings.Contains(output, config.AzureAcrLoginServer()+"/graphql/server") {
		cli.Warningf("Graphql container exists")
		return nil
	}

	c := exec.Command("docker", //nolint:gosec // we are cleaning the storage account string
		"buildx",
		"build",
		"--platform", "linux/amd64",
		"--push",
		"-f", "./infra/docker/standalone-graphql/Dockerfile",
		"-t", config.AzureAcrLoginServer()+"/graphql/server",
		".",
	)

	return executor.Errors(c, "", "Building standalone graphql container")
}

func storageAccValidation(input string, executor execute.Executor) (bool, string) {
	c := exec.Command("az", "storage", "account", "check-name", "--name", input)
	output, err := executor.Output(c, "", "Checking Azure storage account name")
	if err != nil {
		log.Fatal(err)
	}
	return checknameValidation(output)
}

func keyvaultValidation(input string, executor execute.Executor) (bool, string) {
	c := exec.Command("az", "keyvault", "check-name", "--name", input)
	output, err := executor.Output(c, "", "Checking Azure keyvault name")
	if err != nil {
		log.Fatal(err)
	}
	return checknameValidation(output)
}

func acrValidation(input string, executor execute.Executor) (bool, string) {
	c := exec.Command("az", "acr", "check-name", "--name", input)
	output, err := executor.Output(c, "", "Checking acr name")
	if err != nil {
		log.Fatal(err)
	}
	return checknameValidation(output)
}

func checknameValidation(input string) (bool, string) {
	type Result struct {
		Message       string `json:"message"`
		NameAvailable bool   `json:"nameAvailable"`
	}
	var res Result
	err := json.Unmarshal([]byte(input), &res)
	if err != nil {
		panic(err)
	}

	if res.NameAvailable {
		return true, ""
	}
	return false, res.Message
}