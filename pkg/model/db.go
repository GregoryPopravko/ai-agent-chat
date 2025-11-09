package model

type Point struct {
	ID      int                    `json:"id"`
	Vector  []float32              `json:"vector"`
	Payload map[string]interface{} `json:"payload"`
}

type UploadRequest struct {
	Points []Point `json:"points"`
}

type SearchRequest struct {
	Vector      []float32 `json:"vector"`
	Top         int       `json:"top"`
	WithPayload bool      `json:"with_payload"`
}

type SearchResult struct {
	Result []struct {
		ID      int                    `json:"id"`
		Score   float64                `json:"score"`
		Payload map[string]interface{} `json:"payload"`
	} `json:"result"`
}

type Monument struct {
	ID          int                    `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Price       float64                `json:"price"`
	Link        string                 `json:"url"`
	Payload     map[string]interface{} `json:"payload"`
}
