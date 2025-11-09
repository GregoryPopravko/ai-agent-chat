package main

import (
	_ "bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

type InputItem struct {
	ID            int     `json:"id"`
	Link          string  `json:"link"`
	Price         float32 `json:"price"`
	EmbeddingText string  `json:"embedding_text"`
}

type OutputItem struct {
	ID        int       `json:"id"`
	Link      string    `json:"link"`
	Price     float32   `json:"price"`
	Text      string    `json:"embedding_text"`
	Embedding []float32 `json:"embedding"`
}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Set OPENAI_API_KEY env variable.")
	}

	var inputs []InputItem
	data, err := os.ReadFile("/home/ryhor/IdeaProjects/custom/poststone-chat/cmd/embedder/input.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(data, &inputs); err != nil {
		log.Fatal(err)
	}

	client := resty.New()
	var results []OutputItem

	for _, item := range inputs {
		embedding, err := getEmbedding(client, apiKey, item.EmbeddingText)
		if err != nil {
			log.Printf("Error embedding ID %d: %v", item.ID, err)
			continue
		}
		results = append(results, OutputItem{
			ID:        item.ID,
			Link:      item.Link,
			Price:     item.Price,
			Text:      item.EmbeddingText,
			Embedding: embedding,
		})
	}

	outData, _ := json.MarshalIndent(results, "", "  ")
	if err := os.WriteFile("/home/ryhor/IdeaProjects/custom/poststone-chat/cmd/embedder/output.json", outData, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done! Saved to output.json")
}

func getEmbedding(client *resty.Client, apiKey string, text string) ([]float32, error) {
	type requestBody struct {
		Input string `json:"input"`
		Model string `json:"model"`
	}
	type responseBody struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}

	var respBody responseBody
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody{Input: text, Model: "text-embedding-3-small"}).
		SetResult(&respBody).
		Post("https://api.openai.com/v1/embeddings")

	if err != nil {
		print(resp)
		return nil, err
	}

	if len(respBody.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return respBody.Data[0].Embedding, nil
}
