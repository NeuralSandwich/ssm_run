package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of the instance",
	Long: `This checks both the EC2 system and instance checks
	as well as SSM to see if the instance can be controlled via SSM.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("to be implemented")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolP("instance-id", "i", false, "The instance ID to check")
}
