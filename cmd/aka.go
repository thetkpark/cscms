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
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"regexp"
)

type fullUrlRespBody struct {
	ShortenURL string `json:"shortenUrl"`
	URL        string `json:"url"`
	Visit      int64  `json:"visit"`
}

// akaCmd represents the aka command
var akaCmd = &cobra.Command{
	Use:   "aka",
	Short: "Interact with URL shortener service at aka.cscms.me",
	Long: `Interact with URL shortener service at aka.cscms.me`,
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "get" {
			shortenURL := args[1]
			fullURLRegex := regexp.MustCompile("(.*aka.cscms.me/)(.+)")
			matchFullUrl := fullURLRegex.MatchString(args[1])
			if matchFullUrl {
				url := fullURLRegex.FindStringSubmatch(args[1])
				shortenURL = url[2]
			}
			resp, err := http.Get("https://aka.cscms.me/api/originalUrl?url=" + shortenURL)
			if err != nil {
				panic("API Error")
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic("Body reading error")
			}
			var fullUrlData fullUrlRespBody
			json.Unmarshal(body, &fullUrlData)
			fmt.Printf("Full URL: %s\n", fullUrlData.URL)
			fmt.Printf("Shorten URL: %s\n", "https://aka.cscms.me/"+fullUrlData.ShortenURL)
			fmt.Printf("Visited: %d\n", fullUrlData.Visit)
			//browser.OpenURL(fullUrlData.URL)
		}
	},
}

func init() {
	rootCmd.AddCommand(akaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// akaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// akaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
