package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// storageApiKeyCmd represents the storageApiKey command
var storageApiKeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "Add API key for storage.cscms.me",
	Long:  `Add API key for storage.cscms.me`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("required api token argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := args[0]
		viper.Set("apiKey", apiKey)
		err := viper.WriteConfig()
		if err != nil {
			fmt.Printf("error: %s", err.Error())
		}
		fmt.Println("Your API Key is updated")
	},
}

func init() {
	storageCmd.AddCommand(storageApiKeyCmd)
}
