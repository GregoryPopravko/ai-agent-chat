package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"strings"
	"time"
)

const (
	BucketName = "poststone-chat-logs-349447954185-eu-central-1"
	LogPrefix  = "logs"
)

var s3Client *s3.Client

func AppendToS3Log(ip string, entry string) error {
	today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("%s/%s/%s.txt", LogPrefix, today, sanitizeFileName(ip))

	var buffer bytes.Buffer

	existing, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(key),
	})

	if err == nil {
		defer existing.Body.Close()
		if _, err := buffer.ReadFrom(existing.Body); err != nil {
			log.Printf("Failed to read existing log file content: %v", err)
		}
	} else {
		if !strings.Contains(err.Error(), "NoSuchKey") {
			log.Printf("Failed to get object from S3: %v", err)
			return err
		}
		log.Printf("Log file does not exist yet, creating a new one for IP %s", ip)
	}

	buffer.WriteString(entry + "\n")

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buffer.Bytes()),
	})
	if err != nil {
		log.Printf("Failed to write log to S3: %v", err)
		return err
	}

	return nil
}

func sanitizeFileName(ip string) string {
	return strings.ReplaceAll(ip, ":", "_")
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config: %v", err)
	}
	s3Client = s3.NewFromConfig(cfg)
}
