package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/davyj0nes/ssm_run/awsCmds"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// setup for the powershellCmd
func init() {
	rootCmd.AddCommand(runCmd)
}

// powershellCmd represents the powershell command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a powershell command",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runCommand(cmd, args, "powershell"); err != nil {
			fmt.Fprintf(os.Stderr, "Error running powershell command\n%v\n", err)
			os.Exit(1)
		}
	},
}

// runCommand does the heavy lifting for the command
func runCommand(cmd *cobra.Command, args []string, callType string) error {
	if instanceID == "" {
		return errors.New("instance-id cannot be empty")
	}

	cmdID, err := awsCmds.SSMSendCmd(true, instanceID, commandString)
	if err != nil {
		return errors.Wrap(err, "Problem making SSM Command Request")
	}

	time.Sleep(1 * time.Second)
	cmdInfo, err := awsCmds.SSMGetCmd(instanceID, cmdID)
	if err != nil {
		return errors.Wrap(err, "Problem initially getting SSM Command Invocation")
	}

	var loopCounter int

	for {
		if cmdInfo.CommandStatus == "Success" {
			break
		} else if cmdInfo.CommandStatus == "Failed" {
			fmt.Println("Command Failed")
			break
			// return errors.New("Get Command Invocation Failed")
		}

		if loopCounter == 150 {
			fmt.Println("Request taking too long to process")
			fmt.Println("Status:", cmdInfo.CommandStatus)
			return errors.New("get command invocation timed out")
		}

		// have to create a new request each time
		cmdInfo, err = awsCmds.SSMGetCmd(instanceID, cmdID)
		if err != nil {
			return errors.Wrap(err, "Problem getting SSM Command Invocation in loop")
		}
		// user output to see how many checks are needed
		fmt.Fprint(os.Stderr, ".")
		loopCounter++
		// wait 2 seconds between each check
		time.Sleep(2 * time.Second)
	}

	fmt.Println("\nCommand ID:\t", cmdInfo.CommandID)
	fmt.Println("Instance ID:\t", cmdInfo.InstanceID)
	fmt.Println("Command Status:\t", cmdInfo.CommandStatus)
	fmt.Printf("----Output----\nSTDOUT:\n%s\nSTDERR:\n%s\n", cmdInfo.StdOut, cmdInfo.StdErr)
	return nil
}
