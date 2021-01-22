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
	"github.com/pkg/browser"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)
func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
// akaGetCmd represents the akaGet command
var openBrowser bool
var browserOnly bool
var akaGetCmd = &cobra.Command{
	Use:   "get <shorten url>",
	Short: "Get the detail of URL from shorten URL",
	Long: `Get the detail of URL from shorten URL`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shortenURL := args[0]
		fullURLRegex := regexp.MustCompile("(.*aka.cscms.me/)(.+)")
		matchFullUrl := fullURLRegex.MatchString(args[0])
		if matchFullUrl {
			url := fullURLRegex.FindStringSubmatch(args[0])
			shortenURL = url[2]
		}
		resp, err := http.Get("https://aka.cscms.me/api/originalUrl?url=" + shortenURL)
		if err != nil {
			panic("API Error")
		}
		displayURL(resp)
	},
}

func init() {
	akaGetCmd.Flags().BoolVarP(&openBrowser,"open-browser", "b", false, "Open the full url using your default web browser")
	akaGetCmd.Flags().BoolVarP(&browserOnly,"open-only", "o", false, "Open the website of full URL in you default web browser without showing detail about the URL")
	akaCmd.AddCommand(akaGetCmd)
}


func displayURL (resp *http.Response) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Body reading error")
	}
	var fullUrlData fullUrlRespBody
	err = json.Unmarshal(body, &fullUrlData)
	if err != nil {
		panic("JSON reading error")
	}
	if browserOnly || openBrowser {
		browser.OpenURL(fullUrlData.URL)
	}
	if !browserOnly {
		fmt.Printf("Full URL: %s\n", fullUrlData.URL)
		fmt.Printf("Shorten URL: %s\n", "https://aka.cscms.me/"+fullUrlData.ShortenURL)
		fmt.Printf("Visited: %d\n", fullUrlData.Visit)
	}

}