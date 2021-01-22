package cmd

import (
	"github.com/spf13/cobra"
)

// akaCmd represents the aka command
var akaCmd = &cobra.Command{
	Use:   "aka",
	Short: "Interact with URL shortener service at aka.cscms.me",
	Long:  `Interact with URL shortener service at aka.cscms.me`,
	//Run: func(cmd *cobra.Command, args []string) {
	//},
}

func init() {
	rootCmd.AddCommand(akaCmd)

}
