package services

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/Bhuwan-Shahi/ParagraphGenerator/internal/models"
)

type Generator struct {
	wordData *models.WordData
	rand     *rand.Rand
}

func NewGenerator(wordData *models.WordData) *Generator {
	return &Generator{
		wordData: wordData,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// capitalizeFirst capitalizes the first letter of a string
func (g *Generator) capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func (g *Generator) generateSentence(style string) string {
	var patterns []string

	// Different patterns based on style
	switch style {
	case "formal":
		patterns = []string{
			"The %s %s demonstrates %s characteristics.",
			"Furthermore, the %s %s exhibits %s properties.",
			"In accordance with established principles, %s %s %s %s.",
			"The aforementioned %s %s was thoroughly %s and %s.",
			"Subsequently, every %s %s manifests %s qualities.",
			"This particular %s %s systematically %s %s %s.",
			"Numerous %s %s consistently %s %s.",
			"The distinguished %s and %s collectively %s %s %s.",
		}
	case "casual":
		patterns = []string{
			"Hey, the %s %s is totally %s!",
			"You know what? That %s %s is really %s.",
			"So anyway, %s %s just %s %s.",
			"The %s %s was like, super %s and %s.",
			"Honestly, every %s %s kinda %s %s.",
			"This crazy %s %s just %s %s %s.",
			"A bunch of %s %s always %s %s.",
			"The %s and %s totally %s %s %s.",
		}
	default: // general
		patterns = []string{
			"The %s %s %s %s.",
			"A %s %s %s the %s.",
			"%s %s %s %s.",
			"The %s %s was %s and %s.",
			"Every %s %s %s %s.",
			"This %s %s %s %s %s.",
			"Many %s %s %s %s.",
			"The %s and %s %s %s %s.",
		}
	}

	pattern := patterns[g.rand.Intn(len(patterns))]

	// Count placeholders
	placeholders := strings.Count(pattern, "%s")
	words := make([]interface{}, placeholders)

	for i := 0; i < placeholders; i++ {
		words[i] = g.getRandomWord()
	}

	sentence := fmt.Sprintf(pattern, words...)
	return g.capitalizeFirst(sentence)
}

func (g *Generator) getRandomWord() string {
	// Ensure we have words before trying to access them
	totalTypes := 0
	if len(g.wordData.Adjectives) > 0 {
		totalTypes++
	}
	if len(g.wordData.Adverbs) > 0 {
		totalTypes++
	}
	if len(g.wordData.Nouns) > 0 {
		totalTypes++
	}
	if len(g.wordData.Verbs) > 0 {
		totalTypes++
	}

	if totalTypes == 0 {
		return "word" // fallback
	}

	wordType := g.rand.Intn(4)

	switch wordType {
	case 0:
		if len(g.wordData.Adjectives) > 0 {
			return g.wordData.Adjectives[g.rand.Intn(len(g.wordData.Adjectives))]
		}
		fallthrough
	case 1:
		if len(g.wordData.Adverbs) > 0 {
			return g.wordData.Adverbs[g.rand.Intn(len(g.wordData.Adverbs))]
		}
		fallthrough
	case 2:
		if len(g.wordData.Nouns) > 0 {
			return g.wordData.Nouns[g.rand.Intn(len(g.wordData.Nouns))]
		}
		fallthrough
	case 3:
		if len(g.wordData.Verbs) > 0 {
			return g.wordData.Verbs[g.rand.Intn(len(g.wordData.Verbs))]
		}
		fallthrough
	default:
		// Final fallback - try any available word type
		if len(g.wordData.Nouns) > 0 {
			return g.wordData.Nouns[g.rand.Intn(len(g.wordData.Nouns))]
		}
		if len(g.wordData.Adjectives) > 0 {
			return g.wordData.Adjectives[g.rand.Intn(len(g.wordData.Adjectives))]
		}
		if len(g.wordData.Verbs) > 0 {
			return g.wordData.Verbs[g.rand.Intn(len(g.wordData.Verbs))]
		}
		if len(g.wordData.Adverbs) > 0 {
			return g.wordData.Adverbs[g.rand.Intn(len(g.wordData.Adverbs))]
		}
		return "word"
	}
}

func (g *Generator) GenerateParagraph(sentences int, style string) string {
	// Default to random 3-7 sentences if not specified
	if sentences <= 0 {
		sentences = g.rand.Intn(5) + 3 // 3-7 sentences when not specified
	}

	var sentenceList []string
	for i := 0; i < sentences; i++ {
		sentenceList = append(sentenceList, g.generateSentence(style))
	}

	return strings.Join(sentenceList, " ")
}

func (g *Generator) CountWords(text string) int {
	return len(strings.Fields(text))
}

func (g *Generator) GetWordStats() map[string]int {
	return map[string]int{
		"adjectives": len(g.wordData.Adjectives),
		"adverbs":    len(g.wordData.Adverbs),
		"nouns":      len(g.wordData.Nouns),
		"verbs":      len(g.wordData.Verbs),
	}
}

// Additional utility methods
func (g *Generator) GetTotalWords() int {
	return len(g.wordData.Adjectives) + len(g.wordData.Adverbs) +
		len(g.wordData.Nouns) + len(g.wordData.Verbs)
}

func (g *Generator) ValidateWordData() error {
	if g.GetTotalWords() == 0 {
		return fmt.Errorf("no words loaded")
	}
	return nil
}
