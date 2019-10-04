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
	RootCmd.AddCommand(didoperations.CreateControllerDIDCmd)
	RootCmd.AddCommand(didoperations.GetControllerDIDCmd)
	RootCmd.AddCommand(didoperations.UpdateControllerDIDCmd)
	RootCmd.AddCommand(didoperations.DeleteControllerDIDCmd)
	RootCmd.AddCommand(didoperations.CreateDeviceDIDCmd)
	RootCmd.AddCommand(didoperations.GetDeviceDIDCmd)
	RootCmd.AddCommand(didoperations.UpdateDeviceDIDCmd)
	RootCmd.AddCommand(didoperations.DeleteDeviceDIDCmd)
}
