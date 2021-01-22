/*
Copyright © 2021 Sethanant Pipatpakorn <sethanant.p@icloud.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
)

type shortenURLRespBody struct {
	ShortUrl string `json:"shortUrl"`
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
		reqBodyJson, err := json.Marshal(reqBodyMap)
		if err != nil {
			er("JSON body parsing error")
		}

		resp, err := http.Post("https://aka.cscms.me/api/newUrl", "application/json", bytes.NewBuffer(reqBodyJson))
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
			fmt.Printf("Shorten URL: https://aka.cscms.me/%s", respBody.ShortUrl)
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
