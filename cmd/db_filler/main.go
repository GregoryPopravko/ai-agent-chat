package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"poststone-chat/pkg/model"
	"time"
)

func main() {
	data, err := os.ReadFile("/home/ryhor/IdeaProjects/custom/poststone-chat/cmd/db_filler/data.json")
	if err != nil {
		log.Fatalf("failed to read JSON file: %v", err)
	}

	var rawItems []map[string]interface{}
	if err := json.Unmarshal(data, &rawItems); err != nil {
		log.Fatalf("failed to parse JSON: %v", err)
	}

	var points []model.Point
	for _, item := range rawItems {
		idFloat := item["id"].(float64)
		vec := convertToFloat32Slice(item["embedding"])
		delete(item, "embedding")

		points = append(points, model.Point{
			ID:      int(idFloat),
			Vector:  vec,
			Payload: item,
		})
	}

	uploadReq := model.UploadRequest{Points: points}
	body, _ := json.Marshal(uploadReq)

	url := "https://7405179c-d3a7-4d7e-9636-79fbb58db450.eu-central-1-0.aws.cloud.qdrant.io:6333/collections/gravestones/points?wait=true"
	req, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3MiOiJtIn0.wQYVrdcKwdfG_rpQ11JwcZcS6WGPtlUq575AlEOVvbQ")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %s\nResponse: %s\n", resp.Status, string(respBody))
}

func convertToFloat32Slice(input interface{}) []float32 {
	list, ok := input.([]interface{})
	if !ok {
		return nil
	}
	result := make([]float32, len(list))
	for i, v := range list {
		if f, ok := v.(float64); ok {
			result[i] = float32(f)
		}
	}
	return result
}
