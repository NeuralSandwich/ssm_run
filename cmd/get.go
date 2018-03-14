package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/davyj0nes/ssm_run/awsCmds"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

// powershellCmd represents the powershell command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get Output from SSM command",
	Run: func(cmd *cobra.Command, args []string) {
		if err := getCommand(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error running powershell command\n%v\n", err)
			os.Exit(1)
		}
	},
}

// getCommand does the heavy lifting for the command
func getCommand(cmd *cobra.Command, args []string) error {
	cmdInfo, err := awsCmds.SSMGetCmd(instanceID, commandID)
	if err != nil {
		return errors.Wrap(err, "Problem getting SSM Cowmmand Info")
	}

	fmt.Println("\nCommand ID:\t", cmdInfo.CommandID)
	fmt.Println("Instance ID:\t", cmdInfo.InstanceID)
	fmt.Println("Command Status:\t", cmdInfo.CommandStatus)
	fmt.Printf("----Output----\nSTDOUT:\n%s\nSTDERR:\n%s\n", cmdInfo.StdOut, cmdInfo.StdErr)

	return nil
}
