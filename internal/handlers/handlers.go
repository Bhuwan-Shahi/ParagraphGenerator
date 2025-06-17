package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bhuwan-Shahi/ParagraphGenerator/internal/models"
	"github.com/Bhuwan-Shahi/ParagraphGenerator/internal/services"
)

type Handler struct {
	generator *services.Generator
}

func NewHandler(generator *services.Generator) *Handler {
	return &Handler{generator: generator}
}

func (h *Handler) enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func (h *Handler) Generate(w http.ResponseWriter, r *http.Request) {
	h.enableCORS(w)

	if r.Method == "OPTIONS" {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req models.ParagraphRequest
	sentences := 5     // default
	style := "general" // default

	if r.Method == "POST" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
			return
		}
		if req.Sentences > 0 {
			sentences = req.Sentences
		}
		if req.Style != "" {
			style = req.Style
		}
	} else if r.Method == "GET" {
		if s := r.URL.Query().Get("sentences"); s != "" {
			if parsed, err := strconv.Atoi(s); err == nil && parsed > 0 {
				sentences = parsed
			}
		}
		if st := r.URL.Query().Get("style"); st != "" {
			style = st
		}
	}

	// Updated: Allow up to 1000 sentences instead of 50
	if sentences > 1000 {
		sentences = 1000
	}

	// Ensure minimum of 1 sentence
	if sentences < 1 {
		sentences = 1
	}

	paragraph := h.generator.GenerateParagraph(sentences, style)

	response := models.ParagraphResponse{
		Paragraph: paragraph,
		Sentences: sentences,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	h.enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	response := models.InfoResponse{
		Name:        "Paragraph Generator API",
		Version:     "1.0.0",
		Description: "Generate random paragraphs using various word lists",
		Endpoints: map[string]string{
			"/generate": "POST/GET - Generate paragraphs",
			"/info":     "GET - API information",
		},
		Parameters: map[string]string{
			"sentences": "Number of sentences (1-1000, default: 5)", // Updated documentation
			"style":     "Style of paragraph (general, formal, casual)",
		},
		Examples: map[string]string{
			"GET":  "/generate?sentences=3&style=formal",
			"POST": `{"sentences": 4, "style": "casual"}`,
		},
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	// Serve the HTML file from web/static/
	http.ServeFile(w, r, "web/static/index.html")
}

func (h *Handler) Start(port string) error {
	http.HandleFunc("/", h.Root)
	http.HandleFunc("/generate", h.Generate)
	http.HandleFunc("/info", h.Info)

	return http.ListenAndServe(":"+port, nil)
}
