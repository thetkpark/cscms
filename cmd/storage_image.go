package cmd

import "github.com/spf13/cobra"

var storageImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Upload image to storage",
	Long: `Upload new image to the storage. It can be access with "img.cscms.me". 
All the image is cached with CDN to ensure fast delivery across the globe`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	storageCmd.AddCommand(storageImageCmd)
}
