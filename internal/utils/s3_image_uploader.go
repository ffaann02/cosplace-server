package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	config "github.com/ffaann02/cosplace-server/internal/config"
)

// UploadImageToAmazonS3 uploads an image to Amazon S3
// base64Image: the base64-encoded image string
// prefix: use as folder name
// userID: use as a part of file name
func UploadImageToAmazonS3(base64Image string, prefix string, userID string) (string, error) {

	S3Session := config.AmazonS3Storage()
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")

	// Check if the S3 session is initialized
	if S3Session == nil {
		return "", errors.New("AWS S3 session not initialized")
	}

	// Split the base64 string to extract the mime type and data
	parts := strings.Split(base64Image, ",")
	if len(parts) != 2 {
		return "", errors.New("invalid base64 image format")
	}

	// Decode the base64 data
	imageData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 image: %v", err)
	}

	// Get the mime type and file extension
	mimeType := strings.Split(parts[0], ";")[0]
	mimeType = strings.TrimPrefix(mimeType, "data:")
	exts, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(exts) == 0 {
		return "", fmt.Errorf("unable to determine file extension for mime type: %s", mimeType)
	}

	// Generate a hash-based name for the folder and filename
	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("%s/%s-%d", prefix, userID, timestamp)

	// Create the S3 service client
	svc := s3.New(S3Session)

	// Upload the file
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(fileName),
		Body:          bytes.NewReader(imageData),
		ContentType:   aws.String(mimeType),
		ContentLength: aws.Int64(int64(len(imageData))),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %v", err)
	}

	// Generate the file URL
	baseURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com", bucketName, os.Getenv("AWS_REGION"))
	fileURL := fmt.Sprintf("%s/%s", baseURL, fileName)

	return fileURL, nil
}

func UploadImageToAmazonS3WithRunningNumber(base64Image string, prefix string, userID string) (string, error) {
	S3Session := config.AmazonS3Storage()
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")

	if S3Session == nil {
		return "", errors.New("AWS S3 session not initialized")
	}

	parts := strings.Split(base64Image, ",")
	if len(parts) != 2 {
		return "", errors.New("invalid base64 image format")
	}

	imageData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 image: %v", err)
	}

	mimeType := strings.Split(parts[0], ";")[0]
	mimeType = strings.TrimPrefix(mimeType, "data:")
	exts, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(exts) == 0 {
		return "", fmt.Errorf("unable to determine file extension for mime type: %s", mimeType)
	}

	// Generate the initial file name
	timestamp := time.Now().Unix()
	baseFileName := fmt.Sprintf("%s/%s-%d", prefix, userID, timestamp)
	fileName := baseFileName
	runningNumber := 0

	svc := s3.New(S3Session)

	for {
		// Use ListObjectsV2 to check if a file with the same prefix exists
		listOutput, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket:  aws.String(bucketName),
			Prefix:  aws.String(fileName),
			MaxKeys: aws.Int64(1),
		})

		if err != nil {
			return "", fmt.Errorf("failed to list objects in S3: %v", err)
		}

		if len(listOutput.Contents) == 0 {
			// No matching file found, proceed with this file name
			fmt.Printf("No matching file found, using file name: %s\n", fileName)
			break
		}

		// Update file name with running number if a match is found
		runningNumber++
		fileName = fmt.Sprintf("%s-%d", baseFileName, runningNumber)
		fmt.Printf("File name already exists, trying: %s\n", fileName)
	}

	// Upload the file
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(fileName),
		Body:          bytes.NewReader(imageData),
		ContentType:   aws.String(mimeType),
		ContentLength: aws.Int64(int64(len(imageData))),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %v", err)
	}

	// Generate the file URL
	baseURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com", bucketName, os.Getenv("AWS_REGION"))
	fileURL := fmt.Sprintf("%s/%s", baseURL, fileName)

	return fileURL, nil
}
