package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"poststone-chat/pkg/model"
	"poststone-chat/pkg/texts"
	"strings"
)

const (
	IntentCheap     = "CHEAPEST"
	IntentExpensive = "MOST_EXPENSIVE"
	IntentGeneral   = "GENERAL"
	IntentUnknown   = "UNKNOWN"
	IntentSearch    = "SEARCH"
)

func InterpretQuery(query string) (string, string) {
	reqBody := model.OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []model.Message{{
			Role:    "system",
			Content: texts.PromptContent,
		}, {
			Role:    "user",
			Content: query,
		}},
	}

	data, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var parsed model.OpenAIResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		log.Fatal(err)
	}

	if len(parsed.Choices) > 0 {
		response := strings.TrimSpace(parsed.Choices[0].Message.Content)
		if response == IntentCheap {
			return IntentCheap, ""
		} else if response == IntentExpensive {
			return IntentExpensive, ""
		} else if strings.HasPrefix(response, IntentGeneral+":") {
			return IntentGeneral, strings.TrimSpace(strings.TrimPrefix(response, IntentGeneral+":"))
		} else if strings.HasPrefix(response, IntentSearch+":") {
			return IntentSearch, strings.TrimSpace(strings.TrimPrefix(response, IntentSearch+":"))
		} else {
			return IntentUnknown, ""
		}
	}

	return IntentUnknown, ""
}

func EmbedUserPhrase(phrase string) ([]float32, error) {
	reqBody := map[string]interface{}{
		"input": phrase,
		"model": "text-embedding-3-small",
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result model.OpenAIEmbeddingResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, nil
	}

	return result.Data[0].Embedding, nil
}

func PresentOptionsToUser(query string, monuments []model.Monument) (string, error) {
	content := "Внимание: Ты должен обязательно прикрепить к каждому памятнику его ссылку из поля 'Ссылка'. " +
		"Нельзя пропускать ссылки. Придумай интересную подачу каждого варианта.\n\n"
	content += fmt.Sprintf("Запрос пользователя: %s\n\nПодходящие памятники:\n", query)
	for i, m := range monuments {
		content += fmt.Sprintf("%d. %s\nОписание: %s\nЦена: %.2f BYN\nСсылка: %s\n\n", i+1, m.Title, m.Description, m.Price, m.Link)
	}

	reqBody := model.OpenAIRequest{
		Model: "gpt-4-turbo",
		Messages: []model.Message{{
			Role:    "system",
			Content: texts.PresentPrompt,
		}, {
			Role:    "user",
			Content: content,
		}},
	}

	data, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var parsed model.OpenAIResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", err
	}

	if len(parsed.Choices) > 0 {
		return parsed.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no choices returned")
}

func AnswerGeneralQuestion(userQuestion string) (string, error) {
	reqBody := model.OpenAIRequest{
		Model: "gpt-4-turbo",
		Messages: []model.Message{{
			Role:    "system",
			Content: texts.GeneralQuestionPrompt,
		}, {
			Role:    "user",
			Content: userQuestion,
		}},
	}

	data, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var parsed model.OpenAIResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", err
	}

	if len(parsed.Choices) > 0 {
		return parsed.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no choices returned")
}
