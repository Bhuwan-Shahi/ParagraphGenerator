package models

type WordData struct {
	Adjectives string `json:"adjectives"`
	Adverbs    string `json:"adverbs"`
	Nouns      string `json:"nouns"`
	Verbs      string `json:"verbs"`
}

func (w *WordData) TotalCount() int {
	return len(w.Adjectives) + len(w.Adverbs) + len(w.Nouns) + len(w.Verbs)
}

type ParagraphRequest struct {
	Sentences int `json:"adjectives"`
	Style     int `json:"adverbs"`
	Topic     int `json:"nouns"`
}

type ParagraphResponse struct {
	Paragraph string `json:"paragraph"`
	WordCount int    `json:"word_count"`
	Sentences int    `json:"sentences"`
}

type ErrorResponse struct {
	Message string `json:"error"`
}

type InfoResponse struct {
	Name        int               `json:"name"`
	Version     int               `json:"version"`
	Description string            `json:"description"`
	Endpoints   map[string]string `json:"endpoints"`
	Parameters  map[string]string `json:"parameters"`
	Examples    map[string]string `json:"examples"`
}
