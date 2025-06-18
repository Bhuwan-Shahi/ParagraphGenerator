package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Bhuwan-Shahi/ParagraphGenerator/internal/models"
)

type WordLoader struct {
	dataPath string
}

func NewWordLoader(dataPath string) *WordLoader {
	return &WordLoader{dataPath: dataPath}
}

func (l *WordLoader) loadWordsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filepath.Join(l.dataPath, filename))
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

func (l *WordLoader) LoadAll() (*models.WordData, error) {
	wordData := &models.WordData{}

	var err error

	wordData.Adjectives, err = l.loadWordsFromFile("adjectives.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to load adjectives: %v", err)
	}

	wordData.Adverbs, err = l.loadWordsFromFile("adverbs.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to load adverbs: %v", err)
	}

	wordData.Nouns, err = l.loadWordsFromFile("nouns.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to load nouns: %v", err)
	}

	wordData.Verbs, err = l.loadWordsFromFile("verbs.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to load verbs: %v", err)
	}

	return wordData, nil
}
