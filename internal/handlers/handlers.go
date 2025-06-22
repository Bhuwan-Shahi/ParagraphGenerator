package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Bhuwan-Shahi/ParagraphGenerator/internal/models"
	"github.com/Bhuwan-Shahi/ParagraphGenerator/internal/services"
)

type Handler struct {
	generator           *services.Generator
	rand                *rand.Rand
	programmingQuotes   []string
	inspirationalQuotes []string
}

func NewHandler(generator *services.Generator) *Handler {
	h := &Handler{
		generator: generator,
		rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	// Load quotes from files
	h.loadQuotes()

	return h
}

// loadQuotes reads quotes from text files
func (h *Handler) loadQuotes() {
	// Load programming quotes
	if quotes, err := h.readQuotesFromFile("data/programming.txt"); err != nil {
		fmt.Printf("Warning: Could not load programming quotes: %v\n", err)
		// Fallback to a few hardcoded quotes
		h.programmingQuotes = []string{
			"The only way to do great work is to love what you do.",
			"Code is like humor. When you have to explain it, it's bad.",
			"Programming isn't about what you know; it's about what you can figure out.",
		}
	} else {
		h.programmingQuotes = quotes
		fmt.Printf("Loaded %d programming quotes\n", len(h.programmingQuotes))
	}

	// Load inspirational quotes
	if quotes, err := h.readQuotesFromFile("data/quotes.txt"); err != nil {
		fmt.Printf("Warning: Could not load inspirational quotes: %v\n", err)
		// Fallback to a few hardcoded quotes
		h.inspirationalQuotes = []string{
			"The only way to do great work is to love what you do.",
			"Success is not final, failure is not fatal: it is the courage to continue that counts.",
			"Believe you can and you're halfway there.",
		}
	} else {
		h.inspirationalQuotes = quotes
		fmt.Printf("Loaded %d inspirational quotes\n", len(h.inspirationalQuotes))
	}
}

// readQuotesFromFile reads quotes from a text file, one quote per line
func (h *Handler) readQuotesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %w", filename, err)
	}
	defer file.Close()

	var quotes []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments (lines starting with #)
		if line != "" && !strings.HasPrefix(line, "#") {
			quotes = append(quotes, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filename, err)
	}

	if len(quotes) == 0 {
		return nil, fmt.Errorf("no quotes found in file %s", filename)
	}

	return quotes, nil
}

func (h *Handler) enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// Programming endpoint - returns programming-related quotes
func (h *Handler) Programming(w http.ResponseWriter, r *http.Request) {
	h.enableCORS(w)

	if r.Method == "OPTIONS" {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	count := 1
	if c := r.URL.Query().Get("count"); c != "" {
		if parsed, err := strconv.Atoi(c); err == nil && parsed > 0 {
			count = parsed
		}
	}

	// Limit count to reasonable number
	if count > 20 {
		count = 20
	}

	var selectedQuotes []string
	if len(h.programmingQuotes) == 0 {
		http.Error(w, `{"error": "No programming quotes available"}`, http.StatusInternalServerError)
		return
	}

	if count >= len(h.programmingQuotes) {
		selectedQuotes = h.programmingQuotes
	} else {
		// Get random quotes without repetition
		indices := h.rand.Perm(len(h.programmingQuotes))[:count]
		for _, idx := range indices {
			selectedQuotes = append(selectedQuotes, h.programmingQuotes[idx])
		}
	}

	response := map[string]interface{}{
		"category":        "programming",
		"count":           len(selectedQuotes),
		"total_available": len(h.programmingQuotes),
		"quotes":          selectedQuotes,
	}

	json.NewEncoder(w).Encode(response)
}

// Quotes endpoint - returns general inspirational quotes
func (h *Handler) Quotes(w http.ResponseWriter, r *http.Request) {
	h.enableCORS(w)

	if r.Method == "OPTIONS" {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get count parameter (default: 1)
	count := 1
	if c := r.URL.Query().Get("count"); c != "" {
		if parsed, err := strconv.Atoi(c); err == nil && parsed > 0 {
			count = parsed
		}
	}

	// Limit count to reasonable number
	if count > 20 {
		count = 20
	}

	var selectedQuotes []string
	if len(h.inspirationalQuotes) == 0 {
		http.Error(w, `{"error": "No inspirational quotes available"}`, http.StatusInternalServerError)
		return
	}

	if count >= len(h.inspirationalQuotes) {
		selectedQuotes = h.inspirationalQuotes
	} else {
		// Get random quotes without repetition
		indices := h.rand.Perm(len(h.inspirationalQuotes))[:count]
		for _, idx := range indices {
			selectedQuotes = append(selectedQuotes, h.inspirationalQuotes[idx])
		}
	}

	response := map[string]interface{}{
		"category":        "inspirational",
		"count":           len(selectedQuotes),
		"total_available": len(h.inspirationalQuotes),
		"quotes":          selectedQuotes,
	}

	json.NewEncoder(w).Encode(response)
}

// Your existing methods remain unchanged
func (h *Handler) Generate(w http.ResponseWriter, r *http.Request) {
	h.enableCORS(w)

	if r.Method == "OPTIONS" {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req models.ParagraphRequest
	sentences := 20    // default
	style := "general" // default

	switch r.Method {
	case "POST":
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
			return
		}
		if req.Sentences > 0 {
			sentences = req.Sentences
		}
		if req.Style != 0 {
			style = strconv.Itoa(req.Style)
		}
	case "GET":
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

	// Count actual words in the generated paragraph
	wordCount := len(strings.Fields(paragraph))

	response := models.ParagraphResponse{
		Paragraph: paragraph,
		Sentences: sentences,
		WordCount: wordCount, // Make sure this field exists in your model
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	h.enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	response := models.InfoResponse{
		Name:        "Enhanced Paragraph Generator API",
		Version:     "1.1.0",
		Description: "Generate random paragraphs and get inspirational quotes",
		Endpoints: map[string]string{
			"/generate":    "POST/GET - Generate paragraphs",
			"/programming": "GET - Get programming quotes",
			"/quotes":      "GET - Get inspirational quotes",
			"/health":      "GET - Health check",
			"/info":        "GET - API information",
		},
		Parameters: map[string]string{
			"sentences": "Number of sentences (1-1000, default: 5)",
			"style":     "Style of paragraph (general, formal, casual)",
			"count":     "Number of quotes to return (1-20, default: 1)",
		},
		Examples: map[string]string{
			"Generate":    "/generate?sentences=3&style=formal",
			"Programming": "/programming?count=3",
			"Quotes":      "/quotes?count=5",
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
	http.HandleFunc("/programming", h.Programming)
	http.HandleFunc("/quotes", h.Quotes)
	http.HandleFunc("/info", h.Info)

	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET  /              - Web interface")
	fmt.Println("  GET  /generate      - Generate paragraphs")
	fmt.Printf("  GET  /programming   - Get programming quotes (%d loaded)\n", len(h.programmingQuotes))
	fmt.Printf("  GET  /quotes        - Get inspirational quotes (%d loaded)\n", len(h.inspirationalQuotes))
	fmt.Println("  GET  /info          - API information")

	return http.ListenAndServe(":"+port, nil)
}
