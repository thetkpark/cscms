package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/pkg/browser"

	"github.com/spf13/cobra"
)

// akaGetCmd represents the akaGet command
var openBrowser, browserOnly bool

var akaGetCmd = &cobra.Command{
	Use:   "get <shorten url>",
	Short: "Get the original URL from shorten URL",
	Long: `Get the original of URL from shorten URL`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shortenURL := args[0]
		fullURLRegex := regexp.MustCompile("(.*aka.cscms.me/)(.+)")
		matchFullURL := fullURLRegex.MatchString(args[0])
		if matchFullURL {
			url := fullURLRegex.FindStringSubmatch(args[0])
			shortenURL = url[2]
		}
		resp, err := http.Get("https://aka.cscms.me/api/originalUrl?url=" + shortenURL)
		if err != nil {
			er("Cannot connect to aka.cscms.me. Please check your internet connection")
		}
		if resp.StatusCode == 200 {
			displayURL(resp)
		} else {
			switch resp.StatusCode {
			case 400:
				er("The URL input is not valid")
			case 404:
				er("Original URL not found.")
			case 500:
				er("API Server Error")
			}
		}
	},
}

func init() {
	akaGetCmd.Flags().BoolVarP(&openBrowser,"open-browser", "b", false, "Open the full url using your default web browser")
	akaGetCmd.Flags().BoolVarP(&browserOnly,"open-only", "o", false, "Open the website of full URL in you default web browser without showing detail about the URL")
	akaCmd.AddCommand(akaGetCmd)
}

type fullURLRespBody struct {
	ShortenURL string `json:"shortenUrl"`
	URL        string `json:"url"`
	Visit      int64  `json:"visit"`
}

func displayURL (resp *http.Response) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		er("Cannot read response body from API calling")
	}
	var fullURLData fullURLRespBody
	err = json.Unmarshal(body, &fullURLData)
	if err != nil {
		er("Cannot parse JSON into struct")
	}
	if browserOnly || openBrowser {
		browser.OpenURL(fullURLData.URL)
	}
	if !browserOnly {
		fmt.Printf("Full URL: %s\n", fullURLData.URL)
		fmt.Printf("Shorten URL: %s\n", "https://aka.cscms.me/"+fullURLData.ShortenURL)
		fmt.Printf("Visited: %d\n", fullURLData.Visit)
	}

}