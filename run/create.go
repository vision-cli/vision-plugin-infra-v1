package run

import (
	"embed"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vision-cli/common/execute"
	"github.com/vision-cli/common/tmpl"
	"github.com/vision-cli/vision-plugin-infra-v1/azure"
	"github.com/vision-cli/vision-plugin-infra-v1/placeholders"
)

//go:embed all:_templates
var templateFiles embed.FS

const (
	ProviderAzure = "azure"
	ProviderAws   = "aws"
	ProviderGcp   = "gcp"
	goTemplateDir = "_templates"
)

//go:embed all:_templates
var templateFilesAz embed.FS

var templ_writer = tmpl.NewOsTmpWriter()

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

func Create(p *placeholders.Placeholders, executor execute.Executor, t tmpl.TmplWriter) error {

	if err := tmpl.GenerateFS(templateFiles, goTemplateDir, p.Name, p, false, t); err != nil {
		return fmt.Errorf("generating structure from the template: %w", err)
	}

	if err := azure.EngageAzure(executor); err != nil {
		return fmt.Errorf("engaging Azure: %v", err)
	}

	return nil
}

func bindFlags(cmd *cobra.Command, args []string) error {
	if len(args) > 0 && args[0] == ProviderAzure {
		if err := setFlagConfig(cmd,
			FlagResourceGroup,
			placeholders.SetAzureResourceGroup,
			placeholders.GetAzureResourceGroup,
			placeholders.GetDefaultAzureResourceGroupName); err != nil {
			return err
		}

		if err := setFlagConfig(cmd,
			FlagLocation,
			placeholders.SetAzureLocation,
			placeholders.GetAzureLocation,
			placeholders.GetDefaultAzureLocation); err != nil {
			return err
		}

		return placeholders.SaveConfig()
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
	// silent, err := cmd.Flags().GetBool(config.FlagSilent)
	// if err != nil {
	// 	return err
	// }

	provider := args[0]

	fmt.Printf("PROVIDER CHECKS: %v", provider)

	if !executor.CommandExists("terraform") {
		return fmt.Errorf("terraform is not installed")
	}

	if !executor.CommandExists("docker") {
		return fmt.Errorf("docker is not installed")
	}

	switch provider {
	case ProviderAzure:
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
		ProjectName:    placeholders.GetProjectName(),
		ResourceGroup:  placeholders.GetAzureResourceGroup(),
		Location:       placeholders.GetAzureLocation(),
		StorageAccount: placeholders.GetAzureStorageAccount(),
		Keyvault:       placeholders.GetAzureKeyvault(),
		Acr:            placeholders.GetAzureAcr(),
		AppName:        placeholders.GetDefaultAzureApp(),
	}
	if err := tmpl.GenerateFS(templateFilesAz, goTemplateDir, placeholders.InfraDirectory(), p, false, templ_writer); err != nil {
		return fmt.Errorf("generating the project structure from the template: %s", err)
	}
	return nil
}
