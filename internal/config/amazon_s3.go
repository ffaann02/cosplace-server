package config

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var S3Session *session.Session

// InitAmazonS3 initializes the AWS S3 session
func InitAmazonS3() {
	// Get AWS configuration from environment variables
	region := os.Getenv("AWS_REGION")
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	fmt.Println("Initializing AWS S3 session...")
	fmt.Printf("Region: %s\n", region)

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	S3Session = sess
	fmt.Println("AWS S3 session initialized successfully")
}

func AmazonS3Storage() *session.Session {
	return S3Session
}
