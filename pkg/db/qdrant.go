package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"poststone-chat/pkg/model"
)

func SearchByVector(vector []float32) []model.Monument {
	qdrantURL := os.Getenv("QDRANT_URL")
	collection := os.Getenv("QDRANT_COLLECTION")
	key := os.Getenv("QDRANT_API_KEY")

	data, _ := json.Marshal(model.SearchRequest{
		Vector:      vector,
		Top:         3,
		WithPayload: true,
	})
	url := fmt.Sprintf("%s/collections/%s/points/search", qdrantURL, collection)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("api-key", key)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result model.SearchResult
	_ = json.Unmarshal(body, &result)

	var monuments []model.Monument
	for _, r := range result.Result {
		price := r.Payload["price"].(float64)
		link := r.Payload["link"].(string)
		title := r.Payload["Title"].(string)
		description := r.Payload["Description"].(string)
		monuments = append(monuments, model.Monument{
			ID:          r.ID,
			Title:       title,
			Description: description,
			Price:       price,
			Link:        link,
			Payload:     r.Payload})
	}

	return monuments
}

func Cheapest() []model.Monument {
	return []model.Monument{
		{
			Title: "Небольшой памятник из гранита",
			Description: "Памятник в небольшую цену из настоящего Карельского габбро-диабаза с размерами стелы 80х40 сантиметров. " +
				"Идеально подходит для гравировки портрета.Так же, памятник может быть выполнен из гранита другого цвета, что не повлияет на его цену.",
			Price: 760,
			Link:  "https://poststone.by/home/67-nebolshoj-pamyatnik-iz-granita.html",
		},
		{
			Title: "Памятник «Знак веры»",
			Description: "Строгий памятник из чёрного габбро-диабаза с выступающим крестом и аккуратными веерами у основания. " +
				"Символ веры, уважения и памяти. Подходит для тех, кто ищет сдержанный, но выразительный знак прощания",
			Price: 1100,
			Link:  "https://poststone.by/home/195-pamyatnik-znak-very.html",
		},
		{
			Title: "Памятник «Спокойствие»",
			Description: "Памятник с мягкими плавными линиями и резным основанием в виде вееров. Изготовлен из натурального габбро-диабаза. " +
				"Универсальный вариант, передающий светлую память без излишеств.",
			Price: 1100,
			Link:  "https://poststone.by/home/196-pamyatnik-spokoystvie.html",
		},
	}
}

func MostExpensive() []model.Monument {
	return []model.Monument{
		{
			Title: "Памятник по индивидуальному проекту из карьерного камня и валуна",
			Description: "Мощный памятник из габродиабаза: элегантный, брутальный, с бронзовыми буквами. " +
				"Природный валун, слегка обработанный — живая скала в граните.",
			Price: 18000,
			Link:  "https://poststone.by/home/150-kompleks-iz-karernogo-kamnya-i-valuna.html",
		},
		{
			Title: "Памятник по индивидуальному проекту из натурального валуна",
			Description: "Выдающийся образец надгробной пластики: благородный природный камень с живой текстурой, " +
				"сочетающий строгость форм и необузданную дикость материи.",
			Price: 15000,
			Link:  "https://poststone.by/home/149-pamyatnik-iz-naturalnogo-valuna.html",
		},
		{
			Title: "Памятник из гранита «Слава Веков»",
			Description: "Памятник «Слава Веков» с классической формой и двумя отдельными надгробными плитами. " +
				"Идеален для двойного или большего захоронения, воплощая монументальность и консервативные традиции.",
			Price: 11280,
			Link:  "https://poststone.by/home/36-pamyatnik-slava-vekov.html",
		},
	}
}
