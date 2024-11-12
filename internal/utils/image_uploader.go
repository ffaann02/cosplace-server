// utils/image_uploader.go
package utils

import (
	"context"
	// "crypto/md5"
	"encoding/base64"
	// "encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	imgBB "github.com/JohnNON/ImgBB"
)

// UploadImageToImgBB uploads a base64 image string to ImgBB and returns the hosted image URL.
func UploadImageToImgBB(userID string, base64Image string) (string, error) {
	apiKey := os.Getenv("IMGBB_API_KEY")
	if apiKey == "" {
		return "", errors.New("ImgBB API key is missing")
	}

	// Strip the "data:image/jpeg;base64," prefix if it exists
	if strings.HasPrefix(base64Image, "data:image/") {
		parts := strings.SplitN(base64Image, ",", 2)
		if len(parts) == 2 {
			base64Image = parts[1]
		} else {
			return "", errors.New("invalid base64 image format")
		}
	}

	// Validate the base64 string
	if _, err := base64.StdEncoding.DecodeString(base64Image); err != nil {
		return "", errors.New("invalid base64 string")
	}

	// Generate a hash-based name for the image
	timestamp := time.Now().Unix()
	imageHashName := fmt.Sprintf("%s-%d", userID, timestamp)

	// Create an ImgBB image instance
	img, err := imgBB.NewImage(imageHashName, 60, base64Image)
	if err != nil {
		return "", err
	}

	// Set up an HTTP client with a timeout
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Create an ImgBB client
	imgBBClient := imgBB.NewClient(httpClient, apiKey)

	// Upload the image to ImgBB
	resp, err := imgBBClient.Upload(context.Background(), img)
	if err != nil {
		return "", err
	}

	// Return the URL of the uploaded image
	return resp.Data.URL, nil
}

// hashSum creates an MD5 hash for naming or caching purposes
// func hashSum(b []byte) string {
// 	sum := md5.Sum(b)
// 	return hex.EncodeToString(sum[:])
// }
