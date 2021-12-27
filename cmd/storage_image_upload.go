package cmd

import "github.com/spf13/cobra"

var storageImageUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload to storage",
	Long:  "Upload to storage",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	storageImageCmd.AddCommand(storageImageUploadCmd)
}
