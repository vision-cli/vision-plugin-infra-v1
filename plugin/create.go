package plugin

import (
	"embed"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/atoscerebro/jarvis/config"
	"github.com/atoscerebro/jarvis/execute"
	"github.com/atoscerebro/jarvis/tmpl"
)

const (
	ProviderAzure = "azure"
	ProviderAws   = "aws"
	ProviderGcp   = "gcp"
	goTemplateDir = "_templates"
)

//go:embed all:_templates
var templateFilesAz embed.FS

var createCmd = &cobra.Command{
	Use:   "create [aws|azure|gcp]",
	Short: "Create the infra assets",
	Long: `Create the infra folder with terraform assets for the cloud provider selected.
You need to have terraform installed and the cloud provider CLI installed and configured. For example on a mac for Azure

brew update && brew install terraform && brew install azure-cli

Create will also
- create a resource group if it doesnt exist
- create a storage account for the terraform state
- run terraform init and terraform apply using the service principle provided.
- A Github secret will be created`,
	ValidArgs: []string{ProviderAzure, ProviderAws, ProviderGcp},
	Args:      cobra.ExactValidArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := bindFlags(cmd, args); err != nil {
			return err
		}
		osExec := execute.NewOsExecutor()
		return providerChecks(cmd, args, osExec)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return createFolderStructure()
	},
}

func bindFlags(cmd *cobra.Command, args []string) error {
	if len(args) > 0 && args[0] == ProviderAzure {
		if err := setFlagConfig(cmd,
			FlagResourceGroup,
			config.SetAzureResourceGroup,
			config.AzureResourceGroup,
			config.DefaultAzureResourceGroupName); err != nil {
			return err
		}

		if err := setFlagConfig(cmd,
			FlagLocation,
			config.SetAzureLocation,
			config.AzureLocation,
			config.DefaultAzureLocation); err != nil {
			return err
		}

		return config.SaveConfig()
	}
	return nil
}

func setFlagConfig(cmd *cobra.Command, flag string, setter func(string), getter func() string, def func() string) error {
	f, err := cmd.Flags().GetString(flag)
	if err != nil {
		return fmt.Errorf("flag [%s]: %w", FlagLocation, err)
	}
	if f != "" {
		setter(f)
	} else if getter() == "" {
		setter(def())
	}

	return nil
}

func providerChecks(cmd *cobra.Command, args []string, executor execute.Executor) error {
	silent, err := cmd.Flags().GetBool(config.FlagSilent)
	if err != nil {
		return err
	}

	provider := args[0]

	if !executor.CommandExists("terraform") {
		return fmt.Errorf("terraform is not installed")
	}

	if !executor.CommandExists("docker") {
		return fmt.Errorf("docker is not installed")
	}

	switch provider {
	case ProviderAzure:
		return AzureChecks(silent, executor)
	case ProviderAws:
	case ProviderGcp:
	}

	return nil
}

func createFolderStructure() error {
	type Placeholders struct {
		ProjectName    string
		ResourceGroup  string
		Location       string
		StorageAccount string
		Keyvault       string
		Acr            string
		AppName        string
	}
	p := Placeholders{
		ProjectName:    tmpl.Kebab(config.ProjectName()),
		ResourceGroup:  config.AzureResourceGroup(),
		Location:       config.AzureLocation(),
		StorageAccount: config.AzureStorageAccount(),
		Keyvault:       config.AzureKeyvault(),
		Acr:            config.AzureAcr(),
		AppName:        config.DefaultAzureApp(),
	}
	if err := tmpl.GenerateFS(templateFilesAz, goTemplateDir, config.InfraDirectory(), p, false); err != nil {
		return fmt.Errorf("generating the project structure from the template: %s", err)
	}
	return nil
}