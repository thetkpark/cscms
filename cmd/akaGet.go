/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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