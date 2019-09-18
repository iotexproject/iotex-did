package cmd

import (
	"github.com/spf13/cobra"

	"github.com/iotexproject/iotex-DID/cmd/didoperations"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "did",
	Short: "Command-line interface for did operations",
	Long:  "Command-line interface for did operations",
}

func init() {
	RootCmd.AddCommand(didoperations.CreateDIDCmd)
	RootCmd.AddCommand(didoperations.GetDIDCmd)
	RootCmd.AddCommand(didoperations.UpdateDIDCmd)
	RootCmd.AddCommand(didoperations.DeleteDIDCmd)
}
