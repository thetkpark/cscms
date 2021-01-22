package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

type shortenURLRespBody struct {
	ShortURL string `json:"shortUrl"`
}

type shortenURLErrorBody struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

// createCmd represents the create command
var prefer string
var long bool
var akaCreateCmd = &cobra.Command{
	Use:   "create <url>",
	Args: cobra.ExactArgs(1),
	Short: "Create a shorten URL",
	Long: "Create a shorten URL from aka.cscms.me",
	Run: func(cmd *cobra.Command, args []string) {
		reqBodyMap := map[string]string{
			"url": args[0],
		}
		if len(prefer) > 0 {
			reqBodyMap["prefer"] = prefer
		}
		reqBodyJSON, err := json.Marshal(reqBodyMap)
		if err != nil {
			er("JSON body parsing error")
		}

		resp, err := http.Post("https://aka.cscms.me/api/newUrl", "application/json", bytes.NewBuffer(reqBodyJSON))
		if err != nil {
			er("Cannot connect to aka.cscms.me API server")
		}
		body, err := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == 201 {
			defer resp.Body.Close()
			var respBody shortenURLRespBody
			err = json.Unmarshal(body, &respBody)
			if err != nil {
				er("Cannot parse JSON response")
			}
			if long {
				fmt.Printf("Original URL: %s\n", args[0])
			}
			fmt.Printf("Shorten URL: https://aka.cscms.me/%s", respBody.ShortURL)
		} else if resp.StatusCode == 400 {
			var respErrorBody shortenURLErrorBody
			err = json.Unmarshal(body, &respErrorBody)
			if err != nil {
				er("Cannot parse JSON Error response")
			}
			er(respErrorBody.Error)
		} else {
			er("API Server Error")
		}
	},
}

func init() {
	akaCreateCmd.Flags().StringVarP(&prefer, "prefer", "p", "", "Set string that your perfer to use as a shorten URL")
	akaCreateCmd.Flags().BoolVarP(&long, "long", "l", false, "Output in more detail which includes shorten URL and original URL")
	akaCmd.AddCommand(akaCreateCmd)
}
