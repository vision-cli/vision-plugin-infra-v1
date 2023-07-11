package plugin

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/vision-cli/vision/config"
	"github.com/vision-cli/vision/flag"
)

const (
	FlagResourceGroup  = "resource-group"
	FlagLocation       = "location"
	FlagStorageAccount = "storage-account"
	Command            = "infra"
)

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().AddFlagSet(flag.ConfigFlagset())
	createCmd.Flags().StringP(FlagResourceGroup, "p", "", "Azure resource group")
	createCmd.Flags().StringP(FlagLocation, "l", "", "Azure resource group location")
	createCmd.Flags().StringP(FlagStorageAccount, "s", "", "Azure storage account")
}

var RootCmd = &cobra.Command{
	Use:   Command,
	Short: "Manage project infra via CLI",
	Long:  "Manage configuration and defaults of the project's infra folder and cloud provider via a CLI",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeConfig(cmd)
	},
}

func initializeConfig(cmd *cobra.Command) error {
	var path string
	dir, err := os.Getwd()
	if err != nil {
		path = ""
	} else {
		path = filepath.Base(dir)
	}

	// load the project config file if it exists, otherwise prompt the user to create one
	return config.LoadConfig(cmd.Flags(), flag.IsSilent(cmd.Flags()), config.ConfigFilename, path)
}