package main

import (
	"log"

	"github.com/Bhuwan-Shahi/ParagraphGenerator/internal/config"
	"github.com/Bhuwan-Shahi/ParagraphGenerator/internal/handlers"
	"github.com/Bhuwan-Shahi/ParagraphGenerator/internal/services"
	"github.com/Bhuwan-Shahi/ParagraphGenerator/internal/utils"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize word loader
	loader := utils.NewWordLoader(cfg.DataPath)

	// Load word data
	wordData, err := loader.LoadAll()
	if err != nil {
		log.Fatal("Failed to load word data:", err)
	}

	// Initialize generator service
	generator := services.NewGenerator(wordData)

	// Initialize handlers
	handler := handlers.NewHandler(generator)

	// Start server
	log.Printf("🚀 Paragraph Generator API starting on port %s", cfg.Port)
	log.Printf("📖 Loaded %d total words", wordData.TotalCount())
	log.Printf("🌐 Server running at http://localhost:%s", cfg.Port)

	if err := handler.Start(cfg.Port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
