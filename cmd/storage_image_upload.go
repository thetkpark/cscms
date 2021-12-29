package cmd

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"time"
)

type Image struct {
	ID               uint      `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	OriginalFilename string    `json:"original_filename"`
	FileSize         uint64    `json:"file_size"`
	FilePath         string    `json:"file_path"`
	UserID           uint      `json:"user_id"`
}

var storageImageUploadCmd = &cobra.Command{
	Use:       "upload",
	Short:     "Upload image to storage",
	Long:      "Upload image to storage",
	ValidArgs: []string{"image file"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("select one image file to upload")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Check api key
		apiKey := viper.GetString("apiKey")
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Println("Error opening the file")
			return
		}
		defer file.Close()

		// Check if fileInfo is valid
		fileInfo, err := file.Stat()
		if err != nil {
			return
		}
		// Check file size
		if fileInfo.Size() > 5<<20 {
			fmt.Println("the limited image size is 5MB")
			return
		}

		// Check file content type
		mimeType, err := getFileContentType(file)
		if err != nil {
			return
		}
		if !isSupportedMimeType(mimeType) {
			fmt.Printf("\n%s is not supported", mimeType)
			return
		}

		// Sending HTTP request to api server
		result := Image{}
		client := resty.New()
		res, err := client.R().SetFile("image", args[0]).
			SetHeader("x-api-key", apiKey).
			SetResult(&result).
			Post("https://storage.cscms.me/api/image")
		if err != nil {
			fmt.Println("Error uploading image to server. Please check your internet connection")
			return
		}

		if res.IsError() {
			fmt.Printf("\nThere is an error with code %d", res.StatusCode())
			return
		}

		fmt.Printf("https://img.cscms.me/%s", result.FilePath)
	},
}

func init() {
	storageImageCmd.AddCommand(storageImageUploadCmd)
}

func getFileContentType(out *os.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	if _, err := out.Read(buffer); err != nil {
		return "", nil
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}

func isSupportedMimeType(mimeType string) bool {
	supportedType := [11]string{"image/png", "image/jpeg", "image/gif", "image/x-icon", "image/heic", "image/webp", "image/tiff", "image/svg+xml", "image/bmp", "image/apng", "image/avif"}
	for _, ty := range supportedType {
		if mimeType == ty {
			return true
		}
	}
	return false
}
