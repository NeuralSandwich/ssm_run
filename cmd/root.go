package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	instanceID       string
	commandString    string
	commandID        string
	outputBucketName string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ssm_run",
	Short: "Allows you to run commands over AWS SSM",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&instanceID, "instance-id", "i", "", "The instance ID to run the command on")
	rootCmd.PersistentFlags().StringVarP(&commandString, "command-string", "c", "", "The command to run on the instance")
	rootCmd.PersistentFlags().StringVarP(&commandID, "command-id", "", "", "The SSM command ID to pull")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
