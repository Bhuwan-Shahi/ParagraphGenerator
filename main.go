package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type WordData struct {
	Adjectives []string `json:"adjectives"`
	Adverbs    []string `json:"adverbs"`
	Nouns      []string `json:"nouns"`
	Verbs      []string `json:"verbs"`
}

type ParagraphRequest struct {
	Sentences int    `json:"sentences,omitempty"`
	Style     string `json:"style,omitempty"`
	Topic     string `json:"topic,omitempty"`
}

type ParagraphResponse struct {
	Paragraph string `json:"paragraph"`
	WordCount int    `json:"word_count"`
	Sentences int    `json:"sentences"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var wordData WordData

// Load words from text files
func loadWordsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			// Split by comma and clean each word
			parts := strings.Split(line, ",")
			for _, part := range parts {
				word := strings.TrimSpace(part)
				if word != "" {
					words = append(words, word)
				}
			}
		}
	}
	return words, scanner.Err()
}

// Initialize word data from files
func initWordData() error {
	var err error

	wordData.Adjectives, err = loadWordsFromFile("data/adjectives.txt")
	if err != nil {
		return fmt.Errorf("failed to load adjectives: %v", err)
	}

	wordData.Adverbs, err = loadWordsFromFile("data/adverbs.txt")
	if err != nil {
		return fmt.Errorf("failed to load adverbs: %v", err)
	}

	wordData.Nouns, err = loadWordsFromFile("data/nouns.txt")
	if err != nil {
		return fmt.Errorf("failed to load nouns: %v", err)
	}

	wordData.Verbs, err = loadWordsFromFile("data/verbs.txt")
	if err != nil {
		return fmt.Errorf("failed to load verbs: %v", err)
	}

	log.Printf("Loaded %d adjectives, %d adverbs, %d nouns, %d verbs",
		len(wordData.Adjectives), len(wordData.Adverbs),
		len(wordData.Nouns), len(wordData.Verbs))

	return nil
}

// Generate a random sentence
func generateSentence(style string) string {
	patterns := []string{
		"The %s %s %s %s.",
		"A %s %s %s the %s.",
		"%s %s %s %s.",
		"The %s %s was %s and %s.",
		"Every %s %s %s %s.",
		"This %s %s %s %s %s.",
		"Many %s %s %s %s.",
		"The %s and %s %s %s %s.",
	}

	pattern := patterns[rand.Intn(len(patterns))]

	// Count placeholders
	placeholders := strings.Count(pattern, "%s")
	words := make([]interface{}, placeholders)

	for i := 0; i < placeholders; i++ {
		switch rand.Intn(4) {
		case 0:
			words[i] = wordData.Adjectives[rand.Intn(len(wordData.Adjectives))]
		case 1:
			words[i] = wordData.Adverbs[rand.Intn(len(wordData.Adverbs))]
		case 2:
			words[i] = wordData.Nouns[rand.Intn(len(wordData.Nouns))]
		case 3:
			words[i] = wordData.Verbs[rand.Intn(len(wordData.Verbs))]
		}
	}

	sentence := fmt.Sprintf(pattern, words...)
	return strings.Title(sentence)
}

// Generate a paragraph
func generateParagraph(sentences int, style string) string {
	if sentences <= 0 {
		sentences = rand.Intn(5) + 3 // 3-7 sentences
	}

	var sentenceList []string
	for i := 0; i < sentences; i++ {
		sentenceList = append(sentenceList, generateSentence(style))
	}

	return strings.Join(sentenceList, " ")
}

// Count words in text
func countWords(text string) int {
	return len(strings.Fields(text))
}

// CORS middleware
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// Generate paragraph endpoint
func generateHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req ParagraphRequest
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

	// Limit sentences to reasonable range
	if sentences > 20 {
		sentences = 20
	}

	paragraph := generateParagraph(sentences, style)
	wordCount := countWords(paragraph)

	response := ParagraphResponse{
		Paragraph: paragraph,
		WordCount: wordCount,
		Sentences: sentences,
	}

	json.NewEncoder(w).Encode(response)
}

// Health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	status := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"words_loaded": map[string]int{
			"adjectives": len(wordData.Adjectives),
			"adverbs":    len(wordData.Adverbs),
			"nouns":      len(wordData.Nouns),
			"verbs":      len(wordData.Verbs),
		},
	}

	json.NewEncoder(w).Encode(status)
}

// API info endpoint
func infoHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	info := map[string]interface{}{
		"name":        "Paragraph Generator API",
		"version":     "1.0.0",
		"description": "Generate random paragraphs using various word lists",
		"endpoints": map[string]string{
			"/generate": "POST/GET - Generate paragraphs",
			"/health":   "GET - Health check",
			"/info":     "GET - API information",
		},
		"parameters": map[string]string{
			"sentences": "Number of sentences (1-20, default: 5)",
			"style":     "Style of paragraph (general, formal, casual)",
		},
		"examples": map[string]string{
			"GET":  "/generate?sentences=3&style=formal",
			"POST": `{"sentences": 4, "style": "casual"}`,
		},
	}

	json.NewEncoder(w).Encode(info)
}

// Root handler with simple HTML interface
func rootHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Paragraph Generator API</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
        .container { background: #f5f5f5; padding: 20px; border-radius: 8px; }
        button { background: #007cba; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background: #005a87; }
        .output { background: white; padding: 15px; margin: 10px 0; border-radius: 4px; border-left: 4px solid #007cba; }
        input, select { padding: 8px; margin: 5px; border: 1px solid #ddd; border-radius: 4px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🔤 Paragraph Generator API</h1>
        <p>Generate random paragraphs using AI-powered word combinations.</p>
        
        <div>
            <label>Sentences: </label>
            <input type="number" id="sentences" value="5" min="1" max="20">
            
            <label>Style: </label>
            <select id="style">
                <option value="general">General</option>
                <option value="formal">Formal</option>
                <option value="casual">Casual</option>
            </select>
            
            <button onclick="generateParagraph()">Generate Paragraph</button>
        </div>
        
        <div id="output" class="output" style="display:none;">
            <h3>Generated Paragraph:</h3>
            <p id="paragraph"></p>
            <small id="stats"></small>
        </div>
        
        <h3>API Endpoints:</h3>
        <ul>
            <li><strong>GET/POST /generate</strong> - Generate paragraphs</li>
            <li><strong>GET /health</strong> - Health check</li>
            <li><strong>GET /info</strong> - API information</li>
        </ul>
        
        <h3>Example Usage:</h3>
        <pre>
GET /generate?sentences=3&style=formal
POST /generate
{
  "sentences": 4,
  "style": "casual"
}
        </pre>
    </div>

    <script>
        async function generateParagraph() {
            const sentences = document.getElementById('sentences').value;
            const style = document.getElementById('style').value;
            
            try {
                const response = await fetch('/generate', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        sentences: parseInt(sentences),
                        style: style
                    })
                });
                
                const data = await response.json();
                
                document.getElementById('paragraph').textContent = data.paragraph;
				document.getElementById('stats').textContent = 
					'Words: ' + data.word_count + ' | Sentences: ' + data.sentences;
                document.getElementById('output').style.display = 'block';
                
            } catch (error) {
                alert('Error generating paragraph: ' + error.message);
            }
        }
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Load word data
	if err := initWordData(); err != nil {
		log.Fatal("Failed to initialize word data:", err)
	}

	// Setup routes
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/generate", generateHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/info", infoHandler)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Paragraph Generator API starting on port %s", port)
	log.Printf("📖 Loaded %d total words across all categories",
		len(wordData.Adjectives)+len(wordData.Adverbs)+len(wordData.Nouns)+len(wordData.Verbs))
	log.Printf("🌐 Server running at http://localhost:%s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
