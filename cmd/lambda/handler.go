package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"poststone-chat/pkg/service"
	"poststone-chat/pkg/storage"
	"strings"
	"time"
)

func handleRequest(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ip := req.RequestContext.Identity.SourceIP
	if ip == "" {
		ip = "unknown"
	}

	userMessage := req.Body
	if strings.TrimSpace(userMessage) == "" {
	}
	if err := storage.AppendToS3Log(ip, fmt.Sprintf("%s %s", time.Now(), userMessage)); err != nil {
		log.Printf("Failed to log message %v", err)
	}

	answer, err := service.Answer(userMessage)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,x-requested-with",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "OPTIONS,POST,PATCH",
			},
			Body: "Internal error: " + err.Error(),
		}, nil
	}

	if err = storage.AppendToS3Log(ip, fmt.Sprintf("%s %s", time.Now(), answer)); err != nil {
		log.Printf("Failed to log message %v", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,x-requested-with",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "OPTIONS,POST,PATCH",
		},
		Body: answer,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
